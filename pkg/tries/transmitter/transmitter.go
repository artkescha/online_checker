package transmitter

import (
	"github.com/artkescha/grader_api/send_solution"
    "google.golang.org/protobuf/proto"

	"github.com/nats-io/nats.go"
)

//go:generate mockgen -destination=./transmitter_mock.go -package=transmitter . Transmitter

type Transmitter interface {
	Transmit(topic string, try *send_solution.Try) error
}

type Publisher struct {
	natsConn *nats.Conn
}

func New(natsConn *nats.Conn) *Publisher {
	return &Publisher{natsConn: natsConn}
}

func (s *Publisher) Transmit(topic string, try *send_solution.Try) error {
	data, err := proto.Marshal(try)

	if err != nil {
		return err
	}

	err = s.natsConn.Publish(topic, data)
	if err != nil {
		return err
	}

	return nil
}
