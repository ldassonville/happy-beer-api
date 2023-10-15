package dispensers

import (
	"context"
	"errors"
	"fmt"
	"github.com/ldassonville/beer-puller-api/internal/beer"
	"github.com/ldassonville/beer-puller-api/internal/dispensers/storage"
	"github.com/ldassonville/beer-puller-api/internal/records"
	"github.com/ldassonville/beer-puller-api/pkg/api"
	"github.com/ldassonville/beer-puller-api/pkg/core/event"
	"github.com/sirupsen/logrus"
	"strings"
	"time"
)

type Service struct {
	dao storage.Storage

	dispatcher event.Dispatcher
	recordsSvc *records.Service
}

func NewService(dao storage.Storage, recordsSvc *records.Service) *Service {
	return &Service{
		dao:        dao,
		recordsSvc: recordsSvc,
		dispatcher: event.Dispatcher{},
	}
}

func (s *Service) Search(ctx context.Context) ([]*api.Dispenser, error) {

	dispensers, err := s.dao.Search(ctx)

	return dispensers, err

}

func (s *Service) GetByRef(ctx context.Context, ref string) (*api.Dispenser, error) {

	return s.dao.GetByRef(ctx, ref)

}

func (s *Service) addRecord(ctx context.Context, msg string, dispenser *api.Dispenser) {
	// Record event
	s.recordsSvc.Create(ctx, &api.Record{
		Message: msg,
		Target:  dispenser,
	})
}

func (s *Service) Update(ctx context.Context, dispenser *api.Dispenser) (*api.Dispenser, error) {

	dispenser, err := s.dao.Update(ctx, dispenser)
	if err != nil {
		logrus.WithError(err).Warnf("fail to update dispenser")
		return nil, err
	}
	s.dispatcher.Dispatch(ctx, "UPDATE", dispenser)
	return dispenser, err
}

func (s *Service) Create(ctx context.Context, dispenserEditable *api.DispenserEditable) (*api.Dispenser, error) {

	dispenser := &api.Dispenser{
		DispenserEditable: *dispenserEditable,
		Metadata: &api.Metadata{
			CreatedAt: time.Now(),
		},
	}
	s.addRecord(context.Background(), fmt.Sprintf("dispenser %s request have been receive", dispenserEditable.Beer), dispenser)

	shouldSimulateError := s.preCreateHook(dispenser)

	dispenser, err := s.dao.Create(ctx, dispenser)
	if err != nil {
		logrus.WithError(err).Warnf("fail to create dispenser")
		return nil, err
	}

	// Simulate a creation error
	if shouldSimulateError {
		// Add business record
		s.addRecord(ctx, "simulate dispenser creation error", dispenser)
		s.dispatcher.Dispatch(ctx, "ERROR", dispenser)

		return nil, errors.New("creation_error")
	}

	// Add business record
	s.dispatcher.Dispatch(ctx, "CREATE", dispenser)
	s.postCreateHook(dispenser)

	return dispenser, err
}

func (s *Service) preCreateHook(dispenser *api.Dispenser) bool {

	dispenser.State = api.DispenserReady
	dispenser.Status = &api.DispenserStatus{
		Status: api.InternalStatusOk,
	}

	switch dispenser.Beer {
	case beer.LatencyBeer:
		dispenser.State = api.DispenserNone
		dispenser.Status.Status = api.InternalStatusProcessing
		dispenser.Status.Hidden = true
		dispenser.Status.Info = "Simulate sync request with latency. Response not delivered yet"
	case beer.StolenBeer:
	case beer.EasyBeer:
	case beer.LazyBeer:
		s.addRecord(context.Background(), "the dispenser is refreshing...", dispenser)
		dispenser.State = api.DispenserRefreshing
		dispenser.Status.Status = api.InternalStatusOk
		dispenser.Status.Info = "Simulate async ready process"
	case beer.FatalBeer:
		// handle error
		s.addRecord(context.Background(), "the dispenser request is in failure", dispenser)
		dispenser.State = api.DispenserNone
		dispenser.Status.Status = api.InternalStatusFatal
		dispenser.Status.Info = "Simulate request in failure"
		dispenser.Status.Hidden = true
		return true
	}
	return false
}

