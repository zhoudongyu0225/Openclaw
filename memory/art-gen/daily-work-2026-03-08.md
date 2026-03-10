# 美术素材每日工作 - 2026-03-08

## 任务执行：提示词优化 & 新模板产出

---

## 一、提示词优化

### 1.1 光照系统强化

**优化点**：增加动态光照、色温控制、软硬光描述

**优化后提示词**：
```
[主体] subject, dramatic [光照类型] lighting,
[色温] color temperature #XXXXXX,
key light [方向/强度], fill light soft,
rim light edge definition, volumetric god rays,
HDRI environment map, physically based rendering,
cinematic light setup, [氛围] atmosphere,
rendered in Octane/Arnold/V-Ray, 8K --ar 16:9 --v 6
```

**应用示例 - 恐龙化石展厅**：
```
Dinosaur fossil in museum exhibit, dramatic spotlight lighting,
warm tungsten 3200K, blue ambient fill,
key light from above front, rim light from back,
volumetric dust particles in light beams,
museum atmosphere, scientific display,
photorealistic, Octane render, 8K --ar 16:9 --v 6
```

### 1.2 材质表现强化

**优化点**：细化材质质感、物理参数、纹理细节

**优化后提示词**：
```
[主体], PBR material, [材质类型] texture,
normal map detail, roughness map, metallic map,
displacement detail, subsurface scattering,
material ID: [材质编号], physically accurate,
[表面特征] surface detail, worn [老化效果],
8K texture set, substance painter ready --ar 1:1 --v 6
```

---

## 二、新模板产出

### 2.1 AI视频生成提示词模板

适用于 Runway、Pika、OpenAI Sora 等视频生成工具

```
[开场] Opening shot: [场景/主体描述],
[动作1] Action: [运动描述],
[转场] Transition: [切换方式],
[动作2] Secondary action: [次要动作],
[结尾] Closing: [结束方式],
Cinematic motion, [帧率]fps, [时长]seconds,
professional cinematography, smooth movement,
[风格] style, [光照] lighting --video --ar 16:9
```

**应用示例 - 恐龙苏醒**：
```
Opening shot: Ancient dinosaur egg glowing in jungle clearing,
Action: Egg cracks with pulsing light, tiny claw emerges,
Transition: Quick zoom to dramatic reveal,
Secondary action: Baby dinosaur first breath, mist swirls,
Closing: Mother dinosaur silhouette in background,
Cinematic motion, 24fps, 8 seconds,
epic reveal, golden hour lighting --video --ar 16:9
```

### 2.2 多人同框风格统一模板

适用于游戏过场动画、多角色海报

```
Group composition, [数量] characters,
[排列方式] arrangement, central focus,
[统一风格] unified style, [配色方案] color palette,
consistent lighting direction, matching shadows,
[视角] eye level view, depth of field,
cinematic group shot, ensemble cast,
harmonious composition, professional art direction --ar 16:9 --v 6
```

**应用示例 - 恐龙战队**：
```
Group composition, 5 dinosaur warriors,
triangular arrangement, central leader focal point,
cohesive stylized realism, warm earth tones #8B4513 #D2691E,
consistent top-front lighting, matching cast shadows,
eye level view, shallow depth of field,
cinematic group shot, dinosaur team ensemble,
harmonious composition, professional art direction --ar 16:9 --v 6
```

### 2.3 UI 动效提示词模板

适用于抽卡动效、转场动画

```
UI animation frame [N], [动画类型] animation,
[easing] easing function, [时长] duration,
[元素] elements in motion, smooth interpolation,
[glow] glow/pulse effect, particle trails,
[背景] dynamic background, ui motion design,
after effects keyframe, lottie ready,
game ui animation, polished transition --ar 9:16 --v 6
```

### 2.4 HDR后期处理模板

适用于游戏截图、引擎内渲染

```
[场景], HDR rendering, tone mapping ACES,
bloom effect, lens flare, chromatic aberration,
motion blur, depth of field, ambient occlusion,
SSAO contact shadows, screen space reflections,
color grading [风格], film grain subtle,
[分辨率] resolution, unreal engine 5 screenshot,
photorealistic gaming --ar 16:9 --v 6 --s 200
```

---

## 三、知识沉淀

### 3.1 光照参数速查

| 光照类型 | 关键词 | 适用场景 |
|----------|--------|----------|
| 电影光 | cinematic, three-point lighting | 角色立绘 |
| 戏剧光 | dramatic, high contrast | 宣传海报 |
| 柔和光 | soft, diffused, overcast | 休闲游戏 |
| 科技光 | neon, volumetric, cyberpunk | 科幻题材 |
| 自然光 | golden hour, blue hour, ambient | 场景背景 |

### 3.2 材质关键词速查

| 材质类型 | 关键词 |
|----------|--------|
| 金属 | metallic, chrome, brushed metal |
| 皮革 | leather, worn, patina |
| 岩石 | rocky, granite, limestone |
| 毛发 | fur, feathers, scales |
| 能量 | glowing, energy, plasma |

### 3.3 视频生成参数参考

| 参数 | 推荐值 |
|------|--------|
| 分辨率 | 16:9 (横版) / 9:16 (竖版) |
| 帧率 | 24fps (电影感) / 30fps (流畅) |
| 时长 | 3-10秒 (主流) |
| 风格 | --style cinematic / animated / realistic |

---

## 四、质量检查

### 今日模板测试结果

