package socketservice

import (
	"encoding/json"
	"github.com/getclasslabs/go-chat/internal/domains"
	"github.com/getclasslabs/go-tools/pkg/tracer"
	"github.com/gorilla/websocket"
)

func WriteText(i *tracer.Infos, s *websocket.Conn, message *domains.Message) error {
	return Write(i, s, message, 1)
}


func Write(i *tracer.Infos, s *websocket.Conn, message *domains.Message, msgType int) error {
	i.TraceIt("writing")
	defer i.Span.Finish()

	msgStr, err := json.Marshal(&message)
	if err != nil {
		i.LogError(err)
		return err
	}

	err = write(s, msgType, msgStr)
	if err != nil {
		i.LogError(err)
		return err
	}

	return nil
}

func write(s *websocket.Conn, msgType int , msg []byte) error{
	return s.WriteMessage(msgType, msg)
}

func Read(i *tracer.Infos, s *websocket.Conn) (int, *domains.Message, error) {
	i.TraceIt("reading")
	defer i.Span.Finish()

	msgType, msg, err := s.ReadMessage()
	if err != nil {
		i.LogError(err)
		return 0, nil, err
	}

	var m domains.Message
	err = json.Unmarshal(msg, &m)
	if err != nil{
		i.LogError(err)
		_ = write(s, websocket.TextMessage, []byte("wrong message format, err: " + err.Error()))
	}
	return msgType, &m, err
}
