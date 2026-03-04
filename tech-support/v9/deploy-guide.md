# 部署指南

## 一、本地开发

### 1.1 环境要求
- Node.js 18+
- npm / yarn
- Git

### 1.2 启动开发服务器
```bash
cd trial-game
npx serve .
# 或
python -m http.server 8000
```

### 1.3 访问
http://localhost:3000

---

## 二、构建优化

### 2.1 HTML优化
- 压缩CSS/JS
- 内联关键CSS
- 延迟加载非关键资源

### 2.2 资源优化
- 图片压缩（TinyPNG）
- 精灵图合并
- 字体子集化

---

## 三、部署平台

### 3.1 Vercel（推荐）
```bash
npm i -g vercel
vercel
# 自动部署到 https://xxx.vercel.app
```

### 3.2 Cloudflare Pages
1. 连接到GitHub仓库
2. 构建命令：空或 `echo "static"`
3. 输出目录：./

### 3.3 GitHub Pages
```bash
git checkout -b gh-pages
cp -r trial-game/* .
git commit -m "deploy"
git push origin gh-pages
```

---

## 四、域名配置

### 4.1 Vercel
- 自动配置SSL
- 支持自定义域名

### 4.2 Cloudflare
- 添加域名
- 配置DNS解析

---

## 五、监控维护

### 5.1 访问统计
- Vercel Analytics
- Google Analytics

### 5.2 错误监控
- Sentry
- LogRocket
