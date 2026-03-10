using System;
using System.Collections.Generic;
using System.Threading.Tasks;

namespace DanmakuGameClient
{
    /// <summary>
    /// 社交系统客户端 - 好友、公会、黑名单
    /// </summary>
    public class SocialClient
    {
        private readonly GameClient _gameClient;
        private List<Friend> _friends = new List<Friend>();
        private List<Guild> _guilds = new List<Guild>();
        private Guild _currentGuild;
        private List<BlacklistEntry> _blacklist = new List<BlacklistEntry>();
        private List<RecommendFriend> _recommendFriends = new List<RecommendFriend>();
        private List<RecommendGuild> _recommendGuilds = new List<RecommendGuild>();

        // 回调
        public Action<List<Friend>> OnFriendsUpdated;
        public Action<Friend> OnFriendAdded;
        public Action<string> OnFriendRemoved;
        public Action<Friend> OnFriendOnline;
        public Action<Friend> OnFriendOffline;
        public Action<FriendInvite> OnFriendInviteReceived;
        public Action<Guild> OnGuildUpdated;
        public Action<Guild> OnGuildCreated;
        public Action OnGuildLeft;
        public Action<List<GuildMember>> OnGuildMembersUpdated;
        public Action<GuildMember> OnGuildMemberJoined;
        public Action<string> OnGuildMemberLeft;
        public Action<GuildApply> OnGuildApplyReceived;
        public Action<GuildInvite> OnGuildInviteReceived;
        public Action<List<BlacklistEntry>> OnBlacklistUpdated;
        public Action<string> OnBlacklistAdded;
        public Action<string> OnBlacklistRemoved;
        public Action<string> OnError;

        public SocialClient(GameClient gameClient)
        {
            _gameClient = gameClient;
            RegisterHandlers();
        }

        #region 消息注册

        private void RegisterHandlers()
        {
            // 好友相关
            _gameClient.RegisterHandler<SCFriendsList>(OnFriendsList);
            _gameClient.RegisterHandler<SCFriendAdded>(OnFriendAddedNotify);
            _gameClient.RegisterHandler<SCFriendRemoved>(OnFriendRemovedNotify);
            _gameClient.RegisterHandler<SCFriendOnline>(OnFriendOnlineNotify);
            _gameClient.RegisterHandler<SCFriendOffline>(OnFriendOfflineNotify);
            _gameClient.RegisterHandler<SCFriendInvite>(OnFriendInviteNotify);
            _gameClient.RegisterHandler<SCFriendInviteResult>(OnFriendInviteResult);
            _gameClient.RegisterHandler<SCFriendStatusChanged>(OnFriendStatusChanged);
            _gameClient.RegisterHandler<SCRecommendFriends>(OnRecommendFriends);

            // 公会相关
            _gameClient.RegisterHandler<SCGuildInfo>(OnGuildInfo);
            _gameClient.RegisterHandler<SCGuildCreated>(OnGuildCreatedNotify);
            _gameClient.RegisterHandler<SCGuildLeft>(OnGuildLeftNotify);
            _gameClient.RegisterHandler<SCGuildMembers>(OnGuildMembers);
            _gameClient.RegisterHandler<SCGuildMemberJoined>(OnGuildMemberJoinedNotify);
            _gameClient.RegisterHandler<SCGuildMemberLeft>(OnGuildMemberLeftNotify);
            _gameClient.RegisterHandler<SCGuildApply>(OnGuildApplyNotify);
            _gameClient.RegisterHandler<SCGuildApplyResult>(OnGuildApplyResult);
            _gameClient.RegisterHandler<SCGuildInvite>(OnGuildInviteNotify);
            _gameClient.RegisterHandler<SCGuildInviteResult>(OnGuildInviteResult);
            _gameClient.RegisterHandler<SCGuildNoticeUpdated>(OnGuildNoticeUpdated);
            _gameClient.RegisterHandler<SCGuildLevelUp>( OnGuildLevelUp);
            _gameClient.RegisterHandler<SCRecommendGuilds>(OnRecommendGuilds);

            // 黑名单相关
            _gameClient.RegisterHandler<SCBlacklist>(OnBlacklist);
            _gameClient.RegisterHandler<SCBlacklistAdded>(OnBlacklistAddedNotify);
            _gameClient.RegisterHandler<SCBlacklistRemoved>(OnBlacklistRemovedNotify);

            // 错误
            _gameClient.RegisterHandler<SCSocialError>(OnSocialError);
        }

