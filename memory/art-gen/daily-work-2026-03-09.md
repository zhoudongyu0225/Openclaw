# 美术素材每日工作 - 2026-03-09

## 任务执行：提示词优化 & 新模板产出

---

## 一、提示词优化

### 1.1 角色表情与情感表达强化

**优化点**：细化面部表情、情感层次微表情引导

**优化后提示词**：
```
[角色类型] character, [性别/年龄] demographic,
[表情类型] facial expression, [情感核心] core emotion,
[微表情] micro-expression, [眼神] eye direction/gaze,
[眉毛/嘴角] eyebrow/mouth detail, [身体语言] body language,
[心理状态] psychological state, [内心独白感] inner narrative feeling,
portrait style, [视角] portrait angle,
studio lighting, catch lights in eyes,
emotional depth, character study,
cinematic portrait, [风格] style --ar 3:4 --v 6
```

**应用示例 - 恐龙饲养员**：
```
Middle-aged dinosaur caretaker, weary but determined expression,
tired eyes with hope beneath, slight frown of concern,
eyebrows slightly furrowed, hands clasped nervously,
protective instinct, portrait from chest up,
soft window lighting, subtle catch lights,
emotional depth, character study,
cinematic portrait, documentary realism --ar 3:4 --v 6
```

### 1.2 机械/科技元素强化

**优化点**：科技质感、界面元素、发光部件描述

**优化后提示词**：
```
[主体], sci-fi [设备类型] design,
[风格] tech aesthetic: [赛博/复古/极简],
holographic UI elements, LED glow strips,
[发光颜色] emission color #XXXXXX,
metallic panels, exposed circuitry,
[功能性细节] functional details, vents and grills,
brushed metal texture, scratches and wear,
[科技感] futuristic feel, [时代背景] era,
screen displays: [界面内容], antenna/sensors,
sci-fi prop, detailed mechanical design,
PBR textures, 8K, [渲染引擎] --ar 1:1 --v 6
```

**应用示例 - 恐龙追踪器**：
```
Handheld dinosaur tracking device, rugged sci-fi design,
cyberpunk aesthetic, holographic display floating above,
blue emission #00BFFF, scratched aluminum body,
exposed circuit boards on back panel,
rubber grip texture, functional buttons and dials,
LCD screen showing dinosaur silhouette map,
satellite dish antenna, LED status lights,
futuristic feel, year 2085 setting,
GPS coordinates, heart rate monitor display,
sci-fi prop, detailed mechanical design,
Octane render, 8K --ar 1:1 --v 6
```

---

## 二、新模板产出

### 2.1 连续帧/故事板提示词模板

适用于生成连贯的动画帧、漫画分镜、游戏过场

```
Frame [序号]: [场景描述],
Character emotion: [角色情感],
Action: [动作描述],
Camera: [镜头运动],
[序号+1] frame connects with smooth transition,
Frame [序号+1]: [下一帧场景],
[连接方式] transition: [淡入/切换/叠化],
continuing the narrative,
consistent character design, matching style,
storyboard sequence, [作品类型] --ar 16:9 --v 6
```

**应用示例 - 恐龙破壳**：
```
Frame 1: Giant dinosaur egg in jungle clearing, calm atmosphere,
Character emotion: anticipation, gentle glow from within,
Action: Egg begins to crack, small fissure appearing,
Camera: Close-up static shot,
Frame 2 connects with smooth transition,
Frame 2: Crack widens, tiny claw pushes through shell,
dissolve transition: glow intensifies,
continuing the narrative,
consistent dinosaur design, photorealistic style,
storyboard sequence, animated film --ar 16:9 --v 6
```

### 2.2 风格迁移/艺术模仿模板

适用于将主体转换为特定艺术风格

```
[主体描述] subject, in the style of [艺术家/流派],
[风格特征] artistic movement: [印象派/波普/水墨...],
[笔触特点] brushwork: [细腻/粗犷/点彩...],
[色彩倾向] color approach: [鲜艳/柔和/高对比...],
[质感] texture: [油画布/宣纸/漫画纸...],
[媒介感] medium: [油画/水彩/版画...],
art reproduction, [作品类型] art,
faithful to [艺术家名] technique,
museum quality, [尺寸] --ar 3:2 --v 6
```

