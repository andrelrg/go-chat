package domains

import (
	"errors"
	"github.com/getclasslabs/go-chat/internal/services/messages"
	"github.com/getclasslabs/go-chat/internal/services/rooms"
	"github.com/getclasslabs/go-tools/pkg/tracer"
)

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

func (m *Message) Action(i *tracer.Infos) error {
	var err error
	switch m.Kind {
	case Normal:
		err = messages.CircuitSave(i, m)
		if err != nil {
			return err
		}
		break
	case Deleting:
		err = messages.Delete(i, m)
		if err != nil {
			return err
		}
		break
	case Blocking:
		err = rooms.Block(i, m.RoomID)
		if err != nil {
			return err
		}
		break
	default:
		return errors.New("kind doesn't exists")
	}
}
