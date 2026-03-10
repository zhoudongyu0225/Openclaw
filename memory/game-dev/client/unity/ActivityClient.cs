// ============================================================
// 活动客户端 - Activity Client
// 弹幕游戏 Unity 客户端
// ============================================================

using System;
using System.Collections.Generic;
using System.Linq;
using UnityEngine;

namespace DanmakuGameClient
{
    // ============================================================
    // 活动类型
    // ============================================================
    
    public enum ActivityType
    {
        Daily = 1,          // 每日活动
        Weekly = 2,         // 每周活动
        Season = 3,         // 赛季活动
        Festival = 4,       // 节日活动
        LimitTime = 5,      // 限时活动
        Login = 6,          // 登录活动
        Recharge = 7,       // 充值活动
        Consumption = 8,   // 消费活动
        Battle = 9,         // 战斗活动
        Guild = 10,         // 公会活动
    }

    // ============================================================
    // 活动状态
    // ============================================================
    
    public enum ActivityStatus
    {
        Upcoming = 1,   // 即将开始
        Active = 2,     // 进行中
        Ended = 3,      // 已结束
    }

    // ============================================================
    // 活动奖励状态
    // ============================================================
    
    public enum RewardStatus
    {
        Locked = 1,     // 锁定 (未达成)
        Available = 2,  // 可领取
        Claimed = 3,    // 已领取
    }

    // ============================================================
    // 活动信息
    // ============================================================
    
    [Serializable]
    public class ActivityInfo
    {
        public string id;                    // 活动ID
        public string title;                  // 活动标题
        public string description;             // 活动描述
        public ActivityType type;             // 活动类型
        public ActivityStatus status;          // 活动状态
        public long startTime;                // 开始时间
        public long endTime;                  // 结束时间
        public string icon;                   // 图标
        public string banner;                 // 横幅
        public bool showCountdown;            // 显示倒计时
        public int priority;                  // 优先级
        public Dictionary<string, string> extData; // 扩展数据
    }

    // ============================================================
    // 活动奖励
    // ============================================================
    
    [Serializable]
    public class ActivityReward
    {
        public string rewardId;               // 奖励ID
        public string name;                   // 奖励名称
        public int type;                      // 奖励类型 (1=金币,2=钻石,3=道具,4=皮肤)
        public int itemId;                    // 道具ID
        public int count;                    // 数量
        public string icon;                   // 图标
        public long requirement;              // 领取条件 (数值)
        public RewardStatus status;           // 状态
    }

    // ============================================================
    // 活动进度
    // ============================================================
    
    [Serializable]
    public class ActivityProgress
    {
        public string activityId;             // 活动ID
        public long currentValue;             // 当前进度值
        public long targetValue;              // 目标值
        public float progress;                // 进度百分比 (0-1)
        public int claimedRewards;           // 已领取奖励数
        public int totalRewards;              // 总奖励数
    }

    // ============================================================
    // 活动数据
    // ============================================================
    
    [Serializable]
    public class ActivityData
    {
        public ActivityInfo info;                           // 活动信息
        public List<ActivityReward> rewards;                // 奖励列表
        public ActivityProgress progress;                   // 进度
        public int remainingSeconds;                        // 剩余秒数
    }

    // ============================================================
    // 活动客户端主类
    // ============================================================
    
    public class ActivityClient : MonoBehaviour
    {
        private static ActivityClient _instance;
        public static ActivityClient Instance
        {
            get
            {
                if (_instance == null)
                {
                    _instance = FindObjectOfType<ActivityClient>();
                }
                return _instance;
            }
        }

        // 活动缓存
        private Dictionary<string, ActivityData> _activities = new Dictionary<string, ActivityData>();
        private List<ActivityType> _subscribedTypes = new List<ActivityType>();
        
        // 定时刷新
        private float _lastRefreshTime;
        private const float REFRESH_INTERVAL = 30f;
        
        // 事件
        public event Action<string> OnActivityUpdated;      // 活动更新
        public event Action<string> OnRewardClaimed;        // 奖励领取
        public event Action<string> OnActivityEnded;         // 活动结束
        public event Action<string> OnError;

        // ============================================================
        // 初始化
        // ============================================================
        
        private void Awake()
        {
            _instance = this;
        }

