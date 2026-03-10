using System;
using System.Collections.Generic;
using System.Threading;
using System.Threading.Tasks;

namespace DanmakuGameClient
{
    /// <summary>
    /// 观战系统客户端
    /// </summary>
    public class SpectatorClient
    {
        private readonly GameClient _gameClient;
        private SpectatorRoom _currentRoom;
        private List<Spectator> _spectators = new List<Spectator>();
        private bool _isSpectating;
        private int _cameraAngle;
        private bool _showDanmaku = true;
        private bool _showStats = true;

        // 观战状态回调
        public Action<SpectatorRoom> OnSpectatorRoomCreated;
        public Action<Spectator> OnJoinedSpectator;
        public Action<Spectator> OnSpectatorLeft;
        public Action<List<Spectator>> OnSpectatorsUpdated;
        public Action<BattleReplay> OnReplayReceived;
        public Action<string> OnError;
        public Action OnSpectatorLeftRoom;

        public SpectatorClient(GameClient gameClient)
        {
            _gameClient = gameClient;
            RegisterHandlers();
        }

        #region 消息注册

        private void RegisterHandlers()
        {
            _gameClient.RegisterHandler<SCSpectatorRoomCreated>(OnSpectatorRoomCreated);
            _gameClient.RegisterHandler<SCJoinSpectator>(OnJoinSpectator);
            _gameClient.RegisterHandler<SCSpectatorLeft>(OnSpectatorLeft);
            _gameClient.RegisterHandler<SCSpectatorsUpdated>(OnSpectatorsUpdated);
            _gameClient.RegisterHandler<SCSpectatorViewUpdate>(OnSpectatorViewUpdate);
            _gameClient.RegisterHandler<SCBattleReplay>(OnBattleReplay);
            _gameClient.RegisterHandler<SCSpectatorError>(OnSpectatorError);
            _gameClient.RegisterHandler<SCSpectatorLeftRoom>(OnSpectatorLeftRoom);
        }

        #endregion

        #region 观战房间操作

        /// <summary>
        /// 创建观战房间
        /// </summary>
        public async Task<SpectatorRoom> CreateSpectatorRoomAsync(string battleRoomId)
        {
            var req = new CSCreateSpectatorRoom
            {
                BattleRoomId = battleRoomId
            };

            var result = await _gameClient.SendRequestAsync<CSCreateSpectatorRoom, SCSpectatorRoomCreated>(req);
            return result?.Room;
        }

        /// <summary>
        /// 加入观战
        /// </summary>
        public async Task<Spectator> JoinSpectatorAsync(string battleRoomId)
        {
            var req = new CSJoinSpectator
            {
                BattleRoomId = battleRoomId
            };

            var result = await _gameClient.SendRequestAsync<CSJoinSpectator, SCJoinSpectator>(req);
            _isSpectating = result != null;
            return result?.Spectator;
        }

        /// <summary>
        /// 离开观战
        /// </summary>
        public async Task LeaveSpectatorAsync()
        {
            var req = new CSLeaveSpectator();
            await _gameClient.SendRequestAsync<CSLeaveSpectator, SCLeaveSpectator>(req);
            _isSpectating = false;
            _currentRoom = null;
        }

        /// <summary>
        /// 获取观战房间信息
        /// </summary>
        public async Task<SpectatorRoom> GetSpectatorRoomAsync(string battleRoomId)
        {
            var req = new CSGetSpectatorRoom
            {
                BattleRoomId = battleRoomId
            };

            var result = await _gameClient.SendRequestAsync<CSGetSpectatorRoom, SCSpectatorRoom>(req);
            return result?.Room;
        }

        /// <summary>
        /// 获取观战列表
        /// </summary>
        public async Task<List<Spectator>> GetSpectatorsAsync(string battleRoomId)
        {
            var req = new CSGetSpectators
            {
                BattleRoomId = battleRoomId
            };

            var result = await _gameClient.SendRequestAsync<CSGetSpectators, SCSpectators>(req);
            return result?.Spectators ?? new List<Spectator>();
        }

        #endregion

        #region 视角控制

        /// <summary>
        /// 切换视角
        /// </summary>
        public async Task SetCameraAngleAsync(int angle)
        {
            var req = new CSSpectatorCameraAngle
            {
                Angle = angle
            };

            await _gameClient.SendRequestAsync<CSSpectatorCameraAngle, SCSpectatorCameraAngle>(req);
            _cameraAngle = angle;
        }

        /// <summary>
        /// 切换到自由视角
        /// </summary>
        public async Task SetFreeCameraAsync()
        {
            await SetCameraAngleAsync(1);
        }

        /// <summary>
        /// 切换到跟随视角
        /// </summary>
        public async Task SetFollowCameraAsync()
        {
            await SetCameraAngleAsync(0);
        }

