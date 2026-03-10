// ============================================
// Flutter 客户端示例 - 弹幕游戏
// 对接后端 WebSocket + JSON 协议
// ============================================

import 'dart:async';
import 'dart:convert';
import 'package:flutter/material.dart';
import 'package:web_socket_channel/web_socket_channel.dart';

/// 消息ID枚举
class MsgID {
  static const int csCreateRoom = 1001;
  static const int scCreateRoom = 1002;
  static const int csJoinRoom = 1003;
  static const int scJoinRoom = 1004;
  static const int csLeaveRoom = 1005;
  static const int scLeaveRoom = 1006;
  static const int csListRooms = 1007;
  static const int scListRooms = 1008;
  static const int csStartGame = 2001;
  static const int scStartGame = 2002;
  static const int csPlaceTower = 2003;
  static const int scPlaceTower = 2004;
  static const int csUpgradeTower = 2005;
  static const int scUpgradeTower = 2006;
  static const int csFrameInput = 3001;
  static const int scFrameSync = 3002;
  static const int csSendGift = 4001;
  static const int scReceiveGift = 4002;
  static const int csSendDanmaku = 4003;
  static const int scReceiveDanmaku = 4004;
  static const int csHeartbeat = 5001;
  static const int scHeartbeat = 5002;
  static const int scError = 9001;
}

/// 错误码
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
}

/// 基础消息结构
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

// ============ 房间相关模型 ============

class CreateRoomReq {
  final String playerId;
  final String playerName;
  final String mode;
  final int maxPlayers;

  CreateRoomReq({
    required this.playerId,
    required this.playerName,
    this.mode = 'classic',
    this.maxPlayers = 2,
  });

  Map<String, dynamic> toJson() => {
    'player_id': playerId,
    'player_name': playerName,
    'mode': mode,
    'max_players': maxPlayers,
  };
}

class CreateRoomResp {
  final String roomId;
  final String hostId;
  final String mode;
  final int maxPlayers;
  final int createdAt;

  CreateRoomResp.fromJson(Map<String, dynamic> json)
      : roomId = json['room_id'] as String,
        hostId = json['host_id'] as String,
        mode = json['mode'] as String,
        maxPlayers = json['max_players'] as int,
        createdAt = json['created_at'] as int;
}

class JoinRoomReq {
  final String roomId;
  final String playerId;
  final String playerName;

  JoinRoomReq({
    required this.roomId,
    required this.playerId,
    required this.playerName,
  });

  Map<String, dynamic> toJson() => {
    'room_id': roomId,
    'player_id': playerId,
    'player_name': playerName,
  };
}

class JoinRoomResp {
  final bool success;
  final String roomId;
  final List<String> playerIds;
  final List<String> playerNames;
  final String mode;

  JoinRoomResp.fromJson(Map<String, dynamic> json)
      : success = json['success'] as bool,
        roomId = json['room_id'] as String,
        playerIds = List<String>.from(json['player_ids'] ?? []),
        playerNames = List<String>.from(json['player_names'] ?? []),
        mode = json['mode'] as String;
}

class RoomInfo {
  final String roomId;
  final String hostId;
  final String mode;
  final int currentPlayers;
  final int maxPlayers;
  final String status;

  RoomInfo.fromJson(Map<String, dynamic> json)
      : roomId = json['room_id'] as String,
        hostId = json['host_id'] as String,
        mode = json['mode'] as String,
        currentPlayers = json['current_players'] as int,
        maxPlayers = json['max_players'] as int,
        status = json['status'] as String;
}

// ============ 游戏相关模型 ============

class PlaceTowerReq {
  final String roomId;
  final String playerId;
  final String towerType;
  final double x;
  final double y;

  PlaceTowerReq({
    required this.roomId,
    required this.playerId,
    required this.towerType,
    required this.x,
    required this.y,
  });

  Map<String, dynamic> toJson() => {
    'room_id': roomId,
    'player_id': playerId,
    'tower_type': towerType,
    'x': x,
    'y': y,
  };
}

class PlaceTowerResp {
  final bool success;
  final String towerId;
  final String towerType;
  final double x;
  final double y;
  final int level;
  final int cost;

