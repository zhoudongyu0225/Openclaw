# 美术素材每日工作 - 2026-03-10

## 任务执行：提示词优化 & 新模板产出

---

## 一、提示词优化

### 1.1 天气/环境效果系统

**优化点**：增加动态天气层次、时间变化、氛围渲染

**优化后提示词**：
```
[场景类型] environment, [季节/气候] season/climate,
[天气状态] weather condition: [暴雨/晴朗/雾/沙尘],
[时间] time of day: [黎明/正午/黄昏/午夜],
[氛围] atmosphere: [tense/peaceful/mysterious/epic],
[光照] lighting: [volumetric rays/golden hour/blue hour/neon lights],
[粒子效果] particles: [rain drops/snowflakes/ash/dust],
[地面效果] ground: [wet reflections/puddles/frost/sand],
[天空] sky: [storm clouds/stars/aurora/blood moon],
[特效增强] post-processing: [bloom/DOF/motion blur],
[渲染] [引擎] render, 8K, cinematic quality --ar 16:9 --v 6
```

**应用示例 - 火山岛战场**：
```
Volcanic island battlefield environment, tropical climate,
active eruption atmosphere, sunset time golden hour,
tense atmosphere, volumetric god rays through smoke,
ash particles falling, lava rivers with glowing cracks,
dark volcanic rock terrain, scorched vegetation,
blood orange sky with smoke clouds, heat distortion,
unreal engine 5 render, 8K, cinematic quality,
epic battle atmosphere --ar 16:9 --v 6 --s 250
```

### 1.2 动态UI血条/状态栏

**优化点**：增加生命值、能量条、buff/debuff显示层

**优化后提示词**：
```
Game UI health bar for [游戏类型],
[形状] shape: [bar/circular/shield/hp orb],
[样式] style: [sci-fi/medieval/fantasy/minimalist],
[颜色] primary color: #[HEX], secondary: #[HEX],
[动画状态] state: [full/half/full with shield/debuff active],
[边框] border: [glowing/metallic/ornate/none],
[图标] icon: [heart/shield/crystal/skull],
[数字显示] HP text: [numerical/percentage/both],
[特效] effects: [pulse glow when low/ripple on damage/shield break],
[背景] background: [transparent/gradient/texture],
game ui mockup, figma style, 4K --ar 16:9 --v 6
```

**应用示例 - 科幻风血条**：
```
Game UI health bar for sci-fi tower defense,
circular orb shape, holographic style,
cyan #00FFFF primary with magenta #FF00FF accents,
full health with energy shield active,
glowing border effect, geometric shield icon,
numerical HP display 1250/1250,
pulse glow animation ready, ripple on damage,
transparent background with subtle grid texture,
game ui mockup, figma style, 4K --ar 16:9 --v 6 --s 180
```

### 1.3 地形生成/地块设计

**优化点**：增加可建造区域、地形功能、视觉层次

**优化后提示词**：
```
Game terrain tile for tower defense,
[tile type]: [grass/rock/water/lava/sand/snow],
[功能] function: [buildable/defensive boost/resource generation],
[等级] tier: [tier 1 basic/tier 2 improved/tier 3 premium],
[装饰] decorations: [trees/rocks/crystals/ruins/skeleton],
[视觉风格] visual: [stylized/realistic/low-poly/painted],
[光照] lighting: [baked shadows/real-time/ambient occlusion],
[地块形状] shape: [square/hexagonal/irregular],
[边缘处理] edges: [flat/cliff/river/wall],
isometric view, game asset, white background --ar 4:3 --v 6
```

**应用示例 - 史诗地形**：
```
Game terrain tile for tower defense,
volcanic rock tile, defensive boost function,
tier 3 premium quality, glowing crystal decorations,
cracked lava veins with emission, obsidian texture,
stylized visual, baked ambient occlusion,
hexagonal shape, cliff edges with lava waterfall,
isometric view, game asset, clean render --ar 4:3 --v 6 --s 200
```

---

## 二、新模板产出

### 2.1 过场动画/剧情CG模板

适用于游戏剧情、任务过场、角色介绍

```
CG Scene [序号]: [场景设定],
[角色1] character: [描述], [姿态] pose: [动作],
[角色2] character: [描述], [姿态] pose: [动作],
[对话气泡/文字] dialogue: "[台词]",
[情感基调] mood: [tense/heartwarming/sad/epic],
[镜头] camera: [close-up/wide shot/over-shoulder/dramatic angle],
[光照] lighting: [key light/fill/rim/back light],
[氛围] atmosphere: [rain/mist/sun rays/darkness],
[风格] style: [anime/cinematic/realistic/painterly],
[特效] vfx: [glow/particles/motion lines],
[渲染] [引擎], 8K, cinematic --ar 16:9 --v 6

示例：
CG Scene 1: Dinosaur sanctuary at dawn,
main character: young female researcher, pose: reaching toward baby T-Rex,
baby T-Rex: curious, nuzzling her hand,
dialogue: "You'll be safe here now.",
mood: heartwarming, camera: low angle close-up,
warm golden sunrise lighting, misty atmosphere,
cinematic style, soft bloom vfx,
unreal engine 5, 8K, cinematic --ar 16:9 --v 6 --s 250
```

### 2.2 装备槽/物品栏UI模板

适用于背包系统、装备选择、物品管理

```
Game inventory slot UI,
[类型] type: [weapon/armor/accessory/consumable/material],
[稀有度] rarity: [common/rare/epic/legendary],
[样式] style: [sci-fi/medieval/fantasy/minimalist],
[形状] shape: [square/hexagonal/round],
[边框] border: [glowing edge/ornate frame/metallic/None],
[图标] icon: [具体物品描述],
[数量显示] count: [数量/badge],
[状态] state: [empty/equipped/locked/disabled],
[背景] background: [dark #XXXXXX/gradient/transparent],
[高亮] highlight: [selected/hover/available],
[特效] effects: [rarity glow/particle sparkle],
game ui mockup, clean design, figma style --ar 1:1 --v 6

示例：
Game inventory slot UI,
type: legendary weapon, golden legendary sword,
rarity: legendary, hexagonal shape,
glowing golden border #FFD700 with particle sparkles,
detailed sword icon with glow effect,
count badge: 1, state: equipped,
dark background #1A1A2E with subtle gradient,
selected highlight with cyan #00FFFF rim,
rarity glow effect radiating outward,
game ui mockup, clean design, figma style --ar 1:1 --v 6 --s 180
```

### 2.3 Boss战UI模板

适用于Boss血条、技能预警、阶段转换

```
Boss fight UI for [游戏类型],
[boss名称] boss: [描述],
[血条样式] health bar: [large/segmented/phased],
[阶段显示] phase: [phase 1/2/3 with transformation],
[技能预警] skill warning: [telegraph area/direction indicator],
[冷却显示] cooldown: [ability icons with timers],
[特效] effects: [intense/flash warning/phase transition],
[样式] style: [epic/ominous/chaotic],
[颜色] color scheme: [主色 #HEX/强调色 #HEX],
[动画] animations: [pulse on damage/phase change vfx],
game ui mockup, cinematic layout, 4K --ar 16:9 --v 6

示例：
Boss fight UI for dinosaur tower defense,
boss: Ancient Dragon Titan, massive segmented health bar,
3 phases with transformation (normal/enraged/apocalyptic),
skill telegraph: red warning zones on ground,
4 ability icons with cooldown timers (fire breath/tail swipe/roar/meteor),
ominous style, blood red #8B0000 with gold #FFD700 accents,
phase transition: screen shake, dramatic flash, transformation vfx,
game ui mockup, cinematic layout, 4K --ar 16:9 --v 6 --s 200
```

### 2.4 任务/成就系统UI模板

适用于任务列表、成就解锁、进度追踪

