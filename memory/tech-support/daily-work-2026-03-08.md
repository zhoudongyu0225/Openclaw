# 技术支持日报 - 2026-03-08

## 服务状态

| 项目 | 状态 |
|------|------|
| Gateway | 运行中 (pid 347514) |
| RPC Probe | 正常 |
| Port | 18789 (loopback) |
| Systemd | enabled |

## 系统资源

- **磁盘**: 40G / 22G used / 19G available (55%)
- **内存**: 1.9Gi total, 1.0Gi used, 523Mi free, 929Mi available
- **Swap**: 8.0Gi total, 413Mi used
- **Load**: 0.64, 0.23, 0.12 (12天运行)
- **Uptime**: 12 days, 13:23

## 技术优化建议

### ⚠️ 需关注

1. **Node版本管理器依赖**: Gateway 使用 nvm 管理的 Node 22.22.0
   - 服务配置警告: 升级后可能失效
   - 建议: 安装系统 Node 22+ 迁移

2. **可选API Key缺失**:
   - voyage (embedding)
   - mistral (模型)
   - 不影响基本功能，如需使用可配置

### ✅ 正常

- Gateway 绑定 loopback，安全
- 磁盘空间充足
- 内存充足
- Load 低

## 活跃会话

- cron 定时任务多个运行中
- feishu group 活跃

---
记录时间: 2026-03-08 23:11 UTC
