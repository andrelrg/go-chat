package messages

import (
	"errors"
	"github.com/getclasslabs/go-chat/internal/domains"
	"github.com/getclasslabs/go-chat/internal/repositories/messages"
	"github.com/getclasslabs/go-chat/internal/repositories/users"
	"github.com/getclasslabs/go-tools/pkg/tracer"
	"time"
)

var (
	Times = 5
	PauseInSecs = 1
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
	id, err = rep.Create(i, m)
	if err != nil {
		return err
	}

	m.ID = id
	return nil
}

func CircuitSave(i *tracer.Infos, domain *domains.Message) error {
	i.TraceIt("circuit breaker")
	defer i.Span.Finish()
	var err error
	success := false

	for tries := 0; tries <= Times; tries++ {
		err = Save(i, domain)
		if err != nil {
			i.LogError(err)
			time.Sleep(time.Duration(PauseInSecs) * time.Second)
			continue
		}
		success = true
		break
	}
	if !success{
		return err
	}

	return nil
}

func Delete(i *tracer.Infos, domain *domains.Message) error {
	i.TraceIt("deleting message")
	defer i.Span.Finish()

	return errors.New("not implemented yet")
}

func Action(i *tracer.Infos, m *domains.Message) error {
	var err error
	switch m.Kind {
	case domains.Normal:
		err = CircuitSave(i, m)
		if err != nil {
			return err
		}
		break
	case domains.Deleting:
		err = Delete(i, m)
		if err != nil {
			return err
		}
		break
	case domains.Blocking:
		//err = rooms.Block(i, m.RoomID)
		//if err != nil {
		//	return err
		//}
		break
	default:
		return errors.New("kind doesn't exists")
	}
	return nil
}
