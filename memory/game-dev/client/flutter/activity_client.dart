import 'dart:convert';
import 'package:http/http.dart' as http;

/// 活动类型
enum ActivityType {
  daily,
  weekly,
  limited,
  event,
  battle,
  gift,
  recharge,
}

/// 活动状态
enum ActivityState {
  pending,
  active,
  ended,
  reward,
}

/// 活动奖励
class ActivityReward {
  final String type;
  final int value;
  final String? itemId;
  final int itemCount;

  ActivityReward({
    required this.type,
    required this.value,
    this.itemId,
    required this.itemCount,
  });

  factory ActivityReward.fromJson(Map<String, dynamic> json) {
    return ActivityReward(
      type: json['type'] ?? 'coin',
      value: json['value'] ?? 0,
      itemId: json['item_id'],
      itemCount: json['item_count'] ?? 1,
    );
  }

  Map<String, dynamic> toJson() => {
    'type': type,
    'value': value,
    'item_id': itemId,
    'item_count': itemCount,
  };

  String get displayName {
    switch (type) {
      case 'coin':
        return '$value 金币';
      case 'diamond':
        return '$value 钻石';
      case 'item':
        return '$itemCount 个 $itemId';
      default:
        return '$type: $value';
    }
  }
}

/// 活动条件
class ActivityCondition {
  final String type;
  final String target;
  final int value;
  final String description;

  ActivityCondition({
    required this.type,
    required this.target,
    required this.value,
    required this.description,
  });

  factory ActivityCondition.fromJson(Map<String, dynamic> json) {
    return ActivityCondition(
      type: json['type'] ?? '',
      target: json['target'] ?? '',
      value: json['value'] ?? 0,
      description: json['description'] ?? '',
    );
  }

  Map<String, dynamic> toJson() => {
    'type': type,
    'target': target,
    'value': value,
    'description': description,
  };
}

/// 活动信息
class Activity {
  final String id;
  final String name;
  final ActivityType type;
  final ActivityState state;
  final int startTime;
  final int endTime;
  final List<ActivityReward> rewards;
  final List<ActivityCondition> conditions;
  final String content;
  final String icon;
  final int sort;
  final Map<String, dynamic>? config;

  Activity({
    required this.id,
    required this.name,
    required this.type,
    required this.state,
    required this.startTime,
    required this.endTime,
    required this.rewards,
    required this.conditions,
    required this.content,
    required this.icon,
    required this.sort,
    this.config,
  });

  factory Activity.fromJson(Map<String, dynamic> json) {
    return Activity(
      id: json['id'] ?? '',
      name: json['name'] ?? '',
      type: ActivityType.values[json['type'] ?? 0],
      state: ActivityState.values[json['state'] ?? 0],
      startTime: json['start_time'] ?? 0,
      endTime: json['end_time'] ?? 0,
      rewards: (json['rewards'] as List<dynamic>?)
          ?.map((e) => ActivityReward.fromJson(e))
          .toList() ?? [],
      conditions: (json['conditions'] as List<dynamic>?)
          ?.map((e) => ActivityCondition.fromJson(e))
          .toList() ?? [],
      content: json['content'] ?? '',
      icon: json['icon'] ?? '',
      sort: json['sort'] ?? 0,
      config: json['config'],
    );
  }

  Map<String, dynamic> toJson() => {
    'id': id,
    'name': name,
    'type': type.index,
    'state': state.index,
    'start_time': startTime,
    'end_time': endTime,
    'rewards': rewards.map((e) => e.toJson()).toList(),
    'conditions': conditions.map((e) => e.toJson()).toList(),
    'content': content,
    'icon': icon,
    'sort': sort,
    'config': config,
  };

  bool get isActive => state == ActivityState.active;
  bool get isEnded => state == ActivityState.ended;

  Duration get remainingTime {
    final now = DateTime.now().millisecondsSinceEpoch;
    final remaining = endTime - now;
    return Duration(milliseconds: remaining > 0 ? remaining : 0);
  }

  String get remainingTimeText {
    final remaining = remainingTime;
    if (remaining.inDays > 0) {
      return '${remaining.inDays}天 ${remaining.inHours % 24}小时';
    } else if (remaining.inHours > 0) {
      return '${remaining.inHours}小时 ${remaining.inMinutes % 60}分钟';
    } else if (remaining.inMinutes > 0) {
      return '${remaining.inMinutes}分钟';
    } else {
      return '即将结束';
    }
  }