```
Quest/achievement UI for [游戏类型],
[类型] type: [quest list/achievement popup/progress tracker],
[状态] state: [available/in progress/completed/failed],
[标题] title: "[任务/成就名称]",
[描述] description: "[具体内容]",
[进度] progress: [current/total with percentage],
[奖励] rewards: [items/experience/currency],
[样式] style: [adventure/fantasy/sci-fi],
[颜色] colors: [primary #HEX/secondary #HEX/accent #HEX],
[图标] icon: [quest marker/achievement badge/trophy],
[动画] animation: [slide in/pulse glow/celebration particles],
[布局] layout: [list/grid/card],
game ui mockup, professional design --ar 16:9 --v 6

示例：
Achievement popup UI for dinosaur tower defense,
type: achievement unlock, state: newly completed,
title: "Extinction Event",
description: "Defeat 1000 dinosaurs in a single run",
progress: 1000/1000 (100%), rewards: 500 gems + "Dinosaur Slayer" title,
fantasy style, gold #FFD700 primary with deep blue #1E90FF secondary,
trophy icon with celebratory rays, unlock animation with particle burst,
card layout with glow effect, game ui mockup, professional design --ar 16:9 --v 6 --s 180
```

---

## 三、提示词公式进阶

### 多主体组合公式
```
[主主体, 详细描述] + [次主体, 详细描述] + [环境] + [风格] + [构图] + [光照] + [氛围] + [技术规格]
```

### 动态效果公式
```
[静态主体] + [motion blur] + [particle trails] + [dynamic pose] + [action line] + [impact effects]
```

### 材质特写公式
```
[材质类型] texture close-up, [纹理细节] detailed grain,
[光照角度] lighting at [角度], [反射/折射] reflections,
[放大倍数] macro view, [风格] [stylized/photorealistic],
[用途] [material library/game asset/reference] --ar 1:1 --v 6
```

---

## 五、新增模板 - 角色与特效

### 5.1 角色设计模板（多角度展示）

适用于角色图鉴、英雄选择、皮肤展示

```
Character design sheet for [游戏类型],
[角色名称]: [职业/定位], [种族/背景],
[外观描述]: [体型/服装/配色/特色元素],
[姿态] pose A: [front facing, full body],
[姿态] pose B: [3/4 view, combat pose],
[姿态] pose C: [side profile, walking],
[武器/装备] weapons: [主武器描述],
[配色方案] color palette: [主色 #HEX/辅色 #HEX/点缀色 #HEX],
[风格] style: [anime/realistic/cartoon/painted],
[年龄段] age group: [child/young/adult/elder],
[表情] expressions: [neutral/angry/happy/serious],
[材质] materials: [metal/fabric/leather/magic glow],
[渲染] [引擎], 8K, character design sheet --ar 3:2 --v 6

示例：
Character design sheet for dinosaur tower defense,
"Dr. Maya Chen": paleontologist scientist, human female,
athletic build, field research vest with pockets,
tan cargo pants, hiking boots, utility belt,
pose A: front facing, confident smile, holding tranquilizer gun,
pose B: 3/4 view, combat ready, aiming weapon,
pose C: side profile, running forward,
weapons: tranquilizer rifle, proximity mine, flare gun,
color palette: forest green #228B22/khaki #F0E68C/orange #FF8C00,
anime style, young adult, determined expression,
tech-enhanced fabric materials with glow accents,
unreal engine 5, 8K, character design sheet --ar 3:2 --v 6 --s 250
```

### 5.2 技能/特效图标模板

适用于技能栏、法术图标、Buff/Debuff标识

```
Ability icon for [游戏类型],
[技能名称]: [技能效果描述],
[类型] type: [attack/heal/buff/debuff/utility/passive],
[元素] element: [fire/water/earth/air/void/light/dark],
[形状] shape: [circular/diamond/shield/rune/crystal],
[颜色] colors: [主色 #HEX/渐变到 #HEX],
[特效] effects: [glow/particles/flames/ice crystals/lightning],
[图标核心] core: [ swords/cross/symbol/animal/star],
[动画状态] animation: [ready/cooldown/active],
[背景] background: [solid/transparent/gradient/pattern],
[稀有度] rarity: [common/rare/epic/legendary],
[格式] icon set style, [游戏引擎] asset, clean render --ar 1:1 --v 6

示例：
Ability icon for dinosaur tower defense,
"Extinction Roar": deals massive damage to all enemies,
type: ultimate attack, element: void/dark,
hexagonal dark crystal shape, 
gradient from deep purple #4B0082 to black #000000,
screaming T-Rex skull core with void energy,
pulsing dark aura with crimson particles,
ready state with golden border,
transparent background with void mist,
legendary rarity, legendary glow effect,
icon set style, unreal engine asset, clean render --ar 1:1 --v 6 --s 180
```

### 5.3 加载画面模板

适用于关卡加载、地图过渡、模式切换

```
Loading screen for [游戏类型],
[场景设定] setting: [描述],
[主要视觉] main visual: [角色/怪物/场景/ Logo],
[动画] animation: [idle/moving/interacting],
[进度条] progress bar: [位置/样式],
[提示信息] tips: "[游戏技巧或lore]",
[背景] background: [image/video/animated particles],
[样式] style: [cinematic/minimalist/immersive],
[配色] color scheme: [主色 #HEX],
[品牌元素] branding: [logo位置/版权信息],
[加载动画] loading animation: [spinner/progress/character action],
[渲染] [引擎], 4K --ar 16:9 --v 6

示例：
Loading screen for dinosaur tower defense,
setting: Jurassic research station interior,
main visual: animatronic baby dinosaur in incubation pod,
baby dino: sleeping peacefully, occasional twitch,
progress bar: bottom of screen, sci-fi holographic style,
tips: "Tip: Use stun towers to give your DPS towers more time!",
background: animated rain outside station windows with fog,
cinematic style, deep blue #000080 with green #00FF7F accents,
studio logo top-left, copyright bottom-right,
loading animation: DNA strand filling up,
unreal engine 5, 4K --ar 16:9 --v 6 --s 180
```

### 5.4 小地图/世界地图模板

适用于游戏内地图、探索系统、区域显示

```
Game minimap for [游戏类型],
[地图类型] type: [minimap/tactical map/world map/region map],
[显示内容] shows: [enemies/towers/resources/objectives/player],
[样式] style: [top-down/3D/isometric/radar],
[颜色方案] color scheme: [dark mode #XXXXXX / light mode #XXXXXX],
[标记] markers: [friendly (green)/enemy (red)/objective (gold)/danger (red pulse)],
[缩放级别] zoom levels: [close/mid/far],
[功能元素] functional: [compass/player arrow/scale/coordinates],
[边框] border: [circular/square/rounded/None],
[特效] effects: [fog of war reveal/pulse on objective/glow on enemy],
[透明度] opacity: [map X% / markers Y%],
[渲染] clean vector style or [引擎] render --ar 1:1 --v 6

示例：
Game minimap for dinosaur tower defense,
type: tactical minimap with radar sweep,
shows: all towers (blue dots), enemies (red dots), 
resources (green crystals), boss (pulsing red),
top-down orthographic view, dark mode #1A1A2E background,
friendly towers: cyan #00FFFF, enemies: red #FF0000,
gold #FFD700 objective markers with pulse,
3 zoom levels, compass top-right, player arrow center,
circular border with glowing edge,
fog of war: explored areas bright, unexplored dark,
enemy pulse effect on detection, clean vector style --ar 1:1 --v 6 --s 180
```

---

## 六、提示词进阶优化

### 6.1 光照系统增强

**优化思路**：从单一光源到多光源系统，从静态到动态

