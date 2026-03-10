# 技术支持 - 专家级知识库

## 职责定位
- 全站工程师中台
- 服务各项目技术需求
- 工具开发
- 运维支持

## 当前任务

### 阶段一：配合广告创意
- 试玩游戏 HTML 开发
- 投放落地页
- 数据追踪

### 阶段二：全站中台
- 通用工具开发
- CI/CD 建设
- 监控系统

## 技术栈

### 前后端
| 方向 | 技术 | 用途 |
|------|------|------|
| 前端 | HTML/CSS/JS | 静态页面 |
| 前端框架 | React/Vue | 复杂交互 |
| 后端 | Node.js/Go | API 服务 |
| 数据库 | MongoDB/MySQL | 数据存储 |

### 运维
| 方向 | 技术 | 用途 |
|------|------|------|
| 容器 | Docker | 环境隔离 |
| 部署 | Vercel/Netlify | 前端托管 |
| 部署 | Docker Compose | 本地开发 |
| CI/CD | GitHub Actions | 自动化 |

### 开发工具
| 用途 | 工具 |
|------|------|
| 代码 | VS Code |
| 版本控制 | Git |
| 调试 | Chrome DevTools |
| API 测试 | Postman/Apifox |

## 试玩游戏开发规范

### 目录结构
```
trial-game/
├── index.html      # 主入口
├── css/
│   └── style.css   # 样式
├── js/
│   ├── game.js     # 游戏主逻辑
│   ├── battle.js   # 战斗系统
│   ├── enemy.js    # 敌人生成
│   ├── tower.js    # 防御塔
│   └── storage.js  # 本地存储
├── assets/
│   ├── images/     # 图片资源
│   └── audio/      # 音效
└── README.md       # 说明文档
```

### 代码规范

#### HTML
- 使用语义化标签
- class 命名：BEM 风格
- 移动端适配

#### CSS
- 使用 CSS 变量
- Flex/Grid 布局
- 响应式设计

#### JavaScript
- ES6+ 语法
- 模块化（import/export）
- 注释清晰
- 变量命名有意义

### 性能优化

#### 加载优化
- 压缩 CSS/JS
- 图片懒加载
- CDN 加速

#### 渲染优化
- 减少重排重绘
- 使用 transform/opacity
- requestAnimationFrame

#### 内存优化
- 及时清理事件监听
- 对象池复用
- 避免内存泄漏

## 部署指南

### Vercel 部署
```bash
# 安装 vercel
npm i -g vercel

# 登录
vercel login

# 部署
vercel --prod
```

### GitHub Actions 自动部署
```yaml
name: Deploy
on:
  push:
    branches: [main]
jobs:
  deploy:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v2
      - uses: amondnet/vercel-action@v20
        with:
          vercel-token: ${{ secrets.VERCEL_TOKEN }}
          vercel-org-id: ${{ secrets.ORG_ID }}
          vercel-project-id: ${{ secrets.PROJECT_ID }}
```

## 监控与日志

### 前端监控
- 错误收集：window.onerror
- 性能：Performance API
- 用户行为：埋点

### 服务监控
- 状态检查：Health Check API
- 资源：CPU/内存/磁盘
- 日志：结构化 JSON 日志

### 告警
- 错误率 > 5%
- 响应时间 > 2s
- 磁盘 > 80%

## 常用脚本

### 一键部署
```bash
#!/usr/bin/env bash
set -e

echo "=== 开始部署 ==="

# 构建
npm run build

# 部署
vercel --prod --yes

echo "=== 部署完成 ==="
```

### 数据库备份
```bash
#!/usr/bin/env bash
DATE=$(date +%Y%m%d)
mongodump --db myapp --out backup/$DATE
```

### 日志分析
```bash
# 统计错误
grep -i error logs/app.log | wc -l

# 查看实时日志
tail -f logs/app.log
```

## 常见问题

### 前端页面打不开
1. 检查控制台错误
2. 检查网络请求
3. 刷新缓存
4. 查看部署状态

### 部署失败
1. 检查依赖是否完整
2. 检查环境变量
3. 查看部署日志
4. 回滚上一个版本

### 性能问题
1. 使用 Chrome DevTools Profiler
2. 检查网络请求
3. 查看代码是否有死循环
4. 优化图片/资源

### 数据库连接失败
1. 检查连接字符串
2. 检查数据库服务状态
3. 检查防火墙
4. 查看连接池

## 安全规范

### 代码安全
- 不在前端存敏感信息
- 输入验证
- 输出转义
- 依赖更新

### 部署安全
- 密钥放环境变量
- 最小权限原则
- 定期更换密码
- 开启 2FA

## 学习资源

### 前端
- MDN Web Docs
- CSS-Tricks
- JavaScript.info

### 运维
- Docker 文档
- Vercel 文档
- GitHub Actions 文档

### 工具
- Stack Overflow
- GitHub Issues
- 掘金/知乎技术专栏
