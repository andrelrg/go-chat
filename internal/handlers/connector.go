package handlers

import (
	"github.com/getclasslabs/go-chat/internal/services/messages"
	"github.com/getclasslabs/go-chat/internal/services/rooms"
	"github.com/getclasslabs/go-chat/internal/services/socketservice"
	"github.com/getclasslabs/go-tools/pkg/request"
	"github.com/getclasslabs/go-tools/pkg/tracer"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"net/http"
)

func Connect(w http.ResponseWriter, r *http.Request) {
	i := r.Context().Value(request.ContextKey).(*tracer.Infos)
	i.TraceIt(spanName)
	defer i.Span.Finish()

	socket, err := socketservice.NewSocket(w, r)
	if err != nil {
		i.LogError(err)
		return
	}

	roomIdentifier := mux.Vars(r)["room"]

	room, err := rooms.Enter(i, socket, roomIdentifier)
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	for {
		msgType, message, err := socket.Read(i)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err) {
				socket.Remove(room)
				return
			}
			continue
		}

		if msgType == websocket.CloseMessage {
			socket.Remove(room)
			return
		}

		message.RoomID = room
		messages.Action(i, message)

		for _, s := range socketservice.Connected[room] {
			err = s.Write(i, message, msgType)
			if err != nil {
				continue
			}
		}
	}
}