**应用示例 - 梵高风格恐龙**：
```
T-Rex dinosaur portrait, in the style of Vincent van Gogh,
Impressionist artistic movement,
swirling brushwork patterns, visible brush strokes,
vibrant color palette: blues and yellows contrasting,
impasto texture, oil on canvas feel,
post-impressionism technique,
art reproduction, fine art painting,
faithful to Van Gogh's Starry Night energy,
museum quality, 24x18 inches --ar 3:2 --v 6
```

### 2.3 批量变体生成模板

适用于生成同一主体的多种变体（配色/角度/表情）

```
[主体], base design for variant generation,
[变体类型1]: [第一变体描述],
[变体类型2]: [第二变体描述],
[变体类型3]: [第三变体描述],
consistent core design, customizable elements,
color palette A/B/C options,
angle variants: front/3quarter/side,
expression variants: neutral/angry/happy,
turnaround sheet, orthographic views,
[用途] for [游戏/动画/商品] use,
modular design, easy to modify,
professional concept art, design sheet --ar 16:9 --v 6
```

**应用示例 - 恐龙IP设计变体**：
```
Friendly dinosaur mascot, base design for variant generation,
Variant A: Baby form, round features, big eyes, curious expression,
Variant B: Adult form, streamlined body, confident stance,
Variant C: armored variant, protective gear, battle-ready,
consistent core design, customizable color zones,
three color palette options: Jungle Green, Desert Orange, Volcanic Red,
angle variants: front view, 3/4 view, side profile,
expression variants: neutral, happy, surprised,
turnaround sheet, orthographic views,
for mobile game use, character asset design,
modular design, easy to modify,
professional concept art, character design sheet --ar 16:9 --v 6
```

---

## 三、技巧总结

### 3.1 本日最佳实践

1. **情感表达**：用"micro-expression"引导微表情，比单纯"happy/sad"更生动
2. **科技元素**：具体到"emission color #XXXXXX"和"holographic UI"才能保证细节
3. **连续叙事**：明确标注帧号和转场方式，保证生成连贯
4. **风格迁移**：加入"faithful to [artist]"比单纯描述风格更准确

### 3.2 验证方法

- 用 --upbeta 或 --no watermark 检查细节保留度
- 批量生成时固定 seed 测试一致性
- 风格迁移用低分辨率先验验证风格捕捉

---

*任务完成时间：2026-03-09 02:08 UTC*

---

## 今日cron执行确认

✅ **状态**：已完成  
⏰ **执行时间**：2026-03-09 09:12 UTC  
📋 **内容**：提示词优化(表情/机械)、3个新模板(故事板/风格迁移/批量变体)

---

## 2026-03-09 扩展（14:17 UTC）

### 🎯 主题：动态效果表现 & 材质细节强化 & 视角创意

---

### 🆕 动态效果表现模板

#### 运动模糊/速度线效果
```
[主体] [运动类型] with motion blur effect,
[速度感] speed: [数值] km/h feeling,
[运动方向] motion direction: [左→右/上→下/旋转],
[模糊强度] blur intensity: [轻微/中等/极致],
[速度线] speed lines: [数量] [颜色] [长度],
[残影] motion trail: [长度] [透明度],
[环境互动] environment interaction: [扬尘/气流/水花],
[焦点清晰] focal point sharp, action frozen,
[视角] perspective: [低角度/高角度/跟随],
[帧率模拟] cinematic frame: [24fps/60fps] look,
[后期感] post-processing: [电影感调色],
[风格] style: [写实/动漫/抽象],
dynamic shot, speed action,
frozen moment, professional photography --ar 16:9 --v 6
```

**应用示例 - 恐龙冲刺**：
```
T-Rex dinosaur sprinting with motion blur effect,
speed: 50km/h raw power feeling,
motion direction: left to right across frame,
blur intensity: extreme, legs blurred,
speed lines: 20 orange-gold lines, long trails,
motion trail: dust cloud behind, 3 body-length trail,
environment interaction: dirt flying, leaves swirling,
focal point sharp: head and upper body crystal clear,
perspective: low angle tracking shot,
cinematic frame: 24fps action movie look,
post-processing: anamorphic flare, orange-teal grade,
style: photorealistic,
dynamic dinosaur shot, speed action,
frozen moment, cinematic wildlife photography --ar 16:9 --v 6
```

