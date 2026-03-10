# MEMORY.md - 长期记忆

## 重要原则

- 分身记忆：不同业务用不同分身，广告创意用 `memory/ad-creative/`，每次启动先读对应分身的 README.md

## 记忆系统架构

### 三层记忆（参考 EverMemOS）

| 类型 | 位置 | 功能 |
|------|------|------|
| Working Memory | `~/.openclaw/memory/working/` | 当前会话上下文 |
| MemCell | `~/.openclaw/memory/memcell/` | 原子记忆单元 |
| Episode | `~/.openclaw/memory/episode/` | 主题聚合记忆 |
| Profile | `~/.openclaw/memory/profile/` | 用户/分身画像 |

### 使用方法

```bash
# 提取记忆
~/.openclaw/memory/memory-system.sh memorize <session_id> <content>

# 搜索记忆
~/.openclaw/memory/memory-system.sh recall <query>

# 查看状态
~/.openclaw/memory/memory-system.sh status
```

## 分身列表

| 分身 | 目录 | 业务 |
|------|------|------|
| main | memory/main/ | 龙虾统筹 - 整体规划 |
| content-ops | memory/content-ops/ | 眠心诊所 - 内容运营 |
| ecommerce | memory/ecommerce/ | 氢美健康 - 电商运营 |
| game-dev | memory/game-dev/ | 弹幕游戏 - 游戏开发 |
| ad-creative | memory/ad-creative/ | 游戏广告 - 广告创意 |
| art-gen | memory/art-gen/ | 美术素材 - AI生成 |
| tech-support | memory/tech-support/ | 技术支持 |

**当前：ad-creative（广告创意）**

---

## 项目详情（2026-03-02 更新）

### 1. 眠心诊所（content-ops）
- 运营主体：成都高新眠心诊所
- 核心竞争力：**一口价治疗套餐 7500 元/月，确定周期治愈**
- 营销策略：先触达核心用户，再线上转化泛用户
- 渠道：公众号、小红书、抖音

### 2. 氢美健康（ecommerce）
- 公司名：氢美健康
- 产品：氢镁胶囊，主要成分氢化镁
- 卖点：改善睡眠/清除体内垃圾
- 法人：丁文江院士儿子
- 网站要求：高端简约
- 运营：线上 + 线下

### 3. 弹幕游戏（game-dev）
- 平台：抖音、快手（国内）→ tiktok（海外）
- 玩法：恐龙攻守塔防，双主播 PK
- 胜负：基地被恐龙破坏即失败
- 后端：主播 PK、SDK、API、通信、玩家管理、GM 后台
- 技术：UE4 + Go + MongoDB + Redis + Protobuf

### 4. 游戏广告（ad-creative）
- 万国觉醒：不带主角的埃及主题"模拟城市"
- 万龙觉醒：类 TopHeroes，ARPG + 模拟经营
- 工具：广大大

### 5. 美术素材（art-gen）
- 职责：UI、原画、3D模型、贴图、动画、特效、视频
- 当前：配合弹幕游戏
- UI 参考：《红警 OL》手游
- 方向：3D 动作独立游戏

### 6. 技术支持（tech-support）
- 职责：工具开发、IT 运维
- 当前：配合广告创意开发试玩 HTML

---

## 每日汇报要求（2026-02-27 新增）

- 每天 12:00 自动执行 cron 任务汇报
- 检查所有6个分身的状态
- 汇报格式：每个分身的状态、实际完成内容、遇到的问题
- 原则：诚实汇报，不编造进度
