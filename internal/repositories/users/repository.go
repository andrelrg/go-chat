package users

import (
	"github.com/getclasslabs/go-chat/internal/repositories"
	"github.com/getclasslabs/go-tools/pkg/db"
	"github.com/getclasslabs/go-tools/pkg/tracer"
)

type User struct {
	db        db.Database
	traceName string
}

func NewUser() *User {
	return &User{
		db:        repositories.Db,
		traceName: "room repository",
	}
}

func (r *User) Create(i *tracer.Infos, identifier string) (int64, error) {
	i.TraceIt(r.traceName)
	defer i.Span.Finish()

	q := "INSERT INTO users(identifier) VALUES(?)"

	_, err := r.db.Insert(i, q, identifier)
	if err != nil {
		i.LogError(err)
		return 0, err
	}

	return r.Exists(i, identifier), nil
}

func (r *User) Exists(i *tracer.Infos, identifier string) int64 {
	i.TraceIt(r.traceName)
	defer i.Span.Finish()

	q := "SELECT id FROM users WHERE identifier = ?"

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
