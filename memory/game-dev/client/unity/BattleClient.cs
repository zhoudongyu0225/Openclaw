using System;
using System.Collections;
using System.Collections.Generic;
using UnityEngine;
using UnityEngine.Networking;

namespace Danmaku.Battle
{
    /// <summary>
    /// 战斗类型
    /// </summary>
    public enum BattleType
    {
        PVE = 0,           // 人机对战
        PVP = 1,           // 玩家对战
        BOSS = 2,          // Boss挑战
        SURVIVAL = 3,      // 生存模式
        PRACTICE = 4,      // 练习模式
    }

    /// <summary>
    /// 战斗状态
    /// </summary>
    public enum BattleState
    {
        Waiting = 0,       // 等待开始
        Loading = 1,       // 加载中
        Playing = 2,       // 进行中
        Paused = 3,        // 暂停
        Finished = 4,      // 已结束
    }

    /// <summary>
    /// 难度等级
    /// </summary>
    public enum DifficultyLevel
    {
        Easy = 0,          // 简单
        Normal = 1,        // 普通
        Hard = 2,          // 困难
        Lunatic = 3,      // 疯狂
    }

    /// <summary>
    /// 玩家信息
    /// </summary>
    [Serializable]
    public class BattlePlayer
    {
        public string playerId;
        public string playerName;
        public int level;
        public int score;
        public int hp;
        public int maxHp;
        public int bomb;
        public int lives;
        public float power;
        public float graze;
        public long lastHitTime;
        public bool isDead;
        public Vector2 position;
    }

    /// <summary>
    /// 子弹信息
    /// </summary>
    [Serializable]
    public class BulletInfo
    {
        public string bulletId;
        public int type;
        public Vector2 position;
        public Vector2 velocity;
        public float rotation;
        public int color;
        public float damage;
        public bool isPlayer;
    }

    /// <summary>
    /// 敌人信息
    /// </summary>
    [Serializable]
    public class EnemyInfo
    {
        public string enemyId;
        public int enemyType;
        public Vector2 position;
        public int hp;
        public int maxHp;
        public float speed;
        public int score;
        public bool isBoss;
        public List<BulletInfo> bullets;
    }

    /// <summary>
    /// Boss信息
    /// </summary>
    [Serializable]
    public class BossInfo
    {
        public string bossId;
        public string enemyId;
        public string name;
        public int hp;
        public int maxHp;
        public int phase;
        public int maxPhase;
        public float spellCardTime;
        public float spellCardRemaining;
        public string spellCardName;
    }

    /// <summary>
    /// 道具信息
    /// </summary>
    [Serializable]
    public class ItemInfo
    {
        public string itemId;
        public int type;
        public Vector2 position;
        public float value;
    }

    /// <summary>
    /// 战斗结果
    /// </summary>
    [Serializable]
    public class BattleResult
    {
        public string playerId;
        public int score;
        public int rank;
        public int killCount;
        public int bossKillCount;
        public float graze;
        public float maxCombo;
        public float accuracy;
        public int time;
        public bool isWin;
    }

    /// <summary>
    /// 战斗帧数据
    /// </summary>
    [Serializable]
    public class BattleFrame
    {
        public int frameIndex;
        public List<PlayerAction> playerActions;
        public List<BulletInfo> newBullets;
        public List<string> removedBullets;
        public List<EnemyInfo> newEnemies;
        public List<string> removedEnemies;
        public List<ItemInfo> newItems;
    }

    /// <summary>
    /// 玩家操作
    /// </summary>
    [Serializable]
    public class PlayerAction
    {
        public string playerId;
        public Vector2 moveDirection;
        public bool shoot;
        public bool bomb;
        public bool focus;
        public bool pause;
    }

    /// <summary>
    /// 战斗初始化请求
    /// </summary>
    [Serializable]
    public class CSStartBattle
    {
        public string playerId;
        public BattleType battleType;
        public DifficultyLevel difficulty;
        public int stageId;
        public string roomId;
    }

