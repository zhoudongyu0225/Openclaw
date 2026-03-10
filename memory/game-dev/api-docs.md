# 弹幕游戏 - API 接口文档

## 基础信息

| 项目 | 值 |
|------|-----|
| 基础URL | `https://api.bulletgame.com` |
| WebSocket | `wss://ws.bulletgame.com` |
| 协议 | HTTP/1.1 + WebSocket |
| 数据格式 | JSON / Protobuf |
| 字符编码 | UTF-8 |

---

## 认证

### 抖音授权登录

```
GET /api/auth/douyin/login
```

**请求参数**
| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| redirect_uri | string | 是 | 授权回调地址 |

**响应**
```json
{
  "code": 0,
  "data": {
    "auth_url": "https://open.douyin.com/oauth/authorize/..."
  }
}
```

### 授权回调

```
GET /api/auth/douyin/callback
```

**请求参数**
| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| code | string | 是 | 授权码 |
| state | string | 是 | 状态令牌 |

**响应**
```json
{
  "code": 0,
  "data": {
    "access_token": "xxxxx",
    "refresh_token": "xxxxx",
    "expires_in": 86400,
    "open_id": "xxxxx"
  }
}
```

### Token 刷新

```
POST /api/auth/douyin/refresh
```

**请求头**
```
Authorization: Bearer <access_token>
```

**响应**
```json
{
  "code": 0,
  "data": {
    "access_token": "xxxxx",
    "expires_in": 86400
  }
}
```

---

## 用户接口

### 获取用户信息

```
GET /api/user/info
```

**请求头**
```
Authorization: Bearer <access_token>
```

**响应**
```json
{
  "code": 0,
  "data": {
    "user_id": "user_123",
    "nickname": "主播A",
    "avatar": "https://example.com/avatar.png",
    "level": 10,
    "exp": 5000,
    "coins": 1000,
    "diamonds": 50,
    "statistics": {
      "total_games": 100,
      "win_games": 70,
      "total_kills": 500,
      "total_damage": 1000000
    }
  }
}
```

### 更新用户设置

```
PUT /api/user/settings
```

**请求体**
```json
{
  "nickname": "新名字",
  "avatar": "https://example.com/new_avatar.png",
  "settings": {
    "music_volume": 80,
    "sfx_volume": 90,
    "notification_enabled": true
  }
}
```

---

## 房间接口

### 创建房间

```
POST /api/room/create
```

**请求头**
```
Authorization: Bearer <access_token>
```

**请求体**
```json
{
  "room_name": "我的房间",
  "game_mode": "rank",
  "map_id": "map_01",
  "max_players": 8,
  "password": "optional_password"
}
```

**响应**
```json
{
  "code": 0,
  "data": {
    "room_id": "room_20260307120000",
    "room_name": "我的房间",
    "game_mode": "rank",
    "host_id": "user_123",
    "status": "waiting",
    "players": [
      {
        "user_id": "user_123",
        "nickname": "主播A",
        "is_ready": true,
        "is_owner": true
      }
    ],
    "websocket_url": "wss://ws.bulletgame.com/room/room_20260307120000"
  }
}
```

### 加入房间

```
POST /api/room/join
```

**请求体**
```json
{
  "room_id": "room_20260307120000",
  "password": "optional_password"
}
```

**响应**
```json
{
  "code": 0,
  "data": {
    "room_id": "room_20260307120000",
    "room_name": "我的房间",
    "players": [...],
    "websocket_url": "wss://ws.bulletgame.com/room/room_20260307120000"
  }
}
```

### 离开房间

```
POST /api/room/leave
```

**请求体**
```json
{
  "room_id": "room_20260307120000"
}
```

### 房间列表

```
GET /api/room/list
```

**查询参数**
| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| game_mode | string | 否 | 筛选模式 (rank/normal/practice) |
| page | int | 否 | 页码 (默认1) |
| page_size | int | 否 | 每页数量 (默认20) |

**响应**
```json
{
  "code": 0,
  "data": {
    "rooms": [
      {
        "room_id": "room_xxx",
        "room_name": "房间名",
        "game_mode": "rank",
        "host_id": "user_123",
        "status": "waiting",
        "current_players": 4,
        "max_players": 8
      }
    ],
    "total": 100,
    "page": 1,
    "page_size": 20
  }
}
```

### 房间详情

```
GET /api/room/{room_id}
```

**响应**
```json
{
  "code": 0,
  "data": {
    "room_id": "room_xxx",
    "room_name": "房间名",
    "game_mode": "rank",
    "map_id": "map_01",
    "status": "playing",
    "players": [...],
    "settings": {
      "max_players": 8,
      "time_limit": 300,
      "win_score": 1000
    },
    "created_at": "2026-03-07T12:00:00Z"
  }
}
```

---

## 游戏接口 (WebSocket)

### 连接

```
WSS /ws/game/{room_id}?token=<access_token>
```

### 消息格式