        private void Update()
        {
            // 定时刷新活动状态
            if (Time.time - _lastRefreshTime > REFRESH_INTERVAL)
            {
                _lastRefreshTime = Time.time;
                RefreshActivityStates();
            }
        }

        // ============================================================
        // 请求活动列表
        // ============================================================
        
        public void RequestActivityList(ActivityType type = ActivityType.Daily)
        {
            var req = new CSGetActivityList
            {
                Type = (int)type
            };
            
            NetworkManager.Instance.Send(req, (sc) =>
            {
                var response = sc as SCActivityList;
                if (response != null)
                {
                    ProcessActivityList(response);
                }
            });
        }

        // ============================================================
        // 请求所有活动
        // ============================================================
        
        public void RequestAllActivities()
        {
            var req = new CSGetAllActivities();
            NetworkManager.Instance.Send(req, (sc) =>
            {
                var response = sc as SCActivityList;
                if (response != null)
                {
                    ProcessActivityList(response);
                }
            });
        }

        // ============================================================
        // 请求活动详情
        // ============================================================
        
        public void RequestActivityDetail(string activityId)
        {
            var req = new CSGetActivityDetail
            {
                ActivityId = activityId
            };
            
            NetworkManager.Instance.Send(req, (sc) =>
            {
                var response = sc as SCActivityDetail;
                if (response != null)
                {
                    ProcessActivityDetail(response);
                }
            });
        }

        // ============================================================
        // 领取活动奖励
        // ============================================================
        
        public void ClaimReward(string activityId, string rewardId)
        {
            var req = new CSClaimActivityReward
            {
                ActivityId = activityId,
                RewardId = rewardId
            };
            
            NetworkManager.Instance.Send(req, (sc) =>
            {
                var response = sc as SCClaimActivityRewardResult;
                if (response != null)
                {
                    if (response.Success)
                    {
                        // 更新本地状态
                        if (_activities.TryGetValue(activityId, out var activity))
                        {
                            var reward = activity.rewards.FirstOrDefault(r => r.rewardId == rewardId);
                            if (reward != null)
                            {
                                reward.status = RewardStatus.Claimed;
                            }
                            activity.progress.claimedRewards++;
                        }
                        
                        OnRewardClaimed?.Invoke(activityId);
                    }
                    else
                    {
                        OnError?.Invoke(response.Message);
                    }
                }
            });
        }

        // ============================================================
        // 订阅活动类型
        // ============================================================
        
        public void SubscribeActivityType(ActivityType type)
        {
            if (!_subscribedTypes.Contains(type))
            {
                _subscribedTypes.Add(type);
            }
            
            var req = new CSSubscribeActivity
            {
                Type = (int)type,
                Subscribe = true
            };
            
            NetworkManager.Instance.Send(req);
        }

        // ============================================================
        // 取消订阅活动类型
        // ============================================================
        
        public void UnsubscribeActivityType(ActivityType type)
        {
            _subscribedTypes.Remove(type);
            
            var req = new CSSubscribeActivity
            {
                Type = (int)type,
                Subscribe = false
            };
            
            NetworkManager.Instance.Send(req);
        }

        // ============================================================
        // 获取活动
        // ============================================================
        
        public ActivityData GetActivity(string activityId)
        {
            if (_activities.TryGetValue(activityId, out var activity))
            {
                return activity;
            }
            return null;
        }

        // ============================================================
        // 获取进行中的活动
        // ============================================================
        
        public List<ActivityData> GetActiveActivities()
        {
            return _activities.Values
                .Where(a => a.info.status == ActivityStatus.Active)
                .OrderByDescending(a => a.info.priority)
                .ToList();
        }

        // ============================================================
        // 获取可领取奖励的活动
        // ============================================================
        
        public List<ActivityData> GetClaimableActivities()
        {
            return _activities.Values
                .Where(a => a.info.status == ActivityStatus.Active && 
                           a.rewards.Any(r => r.status == RewardStatus.Available))
                .ToList();
        }

        // ============================================================
        // 获取活动数量
        // ============================================================
        
        public int GetActiveCount()
        {
            return _activities.Values.Count(a => a.info.status == ActivityStatus.Active);
        }

        // ============================================================
        // 获取可领取数量
        // ============================================================
        
        public int GetClaimableCount()
        {
            return _activities.Values
                .Sum(a => a.rewards.Count(r => r.status == RewardStatus.Available));
        }

