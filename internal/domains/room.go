package domains

type Room struct {
	ID         int    `json:"ID,omitempty"`
	Identifier string `json:"roomIdentifier,omitempty"`
}
