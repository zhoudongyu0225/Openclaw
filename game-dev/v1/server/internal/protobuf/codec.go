package protobuf

import (
    "encoding/binary"
    "errors"
)

var (
    ErrInvalidMessage = errors.New("invalid message")
)

type Message struct {
    ID   uint32
    Data []byte
}

func Encode(msg *Message) ([]byte, error) {
    buf := make([]byte, 6+len(msg.Data))
    binary.BigEndian.PutUint32(buf[0:4], msg.ID)
    binary.BigEndian.PutUint16(buf[4:6], uint16(len(msg.Data)))
    copy(buf[6:], msg.Data)
    return buf, nil
}

func Decode(buf []byte) (*Message, error) {
    if len(buf) < 6 {
        return nil, ErrInvalidMessage
    }
    id := binary.BigEndian.Uint32(buf[0:4])
    length := binary.BigEndian.Uint16(buf[4:6])
    if len(buf) < 6+int(length) {
        return nil, ErrInvalidMessage
    }
    return &Message{
        ID:   id,
        Data: buf[6 : 6+length],
    }, nil
}

// 消息ID定义
const (
    MsgIDCreateRoom    uint32 = 1001
    MsgIDJoinRoom      uint32 = 1002
    MsgIDLeaveRoom     uint32 = 1003
    MsgIDRoomList      uint32 = 1004
    
    MsgIDGameStart     uint32 = 2001
    MsgIDGameState     uint32 = 2002
    MsgIDPlayerAction  uint32 = 2003
    MsgIDGameOver      uint32 = 2004
    
    MsgIDHeartbeat     uint32 = 3001
)
