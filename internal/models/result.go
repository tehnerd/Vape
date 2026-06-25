package models

import "encoding/json"

type MeasurementResult struct {
	MsmID     int     `json:"msm_id"`
	ProbeID   int     `json:"prb_id"`
	Timestamp int64   `json:"timestamp"`
	Type      string  `json:"type"`
	From      string  `json:"from"`
	DstAddr   string  `json:"dst_addr"`
	DstName   string  `json:"dst_name"`
	AF        int     `json:"af"`
	Proto     string  `json:"proto"`
	// Common result fields
	Error   *ResultError `json:"error,omitempty"`
	Stored  int64        `json:"stored"`
	GroupID int          `json:"group_id"`
	// Ping specific
	Avg      float64      `json:"avg,omitempty"`
	Min      float64      `json:"min,omitempty"`
	Max      float64      `json:"max,omitempty"`
	Rcvd     int          `json:"rcvd,omitempty"`
	Sent     int          `json:"sent,omitempty"`
	Dup      int          `json:"dup,omitempty"`
	TTL      int          `json:"ttl,omitempty"`
	Size     int          `json:"size,omitempty"`
	Results  []PingResult `json:"-"`
	// Traceroute specific
	ParisID int             `json:"paris_id,omitempty"`
	Hops    []TracerouteHop `json:"-"`
	// DNS specific
	DNSResult *DNSResult `json:"resultset,omitempty"`
	// SSL specific
	SSLCerts     []SSLCert `json:"cert,omitempty"`
	ServerCipher string    `json:"server_cipher,omitempty"`
	RT           float64   `json:"rt,omitempty"`
	TLSVersion   string    `json:"ver,omitempty"`
	TTC          float64   `json:"ttc,omitempty"`
	// HTTP specific
	HTTPResult *HTTPResult `json:"-"`
	// NTP specific
	NTPResult *NTPResult `json:"-"`
}

type ResultError struct {
	Timeout int    `json:"timeout,omitempty"`
	Message string `json:"message,omitempty"`
}

// TimeoutFlag represents the RIPE Atlas "x" field, which marks a lost packet or
// hop timeout. The API encodes it as the string "*", though some result types
// may use a JSON boolean, so both forms are accepted.
type TimeoutFlag bool

func (t *TimeoutFlag) UnmarshalJSON(data []byte) error {
	var b bool
	if err := json.Unmarshal(data, &b); err == nil {
		*t = TimeoutFlag(b)
		return nil
	}
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	*t = TimeoutFlag(s != "")
	return nil
}

type PingResult struct {
	RTT     float64     `json:"rtt,omitempty"`
	TTL     int         `json:"ttl,omitempty"`
	Error   string      `json:"error,omitempty"`
	Timeout TimeoutFlag `json:"x,omitempty"`
}

type TracerouteHop struct {
	Hop    int             `json:"hop"`
	Error  string          `json:"error,omitempty"`
	Result []TracerouteRTT `json:"result,omitempty"`
}

type TracerouteRTT struct {
	From    string  `json:"from,omitempty"`
	TTL     int     `json:"ttl,omitempty"`
	Size    int     `json:"size,omitempty"`
	RTT     float64     `json:"rtt,omitempty"`
	Error   string      `json:"err,omitempty"`
	Timeout TimeoutFlag `json:"x,omitempty"`
}

type DNSResult struct {
	Time       float64    `json:"rt,omitempty"`
	Size       int        `json:"size,omitempty"`
	Abuf       string     `json:"abuf,omitempty"`
	Answers    []DNSAnswer `json:"answers,omitempty"`
	NSID       string     `json:"nsid,omitempty"`
	SubID      int        `json:"subid,omitempty"`
	SubMax     int        `json:"submax,omitempty"`
	DstAddr    string     `json:"dst_addr,omitempty"`
	Error      *DNSError  `json:"error,omitempty"`
}

type DNSAnswer struct {
	Type   string `json:"type"`
	Name   string `json:"name"`
	Rdata  string `json:"rdata"`
	TTL    int    `json:"ttl"`
}

type DNSError struct {
	Timeout int    `json:"timeout,omitempty"`
	Message string `json:"message,omitempty"`
}

// SSLCert represents a single certificate in the chain.
// The RIPE Atlas API returns cert as an array of PEM strings.
type SSLCert struct {
	Raw string
}

func (s *SSLCert) UnmarshalJSON(data []byte) error {
	var raw string
	if err := json.Unmarshal(data, &raw); err == nil {
		s.Raw = raw
		return nil
	}
	// Fallback: try as object with known fields
	var obj struct {
		SubjectCN string `json:"subject_cn"`
		IssuerCN  string `json:"issuer_cn"`
	}
	if err := json.Unmarshal(data, &obj); err != nil {
		return err
	}
	s.Raw = obj.SubjectCN
	return nil
}

type HTTPResult struct {
	HeaderSize int     `json:"hsize,omitempty"`
	BodySize   int     `json:"bsize,omitempty"`
	Header     string  `json:"header,omitempty"`
	Body       string  `json:"body,omitempty"`
	RT         float64 `json:"rt,omitempty"`
	ResCode    int     `json:"res,omitempty"`
	Method     string  `json:"method,omitempty"`
	Version    string  `json:"ver,omitempty"`
	DnsError   string  `json:"dnserr,omitempty"`
	Error      string  `json:"err,omitempty"`
}

type NTPResult struct {
	FinalTimestamp float64 `json:"final-ts,omitempty"`
	Offset         float64 `json:"offset,omitempty"`
	Poll           int     `json:"poll,omitempty"`
	Precision      int     `json:"precision,omitempty"`
	RefID          string  `json:"ref-id,omitempty"`
	RefTS          float64 `json:"ref-ts,omitempty"`
	RootDelay      float64 `json:"root-delay,omitempty"`
	Stratum        int     `json:"stratum,omitempty"`
	Version        int     `json:"version,omitempty"`
	Li             string  `json:"li,omitempty"`
	Mode           string  `json:"mode,omitempty"`
	Error          string  `json:"error,omitempty"`
}

// UnmarshalJSON handles the polymorphic "result" field which maps to different
// types depending on the measurement type (ping, traceroute, http, ntp).
func (m *MeasurementResult) UnmarshalJSON(data []byte) error {
	type Alias MeasurementResult
	aux := &struct {
		Result json.RawMessage `json:"result,omitempty"`
		*Alias
	}{
		Alias: (*Alias)(m),
	}
	if err := json.Unmarshal(data, aux); err != nil {
		return err
	}

	if len(aux.Result) == 0 {
		return nil
	}

	switch m.Type {
	case "traceroute":
		return json.Unmarshal(aux.Result, &m.Hops)
	case "ping":
		return json.Unmarshal(aux.Result, &m.Results)
	case "http":
		m.HTTPResult = &HTTPResult{}
		return json.Unmarshal(aux.Result, m.HTTPResult)
	case "ntp":
		m.NTPResult = &NTPResult{}
		return json.Unmarshal(aux.Result, m.NTPResult)
	}
	return nil
}

type ResultListResponse struct {
	Results []MeasurementResult `json:"results,omitempty"`
}
