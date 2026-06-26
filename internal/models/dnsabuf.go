package models

import (
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"net"
	"strings"
)

// dnsType maps the most common DNS resource record type codes to their textual
// representation. Unknown types fall back to "TYPEnnn" (RFC 3597 style).
var dnsType = map[uint16]string{
	1:   "A",
	2:   "NS",
	5:   "CNAME",
	6:   "SOA",
	12:  "PTR",
	15:  "MX",
	16:  "TXT",
	28:  "AAAA",
	33:  "SRV",
	35:  "NAPTR",
	43:  "DS",
	46:  "RRSIG",
	47:  "NSEC",
	48:  "DNSKEY",
	257: "CAA",
}

func dnsTypeString(t uint16) string {
	if s, ok := dnsType[t]; ok {
		return s
	}
	return fmt.Sprintf("TYPE%d", t)
}

// parseAbuf decodes the base64-encoded raw DNS message returned by RIPE Atlas in
// the "abuf" field and extracts the answer section. RIPE Atlas does not decode
// the response for us; the answer records only exist inside this wire-format
// message (the count is mirrored in the ANCOUNT header field).
func parseAbuf(abuf string) ([]DNSAnswer, error) {
	msg, err := base64.StdEncoding.DecodeString(abuf)
	if err != nil {
		return nil, fmt.Errorf("decode abuf: %w", err)
	}
	if len(msg) < 12 {
		return nil, fmt.Errorf("dns message too short: %d bytes", len(msg))
	}

	qdCount := binary.BigEndian.Uint16(msg[4:6])
	anCount := binary.BigEndian.Uint16(msg[6:8])

	off := 12
	// Skip the question section: each question is name + qtype(2) + qclass(2).
	for i := 0; i < int(qdCount); i++ {
		_, n, err := readName(msg, off)
		if err != nil {
			return nil, err
		}
		off = n + 4
		if off > len(msg) {
			return nil, fmt.Errorf("truncated question section")
		}
	}

	answers := make([]DNSAnswer, 0, anCount)
	for i := 0; i < int(anCount); i++ {
		name, n, err := readName(msg, off)
		if err != nil {
			return nil, err
		}
		off = n
		if off+10 > len(msg) {
			return nil, fmt.Errorf("truncated answer header")
		}
		rrType := binary.BigEndian.Uint16(msg[off : off+2])
		ttl := binary.BigEndian.Uint32(msg[off+4 : off+8])
		rdLen := int(binary.BigEndian.Uint16(msg[off+8 : off+10]))
		off += 10
		if off+rdLen > len(msg) {
			return nil, fmt.Errorf("truncated rdata")
		}
		rdata := parseRdata(msg, off, rdLen, rrType)
		off += rdLen

		answers = append(answers, DNSAnswer{
			Type:  dnsTypeString(rrType),
			Name:  name,
			Rdata: rdata,
			TTL:   int(ttl),
		})
	}
	return answers, nil
}

// readName decodes a (possibly compressed) DNS name starting at off and returns
// the dotted name plus the offset immediately after the name in the record
// stream (compression pointers do not advance that offset beyond the pointer).
func readName(msg []byte, off int) (string, int, error) {
	var labels []string
	origOff := off
	jumped := false
	// Guard against pointer loops.
	for ptrs := 0; ; {
		if off >= len(msg) {
			return "", 0, fmt.Errorf("name out of bounds")
		}
		b := int(msg[off])
		switch {
		case b == 0:
			off++
			if !jumped {
				origOff = off
			}
			return strings.Join(labels, "."), origOff, nil
		case b&0xc0 == 0xc0:
			if off+1 >= len(msg) {
				return "", 0, fmt.Errorf("truncated compression pointer")
			}
			if !jumped {
				origOff = off + 2
			}
			off = (b&0x3f)<<8 | int(msg[off+1])
			jumped = true
			ptrs++
			if ptrs > len(msg) {
				return "", 0, fmt.Errorf("compression pointer loop")
			}
		default:
			if off+1+b > len(msg) {
				return "", 0, fmt.Errorf("truncated label")
			}
			labels = append(labels, string(msg[off+1:off+1+b]))
			off += 1 + b
		}
	}
}

// parseRdata renders the resource record data into a human-readable string for
// common record types, falling back to hex for anything unrecognized.
func parseRdata(msg []byte, off, rdLen int, rrType uint16) string {
	data := msg[off : off+rdLen]
	switch dnsTypeString(rrType) {
	case "A":
		if rdLen == 4 {
			return net.IP(data).String()
		}
	case "AAAA":
		if rdLen == 16 {
			return net.IP(data).String()
		}
	case "CNAME", "NS", "PTR":
		if name, _, err := readName(msg, off); err == nil {
			return name
		}
	case "MX":
		if rdLen >= 3 {
			pref := binary.BigEndian.Uint16(data[:2])
			if name, _, err := readName(msg, off+2); err == nil {
				return fmt.Sprintf("%d %s", pref, name)
			}
		}
	case "TXT":
		var parts []string
		i := 0
		for i < rdLen {
			l := int(data[i])
			i++
			if i+l > rdLen {
				break
			}
			parts = append(parts, string(data[i:i+l]))
			i += l
		}
		return strings.Join(parts, " ")
	case "SOA":
		mname, n, err := readName(msg, off)
		if err == nil {
			if rname, n2, err := readName(msg, n); err == nil && n2+20 <= off+rdLen {
				serial := binary.BigEndian.Uint32(msg[n2 : n2+4])
				return fmt.Sprintf("%s %s %d", mname, rname, serial)
			}
		}
	}
	return fmt.Sprintf("%x", data)
}