        // ============================================================
        // 刷新活动状态
        // ============================================================
        
        private void RefreshActivityStates()
        {
            long now = DateTimeOffset.UtcNow.ToUnixTimeSeconds();
            
            foreach (var kvp in _activities)
            {
                var activity = kvp.Value;
                var oldStatus = activity.info.status;
                
                // 更新状态
                if (now < activity.info.startTime)
                {
                    activity.info.status = ActivityStatus.Upcoming;
                    activity.remainingSeconds = (int)(activity.info.startTime - now);
                }
                else if (now >= activity.info.startTime && now < activity.info.endTime)
                {
                    activity.info.status = ActivityStatus.Active;
                    activity.remainingSeconds = (int)(activity.info.endTime - now);
                }
                else
                {
                    activity.info.status = ActivityStatus.Ended;
                    activity.remainingSeconds = 0;
                }
                
                // 状态变化通知
                if (oldStatus != activity.info.status)
                {
                    if (activity.info.status == ActivityStatus.Ended)
                    {
                        OnActivityEnded?.Invoke(kvp.Key);
                    }
                    else
                    {
                        OnActivityUpdated?.Invoke(kvp.Key);
                    }
                }
            }
        }

        // ============================================================
        // 处理活动列表
        // ============================================================
        
        private void ProcessActivityList(SCActivityList response)
        {
            if (response.Activities == null) return;
            
            foreach (var proto in response.Activities)
            {
                var activity = ConvertToActivityData(proto);
                _activities[activity.info.id] = activity;
                OnActivityUpdated?.Invoke(activity.info.id);
            }
        }

        // ============================================================
        // 处理活动详情
        // ============================================================
        
        private void ProcessActivityDetail(SCActivityDetail response)
        {
            if (response.Activity != null)
            {
                var activity = ConvertToActivityData(response.Activity);
                _activities[activity.info.id] = activity;
                OnActivityUpdated?.Invoke(activity.info.id);
            }
        }

        // ============================================================
        // 转换为活动数据
        // ============================================================
        
        private ActivityData ConvertToActivityData(ActivityProto proto)
        {
            var data = new ActivityData
            {
                info = new ActivityInfo
                {
                    id = proto.Id,
                    title = proto.Title,
                    description = proto.Description,
                    type = (ActivityType)proto.Type,
                    status = (ActivityStatus)proto.Status,
                    startTime = proto.StartTime,
                    endTime = proto.EndTime,
                    icon = proto.Icon,
                    banner = proto.Banner,
                    showCountdown = proto.ShowCountdown,
                    priority = proto.Priority
                },
                rewards = new List<ActivityReward>(),
                progress = new ActivityProgress
                {
                    activityId = proto.Id,
                    currentValue = proto.CurrentValue,
                    targetValue = proto.TargetValue,
                    progress = proto.TargetValue > 0 ? (float)proto.CurrentValue / proto.TargetValue : 0,
                    claimedRewards = proto.ClaimedRewards,
                    totalRewards = proto.TotalRewards
                }
            };
            
            // 计算剩余时间
            long now = DateTimeOffset.UtcNow.ToUnixTimeSeconds();
            if (now < data.info.startTime)
            {
                data.remainingSeconds = (int)(data.info.startTime - now);
            }
            else if (now < data.info.endTime)
            {
                data.remainingSeconds = (int)(data.info.endTime - now);
            }
            else
            {
                data.remainingSeconds = 0;
            }
            
            // 转换奖励
            if (proto.Rewards != null)
            {
                foreach (var r in proto.Rewards)
                {
                    data.rewards.Add(new ActivityReward
                    {
                        rewardId = r.RewardId,
                        name = r.Name,
                        type = r.Type,
                        itemId = r.ItemId,
                        count = r.Count,
                        icon = r.Icon,
                        requirement = r.Requirement,
                        status = (RewardStatus)r.Status
                    });
                }
            }
            
            return data;
        }

        // ============================================================
        // 处理活动进度更新推送
        // ============================================================
        