  String get typeName {
    switch (type) {
      case ActivityType.daily:
        return '每日活动';
      case ActivityType.weekly:
        return '每周活动';
      case ActivityType.limited:
        return '限时活动';
      case ActivityType.event:
        return '事件活动';
      case ActivityType.battle:
        return '对战活动';
      case ActivityType.gift:
        return '礼物活动';
      case ActivityType.recharge:
        return '充值活动';
    }
  }

  String get stateName {
    switch (state) {
      case ActivityState.pending:
        return '即将开始';
      case ActivityState.active:
        return '进行中';
      case ActivityState.ended:
        return '已结束';
      case ActivityState.reward:
        return '领奖中';
    }
  }
}

/// 玩家活动进度
class PlayerActivityProgress {
  final String playerId;
  final String activityId;
  final int progress;
  final bool claimed;
  final int? claimedAt;
  final Map<String, dynamic>? extraData;

  PlayerActivityProgress({
    required this.playerId,
    required this.activityId,
    required this.progress,
    required this.claimed,
    this.claimedAt,
    this.extraData,
  });

  factory PlayerActivityProgress.fromJson(Map<String, dynamic> json) {
    return PlayerActivityProgress(
      playerId: json['player_id'] ?? '',
      activityId: json['activity_id'] ?? '',
      progress: json['progress'] ?? 0,
      claimed: json['claimed'] ?? false,
      claimedAt: json['claimed_at'],
      extraData: json['extra_data'],
    );
  }

  Map<String, dynamic> toJson() => {
    'player_id': playerId,
    'activity_id': activityId,
    'progress': progress,
    'claimed': claimed,
    'claimed_at': claimedAt,
    'extra_data': extraData,
  };
}

/// 活动列表响应
class ActivityListResponse {
  final List<Activity> activities;
  final int total;
  final int page;
  final int pageSize;

  ActivityListResponse({
    required this.activities,
    required this.total,
    required this.page,
    required this.pageSize,
  });

  factory ActivityListResponse.fromJson(Map<String, dynamic> json) {
    return ActivityListResponse(
      activities: (json['activities'] as List<dynamic>?)
          ?.map((e) => Activity.fromJson(e))
          .toList() ?? [],
      total: json['total'] ?? 0,
      page: json['page'] ?? 1,
      pageSize: json['page_size'] ?? 20,
    );
  }
}

/// 活动详情响应
class ActivityDetailResponse {
  final Activity activity;
  final PlayerActivityProgress? progress;
  final List<ActivityReward> claimableRewards;

  ActivityDetailResponse({
    required this.activity,
    this.progress,
    required this.claimableRewards,
  });

  factory ActivityDetailResponse.fromJson(Map<String, dynamic> json) {
    return ActivityDetailResponse(
      activity: Activity.fromJson(json['activity'] ?? {}),
      progress: json['progress'] != null
          ? PlayerActivityProgress.fromJson(json['progress'])
          : null,
      claimableRewards: (json['claimable_rewards'] as List<dynamic>?)
          ?.map((e) => ActivityReward.fromJson(e))
          .toList() ?? [],
    );
  }
}

/// 领取奖励响应
class ClaimRewardResponse {
  final bool success;
  final String message;
  final List<ActivityReward> rewards;
  final Map<String, dynamic>? newBalance;

  ClaimRewardResponse({
    required this.success,
    required this.message,
    required this.rewards,
    this.newBalance,
  });

  factory ClaimRewardResponse.fromJson(Map<String, dynamic> json) {
    return ClaimRewardResponse(
      success: json['success'] ?? false,
      message: json['message'] ?? '',
      rewards: (json['rewards'] as List<dynamic>?)
          ?.map((e) => ActivityReward.fromJson(e))
          .toList() ?? [],
      newBalance: json['new_balance'],
    );
  }
}

/// 活动客户端
class ActivityClient {
  final String baseUrl;
  final String token;
  final http.Client _client;

  ActivityClient({
    required this.baseUrl,
    required this.token,
    http.Client? client,
  }) : _client = client ?? http.Client();

  Map<String, String> get _headers => {
    'Content-Type': 'application/json',
    'Authorization': 'Bearer $token',
  };