    /// <summary>
    /// 战斗初始化响应
    /// </summary>
    [Serializable]
    public class SCStartBattle
    {
        public bool success;
        public string message;
        public string battleId;
        public BattleState state;
        public List<BattlePlayer> players;
        public int stageId;
        public int enemyLevel;
    }

    /// <summary>
    /// 玩家操作请求
    /// </summary>
    [Serializable]
    public class CSPlayerAction
    {
        public string battleId;
        public Vector2 moveDirection;
        public bool shoot;
        public bool bomb;
        public bool focus;
        public bool pause;
    }

    /// <summary>
    /// 战斗状态同步
    /// </summary>
    [Serializable]
    public class SCBattleSync
    {
        public int frameIndex;
        public BattleState state;
        public List<BattlePlayer> players;
        public List<BulletInfo> bullets;
        public List<EnemyInfo> enemies;
        public BossInfo boss;
        public List<ItemInfo> items;
    }

    /// <summary>
    /// 战斗结果响应
    /// </summary>
    [Serializable]
    public class SCBattleResult
    {
        public bool success;
        public string message;
        public BattleResult result;
        public List<BattleResult> rankings;
    }

    /// <summary>
    /// 使用道具请求
    /// </summary>
    [Serializable]
    public class CSUseItem
    {
        public string battleId;
        public int itemType;
    }

    /// <summary>
    /// 使用道具响应
    /// </summary>
    [Serializable]
    public class SCUseItem
    {
        public bool success;
        public string message;
        public int itemType;
        public int remainingCount;
    }

    /// <summary>
    /// 战斗客户端
    /// </summary>
    public class BattleClient : MonoBehaviour
    {
        [Header("服务器配置")]
        [SerializeField] private string serverUrl = "http://localhost:8080";
        [SerializeField] private float syncInterval = 0.033f; // 30fps

        [Header("战斗配置")]
        [SerializeField] private int maxPlayers = 4;
        [SerializeField] private bool enableInterpolation = true;
        [SerializeField] private float interpolationDelay = 0.1f;

        private string authToken;
        private string playerId;
        private string currentBattleId;
        private BattleState currentState = BattleState.Waiting;
        
        private List<BattlePlayer> players = new List<BattlePlayer>();
        private List<BulletInfo> bullets = new List<BulletInfo>();
        private List<EnemyInfo> enemies = new List<EnemyInfo>();
        private BossInfo boss;
        private List<ItemInfo> items = new List<ItemInfo>();
        
        private int localFrameIndex = 0;
        private int serverFrameIndex = 0;
        private List<BattleFrame> frameBuffer = new List<BattleFrame>();
        
        private bool isProcessing = false;
        private Coroutine syncCoroutine;

        public event Action<BattleState> OnStateChanged;
        public event Action<List<BattlePlayer>> OnPlayersUpdated;
        public event Action<List<BulletInfo>> OnBulletsUpdated;
        public event Action<List<EnemyInfo>> OnEnemiesUpdated;
        public event Action<BossInfo> OnBossUpdated;
        public event Action<List<ItemInfo>> OnItemsUpdated;
        public event Action<BattleResult> OnBattleFinished;

        public string BattleId => currentBattleId;
        public BattleState State => currentState;
        public List<BattlePlayer> Players => players;
        public List<BulletInfo> Bullets => bullets;
        public List<EnemyInfo> Enemies => enemies;
        public BossInfo Boss => boss;
        public List<ItemInfo> Items => items;

        /// <summary>
        /// 初始化战斗客户端
        /// </summary>
        public void Initialize(string token, string playerId)
        {
            this.authToken = token;
            this.playerId = playerId;
        }

        /// <summary>
        /// 开始战斗
        /// </summary>
        public void StartBattle(BattleType type, DifficultyLevel difficulty, int stageId, string roomId = null)
        {
            StartCoroutine(StartBattleCoroutine(type, difficulty, stageId, roomId));
        }