        #endregion

        #region 好友系统

        /// <summary>
        /// 获取好友列表
        /// </summary>
        public async Task<List<Friend>> GetFriendsAsync()
        {
            var req = new CSGetFriends();
            var result = await _gameClient.SendRequestAsync<CSGetFriends, SCFriendsList>(req);
            _friends = result?.Friends ?? new List<Friend>();
            return _friends;
        }

        /// <summary>
        /// 添加好友
        /// </summary>
        public async Task<bool> AddFriendAsync(string userId, string message = "")
        {
            var req = new CSAddFriend
            {
                UserId = userId,
                Message = message
            };
            var result = await _gameClient.SendRequestAsync<CSAddFriend, SCAddFriendResult>(req);
            return result?.Success ?? false;
        }

        /// <summary>
        /// 删除好友
        /// </summary>
        public async Task<bool> RemoveFriendAsync(string userId)
        {
            var req = new CSRemoveFriend { UserId = userId };
            var result = await _gameClient.SendRequestAsync<CSRemoveFriend, SCRemoveFriendResult>(req);
            return result?.Success ?? false;
        }

        /// <summary>
        /// 搜索好友
        /// </summary>
        public async Task<List<PlayerBrief>> SearchFriendsAsync(string keyword)
        {
            var req = new CSSearchFriends { Keyword = keyword };
            var result = await _gameClient.SendRequestAsync<CSSearchFriends, SCSearchFriendsResult>(req);
            return result?.Players ?? new List<PlayerBrief>();
        }

        /// <summary>
        /// 发送好友邀请
        /// </summary>
        public async Task<bool> SendFriendInviteAsync(string userId, string message = "")
        {
            var req = new CSFriendInvite
            {
                UserId = userId,
                Message = message
            };
            var result = await _gameClient.SendRequestAsync<CSFriendInvite, SCFriendInviteResult>(req);
            return result?.Success ?? false;
        }

        /// <summary>
        /// 响应好友邀请
        /// </summary>
        public async Task<bool> RespondFriendInviteAsync(string inviteId, bool accept)
        {
            var req = new CSFriendInviteRespond
            {
                InviteId = inviteId,
                Accept = accept
            };
            var result = await _gameClient.SendRequestAsync<CSFriendInviteRespond, SCFriendInviteRespondResult>(req);
            return result?.Success ?? false;
        }

        /// <summary>
        /// 获取推荐好友
        /// </summary>
        public async Task<List<RecommendFriend>> GetRecommendFriendsAsync()
        {
            var req = new CSGetRecommendFriends();
            var result = await _gameClient.SendRequestAsync<CSGetRecommendFriends, SCRecommendFriends>(req);
            _recommendFriends = result?.Friends ?? new List<RecommendFriend>();
            return _recommendFriends;
        }

        /// <summary>
        /// 发送好友消息
        /// </summary>
        public async Task<bool> SendFriendMessageAsync(string friendId, string message)
        {
            var req = new CSFriendMessage
            {
                FriendId = friendId,
                Message = message
            };
            var result = await _gameClient.SendRequestAsync<CSFriendMessage, SCFriendMessageResult>(req);
            return result?.Success ?? false;
        }

