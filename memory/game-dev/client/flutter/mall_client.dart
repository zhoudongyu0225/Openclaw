// ============================================
// 商城客户端 - Flutter UI
// ============================================

import 'dart:async';
import 'dart:convert';
import 'package:flutter/material.dart';
import 'package:web_socket_channel/web_socket_channel.dart';

/// 商城客户端
class MallClient {
  WebSocketChannel? _channel;
  final String serverUrl;
  final String token;
  
  final _itemsController = StreamController<List<MallItem>>.broadcast();
  final _balanceController = StreamController<Map<String, int>>.broadcast();
  
  Stream<List<MallItem>> get items => _itemsController.stream;
  Stream<Map<String, int>> get balance => _balanceController.stream;
  
  MallClient({required this.serverUrl, required this.token});
  
  void connect() {
    _channel = WebSocketChannel.connect(
      Uri.parse('$serverUrl?token=$token&module=mall'),
    );
    _channel!.stream.listen(_onMessage);
  }
  
  void _onMessage(dynamic data) {
    final json = jsonDecode(data);
    final msgId = json['msg_id'];
    final jsonData = json['data'];
    
    switch (msgId) {
      case 20001: // SCMallItems
        final items = (jsonData['items'] as List)
            .map((i) => MallItem.fromJson(i))
            .toList();
        _itemsController.add(items);
        break;
      case 20002: // SCPurchaseResult
        final balance = Map<String, int>.from(jsonData['balance']);
        _balanceController.add(balance);
        break;
      case 20003: // SCBalance
        final balance = Map<String, int>.from(jsonData);
        _balanceController.add(balance);
        break;
    }
  }
  
  /// 获取商城商品
  void getMallItems(MallType type) {
    _channel?.sink.add(jsonEncode({
      'msg_id': 20001,
      'data': {'type': type.name},
    }));
  }
  
  /// 购买商品
  void purchase(String itemId, int quantity) {
    _channel?.sink.add(jsonEncode({
      'msg_id': 20002,
      'data': {'item_id': itemId, 'quantity': quantity},
    }));
  }
  
  /// 获取余额
  void getBalance() {
    _channel?.sink.add(jsonEncode({
      'msg_id': 20003,
      'data': {},
    }));
  }
  
  /// 获取今日消费
  void getTodaySpend() {
    _channel?.sink.add(jsonEncode({
      'msg_id': 20004,
      'data': {},
    }));
  }
  
  /// 获取购买历史
  void getPurchaseHistory({int page = 1, int pageSize = 20}) {
    _channel?.sink.add(jsonEncode({
      'msg_id': 20005,
      'data': {'page': page, 'page_size': pageSize},
    }));
  }
  
  void dispose() {
    _channel?.sink.close();
    _itemsController.close();
    _balanceController.close();
  }
}

/// 商城类型
enum MallType {
  gift,     // 礼物商城
  item,     // 道具商城
  skin,     // 皮肤商城
  random,   // 随机商城
  honor,    // 荣誉商城
}

/// 货币类型
enum CurrencyType {
  gold,   // 金币
  gem,    // 钻石
  honor,  // 荣誉点
  credit, // 积分
}

/// 商城商品
class MallItem {
  final String id;
  final String name;
  final String description;
  final MallType type;
  final int price;
  final CurrencyType currency;
  final int stock;
  final int? discount;
  final int? purchaseLimit;
  final int? minLevel;
  final int? vipLevel;
  final PurchaseLimitType? limitType;
  final String? imageUrl;
  final List<String> tags;
  final int startTime;
  final int endTime;
  
  MallItem({
    required this.id,
    required this.name,
    required this.description,
    required this.type,
    required this.price,
    required this.currency,
    required this.stock,
    this.discount,
    this.purchaseLimit,
    this.minLevel,
    this.vipLevel,
    this.limitType,
    this.imageUrl,
    required this.tags,
    required this.startTime,
    required this.endTime,
  });
  
