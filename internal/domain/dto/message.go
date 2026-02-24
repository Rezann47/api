package dto

import (
	"time"

	"github.com/google/uuid"
)

// ─── Request ──────────────────────────────────────────────

type SendMessageReq struct {
	ReceiverID uuid.UUID `json:"receiver_id" binding:"required"`
	Content    string    `json:"content"     binding:"required,min=1,max=2000"`
}

// ─── Response ─────────────────────────────────────────────

type MessageRes struct {
	ID         uuid.UUID  `json:"id"`
	SenderID   uuid.UUID  `json:"sender_id"`
	ReceiverID uuid.UUID  `json:"receiver_id"`
	Content    string     `json:"content"`
	IsRead     bool       `json:"is_read"`
	ReadAt     *time.Time `json:"read_at"`
	CreatedAt  time.Time  `json:"created_at"`
	IsMine     bool       `json:"is_mine"` // frontend kolaylığı için
}

// Konuşma listesi için — son mesajla birlikte
type ConversationRes struct {
	PeerID      uuid.UUID `json:"peer_id"`
	PeerName    string    `json:"peer_name"`
	PeerAvatar  int16     `json:"peer_avatar_id"`
	PeerOnline  bool      `json:"peer_online"`
	LastMessage string    `json:"last_message"`
	LastAt      time.Time `json:"last_at"`
	UnreadCount int64     `json:"unread_count"`
}

// Okunmamış toplam sayısı
type UnreadCountRes struct {
	Count int64 `json:"count"`
}