        /// <summary>
        /// 查看好友详情
        /// </summary>
        public async Task<FriendDetail> GetFriendDetailAsync(string friendId)
        {
            var req = new CSGetFriendDetail { FriendId = friendId };
            var result = await _gameClient.SendRequestAsync<CSGetFriendDetail, SCGetFriendDetailResult>(req);
            return result?.Detail;
        }

        #endregion

        #region 公会系统

        /// <summary>
        /// 创建公会
        /// </summary>
        public async Task<Guild> CreateGuildAsync(string name, string icon, string notice)
        {
            var req = new CSCreateGuild
            {
                Name = name,
                Icon = icon,
                Notice = notice
            };
            var result = await _gameClient.SendRequestAsync<CSCreateGuild, SCGuildCreated>(req);
            if (result?.Guild != null)
            {
                _currentGuild = result.Guild;
            }
            return result?.Guild;
        }

        /// <summary>
        /// 加入公会
        /// </summary>
        public async Task<bool> JoinGuildAsync(string guildId, string message = "")
        {
            var req = new CSJoinGuild
            {
                GuildId = guildId,
                Message = message
            };
            var result = await _gameClient.SendRequestAsync<CSJoinGuild, SCJoinGuildResult>(req);
            return result?.Success ?? false;
        }

        /// <summary>
        /// 离开公会
        /// </summary>
        public async Task LeaveGuildAsync()
        {
            var req = new CSLeaveGuild();
            await _gameClient.SendRequestAsync<CSLeaveGuild, SCLeaveGuildResult>(req);
            _currentGuild = null;
        }

        /// <summary>
        /// 获取公会信息
        /// </summary>
        public async Task<Guild> GetGuildInfoAsync(string guildId)
        {
            var req = new CSGetGuildInfo { GuildId = guildId };
            var result = await _gameClient.SendRequestAsync<CSGetGuildInfo, SCGuildInfo>(req);
            return result?.Guild;
        }

        /// <summary>
        /// 获取我的公会信息
        /// </summary>
        public async Task<Guild> GetMyGuildAsync()
        {
            var req = new CSGetMyGuild();
            var result = await _gameClient.SendRequestAsync<CSGetMyGuild, SCGuildInfo>(req);
            _currentGuild = result?.Guild;
            return _currentGuild;
        }

        /// <summary>
        /// 获取公会成员列表
        /// </summary>
        public async Task<List<GuildMember>> GetGuildMembersAsync(string guildId)
        {
            var req = new CSGetGuildMembers { GuildId = guildId };
            var result = await _gameClient.SendRequestAsync<CSGetGuildMembers, SCGuildMembers>(req);
            return result?.Members ?? new List<GuildMember>();
        }

        /// <summary>
        /// 搜索公会
        /// </summary>
        public async Task<List<Guild>> SearchGuildsAsync(string keyword)
        {
            var req = new CSSearchGuilds { Keyword = keyword };
            var result = await _gameClient.SendRequestAsync<CSSearchGuilds, SCSearchGuildsResult>(req);
            return result?.Guilds ?? new List<Guild>();
        }

        /// <summary>
        /// 申请加入公会
        /// </summary>
        public async Task<bool> ApplyGuildAsync(string guildId, string message = "")
        {
            var req = new CSGuildApply
            {
                GuildId = guildId,
                Message = message
            };
            var result = await _gameClient.SendRequestAsync<CSGuildApply, SCGuildApplyResult>(req);
            return result?.Success ?? false;
        }

        /// <summary>
        /// 响应公会申请
        /// </summary>
        public async Task<bool> RespondGuildApplyAsync(string applyId, bool accept)
        {
            var req = new CSGuildApplyRespond
            {
                ApplyId = applyId,
                Accept = accept
            };
            var result = await _gameClient.SendRequestAsync<CSGuildApplyRespond, SCGuildApplyRespondResult>(req);
            return result?.Success ?? false;
        }