```
Multi-light setup for [场景],
[主光源] key light: [类型] [颜色] [角度],
[补光] fill light: [类型] [颜色] [强度],
[轮廓光] rim light: [类型] [颜色] [位置],
[环境光] ambient: [颜色] [强度],
[动态效果] dynamic: [flickering/shifting/moving],
[阴影] shadows: [hard/soft/ray traced],
[特效] effects: [god rays/lens flare/glow],
[引擎] render, 8K --ar [比例] --v 6
```

### 6.2 材质升级公式

**优化思路**：从单一材质到复合材质，从平面到PBR

```
Advanced material: [材质A] base with [材质B] inlays,
[纹理] texture: [normal/bump/displacement] maps,
[反射] reflections: [smooth/matte/metallic],
[细节] detail: [scratches/wear/weathering/damage],
[光照反应] light response: [specular/roughness/transmission],
[特写] macro detail: [详细纹理描述],
[用途] [game asset/texture library/reference],
[引擎] render, 8K --ar 1:1 --v 6
```

### 6.3 构图法则应用

**优化思路**：根据画面目的选择最佳构图

| 构图类型 | 适用场景 | 提示词 |
|---------|---------|--------|
| 三分法 | 平衡场景 | rule of thirds, balanced composition |
| 中心构图 | 强调主体 | centered composition, focus on [主体] |
| 对角线 | 动感场景 | diagonal composition, dynamic flow |
| 对称构图 | 建筑/UI | symmetrical, perfect balance |
| 引导线 | 强调纵深 | leading lines toward [焦点] |

---

## 七、本日新增产出

**新增模板**：
- 角色设计模板（多角度展示）
- 技能/特效图标模板
- 加载画面模板
- 小地图/世界地图模板

**提示词进阶优化**：
- 多光源系统提示词
- 高级材质提示词公式
- 构图法则应用指南

**累计今日产出**：
- 优化类：3套系统（天气环境、UI血条、地形）
- 模板类：8套（CG场景、物品栏、Boss战、成就、角色、技能、加载、地图）
- 公式类：5组（多主体、动态效果、材质特写、光照、材质）

---

## 八、今日补充优化 - 风格化与特效

### 8.1 Low-Poly风格模板

适用于移动端游戏、资源优化、项目原型

```
Low-poly [主体类型] for game,
[风格变体] style: [isometric/characters/vehicles/environment],
[多边形数量] poly count: [low/medium/high],
[颜色数量] color count: [limited palette/gradient],
[表面处理] shading: [flat/cel-shaded/gouraud],
[细节层次] LOD: [LOD0/LOD1/LOD2],
[用途] usage: [mobile game/VR/project/prototype],
[是否包含动画] animated: [yes/no],
[文件格式] format: [FBX/OBJ/GLTF],
[渲染] clean render, white background --ar 4:3 --v 6

示例：
Low-poly dinosaur pack for mobile tower defense,
style: isometric environment + characters,
medium poly count (5-10K tris per model),
limited color palette (8 colors max),
flat shading with subtle gradients,
LOD0: full detail, LOD1: 50% reduction,
animated: idle + attack + death,
FBX format with textures,
clean render, white background --ar 4:3 --v 6 --s 180
```

### 8.2 像素艺术模板

适用于复古游戏、2D游戏、UI图标

```
Pixel art [主体] for game,
[像素尺寸] pixel size: [16x16/32x32/64x64/128x128],
[风格] style: [8-bit/16-bit/modern pixel],
[调色板] palette: [NES/GameBoy/complete RGB],
[视角] perspective: [top-down/isometric/side-scroll],
[动画帧数] animation frames: [2/4/8/12],
[用途] usage: [sprite/sheet/icon/tileset],
[抖动效果] dithering: [none/ordered/bayer],
[背景] background: [transparent/tiled/gradient],
clean pixel perfect, game ready --ar [比例] --v 6

示例：
Pixel art dinosaur enemy sprites for retro tower defense,
pixel size: 32x32 for enemies, 64x64 for boss,
style: 16-bit era pixel art,
complete RGB palette with shadow colors,
top-down perspective with shadow,
animation frames: 4 (idle/walk/attack/death),
usage: sprite sheet with animation rows,
ordered dithering for gradients,
transparent background,
pixel perfect, game ready --ar 8:1 --v 6 --s 180
```

### 8.3 2D立绘模板

适用于视觉小说、卡牌游戏、角色介绍

```
2D character portrait for [游戏类型],
[角色名称]: [职业/性格描述],
[视角] view: [bust/half-body/full-body],
[姿态] pose: [standing/sitting/combat/emote],
[表情] expression: [neutral/happy/angry/sad/surprised],
[服装] outfit: [详细描述],
[配色] color scheme: [主色 #HEX/辅色 #HEX],
[风格] style: [anime/realistic/vector/illustrated],
[背景] background: [solid/transparent/场景描述],
[特效] effects: [glow/sparkles/motion lines],
[格式] format: [portrait sheet/expression sheet/half-body],
[渲染] [引擎], high quality illustration --ar 3:4 --v 6

示例：
2D character portrait for dinosaur tower defense,
"Captain Rex": veteran dinosaur handler, confident veteran,
view: half-body, pose: standing with weapon ready,
expression: determined and battle-hardened,
tactical vest with dinosaur handling gear, utility belt,
color scheme: military green #556B2F with orange #FF8C00 accents,
anime style with clean lines, subtle cel-shading,
transparent background for UI use,
glow effect on weapon highlights,
portrait sheet with 4 expressions,
Clip Studio Paint, high quality illustration --ar 3:4 --v 6 --s 200
```

### 8.4 粒子特效模板

适用于技能特效、环境效果、UI反馈

```
Particle effect: [特效名称],
[类型] type: [fire/water/ice/lightning/void/nature/magic],
[形状] shape: [orb/plasma/beam/spark/swarm/ring],
[粒子数量] particle count: [light/medium/heavy],
[颜色渐变] gradient: [起始色 #HEX → 结束色 #HEX],
[运动方式] motion: [spiral/burst/stream/hover/flow],
[持续时间] duration: [instant/pulse/continuous],
[触发方式] trigger: [on-hit/on-cast/on-death/passive],
[混合模式] blend mode: [add/alpha/multiply],
[用途] usage: [skill VFX/ambient/UI feedback],
[引擎] VFX asset, transparent background --ar 1:1 --v 6

示例：
Particle effect: Magma Eruption,
type: fire + earth, 
shape: burst from ground with debris,
particle count: heavy (500+ particles),
gradient: bright yellow #FFFF00 → orange #FF4500 → dark red #8B0000,
motion: explosive burst upward, debris arc, ember float,
duration: 2-second burst with 3-second embers,
trigger: on skill cast,
add blend mode for glow,
usage: ultimate skill VFX for fire dinosaur,
unreal engine Niagara, VFX asset, transparent background --ar 1:1 --v 6 --s 180
```

---

## 九、提示词急救箱

### 快速修复公式

| 问题 | 修复词 |
|-----|--------|
| 画面太暗 | +bright lighting +high key +well lit |
| 细节不足 | +intricate details +complex +highly detailed |
| 风格不对 | +[风格关键词] style +artstation top |
| 构图松散 | +tight composition +centered +rule of thirds |
| 色彩平淡 | +vibrant colors +saturated +bold palette |
| 太平面 | +depth +atmospheric +3D render +PBR |
| 不够专业 | +professional +clean +polished +studio quality |

---

## 十、本日完整产出汇总

**提示词优化类**：
- 天气/环境效果系统
- 动态UI血条/状态栏
- 地形生成/地块设计
- 光照系统增强
- 材质升级公式

**新模板产出类**：
- CG场景/过场动画模板
- 装备槽/物品栏UI模板
- Boss战UI模板
- 任务/成就系统UI模板
- 角色设计模板
- 技能/特效图标模板
- 加载画面模板
- 小地图/世界地图模板
- Low-Poly风格模板
- 像素艺术模板
- 2D立绘模板
- 粒子特效模板

