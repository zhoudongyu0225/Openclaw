package main

import (
    "log"
    "net/http"
    "github.com/gorilla/websocket"
    "barrage-game/internal/config"
    "barrage-game/internal/handler"
    "barrage-game/internal/room"
)

var upgrader = websocket.Upgrader{
    CheckOrigin: func(r *http.Request) bool { return true },
}

func main() {
    cfg := config.Load()
    
    // 创建房间管理器
    roomMgr := room.NewManager()
    
    // WebSocket handler
    http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
        handler.HandleWS(w, r, roomMgr)
    })
    
    // API
    http.HandleFunc("/api/room/create", handler.CreateRoom)
    http.HandleFunc("/api/room/join", handler.JoinRoom)
    http.HandleFunc("/api/room/list", handler.ListRooms)
    
    log.Printf("Server starting on :%d", cfg.Port)
    http.ListenAndServe(":"+cfg.Port, nil)
}
