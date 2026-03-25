package antiraid

import (
	"context"
	"sync"
	"time"

	"sentinel-adaptive/internal/modules/audit"
	"sentinel-adaptive/internal/playbook"
	"sentinel-adaptive/internal/utils"

	"github.com/bwmarrin/discordgo"
)

type Module struct {
	mu       sync.Mutex
	counters map[string]*utils.JoinCounter
	windows  map[string]time.Duration
	playbook *playbook.Engine
	audit    *audit.Logger
}

func New(playbookEngine *playbook.Engine, auditLogger *audit.Logger) *Module {
	return &Module{
		counters: make(map[string]*utils.JoinCounter),
		windows:  make(map[string]time.Duration),
		playbook: playbookEngine,
		audit:    auditLogger,
	}
}

func (m *Module) HandleJoin(ctx context.Context, session *discordgo.Session, event *discordgo.GuildMemberAdd, raidJoins, raidWindowSeconds int) bool {
	guildID := event.GuildID
	if guildID == "" && event.Member != nil {
		guildID = event.Member.GuildID
	}
	if guildID == "" {
		return false
	}

	window := time.Duration(raidWindowSeconds) * time.Second
	counter := m.getCounter(guildID, window)
	count := counter.Add(time.Now())
	if count < raidJoins {
		return false
	}

	userID := ""
	if event.Member != nil && event.Member.User != nil {
		userID = event.Member.User.ID
	}
	m.audit.Log(ctx, audit.LevelWarn, guildID, userID, "anti_raid", "raid threshold reached")
	return true
}

func (m *Module) getCounter(guildID string, window time.Duration) *utils.JoinCounter {
	m.mu.Lock()
	defer m.mu.Unlock()
	if prev, ok := m.windows[guildID]; !ok || prev != window {
		m.counters[guildID] = utils.NewJoinCounter(window)
		m.windows[guildID] = window
	}
	return m.counters[guildID]
}
