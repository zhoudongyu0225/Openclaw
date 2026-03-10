// ============================================
// Unity 商城客户端 - 弹幕游戏
// 对接后端商城系统 WebSocket + Protobuf 协议
// ============================================

using System;
using System.Collections.Generic;
using UnityEngine;
using UnityEngine.UI;

namespace DanmakuGame.Client
{
    // ============================================
    // 商城消息ID
    // ============================================
    public enum MallMsgID : ushort
    {
        CS_GET_MALL_ITEMS = 6101,
        SC_MALL_ITEMS = 6102,
        CS_PURCHASE = 6103,
        SC_PURCHASE_RESULT = 6104,
        CS_GET_BALANCE = 6105,
        SC_BALANCE = 6106,
        CS_GET_TODAY_SPEND = 6107,
        SC_TODAY_SPEND = 6108,
        CS_GET_PURCHASE_HISTORY = 6109,
        SC_PURCHASE_HISTORY = 6110,
    }

    // ============================================
    // 商城类型
    // ============================================
    public enum MallType : byte
    {
        Gift = 1,      // 礼物商城
        Item = 2,      // 道具商城
        Skin = 3,      // 皮肤商城
        Random = 4,    // 随机商城
        Honor = 5,     // 荣誉商城
    }

    // ============================================
    // 货币类型
    // ============================================
    public enum CurrencyType : byte
    {
        Gold = 1,    // 金币
        Gem = 2,     // 钻石
        Honor = 3,   // 荣誉点
        Credit = 4,  // 积分
    }

    // ============================================
    // 商品数据模型
    // ============================================
    [Serializable]
    public class MallItem
    {
        public string item_id;
        public string item_name;
        public string description;
        public int item_type;      // 1=金币, 2=钻石, 3=道具, 4=皮肤, 5=礼物
        public MallType mall_type;
        public CurrencyType currency_type;
        public int price;          // 价格
        public int original_price; // 原价(折扣用)
        public float discount;     // 折扣率 (0.0-1.0)
        public int stock;           // 库存 (-1表示无限)
        public int purchase_limit;  // 购买限制 (-1表示无限)
        public int purchased_count;// 已购买数量
        public int required_level; // 所需等级
        public bool is_vip_only;   // VIP专属
        public bool is_hot;        // 热门
        public bool is_new;        // 新品
        public long start_time;    // 上架时间
        public long end_time;      // 下架时间
        public string icon_url;    // 图标URL
    }

    // ============================================
    // 玩家余额
    // ============================================
    [Serializable]
    public class PlayerBalance
    {
        public int gold;     // 金币
        public int gem;      // 钻石
        public int honor;    // 荣誉点
        public int credit;   // 积分
        public int vip_level; // VIP等级
    }

    // ============================================
    // 购买历史记录
    // ============================================
    [Serializable]
    public class PurchaseRecord
    {
        public string order_id;
        public string item_id;
        public string item_name;
        public int price;
        public CurrencyType currency_type;
        public int quantity;
        public long purchase_time;
    }

    // ============================================
    // 获取商品列表请求/响应
    // ============================================
    [Serializable]
    public class CSGetMallItemsReq
    {
        public string player_id;
        public MallType mall_type;
    }

    [Serializable]
    public class SCMallItemsResp
    {
        public bool success;
        public string error_msg;
        public MallType mall_type;
        public List<MallItem> items;
    }

    // ============================================
    // 购买请求/响应
    // ============================================
    [Serializable]
    public class CSPurchaseReq
    {
        public string player_id;
        public string item_id;
        public int quantity = 1;
    }

    [Serializable]
    public class SCPurchaseResp
    {
        public bool success;
        public string error_msg;
        public string order_id;
        public string item_id;
        public string item_name;
        public int price;
        public CurrencyType currency_type;
        public int quantity;
        public PlayerBalance new_balance;
    }

    // ============================================
    // 获取余额请求/响应
    // ============================================
    [Serializable]
    public class CSGetBalanceReq
    {
        public string player_id;
    }

    [Serializable]
    public class SCBalanceResp
    {
        public bool success;
        public PlayerBalance balance;
    }

    // ============================================
    // 今日消费请求/响应
    // ============================================
    [Serializable]
    public class CSGetTodaySpendReq
    {
        public string player_id;
    }

