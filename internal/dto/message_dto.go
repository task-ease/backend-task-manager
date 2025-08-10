package dto

import (
	"mime/multipart"

	"github.com/google/uuid"
)

type UploadImage struct {
	Form    *multipart.Form `json:"form"`
	ChatId  string          `json:"chatId"`
	UserId  uuid.UUID       `json:"userId"`
	Content string          `json:"content"`
}