#### 粒子系统/灰尘烟雾
```
[场景] with [粒子类型] particle system,
[粒子数量] density: [稀疏/中等/密集],
[粒子大小] particle size: [微米级/毫米级/厘米级],
[颜色渐变] color gradient: [起始色] → [中间色] → [消散色],
[运动模式] motion: [飘动/升腾/扩散/旋涡],
[光源互动] light interaction: [透射/反射/自发光],
[层次感] depth layers: [前景/中景/远景] 各[数量]粒子,
[动态模糊] particle motion blur: [开/关],
[混合模式] blend mode: [ additive / screen / normal],
[引擎预设] engine preset: Unity Particle System / Unreal Niagara,
[性能考量] performance: [轻量/中等/重度] GPU负载,
atmospheric particles, environmental VFX,
game-ready particles, cinematic density --ar 16:9 --v 6
```

---

### 🆕 材质细节强化模板

#### 金属质感强化
```
[物体] with [类型] metal material,
[金属类型] metal type: [铬/铜/钢/金/青铜/氧化],
[表面处理] surface finish: [抛光/拉丝/磨砂/做旧],
[反射] reflections: [镜面/模糊/环境贴图],
[氧化层] patina/oxidation: [程度] [颜色],
[划痕] scratches: [密度] [深度] [方向],
[指纹/使用痕迹] wear marks: [位置] [明显程度],
[环境因素] environmental exposure: [海边/工业/室内],
[微距细节] macro details: [纹理] [瑕疵],
[光线角度] lighting angle: [主光位置],
[质感关键词] texture keywords: brushed, weathered, authentic,
PBR material study, 8K resolution,
realistic metal, industrial texture --ar 4:3 --v 6
```

#### 有机材质/皮肤纹理
```
[生物/部位] with organic material detail,
[材质类型] material: [皮肤/鳞片/角质/软骨],
[纹理细节] texture: [毛孔/鳞片排列/裂纹/褶皱],
[颜色变化] color variation: [基色] with [斑点/渐变/深浅],
[次表面散射] subsurface scattering: [程度] [颜色透出],
[湿度/油脂] moisture/oil: [干燥/适中/湿润],
[血管/肌理] underlying details: [血管可见/肌肉线条/骨骼形状],
[微距特写] macro shot: close-up detail,
[光源] lighting: [柔光箱/自然光/背光],
[放大倍数] magnification: [1x/5x/10x] macro,
[风格] style: [科学准确/风格化/特写艺术],
biological material study, anatomical detail,
organic texture, skin/scale reference --ar 4:3 --v 6
```

---

### 🆕 视角创意/特殊镜头模板

#### 昆虫视角/微距摄影
```
[常见主体] from insect's eye view,
[视角类型] perspective: extreme macro, bug's vision,
[焦距] focal length: [数值]mm macro lens,
[景深] depth of field: [极浅/浅/中] focus stack,
[透视变形] perspective distortion: extreme close-up,
[清晰范围] sharp area: [精确到毫米] focus zone,
[背景虚化] background: [完全虚化/抽象色块],
[光线] light: [逆光/侧光] creating rim light,
[比例感] scale: [主体] appears [巨大/正常/微型],
[细节暴露] exposed details: normally hidden textures visible,
[视觉冲击] visual impact: alien perspective, new viewpoint,
macro photography, insect scale,
extreme close-up, scientific visualization --ar 4:3 --v 6
```

**应用示例 - 微观恐龙鳞片**：
```
T-Rex dinosaur skin from insect's eye view,
perspective: extreme macro, reptilian scale close-up,
focal length: 100mm macro lens,
depth of field: extremely shallow, focus stack required,
perspective distortion:鳞片表面起伏清晰可见,
sharp area: single scale edge in focus, 2mm zone,
background: completely blurred, abstract green tones,
light: backlight creating translucent scale edges,
scale: normal scale appears monumental under lens,
details exposed: individual scale edges, texture gaps, color variation,
visual impact: alien perspective, prehistoric texture revealed,
macro dinosaur skin, reptile scale study,
extreme close-up, paleontological visualization --ar 4:3 --v 6
```

