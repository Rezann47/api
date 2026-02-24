package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/Rezann47/YksKoc/internal/domain/apperror"
	"github.com/Rezann47/YksKoc/internal/domain/entity"
	"github.com/Rezann47/YksKoc/internal/middleware"
	"github.com/Rezann47/YksKoc/internal/service"
	"github.com/Rezann47/YksKoc/pkg/response"
)

type StudyPlanHandler struct {
	svc service.StudyPlanService
}

func NewStudyPlanHandler(svc service.StudyPlanService) *StudyPlanHandler {
	return &StudyPlanHandler{svc: svc}
}

// POST /study-plans — öğrenci kendi planını ekler
func (h *StudyPlanHandler) Create(c *gin.Context) {
	userID := middleware.GetUserID(c)

	var req struct {
		Title    string  `json:"title"`
		PlanDate string  `json:"plan_date" binding:"required"` // "2026-02-22"
		Note     *string `json:"note"`
		Items    []struct {
			SubjectID       string  `json:"subject_id" binding:"required"`
			TopicID         *string `json:"topic_id"`
			DurationMinutes int     `json:"duration_minutes"`
			DisplayOrder    int16   `json:"display_order"`
		} `json:"items"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err.Error())
		return
	}

	date, err := time.Parse("2006-01-02", req.PlanDate)
	if err != nil {
		response.ValidationError(c, "plan_date formatı geçersiz (YYYY-MM-DD)")
		return
	}

	items := make([]service.CreateStudyPlanItemInput, 0, len(req.Items))
	for _, it := range req.Items {
		sid, err := uuid.Parse(it.SubjectID)
		if err != nil {
			response.ValidationError(c, "geçersiz subject_id")
			return
		}
		var topicID *uuid.UUID
		if it.TopicID != nil {
			tid, err := uuid.Parse(*it.TopicID)
			if err != nil {
				response.ValidationError(c, "geçersiz topic_id")
				return
			}
			topicID = &tid
		}
		items = append(items, service.CreateStudyPlanItemInput{
			SubjectID:       sid,
			TopicID:         topicID,
			DurationMinutes: it.DurationMinutes,
			DisplayOrder:    it.DisplayOrder,
		})
	}

	plan, err := h.svc.Create(c.Request.Context(), service.CreateStudyPlanInput{
		UserID:    userID,
		CreatedBy: userID,
		Title:     req.Title,
		PlanDate:  date,
		Note:      req.Note,
		Items:     items,
	})
	if err != nil {
		response.Error(c, err)
		return
	}
	c.JSON(http.StatusCreated, gin.H{"success": true, "data": toStudyPlanDTO(plan)})
}

// POST /instructor/students/:studentID/study-plans — koç öğrencisine plan ekler
func (h *StudyPlanHandler) CreateForStudent(c *gin.Context) {
	instructorID := middleware.GetUserID(c)
	studentIDStr := c.Param("studentID")
	studentID, err := uuid.Parse(studentIDStr)
	if err != nil {
		response.Error(c, apperror.NewBadRequest("geçersiz öğrenci ID"))
		return
	}

	var req struct {
		Title    string  `json:"title"`
		PlanDate string  `json:"plan_date" binding:"required"`
		Note     *string `json:"note"`
		Items    []struct {
			SubjectID       string  `json:"subject_id" binding:"required"`
			TopicID         *string `json:"topic_id"`
			DurationMinutes int     `json:"duration_minutes"`
			DisplayOrder    int16   `json:"display_order"`
		} `json:"items"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err.Error())
		return
	}

	date, err := time.Parse("2006-01-02", req.PlanDate)
	if err != nil {
		response.ValidationError(c, "plan_date formatı geçersiz (YYYY-MM-DD)")
		return
	}

	items := make([]service.CreateStudyPlanItemInput, 0, len(req.Items))
	for _, it := range req.Items {
		sid, err := uuid.Parse(it.SubjectID)
		if err != nil {
			response.ValidationError(c, "geçersiz subject_id")
			return
		}
		var topicID *uuid.UUID
		if it.TopicID != nil {
			tid, err := uuid.Parse(*it.TopicID)
			if err != nil {
				response.ValidationError(c, "geçersiz topic_id")
				return
			}
			topicID = &tid
		}
		items = append(items, service.CreateStudyPlanItemInput{
			SubjectID:       sid,
			TopicID:         topicID,
			DurationMinutes: it.DurationMinutes,
			DisplayOrder:    it.DisplayOrder,
		})
	}

	plan, err := h.svc.Create(c.Request.Context(), service.CreateStudyPlanInput{
		UserID:    studentID,
		CreatedBy: instructorID,
		Title:     req.Title,
		PlanDate:  date,
		Note:      req.Note,
		Items:     items,
	})
	if err != nil {
		response.Error(c, err)
		return
	}
	c.JSON(http.StatusCreated, gin.H{"success": true, "data": toStudyPlanDTO(plan)})
}

// GET /study-plans?date=2026-02-22
func (h *StudyPlanHandler) ListByDate(c *gin.Context) {
	userID := middleware.GetUserID(c)
	dateStr := c.Query("date")
	if dateStr == "" {
		dateStr = time.Now().Format("2006-01-02")
	}
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		response.ValidationError(c, "date formatı geçersiz (YYYY-MM-DD)")
		return
	}
	plans, err := h.svc.ListByDate(c.Request.Context(), userID, date)
	if err != nil {
		response.Error(c, err)
		return
	}
	dtos := make([]studyPlanDTO, 0, len(plans))
	for _, p := range plans {
		dtos = append(dtos, toStudyPlanDTO(p))
	}
	response.OK(c, dtos)
}

