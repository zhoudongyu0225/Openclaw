// ============================================
// Flutter 社交客户端 - 弹幕游戏
// 包含聊天、好友、公会、礼物、弹幕等功能
// ============================================

import 'dart:async';
import 'dart:convert';
import 'dart:math';
import 'package:flutter/material.dart';
import 'package:web_socket_channel/web_socket_channel.dart';

// ==================== 常量定义 ====================

class MsgID {
  // 房间相关
  static const int csCreateRoom = 1001;
  static const int scCreateRoom = 1002;
  static const int csJoinRoom = 1003;
  static const int scJoinRoom = 1004;
  static const int csLeaveRoom = 1005;
  static const int scLeaveRoom = 1006;
  static const int csListRooms = 1007;
  static const int scListRooms = 1008;

  // 游戏相关
  static const int csStartGame = 2001;
  static const int scStartGame = 2002;
  static const int csPlaceTower = 2003;
  static const int scPlaceTower = 2004;
  static const int csUpgradeTower = 2005;
  static const int scUpgradeTower = 2006;
  static const int csSellTower = 2007;
  static const int scSellTower = 2008;

  // 帧同步
  static const int csFrameInput = 3001;
  static const int scFrameSync = 3002;

  // 礼物弹幕
  static const int csSendGift = 4001;
  static const int scReceiveGift = 4002;
  static const int csSendDanmaku = 4003;
  static const int scReceiveDanmaku = 4004;

  // 聊天
  static const int csSendChat = 5001;
  static const int scReceiveChat = 5002;
  static const int csQueryHistory = 5003;
  static const int scQueryHistory = 5004;

  // 好友
  static const int csAddFriend = 6001;
  static const int scAddFriend = 6002;
  static const int csRemoveFriend = 6003;
  static const int scRemoveFriend = 6004;
  static const int csGetFriends = 6005;
  static const int scGetFriends = 6006;
  static const int csBlockPlayer = 6007;
  static const int scBlockPlayer = 6008;

  // 心跳
  static const int csHeartbeat = 9001;
  static const int scHeartbeat = 9002;

  // 错误
  static const int scError = 9999;
}

class ErrorCode {
  static const int success = 0;
  static const int invalidToken = 1001;
  static const int roomNotFound = 2001;
  static const int roomFull = 2002;
  static const int roomAlreadyStarted = 2003;
  static const int notRoomOwner = 2004;
  static const int invalidPosition = 3001;
  static const int notEnoughMoney = 3002;
  static const int towerNotFound = 3003;
  static const int friendNotFound = 4001;
  static const int friendAlreadyExists = 4002;
  static const int cannotBlockSelf = 4003;
}

class ChatChannel {
  static const int world = 0;
  static const int guild = 1;
  static const int private = 2;
  static const int system = 3;
  static const int battle = 4;
}

// ==================== 消息模型 ====================

class WSMessage {
  final int msgId;
  final Map<String, dynamic> data;

  WSMessage({required this.msgId, required this.data});

  factory WSMessage.fromJson(Map<String, dynamic> json) {
    return WSMessage(
      msgId: json['msg_id'] as int,
      data: json['data'] as Map<String, dynamic>? ?? {},
    );
  }

  Map<String, dynamic> toJson() => {
    'msg_id': msgId,
    'data': data,
  };
}

String encodeMessage(WSMessage msg) => jsonEncode(msg.toJson());

WSMessage decodeMessage(String data) {
  return WSMessage.fromJson(jsonDecode(data) as Map<String, dynamic>);
}

// ==================== 聊天消息 ====================

class ChatMessage {
  final String id;
  final int channel;
  final String senderId;
  final String senderName;
  final String content;
  final int timestamp;
  final String? targetId;
  final String? guildId;
  final bool isSystem;
  final Map<String, dynamic>? extra;

  ChatMessage({
    required this.id,
    required this.channel,
    required this.senderId,
    required this.senderName,
    required this.content,
    required this.timestamp,
    this.targetId,
    this.guildId,
    this.isSystem = false,
    this.extra,
  });

