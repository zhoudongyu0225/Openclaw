// 邮件系统 - danmaku_game/server/mail.go
package main

import (
	"encoding/json"
	"errors"
	"time"
)

// MailType 邮件类型
type MailType int

const (
	MailTypeSystem   MailType = 1 // 系统邮件
	MailTypePlayer   MailType = 2 // 玩家邮件
	MailTypeGift     MailType = 3 // 礼物邮件
	MailTypeAuction  MailType = 4 // 拍卖邮件
	MailTypeGM       MailType = 5 // GM邮件
	MailTypeActivity MailType = 6 // 活动邮件
)

// MailAttachment 邮件附件
type MailAttachment struct {
	ItemID   string `json:"item_id"`   // 道具ID
	ItemName string `json:"item_name"` // 道具名称
	Count    int    `json:"count"`     // 数量
	Type     int    `json:"type"`      // 道具类型
}

// Mail 邮件结构
type Mail struct {
	MailID     string           `json:"mail_id"`     // 邮件ID
	ReceiverID int64            `json:"receiver_id"` // 接收者ID
	SenderID   int64            `json:"sender_id"`   // 发送者ID (0=系统)
	SenderName string           `json:"sender_name"` // 发送者名称
	Title      string           `json:"title"`       // 邮件标题
	Content    string           `json:"content"`     // 邮件内容
	Type       MailType         `json:"type"`        // 邮件类型
	Attachments []MailAttachment `json:"attachments"` // 附件列表
	ReadStatus bool              `json:"read_status"` // 已读状态
	Claimed    bool              `json:"claimed"`     // 是否已领取附件
	ExpiredAt  int64             `json:"expired_at"`  // 过期时间戳
	CreatedAt  int64             `json:"created_at"`  // 创建时间戳
}

// MailSystem 邮件系统
type MailSystem struct {
	db           *Database
	cache        *Cache
	mailIDGen    *IDGenerator
	maxMailCount int // 最大邮件数量
	mailExpire   time.Duration
}

// NewMailSystem 创建邮件系统
func NewMailSystem(db *Database, cache *Cache) *MailSystem {
	return &MailSystem{
		db:            db,
		cache:         cache,
		mailIDGen:     NewIDGenerator("mail"),
		maxMailCount:  100,
		mailExpire:    7 * 24 * time.Hour, // 7天过期
	}
}

// SendSystemMail 发送系统邮件
func (m *MailSystem) SendSystemMail(receiverID int64, title, content string, attachments []MailAttachment) (string, error) {
	mail := &Mail{
		MailID:      m.mailIDGen.Generate(),
		ReceiverID:  receiverID,
		SenderID:    0,
		SenderName:  "系统",
		Title:       title,
		Content:     content,
		Type:        MailTypeSystem,
		Attachments: attachments,
		ReadStatus:  false,
		Claimed:     false,
		ExpiredAt:   time.Now().Add(m.mailExpire).Unix(),
		CreatedAt:   time.Now().Unix(),
	}

	return m.mailIDGen.Generate(), m.saveMail(mail)
}

// SendPlayerMail 发送玩家邮件
func (m *MailSystem) SendPlayerMail(senderID int64, senderName string, receiverID int64, title, content string, attachments []MailAttachment) (string, error) {
	mail := &Mail{
		MailID:      m.mailIDGen.Generate(),
		ReceiverID:  receiverID,
		SenderID:    senderID,
		SenderName:  senderName,
		Title:       title,
		Content:     content,
		Type:        MailTypePlayer,
		Attachments: attachments,
		ReadStatus:  false,
		Claimed:     false,
		ExpiredAt:   time.Now().Add(m.mailExpire).Unix(),
		CreatedAt:   time.Now().Unix(),
	}

	return m.mailIDGen.Generate(), m.saveMail(mail)
}

// SendGMGMail 发送GM邮件
func (m *MailSystem) SendGMGMail(receiverID int64, title, content string, attachments []MailAttachment) (string, error) {
	mail := &Mail{
		MailID:      m.mailIDGen.Generate(),
		ReceiverID:  receiverID,
		SenderID:    0,
		SenderName:  "GM",
		Title:       title,
		Content:     content,
		Type:        MailTypeGM,
		Attachments: attachments,
		ReadStatus:  false,
		Claimed:     false,
		ExpiredAt:   time.Now().Add(m.mailExpire * 2).Unix(), // GM邮件30天
		CreatedAt:   time.Now().Unix(),
	}

	return m.mailIDGen.Generate(), m.saveMail(mail)
}

