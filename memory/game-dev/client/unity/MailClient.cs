// ============================================
// Unity 邮件客户端 - 弹幕游戏
// 对接后端邮件系统 WebSocket + Protobuf 协议
// ============================================

using System;
using System.Collections.Generic;
using UnityEngine;
using UnityEngine.UI;

namespace DanmakuGame.Client
{
    // ============================================
    // 邮件消息ID
    // ============================================
    public enum MailMsgID : ushort
    {
        CS_MAIL_LIST = 6001,
        SC_MAIL_LIST = 6002,
        CS_READ_MAIL = 6003,
        SC_READ_MAIL = 6004,
        CS_CLAIM_ATTACHMENTS = 6005,
        SC_CLAIM_ATTACHMENTS = 6006,
        CS_DELETE_MAIL = 6007,
        SC_DELETE_MAIL = 6008,
        CS_BATCH_DELETE_READ = 6009,
        SC_BATCH_DELETE_READ = 6010,
        CS_MAIL_UNREAD_COUNT = 6011,
        SC_MAIL_UNREAD_COUNT = 6012,
    }

    // ============================================
    // 邮件类型
    // ============================================
    public enum MailType : byte
    {
        System = 1,      // 系统邮件
        Player = 2,      // 玩家邮件
        Gift = 3,       // 礼物邮件
        Auction = 4,    // 拍卖邮件
        GM = 5,         // GM邮件
        Activity = 6,   // 活动邮件
    }

    // ============================================
    // 邮件数据模型
    // ============================================
    [Serializable]
    public class MailInfo
    {
        public string mail_id;
        public string title;
        public string content;
        public string sender_name;
        public MailType mail_type;
        public long send_time;
        public long expire_time;
        public bool is_read;
        public bool has_attachment;
        public List<MailAttachment> attachments;
    }

    [Serializable]
    public class MailAttachment
    {
        public string item_id;
        public string item_name;
        public int item_type;  // 1=金币, 2=钻石, 3=道具, 4=皮肤
        public int count;
        public bool is_claimed;
    }

    // ============================================
    // 邮件列表请求/响应
    // ============================================
    [Serializable]
    public class CSMailListReq
    {
        public string player_id;
        public int page = 1;
        public int page_size = 20;
    }

    [Serializable]
    public class SCMailListResp
    {
        public bool success;
        public string error_msg;
        public List<MailInfo> mails;
        public int total_count;
        public int unread_count;
    }

    // ============================================
    // 读取邮件请求/响应
    // ============================================
    [Serializable]
    public class CSReadMailReq
    {
        public string player_id;
        public string mail_id;
    }

    [Serializable]
    public class SCReadMailResp
    {
        public bool success;
        public string error_msg;
        public MailInfo mail;
    }

    // ============================================
    // 领取附件请求/响应
    // ============================================
    [Serializable]
    public class CSClaimAttachmentsReq
    {
        public string player_id;
        public string mail_id;
    }

    [Serializable]
    public class SCClaimAttachmentsResp
    {
        public bool success;
        public string error_msg;
        public List<MailAttachment> claimed_items;
        public int gold_change;    // 金币变化
        public int gem_change;     // 钻石变化
    }

    // ============================================
    // 删除邮件请求/响应
    // ============================================
    [Serializable]
    public class CSDeleteMailReq
    {
        public string player_id;
        public string mail_id;
    }

    [Serializable]
    public class SCDeleteMailResp
    {
        public bool success;
        public string error_msg;
    }

    // ============================================
    // 批量删除已读邮件请求/响应
    // ============================================
    [Serializable]
    public class CSBatchDeleteReadMailsReq
    {
        public string player_id;
    }

    [Serializable]
    public class SCBatchDeleteReadMailsResp
    {
        public bool success;
        public string error_msg;
        public int deleted_count;
    }

    // ============================================
    // 未读数量请求/响应
    // ============================================
    [Serializable]
    public class CSUnreadCountReq
    {
        public string player_id;
    }