  factory ChatMessage.fromJson(Map<String, dynamic> json) {
    return ChatMessage(
      id: json['id'] as String? ?? '',
      channel: json['channel'] as int? ?? 0,
      senderId: json['sender_id'] as String? ?? '',
      senderName: json['sender_name'] as String? ?? '',
      content: json['content'] as String? ?? '',
      timestamp: json['timestamp'] as int? ?? 0,
      targetId: json['target_id'] as String?,
      guildId: json['guild_id'] as String?,
      isSystem: json['is_system'] as bool? ?? false,
      extra: json['extra'] as Map<String, dynamic>?,
    );
  }

  Map<String, dynamic> toJson() => {
    'id': id,
    'channel': channel,
    'sender_id': senderId,
    'sender_name': senderName,
    'content': content,
    'timestamp': timestamp,
    if (targetId != null) 'target_id': targetId,
    if (guildId != null) 'guild_id': guildId,
    'is_system': isSystem,
    if (extra != null) 'extra': extra,
  };
}

// ==================== 好友模型 ====================

class Friend {
  final String playerId;
  final String playerName;
  final int level;
  final int vipLevel;
  final bool online;
  final int lastSeen;
  final String? guildName;

  Friend({
    required this.playerId,
    required this.playerName,
    required this.level,
    this.vipLevel = 0,
    this.online = false,
    this.lastSeen = 0,
    this.guildName,
  });

  factory Friend.fromJson(Map<String, dynamic> json) {
    return Friend(
      playerId: json['player_id'] as String,
      playerName: json['player_name'] as String,
      level: json['level'] as int? ?? 1,
      vipLevel: json['vip_level'] as int? ?? 0,
      online: json['online'] as bool? ?? false,
      lastSeen: json['last_seen'] as int? ?? 0,
      guildName: json['guild_name'] as String?,
    );
  }

  Map<String, dynamic> toJson() => {
    'player_id': playerId,
    'player_name': playerName,
    'level': level,
    'vip_level': vipLevel,
    'online': online,
    'last_seen': lastSeen,
    if (guildName != null) 'guild_name': guildName,
  };
}

class FriendRequest {
  final String requestId;
  final String fromId;
  final String fromName;
  final int timestamp;
  final String status; // pending, accepted, rejected

  FriendRequest({
    required this.requestId,
    required this.fromId,
    required this.fromName,
    required this.timestamp,
    this.status = 'pending',
  });

  factory FriendRequest.fromJson(Map<String, dynamic> json) {
    return FriendRequest(
      requestId: json['request_id'] as String,
      fromId: json['from_id'] as String,
      fromName: json['from_name'] as String,
      timestamp: json['timestamp'] as int,
      status: json['status'] as String? ?? 'pending',
    );
  }
}

// ==================== 礼物模型 ====================

class Gift {
  final String type;
  final String name;
  final int value;
  final String icon;

  static const List<Gift> gifts = [
    Gift(type: 'coin', name: '金币', value: 1, icon: '🪙'),
    Gift(type: 'star', name: '星星', value: 10, icon: '⭐'),
    Gift(type: 'rocket', name: '火箭', value: 100, icon: '🚀'),
    Gift(type: 'car', name: '跑车', value: 500, icon: '🚗'),
    Gift(type: 'plane', name: '飞机', value: 1000, icon: '✈️'),
    Gift(type: 'bang', name: '棒', value: 2000, icon: '💥'),
  ];

  static Gift? getByType(String type) {
    try {
      return gifts.firstWhere((g) => g.type == type);
    } catch (_) {
      return null;
    }
  }
}

class GiftMessage {
  final String senderId;
  final String senderName;
  final String giftType;
  final int value;
  final int timestamp;

  GiftMessage({
    required this.senderId,
    required this.senderName,
    required this.giftType,
    required this.value,
    required this.timestamp,
  });

