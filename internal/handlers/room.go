package handlers

import (
	"encoding/json"
	"github.com/getclasslabs/go-chat/internal/domains"
	"github.com/getclasslabs/go-chat/internal/services/rooms"
	"github.com/getclasslabs/go-tools/pkg/request"
	"github.com/getclasslabs/go-tools/pkg/tracer"
	"net/http"
)

func CreateRoom(w http.ResponseWriter, r *http.Request) {
	i := r.Context().Value(request.ContextKey).(*tracer.Infos)
	i.TraceIt(spanName)
	defer i.Span.Finish()

	room := domains.Room{}
	err := json.NewDecoder(r.Body).Decode(&room)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = rooms.Create(i, &room)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}
