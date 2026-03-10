// ============================================
// 签到客户端 - Flutter UI
// ============================================

import 'dart:async';
import 'dart:convert';
import 'package:flutter/material.dart';
import 'package:web_socket_channel/web_socket_channel.dart';

/// 签到客户端
class SignInClient {
  WebSocketChannel? _channel;
  final String serverUrl;
  final String token;
  
  final _statusController = StreamController<SignInStatus>.broadcast();
  final _calendarController = StreamController<SignInCalendar>.broadcast();
  
  Stream<SignInStatus> get status => _statusController.stream;
  Stream<SignInCalendar> get calendar => _calendarController.stream;
  
  SignInClient({required this.serverUrl, required this.token});
  
  void connect() {
    _channel = WebSocketChannel.connect(
      Uri.parse('$serverUrl?token=$token&module=signin'),
    );
    _channel!.stream.listen(_onMessage);
  }
  
  void _onMessage(dynamic data) {
    final json = jsonDecode(data);
    final msgId = json['msg_id'];
    final jsonData = json['data'];
    
    switch (msgId) {
      case 30001: // SCSignInStatus
        _statusController.add(SignInStatus.fromJson(jsonData));
        break;
      case 30002: // SCSignInCalendar
        _calendarController.add(SignInCalendar.fromJson(jsonData));
        break;
      case 30003: // SCSignInResult
        // 签到结果
        break;
    }
  }
  
  /// 获取签到状态
  void getSignInStatus() {
    _channel?.sink.add(jsonEncode({
      'msg_id': 30001,
      'data': {},
    }));
  }
  
  /// 获取签到日历
  void getSignInCalendar({int? month}) {
    _channel?.sink.add(jsonEncode({
      'msg_id': 30002,
      'data': month != null ? {'month': month} : {},
    }));
  }
  
  /// 签到
  void signIn() {
    _channel?.sink.add(jsonEncode({
      'msg_id': 30003,
      'data': {},
    }));
  }
  
  /// 获取签到排行榜
  void getSignInRank({int page = 1, int pageSize = 20}) {
    _channel?.sink.add(jsonEncode({
      'msg_id': 30004,
      'data': {'page': page, 'page_size': pageSize},
    }));
  }
  
  void dispose() {
    _channel?.sink.close();
    _statusController.close();
    _calendarController.close();
  }
}

/// 签到状态
class SignInStatus {
  final bool canSignIn;
  final bool hasSignedInToday;
  final int consecutiveDays;
  final int totalDays;
  final int monthDays;
  final List<SignInReward> todayRewards;
  final List<SignInReward> tomorrowRewards;
  
  SignInStatus({
    required this.canSignIn,
    required this.hasSignedInToday,
    required this.consecutiveDays,
    required this.totalDays,
    required this.monthDays,
    required this.todayRewards,
    required this.tomorrowRewards,
  });
  
  factory SignInStatus.fromJson(Map<String, dynamic> json) {
    return SignInStatus(
      canSignIn: json['can_sign_in'] ?? false,
      hasSignedInToday: json['has_signed_in_today'] ?? false,
      consecutiveDays: json['consecutive_days'] ?? 0,
      totalDays: json['total_days'] ?? 0,
      monthDays: json['month_days'] ?? 0,
      todayRewards: (json['today_rewards'] as List?)
          ?.map((r) => SignInReward.fromJson(r))
          .toList() ?? [],
      tomorrowRewards: (json['tomorrow_rewards'] as List?)
          ?.map((r) => SignInReward.fromJson(r))
          .toList() ?? [],
    );
  }
}

/// 签到奖励
class SignInReward {
  final String itemId;
  final String itemName;
  final int count;
  final RewardType type;
  final bool isVipOnly;
  
  SignInReward({
    required this.itemId,
    required this.itemName,
    required this.count,
    required this.type,
    required this.isVipOnly,
  });
  
  factory SignInReward.fromJson(Map<String, dynamic> json) {
    return SignInReward(
      itemId: json['item_id'],
      itemName: json['item_name'],
      count: json['count'],
      type: RewardType.values.firstWhere(
        (e) => e.name == json['type'],
        orElse: () => RewardType.item,
      ),
      isVipOnly: json['is_vip_only'] ?? false,
    );
  }
}