  PlaceTowerResp.fromJson(Map<String, dynamic> json)
      : success = json['success'] as bool,
        towerId = json['tower_id'] as String,
        towerType = json['tower_type'] as String,
        x = (json['x'] as num).toDouble(),
        y = (json['y'] as num).toDouble(),
        level = json['level'] as int,
        cost = json['cost'] as int;
}

class GameStartResp {
  final String roomId;
  final int wave;
  final int startTime;
  final Map<String, int> playersMoney;

  GameStartResp.fromJson(Map<String, dynamic> json)
      : roomId = json['room_id'] as String,
        wave = json['wave'] as int,
        startTime = json['start_time'] as int,
        playersMoney = Map<String, int>.from(json['players_money'] ?? {});
}

// ============ 帧同步模型 ============

class FrameInput {
  final String playerId;
  final int frameId;
  final String action;
  final List<String> args;

  FrameInput({
    required this.playerId,
    required this.frameId,
    required this.action,
    this.args = const [],
  });

  Map<String, dynamic> toJson() => {
    'player_id': playerId,
    'frame_id': frameId,
    'action': action,
    'args': args,
  };
}

class Tower {
  final String id;
  final String ownerId;
  final String type;
  final int level;
  final double x;
  final double y;
  final double range;
  final double attackSpeed;
  final int damage;
  final int cost;

  Tower.fromJson(Map<String, dynamic> json)
      : id = json['id'] as String,
        ownerId = json['owner_id'] as String,
        type = json['type'] as String,
        level = json['level'] as int,
        x = (json['x'] as num).toDouble(),
        y = (json['y'] as num).toDouble(),
        range = (json['range'] as num).toDouble(),
        attackSpeed = (json['attack_speed'] as num).toDouble(),
        damage = json['damage'] as int,
        cost = json['cost'] as int;
}

class Enemy {
  final String id;
  final String type;
  final double x;
  final double y;
  final int hp;
  final int maxHp;
  final double speed;
  final double progress;
  final bool slowed;

  Enemy.fromJson(Map<String, dynamic> json)
      : id = json['id'] as String,
        type = json['type'] as String,
        x = (json['x'] as num).toDouble(),
        y = (json['y'] as num).toDouble(),
        hp = json['hp'] as int,
        maxHp = json['max_hp'] as int,
        speed = (json['speed'] as num).toDouble(),
        progress = (json['progress'] as num).toDouble(),
        slowed = json['slowed'] as bool? ?? false;
}

class Projectile {
  final String id;
  final String towerId;
  final String targetId;
  final double x;
  final double y;
  final double vx;
  final double vy;
  final int damage;
  final String effect;

  Projectile.fromJson(Map<String, dynamic> json)
      : id = json['id'] as String,
        towerId = json['tower_id'] as String,
        targetId = json['target_id'] as String,
        x = (json['x'] as num).toDouble(),
        y = (json['y'] as num).toDouble(),
        vx = (json['vx'] as num).toDouble(),
        vy = (json['vy'] as num).toDouble(),
        damage = json['damage'] as int,
        effect = json['effect'] as String;
}

class GameState {
  final int wave;
  final int timestamp;
  final Map<String, int> playersMoney;
  final List<Tower> towers;
  final List<Enemy> enemies;
  final List<Projectile> projectiles;

  GameState.fromJson(Map<String, dynamic> json)
      : wave = json['wave'] as int,
        timestamp = json['timestamp'] as int,
        playersMoney = Map<String, int>.from(json['players_money'] ?? {}),
        towers = (json['towers'] as List?)
            ?.map((e) => Tower.fromJson(e as Map<String, dynamic>))
            .toList() ?? [],
        enemies = (json['enemies'] as List?)
            ?.map((e) => Enemy.fromJson(e as Map<String, dynamic>))
            .toList() ?? [],
        projectiles = (json['projectiles'] as List?)
            ?.map((e) => Projectile.fromJson(e as Map<String, dynamic>))
            .toList() ?? [];
}

class FrameSync {
  final int frameId;
  final List<FrameInput> inputs;
  final GameState state;

