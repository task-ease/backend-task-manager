package handlers

import (
	"backend-task-manager/infrastructure/auth"
	"backend-task-manager/internal/domain/rules"
	"backend-task-manager/internal/dto/request"
	"backend-task-manager/internal/enums"
	"backend-task-manager/internal/middleware"
	"backend-task-manager/internal/usecase"
	"backend-task-manager/mixins"
	"net/http"

	"github.com/gin-gonic/gin"
)

type DocumentHandler struct {
	documentUC  *usecase.DocumentUsecase
	workspaceUC *usecase.WorkSpaceUsecase
	projectUC   *usecase.ProjectUseCase
}

func NewDocumentHandler(documentUC *usecase.DocumentUsecase, workspaceUC *usecase.WorkSpaceUsecase, projectUC *usecase.ProjectUseCase) *DocumentHandler {
	return &DocumentHandler{documentUC, workspaceUC, projectUC}
}

func (h *DocumentHandler) RegisterRoutes(router *gin.RouterGroup) {
	authService := auth.NewJWTService()

	protected := router.Group("/doc", middleware.JWTMiddleware(authService))

	protected.GET("/:documentId", middleware.AccessMiddleware(enums.ParamDocument, h.documentUC, rules.AllDocumentRoles), h.getDoc)
	protected.GET("/id/:workspaceId", middleware.AccessMiddleware(enums.ParamWorkspace, h.workspaceUC, rules.CanEditWorkspace), h.getIdByNameAndAndWorkspace)
	protected.GET("/all/:workspaceId", middleware.AccessMiddleware(enums.ParamWorkspace, h.workspaceUC, rules.CanEditWorkspace), h.getDocs)
	protected.GET("/access/:documentId", middleware.AccessMiddleware(enums.ParamDocument, h.documentUC, rules.AllDocumentRoles), h.checkUserAccess)
	protected.GET("/search/:workspaceId", middleware.AccessMiddleware(enums.ParamWorkspace, h.workspaceUC, rules.AllWorkspaceRoles), h.getDocsByName)

	protected.POST("/:workspaceId", middleware.AccessMiddleware(enums.ParamWorkspace, h.workspaceUC, rules.CanEditWorkspace), h.createDoc)

	protected.PATCH("/name/:documentId", middleware.AccessMiddleware(enums.ParamDocument, h.documentUC, rules.CanEdit), h.updateDocName)
	protected.PATCH("/parent/:documentId", middleware.AccessMiddleware(enums.ParamDocument, h.documentUC, rules.CanEdit), h.updateDocParent)
	protected.PATCH("/content/:documentId", middleware.AccessMiddleware(enums.ParamDocument, h.documentUC, rules.CanEdit), h.updateDocContent)
	protected.PATCH("/user/edit/:documentId", middleware.AccessMiddleware(enums.ParamDocument, h.documentUC, rules.CanEdit), h.updateDocUserEditPermissions)
	protected.PATCH("/visibility/:documentId", middleware.AccessMiddleware(enums.ParamDocument, h.documentUC, rules.CanEdit), h.updateDocVisibility)
}

func (h *DocumentHandler) checkUserAccess(c *gin.Context) {
	role, exists := c.Get("role")

	if !exists {
		c.Status(http.StatusUnauthorized)
		return
	}

	c.JSON(http.StatusOK, gin.H{"role": role})
}

func (h *DocumentHandler) createDoc(c *gin.Context) {
	workspaceId, err := mixins.ParamToUUID(c, string(enums.ParamWorkspace))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userId, err := mixins.ParseContextUserId(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var dto request.CreateDocsRequest
	if err = c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id, err := h.documentUC.CreateDoc(c, userId, workspaceId, dto)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"id": id})
}

func (h *DocumentHandler) updateDocName(c *gin.Context) {
	documentId, err := mixins.ParamToUUID(c, string(enums.ParamDocument))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	name := c.Query("name")

	if err = h.documentUC.UpdateDocumentName(c.Request.Context(), documentId, name); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}

func (h *DocumentHandler) updateDocParent(c *gin.Context) {
	parentId, err := mixins.QueryToUUIDCanBeNull(c, "parentId")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	documentId, err := mixins.ParamToUUID(c, string(enums.ParamDocument))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.documentUC.UpdateDocParent(c.Request.Context(), documentId, parentId); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}

func (h *DocumentHandler) updateDocContent(c *gin.Context) {
	userId, err := mixins.ParseContextUserId(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	documentId, err := mixins.ParamToUUID(c, string(enums.ParamDocument))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	content := c.Query("content")
	if err = h.documentUC.UpdateDocContent(c.Request.Context(), documentId, userId, content); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}

func (h *DocumentHandler) updateDocUserEditPermissions(c *gin.Context) {
	documentId, err := mixins.ParamToUUID(c, string(enums.ParamDocument))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var req request.UpdateDocUserEditPermissionsRequest
	if err = c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err = h.documentUC.UpdateDocUserEditPermissions(c.Request.Context(), documentId, req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}

func (h *DocumentHandler) updateDocVisibility(c *gin.Context) {
	documentId, err := mixins.ParamToUUID(c, string(enums.ParamDocument))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	to := c.Query("to")

	if err = h.documentUC.UpdateDocVisibility(c.Request.Context(), documentId, enums.DocumentVisibilityTypes(to)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}

func (h *DocumentHandler) getDocs(c *gin.Context) {
	workspaceId, err := mixins.ParamToUUID(c, string(enums.ParamWorkspace))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userId, err := mixins.ParseContextUserId(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	docs, err := h.documentUC.GetAllByUserAndWorkspaceId(c.Request.Context(), userId, workspaceId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"docs": docs})
}

func (h *DocumentHandler) getDoc(c *gin.Context) {
	documentId, err := mixins.ParamToUUID(c, string(enums.ParamDocument))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	document, err := h.documentUC.GetDocument(c.Request.Context(), documentId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"doc": document})
}

func (h *DocumentHandler) getDocsByName(c *gin.Context) {
	workspaceId, err := mixins.ParamToUUID(c, string(enums.ParamWorkspace))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userId, err := mixins.ParseContextUserId(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	name := c.Query("name")

	docs, err := h.documentUC.GetDocumentsByName(c.Request.Context(), userId, workspaceId, name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"docs": docs})
}

func (h *DocumentHandler) getIdByNameAndAndWorkspace(c *gin.Context) {
	workspaceId, err := mixins.ParamToUUID(c, string(enums.ParamWorkspace))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	name := c.Query("name")

	id, err := h.documentUC.GetIdByNameAndAndWorkspace(c.Request.Context(), name, workspaceId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"id": id})
}