#### 航拍/卫星视角
```
[地形/场景] from aerial/satellite view,
[高度] altitude: [数值] meters/feet,
[视角类型] perspective: [航拍/卫星/无人机],
[覆盖范围] coverage: [宽域] x [距离],
[细节级别] detail level: [可辨识物体/轮廓/抽象],
[光线角度] sun angle: [时间] shadow length,
[季节/植被] season: [夏季绿/秋季黄/冬季白],
[云层] cloud cover: [无/稀疏/覆盖],
[标注感] annotation readiness: [可添加标签/纯视觉],
[用途] for [游戏小地图/概念图/参考],
[风格] style: [写实航拍/概念艺术/卫星影像],
aerial perspective, top-down view,
bird's eye, drone footage look,
landscape survey, terrain study --ar 16:9 --v 6
```

---

### 🆕 情绪板/Moodboard 模板

#### 氛围参考情绪板
```
[项目/场景] moodboard, [情绪关键词] emotional keywords,
[参考元素1] visual reference: [风格/色调/质感],
[参考元素2] color mood: [主色调] + [对比色],
[参考元素3] lighting mood: [光照风格] + [光质],
[参考元素4] composition: [构图法则],
[参考元素5] texture mood: [材质触感],
[参考元素6] era/style: [时代感],
[参考元素7] cultural reference: [文化元素],
[参考元素8] mood: [整体氛围描述],
[排版] layout: [网格/自由/分层],
[用途] for [概念设计/UI风格/品牌],
[格式] format: [竖版/横版/方形],
mood board, design reference,
visual direction, art direction --ar 3:2 --v 6
```

---

### 🆕 UI交互状态模板

#### 按钮/组件多状态
```
[组件名称] UI button with states,
[默认状态] default: [描述] style, [颜色],
[悬停状态] hover: [变化描述], glow effect,
[按下状态] pressed: [变化描述], scale down,
[禁用状态] disabled: [描述], reduced opacity,
[加载状态] loading: [动画描述], spinner/progress,
[选中状态] selected: [描述], highlight border,
[焦点状态] focus: [描述], accessibility ring,
[图标] icon: [图标描述] left-aligned,
[文字] text: [按钮文案],
[圆角] corner radius: [数值]px,
[尺寸] size: [宽度]x[高度]px,
[风格] style: [扁平/拟物/Neumorphism],
game UI component, interactive button,
user interface design, HUD element --ar 4:3 --v 6
```

---

### 📊 本轮扩展统计

| 类型 | 数量 |
|------|------|
| 动态效果 | 2 (运动模糊/粒子) |
| 材质细节 | 2 (金属/有机) |
| 视角创意 | 2 (微距/航拍) |
| 情绪板 | 1 |
| UI交互状态 | 1 |
| **本轮总计** | **8** |

---

### 📋 技巧沉淀 - 动态效果

- [ ] 速度感与模糊强度匹配
- [ ] 粒子层次分明不杂乱
- [ ] 焦点清晰背景模糊
- [ ] 运动方向明确

---

### 📋 技巧沉淀 - 材质细节

- [ ] 金属反射与环境匹配
- [ ] 划痕方向统一
- [ ] 有机材质避免过度平滑
- [ ] 微距展示肉眼不可见细节

---

### 📋 技巧沉淀 - 视角创意

- [ ] 新视角带来视觉冲击
- [ ] 航拍比例感合理
- [ ] 微距清晰范围精确控制
- [ ] 情绪板元素协调统一

---

## 三、晚间扩展优化 (22:17 UTC)

### 3.1 毛发/羽毛精细化描写

**优化点**：毛发质感层次、羽毛结构、光效散射

**优化后提示词**：
```
[主体] with [毛发/羽毛类型] fur/plumage,
[密度] density: [稀疏/浓密/蓬松],
[长度] length: [短绒/长毛/羽毛层],
[颜色渐变] color gradient: [基色到末梢],
[光泽] sheen: [丝光/哑光/金属光],
[质感] texture: [丝滑/卷曲/粗糙],
individual hair strands visible,
fur detail maps, pore-level detail,
subsurface scattering on light,
[环境光影响] ambient occlusion in fur,
detailed wildlife photography,
8K, macro lens detail --ar 3:4 --v 6
```

