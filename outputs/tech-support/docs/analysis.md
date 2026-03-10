# 氢美健康（逆动力）网站性能优化分析

**分析日期**: 2026-03-07  
**文件**: projects/nidongli/index-v4.html  
**原始大小**: 24KB, 685行

---

## 🔍 发现的问题

### 1. Google Fonts 阻塞渲染
- 使用 `fonts.googleapis.com` 加载字体
- 国内访问极慢，会阻塞页面渲染
- 缺少 `font-display: swap`

### 2. 外部图片无懒加载
- 直接使用 Unsplash 图片 URL
- 没有使用 `loading="lazy"` 属性
- 没有设置 width/height，导致 CLS (布局偏移)

### 3. 缺少关键优化标签
- 没有 `<meta name="description">` SEO
- 没有预加载关键资源 (`preload`)
- 没有预连接外部域名 (`preconnect`)

---

## ✅ 优化方案

### 优化 1: 使用国内可访问的字体服务 + font-display
```html
<!-- 替换 Google Fonts -->
<link rel="preconnect" href="https://fonts.font.im">
<link href="https://fonts.font.im/css/css2?family=Noto+Sans+SC:wght@300;400;500&family=Playfair+Display:wght@400;500;600&family=Inter:wght@300;400;500;600&display=swap" rel="stylesheet">
<!-- 或使用系统字体栈作为 fallback -->
```

### 优化 2: 图片懒加载 + 尺寸属性
```html
<img loading="lazy" decoding="async" width="400" height="300" src="...">
```

### 优化 3: 预连接 + SEO + Meta 标签
```html
<link rel="preconnect" href="https://images.unsplash.com">
<meta name="description" content="逆动力氢镁胶囊...">
```

---

## 📁 输出文件

- `index-v4-optimized.html` - 优化后的版本 (已复制到 projects/nidongli/index-v5.html)

---

## 🚀 3 点核心优化建议（已实现）

### 优化 1: 替换 Google Fonts 为国内可用源 + 系统字体 fallback
**问题**: fonts.googleapis.com 在国内访问极慢，阻塞渲染
**解决**: 
- 使用 font.im 镜像源
- 添加系统字体栈作为第一选择 (-apple-system, BlinkMacSystemFont 等)
- 使用 JS 监听字体加载完成后切换

```html
<!-- 之前 -->
<link href="https://fonts.googleapis.com/css2?family=..." rel="stylesheet">

<!-- 优化后 -->
<link rel="preconnect" href="https://fonts.font.im" crossorigin>
<link href="https://fonts.font.im/css2?family=..." rel="stylesheet">

<!-- CSS 中添加系统字体 -->
font-family: -apple-system, BlinkMacSystemFont, 'Inter', 'Noto Sans SC', sans-serif;
```

### 优化 2: 图片添加懒加载 + 尺寸属性
**问题**: 图片没有指定宽高 → CLS (布局偏移)；非首屏图片没有懒加载 → 浪费带宽
**解决**:
- 首屏图片: `loading="eager" decoding="async"`
- 非首屏图片: `loading="lazy" decoding="async"`
- 所有图片添加 `width="" height=""` 属性

```html
<!-- 之前 -->
<img class="hero-img" src="..." alt="氢镁胶囊">

<!-- 优化后 -->
<img class="hero-img" src="..." alt="氢镁胶囊" width="400" height="500" loading="eager" decoding="async">
<img class="product-img" src="..." alt="产品" width="600" height="600" loading="lazy" decoding="async">
```

### 优化 3: 添加 SEO + 预连接
**问题**: 缺少 meta 描述、关键词；没有预连接外部资源
**解决**:
- 添加 meta description、keywords
- 添加 Open Graph 标签
- 预连接 images.unsplash.com 和 fonts.font.im

```html
<meta name="description" content="逆动力氢镁胶囊，中国氢科学研究院研发...">
<meta name="keywords" content="氢镁胶囊,氢健康,抗氧化...">
<link rel="preconnect" href="https://images.unsplash.com">
```

---

## 📊 性能提升预期

| 指标 | 优化前 | 优化后 |
|------|--------|--------|
| 首屏渲染 | 受 Google Fonts 阻塞 | 系统字体即时显示 |
| CLS | 无尺寸属性，有偏移 | 添加宽高，零偏移 |
| LCP | - | 预连接加速 |
| SEO | 无 | 完整 meta 标签 |
| 国内访问 | 慢 | 使用 font.im 镜像 |
