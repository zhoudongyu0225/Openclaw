// 战斗系统
class BattleSystem {
    constructor(game) {
        this.game = game;
        this.wave = 1;
        this.enemies = [];
    }
    
    // 开始波次
    startWave() {
        const count = 5 + this.wave * 2;
        for (let i = 0; i < count; i++) {
            this.spawnEnemy();
        }
    }
    
    // 生成敌人
    spawnEnemy() {
        const types = ['droid', 'raptor', 'tricera', 'trex'];
        const type = types[Math.floor(Math.random() * types.length)];
        
        this.enemies.push({
            type: type,
            hp: this.getHP(type),
            maxHp: this.getHP(type),
            x: -50,
            y: Math.random() * this.game.canvas.height,
            speed: 2 + this.wave * 0.1
        });
    }
    
    getHP(type) {
        const hps = { droid: 30, raptor: 50, tricera: 100, trex: 300 };
        return hps[type];
    }
    
    // 更新
    update() {
        this.enemies.forEach(e => {
            e.x += e.speed;
        });
        
        // 敌人到达终点
        this.enemies = this.enemies.filter(e => e.x < this.game.canvas.width + 50);
    }
}