// GET /study-plans/month?year=2026&month=2
func (h *StudyPlanHandler) ListByMonth(c *gin.Context) {
	userID := middleware.GetUserID(c)
	now := time.Now()
	year := now.Year()
	month := int(now.Month())

	if y := c.Query("year"); y != "" {
		if v, err := strconv.Atoi(y); err == nil {
			year = v
		}
	}
	if m := c.Query("month"); m != "" {
		if v, err := strconv.Atoi(m); err == nil {
			month = v
		}
	}

	plans, err := h.svc.ListByMonth(c.Request.Context(), userID, year, month)
	if err != nil {
		response.Error(c, err)
		return
	}
	dtos := make([]studyPlanDTO, 0, len(plans))
	for _, p := range plans {
		dtos = append(dtos, toStudyPlanDTO(p))
	}
	response.OK(c, dtos)
}

// GET /instructor/students/:studentID/study-plans?date=
func (h *StudyPlanHandler) GetStudentPlans(c *gin.Context) {
	studentIDStr := c.Param("studentID")
	studentID, err := uuid.Parse(studentIDStr)
	if err != nil {
		response.Error(c, apperror.NewBadRequest("geçersiz öğrenci ID"))
		return
	}
	dateStr := c.Query("date")
	if dateStr == "" {
		dateStr = time.Now().Format("2006-01-02")
	}
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		response.ValidationError(c, "date formatı geçersiz")
		return
	}
	plans, err := h.svc.ListByDate(c.Request.Context(), studentID, date)
	if err != nil {
		response.Error(c, err)
		return
	}
	dtos := make([]studyPlanDTO, 0, len(plans))
	for _, p := range plans {
		dtos = append(dtos, toStudyPlanDTO(p))
	}
	response.OK(c, dtos)
}

// DELETE /study-plans/:id
func (h *StudyPlanHandler) Delete(c *gin.Context) {
	userID := middleware.GetUserID(c)
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, apperror.NewBadRequest("geçersiz plan ID"))
		return
	}
	if err := h.svc.Delete(c.Request.Context(), id, userID); err != nil {
		response.Error(c, err)
		return
	}
	response.NoContent(c)
}

// PATCH /study-plans/:id/items/:itemID/complete
func (h *StudyPlanHandler) CompleteItem(c *gin.Context) {
	userID := middleware.GetUserID(c)
	planID, _ := uuid.Parse(c.Param("id"))
	itemID, _ := uuid.Parse(c.Param("itemID"))
	if err := h.svc.CompleteItem(c.Request.Context(), planID, itemID, userID); err != nil {
		response.Error(c, err)
		return
	}
	response.OK(c, gin.H{"completed": true})
}

// PATCH /study-plans/:id/items/:itemID/uncomplete
func (h *StudyPlanHandler) UncompleteItem(c *gin.Context) {
	userID := middleware.GetUserID(c)
	planID, _ := uuid.Parse(c.Param("id"))
	itemID, _ := uuid.Parse(c.Param("itemID"))
	if err := h.svc.UncompleteItem(c.Request.Context(), planID, itemID, userID); err != nil {
		response.Error(c, err)
		return
	}
	response.OK(c, gin.H{"completed": false})
}

// ─── DTO ─────────────────────────────────────────────────

type studyPlanItemDTO struct {
	ID              string  `json:"id"`
	SubjectID       string  `json:"subject_id"`
	SubjectName     string  `json:"subject_name"`
	TopicID         *string `json:"topic_id"`
	TopicName       *string `json:"topic_name"`
	DurationMinutes int     `json:"duration_minutes"`
	DisplayOrder    int16   `json:"display_order"`
	IsCompleted     bool    `json:"is_completed"`
	ExampType       string  `json:"exam_type"`
}

type studyPlanDTO struct {
	ID          string             `json:"id"`
	UserID      string             `json:"user_id"`
	CreatedBy   string             `json:"created_by"`
	CreatorName string             `json:"creator_name"`
	Title       string             `json:"title"`
	PlanDate    string             `json:"plan_date"`
	Note        *string            `json:"note"`
	Items       []studyPlanItemDTO `json:"items"`
}

func toStudyPlanDTO(p *entity.StudyPlan) studyPlanDTO {
	items := make([]studyPlanItemDTO, 0, len(p.Items))
	for _, it := range p.Items {
		dto := studyPlanItemDTO{
			ID:              it.ID.String(),
			SubjectID:       it.SubjectID.String(),
			SubjectName:     it.Subject.Name,
			DurationMinutes: it.DurationMinutes,
			ExampType:       string(it.Subject.ExamType),

			DisplayOrder: it.DisplayOrder,
			IsCompleted:  it.IsCompleted,
		}
		if it.TopicID != nil {
			s := it.TopicID.String()
			dto.TopicID = &s
		}
		if it.Topic != nil {
			dto.TopicName = &it.Topic.Name
		}
		items = append(items, dto)
	}
	return studyPlanDTO{
		ID:          p.ID.String(),
		UserID:      p.UserID.String(),
		CreatedBy:   p.CreatedBy.String(),
		CreatorName: p.Creator.Name,
		Title:       p.Title,
		PlanDate:    p.PlanDate.Format("2006-01-02"),
		Note:        p.Note,
		Items:       items,
	}
}