  /// 获取活动列表
  Future<ActivityListResponse> getActivityList({
    ActivityType? type,
    ActivityState? state,
    int page = 1,
    int pageSize = 20,
  }) async {
    final queryParams = {
      if (type != null) 'type': type.index.toString(),
      if (state != null) 'state': state.index.toString(),
      'page': page.toString(),
      'page_size': pageSize.toString(),
    };

    final uri = Uri.parse('$baseUrl/api/activity/list')
        .replace(queryParameters: queryParams);

    final response = await _client.get(uri, headers: _headers);

    if (response.statusCode == 200) {
      final data = json.decode(response.body);
      return ActivityListResponse.fromJson(data);
    } else {
      throw Exception('获取活动列表失败: ${response.statusCode}');
    }
  }

  /// 获取活动详情
  Future<ActivityDetailResponse> getActivityDetail(String activityId) async {
    final uri = Uri.parse('$baseUrl/api/activity/detail')
        .replace(queryParameters: {'activity_id': activityId});

    final response = await _client.get(uri, headers: _headers);

    if (response.statusCode == 200) {
      final data = json.decode(response.body);
      return ActivityDetailResponse.fromJson(data);
    } else {
      throw Exception('获取活动详情失败: ${response.statusCode}');
    }
  }

  /// 领取活动奖励
  Future<ClaimRewardResponse> claimReward(String activityId) async {
    final uri = Uri.parse('$baseUrl/api/activity/claim');
    final body = json.encode({'activity_id': activityId});

    final response = await _client.post(
      uri,
      headers: _headers,
      body: body,
    );

    if (response.statusCode == 200) {
      final data = json.decode(response.body);
      return ClaimRewardResponse.fromJson(data);
    } else {
      throw Exception('领取奖励失败: ${response.statusCode}');
    }
  }

  /// 获取我的活动进度
  Future<List<PlayerActivityProgress>> getMyProgress({
    ActivityType? type,
  }) async {
    final queryParams = {
      if (type != null) 'type': type.index.toString(),
    };

    final uri = Uri.parse('$baseUrl/api/activity/my_progress')
        .replace(queryParameters: queryParams);

    final response = await _client.get(uri, headers: _headers);

    if (response.statusCode == 200) {
      final data = json.decode(response.body);
      return (data['progress'] as List<dynamic>?)
          ?.map((e) => PlayerActivityProgress.fromJson(e))
          .toList() ?? [];
    } else {
      throw Exception('获取活动进度失败: ${response.statusCode}');
    }
  }

  /// 订阅活动更新
  Future<bool> subscribeActivity(List<String> activityIds) async {
    final uri = Uri.parse('$baseUrl/api/activity/subscribe');
    final body = json.encode({'activity_ids': activityIds});

    final response = await _client.post(
      uri,
      headers: _headers,
      body: body,
    );

    return response.statusCode == 200;
  }

  /// 获取进行中的活动
  Future<List<Activity>> getActiveActivities() async {
    final response = await getActivityList(state: ActivityState.active);
    return response.activities;
  }

  /// 获取可领取奖励的活动
  Future<List<Activity>> getClaimableActivities() async {
    final activities = await getActiveActivities();
    return activities.where((a) => a.conditions.isNotEmpty).toList();
  }

  /// 批量领取多个活动的奖励
  Future<Map<String, ClaimRewardResponse>> batchClaimRewards(
    List<String> activityIds,
  ) async {
    final results = <String, ClaimRewardResponse>{};

    for (final activityId in activityIds) {
      try {
        final result = await claimReward(activityId);
        results[activityId] = result;
      } catch (e) {
        results[activityId] = ClaimRewardResponse(
          success: false,
          message: e.toString(),
          rewards: [],
        );
      }
    }

    return results;
  }

  /// 释放资源
  void dispose() {
    _client.close();
  }
}

/// 活动页面管理器
class ActivityUIManager {
  final ActivityClient _client;
  List<Activity> _activities = [];
  Map<String, PlayerActivityProgress> _progressMap = {};
  ActivityType? _currentFilter;
  bool _isLoading = false;
  String? _error;

  ActivityUIManager(this._client);

  List<Activity> get activities => _activities;
  Map<String, PlayerActivityProgress> get progressMap => _progressMap;
  ActivityType? get currentFilter => _currentFilter;
  bool get isLoading => _isLoading;
  String? get error => _error;

  List<Activity> get filteredActivities {
    if (_currentFilter == null) return _activities;
    return _activities.where((a) => a.type == _currentFilter).toList();
  }

  List<Activity> get activeActivities =>
      _activities.where((a) => a.isActive).toList();

