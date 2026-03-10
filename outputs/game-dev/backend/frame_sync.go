package main

import (
	"fmt"
	"sync"
	"time"
)

// ============================================
// 帧同步系统 (Frame Synchronization)
// ============================================

// 帧消息类型
type FrameMsgType int

const (
	FrameMsgTypeInput FrameMsgType = iota // 玩家输入
	FrameMsgTypeState                      // 状态同步
	FrameMsgTypeEvent                      // 事件
	FrameMsgTypeSync                       // 帧同步
)

// 帧消息
type FrameMessage struct {
	Type     FrameMsgType `json:"type"`
	Frame    int64       `json:"frame"`
	PlayerID string      `json:"playerId"`
	Data     interface{} `json:"data"`
	Timestamp int64     `json:"timestamp"`
}

// 玩家输入
type PlayerInput struct {
	Type      string  `json:"type"`      // place_tower/upgrade_tower/sell_tower/start_wave
	TowerID   string  `json:"towerId"`   // 塔ID
	TowerType string  `json:"towerType"` // 塔类型
	X         float64 `json:"x"`         // 位置X
	Y         float64 `json:"y"`         // 位置Y
}

// 帧状态
type FrameState struct {
	Frame    int64           `json:"frame"`
	Time     time.Time       `json:"time"`
	Players  []*PlayerState  `json:"players"`
	Towers   []*TowerState   `json:"towers"`
	Enemies  []*EnemyState   `json:"enemies"`
	Projectiles []*ProjectileState `json:"projectiles"`
	Danmaku  []*DanmakuState `json:"danmaku"`
}

// 玩家状态
type PlayerState struct {
	ID       string `json:"id"`
	Money    int    `json:"money"`
	Lives    int    `json:"lives"`
	Score    int    `json:"score"`
	Ready    bool   `json:"ready"`
}

// 塔状态
type TowerState struct {
	ID       string  `json:"id"`
	Type     string  `json:"type"`
	X        float64 `json:"x"`
	Y        float64 `json:"y"`
	Level    int     `json:"level"`
	TargetID string  `json:"targetId"` // 当前目标
}

// 敌人状态
type EnemyState struct {
	ID        string  `json:"id"`
	Type      string  `json:"type"`
	X         float64 `json:"x"`
	Y         float64 `json:"y"`
	HP        float64 `json:"hp"`
	MaxHP     float64 `json:"maxHp"`
	Progress  float64 `json:"progress"` // 进度 0-1
	Speed     float64 `json:"speed"`
	Frozen    float64 `json:"frozen"`   // 冰冻时间
}

// 投射物状态
type ProjectileState struct {
	ID        string  `json:"id"`
	Type      string  `json:"type"`
	X         float64 `json:"x"`
	Y         float64 `json:"y"`
	TargetID  string  `json:"targetId"`
	Damage    float64 `json:"damage"`
	Speed     float64 `json:"speed"`
}

// 弹幕状态
type DanmakuState struct {
	ID       string  `json:"id"`
	Content  string  `json:"content"`
	X        float64 `json:"x"`
	Y        float64 `json:"y"`
	Speed    float64 `json:"speed"`
	Color    string  `json:"color"`
}

// 帧同步管理器
type FrameSyncManager struct {
	// 配置
	FPS          int         // 目标帧率
	FrameLimit   int64       // 最大帧数 (0=无限)
	ReplayMode   bool        // 回放模式
	
	// 状态
	FrameCount   int64       // 当前帧数
	Running      bool        // 是否运行
	Paused       bool        // 是否暂停
	
	// 房间
	RoomID       string
	Players      map[string]*FramePlayer // playerID -> player
	
	// 输入缓冲
	InputBuffer  map[int64][]*PlayerInput // frame -> inputs
	
	// 状态历史 (用于回放)
	StateHistory map[int64]*FrameState // frame -> state
	
	// 玩家输入队列
	InputQueue chan *FrameMessage
	
	// 控制
	mu          sync.RWMutex
	wg          sync.WaitGroup
}

// 帧玩家
type FramePlayer struct {
	ID       string
	Name     string
	Input    *PlayerInput
	LastFrame int64
	Ready    bool
}

