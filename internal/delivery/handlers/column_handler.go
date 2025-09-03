package handlers

import (
	"backend-task-manager/infrastructure/auth"
	"backend-task-manager/internal/dto/request"
	"backend-task-manager/internal/enums"
	"backend-task-manager/internal/middleware"
	"backend-task-manager/internal/usecase"
	"backend-task-manager/mixins"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ColumnHandler struct {
	uc *usecase.ColumnUsecase
}

func NewColumnHandler(uc *usecase.ColumnUsecase) *ColumnHandler {
	return &ColumnHandler{uc}
}

func (h *ColumnHandler) RegisterRoutes(router *gin.RouterGroup) {
	authService := auth.NewJWTService()

	protected := router.Group("/columns", middleware.JWTMiddleware(authService))

	protected.GET("/", h.getColumns)
	protected.GET("/templates/:workspaceId", h.getAllColumnTemplates)

	protected.POST("/template", h.createColumnTemplate)

	protected.PUT("/renumber/:workspaceId", h.renumberColumnTemplate)

	protected.PATCH("/template/value", h.updateColumnTemplateValue)
	protected.PATCH("/template/status", h.updateColumnTemplateStatus)
}

func (h *ColumnHandler) getColumns(c *gin.Context) {
	workspaceId, err := mixins.QueryToUUID(c, "workspaceId")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	projectId, err := mixins.QueryToUUIDCanBeNull(c, "projectId")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	sprintId, err := mixins.QueryToUUIDCanBeNull(c, "sprintId")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	columns, err := h.uc.GetColumns(c.Request.Context(), workspaceId, projectId, sprintId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"columns": columns})
}

func (h *ColumnHandler) createColumnTemplate(c *gin.Context) {
	var input request.CreateNewColumnTemplate
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id, err := h.uc.CreateColumnTemplate(c.Request.Context(), input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": id})
}

func (h *ColumnHandler) getAllColumnTemplates(c *gin.Context) {
	workspaceId, err := mixins.ParamToUUID(c, "workspaceId")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	columns, err := h.uc.GetAllColumnTemplates(c.Request.Context(), workspaceId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"columns": columns})
}

func (h *ColumnHandler) updateColumnTemplateValue(c *gin.Context) {
	var input request.UpdateColumnTemplate
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var err error
	switch input.Method {
	case enums.ChangeColumnTemplateName:
		err = h.uc.UpdateColumnTemplateName(c.Request.Context(), input.ColumnTemplateId, input.Value)
	case enums.ChangeColumnTemplateColor:
		err = h.uc.UpdateColumnTemplateColor(c.Request.Context(), input.ColumnTemplateId, input.Value)
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}

func (h *ColumnHandler) updateColumnTemplateStatus(c *gin.Context) {
	var input request.UpdateColumnTemplateStatus
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var err error
	switch input.Method {
	case enums.ChangeColumnTemplateRequired:
		err = h.uc.UpdateColumnTemplateStatusRequired(c.Request.Context(), input.ColumnTemplateId, input.Value)
	case enums.ChangeColumnTemplateActive:
		err = h.uc.UpdateColumnTemplateStatusActive(c.Request.Context(), input.ColumnTemplateId, input.Value)
	case enums.ChangeColumnTemplateDone:
		err = h.uc.UpdateColumnTemplateStatusDone(c.Request.Context(), input.ColumnTemplateId)
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}

func (h *ColumnHandler) renumberColumnTemplate(c *gin.Context) {
	workspaceId, err := mixins.ParamToUUID(c, "workspaceId")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.uc.RenumberColumnTemplatesPositions(c.Request.Context(), workspaceId); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}