  List<Activity> get pendingActivities =>
      _activities.where((a) => a.state == ActivityState.pending).toList();

  /// 加载活动列表
  Future<void> loadActivities({
    ActivityType? type,
    bool refresh = false,
  }) async {
    if (_isLoading) return;

    _isLoading = true;
    _error = null;

    try {
      final response = await _client.getActivityList(
        type: type,
        page: 1,
        pageSize: 50,
      );
      _activities = response.activities;
      _currentFilter = type;

      // 加载进度
      await loadProgress();
    } catch (e) {
      _error = e.toString();
    } finally {
      _isLoading = false;
    }
  }

  /// 加载活动进度
  Future<void> loadProgress() async {
    try {
      final progressList = await _client.getMyProgress(type: _currentFilter);
      _progressMap = {
        for (var p in progressList) p.activityId: p,
      };
    } catch (e) {
      // 忽略进度加载错误
    }
  }

  /// 刷新活动列表
  Future<void> refresh() async {
    await loadActivities(type: _currentFilter, refresh: true);
  }

  /// 切换筛选
  Future<void> setFilter(ActivityType? type) async {
    if (_currentFilter == type) return;
    await loadActivities(type: type);
  }

  /// 领取奖励
  Future<ClaimRewardResponse> claimReward(String activityId) async {
    final result = await _client.claimReward(activityId);
    if (result.success) {
      // 更新本地状态
      final index = _activities.indexWhere((a) => a.id == activityId);
      if (index != -1) {
        final activity = _activities[index];
        _activities[index] = Activity(
          id: activity.id,
          name: activity.name,
          type: activity.type,
          state: ActivityState.reward,
          startTime: activity.startTime,
          endTime: activity.endTime,
          rewards: activity.rewards,
          conditions: activity.conditions,
          content: activity.content,
          icon: activity.icon,
          sort: activity.sort,
          config: activity.config,
        );
      }
    }
    return result;
  }

  /// 获取活动进度
  PlayerActivityProgress? getProgress(String activityId) {
    return _progressMap[activityId];
  }

  /// 检查是否可以领取
  bool canClaim(String activityId) {
    final progress = _progressMap[activityId];
    final activity = _activities.firstWhere(
      (a) => a.id == activityId,
      orElse: () => Activity(
        id: '',
        name: '',
        type: ActivityType.daily,
        state: ActivityState.active,
        startTime: 0,
        endTime: 0,
        rewards: [],
        conditions: [],
        content: '',
        icon: '',
        sort: 0,
      ),
    );

    if (!activity.isActive || progress == null || progress.claimed) {
      return false;
    }

    // 检查是否满足条件
    for (final condition in activity.conditions) {
      if (progress.progress < condition.value) {
        return false;
      }
    }

    return true;
  }

  /// 获取可领取的活动数量
  int get claimableCount {
    return _activities.where((a) => canClaim(a.id)).length;
  }
}

/// 活动条目组件
class ActivityItem extends StatelessWidget {
  final Activity activity;
  final PlayerActivityProgress? progress;
  final VoidCallback? onTap;
  final VoidCallback? onClaim;

  const ActivityItem({
    super.key,
    required this.activity,
    this.progress,
    this.onTap,
    this.onClaim,
  });