    [Serializable]
    public class SCTodaySpendResp
    {
        public int gold_spent;
        public int gem_spent;
        public int honor_spent;
        public int credit_spent;
    }

    // ============================================
    // 购买历史请求/响应
    // ============================================
    [Serializable]
    public class CSGetPurchaseHistoryReq
    {
        public string player_id;
        public int page = 1;
        public int page_size = 20;
    }

    [Serializable]
    public class SCPurchaseHistoryResp
    {
        public bool success;
        public string error_msg;
        public List<PurchaseRecord> records;
        public int total_count;
    }

    // ============================================
    // 商城客户端管理器
    // ============================================
    public class MallClient : MonoBehaviour
    {
        private static MallClient _instance;
        public static MallClient Instance
        {
            get
            {
                if (_instance == null)
                {
                    GameObject obj = new GameObject("MallClient");
                    _instance = obj.AddComponent<MallClient>();
                    DontDestroyOnLoad(obj);
                }
                return _instance;
            }
        }

        // 当前商城类型
        private MallType _currentMallType = MallType.Gift;

        // 商品缓存
        private Dictionary<MallType, List<MallItem>> _cachedItems = new Dictionary<MallType, List<MallItem>>();

        // 余额缓存
        private PlayerBalance _cachedBalance = new PlayerBalance();

        // 事件回调
        public Action<MallType, List<MallItem>> OnItemsReceived;
        public Action<SCPurchaseResp> OnPurchaseSuccess;
        public Action<bool, string> OnPurchaseFailed;
        public Action<PlayerBalance> OnBalanceReceived;
        public Action<SCTodaySpendResp> OnTodaySpendReceived;
        public Action<List<PurchaseRecord>> OnHistoryReceived;

        // 关联的主客户端
        private DanmakuGameClient _mainClient;

        public void Initialize(DanmakuGameClient mainClient)
        {
            _mainClient = mainClient;
        }

        // ============================================
        // 获取商品列表
        // ============================================
        public void GetMallItems(string playerId, MallType mallType)
        {
            _currentMallType = mallType;

            CSGetMallItemsReq req = new CSGetMallItemsReq
            {
                player_id = playerId,
                mall_type = mallType
            };

            _mainClient.SendMessage((ushort)MallMsgID.CS_GET_MALL_ITEMS, req);
        }

        // 处理商品列表响应
        public void HandleMallItemsResp(SCMallItemsResp resp)
        {
            if (resp.success)
            {
                _cachedItems[resp.mall_type] = resp.items;
                OnItemsReceived?.Invoke(resp.mall_type, resp.items);
            }
            else
            {
                OnPurchaseFailed?.Invoke(false, resp.error_msg);
            }
        }

        // ============================================
        // 购买商品
        // ============================================
        public void Purchase(string playerId, string itemId, int quantity = 1)
        {
            CSPurchaseReq req = new CSPurchaseReq
            {
                player_id = playerId,
                item_id = itemId,
                quantity = quantity
            };

            _mainClient.SendMessage((ushort)MallMsgID.CS_PURCHASE, req);
        }

        // 处理购买响应
        public void HandlePurchaseResp(SCPurchaseResp resp)
        {
            if (resp.success)
            {
                // 更新余额缓存
                _cachedBalance = resp.new_balance;

                OnPurchaseSuccess?.Invoke(resp);
            }
            else
            {
                OnPurchaseFailed?.Invoke(false, resp.error_msg);
            }
        }

        // ============================================
        // 获取余额
        // ============================================
        public void GetBalance(string playerId)
        {
            CSGetBalanceReq req = new CSGetBalanceReq
            {
                player_id = playerId
            };

            _mainClient.SendMessage((ushort)MallMsgID.CS_GET_BALANCE, req);
        }

        // 处理余额响应
        public void HandleBalanceResp(SCBalanceResp resp)
        {
            if (resp.success)
            {
                _cachedBalance = resp.balance;
                OnBalanceReceived?.Invoke(resp.balance);
            }
        }

