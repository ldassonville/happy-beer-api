package event

import (
	"context"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"sync"
)

type Dispatcher struct {
	subscribers sync.Map
}

func (s *Dispatcher) Dispatch(ctx context.Context, typ string, content any) {

	data, err := json.Marshal(content)
	if err != nil {
		logrus.WithError(err).Warnf("fail to send %s msg : %v", typ, content)
		return
	}
	s.DispatchEvent(ctx, Event{
		Typ: typ,
		Msg: string(data),
	})
}

func (s *Dispatcher) DispatchEvent(ctx context.Context, event Event) {

	s.subscribers.Range(func(id any, sub any) bool {
		go func(event Event) {
			sub.(*Subscription).c <- event
		}(event)
		return true
	})
}

func (s *Subscription) Id() string {
	return s.id
}

type Subscription struct {
	id       string
	c        chan Event
	leaveFnc func()
}

type Event struct {
	Typ string
	Msg string
}

func (s *Subscription) Chan() chan Event {
	return s.c
}

func (s *Subscription) Leave() {
	s.leaveFnc()
}

func (s *Dispatcher) Subscribe() *Subscription {

	sub := &Subscription{
		id: uuid.NewString(),
		c:  make(chan Event, 10),
	}
	s.subscribers.Store(sub.id, sub)

	sub.leaveFnc = func() {
		s.subscribers.Delete(sub.id)
		close(sub.c)

	}
	return sub
}