**应用示例 - 羽毛恐龙**：
```
Velociraptor with iridescent plumage,
feather density: medium, sleek contour,
feather length: layered plumage on arms,
color gradient: bronze to emerald green sheen,
metallic sheen, smooth texture with barbs visible,
individual feather strands, detail texture,
subsurface light scattering, warm backlight,
feather depth, prehistoric bird study,
BBC wildlife documentary quality,
8K, macro lens detail --ar 3:4 --v 6
```

### 3.2 水体/液体动态效果

**优化点**：流体形态、波纹折射、泡沫悬浮

**优化后提示词**：
```
[主体] with [水体类型] water effect,
[状态] state: [静止/流动/激流/海浪],
[透明度] transparency: [清澈/浑浊/玻璃感],
[反射] reflections: [镜面/漫反射/焦散],
[折射] refraction: [扭曲程度],
[泡沫/水滴] foam/droplets: [密度/大小],
[运动模糊] motion blur on splashes,
Caustics投影, volumetric water volume,
rain/splash particles, wet surface shader,
fluid simulation, Houdini/Nuke ready,
photorealistic liquid, 8K --ar 16:9 --v 6
```

**应用示例 - 恐龙饮水**：
```
T-Rex drinking from river, calm water surface,
state: gentle ripple, transparency: clear with blue tint,
soft reflections of dinosaur form,
subtle refraction of submerged rocks,
small foam droplets around mouth, 
gentle motion blur on water surface,
caustics patterns on riverbed,
volumetric water, splash particles,
wet fur shader on muzzle,
fluid simulation, VFX ready,
photorealistic liquid, 8K --ar 16:9 --v 6
```

### 3.3 破坏/破碎效果

**优化点**：碎片轨迹、裂纹走向、粉尘云

**优化后提示词**：
```
[主体] with [破坏类型] destruction,
[破坏程度] intensity: [轻微/中等/粉碎],
[碎片类型] debris: [碎片/石块/玻璃/金属],
[裂纹样式] crack pattern: [放射状/龟裂/粉碎性],
[灰尘/粉尘] dust/cloud: [密度/颜色],
[碎片运动] debris trajectory: [飞散方向/速度感],
[受损材质] broken material: [断面细节],
rubble pile, shattered fragments,
explosion/impact point, shockwave,
blast wave distortion, smoke volume,
detailed destruction physics,
VFX compositing elements --ar 16:9 --v 6
```

**应用示例 - 恐龙打碎玻璃**：
```
T-Rex smashing through glass wall,
destruction intensity: complete shatter,
debris: shattered glass panels, sparkling fragments,
crack pattern: radial shatter from impact point,
white dust cloud, medium density,
glass fragments flying outward, motion trails,
broken glass edges, sharp facets,
rubble pile at base, impact sparks,
shockwave distorting air, glass dust particles,
detailed destruction physics,
compositing elements included --ar 16:9 --v 6
```

---

## 四、新模板产出（晚间）

### 4.1 产品展示/电商模板

适用于商品图、海报、电商主图

```
[主体] product shot, [视角] viewing angle,
[尺寸] product size: [实际尺寸],
[材质] material: [材质描述],
[场景] setting: [使用场景/纯色背景],
[光照] lighting: [主光/辅光/轮廓光],
[焦点] focus: [景深控制],
[用途] for [电商/海报/目录] use,
[风格] style: [极简/奢华/科技/自然],
studio product photography,
professional retouching,
high resolution, [分辨率] --ar 1:1 --v 6
```

**应用示例 - 恐龙玩具**：
```
Collectible dinosaur toy product shot,
45-degree angle, medium distance,
size: 12 inches tall, poseable joints,
realistic dinosaur skin texture,
white studio background, infinite cove,
three-point lighting: key left, fill right, rim back,
shallow depth of field, focus on head detail,
for e-commerce and catalog use,
clean minimalist style, professional product photography,
professional retouching, 4K --ar 1:1 --v 6
```

### 4.2 科技UI/界面设计模板

适用于App界面、仪表盘、HUD显示

```
[UI类型] interface design,
[设备] device: [手机/平板/车载/AR眼镜],
[主题] theme: [暗黑/光明/赛博/极简],
[颜色] color scheme: [主色/辅色/强调色],
[布局] layout: [网格/瀑布/卡片],
[元素] elements: [按钮/图表/地图/状态栏],
[hud元素] HUD elements: [数据/读数/指示],
glassmorphism / neomorphism / flat design,
responsive layout, accessibility compliant,
pixel perfect, [分辨率] design,
Figma/Sketch export ready --ar 9:16 --v 6
```