  FrameSync.fromJson(Map<String, dynamic> json)
      : frameId = json['frame_id'] as int,
        inputs = (json['inputs'] as List?)
            ?.map((e) => FrameInput(
              playerId: e['player_id'] as String,
              frameId: e['frame_id'] as int,
              action: e['action'] as String,
              args: List<String>.from(e['args'] ?? []),
            ))
            .toList() ?? [],
        state = GameState.fromJson(json['state'] as Map<String, dynamic>);
}

// ============ 礼物弹幕模型 ============

class SendGiftReq {
  final String roomId;
  final String senderId;
  final String senderName;
  final String giftType;

  SendGiftReq({
    required this.roomId,
    required this.senderId,
    required this.senderName,
    required this.giftType,
  });

  Map<String, dynamic> toJson() => {
    'room_id': roomId,
    'sender_id': senderId,
    'sender_name': senderName,
    'gift_type': giftType,
  };
}

class ReceiveGiftResp {
  final String senderId;
  final String senderName;
  final String giftType;
  final int value;
  final int timestamp;

  ReceiveGiftResp.fromJson(Map<String, dynamic> json)
      : senderId = json['sender_id'] as String,
        senderName = json['sender_name'] as String,
        giftType = json['gift_type'] as String,
        value = json['value'] as int,
        timestamp = json['timestamp'] as int;
}

class SendDanmakuReq {
  final String roomId;
  final String senderId;
  final String senderName;
  final String content;

  SendDanmakuReq({
    required this.roomId,
    required this.senderId,
    required this.senderName,
    required this.content,
  });

  Map<String, dynamic> toJson() => {
    'room_id': roomId,
    'sender_id': senderId,
    'sender_name': senderName,
    'content': content,
  };
}

class ReceiveDanmakuResp {
  final String senderId;
  final String senderName;
  final String content;
  final String color;
  final int timestamp;

  ReceiveDanmakuResp.fromJson(Map<String, dynamic> json)
      : senderId = json['sender_id'] as String,
        senderName = json['sender_name'] as String,
        content = json['content'] as String,
        color = json['color'] as String? ?? 'white',
        timestamp = json['timestamp'] as int;
}

// ============ WebSocket 客户端 ============

class GameClient {
  final String serverUrl;
  final double heartbeatInterval;

  WebSocketChannel? _channel;
  Timer? _heartbeatTimer;
  final _messageController = StreamController<WSMessage>.broadcast();

  bool _isConnected = false;
  String? _playerId;
  String? _playerName;
  String? _currentRoomId;

  // Getters
  bool get isConnected => _isConnected;
  String? get playerId => _playerId;
  String? get playerName => _playerName;
  String? get currentRoomId => _currentRoomId;
  Stream<WSMessage> get messages => _messageController.stream;

  // 回调
  Function(CreateRoomResp)? onCreateRoom;
  Function(JoinRoomResp)? onJoinRoom;
  Function(List<RoomInfo>)? onListRooms;
  Function(GameStartResp)? onGameStart;
  Function(FrameSync)? onFrameSync;
  Function(ReceiveGiftResp)? onReceiveGift;
  Function(ReceiveDanmakuResp)? onReceiveDanmaku;
  Function(int, String)? onError;

  GameClient({
    this.serverUrl = 'ws://localhost:8080/ws',
    this.heartbeatInterval = 30,
  });

  Future<void> connect() async {
    try {
      _channel = WebSocketChannel.connect(Uri.parse(serverUrl));
      
      _channel!.stream.listen(
        (data) => _handleMessage(jsonDecode(data as String)),
        onError: (error) => print('WebSocket Error: $error'),
        onDone: () {
          _isConnected = false;
          _heartbeatTimer?.cancel();
        },
      );

      _isConnected = true;
      _startHeartbeat();
      print('Connected to server');
    } catch (e) {
      print('Connection failed: $e');
    }
  }

  void disconnect() {
    _heartbeatTimer?.cancel();
    _channel?.sink.close();
    _channel = null;
    _isConnected = false;
    _currentRoomId = null;
  }

