package repository

import (
	"context"
	"go-postgres-test/internal/domain"
	"go-postgres-test/internal/dto/response"
	"go-postgres-test/internal/entities"
	"go-postgres-test/internal/enums"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type chatRepo struct {
	conn *pgxpool.Pool
}

func NewChatRepo(conn *pgxpool.Pool) domain.ChatRepository {
	return &chatRepo{conn: conn}
}

func (r *chatRepo) AddUserToChat(ctx context.Context, userId uuid.UUID, chatId string, workspaceId uuid.UUID, role enums.UserRoles) error {
	return r.AddUserToChatTx(ctx, r.conn, userId, chatId, workspaceId, role)
}

func (r *chatRepo) AddUserToChatTx(ctx context.Context, exec entities.Execer, userId uuid.UUID, chatId string, workspaceId uuid.UUID, role enums.UserRoles) error {
	_, err := exec.Exec(ctx,
		`INSERT INTO user_chats (user_id, chat_id, workspace_id, role) VALUES ($1, $2, $3, $4)`, userId, chatId, workspaceId, role)
	return err
}

func (r *chatRepo) CreateChat(ctx context.Context, chat *domain.Chat) (err error) {
	return r.CreateChatTx(ctx, r.conn, chat)
}

func (r *chatRepo) CreateChatTx(ctx context.Context, exec entities.Execer, chat *domain.Chat) error {
	now := time.Now().UTC()
	_, err := exec.Exec(ctx, `
		INSERT INTO chats (id, workspace_id, creator_id, type, created_at, updated_at, last_message_time)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`, chat.ID, chat.WorkspaceID, chat.CreatorID, chat.Type, now, now, now)
	return err
}

func (r *chatRepo) GetAllUserChatsTx(ctx context.Context, exec entities.Execer, userId uuid.UUID, workspaceId uuid.UUID) ([]response.GetChats, error) {
	rows, err := exec.Query(ctx, `
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
		LEFT JOIN user_chats uc2 ON uc2.chat_id = uc.chat_id AND uc2.user_id != $1
		LEFT JOIN users u2 ON u2.id = uc2.user_id
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
		chats = append(chats, chat)
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return chats, err
}

func (r *chatRepo) GetChatsBySearch(ctx context.Context, userID uuid.UUID, workspaceId uuid.UUID, value string) ([]response.GetChatsSearch, error) {
	rows, err := r.conn.Query(ctx, `
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

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var chats []response.GetChatsSearch
	for rows.Next() {
		var chat response.GetChatsSearch
		if err := rows.Scan(&chat.UserID, &chat.ID, &chat.Name, &chat.LastMessage); err != nil {
			return nil, err
		}
		chats = append(chats, chat)
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}
	return chats, err
}

func (r *chatRepo) GetParticipantNameByChatIdTx(ctx context.Context, exec entities.Execer, chatId string, userId uuid.UUID, name *string) error {
	err := exec.QueryRow(ctx, `
		SELECT u.username 
		FROM user_chats uc
		JOIN users u ON uc.user_id = u.id
		WHERE uc.chat_id = $1 AND u.id != $2`, chatId, userId).Scan(name)
	return err
}

func (r *chatRepo) GetLastMessageInfo(ctx context.Context, chat *response.GetChats, userID uuid.UUID) error {
	return r.GetLastMessageInfoTx(ctx, r.conn, chat, userID)
}

func (r *chatRepo) GetLastMessageInfoTx(ctx context.Context, exec entities.Execer, chat *response.GetChats, userID uuid.UUID) error {
	err := exec.QueryRow(ctx, `
		SELECT
		    m.content,
		    m.message_type, 
		    m.created_at,
		    COALESCE(ma.file_url, '') as lastMessageAttachment,
		    (
				SELECT COUNT(*) 
				FROM message_reads
				WHERE chat_id = $1 AND user_id = $3 AND read_at IS NULL
			) AS unread_count
		FROM messages m
		LEFT JOIN message_attachments ma
			ON ma.message_id = m.id AND m.message_type = $2
		WHERE m.chat_id = $1
		ORDER BY m.created_at DESC
		LIMIT 1
	`, chat.ID, enums.MessageImage, userID).Scan(
		&chat.LastMessage,
		&chat.LastMessageType,
		&chat.LastMessageTime,
		&chat.LastMessageAttachment,
		&chat.UnreadForUser,
	)

	return err
}

func (r *chatRepo) UpdateChatTime(ctx context.Context, chatId string, at time.Time) error {
	return r.UpdateChatTimeTx(ctx, r.conn, chatId, at)
}

func (r *chatRepo) UpdateChatTimeTx(ctx context.Context, exec entities.Execer, chatId string, at time.Time) error {
	_, err := exec.Exec(ctx,
		`UPDATE chats SET updated_at = $1 WHERE id = $2
	`, at, chatId)
	return err
}
