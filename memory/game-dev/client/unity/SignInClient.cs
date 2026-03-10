// ============================================
// Unity 签到客户端 - 弹幕游戏
// 对接后端签到系统 WebSocket + Protobuf 协议
// ============================================

using System;
using System.Collections.Generic;
using UnityEngine;
using UnityEngine.UI;

namespace DanmakuGame.Client
{
    // ============================================
    // 签到消息ID
    // ============================================
    public enum SignInMsgID : ushort
    {
        CS_GET_SIGNIN_STATUS = 6201,
        SC_SIGNIN_STATUS = 6202,
        CS_GET_SIGNIN_CALENDAR = 6203,
        SC_SIGNIN_CALENDAR = 6204,
        CS_SIGNIN = 6205,
        SC_SIGNIN_RESULT = 6206,
        CS_GET_SIGNIN_RANK = 6207,
        SC_SIGNIN_RANK = 6208,
    }

    // ============================================
    // 签到状态数据模型
    // ============================================
    [Serializable]
    public class SignInStatus
    {
        public bool can_sign_in;       // 今日是否可以签到
        public bool is_signed_today;   // 今日是否已签到
        public int consecutive_days;   // 连续签到天数
        public int total_days;         // 总签到天数
        public int month_days;         // 本月签到天数
        public int current_month;      // 当前月份
        public List<SignInReward> today_rewards;    // 今日奖励
        public List<SignInReward> tomorrow_rewards;  // 明日奖励
    }

    // ============================================
    // 签到奖励数据模型
    // ============================================
    [Serializable]
    public class SignInReward
    {
        public string reward_id;
        public string reward_name;
        public int reward_type;    // 1=金币, 2=钻石, 3=道具
        public int count;          // 数量
        public bool is_vip;        // 是否VIP专属
    }

    // ============================================
    // 签到日历数据模型
    // ============================================
    [Serializable]
    public class SignInCalendar
    {
        public int year;
        public int month;
        public List<SignInDay> days;
    }

    [Serializable]
    public class SignInDay
    {
        public int day;                    // 日期 (1-31)
        public bool is_signed;             // 是否已签到
        public List<SignInReward> rewards; // 签到奖励
        public bool is_extra_reward;       // 是否有额外奖励
    }

    // ============================================
    // 签到排行榜数据模型
    // ============================================
    [Serializable]
    public class SignInRankItem
    {
        public int rank;
        public string player_id;
        public string player_name;
        public string avatar_url;
        public int consecutive_days;   // 连续签到
        public int total_days;         // 总签到天数
    }

    // ============================================
    // 获取签到状态请求/响应
    // ============================================
    [Serializable]
    public class CSGetSignInStatusReq
    {
        public string player_id;
    }

    [Serializable]
    public class SCSignInStatusResp
    {
        public bool success;
        public string error_msg;
        public SignInStatus status;
    }

    // ============================================
    // 获取签到日历请求/响应
    // ============================================
    [Serializable]
    public class CSGetSignInCalendarReq
    {
        public string player_id;
        public int year;
        public int month;
    }

    [Serializable]
    public class SCSignInCalendarResp
    {
        public bool success;
        public string error_msg;
        public SignInCalendar calendar;
    }

    // ============================================
    // 签到请求/响应
    // ============================================
    [Serializable]
    public class CSSignInReq
    {
        public string player_id;
    }

    [Serializable]
    public class SCSignInResultResp
    {
        public bool success;
        public string error_msg;
        public bool is_signed_today;
        public int consecutive_days;
        public int total_days;
        public List<SignInReward> rewards;       // 获得的奖励
        public List<SignInReward> extra_rewards; // 额外奖励(连续签到)
        public int gold_change;                  // 金币变化
        public int gem_change;                   // 钻石变化
    }

    // ============================================
    // 获取签到排行榜请求/响应
    // ============================================
    [Serializable]
    public class CSGetSignInRankReq
    {
        public string player_id;
        public int rank_type; // 1=连续签到排行, 2=总签到排行
        public int page = 1;
        public int page_size = 20;
    }

    [Serializable]
    public class SCSignInRankResp
    {
        public bool success;
        public string error_msg;
        public int rank_type;
        public List<SignInRankItem> ranks;
        public int player_rank;     // 自己的排名
    }

    // ============================================
    // 签到客户端管理器
    // ============================================
    public class SignInClient : MonoBehaviour
    {
        private static SignInClient _instance;
        public static SignInClient Instance
        {
            get
            {
                if (_instance == null)
                {
                    GameObject obj = new GameObject("SignInClient");
                    _instance = obj.AddComponent<SignInClient>();
                    DontDestroyOnLoad(obj);
                }
                return _instance;
            }
        }

        // 签到状态缓存
        private SignInStatus _cachedStatus = null;

        // 签到日历缓存
        private Dictionary<string, SignInCalendar> _cachedCalendars = new Dictionary<string, SignInCalendar>();

        // 事件回调
        public Action<SignInStatus> OnStatusReceived;
        public Action<SignInCalendar> OnCalendarReceived;
        public Action<SCSignInResultResp> OnSignInSuccess;
        public Action<bool, string> OnSignInFailed;
        public Action<List<SignInRankItem>, int> OnRankReceived;

        // 关联的主客户端
        private DanmakuGameClient _mainClient;

        public void Initialize(DanmakuGameClient mainClient)
        {
            _mainClient = mainClient;
        }

        // ============================================
        // 获取签到状态
        // ============================================
        public void GetStatus(string playerId)
        {
            CSGetSignInStatusReq req = new CSGetSignInStatusReq
            {
                player_id = playerId
            };

            _mainClient.SendMessage((ushort)SignInMsgID.CS_GET_SIGNIN_STATUS, req);
        }