    [Serializable]
    public class SCUnreadCountResp
    {
        public int unread_count;
    }

    // ============================================
    // 邮件客户端管理器
    // ============================================
    public class MailClient : MonoBehaviour
    {
        private static MailClient _instance;
        public static MailClient Instance
        {
            get
            {
                if (_instance == null)
                {
                    GameObject obj = new GameObject("MailClient");
                    _instance = obj.AddComponent<MailClient>();
                    DontDestroyOnLoad(obj);
                }
                return _instance;
            }
        }

        // 邮件列表缓存
        private List<MailInfo> _cachedMails = new List<MailInfo>();
        public List<MailInfo> CachedMails => _cachedMails;

        // 未读数量
        private int _unreadCount = 0;
        public int UnreadCount => _unreadCount;

        // 事件回调
        public Action<List<MailInfo>> OnMailListReceived;
        public Action<MailInfo> OnMailRead;
        public Action<List<MailAttachment>, int, int> OnAttachmentsClaimed;
        public Action<int> OnUnreadCountChanged;
        public Action<bool, string> OnOperationFailed;

        // 关联的主客户端
        private DanmakuGameClient _mainClient;

        public void Initialize(DanmakuGameClient mainClient)
        {
            _mainClient = mainClient;
        }

        // ============================================
        // 获取邮件列表
        // ============================================
        public void GetMailList(string playerId, int page = 1, int pageSize = 20)
        {
            CSMailListReq req = new CSMailListReq
            {
                player_id = playerId,
                page = page,
                page_size = pageSize
            };

            _mainClient.SendMessage((ushort)MailMsgID.CS_MAIL_LIST, req);
        }

        // 处理邮件列表响应
        public void HandleMailListResp(SCMailListResp resp)
        {
            if (resp.success)
            {
                _cachedMails = resp.mails;
                _unreadCount = resp.unread_count;
                OnMailListReceived?.Invoke(resp.mails);
                OnUnreadCountChanged?.Invoke(resp.unread_count);
            }
            else
            {
                OnOperationFailed?.Invoke(false, resp.error_msg);
            }
        }

        // ============================================
        // 读取邮件
        // ============================================
        public void ReadMail(string playerId, string mailId)
        {
            CSReadMailReq req = new CSReadMailReq
            {
                player_id = playerId,
                mail_id = mailId
            };

            _mainClient.SendMessage((ushort)MailMsgID.CS_READ_MAIL, req);
        }

        // 处理读取邮件响应
        public void HandleReadMailResp(SCReadMailResp resp)
        {
            if (resp.success)
            {
                // 更新缓存中的邮件状态
                for (int i = 0; i < _cachedMails.Count; i++)
                {
                    if (_cachedMails[i].mail_id == resp.mail.mail_id)
                    {
                        _cachedMails[i] = resp.mail;
                        break;
                    }
                }

                // 更新未读数量
                if (!_cachedMails.Find(m => m.mail_id == resp.mail.mail_id).is_read)
                {
                    _unreadCount = Mathf.Max(0, _unreadCount - 1);
                    OnUnreadCountChanged?.Invoke(_unreadCount);
                }

                OnMailRead?.Invoke(resp.mail);
            }
            else
            {
                OnOperationFailed?.Invoke(false, resp.error_msg);
            }
        }

        // ============================================
        // 领取附件
        // ============================================
        public void ClaimAttachments(string playerId, string mailId)
        {
            CSClaimAttachmentsReq req = new CSClaimAttachmentsReq
            {
                player_id = playerId,
                mail_id = mailId
            };

            _mainClient.SendMessage((ushort)MailMsgID.CS_CLAIM_ATTACHMENTS, req);
        }