func (s *Service) postCreateHook(dispenser *api.Dispenser) error {

	switch dispenser.Beer {

	case beer.FatalBeer:
	case beer.EasyBeer:
	case beer.LazyBeer:
		go func(ref string) {
			time.Sleep(30 * time.Second)
			d, err := s.dao.GetByRef(context.Background(), ref)
			if err != nil {
				logrus.WithError(err).Warnf("fail to change status")
				return
			}

			previousStatus := strings.ToLower(string(d.State))
			d.State = api.DispenserReady
			d.Status.Status = api.InternalStatusOk
			d.Status.Info = "Async process simulation terminated"

			_, err = s.Update(context.Background(), d)
			if err != nil {
				logrus.WithError(err).Warnf("fail to update status at callback")
				return
			}
			s.addRecord(context.Background(), fmt.Sprintf("the dispenser change state from %s to ready", previousStatus), dispenser)
		}(dispenser.Ref)

	case beer.LatencyBeer:

		s.addRecord(context.Background(), fmt.Sprintf("the dispenser request have latency"), dispenser)
		// Take time to create the beer
		time.Sleep(55 * time.Second)

		dispenser.State = api.DispenserReady
		dispenser.Status.Status = api.InternalStatusOk
		dispenser.Status.Info = "Simulation of request with latency terminated"
		dispenser.Status.Hidden = false
		_, err := s.Update(context.Background(), dispenser)
		if err != nil {
			logrus.WithError(err).Warnf("fail to change status")
		}

	case beer.StolenBeer:

		// Stolen beer will disappear un 3min
		go func(ref string) {
			time.Sleep(90 * time.Second)
			s.logicalDeleteByRef(context.Background(), ref, api.InternalStatusStolen, "The dispenser been has been stolen")
			s.addRecord(context.Background(), "the dispenser have been stolen", dispenser)
		}(dispenser.Ref)
	}

	return nil
}

/*
func (s *Service) DeleteByRef(ctx context.Context, ref string) error {

	dispenser, err := s.dao.GetByRef(ctx, ref)
	if err != nil {
		logrus.WithError(err).Warnf("fail to obtain dispenser before remove")
		return err
	}

	if dispenser == nil {
		return errors.New("not_found")
	}

	err = s.dao.DeleteByRef(ctx, ref)
	if err != nil {
		logrus.WithError(err).Warnf("fail to delete dispenser")
		return err
	}
	s.dispatcher.Dispatch(ctx, "DELETE", dispenser)
	s.addRecord(context.Background(), fmt.Sprintf("the dispenser have been remove"), dispenser)

	return err
}*/

func (s *Service) DeleteByRef(ctx context.Context, ref string) error {

	return s.logicalDeleteByRef(ctx, ref, api.InternalStatusReturned, "The dispenser has been returned (api ask for delete) ")
}

func (s *Service) logicalDeleteByRef(ctx context.Context, ref string, status api.InternalStatus, reason string) error {

	dispenser, err := s.dao.GetByRef(ctx, ref)
	if err != nil {
		logrus.WithError(err).Warnf("fail to obtain dispenser before remove")
		return err
	}

	if dispenser == nil {
		return errors.New("not_found")
	}

	dispenser.State = api.DispenserNone
	dispenser.Status.Status = status
	dispenser.Status.Info = reason
	dispenser.Status.Hidden = true

	dispenser, err = s.dao.Update(ctx, dispenser)
	if err != nil {
		logrus.WithError(err).Warnf("fail to mark dispenser as logical deleted")
		return err
	}
	s.dispatcher.Dispatch(ctx, "DELETE", dispenser)
	s.addRecord(context.Background(), fmt.Sprintf("the dispenser have been remove. Reason : %s", reason), dispenser)

	return err
}

func (s *Service) Subscribe() *event.Subscription {
	return s.dispatcher.Subscribe()
}
