// 万龙觉醒 - 试玩 Demo
// 核心：ARPG 摇杆控制 + 技能系统

class Game {
    constructor() {
        this.canvas = document.getElementById('game-canvas');
        this.ctx = this.canvas.getContext('2d');
        this.resize();
        
        this.player = {
            x: 0,
            y: 0,
            vx: 0,
            vy: 0,
            speed: 5,
            direction: 0,
            attacking: false,
            attackFrame: 0
        };
        
        this.enemies = [];
        this.particles = [];
        this.camera = { x: 0, y: 0 };
        
        this.resources = {
            gold: 1000,
            gems: 50,
            power: 100
        };
        
        this.keys = {};
        this.joystickActive = false;
        this.joystickCenter = { x: 0, y: 0 };
        
        this.init();
    }
    
    resize() {
        this.canvas.width = window.innerWidth;
        this.canvas.height = window.innerHeight;
    }
    
    init() {
        // 玩家初始位置
        this.player.x = this.canvas.width / 2;
        this.player.y = this.canvas.height / 2;
        
        // 绑定事件
        window.addEventListener('resize', () => this.resize());
        this.bindControls();
        
        // 生成敌人
        this.spawnEnemies();
        
        // 游戏循环
        this.lastTime = performance.now();
        this.gameLoop();
    }
    
    bindControls() {
        // 键盘
        document.addEventListener('keydown', e => this.keys[e.key] = true);
        document.addEventListener('keyup', e => this.keys[e.key] = false);
        
        // 虚拟摇杆
        const joystick = document.getElementById('joystick');
        joystick.addEventListener('touchstart', e => this.onJoystickStart(e));
        joystick.addEventListener('touchmove', e => this.onJoystickMove(e));
        joystick.addEventListener('touchend', e => this.onJoystickEnd(e));
        
        // 攻击按钮
        document.getElementById('attack-btn').addEventListener('click', () => {
            this.player.attacking = true;
            this.player.attackFrame = 0;
            this.createAttackEffect();
        });
        
        // 技能按钮
        document.getElementById('skill-1').addEventListener('click', () => {
            this.useSkill(1);
        });
        document.getElementById('skill-2').addEventListener('click', () => {
            this.useSkill(2);
        });
    }
    
    onJoystickStart(e) {
        e.preventDefault();
        const touch = e.touches[0];
        const rect = document.getElementById('joystick').getBoundingClientRect();
        this.joystickCenter = {
            x: rect.left + rect.width / 2,
            y: rect.top + rect.height / 2
        };
        this.joystickActive = true;
    }
    
    onJoystickMove(e) {
        e.preventDefault();
        if (!this.joystickActive) return;
        
        const touch = e.touches[0];
        const dx = touch.clientX - this.joystickCenter.x;
        const dy = touch.clientY - this.joystickCenter.y;
        const distance = Math.min(Math.sqrt(dx*dx + dy*dy), 35);
        const angle = Math.atan2(dy, dx);
        
        const knob = document.querySelector('.joystick-knob');
        knob.style.transform = `translate(calc(-50% + ${distance * Math.cos(angle)}px), calc(-50% + ${distance * Math.sin(angle)}px))`;
        
        this.player.vx = (distance / 35) * this.player.speed * Math.cos(angle);
        this.player.vy = (distance / 35) * this.player.speed * Math.sin(angle);
        this.player.direction = angle;
    }
    
    onJoystickEnd(e) {
        this.joystickActive = false;
        this.player.vx = 0;
        this.player.vy = 0;
        
        const knob = document.querySelector('.joystick-knob');
        knob.style.transform = 'translate(-50%, -50%)';
    }
    
    spawnEnemies() {
        // 初始生成一些敌人
        for (let i = 0; i < 5; i++) {
            this.enemies.push({
                x: Math.random() * this.canvas.width,
                y: Math.random() * this.canvas.height,
                hp: 100,
                maxHp: 100,
                type: ['goblin', 'orc', 'dragon'][Math.floor(Math.random() * 3)],
                vx: (Math.random() - 0.5) * 2,
                vy: (Math.random() - 0.5) * 2
            });
        }
        
        // 定时生成新敌人
        setInterval(() => {
            if (this.enemies.length < 10) {
                this.enemies.push({
                    x: Math.random() > 0.5 ? -50 : this.canvas.width + 50,
                    y: Math.random() * this.canvas.height,
                    hp: 100,
                    maxHp: 100,
                    type: ['goblin', 'orc', 'dragon'][Math.floor(Math.random() * 3)],
                    vx: (Math.random() - 0.5) * 2,
                    vy: (Math.random() - 0.5) * 2
                });
            }
        }, 3000);
    }
    
