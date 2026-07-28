package usecase

import (
	"backend-task-manager/internal/domain"
	"backend-task-manager/internal/dto"
	"backend-task-manager/internal/entities"
	"backend-task-manager/internal/enums"
	"backend-task-manager/internal/helper"
	"backend-task-manager/internal/repository"
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type MessageUsecase struct {
	messageRepo domain.MessageRepository
	chatRepo    domain.ChatRepository
	baseRepo    *repository.BaseRepo
}

func NewMessageUsecase(repo domain.MessageRepository, chatRepo domain.ChatRepository, baseRepo *repository.BaseRepo) *MessageUsecase {
	return &MessageUsecase{messageRepo: repo, chatRepo: chatRepo, baseRepo: baseRepo}
}

func (uc *MessageUsecase) AddMessage(ctx context.Context, message entities.Message) ([]uuid.UUID, error) {
	return helper.WithTx(ctx, uc.baseRepo, func(ctx context.Context, exec pgx.Tx) ([]uuid.UUID, error) {
		return uc.AddMessageTx(ctx, exec, message)
	})
}

func (uc *MessageUsecase) AddMessageTx(ctx context.Context, exec entities.Execer, message entities.Message) ([]uuid.UUID, error) {
	message.ID = "m-" + uuid.New().String()

	if err := uc.messageRepo.AddMessageTx(ctx, exec, message); err != nil {
		return []uuid.UUID{}, err
	}

	if err := uc.chatRepo.UpdateChatTimeTx(ctx, exec, message.ChatID, message.UpdatedAt); err != nil {
		return []uuid.UUID{}, err
	}

	return uc.messageRepo.AddMessageReadsTx(ctx, exec, message.ChatID, message.ID, message.SenderID)
}

func (uc *MessageUsecase) GetAllMessages(ctx context.Context, chatId string) ([]entities.Message, error) {
	return uc.messageRepo.GetAllMessages(ctx, chatId)
}

func (uc *MessageUsecase) UploadImage(ctx context.Context, uploadImageDto dto.UploadImage) (entities.Message, error) {
	return helper.WithTx(ctx, uc.baseRepo, func(ctx context.Context, exec pgx.Tx) (entities.Message, error) {
		message := entities.Message{
			ID:          "",
			ChatID:      uploadImageDto.ChatId,
			SenderID:    uploadImageDto.UserId,
			Content:     uploadImageDto.Content,
			MessageType: enums.MessageImage,
			CreatedAt:   time.Now().UTC(),
			UpdatedAt:   time.Now().UTC(),
		}

		unreadUserIds, err := uc.AddMessageTx(ctx, exec, message)
		if err != nil {
			return entities.Message{}, err
		}

		message.UnreadUsersIds = unreadUserIds

		files := uploadImageDto.Form.File["images"]

		var messageAttachments []entities.Attachment
		for _, fileHeader := range files {
			file, err := fileHeader.Open()
			if err != nil {
				continue
			}

			url, err := uc.messageRepo.UploadImage(file)
			closeErr := file.Close()
			if err != nil {
				return entities.Message{}, err
			}
			if closeErr != nil {
				return entities.Message{}, closeErr
			}

			var newAttachment = entities.Attachment{
				ID:         uuid.New(),
				MessageID:  message.ID,
				FileUrl:    url,
				FileType:   enums.MessageImage,
				FileName:   fileHeader.Filename,
				FileSize:   fileHeader.Size,
				UploadedAt: time.Now().UTC(),
				ChatID:     message.ChatID,
			}

			if err = uc.messageRepo.AddAttachmentTx(ctx, exec, newAttachment); err != nil {
				return entities.Message{}, err
			}

			messageAttachments = append(messageAttachments, newAttachment)
		}

		message.Attachments = messageAttachments

		return message, nil
	})
}

func (uc *MessageUsecase) AddAttachment(ctx context.Context, attachment entities.Attachment) error {
	return uc.messageRepo.AddAttachment(ctx, attachment)
}

func (uc *MessageUsecase) SetMessagesRead(ctx context.Context, messages []uuid.UUID, userId uuid.UUID) error {
	return helper.WithTxVoid(ctx, uc.baseRepo, func(ctx context.Context, exec pgx.Tx) error {
		for _, m := range messages {
			err := uc.messageRepo.SetMessageReadTx(ctx, exec, m, userId)
			if err != nil {
				return err
			}
		}

		return nil
	})
}

func (uc *MessageUsecase) GetAllChatImages(ctx context.Context, chatId string) (*[]entities.Attachment, error) {
	return uc.messageRepo.GetAllChatImages(ctx, chatId)
}
