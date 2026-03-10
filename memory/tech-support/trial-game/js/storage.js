// 本地存储系统
class StorageManager {
    constructor() {
        this.prefix = 'barrage_game_';
    }
    
    save(key, value) {
        try {
            localStorage.setItem(this.prefix + key, JSON.stringify(value));
        } catch(e) {
            console.log('Storage not available');
        }
    }
    
    load(key, defaultValue = null) {
        try {
            const data = localStorage.getItem(this.prefix + key);
            return data ? JSON.parse(data) : defaultValue;
        } catch(e) {
            return defaultValue;
        }
    }
    
    // 保存游戏进度
    saveGameProgress(progress) {
        this.save('progress', {
            level: progress.level || 1,
            gold: progress.gold || 1000,
            gems: progress.gems || 50,
            power: progress.power || 100,
            achievements: progress.achievements || []
        });
    }
    
    // 加载游戏进度
    loadGameProgress() {
        return this.load('progress', {
            level: 1,
            gold: 1000,
            gems: 50,
            power: 100,
            achievements: []
        });
    }
}

const storageManager = new StorageManager();