        // 处理签到状态响应
        public void HandleSignInStatusResp(SCSignInStatusResp resp)
        {
            if (resp.success)
            {
                _cachedStatus = resp.status;
                OnStatusReceived?.Invoke(resp.status);
            }
            else
            {
                OnSignInFailed?.Invoke(false, resp.error_msg);
            }
        }

        // ============================================
        // 获取签到日历
        // ============================================
        public void GetCalendar(string playerId, int year, int month)
        {
            CSGetSignInCalendarReq req = new CSGetSignInCalendarReq
            {
                player_id = playerId,
                year = year,
                month = month
            };

            _mainClient.SendMessage((ushort)SignInMsgID.CS_GET_SIGNIN_CALENDAR, req);
        }

        // 处理签到日历响应
        public void HandleSignInCalendarResp(SCSignInCalendarResp resp)
        {
            if (resp.success)
            {
                string key = $"{resp.calendar.year}-{resp.calendar.month}";
                _cachedCalendars[key] = resp.calendar;
                OnCalendarReceived?.Invoke(resp.calendar);
            }
            else
            {
                OnSignInFailed?.Invoke(false, resp.error_msg);
            }
        }

        // ============================================
        // 执行签到
        // ============================================
        public void SignIn(string playerId)
        {
            CSSignInReq req = new CSSignInReq
            {
                player_id = playerId
            };

            _mainClient.SendMessage((ushort)SignInMsgID.CS_SIGNIN, req);
        }

        // 处理签到响应
        public void HandleSignInResultResp(SCSignInResultResp resp)
        {
            if (resp.success)
            {
                // 更新缓存
                if (_cachedStatus != null)
                {
                    _cachedStatus.is_signed_today = resp.is_signed_today;
                    _cachedStatus.consecutive_days = resp.consecutive_days;
                    _cachedStatus.total_days = resp.total_days;
                }

                OnSignInSuccess?.Invoke(resp);
            }
            else
            {
                OnSignInFailed?.Invoke(false, resp.error_msg);
            }
        }

        // ============================================
        // 获取签到排行榜
        // ============================================
        public void GetRank(string playerId, int rankType = 1, int page = 1, int pageSize = 20)
        {
            CSGetSignInRankReq req = new CSGetSignInRankReq
            {
                player_id = playerId,
                rank_type = rankType,
                page = page,
                page_size = pageSize
            };

            _mainClient.SendMessage((ushort)SignInMsgID.CS_GET_SIGNIN_RANK, req);
        }

        // 处理排行榜响应
        public void HandleSignInRankResp(SCSignInRankResp resp)
        {
            if (resp.success)
            {
                OnRankReceived?.Invoke(resp.ranks, resp.player_rank);
            }
            else
            {
                OnSignInFailed?.Invoke(false, resp.error_msg);
            }
        }

        // ============================================
        // 辅助方法
        // ============================================

        // 获取缓存的签到状态
        public SignInStatus GetCachedStatus()
        {
            return _cachedStatus;
        }

        // 获取缓存的日历
        public SignInCalendar GetCachedCalendar(int year, int month)
        {
            string key = $"{year}-{month}";
            if (_cachedCalendars.ContainsKey(key))
            {
                return _cachedCalendars[key];
            }
            return null;
        }

        // 检查今日是否可以签到
        public bool CanSignInToday()
        {
            if (_cachedStatus == null) return false;
            return _cachedStatus.can_sign_in && !_cachedStatus.is_signed_today;
        }

        // 获取连续签到加成百分比
        public int GetConsecutiveBonusPercent()
        {
            if (_cachedStatus == null) return 0;

            // 每7天连续签到获得额外加成
            int consecutive = _cachedStatus.consecutive_days;
            if (consecutive >= 28) return 100;      // 30天签到满
            if (consecutive >= 21) return 80;
            if (consecutive >= 14) return 60;
            if (consecutive >= 7) return 40;
            return 0;
        }

        // 获取签到状态描述
        public string GetStatusDesc()
        {
            if (_cachedStatus == null) return "未登录";
            if (_cachedStatus.is_signed_today) return "今日已签到";
            if (_cachedStatus.can_sign_in) return "可以签到";
            return "签到不可用";
        }

        // 获取奖励类型描述
        public string GetRewardTypeDesc(int rewardType)
        {
            switch (rewardType)
            {
                case 1: return "金币";
                case 2: return "钻石";
                case 3: return "道具";
                default: return "未知";
            }
        }

        // 获取排行榜类型描述
        public string GetRankTypeDesc(int rankType)
        {
            switch (rankType)
            {
                case 1: return "连续签到榜";
                case 2: return "总签到榜";
                default: return "未知";
            }
        }

        // 格式化签到结果消息
        public string FormatSignInResult(SCSignInResultResp resp)
        {
            string msg = $"签到成功！\n";
            msg += $"连续签到: {resp.consecutive_days}天\n";
            msg += $"总签到: {resp.total_days}天\n";

            if (resp.rewards != null && resp.rewards.Count > 0)
            {
                msg += "获得: ";
                foreach (var r in resp.rewards)
                {
                    msg += $"{r.reward_name}x{r.count} ";
                }
                msg += "\n";
            }

            if (resp.extra_rewards != null && resp.extra_rewards.Count > 0)
                {
                msg += "额外奖励: ";
                foreach (var r in resp.extra_rewards)
                {
                    msg += $"{r.reward_name}x{r.count} ";
                }
            }

            return msg;
        }
    }
}
