// ============================================
// Unity 客户端示例 - 弹幕游戏
// 对接后端 WebSocket + Protobuf 协议
// ============================================

using System;
using System.Collections.Generic;
using System.Threading;
using UnityEngine;
using UnityEngine.UI;

namespace DanmakuGame.Client
{
    // ============================================
    // 消息ID枚举 (对应 server/proto/game.proto)
    // ============================================
    public enum MsgID : ushort
    {
        // 房间相关
        CS_CREATE_ROOM = 1001,
        SC_CREATE_ROOM = 1002,
        CS_JOIN_ROOM = 1003,
        SC_JOIN_ROOM = 1004,
        CS_LEAVE_ROOM = 1005,
        SC_LEAVE_ROOM = 1006,
        CS_LIST_ROOMS = 1007,
        SC_LIST_ROOMS = 1008,
        
        // 游戏相关
        CS_START_GAME = 2001,
        SC_START_GAME = 2002,
        CS_PLACE_TOWER = 2003,
        SC_PLACE_TOWER = 2004,
        CS_UPGRADE_TOWER = 2005,
        SC_UPGRADE_TOWER = 2006,
        CS_SELL_TOWER = 2007,
        SC_SELL_TOWER = 2008,
        
        // 帧同步
        CS_FRAME_INPUT = 3001,
        SC_FRAME_SYNC = 3002,
        
        // 礼物弹幕
        CS_SEND_GIFT = 4001,
        SC_RECEIVE_GIFT = 4002,
        CS_SEND_DANMAKU = 4003,
        SC_RECEIVE_DANMAKU = 4004,
        
        // 心跳
        CS_HEARTBEAT = 5001,
        SC_HEARTBEAT = 5002,
        
        // 错误响应
        SC_ERROR = 9001,
    }

    // ============================================
    // 错误码定义
    // ============================================
    public enum ErrorCode : ushort
    {
        Success = 0,
        InvalidToken = 1001,
        RoomNotFound = 2001,
        RoomFull = 2002,
        RoomAlreadyStarted = 2003,
        NotRoomOwner = 2004,
        InvalidPosition = 3001,
        NotEnoughMoney = 3002,
        TowerNotFound = 3003,
    }

    // ============================================
    // 基础消息结构
    // ============================================
    [Serializable]
    public class WSMessage
    {
        public ushort msg_id;
        public string data; // JSON 序列化后的消息体
    }

    // ============================================
    // 房间相关消息体
    // ============================================
    [Serializable]
    public class CreateRoomReq
    {
        public string player_id;
        public string player_name;
        public string mode; // "classic", "boss", "endless"
        public int max_players = 2;
    }

    [Serializable]
    public class CreateRoomResp
    {
        public string room_id;
        public string host_id;
        public string mode;
        public int max_players;
        public long created_at;
    }

    [Serializable]
    public class JoinRoomReq
    {
        public string room_id;
        public string player_id;
        public string player_name;
    }

    [Serializable]
    public class JoinRoomResp
    {
        public bool success;
        public string room_id;
        public string[] player_ids;
        public string[] player_names;
        public string mode;
    }

    [Serializable]
    public class RoomInfo
    {
        public string room_id;
        public string host_id;
        public string mode;
        public int current_players;
        public int max_players;
        public string status; // "waiting", "playing", "ended"
    }

    [Serializable]
    public class ListRoomsResp
    {
        public RoomInfo[] rooms;
    }

    // ============================================
    // 游戏相关消息体
    // ============================================
    [Serializable]
    public class PlaceTowerReq
    {
        public string room_id;
        public string player_id;
        public string tower_type; // "arrow", "cannon", "ice", "lightning", "heal"
        public float x;
        public float y;
    }

    [Serializable]
    public class PlaceTowerResp
    {
        public bool success;
        public string tower_id;
        public string tower_type;
        public float x;
        public float y;
        public int level;
        public int cost;
    }

    [Serializable]
    public class UpgradeTowerReq
    {
        public string room_id;
        public string player_id;
        public string tower_id;
    }

    [Serializable]
    public class UpgradeTowerResp
    {
        public bool success;
        public string tower_id;
        public int new_level;
        public int cost;
    }

