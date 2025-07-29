package repository

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"go-postgres-test/internal/domain"
	"go-postgres-test/internal/entities"
	"go-postgres-test/mixins"
	"io"
	"net/http"
	"net/url"
	"time"
)

type messageRepo struct {
	conn            *pgxpool.Pool
	imageStorageKey string
	chatRepo        domain.ChatRepository
}

func NewMessageRepository(conn *pgxpool.Pool, imageStorageKey string) domain.MessageRepository {
	return &messageRepo{conn: conn, imageStorageKey: imageStorageKey}
}

func (r *messageRepo) AddMessage(message *domain.Message) (result *[]uuid.UUID, err error) {
	ctx := context.Background()

	tx, err := r.conn.Begin(ctx)
	if err != nil {
		empty := make([]uuid.UUID, 0)
		return &empty, err
	}

	defer func() {
		_ = mixins.TXReturn(tx, ctx, err)
	}()

	return r.AddMessageTx(context.Background(), tx, message)
}

func (r *messageRepo) AddMessageTx(ctx context.Context, exec entities.Execer, message *domain.Message) (*[]uuid.UUID, error) {
	message.ID = "m-" + uuid.New().String()

	_, err := exec.Exec(ctx,
		`INSERT INTO messages (id, chat_id, sender_id, content, message_type, created_at, updated_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		message.ID, message.ChatID, message.SenderID, message.Content, message.MessageType, message.CreatedAt, message.UpdatedAt)
	if err != nil {
		return nil, err
	}

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
	`, message.ID, message.ChatID, message.SenderID)
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

	_, err = exec.Exec(ctx,
		`UPDATE chats SET updated_at = $1 WHERE id = $2`, message.UpdatedAt, message.ChatID)
	if err != nil {
		return nil, err
	}

	return &unreadUserIds, nil
}

func (r *messageRepo) GetAllMessages(chatId string) ([]*domain.Message, error) {
	ctx := context.Background()

	tx, err := r.conn.Begin(ctx)
	if err != nil {
		empty := make([]*domain.Message, 0)
		return empty, err
	}

	defer func() {
		_ = mixins.TXReturn(tx, ctx, err)
	}()

	return r.GetAllMessagesTx(ctx, tx, chatId)
}

func (r *messageRepo) GetAllMessagesTx(ctx context.Context, exec entities.Execer, chatId string) ([]*domain.Message, error) {
	rows, err := exec.Query(ctx, `
		SELECT 
		    m.id,
		    m.sender_id, 
		    m.content, 
		    m.message_type, 
		    m.created_at, 
			m.is_read,
			ARRAY_AGG(mr.user_id) FILTER (WHERE mr.read_at IS NULL) AS unread_by
		FROM messages m
		JOIN message_reads mr ON mr.message_id = m.id		
		WHERE m.chat_id = $1
		GROUP BY 
    		m.id, 
    		m.sender_id, 
    		m.content, 
    		m.message_type, 
    		m.created_at, 
    		m.is_read
		ORDER BY m.created_at
		`, chatId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []*domain.Message
	for rows.Next() {
		var message domain.Message
		if err = rows.Scan(&message.ID, &message.SenderID, &message.Content, &message.MessageType, &message.CreatedAt, &message.IsRead, &message.UnreadUsersIds); err != nil {
			return nil, err
		}

		if err := r.GetAttachmentsTx(ctx, exec, message.ID, &message.Attachments); err != nil {
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

func (r *messageRepo) GetAttachments(messageId string, attachments *[]domain.Attachment) error {
	return r.GetAttachmentsTx(context.Background(), r.conn, messageId, attachments)
}

func (r *messageRepo) GetAttachmentsTx(ctx context.Context, exec entities.Execer, messageId string, attachments *[]domain.Attachment) error {
	rows, err := exec.Query(ctx, `
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

func (r *messageRepo) SetMessageRead(message *domain.Message, userId uuid.UUID) error {
	ctx := context.Background()

	tx, err := r.conn.Begin(ctx)
	if err != nil {
		return err
	}

	defer func() {
		_ = mixins.TXReturn(tx, ctx, err)
	}()

	return r.SetMessageReadTx(ctx, tx, message, userId)
}

func (r *messageRepo) SetMessageReadTx(ctx context.Context, exec entities.Execer, message *domain.Message, userId uuid.UUID) error {
	_, err := exec.Exec(ctx, `
		UPDATE message_reads 
		SET read_at = $1 
		WHERE message_id = $2 AND user_id = $3`,
		time.Now().UTC(), message.ID, userId)
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

func (r *messageRepo) SetMessagesRead(messages *[]domain.Message, userId uuid.UUID) error {
	ctx := context.Background()

	tx, err := r.conn.Begin(ctx)
	if err != nil {
		return err
	}

	defer func() {
		_ = mixins.TXReturn(tx, ctx, err)
	}()

	for _, m := range *messages {
		err := r.SetMessageReadTx(ctx, tx, &m, userId)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *messageRepo) GetAllChatImages(chatId string) (*[]domain.Attachment, error) {
	rows, err := r.conn.Query(context.Background(), `
		SELECT id, file_url, file_type, file_name, uploaded_at from  message_attachments 
		WHERE chat_id = $1`, chatId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var images []domain.Attachment
	for rows.Next() {
		var attachment domain.Attachment
		if err := rows.Scan(&attachment.ID, &attachment.FileUrl, &attachment.FileType, &attachment.FileName, &attachment.UploadedAt); err != nil {
			return nil, err
		}
		images = append(images, attachment)
	}

	return &images, err
}
