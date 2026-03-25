package audit

import (
	"context"
	"time"

	"sentinel-adaptive/internal/storage"

	"go.uber.org/zap"
)

const (
	LevelInfo = "INFO"
	LevelWarn = "WARN"
	LevelCrit = "CRIT"
)

type UserResolver func(guildID, userID string) string

type Logger struct {
	store        *storage.Store
	logger       *zap.Logger
	notify       func(context.Context, storage.AuditLog)
	userResolver UserResolver
}

func NewLogger(store *storage.Store, logger *zap.Logger) *Logger {
	return &Logger{store: store, logger: logger}
}

func (l *Logger) SetNotifier(notify func(context.Context, storage.AuditLog)) {
	l.notify = notify
}

func (l *Logger) SetUserResolver(resolver UserResolver) {
	l.userResolver = resolver
}

func (l *Logger) Log(ctx context.Context, level, guildID, userID, event, details string) {
	entry := storage.AuditLog{
		GuildID:   guildID,
		UserID:    userID,
		Level:     level,
		Event:     event,
		Details:   details,
		CreatedAt: time.Now(),
	}
	if l.store != nil {
		_ = l.store.AddAuditLog(ctx, entry)
	}
	if l.notify != nil {
		l.notify(ctx, entry)
	}
	userName := userID
	if l.userResolver != nil && userID != "" {
		if name := l.userResolver(guildID, userID); name != "" {
			userName = name
		}
	}
	l.logger.Info("audit",
		zap.String("level", level),
		zap.String("guild_id", guildID),
		zap.String("user", userName),
		zap.String("event", event),
		zap.String("details", details),
	)
}
