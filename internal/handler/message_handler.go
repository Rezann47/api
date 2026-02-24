package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/Rezann47/YksKoc/internal/domain/dto"
	"github.com/Rezann47/YksKoc/internal/middleware"
	"github.com/Rezann47/YksKoc/internal/service"
	"github.com/Rezann47/YksKoc/pkg/response"
)

type MessageHandler struct{ svc service.MessageService }

func NewMessageHandler(svc service.MessageService) *MessageHandler {
	return &MessageHandler{svc: svc}
}

// POST /messages
func (h *MessageHandler) Send(c *gin.Context) {
	var req dto.SendMessageReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err.Error())
		return
	}
	res, err := h.svc.Send(c.Request.Context(), middleware.GetUserID(c), req)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.Created(c, res)
}

// GET /messages/conversations
func (h *MessageHandler) ListConversations(c *gin.Context) {
	res, err := h.svc.ListConversations(c.Request.Context(), middleware.GetUserID(c))
	if err != nil {
		response.Error(c, err)
		return
	}
	response.OK(c, res)
}

// GET /messages/conversations/:peerID
func (h *MessageHandler) GetConversation(c *gin.Context) {
	peerID, err := uuid.Parse(c.Param("peerID"))
	if err != nil {
		response.ValidationError(c, "geçersiz peer ID")
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "30"))

	res, err := h.svc.GetConversation(c.Request.Context(), middleware.GetUserID(c), peerID, page, limit)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.OK(c, res)
}

// POST /messages/conversations/:peerID/read
func (h *MessageHandler) MarkRead(c *gin.Context) {
	peerID, err := uuid.Parse(c.Param("peerID"))
	if err != nil {
		response.ValidationError(c, "geçersiz peer ID")
		return
	}
	if err := h.svc.MarkRead(c.Request.Context(), middleware.GetUserID(c), peerID); err != nil {
		response.Error(c, err)
		return
	}
	response.NoContent(c)
}

// GET /messages/unread
func (h *MessageHandler) UnreadCount(c *gin.Context) {
	res, err := h.svc.UnreadCount(c.Request.Context(), middleware.GetUserID(c))
	if err != nil {
		response.Error(c, err)
		return
	}
	response.OK(c, res)
}
