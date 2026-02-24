package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/Rezann47/YksKoc/internal/domain/entity"
)

type MessageRepository interface {
	// Mesaj gönder
	Create(ctx context.Context, msg *entity.Message) error

	// İki kullanıcı arasındaki mesaj geçmişi (sayfalı, yeniden eskiye)
	ListConversation(ctx context.Context, userA, userB uuid.UUID, limit, offset int) ([]entity.Message, int64, error)

	// Kullanıcının tüm konuşmaları (son mesajla birlikte)
	ListConversations(ctx context.Context, userID uuid.UUID) ([]ConversationRow, error)

	// Okunmamış mesaj sayısı (isteğe bağlı: belirli bir göndericiden)
	UnreadCount(ctx context.Context, receiverID uuid.UUID, senderID *uuid.UUID) (int64, error)

	// İki kullanıcı arasındaki tüm okunmamışları okundu işaretle
	MarkRead(ctx context.Context, receiverID, senderID uuid.UUID) error
}

// ListConversations için ham satır
type ConversationRow struct {
	PeerID       uuid.UUID
	PeerName     string
	PeerAvatar   int16
	PeerOnline   bool
	PeerLastSeen *time.Time
	LastMessage  string
	LastAt       time.Time
	UnreadCount  int64
}

// ─── Implementation ───────────────────────────────────────

type messageRepository struct{ db *gorm.DB }

func NewMessageRepository(db *gorm.DB) MessageRepository {
	return &messageRepository{db: db}
}

func (r *messageRepository) Create(ctx context.Context, msg *entity.Message) error {
	return r.db.WithContext(ctx).Create(msg).Error
}

func (r *messageRepository) ListConversation(ctx context.Context, userA, userB uuid.UUID, limit, offset int) ([]entity.Message, int64, error) {
	var msgs []entity.Message
	var total int64

	base := r.db.WithContext(ctx).Model(&entity.Message{}).
		Where(
			"(sender_id = ? AND receiver_id = ?) OR (sender_id = ? AND receiver_id = ?)",
			userA, userB, userB, userA,
		)

	base.Count(&total)

	err := base.
		Order("created_at DESC").
		Limit(limit).Offset(offset).
		Find(&msgs).Error

	return msgs, total, err
}

func (r *messageRepository) ListConversations(ctx context.Context, userID uuid.UUID) ([]ConversationRow, error) {
	// Her peer ile son mesajı + okunmamış sayısını çek
	query := `
		WITH peers AS (
			SELECT DISTINCT
				CASE WHEN sender_id = @uid THEN receiver_id ELSE sender_id END AS peer_id
			FROM messages
			WHERE sender_id = @uid OR receiver_id = @uid
		),
		last_msgs AS (
			SELECT DISTINCT ON (
				LEAST(sender_id, receiver_id),
				GREATEST(sender_id, receiver_id)
			)
				CASE WHEN sender_id = @uid THEN receiver_id ELSE sender_id END AS peer_id,
				content AS last_message,
				created_at AS last_at
			FROM messages
			WHERE sender_id = @uid OR receiver_id = @uid
			ORDER BY
				LEAST(sender_id, receiver_id),
				GREATEST(sender_id, receiver_id),
				created_at DESC
		),
		unread_counts AS (
			SELECT sender_id AS peer_id, COUNT(*) AS unread_count
			FROM messages
			WHERE receiver_id = @uid AND is_read = false
			GROUP BY sender_id
		)
		SELECT
			p.peer_id,
			u.name       AS peer_name,
			u.avatar_id  AS peer_avatar,
			(u.last_seen_at IS NOT NULL AND u.last_seen_at > NOW() - INTERVAL '3 minutes') AS peer_online,
			u.last_seen_at AS peer_last_seen,
			lm.last_message,
			lm.last_at,
			COALESCE(uc.unread_count, 0) AS unread_count
		FROM peers p
		JOIN users u ON u.id = p.peer_id AND u.deleted_at IS NULL
		JOIN last_msgs lm ON lm.peer_id = p.peer_id
		LEFT JOIN unread_counts uc ON uc.peer_id = p.peer_id
		ORDER BY lm.last_at DESC
	`

	var rows []ConversationRow
	err := r.db.WithContext(ctx).Raw(query, map[string]interface{}{"uid": userID}).Scan(&rows).Error
	return rows, err
}

func (r *messageRepository) UnreadCount(ctx context.Context, receiverID uuid.UUID, senderID *uuid.UUID) (int64, error) {
	q := r.db.WithContext(ctx).Model(&entity.Message{}).
		Where("receiver_id = ? AND is_read = false", receiverID)
	if senderID != nil {
		q = q.Where("sender_id = ?", *senderID)
	}
	var count int64
	err := q.Count(&count).Error
	return count, err
}

func (r *messageRepository) MarkRead(ctx context.Context, receiverID, senderID uuid.UUID) error {
	now := time.Now()
	return r.db.WithContext(ctx).
		Model(&entity.Message{}).
		Where("receiver_id = ? AND sender_id = ? AND is_read = false", receiverID, senderID).
		Updates(map[string]interface{}{
			"is_read": true,
			"read_at": now,
		}).Error
}