        /// <summary>
        /// 邀请玩家加入公会
        /// </summary>
        public async Task<bool> InviteGuildMemberAsync(string userId)
        {
            var req = new CSGuildInvite
            {
                UserId = userId
            };
            var result = await _gameClient.SendRequestAsync<CSGuildInvite, SCGuildInviteResult>(req);
            return result?.Success ?? false;
        }

        /// <summary>
        /// 响应公会邀请
        /// </summary>
        public async Task<bool> RespondGuildInviteAsync(string inviteId, bool accept)
        {
            var req = new CSGuildInviteRespond
            {
                InviteId = inviteId,
                Accept = accept
            };
            var result = await _gameClient.SendRequestAsync<CSGuildInviteRespond, SCGuildInviteRespondResult>(req);
            return result?.Success ?? false;
        }

        /// <summary>
        /// 修改公会公告
        /// </summary>
        public async Task<bool> UpdateGuildNoticeAsync(string notice)
        {
            var req = new CSUpdateGuildNotice { Notice = notice };
            var result = await _gameClient.SendRequestAsync<CSUpdateGuildNotice, SCUpdateGuildNoticeResult>(req);
            return result?.Success ?? false;
        }

        /// <summary>
        /// 捐赠公会
        /// </summary>
        public async Task<GuildContributeResult> ContributeGuildAsync(int itemId, int count)
        {
            var req = new CSGuildContribute
            {
                ItemId = itemId,
                Count = count
            };
            var result = await _gameClient.SendRequestAsync<CSGuildContribute, SCGuildContributeResult>(req);
            return result;
        }

        /// <summary>
        /// 升级公会
        /// </summary>
        public async Task<bool> UpgradeGuildAsync()
        {
            var req = new CSUpgradeGuild();
            var result = await _gameClient.SendRequestAsync<CSUpgradeGuild, SCUpgradeGuildResult>(req);
            return result?.Success ?? false;
        }

        /// <summary>
        /// 任命公会职位
        /// </summary>
        public async Task<bool> SetGuildMemberRoleAsync(string memberId, int role)
        {
            var req = new CSSetGuildMemberRole
            {
                MemberId = memberId,
                Role = role
            };
            var result = await _gameClient.SendRequestAsync<CSSetGuildMemberRole, SCSetGuildMemberRoleResult>(req);
            return result?.Success ?? false;
        }

        /// <summary>
        /// 踢出公会成员
        /// </summary>
        public async Task<bool> KickGuildMemberAsync(string memberId)
        {
            var req = new CSKickGuildMember { MemberId = memberId };
            var result = await _gameClient.SendRequestAsync<CSKickGuildMember, SCKickGuildMemberResult>(req);
            return result?.Success ?? false;
        }

        /// <summary>
        /// 转让会长
        /// </summary>
        public async Task<bool> TransferGuildLeaderAsync(string memberId)
        {
            var req = new CSTransferGuildLeader { MemberId = memberId };
            var result = await _gameClient.SendRequestAsync<CSTransferGuildLeader, SCTransferGuildLeaderResult>(req);
            return result?.Success ?? false;
        }

        /// <summary>
        /// 获取公会日志
        /// </summary>
        public async Task<List<GuildLog>> GetGuildLogsAsync(string guildId)
        {
            var req = new CSGetGuildLogs { GuildId = guildId };
            var result = await _gameClient.SendRequestAsync<CSGetGuildLogs, SCGuildLogs>(req);
            return result?.Logs ?? new List<GuildLog>();
        }

        /// <summary>
        /// 获取公会排行榜
        /// </summary>
        public async Task<List<Guild>> GetGuildRankListAsync(int page, int pageSize)
        {
            var req = new CSGetGuildRankList
            {
                Page = page,
                PageSize = pageSize
            };
            var result = await _gameClient.SendRequestAsync<CSGetGuildRankList, SCGuildRankList>(req);
            return result?.Guilds ?? new List<Guild>();
        }

