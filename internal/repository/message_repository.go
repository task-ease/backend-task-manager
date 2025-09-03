package repository

import (
	"backend-task-manager/internal/domain"
	"backend-task-manager/internal/entities"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type messageRepo struct {
	conn            *pgxpool.Pool
	imageStorageKey string
	chatRepo        domain.ChatRepository
}

func NewMessageRepository(conn *pgxpool.Pool, imageStorageKey string) domain.MessageRepository {
	return &messageRepo{conn: conn, imageStorageKey: imageStorageKey}
}

func (r *messageRepo) AddMessage(ctx context.Context, message entities.Message) error {
	return r.AddMessageTx(ctx, r.conn, message)
}

func (r *messageRepo) AddMessageTx(ctx context.Context, exec entities.Execer, message entities.Message) error {
	_, err := exec.Exec(ctx,
		`INSERT INTO messages (id, chat_id, sender_id, content, message_type, created_at, updated_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		message.ID, message.ChatID, message.SenderID, message.Content, message.MessageType, message.CreatedAt, message.UpdatedAt)
	return err
}

func (r *messageRepo) AddMessageReads(ctx context.Context, chatId, messageId string, senderId uuid.UUID) ([]uuid.UUID, error) {
	return r.AddMessageReadsTx(ctx, r.conn, chatId, messageId, senderId)
}

func (r *messageRepo) AddMessageReadsTx(ctx context.Context, exec entities.Execer, chatId, messageId string, senderId uuid.UUID) ([]uuid.UUID, error) {
	rows, err := exec.Query(ctx, `
		INSERT INTO message_reads (user_id, message_id, chat_id)
		SELECT uc.user_id, $1, $2
		FROM user_chats uc
		WHERE uc.chat_id = $2
			AND uc.workspace_id = (
				SELECT c.workspace_id
				FROM chats c
				WHERE c.id = $2
			)
			AND uc.user_id != $3
		RETURNING user_id;
	`, messageId, chatId, senderId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var unreadUserIds []uuid.UUID
	for rows.Next() {
		var userID uuid.UUID
		if err := rows.Scan(&userID); err != nil {
			return nil, err
		}
		unreadUserIds = append(unreadUserIds, userID)
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return unreadUserIds, nil
}

func (r *messageRepo) GetAllMessages(ctx context.Context, chatId string) ([]entities.Message, error) {
	rows, err := r.conn.Query(ctx, `
        SELECT 
            m.id,
            m.sender_id, 
            m.content, 
            m.message_type, 
            m.created_at, 
            m.is_read,
            ARRAY_AGG(mr.user_id) FILTER (WHERE mr.read_at IS NULL) AS unread_by,
            COALESCE(JSON_AGG(
                JSON_BUILD_OBJECT(
                    'id', ma.id,
                    'message_id', ma.message_id,
                    'uploaded_at', ma.uploaded_at,
                    'file_size', ma.file_size,
                    'file_name', ma.file_name,
                    'file_type', ma.file_type,
                    'file_url', ma.file_url
                )
            ) FILTER (WHERE ma.id IS NOT NULL), '[]') AS attachments
        FROM messages m
        JOIN message_reads mr ON mr.message_id = m.id
        LEFT JOIN message_attachments ma ON m.id = ma.message_id
        WHERE m.chat_id = $1
        GROUP BY m.id, m.sender_id, m.content, m.message_type, m.created_at, m.is_read
        ORDER BY m.created_at
    `, chatId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []entities.Message
	for rows.Next() {
		var message entities.Message
		var attachmentsJSON []byte

		if err := rows.Scan(
			&message.ID,
			&message.SenderID,
			&message.Content,
			&message.MessageType,
			&message.CreatedAt,
			&message.IsRead,
			&message.UnreadUsersIds,
			&attachmentsJSON,
		); err != nil {
			return nil, err
		}

		if len(attachmentsJSON) > 0 {
			var attachments []entities.Attachment
			if err := json.Unmarshal(attachmentsJSON, &attachments); err != nil {
				return nil, err
			}
			message.Attachments = attachments
		}

		messages = append(messages, message)
	}

	if err := rows.Err(); err != nil {
		return nil, err
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

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)

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

func (r *messageRepo) AddAttachment(ctx context.Context, attachment entities.Attachment) error {
	return r.AddAttachmentTx(ctx, r.conn, attachment)
}

func (r *messageRepo) AddAttachmentTx(ctx context.Context, exec entities.Execer, attachment entities.Attachment) error {
	_, err := exec.Exec(ctx, `
		INSERT INTO message_attachments 
		(id, message_id, file_url, file_type, file_name, file_size, uploaded_at, chat_id)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`, attachment.ID, attachment.MessageID, attachment.FileUrl, attachment.FileType, attachment.FileName, attachment.Size, attachment.UploadedAt, attachment.ChatID)
	return err
}

func (r *messageRepo) SetMessageRead(ctx context.Context, message, userId uuid.UUID) error {
	return r.SetMessageReadTx(ctx, r.conn, message, userId)
}

func (r *messageRepo) SetMessageReadTx(ctx context.Context, exec entities.Execer, message, userId uuid.UUID) error {
	_, err := exec.Exec(ctx, `
		UPDATE message_reads 
		SET read_at = $1 
		WHERE message_id = $2 AND user_id = $3`,
		time.Now().UTC(), message, userId)
	if err != nil {
		return err
	}

	var isRead bool

	err = exec.QueryRow(ctx, `
			SELECT COUNT(*) = 0
			FROM message_reads
			WHERE message_id = $1 AND read_at IS NULL`,
		message.ID).Scan(&isRead)

	if err != nil {
		return err
	}

	if isRead {
		_, err := exec.Exec(ctx, `
		UPDATE messages
		SET is_read = true
		WHERE id = $1
	`, message.ID)
		if err != nil {
			return err
		}
		return nil
	}

	return nil
}

func (r *messageRepo) GetAllChatImages(ctx context.Context, chatId string) (*[]entities.Attachment, error) {
	rows, err := r.conn.Query(ctx, `
		SELECT id, file_url, file_type, file_name, uploaded_at 
		FROM message_attachments 
		WHERE chat_id = $1
	`, chatId)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var images []entities.Attachment
	for rows.Next() {
		var attachment entities.Attachment
		if err := rows.Scan(&attachment.ID, &attachment.FileUrl, &attachment.FileType, &attachment.FileName, &attachment.UploadedAt); err != nil {
			return nil, err
		}
		images = append(images, attachment)
	}

	return &images, err
}
