package usecase

import (
	"context"
	"go-postgres-test/internal/domain"
	"go-postgres-test/internal/entities"
	"go-postgres-test/internal/enums"
	"go-postgres-test/internal/helper"
	"go-postgres-test/internal/repository"
	"go-postgres-test/internal/response"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type ChatUsecase struct {
	chatRepo    domain.ChatRepository
	messageRepo domain.MessageRepository
	baseRepo    *repository.BaseRepo
}

func NewChatUsecase(
	chatRepo domain.ChatRepository,
	messageRepo domain.MessageRepository,
	baseRepo *repository.BaseRepo,
) *ChatUsecase {
	return &ChatUsecase{chatRepo: chatRepo, baseRepo: baseRepo, messageRepo: messageRepo}
}

func (uc *ChatUsecase) AddUserToChat(ctx context.Context, userId uuid.UUID, chatId string, workspaceId uuid.UUID, role enums.ChatRole) error {
	return uc.chatRepo.AddUserToChat(ctx, userId, chatId, workspaceId, role)
}

func (uc *ChatUsecase) CreateChat(ctx context.Context, chat *domain.Chat, participantId uuid.UUID) error {
	return helper.WithTxVoid(ctx, uc.baseRepo, func(ctx context.Context, exec pgx.Tx) error {
		if chat.CreatorID.String() < participantId.String() {
			chat.ID = "c-" + chat.CreatorID.String() + "&" + participantId.String()
		} else {
			chat.ID = "c-" + participantId.String() + "&" + chat.CreatorID.String()
		}

		if err := uc.chatRepo.AddUserToChatTx(ctx, exec, chat.CreatorID, chat.ID, chat.WorkspaceID, enums.ChatUser); err != nil {
			return err
		}

		if err := uc.chatRepo.AddUserToChatTx(ctx, exec, participantId, chat.ID, chat.WorkspaceID, enums.ChatUser); err != nil {
			return err
		}

		systemMessage := entities.Message{
			ChatID:      chat.ID,
			SenderID:    chat.CreatorID,
			MessageType: enums.MessageSystem,
			Content:     "chat created",
			CreatedAt:   time.Now().UTC(),
			UpdatedAt:   time.Now().UTC(),
		}

		return uc.messageRepo.AddMessageTx(ctx, exec, systemMessage)
	})
}

func (uc *ChatUsecase) GetAllUserChats(ctx context.Context, userId uuid.UUID, workspaceId uuid.UUID) ([]response.GetChats, error) {
	return helper.WithTx(ctx, uc.baseRepo, func(ctx context.Context, exec pgx.Tx) ([]response.GetChats, error) {
		return uc.chatRepo.GetAllUserChatsTx(ctx, exec, userId, workspaceId)
	})
}

func (uc *ChatUsecase) GetChatsBySearch(ctx context.Context, userID uuid.UUID, workspaceId uuid.UUID, value string) ([]response.GetChatsSearch, error) {
	return uc.chatRepo.GetChatsBySearch(ctx, userID, workspaceId, value)
}