    [Serializable]
    public class GameStartResp
    {
        public string room_id;
        public int wave; // 当前波次
        public long start_time;
        public Dictionary<string, int> players_money;
    }

    // ============================================
    // 帧同步消息体
    // ============================================
    [Serializable]
    public class FrameInput
    {
        public string player_id;
        public ushort frame_id;
        public string action; // "place_tower", "upgrade_tower", "sell_tower"
        public string[] args;
    }

    [Serializable]
    public class FrameSync
    {
        public ushort frame_id;
        public FrameInput[] inputs;
        public GameState state;
    }

    [Serializable]
    public class GameState
    {
        public int wave;
        public long timestamp;
        public Dictionary<string, int> players_money;
        public Tower[] towers;
        public Enemy[] enemies;
        public Projectile[] projectiles;
    }

    [Serializable]
    public class Tower
    {
        public string id;
        public string owner_id;
        public string type;
        public int level;
        public float x;
        public float y;
        public float range;
        public float attack_speed;
        public int damage;
        public int cost;
    }

    [Serializable]
    public class Enemy
    {
        public string id;
        public string type;
        public float x;
        public float y;
        public int hp;
        public int max_hp;
        public float speed;
        public float progress; // 0-1 进度
        public bool slowed;
    }

    [Serializable]
    public class Projectile
    {
        public string id;
        public string tower_id;
        public string target_id;
        public float x;
        public float y;
        public float vx;
        public float vy;
        public int damage;
        public string effect; // "normal", "ice", "lightning"
    }

    // ============================================
    // 礼物弹幕消息体
    // ============================================
    [Serializable]
    public class SendGiftReq
    {
        public string room_id;
        public string sender_id;
        public string sender_name;
        public string gift_type; // "coin", "star", "rocket", "car", "plane", "bang"
    }

    [Serializable]
    public class ReceiveGiftResp
    {
        public string sender_id;
        public string sender_name;
        public string gift_type;
        public int value;
        public long timestamp;
    }

    [Serializable]
    public class SendDanmakuReq
    {
        public string room_id;
        public string sender_id;
        public string sender_name;
        public string content;
    }

    [Serializable]
    public class ReceiveDanmakuResp
    {
        public string sender_id;
        public string sender_name;
        public string content;
        public string color; // "white", "red", "gold", "rainbow"
        public long timestamp;
    }

    // ============================================
    // WebSocket 客户端管理器
    // ============================================
    public class GameClient : MonoBehaviour
    {
        [Header("服务器配置")]
        [SerializeField] private string serverUrl = "ws://localhost:8080/ws";
        [SerializeField] private float heartbeatInterval = 30f;

        [Header("调试")]
        [SerializeField] private bool enableLog = true;

        // 状态
        public bool IsConnected { get; private set; }
        public string PlayerId { get; private set; }
        public string PlayerName { get; private set; }
        public string CurrentRoomId { get; private set; }

        // 回调
        public Action<CreateRoomResp> OnCreateRoom;
        public Action<JoinRoomResp> OnJoinRoom;
        public Action<RoomInfo[]> OnListRooms;
        public Action<GameStartResp> OnGameStart;
        public Action<FrameSync> OnFrameSync;
        public Action<ReceiveGiftResp> OnReceiveGift;
        public Action<ReceiveDanmakuResp> OnReceiveDanmaku;
        public Action<ErrorCode, string> OnError;

        // 内部组件
        private WebSocketSharp.WebSocket _ws;
        private Timer _heartbeatTimer;
        private Queue<string> _sendQueue = new Queue<string>();
        private object _queueLock = new object();
        private bool _isSending;

        // ============================================
        // 生命周期
        // ============================================
        private void Awake()
        {
            DontDestroyOnLoad(gameObject);
        }

        private void Start()
        {
            Connect();
        }

        private void OnDestroy()
        {
            Disconnect();
        }

        private void Update()
        {
            // 处理发送队列
            ProcessSendQueue();
        }

