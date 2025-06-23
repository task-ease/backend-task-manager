package repository

import (
	"context"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"go-postgres-test/internal/domain"
	"go-postgres-test/internal/response"
	"go-postgres-test/internal/types"
	"time"
)

type chatRepo struct {
	conn        *pgxpool.Pool
	messageRepo domain.MessageRepository
}

func NewChatRepo(conn *pgxpool.Pool, messageRepo domain.MessageRepository) domain.ChatRepository {
	return &chatRepo{
		conn:        conn,
		messageRepo: messageRepo,
	}
}

func (r *chatRepo) CreateChat(chat *domain.Chat, participantId uuid.UUID) error {
	if chat.CreatorID.String() < participantId.String() {
		chat.ID = "c-" + chat.CreatorID.String() + "&" + participantId.String()
	} else {
		chat.ID = "c-" + participantId.String() + "&" + chat.CreatorID.String()
	}

	_, err := r.conn.Exec(
		context.Background(),
		`INSERT INTO chats (id, workspace_id, creator_id, type, created_at, updated_at, last_message_time)
			VALUES ($1, $2, $3, $4, $5, $6, $7)`, chat.ID, chat.WorkspaceID, chat.CreatorID, chat.Type, time.Now(), time.Now(), time.Now())

	if err != nil {
		return err
	}

	if err = r.AddUserToChat(chat.CreatorID, chat.ID, chat.WorkspaceID); err != nil {
		return err
	}

	if err = r.AddUserToChat(participantId, chat.ID, chat.WorkspaceID); err != nil {
		return err
	}

	systemMessage := domain.Message{
		ChatID:      chat.ID,
		SenderID:    chat.CreatorID,
		MessageType: types.MessageSystem,
		Content:     "chat created",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	err = r.messageRepo.AddMessage(&systemMessage)
	return err
}

func (r *chatRepo) AddUserToChat(userId uuid.UUID, chatId string, workspaceId uuid.UUID) error {
	_, err := r.conn.Exec(context.Background(),
		`INSERT INTO user_chats (user_id, chat_id, workspace_id) VALUES ($1, $2, $3)`, userId, chatId, workspaceId)
	return err
}

func (r *chatRepo) GetAllUserChats(userId uuid.UUID, workspaceId uuid.UUID) ([]response.GetChats, error) {
	rows, err := r.conn.Query(context.Background(), `
		SELECT
  			uc.chat_id,
  			uc.muted,
  			uc.pinned,
  			uc.notification,
  			uc.role,
  			c.type,
  			u2.is_online,
  			u2.id
		FROM user_chats uc
		JOIN chats c ON c.id = uc.chat_id
		JOIN user_chats uc2 ON uc2.chat_id = uc.chat_id AND uc2.user_id != $1
		JOIN users u2 ON u2.id = uc2.user_id
		WHERE uc.user_id = $1 AND c.workspace_id = $2`, userId, workspaceId)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var chats []response.GetChats
	for rows.Next() {
		var chat response.GetChats
		if err := rows.Scan(&chat.ID, &chat.Muted, &chat.Pinned, &chat.Notification, &chat.Role, &chat.Type, &chat.IsOnline, &chat.ParticipantID); err != nil {
			return nil, err
		}
		if err = r.getLastMessageInfo(&chat); err != nil {
			return nil, err
		}
		switch chat.Type {
		case types.TypePrivate:
			if err := r.getUserNameById(chat.ID, userId, &chat.Name); err != nil {
				chat.Name = "error"
			}
		}
		chats = append(chats, chat)
	}

	return chats, err
}

func (r *chatRepo) getUserNameById(chatId string, userId uuid.UUID, name *string) error {
	err := r.conn.QueryRow(context.Background(), `
		SELECT u.username 
		FROM user_chats uc
		JOIN users u ON uc.user_id = u.id
		WHERE uc.chat_id = $1 AND u.id != $2`, chatId, userId).Scan(name)
	return err
}

func (r *chatRepo) getLastMessageInfo(chat *response.GetChats) error {
	err := r.conn.QueryRow(context.Background(), `
		SELECT m.content, m.message_type, m.created_at 
		FROM messages m
		WHERE chat_id = $1
		ORDER BY m.created_at DESC
		LIMIT 1`, chat.ID).Scan(&chat.LastMessage, &chat.LastMessageType, &chat.LastMessageTime)
	return err
}

func (r *chatRepo) GetChatsBySearch(userID uuid.UUID, workspaceId uuid.UUID, value string) ([]response.GetChatsSearch, error) {
	rows, err := r.conn.Query(context.Background(), `
		 SELECT
		    u.id,
		    c.chat_id,
    		u.username,
    		c.last_message
		FROM user_workspaces uw
		JOIN users u ON u.id = uw.user_id
		LEFT JOIN LATERAL (
    		SELECT 
        		ch.id AS chat_id,
        		(
            		SELECT content
            		FROM messages
            		WHERE chat_id = ch.id
            		ORDER BY created_at DESC
            		LIMIT 1
        		) AS last_message
    		FROM chats ch
    		JOIN user_chats uc1 ON uc1.chat_id = ch.id AND uc1.user_id = $1
    		JOIN user_chats uc2 ON uc2.chat_id = ch.id AND uc2.user_id = u.id
      		AND ch.workspace_id = $2
		) c ON true
		WHERE uw.workspace_id = $2
  		AND u.id != $1
  		AND u.username ILIKE $3
	LIMIT 10;
	`, userID, workspaceId, "%"+value+"%")
	defer rows.Close()

	if err != nil {
		return nil, err
	}

	var chats []response.GetChatsSearch
	for rows.Next() {
		var chat response.GetChatsSearch
		if err := rows.Scan(&chat.UserID, &chat.ID, &chat.Name, &chat.LastMessage); err != nil {
			return nil, err
		}
		chats = append(chats, chat)
	}
	return chats, err
}