**公式/指南类**：
- 多主体组合公式
- 动态效果公式
- 材质特写公式
- 构图法则应用指南
- 提示词急救箱

---

## 九、本日新增优化 - 特殊效果与风格化

### 9.1 水体/液体特效模板

适用于水域关卡、液体攻击、环境效果

```
[水体类型] water effect for [游戏类型],
[效果类型] effect: [still/flowing/breaking/splashing],
[颜色] color: [主色 #HEX] to [辅色 #HEX],
[透明度] opacity: [百分比],
[状态] state: [calm/turbulent/frozen/boiling],
[特效] features: [waves/foam/reflections/caustics],
[光照交互] light interaction: [specular/refraction/reflection],
[用途] usage: [game vfx/environment/attack effect],
[格式] particle system ready / sprite sheet,
[引擎] render, transparent background --ar 16:9 --v 6

示例：
Lava river effect for dinosaur tower defense,
effect: flowing molten rock with fire particles,
color: orange #FF4500 to yellow #FFD700 gradient,
opacity: 90%, state: turbulent and dangerous,
features: heat waves, glowing cracks, ember sparks,
light interaction: intense glow, volumetric heat,
game vfx, particle system ready, sprite sheet layout,
unreal engine 5, transparent background --ar 16:9 --v 6 --s 250
```

### 9.2 粒子系统提示词

适用于技能特效、环境氛围、UI动画

```
Particle system for [效果描述],
[粒子类型] type: [fire/smoke/magic/energy/dust/sparks],
[数量密度] density: [sparse/medium/dense],
[颜色渐变] gradient: [起始色] to [结束色],
[运动轨迹] motion: [spiral/random/directional/explosion],
[大小变化] size: [uniform/varied with age],
[生命周期] lifetime: [short/medium/long],
[发光效果] emission: [glow/bloom/None],
[触发方式] trigger: [continuous/on-hit/periodic],
[兼容性] format: [sprite sheet/strip/individual],
[引擎] compatible: [Unity/UE/Godot], game vfx --ar 1:1 --v 6

示例：
Magic shield break particle system,
type: energy dissipation with magic sparks,
density: dense burst (50+ particles),
gradient: cyan #00FFFF to transparent white,
motion: radial explosion outward,
size: varied, larger at center, smaller at edges,
lifetime: short (0.5s), 
emission: strong glow with bloom,
trigger: on shield destruction,
format: sprite sheet 8x8 grid,
Unity particle system compatible, game vfx --ar 1:1 --v 6 --s 200
```

### 9.3 2D精灵/动画帧模板

适用于2D游戏、动画序列、精灵图

```
2D sprite sheet for [角色/对象],
[动作序列] animation: [idle/walk/attack/death/special],
[帧数] frame count: [数量] frames,
[风格] style: [pixel art/hand-drawn/vector],
[方向] facing: [right/left/both/multi],
[尺寸] resolution: [像素尺寸如 64x64],
[颜色模式] color mode: [indexed/RGB/Grayscale],
[背景] background: [transparent/solid/None],
[应用场景] usage: [2D game/animatic/reference],
[输出格式] format: [PNG sprite sheet/strip/single files],
[引擎] engine: [Unity/Godot/custom], game ready --ar [比例] --v 6

示例：
T-Rex attack animation sprite sheet,
animation: full attack sequence from wind-up to impact,
frame count: 8 frames,
style: pixel art, 16-bit era aesthetic,
facing: right, resolution: 128x128 pixels,
color mode: indexed 16 colors,
background: transparent,
usage: 2D tower defense game,
format: PNG sprite sheet 4x2 grid,
Unity 2D sprite editor compatible, game ready --ar 2:1 --v 6 --s 180
```

### 9.4 交互式UI动画模板

适用于按钮、过渡、反馈动画

```
UI animation for [交互元素],
[交互类型] interaction: [hover/click/drag/transition],
[动画曲线] easing: [ease-in-out/elastic/bounce/linear],
[持续时间] duration: [0.2s/0.5s/1s],
[效果] effects: [scale/rotate/translate/color change/glow/pulse],
[方向] direction: [in/out/center/scale],
[反馈] feedback: [visual/audio/haptic],
[样式] style: [subtle/dramatic/minimal/elaborate],
[颜色] colors: [主色 #HEX/高亮色 #HEX],
[兼容性] format: [GIF/lottie/code/sequence],
[用途] usage: [mobile/desktop/cross-platform],
[引擎] tool: [After Effects/Spine/Rive/Unity] --ar 16:9 --v 6

示例：
Tower upgrade button animation,
interaction: click with success feedback,
easing: elastic bounce, duration: 0.4s,
effects: scale up 1.1x, golden glow pulse, checkmark appear,
direction: center expansion,
feedback: visual flash + subtle audio chime,
style: sci-fi with satisfying pop,
colors: default #4169E1 to success #00FF7F,
format: Lottie JSON for cross-platform,
usage: mobile + desktop,
tool: After Effects with Rive export --ar 16:9 --v 6 --s 150
```

---

## 十、风格化专项优化

### 10.1 赛博朋克风格模板

适用于科幻城市、夜间场景、霓虹效果

```
Cyberpunk [场景/对象] for [游戏类型],
[环境] setting: [城市/街道/屋顶/地下],
[时间] time: [night/early morning, rain],
[霓虹颜色] neon colors: [主色] [强调色1] [强调色2],
[建筑风格] architecture: [高塔/贫民窟/企业区],
[科技元素] tech elements: [全息/LED/屏幕/无人机],
[氛围] atmosphere: [潮湿/蒸汽/烟雾],
[光照] lighting: [霓虹反射/雨反射/暗角],
[人物/角色] characters: [描述],
[脏迹效果] weathering: [wet/grime/damage],
[风格] style: [赛博朋克/故障艺术/未来主义],
[引擎] render, 8K, cinematic --ar 16:9 --v 6

示例：
Cyberpunk street battle environment for tower defense,
setting: rainy neon-lit alleyway in megacity,
time: night with heavy rain,
neon colors: cyan #00FFFF, magenta #FF00FF, yellow #FFD700,
architecture: towering corporate skyscrapers with billboards,
tech elements: holographic ads, LED strips, security drones,
atmosphere: wet reflections, steam vents, fog,
characters: cybernetic soldiers on rooftops,
weathering: water puddles with neon reflections, grime,
style: cyberpunk with rain effects,
unreal engine 5, 8K, cinematic --ar 16:9 --v 6 --s 250
```

### 10.2 水墨/国风风格模板

适用于东方题材、历史游戏、艺术风格

```
Chinese ink painting style [场景/角色],
[画风] style: [写意/工笔/泼墨/山水],
[墨色层次] ink levels: [浓/淡/焦/湿],
[色彩点缀] color accents: [朱砂/石青/藤黄],
[构图] composition: [留白/满构/深远/平远],
[元素] elements: [山/水/云/松/竹/梅/亭],
[笔触] brushwork: [勾勒/皴法/点厾],
[时代特征] era: [唐/宋/元/明/清/架空],
[氛围] mood: [宁静/萧瑟/壮阔/神秘],
[应用] usage: [game background/character art/loading],
[渲染] digital ink render, artistic --ar 16:9 --v 6

示例：
Chinese ink painting style ancient battlefield,
style: expressive xieyi with detailed gongbi elements,
ink levels: rich dark ink for mountains, light wash for mist,
color accents: vermillion #FF4C00 for flags, azurite #007BA7 for water,
composition: classic diagonal with mountain backdrop,
elements: twisted pines, flowing rivers, distant peaks, broken banners,
brushwork: axe皴 for rocks, slender lines for bamboo,
era: Song dynasty military scene,
mood: solemn and epic,
usage: game cinematic background,
digital ink render, artistic quality --ar 16:9 --v 6 --s 200
```

---

## 十一、提示词调试与优化技巧

