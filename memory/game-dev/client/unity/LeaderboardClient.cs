// ============================================================
// 排行榜客户端 - Leaderboard Client
// 弹幕游戏 Unity 客户端
// ============================================================

using System;
using System.Collections.Generic;
using System.Linq;
using UnityEngine;

namespace DanmakuGameClient
{
    // ============================================================
    // 排行榜类型
    // ============================================================
    
    public enum LeaderboardType
    {
        Level = 1,           // 等级榜
        Gold = 2,             // 金币榜
        Gem = 3,              // 钻石榜
        Combat = 4,          // 战斗力榜
        Kill = 5,             // 击杀榜
        Damage = 6,           // 伤害榜
        Survival = 7,         // 生存榜
        Win = 8,              // 胜场榜
        Guild = 9,            // 公会榜
        SignIn = 10,         // 签到榜
        Rich = 11,            // 富豪榜
        MVP = 12,             // MVP榜
    }

    // ============================================================
    // 排行榜时间范围
    // ============================================================
    
    public enum LeaderboardPeriod
    {
        All = 1,      // 全部
        Daily = 2,    // 今日
        Weekly = 3,   // 本周
        Monthly = 4,  // 本月
        Season = 5,   // 本赛季
    }

    // ============================================================
    // 排行榜条目
    // ============================================================
    
    [Serializable]
    public class LeaderboardEntry
    {
        public int rank;              // 排名
        public string playerId;        // 玩家ID
        public string playerName;      // 玩家名称
        public int level;              // 等级
        public long value;             // 排行榜数值
        public string avatar;          // 头像
        public int titleId;            // 称号ID
        public string guildName;       // 公会名称
        public int vipLevel;           // VIP等级
        public bool isMyFriend;        // 是否好友
        public long change;            // 排名变化 (+上升 -下降 =不变)
    }

    // ============================================================
    // 排行榜数据
    // ============================================================
    
    [Serializable]
    public class LeaderboardData
    {
        public LeaderboardType type;           // 排行榜类型
        public LeaderboardPeriod period;       // 时间范围
        public List<LeaderboardEntry> entries; // 排行榜条目
        public int myRank;                     // 我的排名
        public long myValue;                   // 我的数值
        public int totalCount;                 // 总人数
        public long refreshTime;               // 刷新时间
    }

    // ============================================================
    // 我的排行榜信息
    // ============================================================
    
    [Serializable]
    public class MyLeaderboardInfo
    {
        public List<LeaderboardType> subscribedTypes;   // 订阅的排行榜类型
        public Dictionary<LeaderboardType, int> ranks;   // type -> rank
        public Dictionary<LeaderboardType, long> values; // type -> value
    }

    // ============================================================
    // 排行榜客户端主类
    // ============================================================
    
    public class LeaderboardClient : MonoBehaviour
    {
        private static LeaderboardClient _instance;
        public static LeaderboardClient Instance
        {
            get
            {
                if (_instance == null)
                {
                    _instance = FindObjectOfType<LeaderboardClient>();
                }
                return _instance;
            }
        }

        // 缓存
        private Dictionary<LeaderboardType, LeaderboardData> _cache = new Dictionary<LeaderboardType, LeaderboardData>();
        private Dictionary<LeaderboardType, float> _lastRefreshTime = new Dictionary<LeaderboardType, float>();
        
        // 配置
        private const float CACHE_DURATION = 60f; // 缓存60秒
        private const int PAGE_SIZE = 50;          // 每页50条
        
        // 我的信息
        private MyLeaderboardInfo _myInfo = new MyLeaderboardInfo();

        // 事件
        public event Action<LeaderboardType, LeaderboardData> OnLeaderboardReceived;
        public event Action<LeaderboardType, int> OnMyRankReceived;
        public event Action<string> OnError;

        // ============================================================
        // 初始化
        // ============================================================
        
        private void Awake()
        {
            _instance = this;
            _myInfo.subscribedTypes = new List<LeaderboardType>();
            _myInfo.ranks = new Dictionary<LeaderboardType, int>();
            _myInfo.values = new Dictionary<LeaderboardType, long>();
        }

        // ============================================================
        // 获取排行榜 (缓存版本)
        // ============================================================
        