  @override
  Widget build(BuildContext context) {
    final canClaim = progress != null &&
        !progress!.claimed &&
        activity.conditions.every(
          (c) => progress!.progress >= c.value,
        );

    return Card(
      margin: const EdgeInsets.symmetric(horizontal: 16, vertical: 8),
      child: InkWell(
        onTap: onTap,
        child: Padding(
          padding: const EdgeInsets.all(16),
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              Row(
                children: [
                  // 活动图标
                  Container(
                    width: 48,
                    height: 48,
                    decoration: BoxDecoration(
                      color: _getTypeColor(activity.type),
                      borderRadius: BorderRadius.circular(8),
                    ),
                    child: Center(
                      child: Text(
                        activity.icon,
                        style: const TextStyle(fontSize: 24),
                      ),
                    ),
                  ),
                  const SizedBox(width: 12),
                  // 活动信息
                  Expanded(
                    child: Column(
                      crossAxisAlignment: CrossAxisAlignment.start,
                      children: [
                        Text(
                          activity.name,
                          style: const TextStyle(
                            fontSize: 16,
                            fontWeight: FontWeight.bold,
                          ),
                        ),
                        const SizedBox(height: 4),
                        Text(
                          '${activity.typeName} · ${activity.stateName}',
                          style: TextStyle(
                            fontSize: 12,
                            color: Colors.grey[600],
                          ),
                        ),
                      ],
                    ),
                  ),
                  // 状态标签
                  if (activity.isActive)
                    Container(
                      padding: const EdgeInsets.symmetric(
                        horizontal: 8,
                        vertical: 4,
                      ),
                      decoration: BoxDecoration(
                        color: Colors.green,
                        borderRadius: BorderRadius.circular(4),
                      ),
                      child: Text(
                        activity.remainingTimeText,
                        style: const TextStyle(
                          color: Colors.white,
                          fontSize: 10,
                        ),
                      ),
                    ),
                ],
              ),
              if (activity.conditions.isNotEmpty) ...[
                const SizedBox(height: 12),
                // 进度条
                _buildProgressBar(),
              ],
              if (activity.rewards.isNotEmpty) ...[
                const SizedBox(height: 12),
                // 奖励预览
                _buildRewardsPreview(),
              ],
              if (canClaim) ...[
                const SizedBox(height: 12),
                // 领取按钮
                SizedBox(
                  width: double.infinity,
                  child: ElevatedButton(
                    onPressed: onClaim,
                    style: ElevatedButton.styleFrom(
                      backgroundColor: Colors.orange,
                      foregroundColor: Colors.white,
                    ),
                    child: const Text('领取奖励'),
                  ),
                ),
              ],
            ],
          ),
        ),
      ),
    );
  }

  Widget _buildProgressBar() {
    final progress = this.progress;
    if (progress == null) return const SizedBox();

    // 计算总体进度
    double totalProgress = 0;
    if (activity.conditions.isNotEmpty) {
      totalProgress = activity.conditions
              .map((c) => (progress.progress / c.value).clamp(0.0, 1.0))
              .reduce((a, b) => a + b) /
          activity.conditions.length;
    }

    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Row(
          mainAxisAlignment: MainAxisAlignment.spaceBetween,
          children: [
            Text(
              '进度: ${progress.progress}',
              style: const TextStyle(fontSize: 12),
            ),
            Text(
              '${(totalProgress * 100).toInt()}%',
              style: const TextStyle(fontSize: 12),
            ),
          ],
        ),
        const SizedBox(height: 4),
        LinearProgressIndicator(
          value: totalProgress,
          backgroundColor: Colors.grey[200],
          valueColor: AlwaysStoppedAnimation<Color>(
            canClaim ? Colors.orange : Colors.blue,
          ),
        ),
      ],
    );
  }

  Widget _buildRewardsPreview() {
    return Row(
      children: [
        const Text(
          '奖励: ',
          style: TextStyle(fontSize: 12),
        ),
        Expanded(
          child: Wrap(
            spacing: 4,
            runSpacing: 4,
            children: activity.rewards.take(3).map((reward) {
              return Container(
                padding: const EdgeInsets.symmetric(
                  horizontal: 6,
                  vertical: 2,
                ),
                decoration: BoxDecoration(
                  color: Colors.grey[100],
                  borderRadius: BorderRadius.circular(4),
                ),
                child: Text(
                  reward.displayName,
                  style: const TextStyle(fontSize: 10),
                ),
              );
            }).toList(),
          ),
        ),
        if (activity.rewards.length > 3)
          Text(
            '+${activity.rewards.length - 3}项',
            style: TextStyle(
              fontSize: 10,
              color: Colors.grey[600],
            ),
          ),
      ],
    );
  }

  Color _getTypeColor(ActivityType type) {
    switch (type) {
      case ActivityType.daily:
        return Colors.blue;
      case ActivityType.weekly:
        return Colors.purple;
      case ActivityType.limited:
        return Colors.red;
      case ActivityType.event:
        return Colors.orange;
      case ActivityType.battle:
        return Colors.green;
      case ActivityType.gift:
        return Colors.pink;
      case ActivityType.recharge:
        return Colors.amber;
    }
  }
}

/// 活动详情面板
class ActivityDetailPanel extends StatelessWidget {
  final Activity activity;
  final PlayerActivityProgress? progress;
  final VoidCallback? onClaim;
  final VoidCallback? onClose;

  const ActivityDetailPanel({
    super.key,
    required this.activity,
    this.progress,
    this.onClaim,
    this.onClose,
  });