// GetMailList 获取邮件列表
func (m *MailSystem) GetMailList(playerID int64, offset, limit int) ([]*Mail, error) {
	cacheKey := fmt.Sprintf("mail:list:%d", playerID)
	
	// 尝试从缓存获取
	if cached, err := m.cache.Get(cacheKey); err == nil {
		var mails []*Mail
		if json.Unmarshal([]byte(cached), &mails) == nil {
			return mails, nil
		}
	}

	// 从数据库查询
	query := `SELECT mail_id, receiver_id, sender_id, sender_name, title, content, 
			  type, attachments, read_status, claimed, expired_at, created_at 
			  FROM mails WHERE receiver_id = ? AND expired_at > ? 
			  ORDER BY created_at DESC LIMIT ? OFFSET ?`
	
	rows, err := m.db.Query(query, playerID, time.Now().Unix(), limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var mails []*Mail
	for rows.Next() {
		mail := &Mail{}
		var attachmentsJSON string
		
		err := rows.Scan(&mail.MailID, &mail.ReceiverID, &mail.SenderID, &mail.SenderName,
			&mail.Title, &mail.Content, &mail.Type, &attachmentsJSON, &mail.ReadStatus,
			&mail.Claimed, &mail.ExpiredAt, &mail.CreatedAt)
		if err != nil {
			continue
		}
		
		json.Unmarshal([]byte(attachmentsJSON), &mail.Attachments)
		mails = append(mails, mail)
	}

	// 缓存结果
	if data, err := json.Marshal(mails); err == nil {
		m.cache.SetEX(cacheKey, string(data), 300) // 5分钟缓存
	}

	return mails, nil
}

// GetMail 获取单封邮件
func (m *MailSystem) GetMail(mailID string) (*Mail, error) {
	cacheKey := fmt.Sprintf("mail:%s", mailID)
	
	// 尝试从缓存获取
	if cached, err := m.cache.Get(cacheKey); err == nil {
		var mail Mail
		if json.Unmarshal([]byte(cached), &mail) == nil {
			return &mail, nil
		}
	}

	query := `SELECT mail_id, receiver_id, sender_id, sender_name, title, content, 
			  type, attachments, read_status, claimed, expired_at, created_at 
			  FROM mails WHERE mail_id = ?`
	
	mail := &Mail{}
	var attachmentsJSON string
	
	err := m.db.QueryRow(query, mailID).Scan(&mail.MailID, &mail.ReceiverID, &mail.SenderID,
		&mail.SenderName, &mail.Title, &mail.Content, &mail.Type, &attachmentsJSON,
		&mail.ReadStatus, &mail.Claimed, &mail.ExpiredAt, &mail.CreatedAt)
	if err != nil {
		return nil, err
	}

	json.Unmarshal([]byte(attachmentsJSON), &mail.Attachments)
	
	// 缓存
	if data, err := json.Marshal(mail); err == nil {
		m.cache.SetEX(cacheKey, string(data), 300)
	}

	return mail, nil
}

// ReadMail 读取邮件
func (m *MailSystem) ReadMail(mailID string, playerID int64) (*Mail, error) {
	mail, err := m.GetMail(mailID)
	if err != nil {
		return nil, err
	}

	if mail.ReceiverID != playerID {
		return nil, errors.New("无权访问此邮件")
	}

	if !mail.ReadStatus {
		query := `UPDATE mails SET read_status = true WHERE mail_id = ?`
		_, err = m.db.Exec(query, mailID)
		if err != nil {
			return nil, err
		}
		mail.ReadStatus = true
		
		// 清除缓存
		m.cache.Del(fmt.Sprintf("mail:%s", mailID))
		m.cache.Del(fmt.Sprintf("mail:list:%d", playerID))
	}

	return mail, nil
}

// ClaimAttachments 领取附件
func (m *MailSystem) ClaimAttachments(mailID string, playerID int64, inventory *Inventory) ([]MailAttachment, error) {
	mail, err := m.GetMail(mailID)
	if err != nil {
		return nil, err
	}

	if mail.ReceiverID != playerID {
		return nil, errors.New("无权访问此邮件")
	}

	if mail.Claimed {
		return nil, errors.New("附件已领取")
	}

	if len(mail.Attachments) == 0 {
		return nil, errors.New("没有附件可领取")
	}

	// 检查背包空间
	for _, att := range mail.Attachments {
		if !inventory.CanAddItem(att.ItemID, att.Count) {
			return nil, errors.New("背包空间不足")
		}
	}

	// 发放附件
	for _, att := range mail.Attachments {
		inventory.AddItem(att.ItemID, att.Count)
	}

	// 更新领取状态
	query := `UPDATE mails SET claimed = true WHERE mail_id = ?`
	_, err = m.db.Exec(query, mailID)
	if err != nil {
		return nil, err
	}

	mail.Claimed = true
	
	// 清除缓存
	m.cache.Del(fmt.Sprintf("mail:%s", mailID))
	m.cache.Del(fmt.Sprintf("mail:list:%d", playerID))

	return mail.Attachments, nil
}

// DeleteMail 删除邮件
func (m *MailSystem) DeleteMail(mailID string, playerID int64) error {
	mail, err := m.GetMail(mailID)
	if err != nil {
		return err
	}

	if mail.ReceiverID != playerID {
		return errors.New("无权删除此邮件")
	}

	query := `DELETE FROM mails WHERE mail_id = ?`
	_, err = m.db.Exec(query, mailID)
	if err != nil {
		return err
	}

	// 清除缓存
	m.cache.Del(fmt.Sprintf("mail:%s", mailID))
	m.cache.Del(fmt.Sprintf("mail:list:%d", playerID))

	return nil
}

// BatchDeleteReadMails 批量删除已读邮件
func (m *MailSystem) BatchDeleteReadMails(playerID int64) (int64, error) {
	query := `DELETE FROM mails WHERE receiver_id = ? AND read_status = true AND claimed = true`
	result, err := m.db.Exec(query, playerID)
	if err != nil {
		return 0, err
	}

	count, _ := result.RowsAffected()
	
	// 清除缓存
	m.cache.Del(fmt.Sprintf("mail:list:%d", playerID))

	return count, nil
}

// GetUnreadCount 获取未读邮件数量
func (m *MailSystem) GetUnreadCount(playerID int64) (int, error) {
	cacheKey := fmt.Sprintf("mail:unread:%d", playerID)
	
	// 尝试从缓存获取
	if cached, err := m.cache.Get(cacheKey); err == nil {
		var count int
		if fmt.Sscanf(cached, "%d", &count) == nil {
			return count, nil
		}
	}

	query := `SELECT COUNT(*) FROM mails WHERE receiver_id = ? AND read_status = false AND expired_at > ?`
	var count int
	err := m.db.QueryRow(query, playerID, time.Now().Unix()).Scan(&count)
	if err != nil {
		return 0, err
	}

	// 缓存
	m.cache.SetEX(cacheKey, fmt.Sprintf("%d", count), 60)

	return count, nil
}

// CleanupExpiredMails 清理过期邮件 (定时任务)
func (m *MailSystem) CleanupExpiredMails() (int64, error) {
	query := `DELETE FROM mails WHERE expired_at < ?`
	result, err := m.db.Exec(query, time.Now().Unix())
	if err != nil {
		return 0, err
	}

	count, _ := result.RowsAffected()
	return count, nil
}

// saveMail 保存邮件到数据库
func (m *MailSystem) saveMail(mail *Mail) error {
	attachmentsJSON, _ := json.Marshal(mail.Attachments)

	query := `INSERT INTO mails (mail_id, receiver_id, sender_id, sender_name, title, content, 
			  type, attachments, read_status, claimed, expired_at, created_at) 
			  VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	_, err := m.db.Exec(query, mail.MailID, mail.ReceiverID, mail.SenderID, mail.SenderName,
		mail.Title, mail.Content, mail.Type, attachmentsJSON, mail.ReadStatus,
		mail.Claimed, mail.ExpiredAt, mail.CreatedAt)

	if err == nil {
		// 清除列表缓存
		m.cache.Del(fmt.Sprintf("mail:list:%d", mail.ReceiverID))
		m.cache.Del(fmt.Sprintf("mail:unread:%d", mail.ReceiverID))
	}

	return err
}

// InitMailTable 初始化邮件表
func (m *MailSystem) InitMailTable() error {
	query := `CREATE TABLE IF NOT EXISTS mails (
		mail_id VARCHAR(64) PRIMARY KEY,
		receiver_id BIGINT NOT NULL,
		sender_id BIGINT NOT NULL DEFAULT 0,
		sender_name VARCHAR(64) NOT NULL,
		title VARCHAR(128) NOT NULL,
		content TEXT,
		type TINYINT NOT NULL DEFAULT 1,
		attachments JSON,
		read_status BOOLEAN NOT NULL DEFAULT FALSE,
		claimed BOOLEAN NOT NULL DEFAULT FALSE,
		expired_at BIGINT NOT NULL,
		created_at BIGINT NOT NULL,
		INDEX idx_receiver (receiver_id),
		INDEX idx_expired (expired_at)
	)`

	_, err := m.db.Exec(query)
	return err
}