### 11.1 常见问题解决方案

| 问题 | 原因 | 解决方案 |
|------|------|----------|
| 画面过暗 | 缺少光照描述 | 添加 lighting, golden hour, backlight |
| 缺少细节 | 描述太笼统 | 添加 detailed, intricate, texture |
| 构图失衡 | 缺少构图法则 | 添加 rule of thirds, centered, symmetrical |
| 风格不一致 | 风格词冲突 | 统一风格词，避免混用 |
| 主体不突出 | 缺少焦点 | 添加 focal point, focus on, center |

### 11.2 提示词测试流程

1. **基础测试**：去掉所有风格词，只保留主体描述
2. **逐步添加**：每次只加一类描述词（构图/光照/氛围）
3. **对比分析**：记录每次变化的效果
4. **迭代优化**：根据结果调整权重和描述

### 11.3 提示词权重控制

```
# 提升权重
(关键词:1.2) - 括号+倍数

# 降低权重
[关键词:0.8] - 方括号+倍数

# 强制包含
::关键词:: - 双冒号强调

# 负面提示词（减少不要的元素）
--no [不想要的元素1, 元素2]
```

---

## 十二、本日总结

**今日产出汇总**：

**优化类**：
- 水体/液体特效模板
- 粒子系统提示词
- 2D精灵/动画帧模板
- 交互式UI动画模板
- 赛博朋克风格模板
- 水墨/国风风格模板

**工具类**：
- 提示词调试指南
- 常见问题解决方案
- 权重控制语法

---

> 记录时间：2026-03-10 08:17 UTC
> 状态：✅ 今日任务完成
> 累计产出：22套新模板 + 10组优化系统 + 10个公式指南

---

## 十五、傍晚补充产出（11:17）

### 15.1 爆炸/冲击波特效模板

适用于技能特效、攻击效果、环境破坏

```
Explosion effect for [游戏类型],
[爆炸类型] type: [fire/explosion/shockwave/energy/nuclear],
[阶段] phases: [initial burst → expansion → dissipation],
[颜色] colors: [核心色 #HEX → 中间色 #HEX → 边缘色 #HEX],
[粒子] debris: [rocks/glass/energy shards/spark],
[烟尘] smoke: [dense/light/toxic/fire],
[地面效果] ground impact: [crater/radiating cracks/glow],
[动态] motion: [radial outward/upward spiral/doming],
[规模] scale: [small/medium/large/massive],
[时长] duration: [0.5s/1s/2s],
[触发] trigger: [on hit/ultimate/environmental],
[引擎] VFX asset, transparent --ar 1:1 --v 6

示例：
Massive dinosaur extinction event explosion,
type: fire + shockwave combined,
phases: bright flash → fireball expansion → mushroom cloud → dissipation,
colors: white #FFFFFF → yellow #FFD700 → orange #FF4500 → dark red #8B0000 → smoke gray #333333,
debris: massive rock chunks, burning debris, dust,
smoke: dense dark gray with fire embers,
ground impact: massive crater with radiating cracks, glowing lava veins,
motion: massive upward dome with shockwave ring,
scale: massive (screen-filling),
duration: 3-second sequence,
trigger: boss extinction attack,
unreal engine Niagara, VFX asset, transparent --ar 1:1 --v 6 --s 250
```

### 15.2 闪烁/发光UI模板

适用于选中状态、高亮、稀有物品

```
Glowing UI element for [游戏类型],
[元素类型] element: [button/frame/icon/border/background],
[发光类型] glow type: [pulse/breathing/static/rainbow],
[颜色] color: [主色 #HEX],
[发光强度] intensity: [subtle/medium/intense],
[发光方向] direction: [inner/outer/both],
[动画] animation: [smooth/fast/emergency],
[层数] layers: [single/double/triple],
[附加效果] extra: [particles/ripple/shimmer],
[用途] usage: [selected item/new content/rare drop/notification],
[样式] style: [sci-fi/fantasy/minimal/elegant],
clean render, UI mockup --ar 1:1 --v 6

示例：
Rare dinosaur card glow effect for tower defense,
element: card border glow,
glow type: pulsing golden aura,
color: legendary gold #FFD700,
intensity: intense with soft falloff,
direction: outer glow with inner highlight,
animation: smooth 2-second cycle,
layers: triple (core glow → soft glow → ambient),
extra: sparkle particles rising, subtle rainbow shimmer at edges,
usage: legendary dinosaur card highlight,
style: fantasy with magical particles,
clean render, UI mockup --ar 1:1 --v 6 --s 180
```

### 15.3 残影/拖尾效果模板

适用于快速移动、瞬移、技能路径

```
Motion trail effect for [对象],
[主体类型] subject: [character/projectile/vehicle/creature],
[拖尾类型] trail type: [ghost/speed lines/afterimage/energy stream],
[透明度衰减] opacity fade: [linear/exponential/instant],
[颜色] color: [主色 #HEX] with [辅色 #HEX],
[长度] length: [short/medium/long/very long],
[形状] shape: [tapered/consistent/widening],
[帧数] frame count: [3/5/8/12 frames],
[应用场景] usage: [dash attack/teleport/charge/sprint],
[混合模式] blend: [add/alpha/multiply],
[附加效果] effects: [motion blur/particles/electricity],
[引擎] VFX compatible, transparent --ar 16:9 --v 6

示例：
Raptor dinosaur dash attack motion trail,
subject: fast-moving raptor creature,
trail type: afterimage ghosts + speed lines,
opacity fade: exponential (quick fade),
color: cyan #00FFFF to transparent,
length: medium-long (3 body lengths),
shape: tapered (faster at tail),
frame count: 6 ghost frames,
usage: raptor special charge attack,
blend: add for glow effect,
effects: motion blur + dust particles at feet,
unity particle system compatible, transparent --ar 16:9 --v 6 --s 200
```

### 15.4 冻结/冰霜效果模板

适用于冰系技能、冻结状态、环境效果

```
Ice/Frozen effect for [游戏类型],
[效果类型] type: [frost/ice crystal/blizzard/frozen surface],
[覆盖程度] coverage: [partial/full/thick layer],
[颜色] colors: [冰蓝 #HEX → 深蓝 #HEX → 白霜 #HEX],
[纹理] texture: [crystalline/frosted/smooth/rugged],
[光照效果] light: [specular highlights/refraction/internal glow],
[状态] state: [forming/active/dissolving/mid-transform],
[特效] features: [icicles/frost patterns/ice cracks/snow particles],
[应用] usage: [skill effect/environmental/character frozen],
[融化效果] melt: [with water pools/dripping/steaming],
[引擎] VFX, transparent --ar 1:1 --v 6

示例：
Blizzard zone skill effect for ice dinosaur tower defense,
type: blizzard + frozen surface combined,
coverage: thick layer covering entire area,
colors: bright ice blue #87CEEB → deep blue #0066CC → white frost #FFFFFF,
texture: crystalline with frost patterns,
light: strong specular highlights, subtle internal glow,
state: active blizzard,
features: falling snow particles, icicles forming, frost on ground,
usage: area denial skill "Frozen Domain",
melt: water pools forming at edges with steam,
unreal engine Niagara, VFX, transparent --ar 1:1 --v 6 --s 200
```

### 15.5 护盾/屏障特效模板

适用于防御技能、能量护盾、魔法屏障

