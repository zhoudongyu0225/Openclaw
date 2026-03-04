// 敌人系统
class EnemyManager {
    constructor(game) {
        this.game = game;
        this.enemies = [];
        this.spawnTimer = 0;
        this.spawnInterval = 2000; // 2秒
        this.wave = 1;
    }
    
    update(dt) {
        this.spawnTimer += dt;
        
        // 波次生成
        if (this.spawnTimer >= this.spawnInterval) {
            this.spawnTimer = 0;
            this.spawnWave();
        }
        
        // 更新所有敌人
        this.enemies.forEach(enemy => {
            this.updateEnemy(enemy, dt);
        });
        
        // 移除死亡敌人
        this.enemies = this.enemies.filter(e => e.hp > 0);
    }
    
    spawnWave() {
        const enemyCount = Math.min(this.wave * 2, 20);
        
        for (let i = 0; i < enemyCount; i++) {
            const type = this.getEnemyType();
            const enemy = {
                id: 'enemy_' + Date.now() + '_' + i,
                type: type,
                x: Math.random() * this.game.canvas.width,
                y: -50, // 从顶部生成
                targetX: this.game.player.x + (Math.random() - 0.5) * 200,
                targetY: this.game.player.y + (Math.random() - 0.5) * 200,
                hp: type.hp * (1 + this.wave * 0.1),
                maxHp: type.hp * (1 + this.wave * 0.1),
                speed: type.speed,
                damage: type.damage,
                reward: type.reward,
                state: 'move',
                animFrame: 0
            };
            this.enemies.push(enemy);
        }
        
        // 每10波难度提升
        if (this.enemies.length === 0) {
            this.wave++;
        }
    }
    
    getEnemyType() {
        const types = [
            { name: 'scout', hp: 30, speed: 3, damage: 5, reward: 10, color: '#88ff88' },   // 侦察兵
            { name: 'warrior', hp: 60, speed: 2, damage: 10, reward: 20, color: '#ff8888' }, // 战士
            { name: 'tank', hp: 150, speed: 1, damage: 20, reward: 50, color: '#ff44ff' },   // 坦克
            { name: 'boss', hp: 500, speed: 0.5, damage: 50, reward: 200, color: '#ff0000' } // Boss
        ];
        
        // 随机选择，boss每5波出现
        if (this.wave % 5 === 0 && Math.random() < 0.3) {
            return types[3];
        }
        
        const rand = Math.random();
        if (rand < 0.5) return types[0];
        if (rand < 0.8) return types[1];
        return types[2];
    }
    
    updateEnemy(enemy, dt) {
        // 移动向目标
        const dx = enemy.targetX - enemy.x;
        const dy = enemy.targetY - enemy.y;
        const dist = Math.sqrt(dx * dx + dy * dy);
        
        if (dist > 5) {
            enemy.x += (dx / dist) * enemy.speed;
            enemy.y += (dy / dist) * enemy.speed;
            enemy.state = 'move';
        } else {
            enemy.state = 'idle';
        }
        
        // 攻击玩家
        const pdx = this.game.player.x - enemy.x;
        const pdy = this.game.player.y - enemy.y;
        const pdist = Math.sqrt(pdx * pdx + pdy * pdy);
        
        if (pdist < 50 && enemy.state !== 'attack') {
            enemy.state = 'attack';
            this.game.resources.power -= enemy.damage;
            this.game.updateUI();
            
            // 受伤特效
            this.game.createParticles(this.game.player.x, this.game.player.y, '#ff4444', 5);
        }
        
        enemy.animFrame++;
    }
    
    takeDamage(enemy, damage) {
        enemy.hp -= damage;
        
        if (enemy.hp <= 0) {
            this.game.resources.gold += enemy.reward;
            this.game.resources.power += 5;
            this.game.updateUI();
            this.game.createParticles(enemy.x, enemy.y, '#ffd700', 10);
        }
    }
    
    render(ctx) {
        this.enemies.forEach(enemy => {
            // 敌人身体
            ctx.beginPath();
            ctx.arc(enemy.x, enemy.y, 15 + enemy.hp / 50, 0, Math.PI * 2);
            ctx.fillStyle = enemy.type.color || '#ff4444';
            ctx.fill();
            
            // 敌人类型图标
            const icons = { scout: '🦖', warrior: '🐲', tank: '🦕', boss: '🐉' };
            ctx.font = '20px Arial';
            ctx.fillText(icons[enemy.type.name] || '🦖', enemy.x - 10, enemy.y + 5);
            
            // 血条
            const hpWidth = 30;
            const hpPercent = enemy.hp / enemy.maxHp;
            ctx.fillStyle = '#333';
            ctx.fillRect(enemy.x - hpWidth / 2, enemy.y - 25, hpWidth, 4);
            ctx.fillStyle = hpPercent > 0.5 ? '#44ff44' : '#ffff44';
            ctx.fillRect(enemy.x - hpWidth / 2, enemy.y - 25, hpWidth * hpPercent, 4);
        });
    }
}

// 粒子系统
class ParticleSystem {
    constructor(game) {
        this.game = game;
        this.particles = [];
    }
    
    emit(x, y, color, count = 10) {
        for (let i = 0; i < count; i++) {
            this.particles.push({
                x: x,
                y: y,
                vx: (Math.random() - 0.5) * 10,
                vy: (Math.random() - 0.5) * 10,
                life: 30 + Math.random() * 20,
                maxLife: 50,
                color: color,
                size: 3 + Math.random() * 5
            });
        }
    }
    
    update() {
        this.particles.forEach(p => {
            p.x += p.vx;
            p.y += p.vy;
            p.life--;
            p.size *= 0.95;
        });
        
        this.particles = this.particles.filter(p => p.life > 0);
    }
    
    render(ctx) {
        this.particles.forEach(p => {
            ctx.globalAlpha = p.life / p.maxLife;
            ctx.beginPath();
            ctx.arc(p.x, p.y, p.size, 0, Math.PI * 2);
            ctx.fillStyle = p.color;
            ctx.fill();
        });
        ctx.globalAlpha = 1;
    }
}
