package antiraid

import (
	"context"
	"testing"

	"sentinel-adaptive/internal/modules/audit"
	"sentinel-adaptive/internal/playbook"
	"sentinel-adaptive/internal/storage"

	"github.com/bwmarrin/discordgo"
	"go.uber.org/zap"
)

func TestRaidJoinCounter(t *testing.T) {
	store, _ := storage.New(":memory:")
	_ = store.Migrate()
	auditLogger := audit.NewLogger(store, zap.NewNop())
	playbookEngine := playbook.New(playbook.Config{LockdownMinutes: 1, StrictModeMinutes: 1, ExitStepSeconds: 1}, auditLogger)
	module := New(playbookEngine, auditLogger)

	session := &discordgo.Session{}
	ctx := context.Background()

	join := &discordgo.GuildMemberAdd{GuildID: "g1", Member: &discordgo.Member{GuildID: "g1", User: &discordgo.User{ID: "u1"}}}
	module.HandleJoin(ctx, session, join, 3, 5)
	module.HandleJoin(ctx, session, join, 3, 5)
	module.HandleJoin(ctx, session, join, 3, 5)
	state := playbookEngine.IsLockdown("g1")
	if !state.Lockdown {
		t.Fatalf("expected lockdown")
	}
}

func TestRaidPerGuildThreshold(t *testing.T) {
	store, _ := storage.New(":memory:")
	_ = store.Migrate()
	auditLogger := audit.NewLogger(store, zap.NewNop())
	playbookEngine := playbook.New(playbook.Config{LockdownMinutes: 1, StrictModeMinutes: 1, ExitStepSeconds: 1}, auditLogger)
	module := New(playbookEngine, auditLogger)

	session := &discordgo.Session{}
	ctx := context.Background()

	// guild g1 seuil=5, guild g2 seuil=2
	g1join := &discordgo.GuildMemberAdd{GuildID: "g1", Member: &discordgo.Member{GuildID: "g1", User: &discordgo.User{ID: "u1"}}}
	g2join := &discordgo.GuildMemberAdd{GuildID: "g2", Member: &discordgo.Member{GuildID: "g2", User: &discordgo.User{ID: "u2"}}}

	// 3 joins sur g1 (seuil=5) : pas de lockdown
	module.HandleJoin(ctx, session, g1join, 5, 10)
	module.HandleJoin(ctx, session, g1join, 5, 10)
	module.HandleJoin(ctx, session, g1join, 5, 10)
	if playbookEngine.IsLockdown("g1").Lockdown {
		t.Fatalf("g1 ne devrait pas être en lockdown")
	}

	// 2 joins sur g2 (seuil=2) : lockdown
	module.HandleJoin(ctx, session, g2join, 2, 10)
	module.HandleJoin(ctx, session, g2join, 2, 10)
	if !playbookEngine.IsLockdown("g2").Lockdown {
		t.Fatalf("g2 devrait être en lockdown")
	}
}