  factory MallItem.fromJson(Map<String, dynamic> json) {
    return MallItem(
      id: json['item_id'],
      name: json['name'],
      description: json['description'] ?? '',
      type: MallType.values.firstWhere(
        (e) => e.name == json['type'],
        orElse: () => MallType.item,
      ),
      price: json['price'],
      currency: CurrencyType.values.firstWhere(
        (e) => e.name == json['currency'],
        orElse: () => CurrencyType.gold,
      ),
      stock: json['stock'] ?? -1,
      discount: json['discount'],
      purchaseLimit: json['purchase_limit'],
      minLevel: json['min_level'],
      vipLevel: json['vip_level'],
      limitType: json['limit_type'] != null
          ? PurchaseLimitType.values.firstWhere(
              (e) => e.name == json['limit_type'],
              orElse: () => PurchaseLimitType.none,
            )
          : null,
      imageUrl: json['image_url'],
      tags: List<String>.from(json['tags'] ?? []),
      startTime: json['start_time'] ?? 0,
      endTime: json['end_time'] ?? 0,
    );
  }
  
  bool get isOnSale {
    final now = DateTime.now().millisecondsSinceEpoch;
    return (startTime == 0 || now >= startTime) && 
           (endTime == 0 || now <= endTime);
  }
  
  bool get isOutOfStock => stock == 0;
  
  int get discountedPrice => discount != null 
      ? (price * discount! / 100).round() 
      : price;
  
  bool get hasDiscount => discount != null && discount! < 100;
}

/// 购买限制类型
enum PurchaseLimitType {
  none,
  daily,
  weekly,
  monthly,
}

/// 购买历史
class PurchaseRecord {
  final String id;
  final String itemId;
  final String itemName;
  final int quantity;
  final int totalPrice;
  final CurrencyType currency;
  final int purchasedAt;
  
  PurchaseRecord({
    required this.id,
    required this.itemId,
    required this.itemName,
    required this.quantity,
    required this.totalPrice,
    required this.currency,
    required this.purchasedAt,
  });
  
  factory PurchaseRecord.fromJson(Map<String, dynamic> json) {
    return PurchaseRecord(
      id: json['record_id'],
      itemId: json['item_id'],
      itemName: json['item_name'],
      quantity: json['quantity'],
      totalPrice: json['total_price'],
      currency: CurrencyType.values.firstWhere(
        (e) => e.name == json['currency'],
        orElse: () => CurrencyType.gold,
      ),
      purchasedAt: json['purchased_at'],
    );
  }
}

// ============ 商城页面 ============

class MallPage extends StatefulWidget {
  final MallClient mallClient;
  
  const MallPage({super.key, required this.mallClient});
  
  @override
  State<MallPage> createState() => _MallPageState();
}

class _MallPageState extends State<MallPage> with SingleTickerProviderStateMixin {
  late TabController _tabController;
  MallType _currentType = MallType.gift;
  List<MallItem> _items = [];
  Map<String, int> _balance = {};
  bool _isLoading = false;
  
  final _typeNames = {
    MallType.gift: '礼物',
    MallType.item: '道具',
    MallType.skin: '皮肤',
    MallType.random: '随机',
    MallType.honor: '荣誉',
  };
  
  @override
  void initState() {
    super.initState();
    _tabController = TabController(length: MallType.values.length, vsync: this);
    _tabController.addListener(() {
      setState(() {
        _currentType = MallType.values[_tabController.index];
      });
      _loadItems();
    });
    _loadBalance();
    _loadItems();
  }
  
  void _loadBalance() {
    widget.mallClient.getBalance();
    widget.mallClient.balance.listen((balance) {
      setState(() => _balance = balance);
    });
  }
  
  void _loadItems() {
    setState(() => _isLoading = true);
    widget.mallClient.getMallItems(_currentType);
    
    widget.mallClient.items.listen((items) {
      setState(() {
        _items = items;
        _isLoading = false;
      });
    });
  }
  
