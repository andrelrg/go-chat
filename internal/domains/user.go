package domains

type User struct {
	ID         int64  `json:"ID,omitempty"`
	Identifier string `json:"userIdentifier,omitempty"`
	FullName   string `json:"fullName,omitempty"`
}