        private IEnumerator StartBattleCoroutine(BattleType type, DifficultyLevel difficulty, int stageId, string roomId)
        {
            var request = new CSStartBattle
            {
                playerId = playerId,
                battleType = type,
                difficulty = difficulty,
                stageId = stageId,
                roomId = roomId
            };

            using (UnityWebRequest www = UnityWebRequest.Post(
                $"{serverUrl}/api/battle/start",
                JsonUtility.ToJson(request)))
            {
                www.SetRequestHeader("Content-Type", "application/json");
                www.SetRequestHeader("Authorization", $"Bearer {authToken}");

                yield return www.SendWebRequest();

                if (www.result == UnityWebRequest.Result.Success)
                {
                    var response = JsonUtility.FromJson<SCStartBattle>(www.downloadHandler.text);
                    if (response.success)
                    {
                        currentBattleId = response.battleId;
                        currentState = response.state;
                        players = response.players;
                        
                        syncCoroutine = StartCoroutine(SyncCoroutine());
                        
                        OnStateChanged?.Invoke(currentState);
                    }
                }
                else
                {
                    Debug.LogError($"开始战斗失败: {www.error}");
                }
            }
        }

        /// <summary>
        /// 发送玩家操作
        /// </summary>
        public void SendPlayerAction(Vector2 moveDirection, bool shoot, bool bomb, bool focus, bool pause = false)
        {
            if (string.IsNullOrEmpty(currentBattleId) || currentState != BattleState.Playing)
                return;

            var action = new PlayerAction
            {
                playerId = playerId,
                moveDirection = moveDirection,
                shoot = shoot,
                bomb = bomb,
                focus = focus,
                pause = pause
            };

            StartCoroutine(SendActionCoroutine(action));
        }

        private IEnumerator SendActionCoroutine(PlayerAction action)
        {
            var request = new CSPlayerAction
            {
                battleId = currentBattleId,
                moveDirection = action.moveDirection,
                shoot = action.shoot,
                bomb = action.bomb,
                focus = action.focus,
                pause = action.pause
            };

            using (UnityWebRequest www = UnityWebRequest.Post(
                $"{serverUrl}/api/battle/action",
                JsonUtility.ToJson(request)))
            {
                www.SetRequestHeader("Content-Type", "application/json");
                www.SetRequestHeader("Authorization", $"Bearer {authToken}");
                
                yield return www.SendWebRequest();
            }
        }

        /// <summary>
        /// 同步协程
        /// </summary>
        private IEnumerator SyncCoroutine()
        {
            while (currentState == BattleState.Playing || currentState == BattleState.Loading)
            {
                yield return new WaitForSeconds(syncInterval);
                yield return SyncBattleState();
            }
        }

        /// <summary>
        /// 同步战斗状态
        /// </summary>
        private IEnumerator SyncBattleState()
        {
            if (isProcessing || string.IsNullOrEmpty(currentBattleId))
                yield break;

            isProcessing = true;

            using (UnityWebRequest www = UnityWebRequest.Get(
                $"{serverUrl}/api/battle/sync?battle_id={currentBattleId}&frame={localFrameIndex}"))
            {
                www.SetRequestHeader("Authorization", $"Bearer {authToken}");
                
                yield return www.SendWebRequest();

                if (www.result == UnityWebRequest.Result.Success)
                {
                    var sync = JsonUtility.FromJson<SCBattleSync>(www.downloadHandler.text);
                    
                    if (sync.state != currentState)
                    {
                        currentState = sync.state;
                        OnStateChanged?.Invoke(currentState);
                    }

                    if (sync.players != null)
                    {
                        players = sync.players;
                        OnPlayersUpdated?.Invoke(players);
                    }

                    if (sync.bullets != null)
                    {
                        bullets = sync.bullets;
                        OnBulletsUpdated?.Invoke(bullets);
                    }

                    if (sync.enemies != null)
                    {
                        enemies = sync.enemies;
                        OnEnemiesUpdated?.Invoke(enemies);
                    }

                    if (sync.boss != null)
                    {
                        boss = sync.boss;
                        OnBossUpdated?.Invoke(boss);
                    }

                    if (sync.items != null)
                    {
                        items = sync.items;
                        OnItemsUpdated?.Invoke(items);
                    }

                    serverFrameIndex = sync.frameIndex;
                }
            }

            isProcessing = false;
        }

