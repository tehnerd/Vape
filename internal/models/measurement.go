package models

import "time"

type MeasurementType string

const (
	TypePing       MeasurementType = "ping"
	TypeTraceroute MeasurementType = "traceroute"
	TypeDNS        MeasurementType = "dns"
	TypeSSL        MeasurementType = "sslcert"
	TypeNTP        MeasurementType = "ntp"
	TypeHTTP       MeasurementType = "http"
)

type MeasurementStatusName string

const (
	StatusSpecified MeasurementStatusName = "Specified"
	StatusScheduled MeasurementStatusName = "Scheduled"
	StatusOngoing   MeasurementStatusName = "Ongoing"
	StatusStopped   MeasurementStatusName = "Stopped"
)

type MeasurementStatus struct {
	ID   int                   `json:"id"`
	Name MeasurementStatusName `json:"name"`
}

type Measurement struct {
	ID                int               `json:"id"`
	Type              MeasurementType   `json:"type"`
	Status            MeasurementStatus `json:"status"`
	Target            string            `json:"target"`
	Description       string            `json:"description"`
	AF                int               `json:"af"`
	IsOneoff          bool              `json:"is_oneoff"`
	ResolvedIPs       []string          `json:"resolved_ips"`
	ParticipantCount  int               `json:"participant_count"`
	ProbesRequested   int               `json:"probes_requested"`
	ProbesScheduled   int               `json:"probes_scheduled"`
	CreationTime      int64             `json:"creation_time"`
	StartTime         int64             `json:"start_time"`
	StopTime          int64             `json:"stop_time"`
	Interval          int               `json:"interval"`
	Protocol          string            `json:"protocol"`
	ResolveOnProbe    bool              `json:"resolve_on_probe"`
	CreditsPerResult  int               `json:"credits_per_result"`
	EstimatedResultsPerDay int          `json:"estimated_results_per_day"`
}

func (m *Measurement) GetCreationTime() time.Time {
	return time.Unix(m.CreationTime, 0)
}

func (m *Measurement) GetStartTime() time.Time {
	return time.Unix(m.StartTime, 0)
}

func (m *Measurement) GetStopTime() time.Time {
	if m.StopTime == 0 {
		return time.Time{}
	}
	return time.Unix(m.StopTime, 0)
}

type MeasurementListResponse struct {
	Count    int           `json:"count"`
	Next     string        `json:"next"`
	Previous string        `json:"previous"`
	Results  []Measurement `json:"results"`
}

type MeasurementDefinition struct {
	Type           MeasurementType `json:"type"`
	Target         string          `json:"target,omitempty"`
	AF             int             `json:"af"`
	Description    string          `json:"description,omitempty"`
	Protocol       string          `json:"protocol,omitempty"`
	ResolveOnProbe bool            `json:"resolve_on_probe,omitempty"`
	Packets        int             `json:"packets,omitempty"`
	Size           int             `json:"size,omitempty"`
	// DNS specific
	QueryClass      string `json:"query_class,omitempty"`
	QueryType       string `json:"query_type,omitempty"`
	QueryArgument   string `json:"query_argument,omitempty"`
	UseProbeResolver bool  `json:"use_probe_resolver,omitempty"`
	// HTTP specific
	Method          string `json:"method,omitempty"`
	Path            string `json:"path,omitempty"`
	HeaderBytes     int    `json:"header_bytes,omitempty"`
	// SSL specific
	Port            int    `json:"port,omitempty"`
}

type ProbeSpec struct {
	Requested int    `json:"requested"`
	Type      string `json:"type"`
	Value     string `json:"value"`
	Tags      *Tags  `json:"tags,omitempty"`
}

type Tags struct {
	Include []string `json:"include,omitempty"`
	Exclude []string `json:"exclude,omitempty"`
}

type MeasurementRequest struct {
	Definitions []MeasurementDefinition `json:"definitions"`
	Probes      []ProbeSpec             `json:"probes"`
	IsOneoff    bool                    `json:"is_oneoff"`
	BillTo      string                  `json:"bill_to,omitempty"`
	StartTime   int64                   `json:"start_time,omitempty"`
	StopTime    int64                   `json:"stop_time,omitempty"`
}

type MeasurementCreateResponse struct {
	Measurements []int `json:"measurements"`
}