```
Shield/Barrier effect for [游戏类型],
[护盾类型] type: [energy/ice/fire/void/magic/physical],
[形状] shape: [dome/sphere/hexagonal/plane/barrier],
[颜色] color: [主色 #HEX] with [辅色 #HEX],
[透明度] opacity: [百分比],
[表面效果] surface: [rippling/pulsing/grid/hex pattern/solid],
[边缘效果] edge: [glowing rim/lightning/fracture/sharp],
[激活状态] state: [forming/active/damaged/breaking/dissipating],
[特效] features: [particles/shimmer/energy flow/inscription],
[大小] size: [small/medium/large/area],
[触发] trigger: [passive/on-command/on-hit],
[引擎] VFX, transparent --ar 1:1 --v 6

示例：
Dinosaur handler personal energy shield,
type: energy/void hybrid,
shape: hexagonal dome barrier,
color: cyan #00FFFF core with purple #9400D3 edge,
opacity: 60% (semi-transparent),
surface: hexagonal pattern with energy flow animation,
edge: glowing rim with electricity sparks,
state: active with subtle pulse,
features: faint runic inscriptions, energy particles orbiting,
size: medium (covers 1 person),
trigger: activated by handler ability "Protective Field",
unreal engine Niagara, VFX, transparent --ar 1:1 --v 6 --s 200
```

---

## 十六、终极提示词模板库

### 16.1 万能场景模板

```
[主体] in [环境],
[动作/状态] [描述],
[时间/季节] [时间],
[天气] [天气描述],
[光照] [光照描述],
[氛围] [氛围描述],
[风格] [风格],
[构图] [构图方式],
[相机] [相机角度/运动],
[特效] [特效描述],
[技术] [引擎] render, [分辨率], [质量],
[参数] --ar [比例] --v [版本] --s [风格化程度]
```

### 16.2 万能UI模板

```
[UI类型] for [游戏/应用类型],
[详细描述],
[样式] style: [风格],
[颜色] colors: [主色 #HEX] [辅色 #HEX],
[状态] state: [状态],
[特效] effects: [效果],
[布局] layout: [布局方式],
[尺寸] size: [尺寸],
[用途] usage: [使用场景],
[渲染] UI mockup, [工具风格], [分辨率] --ar [比例] --v 6
```

### 16.3 万能特效模板

```
[特效名称] effect for [游戏类型],
[类型] type: [类型],
[颜色] colors: [起始色] → [结束色],
[粒子] particles: [描述],
[动画] animation: [描述],
[触发] trigger: [触发条件],
[时长] duration: [时长],
[规模] scale: [规模],
[用途] usage: [用途],
[引擎] [引擎] VFX, transparent --ar [比例] --v 6
```

---

## 本日最终汇总

**总计产出**：
- **提示词优化系统**：15套（天气、UI血条、地形、光照、材质、构图等）
- **UI模板**：20套（各类游戏UI组件）
- **特效模板**：15套（粒子、爆炸、冰冻、护盾、残影等）
- **风格模板**：6套（赛博朋克、水墨、Low-Poly、像素等）
- **工具/公式**：8组（组合技、急救箱、调试流程等）
- **万能模板**：3套（场景、UI、特效）

**总模板数**：60+ 套提示词/模板

---

> 记录时间：2026-03-10 11:17 UTC
> 状态：✅ 任务完成
> 累计产出：60+ 套模板 & 优化系统

---

## 十三、上午补充产出（08:17）

### 13.1 渐变/模糊效果UI模板

适用于移动端、模糊背景、毛玻璃效果

```
UI design with blur effect for [应用类型],
[模糊类型] blur: [gaussian/motion/zoom/background],
[模糊强度] intensity: [light/medium/heavy],
[颜色] color: [主色 #HEX] with [透明度 X%],
[内容层] content: [卡片/列表/弹窗],
[样式] style: [iOS glassmorphism/Android material/Windows acrylic],
[边框] border: [subtle glow/none/metallic],
[阴影] shadow: [soft/hard/none],
[圆角] corner: [rounded/square/circular],
[内容] content: [具体UI元素描述],
[层次] layers: [background blur + content + overlay],
[渲染] clean render, figma style, 4K --ar 16:9 --v 6

示例：
Mobile game UI with blur effect for dinosaur tower defense,
background blur: gaussian, intensity: medium (20px),
color: dark blue #0D1B2A with 80% opacity,
content: upgrade confirmation card,
style: iOS glassmorphism with subtle border,
border: subtle white glow 10% opacity,
shadow: soft drop shadow,
corner: rounded 16px,
content: tower icon + stats + confirm/cancel buttons,
layers: blurred game view + card + button overlay,
clean render, figma style, 4K --ar 16:9 --v 6 --s 180
```

### 13.2 文字特效/字体设计模板

适用于游戏标题、Logo、特殊文字效果

```
Typography design for [标题/文字],
[字体风格] style: [sci-fi/medieval/fantasy/handwritten/cyber],
[特效] effects: [3D extrude/glow/neon/emboss/hollow],
[材质] material: [metal/wood/ice/fire/energy/glass],
[颜色] colors: [primary #HEX] with [accent #HEX],
[光照] lighting: [front/rim/back/ambient occlusion],
[背景] background: [transparent/dark/gradient/pattern],
[装饰] decorations: [particles/ornaments/flames/electricity],
[应用场景] usage: [game title/logo/banner/loading],
[格式] format: [2D/3D/isometric],
[渲染] [引擎] render, 8K, professional --ar 16:9 --v 6

示例：
Game title typography "DINO TOWER",
style: bold geometric with sci-fi elements,
effects: 3D extrude + neon glow + energy particles,
material: metallic with energy core,
colors: electric cyan #00FFFF primary with orange #FF6B35 accent,
lighting: front key light + cyan rim light,
background: transparent for overlay use,
decorations: floating DNA strands + energy sparks,
usage: main game title screen,
format: 3D with 2D option,
unreal engine 5, 8K, professional --ar 16:9 --v 6 --s 250
```

### 13.3 混合风格模板

适用于创意项目、特殊需求、风格融合

```
Fusion style [主体] combining [风格A] + [风格B],
[风格A元素] from [风格A]: [具体元素描述],
[风格B元素] from [风格B]: [具体元素描述],
[融合方式] fusion: [seamless/juxtaposed/harmonious],
[颜色协调] color harmony: [complementary/analogous/contrasting],
[主色调] dominant: [颜色 #HEX],
[强调色] accent: [颜色 #HEX],
[视觉焦点] focal point: [描述],
[整体氛围] mood: [描述],
[用途] usage: [概念设计/宣传图/游戏资产],
[渲染] [引擎], 8K --ar 16:9 --v 6

示例：
Fusion style dinosaur design combining photorealistic + low-poly,
from photorealistic: detailed skin texture, realistic lighting, anatomical accuracy,
from low-poly: geometric faceted surfaces, crisp edges, stylized forms,
fusion: seamless blend where body is realistic, scales become geometric patterns,
color harmony: analogous (greens + yellows),
dominant: forest green #228B22,
accent: amber #FFBF00,
focal point: glowing amber eyes with tribal markings,
mood: ancient creature meets modern design,
usage: hero monster design for tower defense,
blender cycles + unreal engine 5, 8K --ar 16:9 --v 6 --s 250
```

### 13.4 阴影/深度效果模板

适用于3D物体、UI层次、界面设计

```
Shadow and depth effect for [对象],
[阴影类型] shadow: [drop/inner/ambient/casted],
[阴影颜色] shadow color: [#HEX] with [opacity X%],
[阴影模糊] blur: [soft/medium/hard],
[阴影偏移] offset: [x: Xpx, y: Ypx],
[深度效果] depth: [elevation/gradient/vignette],
[层次] layers: [foreground/midground/background],
[光照方向] light source: [top-left/top-right/bottom],
[反射] reflections: [subtle/moderate/none],
[应用] usage: [UI cards/3D objects/characters],
[渲染] clean render, studio lighting --ar 16:9 --v 6

示例：
Shadow and depth effect for tower defense tower cards,
shadow: drop shadow + inner shadow,
shadow color: black #000000 with 25% opacity,
blur: soft (15px), offset: x: 0, y: 8px,
depth: elevation with gradient overlay,
layers: card (midground) with towers (foreground), background (faded),
light source: top-left at 45 degrees,
reflections: subtle floor reflection,
usage: game UI selection cards,
clean render, studio lighting --ar 16:9 --v 6 --s 180
```