  @override
  Widget build(BuildContext context) {
    final canClaim = progress != null &&
        !progress!.claimed &&
        activity.conditions.every(
          (c) => progress!.progress >= c.value,
        );

    return Dialog(
      child: Container(
        padding: const EdgeInsets.all(16),
        constraints: const BoxConstraints(maxWidth: 400),
        child: Column(
          mainAxisSize: MainAxisSize.min,
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            // 标题栏
            Row(
              children: [
                Expanded(
                  child: Text(
                    activity.name,
                    style: const TextStyle(
                      fontSize: 18,
                      fontWeight: FontWeight.bold,
                    ),
                  ),
                ),
                IconButton(
                  icon: const Icon(Icons.close),
                  onPressed: onClose,
                ),
              ],
            ),
            const Divider(),
            // 活动时间
            Text(
              '⏰ ${_formatTime(activity.startTime)} - ${_formatTime(activity.endTime)}',
              style: TextStyle(color: Colors.grey[600]),
            ),
            if (activity.isActive) ...[
              const SizedBox(height: 4),
              Text(
                '剩余时间: ${activity.remainingTimeText}',
                style: const TextStyle(color: Colors.orange),
              ),
            ],
            const SizedBox(height: 12),
            // 活动内容
            Text(
              activity.content,
              style: const TextStyle(fontSize: 14),
            ),
            const SizedBox(height: 16),
            // 条件列表
            if (activity.conditions.isNotEmpty) ...[
              const Text(
                '任务条件',
                style: TextStyle(
                  fontSize: 14,
                  fontWeight: FontWeight.bold,
                ),
              ),
              const SizedBox(height: 8),
              ...activity.conditions.map((condition) {
                final currentProgress = progress?.progress ?? 0;
                final isCompleted = currentProgress >= condition.value;

                return Padding(
                  padding: const EdgeInsets.symmetric(vertical: 4),
                  child: Row(
                    children: [
                      Icon(
                        isCompleted
                            ? Icons.check_circle
                            : Icons.radio_button_unchecked,
                        color: isCompleted ? Colors.green : Colors.grey,
                        size: 16,
                      ),
                      const SizedBox(width: 8),
                      Expanded(
                        child: Text(
                          condition.description,
                          style: TextStyle(
                            fontSize: 12,
                            decoration: isCompleted
                                ? TextDecoration.lineThrough
                                : null,
                          ),
                        ),
                      ),
                      Text(
                        '$currentProgress/${condition.value}',
                        style: const TextStyle(fontSize: 12),
                      ),
                    ],
                  ),
                );
              }),
              const SizedBox(height: 16),
            ],
            // 奖励列表
            if (activity.rewards.isNotEmpty) ...[
              const Text(
                '活动奖励',
                style: TextStyle(
                  fontSize: 14,
                  fontWeight: FontWeight.bold,
                ),
              ),
              const SizedBox(height: 8),
              Wrap(
                spacing: 8,
                runSpacing: 8,
                children: activity.rewards.map((reward) {
                  return Container(
                    padding: const EdgeInsets.all(8),
                    decoration: BoxDecoration(
                      color: Colors.grey[100],
                      borderRadius: BorderRadius.circular(8),
                    ),
                    child: Column(
                      mainAxisSize: MainAxisSize.min,
                      children: [
                        Text(
                          reward.displayName,
                          style: const TextStyle(fontSize: 12),
                        ),
                      ],
                    ),
                  );
                }).toList(),
              ),
              const SizedBox(height: 16),
            ],
            // 领取按钮
            if (canClaim)
              SizedBox(
                width: double.infinity,
                child: ElevatedButton(
                  onPressed: onClaim,
                  style: ElevatedButton.styleFrom(
                    backgroundColor: Colors.orange,
                    foregroundColor: Colors.white,
                  ),
                  child: const Text('领取奖励'),
                ),
              ),
            if (progress?.claimed ?? false)
              Container(
                width: double.infinity,
                padding: const EdgeInsets.all(12),
                decoration: BoxDecoration(
                  color: Colors.grey[200],
                  borderRadius: BorderRadius.circular(8),
                ),
                child: const Text(
                  '✅ 已领取',
                  textAlign: TextAlign.center,
                ),
              ),
          ],
        ),
      ),
    );
  }

  String _formatTime(int timestamp) {
    final dt = DateTime.fromMillisecondsSinceEpoch(timestamp);
    return '${dt.month}-${dt.day} ${dt.hour.toString().padLeft(2, '0')}:${dt.minute.toString().padLeft(2, '0')}';
  }
}
