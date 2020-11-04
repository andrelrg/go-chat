package socketservice

import (
	"encoding/json"
	"github.com/getclasslabs/go-chat/internal/domains"
	"github.com/getclasslabs/go-tools/pkg/tracer"
	"github.com/gorilla/websocket"
	"github.com/lithammer/shortuuid"
	"net/http"
	"strconv"
	"time"
)

type SocketUtils struct {
	S  *websocket.Conn
	ID string
}

func NewSocket(w http.ResponseWriter, r *http.Request) (*SocketUtils, error) {
	s, err := Upgrader.Upgrade(w, r, nil)
	if err != nil {
		return nil, err
	}
	return &SocketUtils{
		s,
		shortuuid.New(),
	}, nil
}

func (s *SocketUtils) WriteText(i *tracer.Infos, message *domains.Message) error {
	return s.Write(i, message, 1)
}

func (s *SocketUtils) Write(i *tracer.Infos, message *domains.Message, msgType int) error {
	i.TraceIt("writing")
	defer i.Span.Finish()

	msgStr, err := json.Marshal(&message)
	if err != nil {
		i.LogError(err)
		return err
	}

	err = s.write(msgType, msgStr)
	if err != nil {
		i.LogError(err)
		return err
	}

	return nil
}

func (s *SocketUtils) write(msgType int, msg []byte) error {
	return s.S.WriteMessage(msgType, msg)
}

func (s *SocketUtils) Read(i *tracer.Infos) (int, *domains.Message, error) {
	i.TraceIt("reading")
	defer i.Span.Finish()

	msgType, msg, err := s.S.ReadMessage()
	if err != nil {
		i.LogError(err)
		return 0, nil, err
	}

	var m domains.Message
	err = json.Unmarshal(msg, &m)
	if err != nil {
		i.LogError(err)
		_ = s.write(websocket.TextMessage, []byte("wrong message format, err: "+err.Error()))
	}

	m.CreatedAt = strconv.FormatInt(time.Now().Unix(), 10)
	return msgType, &m, err
}

func (s *SocketUtils) Remove(room int64) {
	l := Connected[room]
	for i, socket := range l {
		if socket.ID == s.ID {
			Connected[room] = append(l[:i], l[i+1:]...)
		}
	}
}