        /// <summary>
        /// 使用道具
        /// </summary>
        public void UseItem(int itemType)
        {
            StartCoroutine(UseItemCoroutine(itemType));
        }

        private IEnumerator UseItemCoroutine(int itemType)
        {
            var request = new CSUseItem
            {
                battleId = currentBattleId,
                itemType = itemType
            };

            using (UnityWebRequest www = UnityWebRequest.Post(
                $"{serverUrl}/api/battle/use_item",
                JsonUtility.ToJson(request)))
            {
                www.SetRequestHeader("Content-Type", "application/json");
                www.SetRequestHeader("Authorization", $"Bearer {authToken}");
                
                yield return www.SendWebRequest();

                if (www.result == UnityWebRequest.Result.Success)
                {
                    var response = JsonUtility.FromJson<SCUseItem>(www.downloadHandler.text);
                    Debug.Log($"使用道具: {response.message}");
                }
            }
        }

        /// <summary>
        /// 离开战斗
        /// </summary>
        public void LeaveBattle()
        {
            if (!string.IsNullOrEmpty(currentBattleId))
            {
                StartCoroutine(LeaveBattleCoroutine());
            }

            if (syncCoroutine != null)
            {
                StopCoroutine(syncCoroutine);
            }

            currentBattleId = null;
            currentState = BattleState.Waiting;
            players.Clear();
            bullets.Clear();
            enemies.Clear();
            boss = null;
            items.Clear();
        }

        private IEnumerator LeaveBattleCoroutine()
        {
            using (UnityWebRequest www = UnityWebRequest.Post(
                $"{serverUrl}/api/battle/leave",
                JsonUtility.ToJson(new { battle_id = currentBattleId, player_id = playerId })))
            {
                www.SetRequestHeader("Content-Type", "application/json");
                www.SetRequestHeader("Authorization", $"Bearer {authToken}");
                
                yield return www.SendWebRequest();
            }
        }

        /// <summary>
        /// 获取战斗结果
        /// </summary>
        public void GetBattleResult()
        {
            StartCoroutine(GetBattleResultCoroutine());
        }

        private IEnumerator GetBattleResultCoroutine()
        {
            using (UnityWebRequest www = UnityWebRequest.Get(
                $"{serverUrl}/api/battle/result?battle_id={currentBattleId}"))
            {
                www.SetRequestHeader("Authorization", $"Bearer {authToken}");
                
                yield return www.SendWebRequest();

                if (www.result == UnityWebRequest.Result.Success)
                {
                    var response = JsonUtility.FromJson<SCBattleResult>(www.downloadHandler.text);
                    if (response.success && response.result != null)
                    {
                        OnBattleFinished?.Invoke(response.result);
                    }
                }
            }
        }

        /// <summary>
        /// 本地插入子弹（客户端预测）
        /// </summary>
        public void AddLocalBullet(BulletInfo bullet)
        {
            bullets.Add(bullet);
        }

        /// <summary>
        /// 本地移除子弹
        /// </summary>
        public void RemoveLocalBullet(string bulletId)
        {
            bullets.RemoveAll(b => b.bulletId == bulletId);
        }

        /// <summary>
        /// 获取本地玩家
        /// </summary>
        public BattlePlayer GetLocalPlayer()
        {
            return players.Find(p => p.playerId == playerId);
        }

        /// <summary>
        /// 获取玩家By ID
        /// </summary>
        public BattlePlayer GetPlayer(string playerId)
        {
            return players.Find(p => p.playerId == playerId);
        }