        // 处理领取附件响应
        public void HandleClaimAttachmentsResp(SCClaimAttachmentsResp resp)
        {
            if (resp.success)
            {
                // 更新缓存中的附件状态
                string mailId = "";
                foreach (var mail in _cachedMails)
                {
                    if (mail.attachments != null)
                    {
                        foreach (var att in mail.attachments)
                        {
                            if (!att.is_claimed)
                            {
                                mailId = mail.mail_id;
                                att.is_claimed = true;
                            }
                        }
                    }
                }

                OnAttachmentsClaimed?.Invoke(
                    resp.claimed_items,
                    resp.gold_change,
                    resp.gem_change
                );
            }
            else
            {
                OnOperationFailed?.Invoke(false, resp.error_msg);
            }
        }

        // ============================================
        // 删除邮件
        // ============================================
        public void DeleteMail(string playerId, string mailId)
        {
            CSDeleteMailReq req = new CSDeleteMailReq
            {
                player_id = playerId,
                mail_id = mailId
            };

            _mainClient.SendMessage((ushort)MailMsgID.CS_DELETE_MAIL, req);
        }

        // 处理删除邮件响应
        public void HandleDeleteMailResp(SCDeleteMailResp resp)
        {
            if (resp.success)
            {
                // 从缓存中移除
                _cachedMails.RemoveAll(m => m.mail_id != "");
            }
            else
            {
                OnOperationFailed?.Invoke(false, resp.error_msg);
            }
        }

        // ============================================
        // 批量删除已读邮件
        // ============================================
        public void BatchDeleteReadMails(string playerId)
        {
            CSBatchDeleteReadMailsReq req = new CSBatchDeleteReadMailsReq
            {
                player_id = playerId
            };

            _mainClient.SendMessage((ushort)MailMsgID.CS_BATCH_DELETE_READ, req);
        }

        // 处理批量删除响应
        public void HandleBatchDeleteReadMailsResp(SCBatchDeleteReadMailsResp resp)
        {
            if (resp.success)
            {
                // 重新获取邮件列表
                Debug.Log($"批量删除了 {resp.deleted_count} 封邮件");
            }
            else
            {
                OnOperationFailed?.Invoke(false, resp.error_msg);
            }
        }

        // ============================================
        // 获取未读数量
        // ============================================
        public void GetUnreadCount(string playerId)
        {
            CSUnreadCountReq req = new CSUnreadCountReq
            {
                player_id = playerId
            };

            _mainClient.SendMessage((ushort)MailMsgID.CS_MAIL_UNREAD_COUNT, req);
        }

        // 处理未读数量响应
        public void HandleUnreadCountResp(SCUnreadCountResp resp)
        {
            _unreadCount = resp.unread_count;
            OnUnreadCountChanged?.Invoke(resp.unread_count);
        }

        // ============================================
        // 辅助方法
        // ============================================

        // 获取过期时间描述
        public string GetExpireTimeDesc(long expireTime)
        {
            long remaining = expireTime - GetCurrentTimestamp();
            if (remaining <= 0) return "已过期";

            if (remaining < 3600) return $"{(remaining / 60)}分钟后过期";
            if (remaining < 86400) return $"{(remaining / 3600)}小时后过期";
            return $"{(remaining / 86400)}天后过期";
        }

        // 获取邮件类型描述
        public string GetMailTypeDesc(MailType type)
        {
            switch (type)
            {
                case MailType.System: return "系统";
                case MailType.Player: return "玩家";
                case MailType.Gift: return "礼物";
                case MailType.Auction: return "拍卖";
                case MailType.GM: return "GM";
                case MailType.Activity: return "活动";
                default: return "未知";
            }
        }

        // 获取当前时间戳(秒)
        private long GetCurrentTimestamp()
        {
            return DateTimeOffset.UtcNow.ToUnixTimeSeconds();
        }

        // 检查邮件是否过期
        public bool IsMailExpired(MailInfo mail)
        {
            return GetCurrentTimestamp() >= mail.expire_time;
        }

        // 获取未读邮件
        public List<MailInfo> GetUnreadMails()
        {
            return _cachedMails.FindAll(m => !m.is_read);
        }

        // 获取有附件的邮件
        public List<MailInfo> GetMailsWithAttachments()
        {
            return _cachedMails.FindAll(m => m.has_attachment);
        }
    }
}