/// 奖励类型
enum RewardType {
  gold,
  gem,
  item,
  equipment,
}

/// 签到日历
class SignInCalendar {
  final int year;
  final int month;
  final List<SignInDay> days;
  
  SignInCalendar({
    required this.year,
    required this.month,
    required this.days,
  });
  
  factory SignInCalendar.fromJson(Map<String, dynamic> json) {
    return SignInCalendar(
      year: json['year'],
      month: json['month'],
      days: (json['days'] as List?)
          ?.map((d) => SignInDay.fromJson(d))
          .toList() ?? [],
    );
  }
}

/// 签到日
class SignInDay {
  final int day;
  final bool signedIn;
  final List<SignInReward> rewards;
  
  SignInDay({
    required this.day,
    required this.signedIn,
    required this.rewards,
  });
  
  factory SignInDay.fromJson(Map<String, dynamic> json) {
    return SignInDay(
      day: json['day'],
      signedIn: json['signed_in'] ?? false,
      rewards: (json['rewards'] as List?)
          ?.map((r) => SignInReward.fromJson(r))
          .toList() ?? [],
    );
  }
}

/// 签到记录
class SignInRecord {
  final int consecutiveDays;
  final int totalDays;
  final int signedAt;
  
  SignInRecord({
    required this.consecutiveDays,
    required this.totalDays,
    required this.signedAt,
  });
  
  factory SignInRecord.fromJson(Map<String, dynamic> json) {
    return SignInRecord(
      consecutiveDays: json['consecutive_days'],
      totalDays: json['total_days'],
      signedAt: json['signed_at'],
    );
  }
}

// ============ 签到页面 ============

class SignInPage extends StatefulWidget {
  final SignInClient signInClient;
  
  const SignInPage({super.key, required this.signInClient});
  
  @override
  State<SignInPage> createState() => _SignInPageState();
}

class _SignInPageState extends State<SignInPage> {
  SignInStatus? _status;
  SignInCalendar? _calendar;
  bool _isLoading = false;
  bool _isSigningIn = false;
  
  @override
  void initState() {
    super.initState();
    _loadData();
  }
  
  void _loadData() {
    setState(() => _isLoading = true);
    widget.signInClient.getSignInStatus();
    widget.signInClient.getSignInCalendar();
    
    widget.signInClient.status.listen((status) {
      setState(() {
        _status = status;
        _isLoading = false;
      });
    });
    
    widget.signInClient.calendar.listen((calendar) {
      setState(() => _calendar = calendar);
    });
  }
  