        /// <summary>
        /// 检查是否存活
        /// </summary>
        public bool IsAlive()
        {
            var player = GetLocalPlayer();
            return player != null && !player.isDead && player.hp > 0;
        }

        /// <summary>
        /// 获取剩余时间
        /// </summary>
        public float GetRemainingTime()
        {
            // 假设战斗时间限制为 180 秒
            float elapsed = localFrameIndex / 30f;
            return Mathf.Max(0, 180f - elapsed);
        }

        private void OnDestroy()
        {
            LeaveBattle();
        }
    }

    /// <summary>
    /// 战斗UI管理器
    /// </summary>
    public class BattleUIManager : MonoBehaviour
    {
        [Header("UI组件")]
        [SerializeField] private GameObject battleHUD;
        [SerializeField] private GameObject pauseMenu;
        [SerializeField] private GameObject resultPanel;
        
        [Header("状态显示")]
        [SerializeField] private Text scoreText;
        [SerializeField] private Text timeText;
        [SerializeField] private Text hpText;
        [SerializeField] private Text bombText;
        [SerializeField] private Text livesText;
        [SerializeField] private Text powerText;
        [SerializeField] private Text grazeText;
        
        [Header]="[SerializeField] private Slider hpSlider;
        [SerializeField] private Slider powerSlider;

        private BattleClient battleClient;

        public void Initialize(BattleClient client)
        {
            battleClient = client;
            
            battleClient.OnStateChanged += OnBattleStateChanged;
            battleClient.OnPlayersUpdated += OnPlayersUpdated;
            battleClient.OnBattleFinished += OnBattleFinished;
        }

        private void OnBattleStateChanged(BattleState state)
        {
            switch (state)
            {
                case BattleState.Loading:
                    battleHUD.SetActive(true);
                    pauseMenu.SetActive(false);
                    resultPanel.SetActive(false);
                    break;
                    
                case BattleState.Playing:
                    battleHUD.SetActive(true);
                    pauseMenu.SetActive(false);
                    resultPanel.SetActive(false);
                    break;
                    
                case BattleState.Paused:
                    battleHUD.SetActive(true);
                    pauseMenu.SetActive(true);
                    resultPanel.SetActive(false);
                    break;
                    
                case BattleState.Finished:
                    battleHUD.SetActive(false);
                    pauseMenu.SetActive(false);
                    resultPanel.SetActive(true);
                    break;
            }
        }

        private void OnPlayersUpdated(List<BattlePlayer> players)
        {
            var localPlayer = players.Find(p => p.playerId == battleClient.GetLocalPlayer()?.playerId);
            if (localPlayer == null) return;

            // 更新UI
            scoreText.text = localPlayer.score.ToString("N0");
            hpText.text = $"{localPlayer.hp}/{localPlayer.maxHp}";
            bombText.text = localPlayer.bomb.ToString();
            livesText.text = localPlayer.lives.ToString();
            powerText.text = localPlayer.power.ToString("F1");
            grazeText.text = localPlayer.graze.ToString("F0");

            // 更新血条
            hpSlider.maxValue = localPlayer.maxHp;
            hpSlider.value = localPlayer.hp;

            // 更新Power条
            powerSlider.maxValue = 100;
            powerSlider.value = localPlayer.power;
        }

        private void OnBattleFinished(BattleResult result)
        {
            resultPanel.SetActive(true);
            
            // 显示结果
            Debug.Log($"战斗结束: 分数={result.score}, 排名={result.rank}");
        }

        /// <summary>
        /// 更新倒计时
        /// </summary>
        private void Update()
        {
            if (battleClient != null && battleClient.State == BattleState.Playing)
            {
                timeText.text = battleClient.GetRemainingTime().ToString("F1");
            }
        }