// 新建帧同步管理器
func NewFrameSyncManager(roomID string, fps int) *FrameSyncManager {
	return &FrameSyncManager{
		RoomID:      roomID,
		FPS:         fps,
		FrameLimit:  0, // 无限
		ReplayMode:  false,
		FrameCount:  0,
		Running:    false,
		Paused:     false,
		Players:    make(map[string]*FramePlayer),
		InputBuffer: make(map[int64][]*PlayerInput),
		StateHistory: make(map[int64]*FrameState),
		InputQueue:  make(chan *FrameMessage, 1000),
	}
}

// 添加玩家
func (fsm *FrameSyncManager) AddPlayer(playerID, name string) {
	fsm.mu.Lock()
	defer fsm.mu.Unlock()
	
	fsm.Players[playerID] = &FramePlayer{
		ID:    playerID,
		Name:  name,
		Ready: false,
	}
}

// 移除玩家
func (fsm *FrameSyncManager) RemovePlayer(playerID string) {
	fsm.mu.Lock()
	defer fsm.mu.Unlock()
	
	delete(fsm.Players, playerID)
}

// 玩家准备
func (fsm *FrameSyncManager) SetPlayerReady(playerID string, ready bool) {
	fsm.mu.Lock()
	defer fsm.mu.Unlock()
	
	if player, ok := fsm.Players[playerID]; ok {
		player.Ready = ready
	}
}

// 发送玩家输入
func (fsm *FrameSyncManager) SendInput(playerID string, input *PlayerInput) {
	msg := &FrameMessage{
		Type:      FrameMsgTypeInput,
		Frame:     fsm.FrameCount,
		PlayerID:  playerID,
		Data:      input,
		Timestamp: time.Now().UnixMilli(),
	}
	
	fsm.InputQueue <- msg
}

// 开始同步
func (fsm *FrameSyncManager) Start() {
	fsm.mu.Lock()
	if fsm.Running {
		fsm.mu.Unlock()
		return
	}
	fsm.Running = true
	fsm.Paused = false
	fsm.FrameCount = 0
	fsm.mu.Unlock()
	
	// 启动输入处理协程
	fsm.wg.Add(1)
	go fsm.processInputs()
	
	// 启动帧同步协程
	fsm.wg.Add(1)
	go fsm.frameLoop()
}

// 停止同步
func (fsm *FrameSyncManager) Stop() {
	fsm.mu.Lock()
	fsm.Running = false
	fsm.mu.Unlock()
	
	fsm.wg.Wait()
}

// 暂停
func (fsm *FrameSyncManager) Pause() {
	fsm.mu.Lock()
	defer fsm.mu.Unlock()
	
	fsm.Paused = true
}

// 恢复
func (fsm *FrameSyncManager) Resume() {
	fsm.mu.Lock()
	defer fsm.mu.Unlock()
	
	fsm.Paused = false
}

// 获取当前帧
func (fsm *FrameSyncManager) GetFrame() int64 {
	fsm.mu.RLock()
	defer fsm.mu.RUnlock()
	
	return fsm.FrameCount
}

// 获取状态历史
func (fsm *FrameSyncManager) GetStateHistory(frame int64) *FrameState {
	fsm.mu.RLock()
	defer fsm.mu.RUnlock()
	
	return fsm.StateHistory[frame]
}

// 输入处理循环
func (fsm *FrameSyncManager) processInputs() {
	defer fsm.wg.Done()
	
	ticker := time.NewTicker(time.Millisecond * 16) // ~60fps
	defer ticker.Stop()
	
	for {
		select {
		case msg, ok := <-fsm.InputQueue:
			if !ok {
				return
			}
			
			if msg.Type == FrameMsgTypeInput {
				fsm.mu.Lock()
				fsm.InputBuffer[msg.Frame] = append(fsm.InputBuffer[msg.Frame], msg.Data.(*PlayerInput))
				
				// 更新玩家最后帧
				if player, ok := fsm.Players[msg.PlayerID]; ok {
					player.LastFrame = msg.Frame
					player.Input = msg.Data.(*PlayerInput)
				}
				fsm.mu.Unlock()
			}
			
		case <-ticker.C:
			// 清理旧输入
			fsm.mu.Lock()
			threshold := fsm.FrameCount - 300 // 保留最近5秒
			for frame := range fsm.InputBuffer {
				if frame < threshold {
					delete(fsm.InputBuffer, frame)
				}
			}
			fsm.mu.Unlock()
		}
	}
}