  void _signIn() async {
    setState(() => _isSigningIn = true);
    widget.signInClient.signIn();
    
    await Future.delayed(const Duration(milliseconds: 500));
    
    setState(() => _isSigningIn = false);
    _loadData();
  }
  
  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text('签到'),
        actions: [
          IconButton(
            icon: const Icon(Icons.leaderboard),
            onPressed: () => _showRankDialog(),
          ),
        ],
      ),
      body: _isLoading
          ? const Center(child: CircularProgressIndicator())
          : _buildContent(),
    );
  }
  
  Widget _buildContent() {
    if (_status == null || _calendar == null) {
      return const Center(child: Text('加载中...'));
    }
    
    return SingleChildScrollView(
      padding: const EdgeInsets.all(16),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.stretch,
        children: [
          _buildStatusCard(),
          const SizedBox(height: 16),
          _buildTodayRewardCard(),
          const SizedBox(height: 16),
          _buildCalendarCard(),
        ],
      ),
    );
  }
  
  Widget _buildStatusCard() {
    return Card(
      child: Padding(
        padding: const EdgeInsets.all(16),
        child: Column(
          children: [
            Row(
              mainAxisAlignment: MainAxisAlignment.spaceAround,
              children: [
                _StatusItem(
                  label: '连续签到',
                  value: '${_status!.consecutiveDays}天',
                  icon: Icons.local_fire_department,
                  color: Colors.orange,
                ),
                _StatusItem(
                  label: '本月签到',
                  value: '${_status!.monthDays}天',
                  icon: Icons.calendar_month,
                  color: Colors.blue,
                ),
                _StatusItem(
                  label: '总签到',
                  value: '${_status!.totalDays}天',
                  icon: Icons.emoji_events,
                  color: Colors.amber,
                ),
              ],
            ),
            const SizedBox(height: 16),
            ElevatedButton.icon(
              onPressed: _status!.canSignIn && !_isSigningIn
                  ? _signIn
                  : null,
              icon: _isSigningIn
                  ? const SizedBox(
                      width: 16,
                      height: 16,
                      child: CircularProgressIndicator(strokeWidth: 2),
                    )
                  : Icon(_status!.hasSignedInToday
                      ? Icons.check
                      : Icons.edit),
              label: Text(_status!.hasSignedInToday
                  ? '今日已签到'
                  : '立即签到'),
              style: ElevatedButton.styleFrom(
                minimumSize: const Size(double.infinity, 48),
                backgroundColor: _status!.hasSignedInToday
                    ? Colors.green
                    : Colors.blue,
              ),
            ),
          ],
        ),
      ),
    );
  }
  
  Widget _buildTodayRewardCard() {
    return Card(
      child: Padding(
        padding: const EdgeInsets.all(16),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            const Text(
              '今日奖励',
              style: TextStyle(
                fontSize: 16,
                fontWeight: FontWeight.bold,
              ),
            ),
            const SizedBox(height: 12),
            Wrap(
              spacing: 8,
              runSpacing: 8,
              children: _status!.todayRewards
                  .map((r) => _RewardChip(reward: r))
                  .toList(),
            ),
            const Divider(height: 24),
            const Text(
              '明日奖励预览',
              style: TextStyle(fontSize: 14, color: Colors.grey),
            ),
            const SizedBox(height: 8),
            Wrap(
              spacing: 8,
              runSpacing: 8,
              children: _status!.tomorrowRewards
                  .map((r) => _RewardChip(reward: r, isPreview: true))
                  .toList(),
            ),
          ],
        ),
      ),
    );
  }
  
  Widget _buildCalendarCard() {
    return Card(
      child: Padding(
        padding: const EdgeInsets.all(16),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Row(
              mainAxisAlignment: MainAxisAlignment.spaceBetween,
              children: [
                Text(
                  '${_calendar!.year}年${_calendar!.month}月',
                  style: const TextStyle(
                    fontSize: 16,
                    fontWeight: FontWeight.bold,
                  ),
                ),
                Row(
                  children: [
                    IconButton(
                      icon: const Icon(Icons.chevron_left),
                      onPressed: () {
                        int month = _calendar!.month - 1;
                        int year = _calendar!.year;
                        if (month < 1) {
                          month = 12;
                          year--;
                        }
                        widget.signInClient.getSignInCalendar(
                          month: year * 100 + month,
                        );
                      },
                    ),
                    IconButton(
                      icon: const Icon(Icons.chevron_right),
                      onPressed: () {
                        int month = _calendar!.month + 1;
                        int year = _calendar!.year;
                        if (month > 12) {
                          month = 1;
                          year++;
                        }
                        widget.signInClient.getSignInCalendar(
                          month: year * 100 + month,
                        );
                      },
                    ),
                  ],
                ),
              ],
            ),
            const SizedBox(height: 8),
            Row(
              mainAxisAlignment: MainAxisAlignment.spaceAround,
              children: ['日', '一', '二', '三', '四', '五', '六']
                  .map((d) => SizedBox(
                        width: 36,
                        child: Text(
                          d,
                          textAlign: TextAlign.center,
                          style: const TextStyle(
                            fontWeight: FontWeight.bold,
                            color: Colors.grey,
                          ),
                        ),
                      ))
                  .toList(),
            ),
            const SizedBox(height: 8),
            _buildCalendarDays(),
          ],
        ),
      ),
    );
  }
  
  Widget _buildCalendarDays() {
    if (_calendar!.days.isEmpty) {
      return const Center(child: Text('本月暂无签到数据'));
    }
    
    final firstDay = _calendar!.days.first;
    final startWeekday = DateTime(
      _calendar!.year,
      _calendar!.month,
      1,
    ).weekday % 7;
    
    final cells = <Widget>[];
    
    // 空白天数
    for (int i = 0; i < startWeekday; i++) {
      cells.add(const SizedBox(width: 36, height: 36));
    }
    
    // 实际天数
    for (final day in _calendar!.days) {
      cells.add(_CalendarDayCell(day: day));
    }
    
    return Wrap(
      spacing: (MediaQuery.of(context).size.width - 32 - 16) / 7 - 36,
      runSpacing: 8,
      children: cells,
    );
  }
  
  void _showRankDialog() {
    showDialog(
      context: context,
      builder: (context) => AlertDialog(
        title: const Text('签到排行榜'),
        content: const SizedBox(
          width: 300,
          height: 400,
          child: Center(child: Text('排行榜功能开发中...')),
        ),
        actions: [
          TextButton(
            onPressed: () => Navigator.pop(context),
            child: const Text('关闭'),
          ),
        ],
      ),
    );
  }
  
  @override
  void dispose() {
    super.dispose();
  }
}

