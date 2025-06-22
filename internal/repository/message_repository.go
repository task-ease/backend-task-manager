package repository

import (
	"context"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"go-postgres-test/internal/domain"
)

type messageRepo struct {
	conn     *pgxpool.Pool
	chatRepo domain.ChatRepository
}

func NewMessageRepository(conn *pgxpool.Pool) domain.MessageRepository {
	return &messageRepo{conn: conn}
}

func (r *messageRepo) AddMessage(message *domain.Message) error {
	message.ID = "m-" + uuid.New().String()
	_, err := r.conn.Exec(context.Background(),
		`INSERT INTO messages (id, chat_id, sender_id, content, message_type, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7)`, message.ID, message.ChatID, message.SenderID, message.Content, message.MessageType, message.CreatedAt, message.UpdatedAt)
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
		messages = append(messages, &message)
	}
	return messages, nil
}
