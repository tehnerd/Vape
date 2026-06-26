package output

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"net"
	"strings"
	"testing"
	"time"
)

func makeTestPEM(t *testing.T) string {
	t.Helper()
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatal(err)
	}
	tmpl := &x509.Certificate{
		SerialNumber: big.NewInt(12345),
		Subject:      pkix.Name{CommonName: "example.com", Organization: []string{"Example Inc"}},
		Issuer:       pkix.Name{CommonName: "Example Root CA"},
		NotBefore:    time.Now(),
		NotAfter:     time.Now().Add(24 * time.Hour),
		DNSNames:     []string{"example.com", "www.example.com"},
		IPAddresses:  []net.IP{net.ParseIP("192.0.2.1")},
		KeyUsage:     x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
	}
	der, err := x509.CreateCertificate(rand.Reader, tmpl, tmpl, &key.PublicKey, key)
	if err != nil {
		t.Fatal(err)
	}
	return string(pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: der}))
}

func TestFormatSSLCertNames(t *testing.T) {
	var buf bytes.Buffer
	formatSSLCert(&buf, 0, makeTestPEM(t), SSLCertNames)
	out := buf.String()
	t.Log("\n" + out)
	// Self-signed test cert: issuer == subject.
	for _, want := range []string{"Subject CN: example.com", "Issuer CN:  example.com", "DNS:www.example.com", "IP:192.0.2.1"} {
		if !strings.Contains(out, want) {
			t.Errorf("names output missing %q", want)
		}
	}
}

func TestFormatSSLCertFull(t *testing.T) {
	var buf bytes.Buffer
	formatSSLCert(&buf, 0, makeTestPEM(t), SSLCertFull)
	out := buf.String()
	t.Log("\n" + out)
	for _, want := range []string{"-----BEGIN CERTIFICATE-----", "Decoded:", "RSA 2048 bits", "Server Auth", "Key Encipherment", "SANs:"} {
		if !strings.Contains(out, want) {
			t.Errorf("full output missing %q", want)
		}
	}
}

func TestParseSSLCertModeAndEscapedNewlines(t *testing.T) {
	if ParseSSLCertMode("full") != SSLCertFull || ParseSSLCertMode("names") != SSLCertNames || ParseSSLCertMode("") != SSLCertNames {
		t.Fatal("ParseSSLCertMode mapping wrong")
	}
	// Atlas-style literal \n escapes should still parse.
	escaped := strings.ReplaceAll(makeTestPEM(t), "\n", "\\n")
	var buf bytes.Buffer
	formatSSLCert(&buf, 0, escaped, SSLCertNames)
	if !strings.Contains(buf.String(), "Subject CN: example.com") {
		t.Errorf("escaped-newline cert did not parse: %s", buf.String())
	}
}