        public LeaderboardData GetLeaderboard(LeaderboardType type, LeaderboardPeriod period = LeaderboardPeriod.All)
        {
            if (_cache.TryGetValue(type, out var data))
            {
                if (data.period == period)
                {
                    float lastTime;
                    if (_lastRefreshTime.TryGetValue(type, out lastTime))
                    {
                        if (Time.time - lastTime < CACHE_DURATION)
                        {
                            return data;
                        }
                    }
                }
            }
            return null;
        }

        // ============================================================
        // 请求排行榜
        // ============================================================
        
        public void RequestLeaderboard(LeaderboardType type, LeaderboardPeriod period = LeaderboardPeriod.All, int page = 1)
        {
            var req = new CSGetLeaderboard
            {
                Type = (int)type,
                Period = (int)period,
                Page = page,
                PageSize = PAGE_SIZE
            };
            
            NetworkManager.Instance.Send(req, (sc) =>
            {
                var response = sc as SCLeaderboardData;
                if (response != null)
                {
                    ProcessLeaderboardResponse(type, period, response);
                }
            });
        }

        // ============================================================
        // 请求我的排名
        // ============================================================
        
        public void RequestMyRank(LeaderboardType type)
        {
            var req = new CSGetMyLeaderboardRank
            {
                Type = (int)type
            };
            
            NetworkManager.Instance.Send(req, (sc) =>
            {
                var response = sc as SCMyLeaderboardRank;
                if (response != null)
                {
                    UpdateMyRank(type, response.Rank, response.Value);
                    OnMyRankReceived?.Invoke(type, response.Rank);
                }
            });
        }

        // ============================================================
        // 请求多个排行榜
        // ============================================================
        
        public void RequestMultipleLeaderboards(List<LeaderboardType> types, LeaderboardPeriod period = LeaderboardPeriod.All)
        {
            foreach (var type in types)
            {
                RequestLeaderboard(type, period);
            }
        }

        // ============================================================
        // 订阅排行榜 (实时更新)
        // ============================================================
        
        public void SubscribeLeaderboard(LeaderboardType type)
        {
            if (!_myInfo.subscribedTypes.Contains(type))
            {
                _myInfo.subscribedTypes.Add(type);
            }
            
            var req = new CSSubscribeLeaderboard
            {
                Type = (int)type,
                Subscribe = true
            };
            
            NetworkManager.Instance.Send(req);
        }

        // ============================================================
        // 取消订阅排行榜
        // ============================================================
        
        public void UnsubscribeLeaderboard(LeaderboardType type)
        {
            _myInfo.subscribedTypes.Remove(type);
            
            var req = new CSSubscribeLeaderboard
            {
                Type = (int)type,
                Subscribe = false
            };
            
            NetworkManager.Instance.Send(req);
        }

        // ============================================================
        // 请求公会排行榜
        // ============================================================
        
        public void RequestGuildLeaderboard(LeaderboardPeriod period = LeaderboardPeriod.All, int page = 1)
        {
            RequestLeaderboard(LeaderboardType.Guild, period, page);
        }

        // ============================================================
        // 获取排行榜前三名
        // ============================================================
        
        public List<LeaderboardEntry> GetTopThree(LeaderboardType type, LeaderboardPeriod period = LeaderboardPeriod.All)
        {
            var data = GetLeaderboard(type, period);
            if (data != null && data.entries != null)
            {
                return data.entries.Take(3).ToList();
            }
            return new List<LeaderboardEntry>();
        }

        // ============================================================
        // 获取我的排名信息
        // ============================================================
        
        public int GetMyRank(LeaderboardType type)
        {
            if (_myInfo.ranks.TryGetValue(type, out var rank))
            {
                return rank;
            }
            return -1;
        }

        // ============================================================
        // 获取我的数值
        // ============================================================
        
        public long GetMyValue(LeaderboardType type)
        {
            if (_myInfo.values.TryGetValue(type, out var value))
            {
                return value;
            }
            return 0;
        }

        // ============================================================
        // 检查是否上榜
        // ============================================================
        
        public bool IsOnBoard(LeaderboardType type)
        {
            return GetMyRank(type) > 0;
        }

        // ============================================================
        // 刷新排行榜
        // ============================================================
        