    useSkill(skillId) {
        // 技能效果
        const skillColors = ['#ff4444', '#00d4ff'];
        for (let i = 0; i < 20; i++) {
            this.particles.push({
                x: this.player.x,
                y: this.player.y,
                vx: (Math.random() -) * 15 0.5,
                vy: (Math.random() - 0.5) * 15,
                life: 30,
                color: skillColors[skillId - 1],
                size: 10 + Math.random() * 10
            });
        }
        
        // 技能冷却视觉
        const btn = document.getElementById(`skill-${skillId}`);
        btn.style.opacity = '0.5';
        setTimeout(() => btn.style.opacity = '1', 1000);
    }
    
    createAttackEffect() {
        // 攻击特效
        for (let i = 0; i < 8; i++) {
            const angle = (Math.PI * 2 / 8) * i + this.player.direction;
            this.particles.push({
                x: this.player.x + Math.cos(angle) * 30,
                y: this.player.y + Math.sin(angle) * 30,
                vx: Math.cos(angle) * 8,
                vy: Math.sin(angle) * 8,
                life: 15,
                color: '#ffd700',
                size: 5
            });
        }
    }
    
    update() {
        // 键盘控制
        if (this.keys['ArrowUp'] || this.keys['w']) this.player.vy = -this.player.speed;
        if (this.keys['ArrowDown'] || this.keys['s']) this.player.vy = this.player.speed;
        if (this.keys['ArrowLeft'] || this.keys['a']) this.player.vx = -this.player.speed;
        if (this.keys['ArrowRight'] || this.keys['d']) this.player.vx = this.player.speed;
        
        // 更新位置
        this.player.x += this.player.vx;
        this.player.y += this.player.vy;
        
        // 边界检测
        this.player.x = Math.max(20, Math.min(this.canvas.width - 20, this.player.x));
        this.player.y = Math.max(20, Math.min(this.canvas.height - 200, this.player.y));
        
        // 更新敌人
        this.enemies.forEach(enemy => {
            // 追踪玩家
            const dx = this.player.x - enemy.x;
            const dy = this.player.y - enemy.y;
            const dist = Math.sqrt(dx*dx + dy*dy);
            
            if (dist > 100) {
                enemy.x += (dx / dist) * enemy.vx;
                enemy.y += (dy / dist) * enemy.vy;
            }
            
            // 攻击检测
            if (dist < 40 && this.player.attacking && this.player.attackFrame < 10) {
                enemy.hp -= 2;
            }
        });
        
        // 移除死亡敌人
        this.enemies = this.enemies.filter(e => e.hp > 0);
        
        // 更新粒子
        this.particles.forEach(p => {
            p.x += p.vx;
            p.y += p.vy;
            p.life--;
        });
        this.particles = this.particles.filter(p => p.life > 0);
        
        // 攻击动画
        if (this.player.attacking) {
            this.player.attackFrame++;
            if (this.player.attackFrame > 15) {
                this.player.attacking = false;
            }
        }
    }
    
    render() {
        // 清空画布
        this.ctx.fillStyle = '#1a1a2e';
        this.ctx.fillRect(0, 0, this.canvas.width, this.canvas.height);
        
        // 绘制背景网格
        this.drawGrid();
        
        // 绘制敌人
        this.enemies.forEach(enemy => this.drawEnemy(enemy));
        
        // 绘制玩家
        this.drawPlayer();
        
        // 绘制粒子
        this.particles.forEach(p => {
            this.ctx.beginPath();
            this.ctx.arc(p.x, p.y, p.size, 0, Math.PI * 2);
            this.ctx.fillStyle = p.color;
            this.ctx.fill();
        });
    }
    
    drawGrid() {
        this.ctx.strokeStyle = 'rgba(255,215,0,0.1)';
        this.ctx.lineWidth = 1;
        
        const gridSize = 50;
        for (let x = 0; x < this.canvas.width; x += gridSize) {
            this.ctx.beginPath();
            this.ctx.moveTo(x, 0);
            this.ctx.lineTo(x, this.canvas.height);
            this.ctx.stroke();
        }
        for (let y = 0; y < this.canvas.height; y += gridSize) {
            this.ctx.beginPath();
            this.ctx.moveTo(0, y);
            this.ctx.lineTo(this.canvas.width, y);
            this.ctx.stroke();
        }
    }
    
    drawPlayer() {
        const { x, y, direction, attacking, attackFrame } = this.player;
        
        // 玩家圆圈
        this.ctx.beginPath();
        this.ctx.arc(x, y, 25, 0, Math.PI * 2);
        this.ctx.fillStyle = '#ffd700';
        this.ctx.fill();
        this.ctx.strokeStyle = '#ff8c00';
        this.ctx.lineWidth = 3;
        this.ctx.stroke();
        
        // 武器指示器
        if (attacking) {
            const weaponLength = 20 + attackFrame * 2;
            this.ctx.beginPath();
            this.ctx.moveTo(x, y);
            this.ctx.lineTo(
                x + Math.cos(direction) * weaponLength,
                y + Math.sin(direction) * weaponLength
            );
            this.ctx.strokeStyle = '#ff4444';
            this.ctx.lineWidth = 5;
            this.ctx.stroke();
        }
        
        // 方向指示
        this.ctx.beginPath();
        this.ctx.arc(
            x + Math.cos(direction) * 15,
            y + Math.sin(direction) * 15,
            5, 0, Math.PI * 2
        );
        this.ctx.fillStyle = '#fff';
        this.ctx.fill();
    }
    
