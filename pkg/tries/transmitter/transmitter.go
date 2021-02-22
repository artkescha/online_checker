package transmitter

import (
	"encoding/json"
	"github.com/artkescha/grader/online_checker/pkg/tries"
	"github.com/nats-io/nats.go"
)

type Transmitter interface {
	Transmit(topic string, try try.Try) error
}

type Publisher struct {
	natsConn *nats.Conn
}

func New(natsConn *nats.Conn) *Publisher {
	return &Publisher{natsConn: natsConn}
}

func (s *Publisher) Transmit(topic string, try try.Try) error {
	data, err := json.Marshal(try)

	if err != nil {
		return err
	}

	err = s.natsConn.Publish(topic, data)

	if err != nil {
		return err
	}

	return nil
}
