package events

import (
	"context"
	"errors"
	"strings"
	"sync"

	uuid "github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
)

type Handler func(interface{}) error

type Manager interface {
	Subscribe(topic string, cb Handler) (id string, err error)
	Unsubscribe(topic, id string) error
	Publish(topic string, payload interface{}) error
}

func NewManager(ctx context.Context, workers int) Manager {
	m := &manager{
		commands:  make(chan command, 64),
		callbacks: make(map[string]map[string]Handler),
		workers:   workers,
	}
	go m.backend(ctx)
	return m
}

type manager struct {
	commands  chan command
	callbacks map[string]map[string]Handler
	workers   int
}

func (m *manager) Subscribe(topic string, cb Handler) (id string, err error) {
	logrus.Debugf("events: subscribe to %v", topic)
	idChan := make(chan string)
	errChan := make(chan error)
	m.commands <- command{
		commandType: subscribe,
		subscribeCommand: &subscribeCommand{
			topic:  topic,
			cb:     cb,
			idChan: idChan,
		},
		error: errChan,
	}
	return <-idChan, <-errChan
}

func (m *manager) Unsubscribe(topic, id string) error {
	logrus.Debugf("events: unsubscribe from %v (id: %v)", topic, id)
	errChan := make(chan error)
	m.commands <- command{
		commandType: unsubscribe,
		unsubscribeCommand: &unsubscribeCommand{
			topic: topic,
			id:    id,
		},
		error: errChan,
	}
	return <-errChan
}

func (m *manager) Publish(topic string, payload interface{}) error {
	logrus.Debugf("events: publish event to %v", topic)
	errChan := make(chan error)
	m.commands <- command{
		commandType: publish,
		publishCommand: &publishCommand{
			topic:   topic,
			payload: payload,
		},
		error: errChan,
	}
	errs := make([]error, 0)
	for err := range errChan {
		errs = append(errs, err)
	}
	if len(errs) > 0 {
		return MultiError(errs)
	}
	return nil
}

type commandType string

const (
	publish     commandType = "publish"
	subscribe   commandType = "subscribe"
	unsubscribe commandType = "unsubscribe"
)

type publishCommand struct {
	topic   string
	payload interface{}
}

type subscribeCommand struct {
	topic  string
	cb     Handler
	idChan chan<- string
}

type unsubscribeCommand struct {
	topic string
	id    string
}

type command struct {
	commandType
	*publishCommand
	*subscribeCommand
	*unsubscribeCommand
	error chan error
}

func (m *manager) backend(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case cmd, ok := <-m.commands:
			if !ok {
				return nil
			}
			switch cmd.commandType {
			case subscribe:
				m.subscribe(cmd)
			case unsubscribe:
				m.unsubscribe(cmd)
			case publish:
				m.publish(cmd)
			}
		}
	}
}

func (m *manager) subscribe(cmd command) {
	defer func() {
		close(cmd.subscribeCommand.idChan)
		close(cmd.error)
	}()
	topicMap, ok := m.callbacks[cmd.subscribeCommand.topic]
	if !ok {
		topicMap = make(map[string]Handler)
	}
	id := uuid.NewV4().String()
	topicMap[id] = cmd.cb
	m.callbacks[cmd.subscribeCommand.topic] = topicMap
	cmd.idChan <- id
}

func (m *manager) unsubscribe(cmd command) {
	defer func() {
		close(cmd.error)
	}()
	topicMap, ok := m.callbacks[cmd.unsubscribeCommand.topic]
	if !ok {
		cmd.error <- errors.New("not found")
		return
	}
	delete(topicMap, cmd.unsubscribeCommand.id)
	m.callbacks[cmd.unsubscribeCommand.topic] = topicMap
}

type workerInput struct {
	cb  Handler
	arg interface{}
}

func (m *manager) publish(cmd command) {
	defer func() {
		close(cmd.error)
	}()
	callbacks, ok := m.callbacks[cmd.publishCommand.topic]
	if !ok {
		return
	}
	var wg sync.WaitGroup
	inputs := make(chan workerInput, m.workers)
	for i := 0; i < m.workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				input, ok := <-inputs
				if !ok {
					return
				}
				err := input.cb(input.arg)
				if err != nil {
					cmd.error <- err
				}
			}
		}()
	}
	for id, cb := range callbacks {
		inputs <- workerInput{cb, cmd.publishCommand.payload}
		logrus.Debugf("events: call handler with id %v", id)
	}
	close(inputs)
	wg.Wait()
}

type MultiError []error

func (e MultiError) Error() string {
	var res strings.Builder
	for _, err := range e {
		res.WriteString(err.Error())
		res.WriteRune('\n')
	}
	return res.String()
}
