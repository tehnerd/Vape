package models

type ProbeStatusName string

const (
	ProbeStatusConnected    ProbeStatusName = "Connected"
	ProbeStatusDisconnected ProbeStatusName = "Disconnected"
	ProbeStatusAbandoned    ProbeStatusName = "Abandoned"
	ProbeStatusNeverSeen    ProbeStatusName = "Never Seen"
)

type ProbeStatus struct {
	ID    int             `json:"id"`
	Name  ProbeStatusName `json:"name"`
	Since string          `json:"since,omitempty"`
}

type Probe struct {
	ID             int          `json:"id"`
	AddressV4      string       `json:"address_v4"`
	AddressV6      string       `json:"address_v6"`
	ASNv4          int          `json:"asn_v4"`
	ASNv6          int          `json:"asn_v6"`
	CountryCode    string       `json:"country_code"`
	Description    string       `json:"description"`
	FirstConnected int64        `json:"first_connected"`
	LastConnected  int64        `json:"last_connected"`
	Latitude       float64      `json:"latitude"`
	Longitude      float64      `json:"longitude"`
	PrefixV4       string       `json:"prefix_v4"`
	PrefixV6       string       `json:"prefix_v6"`
	Status         ProbeStatus  `json:"status"`
	StatusSince    int64        `json:"status_since"`
	Tags           []ProbeTag   `json:"tags"`
	IsAnchor       bool         `json:"is_anchor"`
	IsPublic       bool         `json:"is_public"`
	Geometry       *Geometry    `json:"geometry,omitempty"`
}

type ProbeTag struct {
	Name string `json:"name"`
	Slug string `json:"slug"`
}

type Geometry struct {
	Type        string    `json:"type"`
	Coordinates []float64 `json:"coordinates"`
}

type ProbeListResponse struct {
	Count    int     `json:"count"`
	Next     string  `json:"next"`
	Previous string  `json:"previous"`
	Results  []Probe `json:"results"`
}

type ProbeListParams struct {
	CountryCode string
	ASN         int
	Status      ProbeStatusName
	IsAnchor    *bool
	IsPublic    *bool
	Limit       int
	Offset      int
}
