package mixins

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func ParseContextUserId(c *gin.Context) (uuid.UUID, error) {
	userIdStr, exists := c.Get("userId")

	if !exists {
		return uuid.Nil, errors.New("no userId in context")
	}

	userId, err := uuid.Parse(userIdStr.(string))

	if err != nil {
		return uuid.Nil, err
	}
	return userId, nil
}

func ParseContextCanEdit(c *gin.Context) (bool, error) {
	canEditAny, exists := c.Get("canEdit")
	if !exists {
		return false, errors.New("no canEdit in context")
	}
	canEdit, ok := canEditAny.(bool)
	if !ok {
		return false, errors.New("canEdit must be a boolean")
	}
	return canEdit, nil
}

func QueryToUUID(c *gin.Context, name string) (uuid.UUID, error) {
	valueStr := c.Query(name)
	value, err := uuid.Parse(valueStr)
	if err != nil {
		return uuid.Nil, err
	}
	return value, nil
}

func QueryToUUIDCanBeNull(c *gin.Context, name string) (*uuid.UUID, error) {
	valueStr := c.Query(name)
	if valueStr == "" || valueStr == "null" {
		return nil, nil
	}
	value, err := uuid.Parse(valueStr)
	if err != nil {
		return nil, err
	}
	return &value, nil
}

func ParamToUUID(c *gin.Context, name string) (uuid.UUID, error) {
	value, err := uuid.Parse(c.Param(name))
	if err != nil {
		return uuid.Nil, err
	}
	return value, nil
}
