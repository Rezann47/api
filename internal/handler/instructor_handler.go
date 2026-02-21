package handler

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/Rezann47/YksKoc/internal/domain/dto"
	"github.com/Rezann47/YksKoc/internal/domain/entity"
	"github.com/Rezann47/YksKoc/internal/middleware"
	"github.com/Rezann47/YksKoc/internal/service"
	"github.com/Rezann47/YksKoc/pkg/response"
)

type InstructorHandler struct{ svc service.InstructorService }

func NewInstructorHandler(svc service.InstructorService) *InstructorHandler {
	return &InstructorHandler{svc: svc}
}

func (h *InstructorHandler) AddStudent(c *gin.Context) {
	var req dto.AddStudentReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err.Error())
		return
	}
	if err := h.svc.AddStudent(c.Request.Context(), middleware.GetUserID(c), req); err != nil {
		response.Error(c, err)
		return
	}
	response.NoContent(c)
}

func (h *InstructorHandler) ListStudents(c *gin.Context) {
	pagination := dto.ParsePagination(c)
	res, err := h.svc.ListStudents(c.Request.Context(), middleware.GetUserID(c), pagination)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.OK(c, res)
}

func (h *InstructorHandler) RemoveStudent(c *gin.Context) {
	studentID, err := uuid.Parse(c.Param("studentID"))
	if err != nil {
		response.ValidationError(c, "geçersiz student ID")
		return
	}
	if err := h.svc.RemoveStudent(c.Request.Context(), middleware.GetUserID(c), studentID); err != nil {
		response.Error(c, err)
		return
	}
	response.NoContent(c)
}

func (h *InstructorHandler) GetStudentPomodoros(c *gin.Context) {
	studentID, _ := uuid.Parse(c.Param("studentID"))
	from, _ := time.Parse("2006-01-02", c.DefaultQuery("from", time.Now().AddDate(0, 0, -7).Format("2006-01-02")))
	to, _ := time.Parse("2006-01-02", c.DefaultQuery("to", time.Now().Format("2006-01-02")))
	to = to.Add(24*time.Hour - time.Second)

	res, err := h.svc.GetStudentPomodoros(c.Request.Context(), middleware.GetUserID(c), studentID, from, to)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.OK(c, res)
}

func (h *InstructorHandler) GetStudentProgress(c *gin.Context) {
	studentID, _ := uuid.Parse(c.Param("studentID"))
	res, err := h.svc.GetStudentProgress(c.Request.Context(), middleware.GetUserID(c), studentID)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.OK(c, res)
}

func (h *InstructorHandler) GetStudentExamResults(c *gin.Context) {
	studentID, _ := uuid.Parse(c.Param("studentID"))
	pagination := dto.ParsePagination(c)

	var examType *entity.ExamType
	if et := c.Query("exam_type"); et != "" {
		t := entity.ExamType(et)
		examType = &t
	}

	res, err := h.svc.GetStudentExamResults(c.Request.Context(), middleware.GetUserID(c), studentID, examType, pagination)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.OK(c, res)
}