        // ============================================
        // 连接管理
        // ============================================
        public void Connect()
        {
            if (_ws != null && _ws.ReadyState == WebSocketSharp.WebSocketState.Open)
            {
                Log("Already connected");
                return;
            }

            _ws = new WebSocketSharp.WebSocket(serverUrl);
            _ws.OnOpen += OnWebSocketOpen;
            _ws.OnMessage += OnWebSocketMessage;
            _ws.OnError += OnWebSocketError;
            _ws.OnClose += OnWebSocketClose;
            
            _ws.ConnectAsync();
        }

        public void Disconnect()
        {
            _heartbeatTimer?.Dispose();
            
            if (_ws != null)
            {
                _ws.CloseAsync();
                _ws = null;
            }
            
            IsConnected = false;
        }

        private void OnWebSocketOpen(object sender, EventArgs e)
        {
            IsConnected = true;
            Log("Connected to server");
            
            // 启动心跳
            _heartbeatTimer = new Timer(SendHeartbeat, null, 
                TimeSpan.FromSeconds(heartbeatInterval), 
                TimeSpan.FromSeconds(heartbeatInterval));
        }

        private void OnWebSocketClose(object sender, WebSocketSharp.CloseEventArgs e)
        {
            IsConnected = false;
            Log($"Disconnected: {e.Reason}");
        }

        private void OnWebSocketError(object sender, WebSocketSharp.ErrorEventArgs e)
        {
            Log($"WebSocket Error: {e.Message}");
        }

        // ============================================
        // 消息处理
        // ============================================
        private void OnWebSocketMessage(object sender, WebSocketSharp.MessageEventArgs e)
        {
            try
            {
                var msg = JsonUtility.FromJson<WSMessage>(e.Data);
                HandleMessage(msg);
            }
            catch (Exception ex)
            {
                Log($"Parse error: {ex.Message}");
            }
        }

        private void HandleMessage(WSMessage msg)
        {
            switch ((MsgID)msg.msg_id)
            {
                case MsgID.SC_CREATE_ROOM:
                    var createResp = JsonUtility.FromJson<CreateRoomResp>(msg.data);
                    CurrentRoomId = createResp.room_id;
                    OnCreateRoom?.Invoke(createResp);
                    break;

                case MsgID.SC_JOIN_ROOM:
                    var joinResp = JsonUtility.FromJson<JoinRoomResp>(msg.data);
                    if (joinResp.success)
                    {
                        CurrentRoomId = joinResp.room_id;
                    }
                    OnJoinRoom?.Invoke(joinResp);
                    break;

                case MsgID.SC_LIST_ROOMS:
                    var listResp = JsonUtility.FromJson<ListRoomsResp>(msg.data);
                    OnListRooms?.Invoke(listResp.rooms);
                    break;

                case MsgID.SC_START_GAME:
                    var startResp = JsonUtility.FromJson<GameStartResp>(msg.data);
                    OnGameStart?.Invoke(startResp);
                    break;

                case MsgID.SC_FRAME_SYNC:
                    var frameSync = JsonUtility.FromJson<FrameSync>(msg.data);
                    OnFrameSync?.Invoke(frameSync);
                    break;

                case MsgID.SC_RECEIVE_GIFT:
                    var gift = JsonUtility.FromJson<ReceiveGiftResp>(msg.data);
                    OnReceiveGift?.Invoke(gift);
                    break;

                case MsgID.SC_RECEIVE_DANMAKU:
                    var danmaku = JsonUtility.FromJson<ReceiveDanmakuResp>(msg.data);
                    OnReceiveDanmaku?.Invoke(danmaku);
                    break;

                case MsgID.SC_ERROR:
                    var error = JsonUtility.FromJson<ErrorResponse>(msg.data);
                    OnError?.Invoke((ErrorCode)error.code, error.message);
                    break;

                case MsgID.SC_HEARTBEAT:
                    Log("Heartbeat received");
                    break;
            }
        }

        [Serializable]
        private class ErrorResponse
        {
            public ushort code;
            public string message;
        }

