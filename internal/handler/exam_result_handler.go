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

type ExamResultHandler struct{ svc service.ExamResultService }

func NewExamResultHandler(svc service.ExamResultService) *ExamResultHandler {
	return &ExamResultHandler{svc: svc}
}

func (h *ExamResultHandler) Create(c *gin.Context) {
	var req dto.CreateExamResultReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err.Error())
		return
	}
	res, err := h.svc.Create(c.Request.Context(), middleware.GetUserID(c), req)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.Created(c, res)
}

func (h *ExamResultHandler) List(c *gin.Context) {
	pagination := dto.ParsePagination(c)
	var examType *entity.ExamType
	if et := c.Query("exam_type"); et != "" {
		t := entity.ExamType(et)
		examType = &t
	}
	res, err := h.svc.List(c.Request.Context(), middleware.GetUserID(c), examType, pagination)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.OK(c, res)
}

func (h *ExamResultHandler) GetStats(c *gin.Context) {
	examTypeStr := c.DefaultQuery("exam_type", "TYT")
	res, err := h.svc.GetStats(c.Request.Context(), middleware.GetUserID(c), entity.ExamType(examTypeStr))
	if err != nil {
		response.Error(c, err)
		return
	}
	response.OK(c, res)
}

func (h *ExamResultHandler) Delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.ValidationError(c, "geçersiz ID")
		return
	}
	if err := h.svc.Delete(c.Request.Context(), id, middleware.GetUserID(c)); err != nil {
		response.Error(c, err)
		return
	}
	response.NoContent(c)
}