  factory GiftMessage.fromJson(Map<String, dynamic> json) {
    return GiftMessage(
      senderId: json['sender_id'] as String,
      senderName: json['sender_name'] as String,
      giftType: json['gift_type'] as String,
      value: json['value'] as int,
      timestamp: json['timestamp'] as int,
    );
  }
}

// ==================== 弹幕模型 ====================

class DanmakuMessage {
  final String senderId;
  final String senderName;
  final String content;
  final String color;
  final int timestamp;

  DanmakuMessage({
    required this.senderId,
    required this.senderName,
    required this.content,
    this.color = 'white',
    required this.timestamp,
  });

  factory DanmakuMessage.fromJson(Map<String, dynamic> json) {
    return DanmakuMessage(
      senderId: json['sender_id'] as String,
      senderName: json['sender_name'] as String,
      content: json['content'] as String,
      color: json['color'] as String? ?? 'white',
      timestamp: json['timestamp'] as int,
    );
  }

  Color get colorValue {
    switch (color) {
      case 'red':
        return Colors.red;
      case 'gold':
        return Colors.amber;
      case 'rainbow':
        return Colors.purple;
      default:
        return Colors.white;
    }
  }
}

// ==================== 房间模型 ====================

class RoomInfo {
  final String roomId;
  final String hostId;
  final String mode;
  final int currentPlayers;
  final int maxPlayers;
  final String status;
  final String? title;

  RoomInfo({
    required this.roomId,
    required this.hostId,
    required this.mode,
    required this.currentPlayers,
    required this.maxPlayers,
    required this.status,
    this.title,
  });

  factory RoomInfo.fromJson(Map<String, dynamic> json) {
    return RoomInfo(
      roomId: json['room_id'] as String,
      hostId: json['host_id'] as String,
      mode: json['mode'] as String,
      currentPlayers: json['current_players'] as int,
      maxPlayers: json['max_players'] as int,
      status: json['status'] as String,
      title: json['title'] as String?,
    );
  }
}

// ==================== 社交客户端 ====================

class SocialClient {
  final String serverUrl;
  final int heartbeatInterval;
  final int reconnectDelay;
  final int maxReconnectAttempts;

  WebSocketChannel? _channel;
  Timer? _heartbeatTimer;
  Timer? _reconnectTimer;
  int _reconnectAttempts = 0;

  bool _isConnected = false;
  bool _isConnecting = false;
  String? _playerId;
  String? _playerName;
  String? _currentRoomId;
  int _guildId;

  // 消息流
  final _messageController = StreamController<WSMessage>.broadcast();
  final _chatController = StreamController<ChatMessage>.broadcast();
  final _giftController = StreamController<GiftMessage>.broadcast();
  final _danmakuController = StreamController<DanmakuMessage>.broadcast();
  final _friendController = StreamController<List<Friend>>.broadcast();
  final _connectionController = StreamController<bool>.broadcast();

  // 缓存
  List<ChatMessage> _chatHistory = [];
  List<Friend> _friends = [];
  List<RoomInfo> _rooms = [];

  // Getters
  bool get isConnected => _isConnected;
  bool get isConnecting => _isConnecting;
  String? get playerId => _playerId;
  String? get playerName => _playerName;
  String? get currentRoomId => _currentRoomId;
  int get guildId => _guildId;
  List<ChatMessage> get chatHistory => _chatHistory;
  List<Friend> get friends => _friends;
  List<RoomInfo> get rooms => _rooms;

  // Streams
  Stream<WSMessage> get messages => _messageController.stream;
  Stream<ChatMessage> get chatMessages => _chatController.stream;
  Stream<GiftMessage> get gifts => _giftController.stream;
  Stream<DanmakuMessage> get danmaku => _danmakuController.stream;
  Stream<List<Friend>> get friendList => _friendController.stream;
  Stream<bool> get connectionStatus => _connectionController.stream;

  // Callbacks
  Function(int, String)? onError;
  Function()? onConnect;
  Function()? onDisconnect;
  Function(RoomInfo)? onRoomCreated;
  Function(RoomInfo)? onRoomJoined;
  Function()? onGameStart;