  void _startHeartbeat() {
    _heartbeatTimer = Timer.periodic(
      Duration(seconds: heartbeatInterval.toInt()),
      (_) => sendMessage(MsgID.csHeartbeat, {}),
    );
  }

  void _handleMessage(Map<String, dynamic> json) {
    final msg = WSMessage.fromJson(json);
    _messageController.add(msg);

    switch (msg.msgId) {
      case MsgID.scCreateRoom:
        final resp = CreateRoomResp.fromJson(msg.data);
        _currentRoomId = resp.roomId;
        onCreateRoom?.call(resp);
        break;

      case MsgID.scJoinRoom:
        final resp = JoinRoomResp.fromJson(msg.data);
        if (resp.success) _currentRoomId = resp.roomId;
        onJoinRoom?.call(resp);
        break;

      case MsgID.scListRooms:
        final rooms = (msg.data['rooms'] as List?)
            ?.map((e) => RoomInfo.fromJson(e as Map<String, dynamic>))
            .toList() ?? [];
        onListRooms?.call(rooms);
        break;

      case MsgID.scStartGame:
        final resp = GameStartResp.fromJson(msg.data);
        onGameStart?.call(resp);
        break;

      case MsgID.scFrameSync:
        final sync = FrameSync.fromJson(msg.data);
        onFrameSync?.call(sync);
        break;

      case MsgID.scReceiveGift:
        final gift = ReceiveGiftResp.fromJson(msg.data);
        onReceiveGift?.call(gift);
        break;

      case MsgID.scReceiveDanmaku:
        final danmaku = ReceiveDanmakuResp.fromJson(msg.data);
        onReceiveDanmaku?.call(danmaku);
        break;

      case MsgID.scError:
        final code = msg.data['code'] as int;
        final message = msg.data['message'] as String;
        onError?.call(code, message);
        break;
    }
  }

  // ============ 发送消息 ============

  void sendMessage(int msgId, Map<String, dynamic> data) {
    if (!_isConnected || _channel == null) {
      print('Not connected');
      return;
    }

    final msg = {
      'msg_id': msgId,
      'data': data,
    };

    _channel!.sink.add(jsonEncode(msg));
  }

  // ============ 房间操作 API ============

  void createRoom(String playerId, String playerName, {String mode = 'classic'}) {
    _playerId = playerId;
    _playerName = playerName;

    sendMessage(MsgID.csCreateRoom, CreateRoomReq(
      playerId: playerId,
      playerName: playerName,
      mode: mode,
    ).toJson());
  }

  void joinRoom(String roomId, String playerId, String playerName) {
    _playerId = playerId;
    _playerName = playerName;

    sendMessage(MsgID.csJoinRoom, JoinRoomReq(
      roomId: roomId,
      playerId: playerId,
      playerName: playerName,
    ).toJson());
  }

  void leaveRoom() {
    if (_currentRoomId == null || _playerId == null) return;

    sendMessage(MsgID.csLeaveRoom, {
      'room_id': _currentRoomId,
      'player_id': _playerId,
    });

    _currentRoomId = null;
  }

  void listRooms() {
    sendMessage(MsgID.csListRooms, {});
  }

  // ============ 游戏操作 API ============

  void placeTower(String towerType, double x, double y) {
    if (_currentRoomId == null || _playerId == null) return;

    sendMessage(MsgID.csPlaceTower, PlaceTowerReq(
      roomId: _currentRoomId!,
      playerId: _playerId!,
      towerType: towerType,
      x: x,
      y: y,
    ).toJson());
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

    sendMessage(MsgID.csPlaceTower, {
      'room_id': _currentRoomId,
      'player_id': _playerId,
      'tower_id': towerId,
    });
  }

  // ============ 礼物弹幕 API ============

  void sendGift(String giftType) {
    if (_currentRoomId == null || _playerId == null || _playerName == null) return;

    sendMessage(MsgID.csSendGift, SendGiftReq(
      roomId: _currentRoomId!,
      senderId: _playerId!,
      senderName: _playerName!,
      giftType: giftType,
    ).toJson());
  }