// 帧同步循环
func (fsm *FrameSyncManager) frameLoop() {
	defer fsm.wg.Done()
	
	frameDuration := time.Second / time.Duration(fsm.FPS)
	ticker := time.NewTicker(frameDuration)
	defer ticker.Stop()
	
	for {
		select {
		case <-ticker.C:
			fsm.mu.Lock()
			if !fsm.Running || fsm.Paused {
				fsm.mu.Unlock()
				continue
			}
			
			// 检查帧限制
			if fsm.FrameLimit > 0 && fsm.FrameCount >= fsm.FrameLimit {
				fsm.Running = false
				fsm.mu.Unlock()
				continue
			}
			
			frame := fsm.FrameCount
			fsm.FrameCount++
			fsm.mu.Unlock()
			
			// 获取当前帧输入
			inputs := fsm.getFrameInputs(frame)
			
			// 执行输入
			fsm.executeInputs(inputs)
			
			// 记录状态
			fsm.recordState(frame)
			
			// 广播帧同步消息
			fsm.broadcastFrame(frame)
		}
	}
}

// 获取帧输入
func (fsm *FrameSyncManager) getFrameInputs(frame int64) []*PlayerInput {
	fsm.mu.RLock()
	defer fsm.mu.RUnlock()
	
	return fsm.InputBuffer[frame]
}

// 执行输入
func (fsm *FrameSyncManager) executeInputs(inputs []*PlayerInput) {
	// 这里可以调用游戏引擎的方法处理输入
	for _, input := range inputs {
		switch input.Type {
		case "place_tower":
			// 处理放置塔
			fmt.Printf("Frame %d: Place tower %s at (%.2f, %.2f)\n", 
				fsm.FrameCount, input.TowerType, input.X, input.Y)
			
		case "upgrade_tower":
			// 处理升级塔
			fmt.Printf("Frame %d: Upgrade tower %s\n", fsm.FrameCount, input.TowerID)
			
		case "sell_tower":
			// 处理出售塔
			fmt.Printf("Frame %d: Sell tower %s\n", fsm.FrameCount, input.TowerID)
		}
	}
}

// 记录状态
func (fsm *FrameSyncManager) recordState(frame int64) {
	fsm.mu.Lock()
	defer fsm.mu.Unlock()
	
	state := &FrameState{
		Frame: frame,
		Time:  time.Now(),
	}
	
	// 记录玩家状态
	for _, player := range fsm.Players {
		state.Players = append(state.Players, &PlayerState{
			ID:     player.ID,
			Name:   player.Name,
			Ready:  player.Ready,
			Money:  100,  // TODO: 从游戏状态获取
			Lives:  20,
			Score:  0,
		})
	}
	
	fsm.StateHistory[frame] = state
}

// 广播帧
func (fsm *FrameSyncManager) broadcastFrame(frame int64) {
	fsm.mu.RLock()
	state := fsm.StateHistory[frame]
	fsm.mu.RUnlock()
	
	if state == nil {
		return
	}
	
	// TODO: 广播给所有玩家
	msg := &FrameMessage{
		Type:   FrameMsgTypeSync,
		Frame:  frame,
		Data:   state,
		Timestamp: time.Now().UnixMilli(),
	}
	
	_ = msg // 实际实现需要广播
}

// ============================================
// 回放系统 (Replay System)
// ============================================

// 回放数据
type ReplayData struct {
	RoomID      string          `json:"roomId"`
	FPS         int             `json:"fps"`
	StartTime   time.Time       `json:"startTime"`
	EndTime     time.Time       `json:"endTime"`
	TotalFrames int64           `json:"totalFrames"`
	Players     []*ReplayPlayer `json:"players"`
	Frames      []*ReplayFrame  `json:"frames"`
}

