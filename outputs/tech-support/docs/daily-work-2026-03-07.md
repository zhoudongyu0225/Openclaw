# 技术支持日报 - 2026-03-07

## 服务状态

| 项目 | 状态 |
|------|------|
| Gateway | 运行中 (pid 3932199) |
| 节点服务 | 未安装 |
| Feishu | ✅ OK |
| WeCom | ⚠️ SETUP (未配置) |
| QQ Bot | ⚠️ SETUP (未配置) |

## 安全审计

### 🔴 CRITICAL (3项)
1. **Open groupPolicy + elevated tools** - channels.feishu.groupPolicy="open" 且启用了 elevated tools
2. **Open groupPolicy + runtime/filesystem** - 风险工具暴露
3. **Feishu 安全警告** - 任何成员都可以触发提及

### 🟡 WARN (9项)
- Reverse proxy headers 未信任
- Feishu doc create 可授予权限
- 多用户潜在风险

### 建议
- 设置 `groupPolicy="allowlist"`
- 配置 `tools.fs.workspaceOnly=true`

## 插件状态

| 插件 | 状态 |
|------|------|
| feishu_doc/chat/wiki/drive/bitable | ✅ 正常 |
| wecom | ❌ 加载失败 (TypeError: api.registerHttpHandler is not a function) |

## 更新

- pnpm/npm 有可用更新 (2026.3.7)
- 建议运行 `openclaw update`

## 资源使用

- Session: 162 个活动
- Memory: 0 files, 缓存开启
- 模型: MiniMax-M2.5 (200k ctx)

## 待处理

- [ ] 修复 wecom 插件加载问题
- [ ] 解决安全审计中的 CRITICAL 项
- [ ] 执行 `openclaw update` 更新