        /// <summary>
        /// 获取推荐公会
        /// </summary>
        public async Task<List<RecommendGuild>> GetRecommendGuildsAsync()
        {
            var req = new CSGetRecommendGuilds();
            var result = await _gameClient.SendRequestAsync<CSGetRecommendGuilds, SCRecommendGuilds>(req);
            _recommendGuilds = result?.Guilds ?? new List<RecommendGuild>();
            return _recommendGuilds;
        }

        #endregion

        #region 黑名单系统

        /// <summary>
        /// 获取黑名单
        /// </summary>
        public async Task<List<BlacklistEntry>> GetBlacklistAsync()
        {
            var req = new CSGetBlacklist();
            var result = await _gameClient.SendRequestAsync<CSGetBlacklist, SCBlacklist>(req);
            _blacklist = result?.Entries ?? new List<BlacklistEntry>();
            return _blacklist;
        }

        /// <summary>
        /// 添加黑名单
        /// </summary>
        public async Task<bool> AddBlacklistAsync(string userId)
        {
            var req = new CSAddBlacklist { UserId = userId };
            var result = await _gameClient.SendRequestAsync<CSAddBlacklist, SCAddBlacklistResult>(req);
            return result?.Success ?? false;
        }

        /// <summary>
        /// 移除黑名单
        /// </summary>
        public async Task<bool> RemoveBlacklistAsync(string userId)
        {
            var req = new CSRemoveBlacklist { UserId = userId };
            var result = await _gameClient.SendRequestAsync<CSRemoveBlacklist, SCRemoveBlacklistResult>(req);
            return result?.Success ?? false;
        }

        #endregion

        #region 消息处理

        private void OnFriendsList(SCFriendsList msg)
        {
            _friends = msg.Friends;
            OnFriendsUpdated?.Invoke(_friends);
        }

        private void OnFriendAddedNotify(SCFriendAdded msg)
        {
            _friends.Add(msg.Friend);
            OnFriendAdded?.Invoke(msg.Friend);
        }

        private void OnFriendRemovedNotify(SCFriendRemoved msg)
        {
            _friends.RemoveAll(f => f.UserId == msg.UserId);
            OnFriendRemoved?.Invoke(msg.UserId);
        }

        private void OnFriendOnlineNotify(SCFriendOnline msg)
        {
            var friend = _friends.Find(f => f.UserId == msg.UserId);
            if (friend != null)
            {
                friend.Online = true;
                friend.LastOnlineTime = msg.Timestamp;
                OnFriendOnline?.Invoke(friend);
            }
        }

        private void OnFriendOfflineNotify(SCFriendOffline msg)
        {
            var friend = _friends.Find(f => f.UserId == msg.UserId);
            if (friend != null)
            {
                friend.Online = false;
                friend.LastOnlineTime = msg.Timestamp;
                OnFriendOffline?.Invoke(friend);
            }
        }

        private void OnFriendInviteNotify(SCFriendInvite msg)
        {
            OnFriendInviteReceived?.Invoke(msg.Invite);
        }

        private void OnFriendInviteResult(SCFriendInviteResult msg)
        {
            // 邀请结果
        }

        private void OnFriendStatusChanged(SCFriendStatusChanged msg)
        {
            var friend = _friends.Find(f => f.UserId == msg.UserId);
            if (friend != null)
            {
                friend.Status = msg.Status;
            }
        }

        private void OnRecommendFriends(SCRecommendFriends msg)
        {
            _recommendFriends = msg.Friends;
        }

        private void OnGuildInfo(SCGuildInfo msg)
        {
            _currentGuild = msg.Guild;
            OnGuildUpdated?.Invoke(msg.Guild);
        }

        private void OnGuildCreatedNotify(SCGuildCreated msg)
        {
            _currentGuild = msg.Guild;
            OnGuildCreated?.Invoke(msg.Guild);
        }