        /// <summary>
        /// 切换到上帝视角
        /// </summary>
        public async Task SetGodViewCameraAsync()
        {
            await SetCameraAngleAsync(2);
        }

        /// <summary>
        /// 跟随指定玩家
        /// </summary>
        public async Task FollowPlayerAsync(string userId)
        {
            var req = new CSSpectatorFollowPlayer
            {
                UserId = userId
            };

            await _gameClient.SendRequestAsync<CSSpectatorFollowPlayer, SCSpectatorFollowPlayer>(req);
        }

        #endregion

        #region 显示设置

        /// <summary>
        /// 设置弹幕显示
        /// </summary>
        public async Task SetDanmakuVisibleAsync(bool visible)
        {
            var req = new CSSpectatorDanmaku
            {
                Visible = visible
            };

            await _gameClient.SendRequestAsync<CSSpectatorDanmaku, SCSpectatorDanmaku>(req);
            _showDanmaku = visible;
        }

        /// <summary>
        /// 设置统计显示
        /// </summary>
        public async Task SetStatsVisibleAsync(bool visible)
        {
            var req = new CSSpectatorStats
            {
                Visible = visible
            };

            await _gameClient.SendRequestAsync<CSSpectatorStats, SCSpectatorStats>(req);
            _showStats = visible;
        }

        #endregion

        #region 录像回放

        /// <summary>
        /// 获取战斗回放
        /// </summary>
        public async Task<BattleReplay> GetBattleReplayAsync(string battleRoomId)
        {
            var req = new CSGetBattleReplay
            {
                BattleRoomId = battleRoomId
            };

            var result = await _gameClient.SendRequestAsync<CSGetBattleReplay, SCBattleReplay>(req);
            return result?.Replay;
        }

        /// <summary>
        /// 开始录制观战
        /// </summary>
        public async Task StartRecordingAsync(string battleRoomId)
        {
            var req = new CSStartSpectatorRecording
            {
                BattleRoomId = battleRoomId
            };

            await _gameClient.SendRequestAsync<CSStartSpectatorRecording, SCStartSpectatorRecording>(req);
        }

        /// <summary>
        /// 停止录制观战
        /// </summary>
        public async Task StopRecordingAsync(string battleRoomId)
        {
            var req = new CSStopSpectatorRecording
            {
                BattleRoomId = battleRoomId
            };

            await _gameClient.SendRequestAsync<CSStopSpectatorRecording, SCStopSpectatorRecording>(req);
        }

        #endregion

        #region 消息处理

        private void OnSpectatorRoomCreated(SCSpectatorRoomCreated msg)
        {
            _currentRoom = msg.Room;
            OnSpectatorRoomCreated?.Invoke(msg.Room);
        }

        private void OnJoinSpectator(SCJoinSpectator msg)
        {
            var spectator = msg.Spectator;
            _currentRoom = msg.Room;
            _isSpectating = true;
            OnJoinedSpectator?.Invoke(spectator);
        }

        private void OnSpectatorLeft(SCSpectatorLeft msg)
        {
            OnSpectatorLeftRoom?.Invoke();
            _isSpectating = false;
        }

        private void OnSpectatorsUpdated(SCSpectatorsUpdated msg)
        {
            _spectators = msg.Spectators;
            OnSpectatorsUpdated?.Invoke(msg.Spectators);
        }

        private void OnSpectatorViewUpdate(SCSpectatorViewUpdate msg)
        {
            // 更新观战视角数据
            UpdateView(msg.ViewData);
        }

        private void OnBattleReplay(SCBattleReplay msg)
        {
            OnReplayReceived?.Invoke(msg.Replay);
        }

        private void OnSpectatorError(SCSpectatorError msg)
        {
            OnError?.Invoke(msg.ErrorMessage);
        }

        private void OnSpectatorLeftRoom(SCSpectatorLeftRoom msg)
        {
            _isSpectating = false;
            _currentRoom = null;
            OnSpectatorLeftRoom?.Invoke();
        }

        #endregion

        #region 辅助方法

        private void UpdateView(SpectatorViewData viewData)
        {
            // 更新观战画面数据
        }

        /// <summary>
        /// 获取当前观战状态
        /// </summary>
        public bool IsSpectating => _isSpectating;

        /// <summary>
        /// 获取当前观战房间
        /// </summary>
        public SpectatorRoom CurrentRoom => _currentRoom;

        /// <summary>
        /// 获取观战者列表
        /// </summary>
        public List<Spectator> Spectators => _spectators;

        /// <summary>
        /// 获取当前视角
        /// </summary>
        public int CameraAngle => _cameraAngle;

        /// <summary>
        /// 弹幕是否显示
        /// </summary>
        public bool ShowDanmaku => _showDanmaku;

        /// <summary>
        /// 统计是否显示
        /// </summary>
        public bool ShowStats => _showStats;

        #endregion
    }

    #region 消息定义

    // 请求消息
    public class CSCreateSpectatorRoom
    {
        public string BattleRoomId { get; set; }
    }

