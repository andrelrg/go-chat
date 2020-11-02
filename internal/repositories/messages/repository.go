package messages

import (
	"github.com/getclasslabs/go-chat/internal/domains"
	"github.com/getclasslabs/go-chat/internal/repositories"
	"github.com/getclasslabs/go-tools/pkg/db"
	"github.com/getclasslabs/go-tools/pkg/tracer"
)

type Message struct {
	db        db.Database
	traceName string
}

func NewMessage() *Message {
	return &Message{
		db:        repositories.Db,
		traceName: "message repository",
	}
}

func (r *Message) Create(i *tracer.Infos, m *domains.Message) error {
	i.TraceIt(r.traceName)
	defer i.Span.Finish()

	q := "INSERT INTO messages(by_user, to_room, message) VALUES(?, ?, ?)"

	_, err := r.db.Insert(i, q, m.ByID, m.RoomID, m.Message)
	if err != nil{
		i.LogError(err)
		return err
	}
	return nil
}

func (r *Message) GetFromRoom(i *tracer.Infos, roomID int64) ([]domains.Message, error) {
	i.TraceIt(r.traceName)
	defer i.Span.Finish()

	q := "SELECT " +
	"	m.id as messageID, " +
	"	u.identifier as userIdentifier, " +
	"	u.full_name as fullName, " +
	"	r.identifier as roomIdentifier, " +
	"	m.message, " +
	"	m.created_at as createdAt " +
	"FROM messages m " +
	"INNER JOIN users u ON u.id = m.by_user " +
	"INNER JOIN rooms r ON r.id = m.to_room " +
	"WHERE " +
	"	m.deleted is false AND " +
	"	r.id = ?"

	ret, err := r.db.Fetch(i, q, roomID)
	if err != nil{
		i.LogError(err)
		return nil, err
	}

	var messages []domains.Message
	err = db.Mapper(i, ret, &messages)
	if err != nil{
		i.LogError(err)
		return nil, err
	}

	return messages, nil
}