**客户端发送 (CS)**
```json
{
  "msg_id": 1,
  "data": { ... }
}
```

**服务器推送 (SC)**
```json
{
  "msg_id": 1001,
  "data": { ... }
}
```

### 消息ID定义

#### 客户端 -> 服务器

| msg_id | 消息名 | 说明 |
|--------|--------|------|
| 1 | CS_HEARTBEAT | 心跳 |
| 10 | CS_CREATE_ROOM | 创建房间 |
| 11 | CS_JOIN_ROOM | 加入房间 |
| 12 | CS_LEAVE_ROOM | 离开房间 |
| 13 | CS_ROOM_LIST | 房间列表 |
| 14 | CS_READY | 准备/取消准备 |
| 15 | CS_START_GAME | 开始游戏 |
| 20 | CS_PLACE_TOWER | 放置防御塔 |
| 21 | CS_UPGRADE_TOWER | 升级防御塔 |
| 22 | CS_SELL_TOWER | 出售防御塔 |
| 23 | CS_SKIP_WAVE | 跳过波次 |
| 30 | CS_CHAT | 聊天消息 |

#### 服务器 -> 客户端

| msg_id | 消息名 | 说明 |
|--------|--------|------|
| 1001 | SC_HEARTBEAT_ACK | 心跳回应 |
| 1010 | SC_CREATE_ROOM_ACK | 创建房间回应 |
| 1011 | SC_JOIN_ROOM_ACK | 加入房间回应 |
| 1012 | SC_LEAVE_ROOM_ACK | 离开房间回应 |
| 1013 | SC_ROOM_LIST_ACK | 房间列表回应 |
| 1014 | SC_ROOM_UPDATE | 房间更新 |
| 1015 | SC_PLAYER_JOIN | 玩家加入 |
| 1016 | SC_PLAYER_LEAVE | 玩家离开 |
| 1017 | SC_GAME_START | 游戏开始 |
| 1018 | SC_GAME_OVER | 游戏结束 |
| 1020 | SC_SYNC_STATE | 状态同步 |
| 1021 | SC_WAVE_START | 波次开始 |
| 1022 | SC_WAVE_END | 波次结束 |
| 1023 | SC_TOWER_PLACED | 塔放置成功 |
| 1024 | SC_TOWER_UPGRADED | 塔升级成功 |
| 1025 | SC_TOWER_SOLD | 塔出售成功 |
| 1030 | SC_DANMAKU | 弹幕消息 |
| 1031 | SC_GIFT | 礼物消息 |
| 1100 | SC_ERROR | 错误消息 |

### CS_READY - 准备

```json
{
  "msg_id": 14,
  "data": {
    "is_ready": true
  }
}
```

### CS_START_GAME - 开始游戏

```json
{
  "msg_id": 15,
  "data": {}
}
```

### CS_PLACE_TOWER - 放置防御塔

```json
{
  "msg_id": 20,
  "data": {
    "tower_type": "arrow",
    "position": {
      "x": 5,
      "y": 3
    }
  }
}
```

### CS_UPGRADE_TOWER - 升级防御塔

```json
{
  "msg_id": 21,
  "data": {
    "tower_id": "tower_001"
  }
}
```

### SC_SYNC_STATE - 状态同步

```json
{
  "msg_id": 1020,
  "data": {
    "frame": 1234,
    "battle_state": {
      "wave": 3,
      "money": 500,
      "lives": 10,
      "score": 2500,
      "enemies": [
        {
          "id": "enemy_001",
          "type": "grunt",
          "hp": 100,
          "max_hp": 100,
          "position": {
            "x": 8.5,
            "y": 2.3
          },
          "progress": 0.45
        }
      ],
      "towers": [
        {
          "id": "tower_001",
          "type": "arrow",
          "level": 2,
          "position": {
            "x": 5,
            "y": 3
          },
          "damage": 25,
          "range": 3.5,
          "attack_speed": 1.0
        }
      ],
      "projectiles": [
        {
          "id": "proj_001",
          "type": "arrow",
          "position": {
            "x": 5.2,
            "y": 3.1
          },
          "target_id": "enemy_001",
          "damage": 25
        }
      ]
    }
  }
}
```

### SC_DANMAKU - 弹幕消息

```json
{
  "msg_id": 1030,
  "data": {
    "user_id": "user_456",
    "nickname": "观众A",
    "content": "主播加油!",
    "type": "text",
    "color": "#FFFFFF"
  }
}
```

### SC_GIFT - 礼物消息

```json
{
  "msg_id": 1031,
  "data": {
    "user_id": "user_456",
    "nickname": "观众A",
    "gift_type": "rocket",
    "gift_name": "火箭",
    "count": 1,
    "total_price": 100
  }
}
```

---

## 支付接口

### 创建订单

```
POST /api/pay/create
```

**请求头**
```
Authorization: Bearer <access_token>
```