    public class CSJoinSpectator
    {
        public string BattleRoomId { get; set; }
    }

    public class CSLeaveSpectator { }

    public class CSGetSpectatorRoom
    {
        public string BattleRoomId { get; set; }
    }

    public class CSGetSpectators
    {
        public string BattleRoomId { get; set; }
    }

    public class CSSpectatorCameraAngle
    {
        public int Angle { get; set; }
    }

    public class CSSpectatorFollowPlayer
    {
        public string UserId { get; set; }
    }

    public class CSSpectatorDanmaku
    {
        public bool Visible { get; set; }
    }

    public class CSSpectatorStats
    {
        public bool Visible { get; set; }
    }

    public class CSGetBattleReplay
    {
        public string BattleRoomId { get; set; }
    }

    public class CSStartSpectatorRecording
    {
        public string BattleRoomId { get; set; }
    }

    public class CSStopSpectatorRecording
    {
        public string BattleRoomId { get; set; }
    }

    // 响应消息
    public class SCSpectatorRoomCreated
    {
        public SpectatorRoom Room { get; set; }
    }

    public class SCJoinSpectator
    {
        public Spectator Spectator { get; set; }
        public SpectatorRoom Room { get; set; }
    }

    public class SCLeaveSpectator { }

    public class SCSpectatorRoom
    {
        public SpectatorRoom Room { get; set; }
    }

    public class SCSpectators
    {
        public List<Spectator> Spectators { get; set; }
    }

    public class SCSpectatorsUpdated
    {
        public List<Spectator> Spectators { get; set; }
    }

    public class SCSpectatorViewUpdate
    {
        public SpectatorViewData ViewData { get; set; }
    }

    public class SCBattleReplay
    {
        public BattleReplay Replay { get; set; }
    }

    public class SCSpectatorError
    {
        public string ErrorMessage { get; set; }
    }

    public class SCSpectatorLeftRoom { }

    public class SCSpectatorCameraAngle { }

    public class SCSpectatorFollowPlayer { }

    public class SCSpectatorDanmaku { }

    public class SCSpectatorStats { }

    public class SCStartSpectatorRecording { }

    public class SCStopSpectatorRecording { }

    #endregion

    #region 数据结构

    public class Spectator
    {
        public string UserId { get; set; }
        public string UserName { get; set; }
        public string UserAvatar { get; set; }
        public string RoomId { get; set; }
        public int State { get; set; }
        public long JoinTime { get; set; }
        public int CameraAngle { get; set; }
        public bool ShowDanmaku { get; set; }
        public bool ShowStats { get; set; }
        public string FollowTargetUserId { get; set; }
    }

    public class SpectatorRoom
    {
        public string RoomId { get; set; }
        public string BattleRoomId { get; set; }
        public string HostUserId { get; set; }
        public List<Spectator> Spectators { get; set; }
        public int MaxSpectators { get; set; }
        public bool IsRecording { get; set; }
        public long ViewCount { get; set; }
        public long CreatedAt { get; set; }
    }

    public class SpectatorViewData
    {
        public string BattleRoomId { get; set; }
        public int GameState { get; set; }
        public List<PlayerViewInfo> Players { get; set; }
        public long GameTime { get; set; }
        public int Wave { get; set; }
        public List<DanmakuInfo> DanmakuList { get; set; }
    }

    public class PlayerViewInfo
    {
        public string UserId { get; set; }
        public string UserName { get; set; }
        public int PlayerIndex { get; set; }
        public float PosX { get; set; }
        public float PosY { get; set; }
        public float Health { get; set; }
        public float MaxHealth { get; set; }
        public int Score { get; set; }
        public int Combo { get; set; }
        public List<string> ActiveBuffs { get; set; }
    }

    public class DanmakuInfo
    {
        public string UserId { get; set; }
        public string UserName { get; set; }
        public string Content { get; set; }
        public int Color { get; set; }
        public long Timestamp { get; set; }
    }

    public class BattleReplay
    {
        public string ReplayId { get; set; }
        public string BattleRoomId { get; set; }
        public long StartTime { get; set; }
        public long EndTime { get; set; }
        public List<ReplayFrame> Frames { get; set; }
        public ReplayMetadata Metadata { get; set; }
    }

    public class ReplayFrame
    {
        public long Timestamp { get; set; }
        public int FrameType { get; set; }
        public string Data { get; set; }
    }

    public class ReplayMetadata
    {
        public List<PlayerInfo> Players { get; set; }
        public int WaveCount { get; set; }
        public string WinnerUserId { get; set; }
        public int Duration { get; set; }
    }

    public class PlayerInfo
    {
        public string UserId { get; set; }
        public string UserName { get; set; }
        public int PlayerIndex { get; set; }
        public int FinalScore { get; set; }
        public int Rank { get; set; }
    }

    #endregion
}
