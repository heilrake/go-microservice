package events

import (
	"context"
	"encoding/json"

	"ride-sharing/services/user-service/internal/domain"
	"ride-sharing/shared/contracts"
	"ride-sharing/shared/messaging"
	userPublishers "ride-sharing/shared/messaging/publishers"
	pb "ride-sharing/shared/proto/user"
)

type UserEventPublisher struct {
	rabbitmq *messaging.RabbitMQ
}

func NewUserEventPublisher(rabbitmq *messaging.RabbitMQ) *UserEventPublisher {
	return &UserEventPublisher{
		rabbitmq: rabbitmq,
	}
}

func (p *UserEventPublisher) PublishUserCreated(ctx context.Context, user *domain.User) error {
	return p.publish(ctx, contracts.UserEventCreated, user)
}

func (p *UserEventPublisher) PublishUserUpdated(ctx context.Context, user *domain.User) error {
	return p.publish(ctx, contracts.UserEventUpdated, user)
}

func (p *UserEventPublisher) PublishUserDeleted(ctx context.Context, userID string) error {
	payload := messaging.UserEventData{
		User: &pb.User{
			Id: userID,
		},
	}

	userEventJSON, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	return userPublishers.NewUserPublisher(p.rabbitmq).Publish(ctx, contracts.UserEventDeleted, contracts.AmqpMessage{
		OwnerID: userID,
		Data:    userEventJSON,
	})
}

func (p *UserEventPublisher) publish(ctx context.Context, routingKey string, user *domain.User) error {
	payload := messaging.UserEventData{
		User: &pb.User{
			Id:             user.ID,
			Username:       user.Username,
			Email:          user.Email,
			ProfilePicture: user.ProfilePicture,
		},
	}

	userEventJSON, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	return userPublishers.NewUserPublisher(p.rabbitmq).Publish(ctx, routingKey, contracts.AmqpMessage{
		OwnerID: user.ID,
		Data:    userEventJSON,
	})
}
