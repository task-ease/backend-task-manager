package mixins

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func ParseUserId(c *gin.Context) (uuid.UUID, error) {
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