    drawEnemy(enemy) {
        const colors = {
            goblin: '#44aa44',
            orc: '#aa4444',
            dragon: '#8844aa'
        };
        
        // 敌人本体
        this.ctx.beginPath();
        this.ctx.arc(enemy.x, enemy.y, 20, 0, Math.PI * 2);
        this.ctx.fillStyle = colors[enemy.type] || '#666';
        this.ctx.fill();
        
        // 血条
        const hpWidth = 30;
        const hpPercent = enemy.hp / enemy.maxHp;
        
        this.ctx.fillStyle = '#333';
        this.ctx.fillRect(enemy.x - hpWidth/2, enemy.y - 35, hpWidth, 5);
        
        this.ctx.fillStyle = hpPercent > 0.5 ? '#44ff44' : hpPercent > 0.25 ? '#ffff44' : '#ff4444';
        this.ctx.fillRect(enemy.x - hpWidth/2, enemy.y - 35, hpWidth * hpPercent, 5);
    }
    
    gameLoop() {
        this.update();
        this.render();
        requestAnimationFrame(() => this.gameLoop());
    }
}

// 启动游戏
let game;
document.getElementById('start-btn').addEventListener('click', () => {
    document.getElementById('loading-screen').classList.add('hidden');
    document.getElementById('main-menu').classList.add('hidden');
    document.getElementById('game-screen').classList.remove('hidden');
    
    game = new Game();
});

// 关闭加载动画
setTimeout(() => {
    document.querySelector('.loading-text').textContent = '准备就绪！';
}, 2000);

// 转化弹窗 - 30秒后显示
setTimeout(() => {
    // document.getElementById('conversion-modal').classList.remove('hidden');
}, 30000);

document.querySelector('.close-btn')?.addEventListener('click', () => {
    document.getElementById('conversion-modal').classList.add('hidden');
});

// 升级系统
class UpgradeSystem {
    constructor(game) {
        this.game = game;
        this.buildings = [];
        this.init();
    }
    
    init() {
        // 可建造的建筑
        this.buildingTypes = {
            mine: { name: '金矿', cost: 100, income: 10 },
            farm: { name: '农场', cost: 80, income: 8 },
            barracks: { name: '兵营', cost: 200, troops: 5 },
            tower: { name: '防御塔', cost: 150, defense: 20 }
        };
    }
    
    addBuilding(type) {
        const building = this.buildingTypes[type];
        if (this.game.resources.gold >= building.cost) {
            this.game.resources.gold -= building.cost;
            this.buildings.push({
                type: type,
                ...building,
                level: 1
            });
            this.game.updateUI();
        }
    }
    
    render(ctx) {
        this.buildings.forEach((b, i) => {
            const x = 50 + (i % 4) * 80;
            const y = 150 + Math.floor(i / 4) * 80;
            
            ctx.fillStyle = '#444';
            ctx.fillRect(x, y, 60, 60);
            ctx.fillStyle = '#ffd700';
            ctx.font = '30px Arial';
            ctx.fillText(['⛏️', '🌾', '⚔️', '🗼'][i % 4], x + 15, y + 40);
        });
    }
}

// 成就系统
class AchievementSystem {
    constructor() {
        this.achievements = [
            { id: 'first_kill', name: '首次击杀', desc: '击杀第一个敌人', reward: 50, unlocked: false },
            { id: 'ten_kills', name: '小试牛刀', desc: '击杀10个敌人', reward: 100, unlocked: false },
            { id: 'rich', name: '富甲一方', desc: '拥有1000金币', reward: 200, unlocked: false },
            { id: 'power_500', name: '战力惊人', desc: '战力达到500', reward: 300, unlocked: false }
        ];
    }
    
    check(game) {
        this.achievements.forEach(ach => {
            if (!ach.unlocked) {
                if (ach.id === 'rich' && game.resources.gold >= 1000) this.unlock(ach, game);
                if (ach.id === 'power_500' && game.resources.power >= 500) this.unlock(ach, game);
            }
        });
    }
    
    unlock(ach, game) {
        ach.unlocked = true;
        game.resources.gems += ach.reward;
        this.showNotification(ach.name, ach.reward);
    }
    
    showNotification(name, reward) {
        const notif = document.createElement('div');
        notif.style.cssText = `
            position: fixed; top: 50%; left: 50%; transform: translate(-50%, -50%);
            background: linear-gradient(135deg, #ffd700, #ff8c00); padding: 20px 40px;
            border-radius: 10px; color: #000; font-weight: bold; z-index: 2000;
            animation: fadeInOut 2s forwards;
        `;
        notif.innerHTML = `🏆 成就解锁: ${name}<br>+${reward} 💎`;
        document.body.appendChild(notif);
        setTimeout(() => notif.remove(), 2000);
    }
}

console.log('游戏系统已加载');
