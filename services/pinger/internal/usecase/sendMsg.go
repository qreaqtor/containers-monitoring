package usecase

import (
	"encoding/json"
	"log/slog"

	"github.com/IBM/sarama"
	"github.com/qreaqtor/containers-monitoring/pinger/internal/models"
)

func (p *Pinger) sendMsg(msg models.ContainersMsg) {
	msgJSON, err := json.Marshal(msg)
	if err != nil {
		slog.Error(err.Error())
		return
	}
	message := &sarama.ProducerMessage{
		Topic: p.topic,
		Value: sarama.ByteEncoder(msgJSON),
	}

	p.producer.Input() <- message
}