  SocialClient({
    this.serverUrl = 'ws://localhost:8080/ws',
    this.heartbeatInterval = 30,
    this.reconnectDelay = 3000,
    this.maxReconnectAttempts = 5,
    this._guildId = 0,
  });

  // ==================== 连接管理 ====================

  Future<bool> connect() async {
    if (_isConnected || _isConnecting) {
      return _isConnected;
    }

    _isConnecting = true;

    try {
      _channel = WebSocketChannel.connect(Uri.parse(serverUrl));

      await _channel!.ready;

      _channel!.stream.listen(
        _onMessage,
        onError: _onError,
        onDone: _onDone,
      );

      _isConnected = true;
      _isConnecting = false;
      _reconnectAttempts = 0;
      _startHeartbeat();

      _connectionController.add(true);
      onConnect?.call();

      return true;
    } catch (e) {
      _isConnecting = false;
      _scheduleReconnect();
      return false;
    }
  }

  void disconnect() {
    _heartbeatTimer?.cancel();
    _reconnectTimer?.cancel();
    _channel?.sink.close();
    _channel = null;
    _isConnected = false;
    _connectionController.add(false);
  }

  void _onMessage(dynamic data) {
    try {
      final msg = decodeMessage(data as String);
      _messageController.add(msg);
      _dispatchMessage(msg);
    } catch (e) {
      print('Message parse error: $e');
    }
  }

  void _onError(dynamic error) {
    print('WebSocket error: $error');
    _handleDisconnect();
  }

  void _onDone() {
    _handleDisconnect();
  }

  void _handleDisconnect() {
    if (_isConnected) {
      _isConnected = false;
      _connectionController.add(false);
      onDisconnect?.call();
      _scheduleReconnect();
    }
  }

  void _scheduleReconnect() {
    if (_reconnectAttempts >= maxReconnectAttempts) {
      print('Max reconnect attempts reached');
      return;
    }

    _reconnectTimer?.cancel();
    _reconnectTimer = Timer(Duration(milliseconds: reconnectDelay), () {
      _reconnectAttempts++;
      print('Reconnecting... attempt $_reconnectAttempts');
      connect();
    });
  }

  void _startHeartbeat() {
    _heartbeatTimer?.cancel();
    _heartbeatTimer = Timer.periodic(
      Duration(seconds: heartbeatInterval),
      (_) => sendMessage(MsgID.csHeartbeat, {}),
    );
  }

  // ==================== 消息分发 ====================

  void _dispatchMessage(WSMessage msg) {
    switch (msg.msgId) {
      case MsgID.scReceiveChat:
        final chat = ChatMessage.fromJson(msg.data);
        _chatHistory.add(chat);
        if (_chatHistory.length > 100) {
          _chatHistory.removeAt(0);
        }
        _chatController.add(chat);
        break;

      case MsgID.scReceiveGift:
        final gift = GiftMessage.fromJson(msg.data);
        _giftController.add(gift);
        break;

      case MsgID.scReceiveDanmaku:
        final danmaku = DanmakuMessage.fromJson(msg.data);
        _danmakuController.add(danmaku);
        break;

      case MsgID.scGetFriends:
        _friends = (msg.data['friends'] as List?)
            ?.map((e) => Friend.fromJson(e as Map<String, dynamic>))
            .toList() ?? [];
        _friendController.add(_friends);
        break;

      case MsgID.scListRooms:
        _rooms = (msg.data['rooms'] as List?)
            ?.map((e) => RoomInfo.fromJson(e as Map<String, dynamic>))
            .toList() ?? [];
        break;

      case MsgID.scCreateRoom:
        final room = RoomInfo.fromJson(msg.data);
        _currentRoomId = room.roomId;
        onRoomCreated?.call(room);
        break;

      case MsgID.scJoinRoom:
        if (msg.data['success'] == true) {
          final room = RoomInfo.fromJson(msg.data);
          _currentRoomId = room.roomId;
          onRoomJoined?.call(room);
        }
        break;

      case MsgID.scStartGame:
        onGameStart?.call();
        break;

      case MsgID.scHeartbeat:
        // Heartbeat ack
        break;

      case MsgID.scError:
        final code = msg.data['code'] as int? ?? 0;
        final message = msg.data['message'] as String? ?? 'Unknown error';
        onError?.call(code, message);
        break;
    }
  }

