package rooms

import (
	"errors"
	"github.com/getclasslabs/go-chat/internal/repositories/rooms"
	"github.com/getclasslabs/go-chat/internal/services/messages"
	"github.com/getclasslabs/go-chat/internal/services/socketservice"
	"github.com/getclasslabs/go-tools/pkg/tracer"
)

func create(i *tracer.Infos, identifier string) error {
	i.TraceIt("creating room")
	defer i.Span.Finish()

	repo := rooms.NewRoom()
	return repo.Create(i, identifier)
}

func isValid(i *tracer.Infos, identifier string) int64 {
	i.TraceIt("getting room")
	defer i.Span.Finish()

	repo := rooms.NewRoom()
	return repo.Exists(i, identifier)
}

func Enter(i *tracer.Infos, socket *socketservice.SocketUtils, roomIdentifier string) (int64, error) {
	room := isValid(i, roomIdentifier)
	if room == 0 {
		err := create(i, roomIdentifier)
		if err != nil {

			return 0, err
		}
		room = isValid(i, roomIdentifier)
	}

	socketservice.Connected[room] = append(socketservice.Connected[room], socket)

	msgs, err := messages.GetMessagesFrom(i, room)
	if err != nil{
		i.LogError(err)
	}
	for _, m := range msgs{
		err = socket.WriteText(i, &m)
		if err != nil{
			continue
		}
	}
	return room, nil
}

func Block(i *tracer.Infos, roomID int64) error {
	i.TraceIt("deleting message")
	defer i.Span.Finish()

	return errors.New("not implemented yet")
}