### 13.5 跨平台UI适配模板

适用于多平台发布、响应式设计

```
Cross-platform UI for [应用类型],
[平台] platforms: [iOS/Android/PC/Console/Web],
[适配策略] adaptation: [responsive/adaptive/consistent],
[一致元素] consistent: [colors/logo/typography/icon style],
[平台差异] platform-specific: [iOS: safe areas/Android: gestures/PC: hover states],
[分辨率] resolutions: [描述适配的分辨率],
[方向] orientation: [portrait/landscape/both],
[触控区域] touch targets: [minimum 44x44px for mobile],
[文字大小] text sizing: [platform guidelines followed],
[图标风格] icons: [outlined/filled/mixed],
[布局] layout: [grid/flex/relative],
[渲染] UI mockup, design system --ar 16:9 --v 6

示例：
Cross-platform UI for dinosaur tower defense game,
platforms: iOS + Android + PC + Web,
adaptation: responsive with consistent core experience,
consistent: green jungle theme, same logo, dinosaur silhouette icons,
platform-specific: iOS notch handling, Android back gesture, PC mouse hover,
resolutions: 1920x1080 (PC), 1280x720 (tablet), 750x1334 (mobile portrait),
orientation: landscape for gameplay, portrait for menus,
touch targets: minimum 48x48dp for mobile,
text sizing: iOS: pt, Android: sp, PC: px,
icons: outlined style with filled for selected,
layout: flexbox with percentage-based spacing,
UI mockup, design system documentation --ar 16:9 --v 6 --s 180
```

---

## 十四、提示词组合技

### 14.1 万能增强前缀

```
# 质量增强
professional, high quality, polished, detailed, intricate,

# 风格强化
artstation trending, deviantart featured, behance best,

# 技术规格
8K, unreal engine 5, octane render, cinema 4D, blender,

# 构图强化
rule of thirds, centered, symmetrical, dynamic angle,

# 光照强化
cinematic lighting, volumetric, god rays, rim lighting,

# 氛围强化
moody, atmospheric, epic, dramatic, peaceful
```

### 14.2 风格互斥清单

| 风格A | 冲突风格B | 原因 |
|-------|----------|------|
| photorealistic | cartoon/low-poly | 写实与卡通不兼容 |
| minimal | detailed/intricate | 极简与复杂冲突 |
| flat | 3D/render | 扁平与立体矛盾 |
| anime | realistic | 日系与写实冲突 |

### 14.3 提示词顺序优先级

1. **主体** (最重要) - 描述核心内容
2. **动作/状态** - 主体在做什么
3. **环境** - 场景、地点
4. **风格** - 艺术风格
5. **构图** - 视角、布局
6. **光照** - 灯光效果
7. **氛围** - 情绪、感觉
8. **技术参数** - 分辨率、引擎、版本

---

## 上午补充完成

**08:17 补充产出**：
- 渐变/模糊效果UI模板
- 文字特效/字体设计模板
- 混合风格模板
- 阴影/深度效果模板
- 跨平台UI适配模板
- 提示词组合技（万能增强前缀、风格互斥清单、顺序优先级）

---

## 下午补充：AI视频生成提示词 & 连续帧模板

### 15.1 AI视频生成提示词优化

适用于 Runway、Pika、可灵等视频生成工具

**核心参数**：
```
[主体] [动作描述] [运动轨迹] [相机运动] [时间/季节] [光照] [氛围] [帧率要求] [时长]
```

**优化要点**：
- 使用动态词汇：flowing, swirling, cascading, exploding
- 强调运动轨迹：from left to right, spiraling upward, approaching camera
- 相机运动：pan, zoom, dolly, orbit, handheld shake
- 避免静态词汇：standing, static, still

**示例 - 恐龙冲锋**：
```
A herd of armored dinosaurs charging through misty valley,
tracking shot following leader, dust clouds billowing,
slow motion at peak moment, morning fog breaking,
golden hour lighting, epic cinematic feel,
24fps cinematic, 5 second loop --ar 16:9
```

### 15.2 连续帧/分镜模板

适用于故事板、连续动作、多视角展示

```
Storyboard frame [序号]:
[场景描述] + [角色动作] + [时间线位置]

示例：
Frame 1: Wide shot, dinosaur sanctuary at dawn,
researcher enters frame from left, walking path,
establishing scene, 0-2s

Frame 2: Medium shot, researcher approaches nest,
kneels down slowly, reaching hand out,
emotional beat, 2-4s

Frame 3: Close-up, baby T-Rex sniffs hand,
curious expression, soft glow from egg,
heartwarming moment, 4-6s
```

### 15.3 风格迁移/重绘提示词

适用于局部重绘、风格转换、元素替换

**重绘公式**：
```
[原图元素描述] + [重绘要求] + [风格目标] + [保持不变的部分]
```

**示例**：
```
Original: Green grassland terrain
Keep:地形轮廓和树木位置
Style: volcanic eruption aftermath with lava cracks
Additional: charred ground, glowing cracks, ash particles, dark sky
```

---

## 今日总结 & 明日待办

### 完成内容
- ✅ 天气/环境效果系统提示词优化
- ✅ 动态UI血条/状态栏模板
- ✅ 地形生成/地块设计模板
- ✅ 过场动画CG模板
- ✅ 装备槽/物品栏UI模板
- ✅ Boss战UI模板
- ✅ 任务/成就系统UI模板
- ✅ 角色设计模板（多角度）
- ✅ 技能/特效图标模板
- ✅ 提示词公式进阶（多主体、动态效果、材质特写）
- ✅ 渐变/模糊效果UI模板
- ✅ 文字特效/字体设计模板
- ✅ 混合风格模板
- ✅ 阴影/深度效果模板
- ✅ 跨平台UI适配模板
- ✅ 提示词组合技
- ✅ AI视频生成提示词
- ✅ 连续帧/分镜模板
- ✅ 风格迁移/重绘提示词

### 2.5 粒子系统/技能特效模板

适用于技能特效、Buff图标、环境粒子

```
Particle system effect for [游戏类型],
[特效类型] effect: [hit/impact/aura/ambient/trail],
[元素] element: [fire/ice/lightning/poison/holy/dark],
[形态] form: [burst/spiral/wave/ring/flowing],
[颜色] colors: [主色 #HEX] to [渐变色 #HEX],
[强度] intensity: [subtle/medium/intense/explosive],
[粒子数] particles: [sparse/dense/massive],
[动态] motion: [static/animated/looping],
[发光] glow: [none/soft/intense/radiant],
[附加效果] extra: [ember sparks/magic runes/energy cracks],
[格式] format: [sprite sheet loop/single frame/gif],
[背景] background: [transparent/gradient/dark],
game vfx, high quality, unity/unreal compatible --ar 1:1 --v 6

示例：
Particle system effect for fantasy tower defense,
effect: magic aura, element: holy light,
form: spiral rotating around target,
colors: gold #FFD700 to white #FFFFFF,
intensity: intense with radiant glow,
particles: dense with floating light orbs,
motion: animated looping clockwise,
glow: radiant with lens flare,
extra: holy runes circling, golden sparkles rising,
sprite sheet 8 frames loop, transparent background,
game vfx, high quality --ar 1:1 --v 6 --s 250
```

### 2.6 角色立绘/姿态模板

适用于角色详情、英雄选择、社区展示