        // ============================================
        // 发送消息
        // ============================================
        private void SendMessage(ushort msgId, object data)
        {
            if (!IsConnected)
            {
                Log("Not connected");
                return;
            }

            var msg = new WSMessage
            {
                msg_id = msgId,
                data = JsonUtility.ToJson(data)
            };

            lock (_queueLock)
            {
                _sendQueue.Enqueue(JsonUtility.ToJson(msg));
            }
        }

        private void ProcessSendQueue()
        {
            if (_isSending || _sendQueue.Count == 0) return;

            _isSending = true;
            
            lock (_queueLock)
            {
                while (_sendQueue.Count > 0)
                {
                    var data = _sendQueue.Dequeue();
                    _ws.SendAsync(data, null, null);
                }
            }
            
            _isSending = false;
        }

        // ============================================
        // 心跳
        // ============================================
        private void SendHeartbeat(object state)
        {
            SendMessage((ushort)MsgID.CS_HEARTBEAT, new { });
        }

        // ============================================
        // 房间操作 API
        // ============================================
        public void CreateRoom(string playerId, string playerName, string mode)
        {
            PlayerId = playerId;
            PlayerName = playerName;
            
            var req = new CreateRoomReq
            {
                player_id = playerId,
                player_name = playerName,
                mode = mode
            };
            
            SendMessage((ushort)MsgID.CS_CREATE_ROOM, req);
        }

        public void JoinRoom(string roomId, string playerId, string playerName)
        {
            PlayerId = playerId;
            PlayerName = playerName;
            
            var req = new JoinRoomReq
            {
                room_id = roomId,
                player_id = playerId,
                player_name = playerName
            };
            
            SendMessage((ushort)MsgID.CS_JOIN_ROOM, req);
        }

        public void LeaveRoom()
        {
            if (string.IsNullOrEmpty(CurrentRoomId)) return;
            
            SendMessage((ushort)MsgID.CS_LEAVE_ROOM, new 
            { 
                room_id = CurrentRoomId,
                player_id = PlayerId 
            });
            
            CurrentRoomId = null;
        }

        public void ListRooms()
        {
            SendMessage((ushort)MsgID.CS_LIST_ROOMS, new { });
        }

        // ============================================
        // 游戏操作 API
        // ============================================
        public void PlaceTower(string towerType, float x, float y)
        {
            var req = new PlaceTowerReq
            {
                room_id = CurrentRoomId,
                player_id = PlayerId,
                tower_type = towerType,
                x = x,
                y = y
            };
            
            SendMessage((ushort)MsgID.CS_PLACE_TOWER, req);
        }

        public void UpgradeTower(string towerId)
        {
            var req = new UpgradeTowerReq
            {
                room_id = CurrentRoomId,
                player_id = PlayerId,
                tower_id = towerId
            };
            
            SendMessage((ushort)MsgID.CS_UPGRADE_TOWER, req);
        }

        public void SellTower(string towerId)
        {
            SendMessage((ushort)MsgID.CS_SELL_TOWER, new 
            {
                room_id = CurrentRoomId,
                player_id = PlayerId,
                tower_id = towerId
            });
        }

        // ============================================
        // 礼物弹幕 API
        // ============================================
        public void SendGift(string giftType)
        {
            var req = new SendGiftReq
            {
                room_id = CurrentRoomId,
                sender_id = PlayerId,
                sender_name = PlayerName,
                gift_type = giftType
            };
            
            SendMessage((ushort)MsgID.CS_SEND_GIFT, req);
        }

        public void SendDanmaku(string content)
        {
            var req = new SendDanmakuReq
            {
                room_id = CurrentRoomId,
                sender_id = PlayerId,
                sender_name = PlayerName,
                content = content
            };
            
            SendMessage((ushort)MsgID.CS_SEND_DANMAKU, req);
        }

        // ============================================
        // 调试
        // ============================================
        private void Log(string msg)
        {
            if (enableLog)
            {
                Debug.Log($"[GameClient] {msg}");
            }
        }
    }

    // ============================================
    // Unity 面板控制器示例
    // ============================================
    public class GameController : MonoBehaviour
    {
        [SerializeField] private GameClient client;
        [SerializeField] private Text moneyText;
        [SerializeField] private Text waveText;
        [SerializeField] private Transform towerContainer;
        [SerializeField] private Transform enemyContainer;
        [SerializeField] private GameObject towerPrefab;