        private void OnDestroy()
        {
            if (battleClient != null)
            {
                battleClient.OnStateChanged -= OnBattleStateChanged;
                battleClient.OnPlayersUpdated -= OnPlayersUpdated;
                battleClient.OnBattleFinished -= OnBattleFinished;
            }
        }
    }

    /// <summary>
    /// 子弹渲染器
    /// </summary>
    public class BulletRenderer : MonoBehaviour
    {
        [Header("子弹预制体")]
        [SerializeField] private GameObject[] bulletPrefabs;
        [SerializeField] private float bulletScale = 1f;

        private Dictionary<string, GameObject> activeBullets = new Dictionary<string, GameObject>();
        private BattleClient battleClient;

        public void Initialize(BattleClient client)
        {
            battleClient = client;
            battleClient.OnBulletsUpdated += OnBulletsUpdated;
        }

        private void OnBulletsUpdated(List<BulletInfo> bullets)
        {
            var currentIds = new HashSet<string>();
            
            foreach (var bullet in bullets)
            {
                currentIds.Add(bullet.bulletId);
                
                if (!activeBullets.ContainsKey(bullet.bulletId))
                {
                    SpawnBullet(bullet);
                }
                else
                {
                    UpdateBullet(bullet);
                }
            }

            // 移除消失的子弹
            var toRemove = new List<string>();
            foreach (var id in activeBullets.Keys)
            {
                if (!currentIds.Contains(id))
                {
                    Destroy(activeBullets[id]);
                    toRemove.Add(id);
                }
            }
            
            foreach (var id in toRemove)
            {
                activeBullets.Remove(id);
            }
        }

        private void SpawnBullet(BulletInfo info)
        {
            if (info.type < 0 || info.type >= bulletPrefabs.Length)
                return;

            var bullet = Instantiate(bulletPrefabs[info.type], transform);
            bullet.transform.position = new Vector3(info.position.x, info.position.y, 0);
            bullet.transform.localScale = Vector3.one * bulletScale;
            
            var renderer = bullet.GetComponent<BulletBehavior>();
            if (renderer != null)
            {
                renderer.Initialize(info);
            }

            activeBullets[info.bulletId] = bullet;
        }

        private void UpdateBullet(BulletInfo info)
        {
            if (!activeBullets.ContainsKey(info.bulletId))
                return;

            var bullet = activeBullets[info.bulletId];
            bullet.transform.position = new Vector3(info.position.x, info.position.y, 0);
            bullet.transform.rotation = Quaternion.Euler(0, 0, info.rotation * Mathf.Rad2Deg);
        }

        private void OnDestroy()
        {
            if (battleClient != null)
            {
                battleClient.OnBulletsUpdated -= OnBulletsUpdated;
            }
        }
    }

    /// <summary>
    /// 子弹行为
    /// </summary>
    public class BulletBehavior : MonoBehaviour
    {
        private BulletInfo info;
        private bool isLocalUpdate = false;

        public void Initialize(BulletInfo bulletInfo)
        {
            info = bulletInfo;
        }

        public void SetLocalUpdate(bool local)
        {
            isLocalUpdate = local;
        }

        private void Update()
        {
            if (!isLocalUpdate)
            {
                // 服务器同步模式，不需要本地移动
                return;
            }

            // 本地预测移动
            transform.position += new Vector3(
                info.velocity.x * Time.deltaTime,
                info.velocity.y * Time.deltaTime,
                0
            );

            // 超出屏幕移除
            var pos = transform.position;
            if (pos.x < -20 || pos.x > 20 || pos.y < -15 || pos.y > 15)
            {
                Destroy(gameObject);
            }
        }

        /// <summary>
        /// 碰撞检测
        /// </summary>
        private void OnTriggerEnter2D(Collider2D other)
        {
            if (info.isPlayer)
            {
                // 玩家子弹击中敌人
                var enemy = other.GetComponent<EnemyBehavior>();
                if (enemy != null)
                {
                    enemy.TakeDamage((int)info.damage);
                    Destroy(gameObject);
                }
            }
            else
            {
                // 敌人子弹击中玩家
                var player = other.GetComponent<PlayerBehavior>();
                if (player != null)
                {
                    player.TakeDamage((int)info.damage);
                    Destroy(gameObject);
                }
            }
        }
    }

