package events

import (
	"context"
	"encoding/json"
	"ride-sharing/services/trip-service/internal/domain"
	"ride-sharing/shared/contracts"
	"ride-sharing/shared/messaging"
	tripPublishers "ride-sharing/shared/messaging/publishers"
)

type TripEventPublisher struct {
	rabbitmq *messaging.RabbitMQ
}

func NewTripEventPublisher(rabbitmq *messaging.RabbitMQ) *TripEventPublisher {
	return &TripEventPublisher{
		rabbitmq: rabbitmq,
	}
}

func (p *TripEventPublisher) PublishTripCreated(ctx context.Context, trip *domain.TripModel) error {
	payload := messaging.TripEventData{
		Trip: trip.ToProto(),
	}

	tripEventJSON, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	return tripPublishers.NewTripPublisher(p.rabbitmq).Publish(ctx, contracts.TripEventCreated, contracts.AmqpMessage{
		OwnerID: trip.UserID,
		Data:    tripEventJSON,
	})
}

func (p *TripEventPublisher) PublishTripCompleted(ctx context.Context, trip *domain.TripModel) error {
	payload := messaging.TripCompletedData{
		TripID: trip.ID,
		UserID: trip.UserID,
	}

	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	return p.rabbitmq.PublishMessage(ctx, contracts.PaymentCmdCapturePayment, contracts.AmqpMessage{
		OwnerID: trip.UserID,
		Data:    data,
	})
}
