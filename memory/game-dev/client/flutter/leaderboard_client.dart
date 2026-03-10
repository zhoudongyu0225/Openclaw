// ============================================================
// 排行榜客户端 - Leaderboard Client
// 弹幕游戏 Flutter 客户端
// ============================================================

import 'dart:convert';
import 'package:flutter/material.dart';
import 'package:dio/dio.dart';

// ============================================================
// 枚举定义
// ============================================================

enum LeaderboardType {
  level,     // 等级榜
  gold,     // 金币榜
  gem,       // 钻石榜
  combat,    // 战斗力榜
  kill,      // 击杀榜
  damage,    // 伤害榜
  survival,  // 生存榜
  win,       // 胜场榜
  guild,     // 公会榜
  signIn,    // 签到榜
  rich,      // 富豪榜
  mvp,       // MVP榜
}

enum LeaderboardPeriod {
  all,     // 全部
  daily,   // 今日
  weekly,  // 本周
  monthly, // 本月
  season,  // 本赛季
}

// ============================================================
// 数据模型
// ============================================================

class LeaderboardEntry {
  final int rank;
  final String playerId;
  final String playerName;
  final int level;
  final int value;
  final String? avatar;
  final int? titleId;
  final String? guildName;
  final int vipLevel;
  final bool isMyFriend;
  final int change;

  LeaderboardEntry({
    required this.rank,
    required this.playerId,
    required this.playerName,
    required this.level,
    required this.value,
    this.avatar,
    this.titleId,
    this.guildName,
    this.vipLevel = 0,
    this.isMyFriend = false,
    this.change = 0,
  });

  factory LeaderboardEntry.fromJson(Map<String, dynamic> json) {
    return LeaderboardEntry(
      rank: json['rank'] ?? 0,
      playerId: json['player_id'] ?? '',
      playerName: json['player_name'] ?? '',
      level: json['level'] ?? 0,
      value: json['value'] ?? 0,
      avatar: json['avatar'],
      titleId: json['title_id'],
      guildName: json['guild_name'],
      vipLevel: json['vip_level'] ?? 0,
      isMyFriend: json['is_my_friend'] ?? false,
      change: json['change'] ?? 0,
    );
  }

  Map<String, dynamic> toJson() => {
    'rank': rank,
    'player_id': playerId,
    'player_name': playerName,
    'level': level,
    'value': value,
    'avatar': avatar,
    'title_id': titleId,
    'guild_name': guildName,
    'vip_level': vipLevel,
    'is_my_friend': isMyFriend,
    'change': change,
  };
}

class LeaderboardData {
  final LeaderboardType type;
  final LeaderboardPeriod period;
  final List<LeaderboardEntry> entries;
  final int myRank;
  final int myValue;
  final int totalCount;
  final int refreshTime;

  LeaderboardData({
    required this.type,
    required this.period,
    required this.entries,
    this.myRank = -1,
    this.myValue = 0,
    this.totalCount = 0,
    this.refreshTime = 0,
  });

  factory LeaderboardData.fromJson(Map<String, dynamic> json, LeaderboardType type, LeaderboardPeriod period) {
    return LeaderboardData(
      type: type,
      period: period,
      entries: (json['entries'] as List<dynamic>?)
          ?.map((e) => LeaderboardEntry.fromJson(e))
          .toList() ?? [],
      myRank: json['my_rank'] ?? -1,
      myValue: json['my_value'] ?? 0,
      totalCount: json['total_count'] ?? 0,
      refreshTime: json['refresh_time'] ?? 0,
    );
  }
}

// ============================================================
// 排行榜客户端
// ============================================================

class LeaderboardClient {
  static final LeaderboardClient _instance = LeaderboardClient._internal();
  factory LeaderboardClient() => _instance;
  LeaderboardClient._internal();

  // 网络客户端
  late Dio _dio;
  
  // 缓存
  final Map<LeaderboardType, LeaderboardData> _cache = {};
  final Map<LeaderboardType, DateTime> _lastRefreshTime = {};
  
