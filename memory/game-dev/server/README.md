# 弹幕游戏服务器

## 快速开始

```bash
# 构建
cd server && go build -o danmaku-game ./cmd

# 运行
./danmaku-game

# 测试
go test ./... -v
```

## 项目结构

```
server/
├── cmd/           # 入口
├── internal/      # 内部包
├── bot/           # AI 机器人
├── deploy/       # 部署脚本
├── health/       # 健康检查
├── logger/       # 日志
├── monitor/      # 监控
├── proto/        # 协议
└── security/     # 安全模块
```

## 核心模块

| 模块 | 文件 | 功能 |
|------|------|------|
| 战斗 | battle_tower.go | 塔防战斗 |
| AI | battle_ai.go | 敌人AI |
| 技能 | battle_skill.go | 战斗技能 |
| 帧同步 | frame_sync.go | 同步 |
| WebSocket | websocket.go | 实时通信 |
| 房间 | room_manager.go | 房间管理 |
| 排行榜 | leaderboard.go | 排行 |
| 成就 | achievement_quest.go | 成就任务 |
| 支付 | payment.go | 支付 |
| 缓存 | cache.go | Redis缓存 |
| 熔断 | circuit_breaker.go | 容错 |
| 灰度 | gray_release.go | 灰度发布 |
| 监控 | alerting.go | 告警 |

## API 端口

- HTTP: 8080
- WebSocket: 8081
- 健康检查: 8080/health

## 配置

修改 `config.go` 或使用环境变量。