  // ==================== 消息发送 ====================

  void sendMessage(int msgId, Map<String, dynamic> data) {
    if (!_isConnected || _channel == null) {
      print('Not connected, message not sent: $msgId');
      return;
    }

    final msg = WSMessage(msgId: msgId, data: data);
    _channel!.sink.add(encodeMessage(msg));
  }

  // ==================== 登录/注册 ====================

  void login(String playerId, String playerName) {
    _playerId = playerId;
    _playerName = playerName;

    sendMessage(MsgID.csHeartbeat, {
      'player_id': playerId,
      'player_name': playerName,
    });
  }

  // ==================== 聊天功能 ====================

  void sendChat(int channel, String content, {String? targetId, int? guildId}) {
    if (_playerId == null || _playerName == null) return;

    sendMessage(MsgID.csSendChat, {
      'channel': channel,
      'sender_id': _playerId,
      'sender_name': _playerName,
      'content': content,
      if (targetId != null) 'target_id': targetId,
      if (guildId != null) 'guild_id': guildId,
    });
  }

  void sendWorldChat(String content) {
    sendChat(ChatChannel.world, content);
  }

  void sendGuildChat(String content) {
    sendChat(ChatChannel.guild, content, guildId: _guildId);
  }

  void sendPrivateChat(String targetId, String content) {
    sendChat(ChatChannel.private, content, targetId: targetId);
  }

  void queryChatHistory(int channel, {int limit = 50}) {
    sendMessage(MsgID.csQueryHistory, {
      'channel': channel,
      'limit': limit,
    });
  }

  // ==================== 好友功能 ====================

  void addFriend(String playerId) {
    if (_playerId == null) return;

    sendMessage(MsgID.csAddFriend, {
      'player_id': _playerId,
      'target_id': playerId,
    });
  }

  void removeFriend(String playerId) {
    if (_playerId == null) return;

    sendMessage(MsgID.csRemoveFriend, {
      'player_id': _playerId,
      'target_id': playerId,
    });
  }

  void getFriends() {
    if (_playerId == null) return;

    sendMessage(MsgID.csGetFriends, {
      'player_id': _playerId,
    });
  }

  void blockPlayer(String playerId) {
    if (_playerId == null) return;

    sendMessage(MsgID.csBlockPlayer, {
      'player_id': _playerId,
      'target_id': playerId,
    });
  }

  // ==================== 礼物功能 ====================

  void sendGift(String giftType) {
    if (_currentRoomId == null || _playerId == null || _playerName == null) return;

    sendMessage(MsgID.csSendGift, {
      'room_id': _currentRoomId,
      'sender_id': _playerId,
      'sender_name': _playerName,
      'gift_type': giftType,
    });
  }

  // ==================== 弹幕功能 ====================

  void sendDanmaku(String content, {String color = 'white'}) {
    if (_currentRoomId == null || _playerId == null || _playerName == null) return;

    sendMessage(MsgID.csSendDanmaku, {
      'room_id': _currentRoomId,
      'sender_id': _playerId,
      'sender_name': _playerName,
      'content': content,
      'color': color,
    });
  }

  // ==================== 房间功能 ====================

  void createRoom({String mode = 'classic', String? title, int maxPlayers = 4}) {
    if (_playerId == null || _playerName == null) return;

    sendMessage(MsgID.csCreateRoom, {
      'player_id': _playerId,
      'player_name': _playerName,
      'mode': mode,
      if (title != null) 'title': title,
      'max_players': maxPlayers,
    });
  }

  void joinRoom(String roomId) {
    if (_playerId == null || _playerName == null) return;

    sendMessage(MsgID.csJoinRoom, {
      'room_id': roomId,
      'player_id': _playerId,
      'player_name': _playerName,
    });
  }