        public void RefreshLeaderboard(LeaderboardType type, LeaderboardPeriod period = LeaderboardPeriod.All)
        {
            _cache.Remove(type);
            _lastRefreshTime.Remove(type);
            RequestLeaderboard(type, period);
        }

        // ============================================================
        // 清空缓存
        // ============================================================
        
        public void ClearCache()
        {
            _cache.Clear();
            _lastRefreshTime.Clear();
        }

        // ============================================================
        // 处理排行榜响应
        // ============================================================
        
        private void ProcessLeaderboardResponse(LeaderboardType type, LeaderboardPeriod period, SCLeaderboardData response)
        {
            var data = new LeaderboardData
            {
                type = type,
                period = period,
                entries = new List<LeaderboardEntry>(),
                myRank = response.MyRank,
                myValue = response.MyValue,
                totalCount = response.TotalCount,
                refreshTime = response.RefreshTime
            };
            
            if (response.Entries != null)
            {
                foreach (var entry in response.Entries)
                {
                    data.entries.Add(new LeaderboardEntry
                    {
                        rank = entry.Rank,
                        playerId = entry.PlayerId,
                        playerName = entry.PlayerName,
                        level = entry.Level,
                        value = entry.Value,
                        avatar = entry.Avatar,
                        titleId = entry.TitleId,
                        guildName = entry.GuildName,
                        vipLevel = entry.VipLevel,
                        isMyFriend = entry.IsMyFriend,
                        change = entry.Change
                    });
                }
            }
            
            _cache[type] = data;
            _lastRefreshTime[type] = Time.time;
            
            OnLeaderboardReceived?.Invoke(type, data);
        }

        // ============================================================
        // 更新我的排名
        // ============================================================
        
        private void UpdateMyRank(LeaderboardType type, int rank, long value)
        {
            _myInfo.ranks[type] = rank;
            _myInfo.values[type] = value;
        }

        // ============================================================
        // 处理推送更新
        // ============================================================
        
        public void HandleLeaderboardUpdate(SCLeaderboardUpdate update)
        {
            var type = (LeaderboardType)update.Type;
            
            if (_cache.TryGetValue(type, out var data))
            {
                // 更新对应条目
                var entry = data.entries.FirstOrDefault(e => e.playerId == update.PlayerId);
                if (entry != null)
                {
                    entry.rank = update.Rank;
                    entry.value = update.Value;
                    entry.change = update.Change;
                }
                
                // 重新排序
                data.entries.Sort((a, b) => b.value.CompareTo(a.value));
                for (int i = 0; i < data.entries.Count; i++)
                {
                    data.entries[i].rank = i + 1;
                }
                
                OnLeaderboardReceived?.Invoke(type, data);
            }
        }
    }

    // ============================================================
    // 网络协议 (需要与服务器对应)
    // ============================================================
    
    // 请求排行榜
    public class CSGetLeaderboard
    {
        public int Type { get; set; }
        public int Period { get; set; }
        public int Page { get; set; }
        public int PageSize { get; set; }
    }

    // 排行榜数据响应
    public class SCLeaderboardData
    {
        public int Type { get; set; }
        public int Period { get; set; }
        public List<LeaderboardEntryProto> Entries { get; set; }
        public int MyRank { get; set; }
        public long MyValue { get; set; }
        public int TotalCount { get; set; }
        public long RefreshTime { get; set; }
    }

    // 排行榜条目协议
    public class LeaderboardEntryProto
    {
        public int Rank { get; set; }
        public string PlayerId { get; set; }
        public string PlayerName { get; set; }
        public int Level { get; set; }
        public long Value { get; set; }
        public string Avatar { get; set; }
        public int TitleId { get; set; }
        public string GuildName { get; set; }
        public int VipLevel { get; set; }
        public bool IsMyFriend { get; set; }
        public long Change { get; set; }
    }

    // 请求我的排名
    public class CSGetMyLeaderboardRank
    {
        public int Type { get; set; }
    }

    // 我的排名响应
    public class SCMyLeaderboardRank
    {
        public int Type { get; set; }
        public int Rank { get; set; }
        public long Value { get; set; }
    }

    // 订阅/取消订阅排行榜
    public class CSSubscribeLeaderboard
    {
        public int Type { get; set; }
        public bool Subscribe { get; set; }
    }