  // 配置
  static const int _cacheDurationSeconds = 60;
  static const int _pageSize = 50;
  
  // 我的排名信息
  final Map<LeaderboardType, int> _myRanks = {};
  final Map<LeaderboardType, int> _myValues = {};
  final Set<LeaderboardType> _subscribedTypes = {};
  
  // 事件回调
  Function(LeaderboardType, LeaderboardData)? onLeaderboardReceived;
  Function(LeaderboardType, int)? onMyRankReceived;
  Function(String)? onError;

  // ============================================================
  // 初始化
  // ============================================================
  
  void init(Dio dio) {
    _dio = dio;
  }

  // ============================================================
  // 获取排行榜 (缓存版本)
  // ============================================================
  
  LeaderboardData? getLeaderboard(LeaderboardType type, [LeaderboardPeriod period = LeaderboardPeriod.all]) {
    final cached = _cache[type];
    if (cached != null && cached.period == period) {
      final lastTime = _lastRefreshTime[type];
      if (lastTime != null) {
        final diff = DateTime.now().difference(lastTime).inSeconds;
        if (diff < _cacheDurationSeconds) {
          return cached;
        }
      }
    }
    return null;
  }

  // ============================================================
  // 请求排行榜
  // ============================================================
  
  Future<void> requestLeaderboard(
    LeaderboardType type, [
    LeaderboardPeriod period = LeaderboardPeriod.all,
    int page = 1,
  ]) async {
    try {
      final response = await _dio.get(
        '/api/leaderboard',
        queryParameters: {
          'type': type.index,
          'period': period.index,
          'page': page,
          'page_size': _pageSize,
        },
      );
      
      if (response.statusCode == 200) {
        final data = response.data;
        final leaderboardData = LeaderboardData.fromJson(data, type, period);
        _cache[type] = leaderboardData;
        _lastRefreshTime[type] = DateTime.now();
        
        // 更新我的排名
        _myRanks[type] = leaderboardData.myRank;
        _myValues[type] = leaderboardData.myValue;
        
        onLeaderboardReceived?.call(type, leaderboardData);
      }
    } catch (e) {
      onError?.call(e.toString());
    }
  }

  // ============================================================
  // 请求我的排名
  // ============================================================
  
  Future<void> requestMyRank(LeaderboardType type) async {
    try {
      final response = await _dio.get(
        '/api/leaderboard/my_rank',
        queryParameters: {'type': type.index},
      );
      
      if (response.statusCode == 200) {
        final data = response.data;
        final rank = data['rank'] ?? -1;
        final value = data['value'] ?? 0;
        
        _myRanks[type] = rank;
        _myValues[type] = value;
        
        onMyRankReceived?.call(type, rank);
      }
    } catch (e) {
      onError?.call(e.toString());
    }
  }

  // ============================================================
  // 请求多个排行榜
  // ============================================================
  
  Future<void> requestMultipleLeaderboards(
    List<LeaderboardType> types, [
    LeaderboardPeriod period = LeaderboardPeriod.all,
  ]) async {
    await Future.wait(types.map((t) => requestLeaderboard(t, period)));
  }

  // ============================================================
  // 订阅排行榜
  // ============================================================
  
  Future<void> subscribeLeaderboard(LeaderboardType type) async {
    _subscribedTypes.add(type);
    try {
      await _dio.post(
        '/api/leaderboard/subscribe',
        data: {'type': type.index, 'subscribe': true},
      );
    } catch (e) {
      onError?.call(e.toString());
    }
  }

  // ============================================================
  // 取消订阅
  // ============================================================
  
  Future<void> unsubscribeLeaderboard(LeaderboardType type) async {
    _subscribedTypes.remove(type);
    try {
      await _dio.post(
        '/api/leaderboard/subscribe',
        data: {'type': type.index, 'subscribe': false},
      );
    } catch (e) {
      onError?.call(e.toString());
    }
  }

  // ============================================================
  // 获取前三名
  // ============================================================
  
