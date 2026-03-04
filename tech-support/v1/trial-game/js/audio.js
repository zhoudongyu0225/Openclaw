// 简单的音效系统
class AudioManager {
    constructor() {
        this.context = null;
        this.sounds = {};
    }
    
    init() {
        try {
            this.context = new (window.AudioContext || window.webkitAudioContext)();
        } catch(e) {
            console.log('Audio not supported');
        }
    }
    
    play(type) {
        if (!this.context) return;
        
        const oscillator = this.context.createOscillator();
        const gainNode = this.context.createGain();
        
        oscillator.connect(gainNode);
        gainNode.connect(this.context.destination);
        
        switch(type) {
            case 'attack':
                oscillator.frequency.setValueAtTime(200, this.context.currentTime);
                oscillator.frequency.exponentialRampToValueAtTime(50, this.context.currentTime + 0.1);
                gainNode.gain.setValueAtTime(0.3, this.context.currentTime);
                gainNode.gain.exponentialRampToValueAtTime(0.01, this.context.currentTime + 0.1);
                oscillator.start();
                oscillator.stop(this.context.currentTime + 0.1);
                break;
                
            case 'skill':
                oscillator.frequency.setValueAtTime(400, this.context.currentTime);
                oscillator.frequency.exponentialRampToValueAtTime(800, this.context.currentTime + 0.2);
                gainNode.gain.setValueAtTime(0.3, this.context.currentTime);
                gainNode.gain.exponentialRampToValueAtTime(0.01, this.context.currentTime + 0.3);
                oscillator.start();
                oscillator.stop(this.context.currentTime + 0.3);
                break;
                
            case 'hit':
                oscillator.type = 'square';
                oscillator.frequency.setValueAtTime(100, this.context.currentTime);
                gainNode.gain.setValueAtTime(0.2, this.context.currentTime);
                gainNode.gain.exponentialRampToValueAtTime(0.01, this.context.currentTime + 0.05);
                oscillator.start();
                oscillator.stop(this.context.currentTime + 0.05);
                break;
                
            case 'upgrade':
                oscillator.frequency.setValueAtTime(300, this.context.currentTime);
                oscillator.frequency.setValueAtTime(600, this.context.currentTime + 0.1);
                oscillator.frequency.setValueAtTime(900, this.context.currentTime + 0.2);
                gainNode.gain.setValueAtTime(0.2, this.context.currentTime);
                gainNode.gain.exponentialRampToValueAtTime(0.01, this.context.currentTime + 0.3);
                oscillator.start();
                oscillator.stop(this.context.currentTime + 0.3);
                break;
        }
    }
}

const audioManager = new AudioManager();