    // 排行榜推送更新
    public class SCLeaderboardUpdate
    {
        public int Type { get; set; }
        public string PlayerId { get; set; }
        public int Rank { get; set; }
        public long Value { get; set; }
        public long Change { get; set; }
    }
}

// ============================================================
// 排行榜 UI 管理器
// ============================================================

namespace DanmakuGameClient
{
    public class LeaderboardUIManager : MonoBehaviour
    {
        [Header("排行榜类型按钮")]
        public GameObject[] typeButtons;
        
        [Header("排行榜条目预制体")]
        public GameObject entryPrefab;
        
        [Header("排行榜容器")]
        public Transform entriesContainer;
        
        [Header("我的排名显示")]
        public Text myRankText;
        public Text myValueText;
        
        [Header("顶部前三名")]
        public Transform topThreeContainer;
        
        [Header("Loading")]
        public GameObject loadingObj;

        private LeaderboardType _currentType = LeaderboardType.Level;
        private LeaderboardPeriod _currentPeriod = LeaderboardPeriod.All;
        private List<LeaderboardEntry> _currentEntries = new List<LeaderboardEntry>();

        // ============================================================
        // 初始化
        // ============================================================
        
        private void Start()
        {
            // 绑定事件
            LeaderboardClient.Instance.OnLeaderboardReceived += OnLeaderboardReceived;
            LeaderboardClient.Instance.OnError += OnError;
            
            // 默认请求等级榜
            RequestCurrentLeaderboard();
        }

        private void OnDestroy()
        {
            LeaderboardClient.Instance.OnLeaderboardReceived -= OnLeaderboardReceived;
            LeaderboardClient.Instance.OnError -= OnError;
        }

        // ============================================================
        // 选择排行榜类型
        // ============================================================
        
        public void SelectType(int typeIndex)
        {
            _currentType = (LeaderboardType)typeIndex;
            RequestCurrentLeaderboard();
        }

        // ============================================================
        // 选择时间范围
        // ============================================================
        
        public void SelectPeriod(int periodIndex)
        {
            _currentPeriod = (LeaderboardPeriod)periodIndex;
            RequestCurrentLeaderboard();
        }

        // ============================================================
        // 请求当前排行榜
        // ============================================================
        
        private void RequestCurrentLeaderboard()
        {
            ShowLoading(true);
            LeaderboardClient.Instance.RequestLeaderboard(_currentType, _currentPeriod);
            LeaderboardClient.Instance.RequestMyRank(_currentType);
        }

        // ============================================================
        // 刷新排行榜
        // ============================================================
        
        public void Refresh()
        {
            LeaderboardClient.Instance.RefreshLeaderboard(_currentType, _currentPeriod);
        }

        // ============================================================
        // 处理排行榜数据
        // ============================================================
        
        private void OnLeaderboardReceived(LeaderboardType type, LeaderboardData data)
        {
            if (type != _currentType) return;
            
            ShowLoading(false);
            _currentEntries = data.entries;
            
            // 更新UI
            UpdateEntriesUI(data.entries);
            UpdateMyRankUI(data.myRank, data.myValue);
            UpdateTopThreeUI(data.entries);
        }

        // ============================================================
        // 更新排行榜条目UI
        // ============================================================
        
        private void UpdateEntriesUI(List<LeaderboardEntry> entries)
        {
            // 清空容器
            foreach (Transform child in entriesContainer)
            {
                Destroy(child.gameObject);
            }
            
            // 生成条目
            foreach (var entry in entries)
            {
                var go = Instantiate(entryPrefab, entriesContainer);
                var item = go.GetComponent<LeaderboardEntryItem>();
                item.Initialize(entry);
            }
        }

        // ============================================================
        // 更新我的排名UI
        // ============================================================
        
        private void UpdateMyRankUI(int rank, long value)
        {
            if (myRankText != null)
            {
                myRankText.text = rank > 0 ? $"#{rank}" : "未上榜";
            }
            
            if (myValueText != null)
            {
                myValueText.text = value.ToString("N0");
            }
        }

        // ============================================================
        // 更新前三名UI
        // ============================================================
        