        private void OnGuildLeftNotify(SCGuildLeft msg)
        {
            _currentGuild = null;
            OnGuildLeft?.Invoke();
        }

        private void OnGuildMembers(SCGuildMembers msg)
        {
            OnGuildMembersUpdated?.Invoke(msg.Members);
        }

        private void OnGuildMemberJoinedNotify(SCGuildMemberJoined msg)
        {
            OnGuildMemberJoined?.Invoke(msg.Member);
        }

        private void OnGuildMemberLeftNotify(SCGuildMemberLeft msg)
        {
            OnGuildMemberLeft?.Invoke(msg.UserId);
        }

        private void OnGuildApplyNotify(SCGuildApply msg)
        {
            OnGuildApplyReceived?.Invoke(msg.Apply);
        }

        private void OnGuildApplyResult(SCGuildApplyResult msg)
        {
            // 申请结果
        }

        private void OnGuildInviteNotify(SCGuildInvite msg)
        {
            OnGuildInviteReceived?.Invoke(msg.Invite);
        }

        private void OnGuildInviteResult(SCGuildInviteResult msg)
        {
            // 邀请结果
        }

        private void OnGuildNoticeUpdated(SCGuildNoticeUpdated msg)
        {
            if (_currentGuild != null)
            {
                _currentGuild.Notice = msg.Notice;
            }
        }

        private void OnGuildLevelUp(SCGuildLevelUp msg)
        {
            if (_currentGuild != null)
            {
                _currentGuild.Level = msg.NewLevel;
            }
        }

        private void OnRecommendGuilds(SCRecommendGuilds msg)
        {
            _recommendGuilds = msg.Guilds;
        }

        private void OnBlacklist(SCBlacklist msg)
        {
            _blacklist = msg.Entries;
            OnBlacklistUpdated?.Invoke(_blacklist);
        }

        private void OnBlacklistAddedNotify(SCBlacklistAdded msg)
        {
            _blacklist.Add(new BlacklistEntry { UserId = msg.UserId, BlockTime = msg.Timestamp });
            OnBlacklistAdded?.Invoke(msg.UserId);
        }

        private void OnBlacklistRemovedNotify(SCBlacklistRemoved msg)
        {
            _blacklist.RemoveAll(b => b.UserId == msg.UserId);
            OnBlacklistRemoved?.Invoke(msg.UserId);
        }

        private void OnSocialError(SCSocialError msg)
        {
            OnError?.Invoke(msg.ErrorMessage);
        }

        #endregion

        #region 属性

        public List<Friend> Friends => _friends;
        public List<Guild> Guilds => _guilds;
        public Guild CurrentGuild => _currentGuild;
        public List<BlacklistEntry> Blacklist => _blacklist;
        public List<RecommendFriend> RecommendFriends => _recommendFriends;
        public List<RecommendGuild> RecommendGuilds => _recommendGuilds;

        #endregion
    }

    #region 消息定义 - 好友

    public class CSGetFriends { }
    public class CSAddFriend { public string UserId; public string Message; }
    public class CSRemoveFriend { public string UserId; }
    public class CSSearchFriends { public string Keyword; }
    public class CSFriendInvite { public string UserId; public string Message; }
    public class CSFriendInviteRespond { public string InviteId; public bool Accept; }
    public class CSGetRecommendFriends { }
    public class CSFriendMessage { public string FriendId; public string Message; }
    public class CSGetFriendDetail { public string FriendId; }

    public class SCFriendsList { public List<Friend> Friends; }
    public class SCAddFriendResult { public bool Success; }
    public class SCRemoveFriendResult { public bool Success; }
    public class SCSearchFriendsResult { public List<PlayerBrief> Players; }
    public class SCFriendInviteResult { public bool Success; }
    public class SCFriendInviteRespondResult { public bool Success; }
    public class SCRecommendFriends { public List<RecommendFriend> Friends; }
    public class SCFriendMessageResult { public bool Success; }
    public class SCGetFriendDetailResult { public FriendDetail Detail; }