  void leaveRoom() {
    if (_currentRoomId == null || _playerId == null) return;

    sendMessage(MsgID.csLeaveRoom, {
      'room_id': _currentRoomId,
      'player_id': _playerId,
    });

    _currentRoomId = null;
  }

  void listRooms({String? mode}) {
    sendMessage(MsgID.csListRooms, {
      if (mode != null) 'mode': mode,
    });
  }

  void startGame() {
    if (_currentRoomId == null || _playerId == null) return;

    sendMessage(MsgID.csStartGame, {
      'room_id': _currentRoomId,
      'player_id': _playerId,
    });
  }

  // ==================== 游戏操作 ====================

  void placeTower(String towerType, double x, double y) {
    if (_currentRoomId == null || _playerId == null) return;

    sendMessage(MsgID.csPlaceTower, {
      'room_id': _currentRoomId,
      'player_id': _playerId,
      'tower_type': towerType,
      'x': x,
      'y': y,
    });
  }

  void upgradeTower(String towerId) {
    if (_currentRoomId == null || _playerId == null) return;

    sendMessage(MsgID.csUpgradeTower, {
      'room_id': _currentRoomId,
      'player_id': _playerId,
      'tower_id': towerId,
    });
  }

  void sellTower(String towerId) {
    if (_currentRoomId == null || _playerId == null) return;

    sendMessage(MsgID.csSellTower, {
      'room_id': _currentRoomId,
      'player_id': _playerId,
      'tower_id': towerId,
    });
  }

  // ==================== 清理 ====================

  void dispose() {
    disconnect();
    _messageController.close();
    _chatController.close();
    _giftController.close();
    _danmakuController.close();
    _friendController.close();
    _connectionController.close();
  }
}

// ==================== Flutter UI 组件 ====================

// 聊天窗口组件
class ChatWindow extends StatefulWidget {
  final SocialClient client;
  final int channel;
  final bool showChannelSelector;

  const ChatWindow({
    Key? key,
    required this.client,
    this.channel = ChatChannel.world,
    this.showChannelSelector = true,
  }) : super(key: key);

  @override
  _ChatWindowState createState() => _ChatWindowState();
}

class _ChatWindowState extends State<ChatWindow> {
  final _messageController = TextEditingController();
  final _scrollController = ScrollController();
  List<ChatMessage> _messages = [];
  int _currentChannel = ChatChannel.world;

  @override
  void initState() {
    super.initState();
    _currentChannel = widget.channel;
    _loadHistory();
    _subscribeToMessages();
  }

  void _loadHistory() {
    widget.client.queryChatHistory(_currentChannel);
    setState(() {
      _messages = widget.client.chatHistory
          .where((m) => m.channel == _currentChannel)
          .toList();
    });
  }

  void _subscribeToMessages() {
    widget.client.chatMessages.listen((msg) {
      if (msg.channel == _currentChannel) {
        setState(() {
          _messages.add(msg);
          if (_messages.length > 100) {
            _messages.removeAt(0);
          }
        });
        _scrollToBottom();
      }
    });
  }

  void _scrollToBottom() {
    WidgetsBinding.instance.addPostFrameCallback((_) {
      if (_scrollController.hasClients) {
        _scrollController.animateTo(
          _scrollController.position.maxScrollExtent,
          duration: Duration(milliseconds: 200),
          curve: Curves.easeOut,
        );
      }
    });
  }

  void _sendMessage() {
    final content = _messageController.text.trim();
    if (content.isEmpty) return;

    switch (_currentChannel) {
      case ChatChannel.world:
        widget.client.sendWorldChat(content);
        break;
      case ChatChannel.guild:
        widget.client.sendGuildChat(content);
        break;
    }

    _messageController.clear();
  }

