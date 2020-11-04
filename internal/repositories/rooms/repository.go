package rooms

import (
	"github.com/getclasslabs/go-chat/internal/repositories"
	"github.com/getclasslabs/go-tools/pkg/db"
	"github.com/getclasslabs/go-tools/pkg/tracer"
)

type Room struct {
	db        db.Database
	traceName string
}

func NewRoom() *Room {
	return &Room{
		db:        repositories.Db,
		traceName: "room repository",
	}
}

func (r *Room) Create(i *tracer.Infos, identifier string) error {
	i.TraceIt(r.traceName)
	defer i.Span.Finish()

	q := "INSERT INTO rooms(identifier) VALUES(?)"

	_, err := r.db.Insert(i, q, identifier)
	if err != nil {
		i.LogError(err)
		return err
	}
	return nil
}

func (r *Room) Exists(i *tracer.Infos, identifier string) int64 {
	i.TraceIt(r.traceName)
	defer i.Span.Finish()

	q := "SELECT id FROM rooms WHERE identifier = ?"

	ret, err := r.db.Get(i, q, identifier)
	if err != nil {
		i.LogError(err)
		return 0
	}

	if len(ret) == 0 {
		return 0
	}

	return ret["id"].(int64)
}