    public class SCFriendAdded { public Friend Friend; }
    public class SCFriendRemoved { public string UserId; }
    public class SCFriendOnline { public string UserId; public long Timestamp; }
    public class SCFriendOffline { public string UserId; public long Timestamp; }
    public class SCFriendInvite { public FriendInvite Invite; }
    public class SCFriendStatusChanged { public string UserId; public int Status; }

    #endregion

    #region 消息定义 - 公会

    public class CSCreateGuild { public string Name; public string Icon; public string Notice; }
    public class CSJoinGuild { public string GuildId; public string Message; }
    public class CSLeaveGuild { }
    public class CSGetGuildInfo { public string GuildId; }
    public class CSGetMyGuild { }
    public class CSGetGuildMembers { public string GuildId; }
    public class CSSearchGuilds { public string Keyword; }
    public class CSGuildApply { public string GuildId; public string Message; }
    public class CSGuildApplyRespond { public string ApplyId; public bool Accept; }
    public class CSGuildInvite { public string UserId; }
    public class CSGuildInviteRespond { public string InviteId; public bool Accept; }
    public class CSUpdateGuildNotice { public string Notice; }
    public class CSGuildContribute { public int ItemId; public int Count; }
    public class CSUpgradeGuild { }
    public class CSSetGuildMemberRole { public string MemberId; public int Role; }
    public class CSKickGuildMember { public string MemberId; }
    public class CSTransferGuildLeader { public string MemberId; }
    public class CSGetGuildLogs { public string GuildId; }
    public class CSGetGuildRankList { public int Page; public int PageSize; }
    public class CSGetRecommendGuilds { }

    public class SCGuildCreated { public Guild Guild; }
    public class SCJoinGuildResult { public bool Success; }
    public class SCLeaveGuildResult { }
    public class SCGuildInfo { public Guild Guild; }
    public class SCGuildMembers { public List<GuildMember> Members; }
    public class SCSearchGuildsResult { public List<Guild> Guilds; }
    public class SCGuildApplyResult { public bool Success; }
    public class SCGuildApplyRespondResult { public bool Success; }
    public class SCGuildInviteResult { public bool Success; }
    public class SCGuildInviteRespondResult { public bool Success; }
    public class SCUpdateGuildNoticeResult { public bool Success; }
    public class GuildContributeResult { public int NewContribution; public int NewExp; public int NewLevel; }
    public class SCUpgradeGuildResult { public bool Success; }
    public class SCSetGuildMemberRoleResult { public bool Success; }
    public class SCKickGuildMemberResult { public bool Success; }
    public class SCTransferGuildLeaderResult { public bool Success; }
    public class SCGuildLogs { public List<GuildLog> Logs; }
    public class SCGuildRankList { public List<Guild> Guilds; }
    public class SCRecommendGuilds { public List<RecommendGuild> Guilds; }

    public class SCGuildLeft { }
    public class SCGuildMemberJoined { public GuildMember Member; }
    public class SCGuildMemberLeft { public string UserId; }
    public class SCGuildApply { public GuildApply Apply; }
    public class SCGuildInvite { public GuildInvite Invite; }
    public class SCGuildNoticeUpdated { public string Notice; }
    public class SCGuildLevelUp { public int NewLevel; }

    #endregion

    #region 消息定义 - 黑名单

    public class CSGetBlacklist { }
    public class CSAddBlacklist { public string UserId; }
    public class CSRemoveBlacklist { public string UserId; }

    public class SCBlacklist { public List<BlacklistEntry> Entries; }
    public class SCAddBlacklistResult { public bool Success; }
    public class SCRemoveBlacklistResult { public bool Success; }
    public class SCBlacklistAdded { public string UserId; public long Timestamp; }
    public class SCBlacklistRemoved { public string UserId; }

    #endregion

    public class SCSocialError { public string ErrorMessage; }

