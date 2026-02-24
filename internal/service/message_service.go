package service

import (
	"context"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/Rezann47/YksKoc/internal/domain/apperror"
	"github.com/Rezann47/YksKoc/internal/domain/dto"
	"github.com/Rezann47/YksKoc/internal/domain/entity"
	"github.com/Rezann47/YksKoc/internal/repository"
)

type MessageService interface {
	Send(ctx context.Context, senderID uuid.UUID, req dto.SendMessageReq) (*dto.MessageRes, error)
	GetConversation(ctx context.Context, userID, peerID uuid.UUID, page, limit int) (*dto.PaginatedRes[dto.MessageRes], error)
	ListConversations(ctx context.Context, userID uuid.UUID) ([]dto.ConversationRes, error)
	MarkRead(ctx context.Context, receiverID, senderID uuid.UUID) error
	UnreadCount(ctx context.Context, userID uuid.UUID) (*dto.UnreadCountRes, error)
}

type messageService struct {
	msgRepo repository.MessageRepository
	relRepo repository.InstructorStudentRepository
	log     *zap.Logger
}

func NewMessageService(
	msgRepo repository.MessageRepository,
	relRepo repository.InstructorStudentRepository,
	log *zap.Logger,
) MessageService {
	return &messageService{msgRepo: msgRepo, relRepo: relRepo, log: log}
}

// Send — sadece eğitmen-öğrenci ilişkisi varsa mesaj gönderilebilir
func (s *messageService) Send(ctx context.Context, senderID uuid.UUID, req dto.SendMessageReq) (*dto.MessageRes, error) {
	// İlişki kontrolü (her iki yön de geçerli)
	ok, err := s.isRelated(ctx, senderID, req.ReceiverID)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, apperror.NewForbidden("bu kullanıcıya mesaj gönderemezsin")
	}

	msg := &entity.Message{
		SenderID:   senderID,
		ReceiverID: req.ReceiverID,
		Content:    req.Content,
	}
	if err := s.msgRepo.Create(ctx, msg); err != nil {
		return nil, apperror.NewInternal(err)
	}

	res := mapMsg(msg, senderID)
	return &res, nil
}

func (s *messageService) GetConversation(ctx context.Context, userID, peerID uuid.UUID, page, limit int) (*dto.PaginatedRes[dto.MessageRes], error) {
	// İlişki kontrolü
	ok, err := s.isRelated(ctx, userID, peerID)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, apperror.NewForbidden("bu konuşmaya erişemezsin")
	}

	if limit <= 0 || limit > 50 {
		limit = 30
	}
	if page <= 0 {
		page = 1
	}
	offset := (page - 1) * limit

	msgs, total, err := s.msgRepo.ListConversation(ctx, userID, peerID, limit, offset)
	if err != nil {
		return nil, apperror.NewInternal(err)
	}

	res := make([]dto.MessageRes, len(msgs))
	for i, m := range msgs {
		res[i] = mapMsg(&m, userID)
	}

	paged := dto.NewPaginatedRes(res, total, page, limit)
	return &paged, nil
}

func (s *messageService) ListConversations(ctx context.Context, userID uuid.UUID) ([]dto.ConversationRes, error) {
	rows, err := s.msgRepo.ListConversations(ctx, userID)
	if err != nil {
		return nil, apperror.NewInternal(err)
	}

	res := make([]dto.ConversationRes, len(rows))
	for i, r := range rows {
		res[i] = dto.ConversationRes{
			PeerID:      r.PeerID,
			PeerName:    r.PeerName,
			PeerAvatar:  r.PeerAvatar,
			PeerOnline:  r.PeerOnline,
			LastMessage: r.LastMessage,
			LastAt:      r.LastAt,
			UnreadCount: r.UnreadCount,
		}
	}
	return res, nil
}

func (s *messageService) MarkRead(ctx context.Context, receiverID, senderID uuid.UUID) error {
	return s.msgRepo.MarkRead(ctx, receiverID, senderID)
}

func (s *messageService) UnreadCount(ctx context.Context, userID uuid.UUID) (*dto.UnreadCountRes, error) {
	count, err := s.msgRepo.UnreadCount(ctx, userID, nil)
	if err != nil {
		return nil, apperror.NewInternal(err)
	}
	return &dto.UnreadCountRes{Count: count}, nil
}

// ─── Yardımcılar ──────────────────────────────────────────

// isRelated: iki kullanıcı arasında eğitmen-öğrenci ilişkisi var mı (her iki yön)
func (s *messageService) isRelated(ctx context.Context, a, b uuid.UUID) (bool, error) {
	// a eğitmen, b öğrenci mi?
	ok, err := s.relRepo.Exists(ctx, a, b)
	if err != nil {
		return false, apperror.NewInternal(err)
	}
	if ok {
		return true, nil
	}
	// b eğitmen, a öğrenci mi?
	ok, err = s.relRepo.Exists(ctx, b, a)
	if err != nil {
		return false, apperror.NewInternal(err)
	}
	return ok, nil
}

func mapMsg(m *entity.Message, viewerID uuid.UUID) dto.MessageRes {
	return dto.MessageRes{
		ID:         m.ID,
		SenderID:   m.SenderID,
		ReceiverID: m.ReceiverID,
		Content:    m.Content,
		IsRead:     m.IsRead,
		ReadAt:     m.ReadAt,
		CreatedAt:  m.CreatedAt,
		IsMine:     m.SenderID == viewerID,
	}
}