    /// <summary>
    /// 玩家控制器
    /// </summary>
    public class PlayerController : MonoBehaviour
    {
        [Header("移动设置")]
        [SerializeField] private float moveSpeed = 5f;
        [SerializeField] private float slowSpeed = 2f;
        
        [Header="[SerializeField] private KeyCode shootKey = KeyCode.Z;
        [SerializeField] private KeyCode bombKey = KeyCode.X;
        [SerializeField] private KeyCode focusKey = KeyCode.ShiftLeft;
        [SerializeField] private KeyCode pauseKey = KeyCode.Escape;

        [Header("碰撞")]
        [SerializeField] private Collider2D hitbox;

        private BattleClient battleClient;
        private Vector2 moveDirection;
        private bool isShooting;
        private bool isBombing;
        private bool isFocusing;

        public void Initialize(BattleClient client)
        {
            battleClient = client;
        }

        private void Update()
        {
            if (battleClient == null || battleClient.State != BattleState.Playing)
                return;

            // 读取输入
            moveDirection = Vector2.zero;
            if (Input.GetKey(KeyCode.UpArrow) || Input.GetKey(KeyCode.W))
                moveDirection.y = 1;
            if (Input.GetKey(KeyCode.DownArrow) || Input.GetKey(KeyCode.S))
                moveDirection.y = -1;
            if (Input.GetKey(KeyCode.LeftArrow) || Input.GetKey(KeyCode.A))
                moveDirection.x = -1;
            if (Input.GetKey(KeyCode.RightArrow) || Input.GetKey(KeyCode.D))
                moveDirection.x = 1;

            isShooting = Input.GetKey(shootKey);
            isBombing = Input.GetKeyDown(bombKey);
            isFocusing = Input.GetKey(focusKey);

            // 移动
            float speed = isFocusing ? slowSpeed : moveSpeed;
            transform.Translate(moveDirection * speed * Time.deltaTime);

            // 限制边界
            var pos = transform.position;
            pos.x = Mathf.Clamp(pos.x, -8f, 8f);
            pos.y = Mathf.Clamp(pos.y, -5f, 5f);
            transform.position = pos;

            // 聚焦模式显示判定点
            if (hitbox != null)
            {
                hitbox.enabled = isFocusing;
            }

            // 发送操作到服务器
            battleClient.SendPlayerAction(
                moveDirection,
                isShooting,
                isBombing,
                isFocusing,
                Input.GetKeyDown(pauseKey)
            );
        }
    }

    /// <summary>
    /// 敌人行为
    /// </summary>
    public class EnemyBehavior : MonoBehaviour
    {
        private EnemyInfo info;
        private int currentHp;
        private bool isDead = false;

        public void Initialize(EnemyInfo enemyInfo)
        {
            info = enemyInfo;
            currentHp = enemyInfo.hp;
            transform.position = new Vector3(info.position.x, info.position.y, 0);
        }

        public void TakeDamage(int damage)
        {
            if (isDead) return;

            currentHp -= damage;
            
            if (currentHp <= 0)
            {
                Die();
            }
        }

        private void Die()
        {
            isDead = true;
            // 播放死亡特效
            Destroy(gameObject);
        }

        private void Update()
        {
            if (isDead) return;

            // 简单的移动逻辑，实际应该从服务器同步
            if (!info.isBoss)
            {
                transform.Translate(Vector3.down * info.speed * Time.deltaTime);
            }
        }

        /// <summary>
        /// Boss行为（由服务器控制）
        /// </summary>
        public void SetBossPattern(int patternIndex)
        {
            // Boss技能模式切换
            Debug.Log($"Boss切换模式: {patternIndex}");
        }
    }
}
