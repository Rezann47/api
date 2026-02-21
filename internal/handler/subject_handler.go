package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/Rezann47/YksKoc/internal/domain/dto"
	"github.com/Rezann47/YksKoc/internal/domain/entity"
	"github.com/Rezann47/YksKoc/internal/middleware"
	"github.com/Rezann47/YksKoc/internal/service"
	"github.com/Rezann47/YksKoc/pkg/response"
)

type SubjectHandler struct{ svc service.SubjectService }

func NewSubjectHandler(svc service.SubjectService) *SubjectHandler { return &SubjectHandler{svc: svc} }

func (h *SubjectHandler) ListSubjects(c *gin.Context) {
	var examType *entity.ExamType
	if et := c.Query("exam_type"); et != "" {
		t := entity.ExamType(et)
		examType = &t
	}
	res, err := h.svc.ListSubjects(c.Request.Context(), examType)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.OK(c, res)
}

func (h *SubjectHandler) ListTopics(c *gin.Context) {
	subjectID, err := uuid.Parse(c.Param("subjectID"))
	if err != nil {
		response.ValidationError(c, "geçersiz subject ID")
		return
	}
	res, err := h.svc.ListTopics(c.Request.Context(), subjectID, middleware.GetUserID(c))
	if err != nil {
		response.Error(c, err)
		return
	}
	response.OK(c, res)
}

func (h *SubjectHandler) MarkTopic(c *gin.Context) {
	topicID, err := uuid.Parse(c.Param("topicID"))
	if err != nil {
		response.ValidationError(c, "geçersiz topic ID")
		return
	}
	var req dto.MarkTopicReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err.Error())
		return
	}
	if err := h.svc.MarkTopic(c.Request.Context(), middleware.GetUserID(c), topicID, req.IsCompleted); err != nil {
		response.Error(c, err)
		return
	}
	response.NoContent(c)
}

func (h *SubjectHandler) GetSubjectProgress(c *gin.Context) {
	subjectID, err := uuid.Parse(c.Param("subjectID"))
	if err != nil {
		response.ValidationError(c, "geçersiz subject ID")
		return
	}
	res, err := h.svc.GetSubjectProgress(c.Request.Context(), middleware.GetUserID(c), subjectID)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.OK(c, res)
}

func (h *SubjectHandler) GetAllProgress(c *gin.Context) {
	res, err := h.svc.GetAllProgress(c.Request.Context(), middleware.GetUserID(c))
	if err != nil {
		response.Error(c, err)
		return
	}
	response.OK(c, res)
}
