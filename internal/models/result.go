package models

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
	Results  []PingResult `json:"result,omitempty"`
	// Traceroute specific
	ParisID int                   `json:"paris_id,omitempty"`
	Hops    []TracerouteHop       `json:"result,omitempty"`
	// DNS specific
	DNSResult *DNSResult           `json:"resultset,omitempty"`
	// SSL specific
	SSLResult *SSLResult           `json:"cert,omitempty"`
	// HTTP specific
	HTTPResult *HTTPResult         `json:"result,omitempty"`
	// NTP specific
	NTPResult *NTPResult          `json:"result,omitempty"`
}

type ResultError struct {
	Timeout int    `json:"timeout,omitempty"`
	Message string `json:"message,omitempty"`
}

type PingResult struct {
	RTT     float64 `json:"rtt,omitempty"`
	TTL     int     `json:"ttl,omitempty"`
	Error   string  `json:"error,omitempty"`
	Timeout bool    `json:"x,omitempty"`
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
	RTT     float64 `json:"rtt,omitempty"`
	Error   string  `json:"err,omitempty"`
	Timeout bool    `json:"x,omitempty"`
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

type SSLResult struct {
	SubjectCN    string   `json:"subject_cn,omitempty"`
	IssuerCN     string   `json:"issuer_cn,omitempty"`
	NotBefore    string   `json:"not_before,omitempty"`
	NotAfter     string   `json:"not_after,omitempty"`
	SHA256       string   `json:"sha256fp,omitempty"`
	SHA1         string   `json:"sha1fp,omitempty"`
	CertChain    []string `json:"chain,omitempty"`
	ServerCipher string   `json:"server_cipher,omitempty"`
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

type ResultListResponse struct {
	Results []MeasurementResult `json:"results,omitempty"`
}