        // ============================================
        // 获取今日消费
        // ============================================
        public void GetTodaySpend(string playerId)
        {
            CSGetTodaySpendReq req = new CSGetTodaySpendReq
            {
                player_id = playerId
            };

            _mainClient.SendMessage((ushort)MallMsgID.CS_GET_TODAY_SPEND, req);
        }

        // 处理今日消费响应
        public void HandleTodaySpendResp(SCTodaySpendResp resp)
        {
            OnTodaySpendReceived?.Invoke(resp);
        }

        // ============================================
        // 获取购买历史
        // ============================================
        public void GetPurchaseHistory(string playerId, int page = 1, int pageSize = 20)
        {
            CSGetPurchaseHistoryReq req = new CSGetPurchaseHistoryReq
            {
                player_id = playerId,
                page = page,
                page_size = pageSize
            };

            _mainClient.SendMessage((ushort)MallMsgID.CS_GET_PURCHASE_HISTORY, req);
        }

        // 处理购买历史响应
        public void HandlePurchaseHistoryResp(SCPurchaseHistoryResp resp)
        {
            if (resp.success)
            {
                OnHistoryReceived?.Invoke(resp.records);
            }
            else
            {
                OnPurchaseFailed?.Invoke(false, resp.error_msg);
            }
        }

        // ============================================
        // 辅助方法
        // ============================================

        // 获取当前余额
        public PlayerBalance GetCachedBalance()
        {
            return _cachedBalance;
        }

        // 获取缓存的商品列表
        public List<MallItem> GetCachedItems(MallType mallType)
        {
            if (_cachedItems.ContainsKey(mallType))
            {
                return _cachedItems[mallType];
            }
            return new List<MallItem>();
        }

        // 检查是否可以购买
        public bool CanPurchase(MallItem item)
        {
            // 检查库存 (item.stock != -1 &&
            if item.purchased_count >= item.stock)
            {
                return false;
            }

            // 检查购买限制
            if (item.purchase_limit != -1 && item.purchased_count >= item.purchase_limit)
            {
                return false;
            }

            // 检查等级
            // 需要从玩家数据获取当前等级，这里简化处理
            // if (playerLevel < item.required_level) return false;

            // 检查余额
            switch (item.currency_type)
            {
                case CurrencyType.Gold:
                    return _cachedBalance.gold >= item.price;
                case CurrencyType.Gem:
                    return _cachedBalance.gem >= item.price;
                case CurrencyType.Honor:
                    return _cachedBalance.honor >= item.price;
                case CurrencyType.Credit:
                    return _cachedBalance.credit >= item.price;
            }

            return false;
        }

        // 获取货币类型描述
        public string GetCurrencyDesc(CurrencyType type)
        {
            switch (type)
            {
                case CurrencyType.Gold: return "金币";
                case CurrencyType.Gem: return "钻石";
                case CurrencyType.Honor: return "荣誉";
                case CurrencyType.Credit: return "积分";
                default: return "未知";
            }
        }

        // 获取商城类型描述
        public string GetMallTypeDesc(MallType type)
        {
            switch (type)
            {
                case MallType.Gift: return "礼物商城";
                case MallType.Item: return "道具商城";
                case MallType.Skin: return "皮肤商城";
                case MallType.Random: return "随机商城";
                case MallType.Honor: return "荣誉商城";
                default: return "未知";
            }
        }

        // 计算折扣价格
        public int GetDiscountPrice(MallItem item)
        {
            if (item.discount > 0 && item.discount < 1.0f)
            {
                return (int)(item.original_price * item.discount);
            }
            return item.price;
        }

        // 获取折扣标签
        public string GetDiscountTag(MallItem item)
        {
            if (item.discount > 0 && item.discount < 1.0f)
            {
                return $"{(int)(item.discount * 10)}折";
            }
            return null;
        }

        // 检查商品是否在促销中
        public bool IsOnSale(MallItem item)
        {
            long now = DateTimeOffset.UtcNow.ToUnixTimeSeconds();
            return item.start_time <= now && now <= item.end_time;
        }

        // 获取库存描述
        public string GetStockDesc(MallItem item)
        {
            if (item.stock == -1) return "无限";
            if (item.stock <= 0) return "已售罄";
            if (item.stock <= 10) return $"仅剩{item.stock}";
            return $"库存{item.stock}";
        }
    }
}