        private Dictionary<string, GameObject> _towerObjects = new Dictionary<string, GameObject>();
        private Dictionary<string, GameObject> _enemyObjects = new Dictionary<string, GameObject>();

        private void Start()
        {
            // 设置回调
            client.OnGameStart += OnGameStart;
            client.OnFrameSync += OnFrameSync;
            client.OnReceiveGift += OnReceiveGift;
            client.OnReceiveDanmaku += OnReceiveDanmaku;
            client.OnError += OnError;
            
            // 模拟玩家登录
            var playerId = "player_" + UnityEngine.Random.Range(1000, 9999);
            client.CreateRoom(playerId, "玩家" + playerId, "classic");
        }

        private void OnGameStart(GameStartResp resp)
        {
            Debug.Log($"Game started! Wave: {resp.wave}");
        }

        private void OnFrameSync(FrameSync sync)
        {
            // 更新金币
            if (sync.state.players_money.TryGetValue(client.PlayerId, out int money))
            {
                moneyText.text = $"金币: {money}";
            }
            
            waveText.text = $"波次: {sync.state.wave}";
            
            // 更新塔
            UpdateTowers(sync.state.towers);
            
            // 更新敌人
            UpdateEnemies(sync.state.enemies);
        }

        private void UpdateTowers(Tower[] towers)
        {
            if (towers == null) return;

            // 移除不存在的塔
            var towerIds = new HashSet<string>();
            foreach (var tower in towers)
            {
                towerIds.Add(tower.id);
            }

            foreach (var id in _towerObjects.Keys)
            {
                if (!towerIds.Contains(id))
                {
                    Destroy(_towerObjects[id]);
                    _towerObjects.Remove(id);
                }
            }

            // 更新/创建塔
            foreach (var tower in towers)
            {
                if (_towerObjects.TryGetValue(tower.id, out var obj))
                {
                    // 更新位置和等级
                    obj.transform.position = new Vector3(tower.x, tower.y, 0);
                }
                else
                {
                    // 创建新塔
                    var newTower = Instantiate(towerPrefab, towerContainer);
                    newTower.transform.position = new Vector3(tower.x, tower.y, 0);
                    _towerObjects[tower.id] = newTower;
                }
            }
        }

        private void UpdateEnemies(Enemy[] enemies)
        {
            if (enemies == null) return;

            var enemyIds = new HashSet<string>();
            foreach (var enemy in enemies)
            {
                enemyIds.Add(enemy.id);
            }

            foreach (var id in _enemyObjects.Keys)
            {
                if (!enemyIds.Contains(id))
                {
                    Destroy(_enemyObjects[id]);
                    _enemyObjects.Remove(id);
                }
            }

            foreach (var enemy in enemies)
            {
                if (_enemyObjects.TryGetValue(enemy.id, out var obj))
                {
                    obj.transform.position = new Vector3(enemy.x, enemy.y, 0);
                }
                else
                {
                    // 创建敌人（需要预制体）
                    Debug.Log($"New enemy: {enemy.type} at ({enemy.x}, {enemy.y})");
                }
            }
        }

        private void OnReceiveGift(ReceiveGiftResp gift)
        {
            Debug.Log($"收到礼物: {gift.gift_type} from {gift.sender_name}");
            // 显示礼物特效
        }

        private void OnReceiveDanmaku(ReceiveDanmakuResp danmaku)
        {
            Debug.Log($"弹幕: {danmaku.content}");
            // 显示弹幕
        }

        private void OnError(ErrorCode code, string message)
        {
            Debug.LogError($"Error {code}: {message}");
        }

        // UI 按钮回调
        public void OnClickPlaceArrow()
        {
            client.PlaceTower("arrow", 5f, 5f);
        }

        public void OnClickPlaceCannon()
        {
            client.PlaceTower("cannon", 7f, 5f);
        }

        public void OnClickSendGift()
        {
            client.SendGift("rocket");
        }

        public void OnClickSendDanmaku()
        {
            client.SendDanmaku("666");
        }
    }
}