**应用示例 - 恐龙追踪App**：
```
Mobile app interface design,
device: smartphone, screen size 6.1",
theme: dark nature, Forest green #228B22 accents,
color scheme: dark gray #1A1A1A base, green accents, white text,
layout: bottom navigation, card-based content,
elements: search bar, dinosaur icons, map view, filter chips,
HUD elements: GPS coordinates, distance meter, species identification,
glassmorphism cards with blur, responsive layout,
accessibility: WCAG AA compliant,
pixel perfect, 1170x2532 design,
Figma export ready --ar 9:16 --v 6
```

### 4.3 美食/餐饮摄影模板

适用于菜单、食谱、餐饮广告

```
[菜品描述] dish description,
[烹饪风格] cuisine: [中餐/西餐/日料/融合],
[呈现方式] plating: [摆盘风格],
[餐具] plate/utensils: [材质/风格],
[garnish] garnish: [装饰],
[光泽] glaze/sauce: [光泽感],
[蒸汽] steam: [有无/浓度],
[背景] background: [场景/色彩],
[角度] shooting angle: [俯拍/45度/侧拍],
[氛围] mood: [温馨/高级/家庭/清爽],
food photography, professional styling,
editorial quality --ar 4:5 --v 6
```

**应用示例 - 恐龙形蛋糕**：
```
Dinosaur-shaped birthday cake,
cuisine: modern American fusion,
plating: on wooden board, dinosaur silhouette carved,
plate: white ceramic, minimal,
garnish: edible grass, small fondant flowers,
chocolate glaze with drip effect, shiny,
warm steam rising, light mist,
background: soft gray gradient,
shooting angle: 45-degree hero shot,
mood: celebratory, joyful, party atmosphere,
food photography, professional styling,
Bon Appétit editorial quality --ar 4:5 --v 6
```

### 4.4 汽车/载具设计模板

适用于车辆设计、概念图、宣传图

```
[车型] vehicle type, [品牌风格] brand aesthetic,
[年代] era: [复古/现代/未来],
[定位] segment: [家用/豪华/运动/商用],
[动力] powertrain: [燃油/电动/氢能源/混合动力],
[视角] view: [正面/侧面/45度/俯视/透视],
[姿态] stance: [低趴/高底盘/运动悬架],
[细节] details: [车灯/轮毂/格栅/后视镜],
[环境] setting: [赛道/城市/自然/studio],
[氛围] atmosphere: [激烈/优雅/未来感],
automotive photography, car design,
studio lighting, [渲染引擎] render --ar 16:9 --v 6
```

**应用示例 - 恐龙牵引车**：
```
Heavy-duty dinosaur transport vehicle,
brand aesthetic: rugged industrial,
era: near-future 2085,
segment: commercial utility,
powertrain: hybrid diesel-electric,
view: three-quarter front dynamic,
stance: high clearance, heavy-duty suspension,
details: LED light bar, reinforced cage, mud flaps,
setting: desert research facility, dust environment,
atmosphere: powerful, functional, scientific,
automotive photography, professional vehicle design,
studio lighting, Octane render --ar 16:9 --v 6
```

---

*扩展记录时间：2026-03-09 22:17 UTC*
*执行者：art-gen-daily cron (晚间扩展)*

---

## 📈 累计模板产出（今日更新）

| 时间段 | 优化方向 | 产出模板数 |
|--------|----------|------------|
| 09:12 | 表情/机械/故事板/风格迁移/批量变体 | 5 |
| 14:17 | 动态效果/材质细节/视角/情绪板/UI状态 | 8 |
| 22:17 | 毛发水体/破坏效果/产品展示/UI/美食/载具 | 10 |
| **今日总计** | | **23** |

---

## ✅ 今日任务完成清单

- [x] 提示词优化（动态效果/材质/视角）
- [x] 产出新模板（8个专项模板）
- [x] 晚间扩展优化（毛发水体/破坏效果）
- [x] 晚间新模板产出（4个专项模板）
- [x] 记录到 memory/art-gen/daily-work-2026-03-09.md

---

*本记录由 AI 自动生成 | 持续更新中*