  @override
  Widget build(BuildContext context) {
    return Column(
      children: [
        if (widget.showChannelSelector)
          _buildChannelSelector(),
        Expanded(
          child: ListView.builder(
            controller: _scrollController,
            itemCount: _messages.length,
            itemBuilder: (context, index) {
              final msg = _messages[index];
              return _buildMessage(msg);
            },
          ),
        ),
        _buildInput(),
      ],
    );
  }

  Widget _buildChannelSelector() {
    return Container(
      padding: EdgeInsets.symmetric(horizontal: 8, vertical: 4),
      child: Row(
        children: [
          _buildChannelButton('世界', ChatChannel.world),
          _buildChannelButton('公会', ChatChannel.guild),
          _buildChannelButton('私聊', ChatChannel.private),
        ],
      ),
    );
  }

  Widget _buildChannelButton(String label, int channel) {
    final isSelected = _currentChannel == channel;
    return Padding(
      padding: EdgeInsets.only(right: 8),
      child: ElevatedButton(
        onPressed: () {
          setState(() {
            _currentChannel = channel;
          });
          _loadHistory();
        },
        style: ElevatedButton.styleFrom(
          backgroundColor: isSelected ? Colors.blue : Colors.grey,
        ),
        child: Text(label),
      ),
    );
  }

  Widget _buildMessage(ChatMessage msg) {
    final isSelf = msg.senderId == widget.client.playerId;
    return Container(
      padding: EdgeInsets.symmetric(horizontal: 8, vertical: 4),
      alignment: isSelf ? Alignment.centerRight : Alignment.centerLeft,
      child: Container(
        constraints: BoxConstraints(maxWidth: 280),
        padding: EdgeInsets.all(8),
        decoration: BoxDecoration(
          color: isSelf ? Colors.blue[100] : Colors.grey[200],
          borderRadius: BorderRadius.circular(8),
        ),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Text(
              msg.senderName,
              style: TextStyle(fontWeight: FontWeight.bold, fontSize: 12),
            ),
            Text(msg.content),
          ],
        ),
      ),
    );
  }

  Widget _buildInput() {
    return Container(
      padding: EdgeInsets.all(8),
      child: Row(
        children: [
          Expanded(
            child: TextField(
              controller: _messageController,
              decoration: InputDecoration(
                hintText: '输入消息...',
                border: OutlineInputBorder(),
                contentPadding: EdgeInsets.symmetric(horizontal: 12, vertical: 8),
              ),
              onSubmitted: (_) => _sendMessage(),
            ),
          ),
          SizedBox(width: 8),
          IconButton(
            icon: Icon(Icons.send),
            onPressed: _sendMessage,
          ),
        ],
      ),
    );
  }

  @override
  void dispose() {
    _messageController.dispose();
    _scrollController.dispose();
    super.dispose();
  }
}

// 好友列表组件
class FriendListWidget extends StatefulWidget {
  final SocialClient client;

  const FriendListWidget({Key? key, required this.client}) : super(key: key);

  @override
  _FriendListWidgetState createState() => _FriendListWidgetState();
}

class _FriendListWidgetState extends State<FriendListWidget> {
  List<Friend> _friends = [];
  List<Friend> _onlineFriends = [];
  List<Friend> _offlineFriends = [];

  @override
  void initState() {
    super.initState();
    _loadFriends();
    widget.client.friendList.listen((friends) {
      setState(() {
        _friends = friends;
        _onlineFriends = friends.where((f) => f.online).toList();
        _offlineFriends = friends.where((f) => !f.online).toList();
      });
    });
  }

  void _loadFriends() {
    widget.client.getFriends();
  }

  @override
  Widget build(BuildContext context) {
    return ListView(
      children: [
        if (_onlineFriends.isNotEmpty) ...[
          _buildSectionHeader('在线 (${_onlineFriends.length})'),
          ..._onlineFriends.map((f) => _buildFriendItem(f)),
        ],
        if (_offlineFriends.isNotEmpty) ...[
          _buildSectionHeader('离线 (${_offlineFriends.length})'),
          ..._offlineFriends.map((f) => _buildFriendItem(f)),
        ],
        if (_friends.isEmpty)
          Center(
            child: Padding(
              padding: EdgeInsets.all(16),
              child: Text('暂无好友'),
            ),
          ),
      ],
    );
  }

