package mail

import (
	"sync"
	"time"
)

// Mail represents an in-game mail
type Mail struct {
	ID         string    `json:"id"`
	SenderID   string    `json:"sender_id"`
	SenderName string    `json:"sender_name"`
	ReceiverID string    `json:"receiver_id"`
	Title      string    `json:"title"`
	Content    string    `json:"content"`
	Attachments []*Item  `json:"attachments"`
	IsRead     bool      `json:"is_read"`
	IsClaimed  bool      `json:"is_claimed"`
	ExpireAt   time.Time `json:"expire_at"`
	CreatedAt  time.Time `json:"created_at"`
}

// Item represents an item in mail attachment
type Item struct {
	ItemID   string `json:"item_id"`
	ItemType string `json:"item_type"`
	Count    int    `json:"count"`
}

// MailManager manages all mails
type MailManager struct {
	mu     sync.RWMutex
	mails  map[string]*Mail // mailID -> mail
	inbox  map[string][]string // playerID -> mailIDs
}

// NewMailManager creates a new mail manager
func NewMailManager() *MailManager {
	return &MailManager{
		mails:  make(map[string]*Mail),
		inbox:  make(map[string][]string),
	}
}

// SendMail sends a mail to a player
func (mm *MailManager) SendMail(receiverID, title, content string, attachments []*Item) (*Mail, error) {
	mm.mu.Lock()
	defer mm.mu.Unlock()

	if title == "" || content == "" {
		return nil, ErrInvalidContent
	}

	mail := &Mail{
		ID:         generateID(),
		SenderID:   "SYSTEM",
		SenderName: "系统",
		ReceiverID: receiverID,
		Title:      title,
		Content:    content,
		Attachments: attachments,
		ExpireAt:   time.Now().Add(30 * 24 * time.Hour), // 30 days expire
		CreatedAt:  time.Now(),
	}

	mm.mails[mail.ID] = mail
	mm.inbox[receiverID] = append(mm.inbox[receiverID], mail.ID)

	return mail, nil
}

// SendSystemMail sends a system mail
func (mm *MailManager) SendSystemMail(receiverID, title, content string, attachments []*Item) (*Mail, error) {
	return mm.SendMail(receiverID, title, content, attachments)
}

// BroadcastMail sends mail to multiple players
func (mm *MailManager) BroadcastMail(receiverIDs []string, title, content string, attachments []*Item) []*Mail {
	mm.mu.Lock()
	defer mm.mu.Unlock()

	mails := make([]*Mail, 0, len(receiverIDs))
	for _, rid := range receiverIDs {
		mail := &Mail{
			ID:         generateID(),
			SenderID:   "SYSTEM",
			SenderName: "系统",
			ReceiverID: rid,
			Title:      title,
			Content:    content,
			Attachments: attachments,
			ExpireAt:   time.Now().Add(30 * 24 * time.Hour),
			CreatedAt:  time.Now(),
		}
		mm.mails[mail.ID] = mail
		mm.inbox[rid] = append(mm.inbox[rid], mail.ID)
		mails = append(mails, mail)
	}

	return mails
}

// GetMail gets a mail by ID
func (mm *MailManager) GetMail(mailID string) (*Mail, error) {
	mm.mu.RLock()
	defer mm.mu.RUnlock()

	mail, ok := mm.mails[mailID]
	if !ok {
		return nil, ErrMailNotFound
	}
	return mail, nil
}

// GetInbox gets all mails for a player
func (mm *MailManager) GetInbox(playerID string) ([]*Mail, error) {
	mm.mu.RLock()
	defer mm.mu.RUnlock()

	mailIDs, ok := mm.inbox[playerID]
	if !ok {
		return []*Mail{}, nil
	}

	mails := make([]*Mail, 0, len(mailIDs))
	for _, id := range mailIDs {
		if mail, exists := mm.mails[id]; exists {
			// Check if expired
			if time.Now().Before(mail.ExpireAt) {
				mails = append(mails, mail)
			}
		}
	}

	return mails, nil
}

// GetUnreadCount gets count of unread mails
func (mm *MailManager) GetUnreadCount(playerID string) int {
	mm.mu.RLock()
	defer mm.mu.RUnlock()

	mailIDs, ok := mm.inbox[playerID]
	if !ok {
		return 0
	}

	count := 0
	for _, id := range mailIDs {
		if mail, exists := mm.mails[id]; exists && !mail.IsRead && time.Now().Before(mail.ExpireAt) {
			count++
		}
	}

	return count
}

// MarkAsRead marks a mail as read
func (mm *MailManager) MarkAsRead(mailID, playerID string) error {
	mm.mu.Lock()
	defer mm.mu.Unlock()

	mail, ok := mm.mails[mailID]
	if !ok {
		return ErrMailNotFound
	}

	if mail.ReceiverID != playerID {
		return ErrNoPermission
	}

	mail.IsRead = true
	return nil
}

// ClaimAttachments claims mail attachments
func (mm *MailManager) ClaimAttachments(mailID, playerID string) ([]*Item, error) {
	mm.mu.Lock()
	defer mm.mu.Unlock()

	mail, ok := mm.mails[mailID]
	if !ok {
		return nil, ErrMailNotFound
	}

	if mail.ReceiverID != playerID {
		return nil, ErrNoPermission
	}

	if mail.IsClaimed {
		return nil, ErrAlreadyClaimed
	}

	if len(mail.Attachments) == 0 {
		return nil, ErrNoAttachments
	}

	mail.IsClaimed = true
	return mail.Attachments, nil
}

// DeleteMail deletes a mail
func (mm *MailManager) DeleteMail(mailID, playerID string) error {
	mm.mu.Lock()
	defer mm.mu.Unlock()

	mail, ok := mm.mails[mailID]
	if !ok {
		return ErrMailNotFound
	}

	if mail.ReceiverID != playerID {
		return ErrNoPermission
	}

	// Remove from inbox
	inbox := mm.inbox[playerID]
	for i, id := range inbox {
		if id == mailID {
			mm.inbox[playerID] = append(inbox[:i], inbox[i+1:]...)
			break
		}
	}

	// Remove mail
	delete(mm.mails, mailID)

	return nil
}

// CleanExpired cleans expired mails
func (mm *MailManager) CleanExpired() int {
	mm.mu.Lock()
	defer mm.mu.Unlock()

	now := time.Now()
	count := 0

	for playerID, mailIDs := range mm.inbox {
		validIDs := make([]string, 0)
		for _, id := range mailIDs {
			if mail, exists := mm.mails[id]; exists {
				if now.Before(mail.ExpireAt) {
					validIDs = append(validIDs, id)
				} else {
					delete(mm.mails, id)
					count++
				}
			}
		}
		mm.inbox[playerID] = validIDs
	}

	return count
}

// BatchSendReward sends rewards to multiple players
func (mm *MailManager) BatchSendReward(receiverIDs []string, title, content string, items []*Item) []*Mail {
	return mm.BroadcastMail(receiverIDs, title, content, items)
}

// Mail errors
var (
	ErrInvalidContent  = &MailError{"invalid mail content"}
	ErrMailNotFound    = &MailError{"mail not found"}
	ErrNoPermission    = &MailError{"no permission"}
	ErrAlreadyClaimed  = &MailError{"attachments already claimed"}
	ErrNoAttachments   = &MailError{"no attachments to claim"}
)

type MailError struct {
	msg string
}

func (e *MailError) Error() string {
	return e.msg
}

func generateID() string {
	return time.Now().Format("20060102150405") + randomString(8)
}

func randomString(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[time.Now().UnixNano()%int64(len(letters))]
	}
	return string(b)
}