  void _purchase(MallItem item) {
    showDialog(
      context: context,
      builder: (context) => AlertDialog(
        title: Text('购买 ${item.name}'),
        content: Column(
          mainAxisSize: MainAxisSize.min,
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Text(item.description),
            const SizedBox(height: 16),
            Text('价格: ${item.hasDiscount ? item.discountedPrice : item.price} ${item.currency.name}'),
            if (item.hasDiscount)
              Text(
                '原价: ${item.price}',
                style: const TextStyle(
                  decoration: TextDecoration.lineThrough,
                  color: Colors.grey,
                ),
              ),
          ],
        ),
        actions: [
          TextButton(
            onPressed: () => Navigator.pop(context),
            child: const Text('取消'),
          ),
          ElevatedButton(
            onPressed: () {
              widget.mallClient.purchase(item.id, 1);
              Navigator.pop(context);
            },
            child: const Text('购买'),
          ),
        ],
      ),
    );
  }
  
  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text('商城'),
        bottom: TabBar(
          controller: _tabController,
          isScrollable: true,
          tabs: MallType.values.map((t) => Tab(text: _typeNames[t])).toList(),
        ),
      ),
      body: Column(
        children: [
          _buildBalanceBar(),
          Expanded(
            child: _isLoading
                ? const Center(child: CircularProgressIndicator())
                : _buildItemGrid(),
          ),
        ],
      ),
    );
  }
  
  Widget _buildBalanceBar() {
    return Container(
      padding: const EdgeInsets.all(12),
      color: Colors.grey[100],
      child: Row(
        mainAxisAlignment: MainAxisAlignment.spaceAround,
        children: [
          _BalanceItem(
            icon: Icons.monetization_on,
            iconColor: Colors.amber,
            name: '金币',
            value: _balance['gold'] ?? 0,
          ),
          _BalanceItem(
            icon: Icons.diamond,
            iconColor: Colors.blue,
            name: '钻石',
            value: _balance['gem'] ?? 0,
          ),
          _BalanceItem(
            icon: Icons.star,
            iconColor: Colors.purple,
            name: '荣誉',
            value: _balance['honor'] ?? 0,
          ),
          _BalanceItem(
            icon: Icons.points,
            iconColor: Colors.green,
            name: '积分',
            value: _balance['credit'] ?? 0,
          ),
        ],
      ),
    );
  }
  
  Widget _buildItemGrid() {
    if (_items.isEmpty) {
      return const Center(child: Text('暂无商品'));
    }
    
    return RefreshIndicator(
      onRefresh: () async => _loadItems(),
      child: GridView.builder(
        padding: const EdgeInsets.all(8),
        gridDelegate: const SliverGridDelegateWithFixedCrossAxisCount(
          crossAxisCount: 2,
          childAspectRatio: 0.75,
          crossAxisSpacing: 8,
          mainAxisSpacing: 8,
        ),
        itemCount: _items.length,
        itemBuilder: (context, index) {
          final item = _items[index];
          return _MallItemCard(item: item, onPurchase: () => _purchase(item));
        },
      ),
    );
  }
  
  @override
  void dispose() {
    _tabController.dispose();
    super.dispose();
  }
}

class _BalanceItem extends StatelessWidget {
  final IconData icon;
  final Color iconColor;
  final String name;
  final int value;
  
  const _BalanceItem({
    required this.icon,
    required this.iconColor,
    required this.name,
    required this.value,
  });
  
  @override
  Widget build(BuildContext context) {
    return Row(
      children: [
        Icon(icon, color: iconColor, size: 20),
        const SizedBox(width: 4),
        Text(
          _formatNumber(value),
          style: const TextStyle(fontWeight: FontWeight.bold),
        ),
      ],
    );
  }
  
  String _formatNumber(int n) {
    if (n >= 10000) {
      return '${(n / 10000).toStringAsFixed(1)}w';
    }
    return n.toString();
  }
}