    #region 数据结构

    public class Friend
    {
        public string UserId { get; set; }
        public string UserName { get; set; }
        public string Avatar { get; set; }
        public int Level { get; set; }
        public bool Online { get; set; }
        public int Status { get; set; }
        public long LastOnlineTime { get; set; }
        public long AddTime { get; set; }
        public string Remark { get; set; }
    }

    public class FriendDetail
    {
        public string UserId { get; set; }
        public string UserName { get; set; }
        public string Avatar { get; set; }
        public int Level { get; set; }
        public int VipLevel { get; set; }
        public bool Online { get; set; }
        public string Signature { get; set; }
        public int TotalGames { get; set; }
        public int WinGames { get; set; }
        public int MaxCombo { get; set; }
        public long RegisterTime { get; set; }
        public long LastOnlineTime { get; set; }
    }

    public class PlayerBrief
    {
        public string UserId { get; set; }
        public string UserName { get; set; }
        public string Avatar { get; set; }
        public int Level { get; set; }
        public bool Online { get; set; }
    }

    public class RecommendFriend
    {
        public string UserId { get; set; }
        public string UserName { get; set; }
        public string Avatar { get; set; }
        public int Level { get; set; }
        public string Reason { get; set; }
    }

    public class FriendInvite
    {
        public string InviteId { get; set; }
        public string FromUserId { get; set; }
        public string FromUserName { get; set; }
        public string ToUserId { get; set; }
        public string Message { get; set; }
        public long CreateTime { get; set; }
        public int Status { get; set; }
    }

    public class Guild
    {
        public string GuildId { get; set; }
        public string Name { get; set; }
        public string Icon { get; set; }
        public string Notice { get; set; }
        public int Level { get; set; }
        public int MemberCount { get; set; }
        public int MaxMembers { get; set; }
        public string LeaderId { get; set; }
        public string LeaderName { get; set; }
        public long CreateTime { get; set; }
        public long TotalContribution { get; set; }
        public int Rank { get; set; }
        public long Exp { get; set; }
        public long NextLevelExp { get; set; }
    }

    public class GuildMember
    {
        public string UserId { get; set; }
        public string UserName { get; set; }
        public string Avatar { get; set; }
        public int Level { get; set; }
        public int Role { get; set; }
        public string RoleName { get; set; }
        public long Contribution { get; set; }
        public long WeekContribution { get; set; }
        public long JoinTime { get; set; }
        public long LastActiveTime { get; set; }
        public bool Online { get; set; }
    }

    public class GuildApply
    {
        public string ApplyId { get; set; }
        public string UserId { get; set; }
        public string UserName { get; set; }
        public string Avatar { get; set; }
        public int Level { get; set; }
        public string GuildId { get; set; }
        public string Message { get; set; }
        public long ApplyTime { get; set; }
        public int Status { get; set; }
    }

    public class GuildInvite
    {
        public string InviteId { get; set; }
        public string GuildId { get; set; }
        public string GuildName { get; set; }
        public string FromUserId { get; set; }
        public string FromUserName { get; set; }
        public string ToUserId { get; set; }
        public long CreateTime { get; set; }
        public int Status { get; set; }
    }

    public class RecommendGuild
    {
        public string GuildId { get; set; }
        public string Name { get; set; }
        public string Icon { get; set; }
        public int Level { get; set; }
        public int MemberCount { get; set; }
        public int MaxMembers { get; set; }
        public string Reason { get; set; }
    }

    public class GuildLog
    {
        public long LogId { get; set; }
        public int LogType { get; set; }
        public string UserId { get; set; }
        public string UserName { get; set; }
        public string Content { get; set; }
        public long CreateTime { get; set; }
    }

    public class BlacklistEntry
    {
        public string UserId { get; set; }
        public string UserName { get; set; }
        public string Avatar { get; set; }
        public long BlockTime { get; set; }
    }

    #endregion
}
