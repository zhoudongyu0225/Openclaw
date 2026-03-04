package handler

import (
    "encoding/json"
    "net/http"
    "barrage-game/internal/room"
)

type RoomManager interface {
    CreateRoom() *room.Room
    GetRoom(id string) *room.Room
    DeleteRoom(id string)
}

func HandleWS(w http.ResponseWriter, r *http.Request, mgr RoomManager) {
    // WebSocket升级逻辑
    // 房间加入/退出处理
    // 消息收发
}

func CreateRoom(w http.ResponseWriter, r *http.Request, mgr RoomManager) {
    rRoom := mgr.CreateRoom()
    json.NewEncoder(w).Encode(map[string]string{
        "room_id": rRoom.ID,
    })
}

func JoinRoom(w http.ResponseWriter, r *http.Request, mgr RoomManager) {
    // 加入房间逻辑
}

func ListRooms(w http.ResponseWriter, r *http.Request, mgr RoomManager) {
    // 列出房间逻辑
}