  void sendDanmaku(String content) {
    if (_currentRoomId == null || _playerId == null || _playerName == null) return;

    sendMessage(MsgID.csSendDanmaku, SendDanmakuReq(
      roomId: _currentRoomId!,
      senderId: _playerId!,
      senderName: _playerName!,
      content: content,
    ).toJson());
  }

  void dispose() {
    disconnect();
    _messageController.close();
  }
}

// ============ Flutter 页面示例 ============

class GamePage extends StatefulWidget {
  @override
  _GamePageState createState() => _GamePageState();
}

class _GamePageState extends State<GamePage> {
  late GameClient _client;
  int _money = 0;
  int _wave = 0;
  List<Tower> _towers = [];
  List<Enemy> _enemies = [];
  final List<String> _danmakuList = [];

  @override
  void initState() {
    super.initState();
    _client = GameClient(serverUrl: 'ws://localhost:8080/ws');
    _setupCallbacks();
    _client.connect();
  }

  void _setupCallbacks() {
    _client.onGameStart = (resp) {
      setState(() {
        _wave = resp.wave;
        _money = resp.playersMoney[_client.playerId] ?? 0;
      });
    };

    _client.onFrameSync = (sync) {
      setState(() {
        _wave = sync.state.wave;
        _money = sync.state.playersMoney[_client.playerId!] ?? 0;
        _towers = sync.state.towers;
        _enemies = sync.state.enemies;
      });
    };

    _client.onReceiveGift = (gift) {
      print('收到礼物: ${gift.giftType} from ${gift.senderName}');
    };

    _client.onReceiveDanmaku = (danmaku) {
      setState(() {
        _danmakuList.add('${danmaku.senderName}: ${danmaku.content}');
        if (_danmakuList.length > 10) {
          _danmakuList.removeAt(0);
        }
      });
    };

    _client.onError = (code, message) {
      print('Error $code: $message');
    };
  }

  @override
  void dispose() {
    _client.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: Text('弹幕游戏'),
        actions: [
          Text('金币: $_money'),
          SizedBox(width: 16),
          Text('波次: $_wave'),
        ],
      ),
      body: Stack(
        children: [
          // 游戏画布
          CustomPaint(
            size: Size.infinite,
            painter: GamePainter(_towers, _enemies),
          ),

          // 弹幕显示
          Positioned(
            left: 0,
            bottom: 100,
            child: SizedBox(
              width: 300,
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: _danmakuList.map((d) => 
                  Text(d, style: TextStyle(color: Colors.white))
                ).toList(),
              ),
            ),
          ),
        ],
      ),
      bottomNavigationBar: BottomNavigationBar(
        items: [
          BottomNavigationBarItem(
            icon: Icon(Icons.add_box),
            label: '箭塔',
            onTap: () => _client.placeTower5, 5),
          ),
          BottomNavigationBarItem(
            icon:('arrow',  Icon(Icons.whatshot),
            label: '炮塔',
            onTap: () => _client.placeTower('cannon', 7, 5),
          ),
          BottomNavigationBarItem(
            icon: Icon(Icons.card_giftcard),
            label: '送礼',
            onTap: () => _client.sendGift('rocket'),
          ),
          BottomNavigationBarItem(
            icon: Icon(Icons.message),
            label: '弹幕',
            onTap: () => _client.sendDanmaku('666'),
          ),
        ],
      ),
    );
  }
}

// ============ 游戏画布绘制 ============

class GamePainter extends CustomPainter {
  final List<Tower> towers;
  final List<Enemy> enemies;

  GamePainter(this.towers, this.enemies);

  @override
  void paint(Canvas canvas, Size size) {
    // 绘制塔
    for (final tower in towers) {
      final paint = Paint()..color = Colors.blue;
      canvas.drawCircle(
        Offset(tower.x * 10, tower.y * 10),
        20,
        paint,
      );
    }

    // 绘制敌人
    for (final enemy in enemies) {
      final progress = (enemy.progress * size.width).clamp(0.0, size.width);
      final paint = Paint()..color = Colors.red;
      canvas.drawCircle(
        Offset(progress, enemy.y * 10),
        15,
        paint,
      );
    }
  }

  @override
  bool shouldRepaint(covariant CustomPainter oldDelegate) => true;
}
