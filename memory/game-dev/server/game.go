package main

import (
    "fmt"
    "net/http"
    "encoding/json"
)

// 游戏房间
type Room struct {
    ID      string `json:"id"`
    Players []Player `json:"players"`
    Status  string `json:"status"` // waiting, playing, finished
}

// 玩家
type Player struct {
    ID     string `json:"id"`
    Name   string `json:"name"`
    Team   string `json:"team"` // blue, red
    Score  int    `json:"score"`
}

// 创建房间
func CreateRoom(w http.ResponseWriter, r *http.Request) {
    room := Room{
        ID:     fmt.Sprintf("room_%d", now()),
        Players: []Player{},
        Status:  "waiting",
    }
    
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(room)
}

func now() int64 {
    return time.Now().Unix()
}

func main() {
    http.HandleFunc("/api/room/create", CreateRoom)
    http.ListenAndServe(":8080", nil)
}