class _MallItemCard extends StatelessWidget {
  final MallItem item;
  final VoidCallback onPurchase;
  
  const _MallItemCard({required this.item, required this.onPurchase});
  
  @override
  Widget build(BuildContext context) {
    return Card(
      child: InkWell(
        onTap: item.isOnSale && !item.isOutOfStock ? onPurchase : null,
        child: Padding(
          padding: const EdgeInsets.all(8),
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              Expanded(
                child: Stack(
                  children: [
                    Center(
                      child: item.imageUrl != null
                          ? Image.network(item.imageUrl!)
                          : Icon(
                              _getTypeIcon(),
                              size: 48,
                              color: Colors.grey,
                            ),
                    ),
                    if (item.hasDiscount)
                      Positioned(
                        top: 0,
                        right: 0,
                        child: Container(
                          padding: const EdgeInsets.symmetric(
                            horizontal: 6,
                            vertical: 2,
                          ),
                          decoration: BoxDecoration(
                            color: Colors.red,
                            borderRadius: BorderRadius.circular(4),
                          ),
                          child: Text(
                            '${item.discount}%',
                            style: const TextStyle(
                              color: Colors.white,
                              fontSize: 10,
                              fontWeight: FontWeight.bold,
                            ),
                          ),
                        ),
                      ),
                    if (item.isOutOfStock)
                      Container(
                        color: Colors.black54,
                        child: const Center(
                          child: Text(
                            '已售罄',
                            style: TextStyle(color: Colors.white),
                          ),
                        ),
                      ),
                  ],
                ),
              ),
              const SizedBox(height: 8),
              Text(
                item.name,
                style: const TextStyle(fontWeight: FontWeight.bold),
                maxLines: 1,
                overflow: TextOverflow.ellipsis,
              ),
              const SizedBox(height: 4),
              Row(
                children: [
                  Icon(
                    _getCurrencyIcon(),
                    size: 14,
                    color: _getCurrencyColor(),
                  ),
                  const SizedBox(width: 2),
                  Text(
                    item.hasDiscount
                        ? item.discountedPrice.toString()
                        : item.price.toString(),
                    style: TextStyle(
                      fontWeight: FontWeight.bold,
                      color: item.hasDiscount ? Colors.red : null,
                    ),
                  ),
                ],
              ),
              if (item.tags.isNotEmpty)
                Wrap(
                  spacing: 4,
                  children: item.tags
                      .take(2)
                      .map((t) => Container(
                            padding: const EdgeInsets.symmetric(
                              horizontal: 4,
                              vertical: 1,
                            ),
                            decoration: BoxDecoration(
                              color: Colors.blue[100],
                              borderRadius: BorderRadius.circular(2),
                            ),
                            child: Text(
                              t,
                              style: const TextStyle(fontSize: 10),
                            ),
                          ))
                      .toList(),
                ),
            ],
          ),
        ),
      ),
    );
  }
  
  IconData _getTypeIcon() {
    switch (item.type) {
      case MallType.gift:
        return Icons.card_giftcard;
      case MallType.item:
        return Icons.inventory_2;
      case MallType.skin:
        return Icons.checkroom;
      case MallType.random:
        return Icons.casino;
      case MallType.honor:
        return Icons.workspace_premium;
    }
  }
  
  IconData _getCurrencyIcon() {
    switch (item.currency) {
      case CurrencyType.gold:
        return Icons.monetization_on;
      case CurrencyType.gem:
        return Icons.diamond;
      case CurrencyType.honor:
        return Icons.star;
      case CurrencyType.credit:
        return Icons.points;
    }
  }
  
  Color _getCurrencyColor() {
    switch (item.currency) {
      case CurrencyType.gold:
        return Colors.amber;
      case CurrencyType.gem:
        return Colors.blue;
      case CurrencyType.honor:
        return Colors.purple;
      case CurrencyType.credit:
        return Colors.green;
    }
  }
}