// 回放玩家
type ReplayPlayer struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// 回放帧
type ReplayFrame struct {
	Frame   int64         `json:"frame"`
	Time    time.Time    `json:"time"`
	Inputs  []*PlayerInput `json:"inputs"`
	State   *FrameState  `json:"state"`
}

// 回放管理器
type ReplayManager struct {
	Replays map[string]*ReplayData // roomID -> replay
	mu      sync.RWMutex
}

// 新建回放管理器
func NewReplayManager() *ReplayManager {
	return &ReplayManager{
		Replays: make(map[string]*ReplayData),
	}
}

// 开始录制
func (rm *ReplayManager) StartRecording(roomID string, fps int, players []*ReplayPlayer) *ReplayData {
	rm.mu.Lock()
	defer rm.mu.Unlock()
	
	replay := &ReplayData{
		RoomID:    roomID,
		FPS:       fps,
		StartTime: time.Now(),
		Players:   players,
		Frames:    make([]*ReplayFrame, 0),
	}
	
	rm.Replays[roomID] = replay
	return replay
}

// 添加帧
func (rm *ReplayManager) AddFrame(roomID string, frame *ReplayFrame) {
	rm.mu.Lock()
	defer rm.mu.Unlock()
	
	if replay, ok := rm.Replays[roomID]; ok {
		replay.Frames = append(replay.Frames, frame)
		replay.TotalFrames = int64(len(replay.Frames))
	}
}

// 结束录制
func (rm *ReplayManager) StopRecording(roomID string) *ReplayData {
	rm.mu.Lock()
	defer rm.mu.Unlock()
	
	if replay, ok := rm.Replays[roomID]; ok {
		replay.EndTime = time.Now()
		return replay
	}
	
	return nil
}

// 获取回放
func (rm *ReplayManager) GetReplay(roomID string) *ReplayData {
	rm.mu.RLock()
	defer rm.mu.RUnlock()
	
	return rm.Replays[roomID]
}

// 播放回放
func (rd *ReplayManager) Play(roomID string, speed float64) error {
	replay := rd.GetReplay(roomID)
	if replay == nil {
		return fmt.Errorf("replay not found")
	}
	
	frameDuration := time.Second / time.Duration(float64(replay.FPS)*speed)
	
	for i, frame := range replay.Frames {
		fmt.Printf("Playing frame %d/%d\n", i+1, len(replay.Frames))
		
		// 渲染状态
		_ = frame.State
		
		// 等待下一帧
		time.Sleep(frameDuration)
	}
	
	return nil
}

// ============================================
// 帧同步房间
// ============================================

// 帧同步房间
type FrameSyncRoom struct {
	RoomID       string
	FrameSync    *FrameSyncManager
	Replay       *ReplayManager
	LiveRoom     *LiveRoom
}

// 新建帧同步房间
func NewFrameSyncRoom(roomID string, fps int) *FrameSyncRoom {
	return &FrameSyncRoom{
		RoomID:    roomID,
		FrameSync: NewFrameSyncManager(roomID, fps),
		Replay:    NewReplayManager(),
	}
}

// 初始化回放
func (fsr *FrameSyncRoom) InitReplay(players []*ReplayPlayer) {
	fsr.Replay.StartRecording(fsr.RoomID, fsr.FrameSync.FPS, players)
}

// 添加帧到回放
func (fsr *FrameSyncRoom) AddReplayFrame(frame int64, inputs []*PlayerInput) {
	replayFrame := &ReplayFrame{
		Frame:  frame,
		Time:   time.Now(),
		Inputs: inputs,
	}
	
	fsr.Replay.AddFrame(fsr.RoomID, replayFrame)
}

// 结束回放
func (fsr *FrameSyncRoom) EndReplay() *ReplayData {
	return fsr.Replay.StopRecording(fsr.RoomID)
}

// ============================================
// 全局帧同步管理
// ============================================

var (
	GlobalReplayManager *ReplayManager
)

func init() {
	GlobalReplayManager = NewReplayManager()
}