  List<LeaderboardEntry> getTopThree(
    LeaderboardType type, [
    LeaderboardPeriod period = LeaderboardPeriod.all,
  ]) {
    final data = getLeaderboard(type, period);
    if (data != null) {
      return data.entries.take(3).toList();
    }
    return [];
  }

  // ============================================================
  // 获取我的排名
  // ============================================================
  
  int getMyRank(LeaderboardType type) {
    return _myRanks[type] ?? -1;
  }

  // ============================================================
  // 获取我的数值
  // ============================================================
  
  int getMyValue(LeaderboardType type) {
    return _myValues[type] ?? 0;
  }

  // ============================================================
  // 检查是否上榜
  // ============================================================
  
  bool isOnBoard(LeaderboardType type) {
    return getMyRank(type) > 0;
  }

  // ============================================================
  // 刷新排行榜
  // ============================================================
  
  Future<void> refreshLeaderboard(
    LeaderboardType type, [
    LeaderboardPeriod period = LeaderboardPeriod.all,
  ]) async {
    _cache.remove(type);
    _lastRefreshTime.remove(type);
    await requestLeaderboard(type, period);
  }

  // ============================================================
  // 清空缓存
  // ============================================================
  
  void clearCache() {
    _cache.clear();
    _lastRefreshTime.clear();
  }
}

// ============================================================
// 排行榜页面
// ============================================================

class LeaderboardPage extends StatefulWidget {
  const LeaderboardPage({Key? key}) : super(key: key);

  @override
  State<LeaderboardPage> createState() => _LeaderboardPageState();
}

class _LeaderboardPageState extends State<LeaderboardPage> with SingleTickerProviderStateMixin {
  late TabController _tabController;
  LeaderboardType _currentType = LeaderboardType.level;
  LeaderboardPeriod _currentPeriod = LeaderboardPeriod.all;
  LeaderboardData? _currentData;
  bool _isLoading = false;
  String? _error;

  final List<_TabInfo> _tabs = [
    _TabInfo(LeaderboardType.level, '等级'),
    _TabInfo(LeaderboardType.gold, '金币'),
    _TabInfo(LeaderboardType.gem, '钻石'),
    _TabInfo(LeaderboardType.combat, '战力'),
    _TabInfo(LeaderboardType.kill, '击杀'),
    _TabInfo(LeaderboardType.guild, '公会'),
  ];

  @override
  void initState() {
    super.initState();
    _tabController = TabController(length: _tabs.length, vsync: this);
    _tabController.addListener(_onTabChanged);
    _loadLeaderboard();
  }

  @override
  void dispose() {
    _tabController.dispose();
    super.dispose();
  }

  void _onTabChanged() {
    if (!_tabController.indexIsChanging) {
      setState(() {
        _currentType = _tabs[_tabController.index].type;
      });
      _loadLeaderboard();
    }
  }

