package guild

import (
	"sync"
	"time"
)

// Guild represents a guild/party in the game
type Guild struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	LeaderID    string    `json:"leader_id"`
	Members     []*Member `json:"members"`
	Level       int       `json:"level"`
	Exp         int       `json:"exp"`
	Notice      string    `json:"notice"`
	MaxMembers  int       `json:"max_members"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Member represents a guild member
type Member struct {
	ID        string    `json:"id"`
	PlayerID  string    `json:"player_id"`
	Role      MemberRole `json:"role"` // Leader, Officer, Member
	Contrib   int       `json:"contrib"` // Contribution points
	JoinTime  time.Time `json:"join_time"`
}

// MemberRole guild member roles
type MemberRole int

const (
	RoleLeader MemberRole = iota + 1
	RoleOfficer
	RoleMember
)

// GuildManager manages all guilds
type GuildManager struct {
	mu      sync.RWMutex
	guilds  map[string]*Guild
	players map[string]string // playerID -> guildID
}

// NewGuildManager creates a new guild manager
func NewGuildManager() *GuildManager {
	return &GuildManager{
		guilds:  make(map[string]*Guild),
		players: make(map[string]string),
	}
}

// CreateGuild creates a new guild
func (gm *GuildManager) CreateGuild(name, leaderID string) (*Guild, error) {
	gm.mu.Lock()
	defer gm.mu.Unlock()

	// Check if player is already in a guild
	if _, exists := gm.players[leaderID]; exists {
		return nil, ErrPlayerInGuild
	}

	guild := &Guild{
		ID:         generateID(),
		Name:       name,
		LeaderID:   leaderID,
		Members:    []*Member{{PlayerID: leaderID, Role: RoleLeader, JoinTime: time.Now()}},
		Level:      1,
		MaxMembers: 50,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	gm.guilds[guild.ID] = guild
	gm.players[leaderID] = guild.ID

	return guild, nil
}

// JoinGuild adds a player to a guild
func (gm *GuildManager) JoinGuild(guildID, playerID string) error {
	gm.mu.Lock()
	defer gm.mu.Unlock()

	guild, ok := gm.guilds[guildID]
	if !ok {
		return ErrGuildNotFound
	}

	if len(guild.Members) >= guild.MaxMembers {
		return ErrGuildFull
	}

	if _, exists := gm.players[playerID]; exists {
		return ErrPlayerInGuild
	}

	guild.Members = append(guild.Members, &Member{
		PlayerID: playerID,
		Role:     RoleMember,
		JoinTime: time.Now(),
	})
	gm.players[playerID] = guildID
	guild.UpdatedAt = time.Now()

	return nil
}

// LeaveGuild removes a player from their guild
func (gm *GuildManager) LeaveGuild(playerID string) error {
	gm.mu.Lock()
	defer gm.mu.Unlock()

	guildID, ok := gm.players[playerID]
	if !ok {
		return ErrPlayerNotInGuild
	}

	guild := gm.guilds[guildID]
	for i, m := range guild.Members {
		if m.PlayerID == playerID {
			// Leader cannot leave, must transfer first
			if m.Role == RoleLeader && len(guild.Members) > 1 {
				return ErrLeaderCannotLeave
			}
			guild.Members = append(guild.Members[:i], guild.Members[i+1:]...)
			break
		}
	}

	delete(gm.players, playerID)
	guild.UpdatedAt = time.Now()

	// Delete guild if empty
	if len(guild.Members) == 0 {
		delete(gm.guilds, guildID)
	}

	return nil
}

// GetGuild gets a guild by ID
func (gm *GuildManager) GetGuild(guildID string) (*Guild, error) {
	gm.mu.RLock()
	defer gm.mu.RUnlock()

	guild, ok := gm.guilds[guildID]
	if !ok {
		return nil, ErrGuildNotFound
	}
	return guild, nil
}

// GetPlayerGuild gets the guild a player belongs to
func (gm *GuildManager) GetPlayerGuild(playerID string) (*Guild, error) {
	gm.mu.RLock()
	defer gm.mu.RUnlock()

	guildID, ok := gm.players[playerID]
	if !ok {
		return nil, ErrPlayerNotInGuild
	}
	return gm.guilds[guildID], nil
}

// ListGuilds lists all guilds (with pagination)
func (gm *GuildManager) ListGuilds(page, pageSize int) []*Guild {
	gm.mu.RLock()
	defer gm.mu.RUnlock()

	result := make([]*Guild, 0, len(gm.guilds))
	for _, g := range gm.guilds {
		result = append(result, g)
	}

	start := page * pageSize
	if start >= len(result) {
		return []*Guild{}
	}
	end := start + pageSize
	if end > len(result) {
		end = len(result)
	}

	return result[start:end]
}

// UpdateNotice updates guild notice
func (gm *GuildManager) UpdateNotice(guildID, playerID, notice string) error {
	gm.mu.Lock()
	defer gm.mu.Unlock()

	guild, ok := gm.guilds[guildID]
	if !ok {
		return ErrGuildNotFound
	}

	// Only leader and officers can update notice
	for _, m := range guild.Members {
		if m.PlayerID == playerID && (m.Role == RoleLeader || m.Role == RoleOfficer) {
			guild.Notice = notice
			guild.UpdatedAt = time.Now()
			return nil
		}
	}

	return ErrNoPermission
}

// AddContrib adds contribution points to a member
func (gm *GuildManager) AddContrib(guildID, playerID string, points int) error {
	gm.mu.Lock()
	defer gm.mu.Unlock()

	guild, ok := gm.guilds[guildID]
	if !ok {
		return ErrGuildNotFound
	}

	for _, m := range guild.Members {
		if m.PlayerID == playerID {
			m.Contrib += points
			// Check for level up
			requiredExp := guild.Level * 1000
			if m.Contrib >= requiredExp {
				guild.Level++
			}
			guild.UpdatedAt = time.Now()
			return nil
		}
	}

	return ErrPlayerNotInGuild
}

// TransferLeader transfers guild leadership
func (gm *GuildManager) TransferLeader(guildID, currentLeaderID, newLeaderID string) error {
	gm.mu.Lock()
	defer gm.mu.Unlock()

	guild, ok := gm.guilds[guildID]
	if !ok {
		return ErrGuildNotFound
	}

	if guild.LeaderID != currentLeaderID {
		return ErrNoPermission
	}

	for _, m := range guild.Members {
		if m.PlayerID == newLeaderID {
			m.Role = RoleLeader
		}
		if m.PlayerID == currentLeaderID {
			m.Role = RoleMember
		}
	}

	guild.LeaderID = newLeaderID
	guild.UpdatedAt = time.Now()

	return nil
}

// KickMember kicks a member from guild
func (gm *GuildManager) KickMember(guildID, kickerID, targetID string) error {
	gm.mu.Lock()
	defer gm.mu.Unlock()

	guild, ok := gm.guilds[guildID]
	if !ok {
		return ErrGuildNotFound
	}

	// Only leader and officers can kick
	var kickerRole MemberRole
	for _, m := range guild.Members {
		if m.PlayerID == kickerID {
			kickerRole = m.Role
		}
	}

	if kickerRole != RoleLeader && kickerRole != RoleOfficer {
		return ErrNoPermission
	}

	for i, m := range guild.Members {
		if m.PlayerID == targetID {
			// Cannot kick leader
			if m.Role == RoleLeader {
				return ErrNoPermission
			}
			guild.Members = append(guild.Members[:i], guild.Members[i+1:]...)
			delete(gm.players, targetID)
			guild.UpdatedAt = time.Now()
			return nil
		}
	}

	return ErrPlayerNotInGuild
}

// Guild errors
var (
	ErrGuildNotFound   = &GuildError{"guild not found"}
	ErrPlayerInGuild   = &GuildError{"player already in a guild"}
	ErrPlayerNotInGuild = &GuildError{"player not in a guild"}
	ErrGuildFull       = &GuildError{"guild is full"}
	ErrLeaderCannotLeave = &GuildError{"leader cannot leave, transfer leadership first"}
	ErrNoPermission    = &GuildError{"no permission"}
)

type GuildError struct {
	msg string
}

func (e *GuildError) Error() string {
	return e.msg
}

// Helper to generate unique ID
func generateID() string {
	return time.Now().Format("20060102150405") + randomString(6)
}

func randomString(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[time.Now().UnixNano()%int64(len(letters))]
	}
	return string(b)
}
