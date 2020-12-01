package domains

const (
	Normal   = 0
	Deleting = 1
	Blocking = 2
)

type Message struct {
	ID               int64  `json:"messageID,omitempty"`
	ByID             int64  `json:"byID,omitempty"`
	ByUserIdentifier string `json:"userIdentifier,omitempty"`
	RoomID           int64  `json:"roomID,omitempty"`
	Room             string `json:"roomIdentifier,omitempty"`
	Message          string `json:"message,omitempty"`
	Deleted          bool   `json:"deleted,omitempty"`
	CreatedAt        string `json:"createdAt,omitempty"`
	Kind             int64  `json:"kind,omitempty"`
}