  Widget _buildSectionHeader(String title) {
    return Container(
      padding: EdgeInsets.symmetric(horizontal: 16, vertical: 8),
      color: Colors.grey[200],
      child: Text(
        title,
        style: TextStyle(fontWeight: FontWeight.bold),
      ),
    );
  }

  Widget _buildFriendItem(Friend friend) {
    return ListTile(
      leading: CircleAvatar(
        backgroundColor: friend.online ? Colors.green : Colors.grey,
        child: Text(friend.playerName[0]),
      ),
      title: Text(friend.playerName),
      subtitle: Text('Lv.${friend.level}'),
      trailing: Row(
        mainAxisSize: MainAxisSize.min,
        children: [
          if (friend.guildName != null)
            Text(friend.guildName!, style: TextStyle(fontSize: 12)),
          IconButton(
            icon: Icon(Icons.chat),
            onPressed: () {
              // 打开私聊
            },
          ),
        ],
      ),
      onTap: () {
        // 显示好友详情
      },
    );
  }
}

// 礼物选择器组件
class GiftPicker extends StatelessWidget {
  final Function(String) onGiftSelected;

  const GiftPicker({Key? key, required this.onGiftSelected}) : super(key: key);

  @override
  Widget build(BuildContext context) {
    return Container(
      padding: EdgeInsets.all(16),
      child: Column(
        mainAxisSize: MainAxisSize.min,
        children: [
          Text('选择礼物', style: TextStyle(fontSize: 18, fontWeight: FontWeight.bold)),
          SizedBox(height: 16),
          Wrap(
            spacing: 16,
            runSpacing: 16,
            children: Gift.gifts.map((gift) {
              return InkWell(
                onTap: () {
                  onGiftSelected(gift.type);
                  Navigator.pop(context);
                },
                child: Container(
                  padding: EdgeInsets.all(12),
                  decoration: BoxDecoration(
                    border: Border.all(color: Colors.grey[300]!),
                    borderRadius: BorderRadius.circular(8),
                  ),
                  child: Column(
                    children: [
                      Text(gift.icon, style: TextStyle(fontSize: 32)),
                      SizedBox(height: 4),
                      Text(gift.name),
                      Text('${gift.value}', style: TextStyle(color: Colors.orange)),
                    ],
                  ),
                ),
              );
            }).toList(),
          ),
        ],
      ),
    );
  }
}

// 弹幕渲染器组件
class DanmakuRenderer extends StatefulWidget {
  final Stream<DanmakuMessage> danmakuStream;
  final int maxMessages;

  const DanmakuRenderer({
    Key? key,
    required this.danmakuStream,
    this.maxMessages = 10,
  }) : super(key: key);

  @override
  _DanmakuRendererState createState() => _DanmakuRendererState();
}

class _DanmakuRendererState extends State<DanmakuRenderer> {
  List<DanmakuMessage> _messages = [];

  @override
  void initState() {
    super.initState();
    widget.danmakuStream.listen((msg) {
      setState(() {
        _messages.add(msg);
        if (_messages.length > widget.maxMessages) {
          _messages.removeAt(0);
        }
      });
    });
  }

  @override
  Widget build(BuildContext context) {
    return Positioned(
      left: 0,
      bottom: 100,
      child: SizedBox(
        width: 300,
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: _messages.map((msg) {
            return Text(
              '${msg.senderName}: ${msg.content}',
              style: TextStyle(
                color: msg.colorValue,
                fontSize: 16,
                fontWeight: FontWeight.bold,
                shadows: [
                  Shadow(
                    color: Colors.black,
                    blurRadius: 2,
                  ),
                ],
              ),
            );
          }).toList(),
        ),
      ),
    );
  }
}

// ==================== 工具函数 ====================

String generateId() {
  return DateTime.now().millisecondsSinceEpoch.toString() +
      Random().nextInt(10000).toString();
}
