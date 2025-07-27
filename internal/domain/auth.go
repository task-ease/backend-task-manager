package domain

import "github.com/google/uuid"

type AuthService interface {
	GenerateToken(userID string) (string, error)
	VerifyToken(tokenString string) (string, error)
}

type roleChecker interface {
	CheckWorkspaceRole(userId, workspaceId uuid.UUID) (bool, error)
	CheckProjectRole(userId, projectId uuid.UUID) (bool, error)
}