class _StatusItem extends StatelessWidget {
  final String label;
  final String value;
  final IconData icon;
  final Color color;
  
  const _StatusItem({
    required this.label,
    required this.value,
    required this.icon,
    required this.color,
  });
  
  @override
  Widget build(BuildContext context) {
    return Column(
      children: [
        Icon(icon, color: color, size: 28),
        const SizedBox(height: 4),
        Text(
          value,
          style: const TextStyle(
            fontSize: 18,
            fontWeight: FontWeight.bold,
          ),
        ),
        Text(
          label,
          style: const TextStyle(
            fontSize: 12,
            color: Colors.grey,
          ),
        ),
      ],
    );
  }
}

class _RewardChip extends StatelessWidget {
  final SignInReward reward;
  final bool isPreview;
  
  const _RewardChip({required this.reward, this.isPreview = false});
  
  @override
  Widget build(BuildContext context) {
    return Container(
      padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 6),
      decoration: BoxDecoration(
        color: isPreview ? Colors.grey[200] : Colors.blue[50],
        borderRadius: BorderRadius.circular(16),
        border: Border.all(
          color: isPreview ? Colors.grey[300]! : Colors.blue[200]!,
        ),
      ),
      child: Row(
        mainAxisSize: MainAxisSize.min,
        children: [
          Icon(
            _getIcon(),
            size: 16,
            color: isPreview ? Colors.grey : _getColor(),
          ),
          const SizedBox(width: 4),
          Text(
            '${reward.itemName} x${reward.count}',
            style: TextStyle(
              fontSize: 12,
              color: isPreview ? Colors.grey : Colors.black87,
            ),
          ),
          if (reward.isVipOnly) ...[
            const SizedBox(width: 4),
            const Icon(Icons.vpn_key, size: 12, color: Colors.amber),
          ],
        ],
      ),
    );
  }
  
  IconData _getIcon() {
    switch (reward.type) {
      case RewardType.gold:
        return Icons.monetization_on;
      case RewardType.gem:
        return Icons.diamond;
      case RewardType.item:
        return Icons.inventory_2;
      case RewardType.equipment:
        return Icons.shield;
    }
  }
  
  Color _getColor() {
    switch (reward.type) {
      case RewardType.gold:
        return Colors.amber;
      case RewardType.gem:
        return Colors.blue;
      case RewardType.item:
        return Colors.green;
      case RewardType.equipment:
        return Colors.purple;
    }
  }
}

class _CalendarDayCell extends StatelessWidget {
  final SignInDay day;
  
  const _CalendarDayCell({required this.day});
  
  @override
  Widget build(BuildContext context) {
    return Container(
      width: 36,
      height: 36,
      decoration: BoxDecoration(
        color: day.signedIn ? Colors.green : Colors.transparent,
        shape: BoxShape.circle,
        border: day.signedIn
            ? null
            : Border.all(color: Colors.grey[300]!),
      ),
      child: Center(
        child: Text(
          '${day.day}',
          style: TextStyle(
            color: day.signedIn ? Colors.white : Colors.black87,
            fontWeight: day.signedIn ? FontWeight.bold : FontWeight.normal,
          ),
        ),
      ),
    );
  }
}
