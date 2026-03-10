// ============================================
// 邮件客户端 - Flutter UI
// ============================================

import 'dart:async';
import 'dart:convert';
import 'package:flutter/material.dart';
import 'package:web_socket_channel/web_socket_channel.dart';

/// 邮件客户端
class MailClient {
  WebSocketChannel? _channel;
  final String serverUrl;
  final String token;
  
  final _mailListController = StreamController<List<Mail>>.broadcast();
  final _unreadCountController = StreamController<int>.broadcast();
  
  Stream<List<Mail>> get mailList => _mailListController.stream;
  Stream<int> get unreadCount => _unreadCountController.stream;
  
  MailClient({required this.serverUrl, required this.token});
  
  void connect() {
    _channel = WebSocketChannel.connect(
      Uri.parse('$serverUrl?token=$token&module=mail'),
    );
    _channel!.stream.listen(_onMessage);
  }
  
  void _onMessage(dynamic data) {
    final json = jsonDecode(data);
    final msgId = json['msg_id'];
    final jsonData = json['data'];
    
    switch (msgId) {
      case 10001: // SCMailList
        final mails = (jsonData['mails'] as List)
            .map((m) => Mail.fromJson(m))
            .toList();
        _mailListController.add(mails);
        break;
      case 10002: // SCUnreadCount
        _unreadCountController.add(jsonData['count']);
        break;
    }
  }
  
  /// 获取邮件列表
  void getMailList({int page = 1, int pageSize = 20}) {
    _channel?.sink.add(jsonEncode({
      'msg_id': 10001,
      'data': {'page': page, 'page_size': pageSize},
    }));
  }
  
  /// 读取邮件
  void readMail(String mailId) {
    _channel?.sink.add(jsonEncode({
      'msg_id': 10003,
      'data': {'mail_id': mailId},
    }));
  }
  
  /// 领取附件
  void claimAttachments(String mailId) {
    _channel?.sink.add(jsonEncode({
      'msg_id': 10004,
      'data': {'mail_id': mailId},
    }));
  }
  
  /// 删除邮件
  void deleteMail(String mailId) {
    _channel?.sink.add(jsonEncode({
      'msg_id': 10005,
      'data': {'mail_id': mailId},
    }));
  }
  
  /// 批量删除已读邮件
  void batchDeleteReadMails() {
    _channel?.sink.add(jsonEncode({
      'msg_id': 10006,
      'data': {},
    }));
  }
  
  /// 获取未读数量
  void getUnreadCount() {
    _channel?.sink.add(jsonEncode({
      'msg_id': 10002,
      'data': {},
    }));
  }
  
  void dispose() {
    _channel?.sink.close();
    _mailListController.close();
    _unreadCountController.close();
  }
}

/// 邮件模型
class Mail {
  final String id;
  final String title;
  final String content;
  final MailType type;
  final int senderId;
  final String senderName;
  final List<MailAttachment> attachments;
  final bool isRead;
  final int createdAt;
  final int expireAt;
  
  Mail({
    required this.id,
    required this.title,
    required this.content,
    required this.type,
    required this.senderId,
    required this.senderName,
    required this.attachments,
    required this.isRead,
    required this.createdAt,
    required this.expireAt,
  });
  
  factory Mail.fromJson(Map<String, dynamic> json) {
    return Mail(
      id: json['mail_id'],
      title: json['title'],
      content: json['content'],
      type: MailType.values.firstWhere(
        (e) => e.name == json['type'],
        orElse: () => MailType.system,
      ),
      senderId: json['sender_id'],
      senderName: json['sender_name'],
      attachments: (json['attachments'] as List?)
          ?.map((a) => MailAttachment.fromJson(a))
          .toList() ?? [],
      isRead: json['is_read'] ?? false,
      createdAt: json['created_at'],
      expireAt: json['expire_at'],
    );
  }
  
  bool get isExpired => DateTime.now().millisecondsSinceEpoch > expireAt;
  bool get hasAttachments => attachments.isNotEmpty;
}

/// 邮件类型
enum MailType {
  system,
  player,
  gift,
  auction,
  gm,
  activity,
}