        private void UpdateTopThreeUI(List<LeaderboardEntry> entries)
        {
            if (topThreeContainer == null || entries.Count < 3) return;
            
            var topThree = entries.Take(3).ToList();
            var items = topThreeContainer.GetComponentsInChildren<LeaderboardEntryItem>();
            
            for (int i = 0; i < Math.Min(3, items.Length); i++)
            {
                items[i].Initialize(topThree[i]);
            }
        }

        // ============================================================
        // 显示Loading
        // ============================================================
        
        private void ShowLoading(bool show)
        {
            if (loadingObj != null)
            {
                loadingObj.SetActive(show);
            }
        }

        // ============================================================
        // 错误处理
        // ============================================================
        
        private void OnError(string error)
        {
            Debug.LogError($"Leaderboard Error: {error}");
            ShowLoading(false);
        }
    }

    // ============================================================
    // 排行榜条目 Item
    // ============================================================
    
    public class LeaderboardEntryItem : MonoBehaviour
    {
        [Header("排名")]
        public Text rankText;
        
        [Header("头像")]
        public Image avatarImage;
        
        [Header("玩家名")]
        public Text nameText;
        
        [Header("等级")]
        public Text levelText;
        
        [Header("数值")]
        public Text valueText;
        
        [Header("公会名")]
        public Text guildText;
        
        [Header("VIP标识")]
        public GameObject vipObj;
        
        [Header("好友标识")]
        public GameObject friendObj;
        
        [Header("排名变化")]
        public Text changeText;
        
        [Header("前三名特效")]
        public GameObject[] topThreeEffects;

        // ============================================================
        // 初始化
        // ============================================================
        
        public void Initialize(LeaderboardEntry entry)
        {
            // 排名
            if (rankText != null)
            {
                rankText.text = entry.rank.ToString();
                rankText.color = GetRankColor(entry.rank);
            }
            
            // 玩家名
            if (nameText != null)
            {
                nameText.text = entry.playerName;
            }
            
            // 等级
            if (levelText != null)
            {
                levelText.text = $"Lv.{entry.level}";
            }
            
            // 数值
            if (valueText != null)
            {
                valueText.text = entry.value.ToString("N0");
            }
            
            // 公会
            if (guildText != null)
            {
                guildText.text = entry.guildName;
                guildText.gameObject.SetActive(!string.IsNullOrEmpty(entry.guildName));
            }
            
            // VIP
            if (vipObj != null)
            {
                vipObj.SetActive(entry.vipLevel > 0);
            }
            
            // 好友
            if (friendObj != null)
            {
                friendObj.SetActive(entry.isMyFriend);
            }
            
            // 排名变化
            if (changeText != null)
            {
                if (entry.change > 0)
                {
                    changeText.text = $"↑{entry.change}";
                    changeText.color = Color.green;
                }
                else if (entry.change < 0)
                {
                    changeText.text = $"↓{Math.Abs(entry.change)}";
                    changeText.color = Color.red;
                }
                else
                {
                    changeText.text = "=";
                    changeText.color = Color.gray;
                }
            }
            
            // 前三名特效
            UpdateTopThreeEffects(entry.rank);
            
            // 加载头像
            LoadAvatar(entry.avatar);
        }

        // ============================================================
        // 获取排名颜色
        // ============================================================
        
        private Color GetRankColor(int rank)
        {
            switch (rank)
            {
                case 1: return new Color(1f, 0.84f, 0f); // 金色
                case 2: return new Color(0.8f, 0.8f, 0.8f); // 银色
                case 3: return new Color(0.8f, 0.5f, 0.2f); // 铜色
                default: return Color.white;
            }
        }

        // ============================================================
        // 更新前三名特效
        // ============================================================
        
        private void UpdateTopThreeEffects(int rank)
        {
            for (int i = 0; i < topThreeEffects.Length; i++)
            {
                topThreeEffects[i].SetActive(rank == i + 1);
            }
        }

        // ============================================================
        // 加载头像
        // ============================================================
        
        private void LoadAvatar(string avatarUrl)
        {
            if (avatarImage == null || string.IsNullOrEmpty(avatarUrl)) return;
            
            // TODO: 实现头像加载
            // StartCoroutine(LoadImage(avatarUrl, avatarImage));
        }

        // ============================================================
        // 点击条目
        // ============================================================
        
        public void OnEntryClick()
        {
            // TODO: 查看玩家详情
        }
    }
}