        public void HandleActivityProgressUpdate(SCActivityProgressUpdate update)
        {
            if (_activities.TryGetValue(update.ActivityId, out var activity))
            {
                activity.progress.currentValue = update.CurrentValue;
                activity.progress.progress = activity.progress.targetValue > 0 
                    ? (float)update.CurrentValue / activity.progress.targetValue 
                    : 0;
                
                // 检查可领取状态
                foreach (var reward in activity.rewards)
                {
                    if (reward.status == RewardStatus.Locked && 
                        activity.progress.currentValue >= reward.requirement)
                    {
                        reward.status = RewardStatus.Available;
                    }
                }
                
                OnActivityUpdated?.Invoke(update.ActivityId);
            }
        }

        // ============================================================
        // 获取倒计时字符串
        // ============================================================
        
        public static string GetCountdownString(int seconds)
        {
            if (seconds <= 0) return "已结束";
            
            int days = seconds / 86400;
            int hours = (seconds % 86400) / 3600;
            int minutes = (seconds % 3600) / 60;
            int secs = seconds % 60;
            
            if (days > 0)
            {
                return $"{days}天{hours}小时";
            }
            else if (hours > 0)
            {
                return $"{hours}小时{minutes}分钟";
            }
            else if (minutes > 0)
            {
                return $"{minutes}分{secs}秒";
            }
            else
            {
                return $"{secs}秒";
            }
        }
    }

    // ============================================================
    // 活动 UI 管理器
    // ============================================================
    
    public class ActivityUIManager : MonoBehaviour
    {
        [Header("活动列表容器")]
        public Transform activityListContainer;
        
        [Header("活动条目预制体")]
        public GameObject activityItemPrefab;
        
        [Header("活动详情面板")]
        public GameObject detailPanel;
        
        [Header("Loading")]
        public GameObject loadingObj;
        
        [Header("红点")]
        public GameObject redDot;

        private List<ActivityData> _currentActivities = new List<ActivityData>();

        // ============================================================
        // 初始化
        // ============================================================
        
        private void Start()
        {
            ActivityClient.Instance.OnActivityUpdated += OnActivityUpdated;
            ActivityClient.Instance.OnRewardClaimed += OnRewardClaimed;
            
            // 请求活动列表
            RequestActivities();
        }

        private void OnDestroy()
        {
            ActivityClient.Instance.OnActivityUpdated -= OnActivityUpdated;
            ActivityClient.Instance.OnRewardClaimed -= OnRewardClaimed;
        }

        // ============================================================
        // 请求活动列表
        // ============================================================
        
        public void RequestActivities()
        {
            ShowLoading(true);
            ActivityClient.Instance.RequestAllActivities();
        }

        // ============================================================
        // 刷新
        // ============================================================
        
        public void Refresh()
        {
            RequestActivities();
        }

        // ============================================================
        // 更新UI
        // ============================================================
        
        private void OnActivityUpdated(string activityId)
        {
            ShowLoading(false);
            UpdateActivityList();
            UpdateRedDot();
        }

        // ============================================================
        // 奖励领取回调
        // ============================================================
        
        private void OnRewardClaimed(string activityId)
        {
            // 刷新详情
            ActivityClient.Instance.RequestActivityDetail(activityId);
            UpdateRedDot();
        }

        // ============================================================
        // 更新活动列表
        // ============================================================
        
        private void UpdateActivityList()
        {
            _currentActivities = ActivityClient.Instance.GetActiveActivities();
            
            // 清空容器
            foreach (Transform child in activityListContainer)
            {
                Destroy(child.gameObject);
            }
            
            // 生成活动条目
            foreach (var activity in _currentActivities)
            {
                var go = Instantiate(activityItemPrefab, activityListContainer);
                var item = go.GetComponent<ActivityItem>();
                item.Initialize(activity);
                item.OnClick += OnActivityItemClick;
            }
        }

        // ============================================================
        // 点击活动条目
        // ============================================================
        
        private void OnActivityItemClick(ActivityData activity)
        {
            ShowActivityDetail(activity);
        }

        // ============================================================
        // 显示活动详情
        // ============================================================
        
        private void ShowActivityDetail(ActivityData activity)
        {
            if (detailPanel != null)
            {
                var detail = detailPanel.GetComponent<ActivityDetailPanel>();
                if (detail != null)
                {
                    detail.Show(activity);
                }
            }
        }

        // ============================================================
        // 更新红点
        // ============================================================
        