/// 邮件附件
class MailAttachment {
  final String itemId;
  final String itemName;
  final int count;
  final AttachmentType type;
  
  MailAttachment({
    required this.itemId,
    required this.itemName,
    required this.count,
    required this.type,
  });
  
  factory MailAttachment.fromJson(Map<String, dynamic> json) {
    return MailAttachment(
      itemId: json['item_id'],
      itemName: json['item_name'],
      count: json['count'],
      type: AttachmentType.values.firstWhere(
        (e) => e.name == json['type'],
        orElse: () => AttachmentType.item,
      ),
    );
  }
}

/// 附件类型
enum AttachmentType {
  gold,
  gem,
  item,
  equipment,
  title,
  badge,
}

// ============ 邮件页面 ============

class MailPage extends StatefulWidget {
  final MailClient mailClient;
  
  const MailPage({super.key, required this.mailClient});
  
  @override
  State<MailPage> createState() => _MailPageState();
}

class _MailPageState extends State<MailPage> with SingleTickerProviderStateMixin {
  late TabController _tabController;
  List<Mail> _mails = [];
  bool _isLoading = false;
  
  @override
  void initState() {
    super.initState();
    _tabController = TabController(length: 2, vsync: this);
    _loadMails();
  }
  
  void _loadMails() {
    setState(() => _isLoading = true);
    widget.mailClient.getMailList();
    
    widget.mailClient.mailList.listen((mails) {
      setState(() {
        _mails = mails;
        _isLoading = false;
      });
    });
  }
  
  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text('邮件'),
        bottom: TabBar(
          controller: _tabController,
          tabs: const [
            Tab(text: '收件箱'),
            Tab(text: '已发送'),
          ],
        ),
        actions: [
          PopupMenuButton<String>(
            onSelected: (value) {
              if (value == 'delete_read') {
                widget.mailClient.batchDeleteReadMails();
              }
            },
            itemBuilder: (context) => [
              const PopupMenuItem(value: 'delete_read', child: Text('删除已读')),
            ],
          ),
        ],
      ),
      body: TabBarView(
        controller: _tabController,
        children: [
          _buildMailList(),
          _buildSentList(),
        ],
      ),
    );
  }
  
  Widget _buildMailList() {
    if (_isLoading) {
      return const Center(child: CircularProgressIndicator());
    }
    
    if (_mails.isEmpty) {
      return const Center(child: Text('暂无邮件'));
    }
    
    return RefreshIndicator(
      onRefresh: () async => _loadMails(),
      child: ListView.builder(
        itemCount: _mails.length,
        itemBuilder: (context, index) {
          final mail = _mails[index];
          return _MailListItem(
            mail: mail,
            onTap: () => _showMailDetail(mail),
            onClaim: mail.hasAttachments ? () {
              widget.mailClient.claimAttachments(mail.id);
            } : null,
            onDelete: () {
              widget.mailClient.deleteMail(mail.id);
            },
          );
        },
      ),
    );
  }
  
  Widget _buildSentList() {
    return const Center(child: Text('已发送邮件'));
  }
  
  void _showMailDetail(Mail mail) {
    widget.mailClient.readMail(mail.id);
    
    showModalBottomSheet(
      context: context,
      isScrollControlled: true,
      builder: (context) => DraggableScrollableSheet(
        initialChildSize: 0.7,
        minChildSize: 0.5,
        maxChildSize: 0.95,
        expand: false,
        builder: (context, scrollController) => _MailDetailSheet(
          mail: mail,
          scrollController: scrollController,
          onClaim: mail.hasAttachments ? () {
            widget.mailClient.claimAttachments(mail.id);
            Navigator.pop(context);
          } : null,
          onDelete: () {
            widget.mailClient.deleteMail(mail.id);
            Navigator.pop(context);
          },
        ),
      ),
    );
  }
  
  @override
  void dispose() {
    _tabController.dispose();
    super.dispose();
  }
}

class _MailListItem extends StatelessWidget {
  final Mail mail;
  final VoidCallback onTap;
  final VoidCallback? onClaim;
  final VoidCallback onDelete;
  