  Future<void> _loadLeaderboard() async {
    setState(() {
      _isLoading = true;
      _error = null;
    });

    try {
      final client = LeaderboardClient();
      await client.requestLeaderboard(_currentType, _currentPeriod);
      await client.requestMyRank(_currentType);
      
      setState(() {
        _currentData = client.getLeaderboard(_currentType, _currentPeriod);
        _isLoading = false;
      });
    } catch (e) {
      setState(() {
        _error = e.toString();
        _isLoading = false;
      });
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text('排行榜'),
        bottom: TabBar(
          controller: _tabController,
          tabs: _tabs.map((t) => Tab(text: t.label)).toList(),
          isScrollable: true,
        ),
        actions: [
          PopupMenuButton<LeaderboardPeriod>(
            icon: const Icon(Icons.filter_list),
            onSelected: (period) {
              setState(() => _currentPeriod = period);
              _loadLeaderboard();
            },
            itemBuilder: (context) => [
              const PopupMenuItem(value: LeaderboardPeriod.all, child: Text('全部')),
              const PopupMenuItem(value: LeaderboardPeriod.daily, child: Text('今日')),
              const PopupMenuItem(value: LeaderboardPeriod.weekly, child: Text('本周')),
              const PopupMenuItem(value: LeaderboardPeriod.monthly, child: Text('本月')),
              const PopupMenuItem(value: LeaderboardPeriod.season, child: Text('本赛季')),
            ],
          ),
          IconButton(
            icon: const Icon(Icons.refresh),
            onPressed: _loadLeaderboard,
          ),
        ],
      ),
      body: _buildBody(),
    );
  }

  Widget _buildBody() {
    if (_isLoading) {
      return const Center(child: CircularProgressIndicator());
    }

    if (_error != null) {
      return Center(
        child: Column(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            Text('加载失败: $_error'),
            ElevatedButton(
              onPressed: _loadLeaderboard,
              child: const Text('重试'),
            ),
          ],
        ),
      );
    }

    if (_currentData == null || _currentData!.entries.isEmpty) {
      return const Center(child: Text('暂无数据'));
    }

    return Column(
      children: [
        // 我的排名
        _buildMyRankCard(),
        // 前三名
        _buildTopThree(),
        // 排行榜列表
        Expanded(
          child: ListView.builder(
            itemCount: _currentData!.entries.length,
            itemBuilder: (context, index) {
              return _buildEntryItem(_currentData!.entries[index]);
            },
          ),
        ),
      ],
    );
  }

  Widget _buildMyRankCard() {
    final myRank = LeaderboardClient().getMyRank(_currentType);
    final myValue = LeaderboardClient().getMyValue(_currentType);

    return Container(
      margin: const EdgeInsets.all(16),
      padding: const EdgeInsets.all(16),
      decoration: BoxDecoration(
        color: Colors.blue.shade50,
        borderRadius: BorderRadius.circular(12),
      ),
      child: Row(
        mainAxisAlignment: MainAxisAlignment.spaceAround,
        children: [
          Column(
            children: [
              const Text('我的排名', style: TextStyle(color: Colors.grey)),
              const SizedBox(height: 4),
              Text(
                myRank > 0 ? '#$myRank' : '未上榜',
                style: TextStyle(
                  fontSize: 24,
                  fontWeight: FontWeight.bold,
                  color: myRank > 0 ? Colors.blue : Colors.grey,
                ),
              ),
            ],
          ),
          Column(
            children: [
              const Text('数值', style: TextStyle(color: Colors.grey)),
              const SizedBox(height: 4),
              Text(
                myValue.toString(),
                style: const TextStyle(
                  fontSize: 24,
                  fontWeight: FontWeight.bold,
                ),
              ),
            ],
          ),
        ],
      ),
    );
  }

  Widget _buildTopThree() {
    final topThree = LeaderboardClient().getTopThree(_currentType, _currentPeriod);
    if (topThree.isEmpty) return const SizedBox.shrink();

    return Container(
      height: 120,
      padding: const EdgeInsets.symmetric(horizontal: 16),
      child: Row(
        mainAxisAlignment: MainAxisAlignment.spaceEvenly,
        children: [
          if (topThree.length > 1) _buildTopThreeItem(topThree[1], 2),
          if (topThree.isNotEmpty) _buildTopThreeItem(topThree[0], 1),
          if (topThree.length > 2) _buildTopThreeItem(topThree[2], 3),
        ],
      ),
    );
  }

  Widget _buildTopThreeItem(LeaderboardEntry entry, int rank) {
    final colors = {
      1: Colors.amber,
      2: Colors.grey.shade400,
      3: Colors.brown.shade300,
    };

    return Column(
      mainAxisSize: MainAxisSize.min,
      children: [
        Container(
          width: 50,
          height: 50,
          decoration: BoxDecoration(
            shape: BoxShape.circle,
            border: Border.all(color: colors[rank]!, width: 3),
          ),
          child: Center(
            child: Text(
              entry.playerName.substring(0, 1),
              style: const TextStyle(fontSize: 20, fontWeight: FontWeight.bold),
            ),
          ),
        ),
        const SizedBox(height: 4),
        Text(entry.playerName, overflow: TextOverflow.ellipsis),
        Text(
          '#$rank',
          style: TextStyle(color: colors[rank], fontWeight: FontWeight.bold),
        ),
      ],
    );
  }

  Widget _buildEntryItem(LeaderboardEntry entry) {
    return Card(
      margin: const EdgeInsets.symmetric(horizontal: 16, vertical: 4),
      child: ListTile(
        leading: CircleAvatar(
          backgroundColor: entry.rank <= 3 ? Colors.amber : Colors.grey.shade200,
          child: Text('${entry.rank}'),
        ),
        title: Row(
          children: [
            Text(entry.playerName),
            if (entry.vipLevel > 0) ...[
              const SizedBox(width: 4),
              Icon(Icons.star, color: Colors.purple, size: 16),
              Text('V${entry.vipLevel}', style: const TextStyle(fontSize: 12)),
            ],
          ],
        ),
        subtitle: Text('Lv.${entry.level}'),
        trailing: Column(
          mainAxisAlignment: MainAxisAlignment.center,
          crossAxisAlignment: CrossAxisAlignment.end,
          children: [
            Text(
              entry.value.toString(),
              style: const TextStyle(fontWeight: FontWeight.bold),
            ),
            if (entry.guildName != null)
              Text(
                entry.guildName!,
                style: TextStyle(fontSize: 12, color: Colors.grey.shade600),
              ),
          ],
        ),
      ),
    );
  }
}

