# 技术支持日报 - 2026-03-07

> 注：实际执行日期 2026-03-10 15:22 UTC

## 服务状态检查

### OpenClaw Gateway
- **状态**: ✅ 运行中 (pid: 1183578)
- **端口**: 127.0.0.1:18789
- **RPC Probe**: ✅ OK
- **系统服务**: systemd installed, enabled, running
- **运行时间**: 14 days 5 hours

### 通道状态
| 通道 | 状态 |
|------|------|
| Feishu | ✅ OK (configured) |
| QQ Bot | ⚠️ SETUP (not configured) |
| ADP OpenClaw | ⚠️ SETUP (not configured) |
| WeCom | ❌ OFF |
| DingTalk | ❌ OFF |

## 系统资源

- **磁盘**: 40G / 22G used / 19G available (55%)
- **负载**: 0.72, 0.41, 0.28 (正常)
- **运行时长**: 14 days 5 hours

## ⚠️ 安全警告 (8个)

### 高危
1. **Host-header origin fallback 已启用**: `gateway.controlUi.dangerouslyAllowHostHeaderOriginFallback=true` 启用 DNS rebinding 保护
   - 建议: 关闭此选项或配置 `gateway.controlUi.allowedOrigins`

### 中危
2. **Reverse proxy headers 未信任**: 如通过反向代理暴露 Control UI，需配置 `gateway.trustedProxies`
3. **Feishu doc 创建可授予权限**: `feishu_doc action="create"` 可授予文档访问权限
4. **Extensions 加载但 plugins.allow 未设置**: 发现 5 个扩展但未设置允许列表
5. **扩展插件工具可达**: adp-openclaw, ddingtalk, qqbot, skillhub, wecom 在默认策略下可访问

### 低危
6. **Skillhub 插件警告**: 加载时无安装路径来源记录

## 🚀 技术优化建议

### 立即建议
1. **关闭危险配置**: 编辑 `~/.openclaw/config.yaml`:
   ```yaml
   gateway:
     controlUi:
       dangerouslyAllowHostHeaderOriginFallback: false
   ```

2. **设置插件白名单**: 添加:
   ```yaml
   plugins:
     allow:
       - feishu
       - feishu_doc
       - feishu_chat
       - feishu_wiki
       - feishu_drive
       - feishu_bitable
   ```

3. **升级版本**: 有可用更新 v2026.3.8
   - 运行: `openclaw update`

### 可选优化
- 考虑从 NVM Node 迁移到系统 Node 22+ 提高稳定性

## 操作记录
- [x] 检查 Gateway 状态
- [x] 检查系统资源
- [x] 运行安全审计
- [x] 记录日报

---
记录时间: 2026-03-10 15:22 UTC
