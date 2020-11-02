package rooms

import (
	"github.com/getclasslabs/go-chat/internal/domains"
	"github.com/getclasslabs/go-chat/internal/repositories/rooms"
	"github.com/getclasslabs/go-tools/pkg/tracer"
)

func Create(i *tracer.Infos, room *domains.Room) error {
	i.TraceIt("creating room")
	defer i.Span.Finish()

	repo := rooms.NewRoom()
	return repo.Create(i, room.Identifier)
}

func IsValid(i *tracer.Infos, identifier string) int64 {
	i.TraceIt("getting room")
	defer i.Span.Finish()

	repo := rooms.NewRoom()
	return repo.Exists(i, identifier)
}