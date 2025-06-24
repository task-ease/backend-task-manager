package repository

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"go-postgres-test/internal/domain"
	"io"
	"net/http"
	"net/url"
)

type messageRepo struct {
	conn            *pgxpool.Pool
	imageStorageKey string
	chatRepo        domain.ChatRepository
}

func NewMessageRepository(conn *pgxpool.Pool, imageStorageKey string) domain.MessageRepository {
	return &messageRepo{conn: conn, imageStorageKey: imageStorageKey}
}

func (r *messageRepo) AddMessage(message *domain.Message) error {
	message.ID = "m-" + uuid.New().String()
	_, err := r.conn.Exec(context.Background(),
		`INSERT INTO messages (id, chat_id, sender_id, content, message_type, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7)`, message.ID, message.ChatID, message.SenderID, message.Content, message.MessageType, message.CreatedAt, message.UpdatedAt)
	if err != nil {
		return err
	}
	_, err = r.conn.Exec(context.Background(),
		`UPDATE chats SET updated_at = $1 WHERE id = $2`, message.UpdatedAt, message.ChatID)
	return err
}

func (r *messageRepo) GetAllMessages(chatId string) ([]*domain.Message, error) {
	rows, err := r.conn.Query(context.Background(), `
		SELECT id, sender_id, content, message_type, created_at
		FROM messages
		WHERE chat_id = $1`, chatId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []*domain.Message
	for rows.Next() {
		var message domain.Message
		if err = rows.Scan(&message.ID, &message.SenderID, &message.Content, &message.MessageType, &message.CreatedAt); err != nil {
			return nil, err
		}

		if err := r.getAttachments(message.ID, &message.Attachments); err != nil {
			message.Attachments = nil
		}

		messages = append(messages, &message)
	}
	return messages, nil
}

func (r *messageRepo) UploadImage(image io.Reader) (string, error) {
	imgBytes, err := io.ReadAll(image)
	if err != nil {
		return "", err
	}

	encoded := base64.StdEncoding.EncodeToString(imgBytes)

	data := url.Values{}
	data.Set("key", r.imageStorageKey)
	data.Set("image", encoded)

	resp, err := http.PostForm("https://api.imgbb.com/1/upload", data)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	var result struct {
		Success bool `json:"success"`
		Data    struct {
			URL string `json:"url"`
		} `json:"data"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return "", err
	}
	if !result.Success {
		return "", errors.New("imgbb upload failed")
	}

	return result.Data.URL, nil
}

func (r *messageRepo) AddAttachment(attachment *domain.Attachment) error {
	_, err := r.conn.Exec(context.Background(), `
		INSERT INTO message_attachments 
		(id, message_id, file_url, file_type, file_name, file_size, uploaded_at, chat_id)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
		attachment.ID, attachment.MessageID, attachment.FileUrl, attachment.FileType, attachment.FileName, attachment.Size, attachment.UploadedAt, attachment.ChatID)
	return err
}

func (r *messageRepo) getAttachments(messageId string, attachments *[]domain.Attachment) error {
	rows, err := r.conn.Query(context.Background(), `
		SELECT id, message_id, uploaded_at, file_size, file_name, file_type, file_url
		FROM message_attachments
		WHERE message_id = $1`, messageId)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var a domain.Attachment
		err := rows.Scan(&a.ID, &a.MessageID, &a.UploadedAt, &a.FileSize, &a.FileName, &a.FileType, &a.FileUrl)
		if err != nil {
			return err
		}
		*attachments = append(*attachments, a)
	}

	if err = rows.Err(); err != nil {
		return err
	}

	return nil
}