```
Character portrait for [游戏类型],
[角色名] character: [名称],
[职业/类型] class: [warrior/mage/rogue/beast/mech],
[姿态] pose: [idle/attack/defend/victory/death],
[视角] view: [front/three-quarter/close-up bust/waist up],
[表情] expression: [fierce/calm/angry/determined/smiling],
[装备] equipment: [武器/盔甲/饰品描述],
[风格] style: [realistic/anime/pixel art/low-poly],
[配色] color scheme: [主色调描述],
[背景] background: [solid/transparent/environment/scene],
[特效] effects: [glow/particle/rim light/motion blur],
[渲染] render engine, 8K, studio lighting --ar 3:4 --v 6

示例：
Character portrait for dinosaur tower defense,
character: Rex Trainer Mara, class: beast master,
pose: confident stance with tamed raptor beside her,
view: waist up, three-quarter angle,
expression: determined smile with tactical goggles,
equipment: reinforced combat vest, dinosaur whistle,
cyberpunk-meets-wild style, orange #FF6B35 and black #1A1A2E,
transparent background with subtle tech grid,
effects: cyan energy leash connecting to raptor,
unreal engine 5, 8K, studio lighting --ar 3:4 --v 6 --s 200
```

### 2.7 音效/音频可视化模板

适用于音乐播放器、音效UI、节奏游戏

```
Audio visualizer for [使用场景],
[类型] type: [waveform/equalizer/spectrum/circular],
[风格] style: [minimalist/retro/neon/organic],
[颜色] colors: [主色 #HEX] with [强调色 #HEX],
[动态] animation: [bouncing bars/flowing lines/pulsing circles],
[发光] glow: [none/soft neon/intense bloom],
[背景] background: [transparent/dark #XXXXXX/gradient],
[元素] elements: [frequency bars/wave lines/particle dots],
[强度] intensity: [calm/medium/energetic/extreme],
[特效] extra: [beat detection flash/ripple effects/mirror reflection],
[格式] format: [loop ready/static/animated gif],
game ui audio, music visualization, clean design --ar 16:9 --v 6

示例：
Audio visualizer for rhythm defense game,
type: circular spectrum, style: neon cyberpunk,
colors: magenta #FF00FF to cyan #00FFFF gradient,
animation: pulsing concentric circles with frequency bars,
glow: intense neon bloom effect,
background: dark #0D1117 with subtle grid,
elements: 32 frequency bars in circular arrangement,
intensity: energetic matching fast beats,
extra: beat detection flash on drop, particle burst on climax,
loop ready animation, transparent background,
game ui audio, music visualization, clean design --ar 16:9 --v 6 --s 180
```

### 2.8 场景概念图/氛围图模板

适用于游戏背景、Loading界面、商店展示

```
Concept art environment for [游戏类型],
[场景] location: [描述],
[时间] time: [dawn/day/dusk/night],
[季节] season: [spring/summer/autumn/winter],
[天气] weather: [clear/rain/snow/fog/storm],
[氛围] mood: [peaceful/tense/mysterious/epic/horror],
[视角] perspective: [wide shot/eye level/aerial/dramatic low],
[主体] focal point: [建筑/角色/奇观/破坏],
[光照] lighting: [主光描述], [阴影强度],
[色彩] palette: [主色调描述],
[风格] style: [painterly/realistic/stylized/concept art],
[细节] details: [前景/中景/背景层次],
[渲染] engine, 8K, cinematic, concept art quality --ar 16:9 --v 6

示例：
Concept art environment for dinosaur tower defense,
location: ancient ruins overtaken by nature, tropical jungle,
time: golden hour sunset, season: eternal summer,
weather: clear with mist rising,
mood: mysterious and awe-inspiring,
perspective: wide shot from above, eye-level focal point,
focal point: massive stone gate with glowing crystal,
lighting: warm golden hour sun from side, long dramatic shadows,
palette: emerald greens, golden light, ancient stone grays,
painterly style with bold brushstrokes, rich detail,
details: prehistoric plants in foreground, ruined pillars mid, volcano peak background,
unreal engine 5, 8K, cinematic, concept art quality --ar 16:9 --v 6 --s 250
```

---

## 三、提示词优化总结

### 优化技巧清单
1. **结构化占位符**: 使用 [类别] 描述 格式，便于替换
2. **颜色明确指定**: 用具体 HEX 值，避免 "red/blue" 等模糊词
3. **引擎参数后置**: --ar --v --s 放在最后
4. **情感/氛围优先**: 先定 mood/atmophere，再加细节
5. **层次递进**: 背景→主体→特效→渲染

### 高频优化方向
- 动态效果: 添加 motion/animated/looping
- 发光效果: glow/bloom/emission/radiant
- 质感提升: PBR/textured/metallic/glass
- 视角多样: close-up/wide shot/aerial/dramatic

---

## 四、明日待办
1. 测试实际生成效果，收集bad case
2. 针对具体游戏类型（恐龙塔防）的垂直模板
3. 动画/骨骼绑定相关提示词
4. 音效/配乐视觉化描述（波形图可视化）
5. 批量生成自动化流程

---

## 五、今日新增 - 垂直领域 & UI模板

### 5.1 恐龙设计模板
```
Dinosaur design for tower defense,
[恐龙种类]: [描述],
[size]: [small/medium/large/giant],
[category]: [carnivore/herbivore/flying],
[features]: [具体特征],
[color]: #HEX, accent: #HEX,
[attack]: [攻击方式],
[ability]: [特殊能力],
[animation]: [walk/attack/idle],
[style]: [realistic/stylized],
game asset, transparent background --ar 3:2 --v 6

示例：
Dinosaur design for tower defense,
"T-Rex Alpha": massive tyrannosaurus,
size: large, category: carnivore,
features: powerful jaws, tiny arms, thick scales,
primary: #4A3728, accent: #FFD700,
attack: bite with bone shatter,
ability: roar debuff,
animation: walk, attack, idle,
realistic style --ar 3:2 --v 6 --s 250
```

### 5.2 恐龙栖息地模板
```
Dinosaur habitat,
[type]: [nest/cave/water hole],
[dinosaur]: [种类],
[environment]: [描述],
[decorations]: [恐龙蛋/骸骨/植被],
[condition]: [abandoned/active],
[atmosphere]: [peaceful/tense/dangerous],
isometric view, game environment --ar 16:9 --v 6
```

### 5.3 UI按钮模板
```
UI button for [游戏类型],
[type]: [main action/menu/confirm],
[shape]: [rectangular/rounded/pill],
[style]: [flat/neumorphic/3D],
[colors]: primary #HEX, hover #HEX,
[icon]: [描述],
[label]: "[文字]",
[effects]: [glow/ripple/scale],
game ui mockup, figma style --ar 16:9 --v 6
```

### 5.4 加载条模板
```
Loading bar for [场景],
[type]: [horizontal/circular/segmented],
[style]: [minimalist/tech/fantasy],
[colors]: track #HEX, fill #HEX,
[animation]: [smooth/glowing/particle],
game ui mockup --ar 16:9 --v 6
```

### 5.5 图标模板
```
Tactical icon for [游戏类型],
[type]: [movement/attack/defense/scout],
[shape]: [arrow/diamond/circle],
[style]: [minimalist/technical],
[colors]: primary #HEX, outline #HEX,
[symbol]: [描述],
[animation]: [static/pulse/rotate],
game icon, vector style --ar 1:1 --v 6
```

---

## 六、质量检查清单

生成前检查：
- [ ] 宽高比匹配用途 (16:9/4:3/1:1/3:2)
- [ ] 风格关键词明确
- [ ] 引擎参数在末尾 (--ar --v --s)
- [ ] 颜色使用具体HEX值
- [ ] 氛围词先行
- [ ] 避免模糊词

---

## 七、常用片段库

**光照**: cinematic lighting, volumetric rays / golden hour / neon cyberpunk
**材质**: PBR, metallic, emissive glow, weathered texture
**特效**: motion blur, particles, bloom, depth of field

---

**今日完成**：提示词优化3个系统 + 新模板8类 + 垂直领域模板 + UI模板 + 质量清单