**请求体**
```json
{
  "product_id": "coins_100",
  "payment_channel": "douyin",
  "quantity": 1
}
```

**响应**
```json
{
  "code": 0,
  "data": {
    "order_id": "order_20260307120000",
    "amount": 100,
    "currency": "CNY",
    "payment_url": "https://pay.douyin.com/..."
  }
}
```

### 支付回调

```
POST /api/pay/callback
```

**请求体**
```json
{
  "order_id": "order_xxx",
  "status": "success",
  "amount": 100,
  "transaction_id": "txn_xxx",
  "timestamp": 1709808000,
  "signature": "xxxxx"
}
```

### 订单查询

```
GET /api/pay/order/{order_id}
```

**响应**
```json
{
  "code": 0,
  "data": {
    "order_id": "order_xxx",
    "product_id": "coins_100",
    "status": "success",
    "amount": 100,
    "created_at": "2026-03-07T12:00:00Z",
    "paid_at": "2026-03-07T12:01:00Z"
  }
}
```

---

## 排行榜接口

### 获取排行榜

```
GET /api/leaderboard/{type}
```

**路径参数**
| 参数 | 类型 | 说明 |
|------|------|------|
| type | string | 排行榜类型 (score/wealth/win_rate/kills) |

**查询参数**
| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| page | int | 否 | 页码 |
| page_size | int | 否 | 每页数量 |

**响应**
```json
{
  "code": 0,
  "data": {
    "type": "score",
    "entries": [
      {
        "rank": 1,
        "user_id": "user_123",
        "nickname": "冠军",
        "avatar": "https://...",
        "value": 10000,
        "win_rate": 0.85
      },
      {
        "rank": 2,
        "user_id": "user_456",
        "nickname": "亚军",
        "avatar": "https://...",
        "value": 9000,
        "win_rate": 0.80
      }
    ],
    "my_rank": {
      "rank": 50,
      "value": 1000
    }
  }
}
```

---

## 战斗记录接口

### 获取战斗历史

```
GET /api/battle/history
```

**查询参数**
| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| page | int | 否 | 页码 |
| page_size | int | 否 | 每页数量 |

**响应**
```json
{
  "code": 0,
  "data": {
    "battles": [
      {
        "battle_id": "battle_xxx",
        "room_id": "room_xxx",
        "result": "win",
        "score": 2500,
        "wave": 5,
        "duration": 300,
        "damage_dealt": 5000,
        "enemies_killed": 20,
        "created_at": "2026-03-07T12:00:00Z"
      }
    ],
    "total": 100,
    "page": 1,
    "page_size": 20
  }
}
```

### 战斗详情

```
GET /api/battle/{battle_id}
```

**响应**
```json
{
  "code": 0,
  "data": {
    "battle_id": "battle_xxx",
    "room_id": "room_xxx",
    "players": [...],
    "waves": [
      {
        "wave": 1,
        "enemies_spawned": 10,
        "enemies_killed": 10,
        "damage_dealt": 1000
      }
    ],
    "result": "win",
    "score": 2500,
    "rewards": {
      "coins": 100,
      "exp": 50,
      "diamonds": 5
    }
  }
}
```

---

## 抖音 Webhook

### 事件订阅

```
POST /api/webhook/douyin
```

**请求头**
```
X-Douyin-Signature: sha256=xxxxx
X-Douyin-Timestamp: 1709808000
```

**事件类型**

| event_type | 说明 |
|------------|------|
| gift | 礼物事件 |
| danmaku | 弹幕事件 |
| follow | 关注事件 |
| like | 点赞事件 |
| share | 分享事件 |
| enter | 进入直播间 |

**礼物事件 payload**
```json
{
  "event_type": "gift",
  "room_id": "room_xxx",
  "user": {
    "open_id": "user_xxx",
    "nickname": "观众A"
  },
  "gift": {
    "gift_id": "rocket",
    "gift_name": "火箭",
    "count": 1,
    "total_price": 100
  }
}
```

**弹幕事件 payload**
```json
{
  "event_type": "danmaku",
  "room_id": "room_xxx",
  "user": {
    "open_id": "user_xxx",
    "nickname": "观众A"
  },
  "content": "主播加油!"
}
```

---

## 错误码

| code | 说明 |
|------|------|
| 0 | 成功 |
| 1001 | 参数错误 |
| 1002 | 认证失败 |
| 1003 | 权限不足 |
| 1004 | 资源不存在 |
| 1005 | 房间已满 |
| 1006 | 房间已关闭 |
| 1007 | 游戏进行中 |
| 1008 | 金币不足 |
| 1009 | 塔位置无效 |
| 1010 | 支付失败 |
| 2001 | 服务器内部错误 |

---

## 限流规则

| 接口 | 限制 |
|------|------|
| /api/room/create | 10次/分钟 |
| WebSocket 消息 | 100条/秒 |
| /api/pay/create | 5次/分钟 |
| 其他API | 60次/分钟 |
