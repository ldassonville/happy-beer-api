//go:build integration

package http

import (
	"context"
	"github.com/ldassonville/happy-beer-api/internal/beer"
	"github.com/ldassonville/happy-beer-api/pkg/api"
	apipuller "github.com/ldassonville/happy-beer-api/pkg/client"
	"github.com/sirupsen/logrus"
	"testing"
)

func getClient() apipuller.Client {

	config := &Config{
		ApiUrl: "http://localhost:9000",
	}

	return NewClient(config)
}

func getFirstDispenser(ctx context.Context, client apipuller.Client) (*api.Dispenser, error) {
	dispensers, err := client.SearchDispensers(ctx)
	if err != nil {
		logrus.WithError(err).Errorf("fail to obtain dispensers")
		return nil, err
	}

	if len(dispensers) > 0 {
		return dispensers[0], nil
	}
	return nil, nil
}

func TestClient_GetDispenser(t *testing.T) {

	ctx := context.Background()
	client := getClient()

	dispenser, _ := getFirstDispenser(ctx, client)
	if dispenser == nil {
		logrus.Info("skipping test. base is empty")
		return
	}

	dispensers, err := client.GetDispenser(ctx, dispenser.Ref)
	if err != nil {
		t.Fail()
		return
	}

	println(dispensers)
}

func TestClient_SearchDispensers(t *testing.T) {

	ctx := context.Background()
	client := getClient()

	dispensers, err := client.SearchDispensers(ctx)
	if err != nil {
		t.Fail()
		return
	}

	println(dispensers)
}

func TestClient_CreateDispenser(t *testing.T) {

	ctx := context.Background()
	client := getClient()

	dispenser := &api.DispenserEditable{
		Beer: beer.EasyBeer,
		Size: string(api.DispenserSizeL),
	}

	createdDispenser, err := client.CreateDispenser(ctx, dispenser)
	if err != nil {
		logrus.WithError(err).Error("fail to create dispenser")
		t.Fail()
	}

	println(createdDispenser)
}

func TestClient_DeleteDispenser(t *testing.T) {

	ctx := context.Background()
	client := getClient()

	dispenser, _ := getFirstDispenser(ctx, client)
	if dispenser == nil {
		logrus.Info("skipping test. base is empty")
		return
	}

	dispenserRef := dispenser.Ref
	err := client.DeleteDispenser(ctx, dispenserRef)
	if err != nil {
		logrus.WithError(err).Errorf("fail to deleting dispenser ref %s", dispenserRef)
		t.Fail()
		return
	}
}

func TestClient_UpdateDispenser(t *testing.T) {

	ctx := context.Background()
	client := getClient()

	dispenser, _ := getFirstDispenser(ctx, client)
	if dispenser == nil {
		logrus.Info("skipping test. base is empty")
		return
	}

	dispenser.Size = string(api.DispenserSizeS)

	updatedDispenser, err := client.UpdateDispenser(ctx, dispenser)
	if err != nil {
		logrus.WithError(err).Errorf("fail to update dispenser ref %s", dispenser.Ref)
		t.Fail()
		return
	}

	println(updatedDispenser)

}
