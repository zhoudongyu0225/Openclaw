# HTML5游戏 性能优化方案

## 1. 资源加载优化

### 懒加载
```javascript
// 图片懒加载
const observer = new IntersectionObserver((entries) => {
    entries.forEach(entry => {
        if (entry.isIntersecting) {
            loadImage(entry.target);
        }
    });
});
```

### 资源合并
- 合并小图片为Sprite Sheet
- 合并CSS文件
- 合并JS模块

## 2. 渲染优化

### Canvas优化
```javascript
// 离屏渲染
const offscreen = document.createElement('canvas');
const offCtx = offscreen.getContext('2d');

// 绘制到离屏Canvas
drawGame(offCtx);

// 一次性绘制到主Canvas
ctx.drawImage(offscreen, 0, 0);
```

### requestAnimationFrame
```javascript
let lastTime = 0;
const FPS = 60;
const interval = 1000 / FPS;

function gameLoop(timestamp) {
    requestAnimationFrame(gameLoop);
    
    const delta = timestamp - lastTime;
    if (delta > interval) {
        lastTime = timestamp - (delta % interval);
        update();
        render();
    }
}
```

## 3. 内存优化

### 对象池
```javascript
class ObjectPool {
    constructor(factory, initialSize = 10) {
        this.pool = [];
        this.factory = factory;
        
        for (let i = 0; i < initialSize; i++) {
            this.pool.push(this.factory());
        }
    }
    
    get() {
        return this.pool.pop() || this.factory();
    }
    
    release(obj) {
        this.pool.push(obj);
    }
}
```

## 4. 打包优化

### Vite配置
```javascript
// vite.config.js
export default {
    build: {
        rollupOptions: {
            output: {
                manualChunks: {
                    'vendor': ['vue'],
                    'game': ['./src/game/**/*.js']
                }
            }
        }
    }
}
```
