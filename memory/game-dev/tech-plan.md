# 弹幕游戏技术方案

## 项目架构

### 前端 (UE4)
```
├── 地图系统
│   ├── 网格系统
│   ├── 路径寻路 (A*)
│   └── 视野管理
│
├── 战斗系统
│   ├── 塔防逻辑
│   ├── 兵种AI
│   └── 伤害计算
│
├── 主播系统
│   ├── 视角切换
│   ├── 操作输入
│   └── 指令同步
│
└── 网络同步
    ├── 状态同步
    ├── 预测回滚
    └── 延迟补偿
```

### 后端 (Go)
```
├── 网关层
│   ├── WebSocket
│   ├── 协议解析 (Protobuf)
│   └── 负载均衡
│
├── 业务层
│   ├── 房间管理
│   ├── 匹配系统
│   ├── 战斗逻辑
│   └── 排行榜
│
└── 数据层
    ├── MongoDB (持久化)
    ├── Redis (缓存/状态)
    └── 日志服务
```

---

## 房间系统 (Room System)

### 模块职责
- 房间创建/加入/离开
- 玩家状态管理
- 游戏准备流程
- 房间内广播

### 数据结构

#### 房间 (Room)
```go
type Room struct {
    ID        string           // 房间ID (room_20260307120000)
    Name      string           // 房间名称
    Mode      string           // 游戏模式: rank, normal, practice
    MapID     string           // 地图ID
    Status    string           // waiting, playing, paused, finished
    Players   map[string]*Player
    Settings  RoomSettings
    CreatedAt time.Time
    StartedAt time.Time
    mutex     sync.RWMutex
}

type RoomSettings struct {
    MaxPlayers int    // 最大玩家数 (默认8)
    TimeLimit  int    // 时间限制(秒)
    WinScore   int    // 获胜分数
}
```

#### 玩家 (Player)
```go
type Player struct {
    ID        string
    Name      string
    Team      string  // "blue" or "red"
    IsReady   bool
    IsOwner   bool   // 房主标记
    Conn      *websocket.Conn
    LastActive time.Time
    JoinAt    time.Time
}
```

### 核心流程

#### 1. 创建房间
```
客户端 -> CS_CREATE_ROOM -> 服务器
                               |
                               v
                         创建Room对象
                         分配room_id
                         注册到Manager
                               |
                               v
                         <- SC_CREATE_ROOM_ACK
客户端 <- {room_id, room_info}
```

#### 2. 加入房间
```
客户端 -> CS_JOIN_ROOM(room_id, player_name)
                               |
                               v
                         验证房间存在
                         验证房间未开始
                         分配队伍(红/蓝)
                         添加到房间Players
                               |
                               v
                         <- SC_JOIN_ROOM_ACK
客户端 <- {room_info, self_info}

                         广播 SC_PLAYER_JOIN 给房间内其他玩家
```

#### 3. 准备/开始游戏
```
客户端 -> CS_READY(true/false)
                               |
                               v
                         更新玩家准备状态
                         广播 SC_ROOM_UPDATE
                               |
                               v
房主 -> CS_START_GAME
              |
              v
        验证所有玩家准备
        验证房主身份
        状态 -> playing
              |
              v
        广播 SC_GAME_START
        通知战斗模块开始
```

#### 4. 离开房间
```
客户端 -> CS_LEAVE_ROOM
              |
              v
        从房间移除玩家
        删除player->room映射
              |
              +-> 房间为空? 删除房间
              +-> 游戏进行中? 处理断线
              |
              v
        广播 SC_PLAYER_LEAVE
```

### Manager 核心方法

```go
type Manager struct {
    rooms     map[string]*Room      // roomID -> Room
    playerMap map[string]string     // playerID -> roomID
    roomIndex map[string][]string   // mode -> roomIDs (索引)
}

func NewManager(heartbeat time.Duration) *Manager

func (m *Manager) CreateRoom(req *CreateRoomReq) (*Room, error)
func (m *Manager) JoinRoom(roomID, playerID, playerName string) error
func (m *Manager) LeaveRoom(playerID string) error
func (m *Manager) GetRoom(roomID string) *Room
func (m *Manager) ListRooms() []*Room
func (m *Manager) ListRoomsByMode(mode string) []*Room
func (m *Manager) Broadcast(roomID string, msg *WSMessage) error
func (m *Manager) GetPlayerRoom(playerID string) string
```

### 心跳与超时

| 参数 | 值 | 说明 |
|------|-----|------|
| pingPeriod | 30s | 发送ping间隔 |
| pongWait | 60s | 等待pong超时 |
| writeWait | 10s | 写操作超时 |
| maxMsgSize | 512KB | 单条消息最大 |
| roomCleanup | 30min | 无人房间清理 |

### WebSocket 连接管理

```
                    readPump()              writePump()
    Client <----------------> WS <----------------> Server
         (读取消息)            (ping/pong)          (发送消息)
                              30s周期
```

### Protobuf 消息定义

完整协议定义见 `server/proto/game.proto`

#### 客户端消息 (CS)
- `CS_HEARTBEAT` - 心跳
- `CS_CREATE_ROOM` - 创建房间
- `CS_JOIN_ROOM` - 加入房间
- `CS_LEAVE_ROOM` - 离开房间
- `CS_ROOM_LIST` - 房间列表
- `CS_READY` - 准备/取消准备
- `CS_START_GAME` - 开始游戏

#### 服务器消息 (SC)
- `SC_HEARTBEAT_ACK` - 心跳回应
- `SC_CREATE_ROOM_ACK` - 创建房间回应
- `SC_JOIN_ROOM_ACK` - 加入房间回应
- `SC_LEAVE_ROOM_ACK` - 离开房间回应
- `SC_ROOM_LIST_ACK` - 房间列表回应
- `SC_ROOM_UPDATE` - 房间更新
- `SC_PLAYER_JOIN` - 玩家加入
- `SC_PLAYER_LEAVE` - 玩家离开
- `SC_GAME_START` - 游戏开始
- `SC_ERROR` - 错误消息

---

## 通信协议

### 消息格式 (Protobuf)
```protobuf
message CSMessage {
    uint32 msg_id = 1;
    bytes data = 2;
}

message SCMessage {
    uint32 msg_id = 1;
    bytes data = 2;
}
```

### 核心消息
- `RoomCreate` / `RoomJoin` - 房间
- `MatchStart` - 匹配开始
- `GameStart` - 游戏开始
- `SyncState` - 状态同步
- `PlayerAction` - 玩家操作
- `GameOver` - 游戏结束

## 平台对接

### 抖音
- 小程序SDK
- 登录授权
- 支付
- 分享

### 快手
- 小游戏SDK
- 登录授权
- 支付
- 分享

## 开发计划

### Phase 1 (1-2周)
- [x] 技术选型
- [x] 房间系统 (WebSocket + Manager)
- [x] Protobuf 协议定义
- [ ] 搭建项目骨架
- [ ] 实现基础通信

### Phase 2 (2-3周)
- [x] 房间系统
- [ ] 简单战斗逻辑
- [ ] 基础UI

### Phase 3 (3-4周)
- [ ] 完整塔防逻辑
- [ ] 主播视角切换
- [ ] 积分系统

### Phase 4 (4-5周)
- [ ] 对接抖音SDK
- [ ] 对接快手SDK
- [ ] 测试调优
