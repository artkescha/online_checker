package transmitter

import (
	"github.com/artkescha/grader_api/send_solution"
	"github.com/nats-io/stan.go"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

//go:generate mockgen -destination=./transmitter_mock.go -package=transmitter . Transmitter

type Transmitter interface {
	Transmit(topic string, try *send_solution.Try) error
}

type Publisher struct {
	natsConn stan.Conn
	logger   *zap.SugaredLogger
}

func New(natsConn stan.Conn, logger *zap.SugaredLogger) *Publisher {
	return &Publisher{natsConn: natsConn, logger: logger}
}

func (s *Publisher) Transmit(topic string, try *send_solution.Try) error {
	data, err := proto.Marshal(try)

	if err != nil {
		return err
	}

	_, err = s.natsConn.PublishAsync(topic, data, func(guid string, err error) {
		if err != nil {
			s.logger.Errorf("send publish msg with guid %s failed %s", guid, err)
			return
		}
		s.logger.Debugf("send publish msg with guid %s sucsess", guid)
	})
	if err != nil {
		return err
	}
	return nil
}
