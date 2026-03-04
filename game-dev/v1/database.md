# 弹幕游戏数据库设计

## MongoDB Collections

### 1. players（玩家表）
```javascript
{
    _id: ObjectId,
    platform: "douyin",      // 平台
    openid: "xxx",           // 平台ID
    nickname: "玩家昵称",
    avatar: "头像URL",
    level: 1,                // 等级
    exp: 0,                  // 经验
    gold: 1000,             // 金币
    gems: 10,               // 钻石
    power: 100,              // 战力
    
    // 统计数据
    stats: {
        total_games: 0,
        wins: 0,
        kills: 0,
        damage: 0
    },
    
    created_at: Date,
    updated_at: Date
}
```

### 2. rooms（房间表）
```javascript
{
    _id: ObjectId,
    room_id: "room_xxx",    // 房间ID
    status: "waiting",       // waiting/playing/finished
    
    // 玩家
    players: [{
        player_id: ObjectId,
        team: "blue/red",
        is_ready: true,
        hero: "xxx"          // 选择的英雄
    }],
    
    // 游戏设置
    settings: {
        mode: "rank/normal",
        map: "map_001",
        time_limit: 300      // 秒
    },
    
    // 战绩
    winner: "blue/red",
    stats: {},
    
    created_at: Date,
    finished_at: Date
}
```

### 3. matches（比赛记录）
```javascript
{
    _id: ObjectId,
    match_id: "match_xxx",
    room_id: ObjectId,
    
    players: [{
        player_id: ObjectId,
        team: "blue/red",
        hero: "xxx",
        kills: 0,
        damage: 0,
        gold_earned: 0
    }],
    
    winner: "blue/red",
    duration: 180,           // 秒
    
    created_at: Date
}
```

### 4. anchors（主播表）
```javascript
{
    _id: ObjectId,
    platform: "douyin",
    anchor_id: "xxx",
    nickname: "主播昵称",
    followers: 10000,
    
    // 认证信息
    verified: true,
    verify_info: {},
    
    // 统计数据
    stats: {
        total_games: 0,
        avg_viewers: 0,
        revenue: 0
    },
    
    created_at: Date,
    updated_at: Date
}
```

### 5. items（道具表）
```javascript
{
    _id: ObjectId,
    name: "英雄解锁",
    type: "hero/skin/battle_pass",
    price_gold: 0,
    price_gems: 688,
    description: "解锁新英雄",
    
    // 限制
    limit_type: "none/daily/once",
    limit_value: 0
}
```

---

## Redis Keys

### 1. 会话管理
```
barrage:session:{session_id}  // WebSocket会话
barrage:player:{player_id}:room  // 玩家当前房间
```

### 2. 缓存
```
barrage:room:{room_id}      // 房间状态缓存
barrage:online:count       // 在线人数
barrage:rankings:weekly    // 周排行榜
```

### 3. 计数器
```
barrage:player:{player_id}:games:day_{date}
barrage:anchor:{anchor_id}:viewers
```

---

## 索引设计

### players
```javascript
{ platform: 1, openid: 1 }  // 唯一
{ exp: -1 }                 // 等级榜
{ gold: -1 }                // 财富榜
```

### rooms
```javascript
{ status: 1 }               // 查找可加入房间
{ created_at: -1 }          // 最新房间
```

### matches
```javascript
{ players.player_id: 1 }   // 玩家历史
{ created_at: -1 }          // 最近比赛
```
