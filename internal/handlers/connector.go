package handlers

import (
	"github.com/getclasslabs/go-chat/internal/config"
	"github.com/getclasslabs/go-chat/internal/services/messages"
	"github.com/getclasslabs/go-chat/internal/services/rooms"
	"github.com/getclasslabs/go-chat/internal/services/socketservice"
	"github.com/getclasslabs/go-tools/pkg/request"
	"github.com/getclasslabs/go-tools/pkg/tracer"
	"github.com/gorilla/mux"
	"net/http"
)

func Connect(w http.ResponseWriter, r *http.Request) {
	i := r.Context().Value(request.ContextKey).(*tracer.Infos)
	i.TraceIt(spanName)
	defer i.Span.Finish()

	roomIdentifier := mux.Vars(r)["room"]
	room := rooms.IsValid(i, roomIdentifier)
	if room == 0 {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	socket, err := config.Upgrader.Upgrade(w, r, nil)
	if err != nil {
		i.LogError(err)
		return
	}

	config.Connected[room] = append(config.Connected[room], socket)

	msgs, err := messages.GetMessagesFrom(i, room)
	if err != nil{
		i.LogError(err)
	}
	for _, m := range msgs{
		err = socketservice.WriteText(i, socket, &m)
		if err != nil{
			continue
		}
	}

	for {
		msgType, message, err := socketservice.Read(i, socket)
		if err != nil{
			continue
		}

		message.RoomID = room
		go messages.Save(i, message)

		for _, s := range config.Connected[room] {
			err = socketservice.Write(i, s, message, msgType)
			if err != nil{
				continue
			}
		}
	}
}