        private void UpdateRedDot()
        {
            if (redDot != null)
            {
                int claimableCount = ActivityClient.Instance.GetClaimableCount();
                redDot.SetActive(claimableCount > 0);
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
    }

    // ============================================================
    // 活动条目 Item
    // ============================================================
    
    public class ActivityItem : MonoBehaviour
    {
        [Header("图标")]
        public Image iconImage;
        
        [Header("标题")]
        public Text titleText;
        
        [Header("倒计时")]
        public Text countdownText;
        
        [Header("进度条")]
        public Image progressBar;
        
        [Header("进度文本")]
        public Text progressText;
        
        [Header("可领取标识")]
        public GameObject claimableObj;
        
        [Header("进行中标识")]
        public GameObject activeObj;

        private ActivityData _activity;
        public event Action<ActivityData> OnClick;

        // ============================================================
        // 初始化
        // ============================================================
        
        public void Initialize(ActivityData activity)
        {
            _activity = activity;
            
            // 标题
            if (titleText != null)
            {
                titleText.text = activity.info.title;
            }
            
            // 倒计时
            if (countdownText != null)
            {
                countdownText.text = ActivityClient.GetCountdownString(activity.remainingSeconds);
            }
            
            // 进度
            if (progressBar != null)
            {
                progressBar.fillAmount = activity.progress.progress;
            }
            
            if (progressText != null)
            {
                progressText.text = $"{activity.progress.currentValue}/{activity.progress.targetValue}";
            }
            
            // 可领取
            bool canClaim = activity.rewards.Any(r => r.status == RewardStatus.Available);
            if (claimableObj != null)
            {
                claimableObj.SetActive(canClaim);
            }
            
            // 进行中
            if (activeObj != null)
            {
                activeObj.SetActive(activity.info.status == ActivityStatus.Active && !canClaim);
            }
            
            // 加载图标
            LoadIcon(activity.info.icon);
        }

        // ============================================================
        // 加载图标
        // ============================================================
        
        private void LoadIcon(string iconUrl)
        {
            if (iconImage == null || string.IsNullOrEmpty(iconUrl)) return;
            // TODO: 实现图标加载
        }

        // ============================================================
        // 点击
        // ============================================================
        
        public void OnItemClick()
        {
            OnClick?.Invoke(_activity);
        }
    }

    // ============================================================
    // 活动详情面板
    // ============================================================
    
    public class ActivityDetailPanel : MonoBehaviour
    {
        [Header("标题")]
        public Text titleText;
        
        [Header("描述")]
        public Text descText;
        
        [Header("倒计时")]
        public Text countdownText;
        
        [Header("进度条")]
        public Image progressBar;
        
        [Header("进度文本")]
        public Text progressText;
        
        [Header("奖励容器")]
        public Transform rewardContainer;
        
        [Header("奖励预制体")]
        public GameObject rewardPrefab;

        private ActivityData _activity;

        // ============================================================
        // 显示
        // ============================================================
        
        public void Show(ActivityData activity)
        {
            _activity = activity;
            gameObject.SetActive(true);
            
            // 标题
            if (titleText != null)
            {
                titleText.text = activity.info.title;
            }
            
            // 描述
            if (descText != null)
            {
                descText.text = activity.info.description;
            }
            
            // 倒计时
            if (countdownText != null)
            {
                countdownText.text = ActivityClient.GetCountdownString(activity.remainingSeconds);
            }
            
            // 进度
            if (progressBar != null)
            {
                progressBar.fillAmount = activity.progress.progress;
            }
            
            if (progressText != null)
            {
                progressText.text = $"{activity.progress.currentValue}/{activity.progress.targetValue}";
            }
            
            // 奖励
            UpdateRewards();
        }

        // ============================================================
        // 隐藏
        // ============================================================
        
        public void Hide()
        {
            gameObject.SetActive(false);
        }

        // ============================================================
        // 更新奖励列表
        // ============================================================
        
        private void UpdateRewards()
        {
            if (rewardContainer == null || _activity == null) return;
            
            // 清空
            foreach (Transform child in rewardContainer)
            {
                Destroy(child.gameObject);
            }
            
            // 生成奖励
            foreach (var reward in _activity.rewards)
            {
                var go = Instantiate(rewardPrefab, rewardContainer);
                var item = go.GetComponent<RewardItem>();
                item.Initialize(reward, _activity.progress.currentValue);
                item.OnClaim += OnClaimReward;
            }
        }

        // ============================================================
        // 领取奖励
        // ============================================================
        
        private void OnClaimReward(string rewardId)
        {
            if (_activity != null)
            {
                ActivityClient.Instance.ClaimReward(_activity.info.id, rewardId);
            }
        }
    }

    // ============================================================
    // 奖励条目
    // ============================================================
    
    public class RewardItem : MonoBehaviour
    {
        [Header("图标")]
        public Image iconImage;
        
        [Header("名称")]
        public Text nameText;
        
        [Header("数量")]
        public Text countText;
        
        [Header("进度需求")]
        public Text requirementText;
        
        [Header("领取按钮")]
        public Button claimButton;
        
        [Header("已领取标识")]
        public GameObject claimedObj;
        
        [Header("锁定标识")]
        public GameObject lockedObj;

        private ActivityReward _reward;
        private long _currentValue;
        
        public event Action<string> OnClaim;

        // ============================================================
        // 初始化
        // ============================================================
        
        public void Initialize(ActivityReward reward, long currentValue)
        {
            _reward = reward;
            _currentValue = currentValue;
            
            // 名称
            if (nameText != null)
            {
                nameText.text = reward.name;
            }
            
            // 数量
            if (countText != null)
            {
                countText.text = $"x{reward.count}";
            }
            
            // 需求
            if (requirementText != null)
            {
                requirementText.text = reward.requirement.ToString();
            }
            
            // 状态
            bool canClaim = reward.status == RewardStatus.Available;
            bool claimed = reward.status == RewardStatus.Claimed;
            bool locked = reward.status == RewardStatus.Locked;
            
            if (claimButton != null)
            {
                claimButton.gameObject.SetActive(canClaim);
                claimButton.interactable = canClaim;
            }
            
            if (claimedObj != null)
            {
                claimedObj.SetActive(claimed);
            }
            
            if (lockedObj != null)
            {
                lockedObj.SetActive(locked);
            }
            
            // 加载图标
            LoadIcon(reward.icon);
        }

        // ============================================================
        // 加载图标
        // ============================================================
        
        private void LoadIcon(string iconUrl)
        {
            if (iconImage == null || string.IsNullOrEmpty(iconUrl)) return;
            // TODO: 实现图标加载
        }

        // ============================================================
        // 点击领取
        // ============================================================
        
        public void OnClaimClick()
        {
            OnClaim?.Invoke(_reward.rewardId);
        }
    }

    // ============================================================
    // 网络协议
    // ============================================================
    
    public class CSGetActivityList
    {
        public int Type { get; set; }
    }
    
    public class CSGetAllActivities
    {
    }
    
    public class CSGetActivityDetail
    {
        public string ActivityId { get; set; }
    }
    
    public class CSClaimActivityReward
    {
        public string ActivityId { get; set; }
        public string RewardId { get; set; }
    }
    
    public class CSSubscribeActivity
    {
        public int Type { get; set; }
        public bool Subscribe { get; set; }
    }
    
    public class SCActivityList
    {
        public List<ActivityProto> Activities { get; set; }
    }
    
    public class SCActivityDetail
    {
        public ActivityProto Activity { get; set; }
    }
    
    public class SCClaimActivityRewardResult
    {
        public bool Success { get; set; }
        public string Message { get; set; }
    }
    
    public class SCActivityProgressUpdate
    {
        public string ActivityId { get; set; }
        public long CurrentValue { get; set; }
    }
    
    public class ActivityProto
    {
        public string Id { get; set; }
        public string Title { get; set; }
        public string Description { get; set; }
        public int Type { get; set; }
        public int Status { get; set; }
        public long StartTime { get; set; }
        public long EndTime { get; set; }
        public string Icon { get; set; }
        public string Banner { get; set; }
        public bool ShowCountdown { get; set; }
        public int Priority { get; set; }
        public long CurrentValue { get; set; }
        public long TargetValue { get; set; }
        public int ClaimedRewards { get; set; }
        public int TotalRewards { get; set; }
        public List<RewardProto> Rewards { get; set; }
    }
    
    public class RewardProto
    {
        public string RewardId { get; set; }
        public string Name { get; set; }
        public int Type { get; set; }
        public int ItemId { get; set; }
        public int Count { get; set; }
        public string Icon { get; set; }
        public long Requirement { get; set; }
        public int Status { get; set; }
    }
}
