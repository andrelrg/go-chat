package messages

import (
	"github.com/getclasslabs/go-chat/internal/domains"
	"github.com/getclasslabs/go-chat/internal/repositories/messages"
	"github.com/getclasslabs/go-chat/internal/repositories/users"
	"github.com/getclasslabs/go-tools/pkg/tracer"
)


func GetMessagesFrom(i *tracer.Infos, room int64) ([]domains.Message, error){
	i.TraceIt("getting messages")
	defer i.Span.Finish()

	rep := messages.NewMessage()
	return rep.GetFromRoom(i, room)
}

func Save(i *tracer.Infos, m *domains.Message) error{
	i.TraceIt("saving messages")
	defer i.Span.Finish()
	var err error

	uRepo := users.NewUser()
	id := uRepo.Exists(i, m.ByUserIdentifier)
	if id == 0 {
		id, err = uRepo.Create(i, m.ByUserIdentifier)
		if err != nil {
			return err
		}
	}
	m.ByID = id
	rep := messages.NewMessage()
	return rep.Create(i, m)
}