  const _MailListItem({
    required this.mail,
    required this.onTap,
    this.onClaim,
    required this.onDelete,
  });
  
  @override
  Widget build(BuildContext context) {
    return Card(
      margin: const EdgeInsets.symmetric(horizontal: 8, vertical: 4),
      child: ListTile(
        leading: Stack(
          children: [
            Icon(_getMailIcon(), size: 32),
            if (!mail.isRead)
              Positioned(
                right: 0,
                top: 0,
                child: Container(
                  width: 8,
                  height: 8,
                  decoration: const BoxDecoration(
                    color: Colors.red,
                    shape: BoxShape.circle,
                  ),
                ),
              ),
          ],
        ),
        title: Text(
          mail.title,
          style: TextStyle(
            fontWeight: mail.isRead ? FontWeight.normal : FontWeight.bold,
          ),
        ),
        subtitle: Text(
          mail.senderName,
          style: const TextStyle(fontSize: 12),
        ),
        trailing: Row(
          mainAxisSize: MainAxisSize.min,
          children: [
            if (mail.hasAttachments)
              const Icon(Icons.attach_file, size: 16),
            if (onClaim != null)
              TextButton(
                onPressed: onClaim,
                child: const Text('领取'),
              ),
          ],
        ),
        onTap: onTap,
      ),
    );
  }
  
  IconData _getMailIcon() {
    switch (mail.type) {
      case MailType.system:
        return Icons.system_update;
      case MailType.player:
        return Icons.person;
      case MailType.gift:
        return Icons.card_giftcard;
      case MailType.auction:
        return Icons.gavel;
      case MailType.gm:
        return Icons.admin_panel_settings;
      case MailType.activity:
        return Icons.celebration;
    }
  }
}

class _MailDetailSheet extends StatelessWidget {
  final Mail mail;
  final ScrollController scrollController;
  final VoidCallback? onClaim;
  final VoidCallback onDelete;
  
  const _MailDetailSheet({
    required this.mail,
    required this.scrollController,
    this.onClaim,
    required this.onDelete,
  });
  
  @override
  Widget build(BuildContext context) {
    return Container(
      padding: const EdgeInsets.all(16),
      child: ListView(
        controller: scrollController,
        children: [
          Row(
            children: [
              Expanded(
                child: Text(
                  mail.title,
                  style: const TextStyle(
                    fontSize: 18,
                    fontWeight: FontWeight.bold,
                  ),
                ),
              ),
              IconButton(
                icon: const Icon(Icons.delete_outline),
                onPressed: onDelete,
              ),
            ],
          ),
          const SizedBox(height: 8),
          Text(
            '来自: ${mail.senderName}',
            style: TextStyle(color: Colors.grey[600]),
          ),
          Text(
            '有效期至: ${DateTime.fromMillisecondsSinceEpoch(mail.expireAt)}',
            style: TextStyle(color: Colors.grey[600], fontSize: 12),
          ),
          const Divider(height: 24),
          Text(mail.content),
          if (mail.attachments.isNotEmpty) ...[
            const Divider(height: 24),
            const Text(
              '附件',
              style: TextStyle(fontWeight: FontWeight.bold),
            ),
            const SizedBox(height: 8),
            ...mail.attachments.map((a) => ListTile(
              leading: Icon(_getAttachmentIcon(a.type)),
              title: Text(a.itemName),
              trailing: Text('x${a.count}'),
            )),
            if (onClaim != null)
              ElevatedButton(
                onPressed: onClaim,
                child: const Text('领取附件'),
              ),
          ],
        ],
      ),
    );
  }
  
  IconData _getAttachmentIcon(AttachmentType type) {
    switch (type) {
      case AttachmentType.gold:
        return Icons.monetization_on;
      case AttachmentType.gem:
        return Icons.diamond;
      case AttachmentType.item:
        return Icons.inventory_2;
      case AttachmentType.equipment:
        return Icons.shield;
      case AttachmentType.title:
        return Icons.workspace_premium;
      case AttachmentType.badge:
        return Icons.military_tech;
    }
  }
}