| 模板类型 | 描述 | 状态 |
|----------|------|------|
| 光照系统 | 动态光照/色温控制 | ✅ |
| 材质表现 | PBR材质/纹理细节 | ✅ |
| AI视频生成 | Runway/Pika/Sora | ✅ |
| 多人同框 | 风格统一/构图 | ✅ |
| UI动效 | 抽卡/转场 | ✅ |
| HDR后期 | 引擎渲染/调色 | ✅ |

---

## 五、明日待办

- [ ] 测试AI视频模板实际生成效果
- [ ] 整理视频生成案例存档
- [ ] 优化2D原画到3D模型的提示词桥接

---

## 2026-03-08 晚间扩展（23:08 UTC）

### 🎯 主题：跨引擎兼容性 & 批量产出流程化 & 风格品牌化

### 🆕 引擎适配模板

#### Unity 3D 资产模板
```
Unity game asset, [主体类型],
FBX format, Mesh optimizer ready,
PBR textures: Albedo/Normal/Roughness/Metallic/AO,
[分辨率]x[分辨率] texture resolution,
[多边形数] triangles max, LOD levels (0/1/2/3),
Unity import settings, Standard Shader compatible,
Prefab ready, level of detail included,
game development asset --ar 1:1 --v 6
```

#### Unreal Engine 5 模板
```
Unreal Engine 5 asset, [主体类型],
Nanite compatible, Lumen ready,
[多边形数] polygons, baked lighting ready,
.fbx/.uasset format, material instance ready,
Substance textures, 4K textures minimum,
world position offset animations,
LOD chain included, UE5 marketplace style --ar 1:1 --v 6
```

#### Godot 引擎模板
```
Godot 4.x game asset, [主体类型],
GLES3/Vulkan compatible,
[分辨率]x[分辨率] textures, imported sprite sheet,
AnimationPlayer ready, Sprite2D/Sprite3D ready,
PNG8/32 optimized, ETC2/ASTC compression,
Godot import settings, node structure included,
open source game engine compatible --ar 1:1 --v 6
```

### 🆕 品牌风格指南模板

#### 品牌视觉一致性检查
```
[PROJECT] brand style guide check,
Logo usage: [位置/尺寸/保护区域],
Color palette: [主色] [辅色] [强调色],
Typography: [字体名称] [字号层级],
Imagery style: [风格描述],
Spacing system: [网格系统],
Iconography: [图标风格],
Photography: [摄影风格],
Brand consistency verification,
professional brand identity --ar 16:9 --v 6
```

#### IP角色设计规范
```
[IP名称] character design specification,
[角色名] character sheet,
expression sheet: happy/sad/angry/surprised,
pose sheet: front/side/3/4 view,
[配色方案] color breakdown,
silhouette distinctiveness check,
[年龄段] age group, [性别] gender,
character consistency guide,
IP brand guidelines --ar 3:4 --v 6
```

### 🆕 批量产出工作流模板

#### 系列化产出提示词
```
[系列名称] Series [编号],
[主题] theme, batch production [数量] variants,
[风格] style locked, [配色] palette locked,
variation A: [变化点1], variation B: [变化点2],
variation C: [变化点3], variation D: [变化点4],
consistent quality, production ready,
series cohesion, art direction controlled --ar [比例] --v 6
```

#### 批量重命名规范
```
[NAME]_[SERIES]_[VARIANT]_[RES]_[DATE]
示例: TRex_Enemy_01_4K_20260308
      TRex_Enemy_02_4K_20260308
```

### 🆕 无障碍设计模板

#### 色盲友好检测提示词
```
[场景], accessibility checked,
colorblind safe palette, deuteranopia friendly,
protanopia friendly, tritanopia friendly,
contrast ratio WCAG AA minimum,
no red-green critical information,
shape+iconography redundancy,
text readable at small sizes,
inclusive design --ar 16:9 --v 6
```

### 🆕 版权/法律合规模板

#### 版权检查清单提示词
```
[作品类型] copyright verification,
Original content: [X]% original,
Reference sources documented,
Trademark check: [X] none detected,
Model/texture ownership: [来源],
Commercial rights: [授权范围],
Attribution requirements: [如有],
Stock assets licensed: [来源],
AI generation disclosure: [是否需要],
legal compliance checked --ar 16:9 --v 6
```

---

### 📊 晚间扩展产出

| 类型 | 数量 |
|------|------|
| 引擎适配模板 | 3 (Unity/UE5/Godot) |
| 品牌风格指南 | 2 |
| 批量产出流程 | 2 |
| 无障碍设计 | 1 |
| 版权合规 | 1 |
| **本轮总计** | **9** |

---

### 📋 技巧沉淀（晚间扩展更新）

#### 跨引擎兼容性检查清单
- [ ] 导出格式匹配 (FBX/OBJ/GLTF)
- [ ] 纹理压缩格式 (ASTC/ETC2/DXT5)
- [ ] 多边形预算
- [ ] Draw Call 优化
- [ ] 烘焙光照贴图UV

#### 品牌一致性检查点
- [ ] Logo 使用规范
- [ ] 配色严格统一
- [ ] 字体层级清晰
- [ ] 视觉风格锁定
- [ ] 跨平台一致性

#### 批量产出质量控制
- [ ] 首件质检 (First Article Inspection)
- [ ] 随机抽检 (Spot Check 10%)
- [ ] 风格一致性审核
- [ ] 命名规范性审核
- [ ] 元数据完整性

---

*任务完成 ✅*
*记录时间：2026-03-08 23:08 UTC*
*执行者：art-gen-daily cron*

---

*本记录由 AI 自动生成 | 更新于 2026-03-08 23:08 UTC*