class _TabInfo {
  final LeaderboardType type;
  final String label;
  
  _TabInfo(this.type, this.label);
}

// ============================================================
// 排行榜条目组件
// ============================================================

class LeaderboardEntryWidget extends StatelessWidget {
  final LeaderboardEntry entry;
  final VoidCallback? onTap;

  const LeaderboardEntryWidget({
    Key? key,
    required this.entry,
    this.onTap,
  }) : super(key: key);

  @override
  Widget build(BuildContext context) {
    return Card(
      margin: const EdgeInsets.symmetric(horizontal: 16, vertical: 4),
      child: InkWell(
        onTap: onTap,
        child: Padding(
          padding: const EdgeInsets.all(12),
          child: Row(
            children: [
              // 排名
              Container(
                width: 40,
                height: 40,
                decoration: BoxDecoration(
                  color: _getRankColor(entry.rank),
                  shape: BoxShape.circle,
                ),
                child: Center(
                  child: Text(
                    '${entry.rank}',
                    style: const TextStyle(
                      color: Colors.white,
                      fontWeight: FontWeight.bold,
                    ),
                  ),
                ),
              ),
              const SizedBox(width: 12),
              // 头像和名称
              Expanded(
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Row(
                      children: [
                        Text(
                          entry.playerName,
                          style: const TextStyle(
                            fontWeight: FontWeight.bold,
                            fontSize: 16,
                          ),
                        ),
                        if (entry.isMyFriend)
                          const Padding(
                            padding: EdgeInsets.only(left: 4),
                            child: Icon(Icons.favorite, color: Colors.red, size: 16),
                          ),
                      ],
                    ),
                    const SizedBox(height: 4),
                    Text('Lv.${entry.level}'),
                  ],
                ),
              ),
              // 数值和公会
              Column(
                crossAxisAlignment: CrossAxisAlignment.end,
                children: [
                  Text(
                    entry.value.toString(),
                    style: const TextStyle(
                      fontWeight: FontWeight.bold,
                      fontSize: 18,
                    ),
                  ),
                  if (entry.guildName != null)
                    Text(
                      entry.guildName!,
                      style: TextStyle(
                        fontSize: 12,
                        color: Colors.grey.shade600,
                      ),
                    ),
                ],
              ),
            ],
          ),
        ),
      ),
    );
  }

  Color _getRankColor(int rank) {
    switch (rank) {
      case 1: return Colors.amber;
      case 2: return Colors.grey.shade400;
      case 3: return Colors.brown.shade300;
      default: return Colors.blue;
    }
  }
}
