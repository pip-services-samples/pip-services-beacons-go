package test_persistence

import (
	"testing"

	bpersist "github.com/pip-services-samples/pip-data-microservice-go/persistence"
	cconf "github.com/pip-services3-go/pip-services3-commons-go/config"
)

func TestBeaconsFilePersistence(t *testing.T) {
	var persistence *bpersist.BeaconsFilePersistence
	var fixture *BeaconsPersistenceFixture

	persistence = bpersist.NewBeaconsFilePersistence("../../temp/beacons.test.json")
	persistence.Configure(cconf.NewEmptyConfigParams())
	fixture = NewBeaconsPersistenceFixture(persistence)

	opnErr := persistence.Open("")
	if opnErr == nil {
		persistence.Clear("")
	}

	defer persistence.Close("")

	t.Run("BeaconsFilePersistence:CRUD Operations", fixture.TestCrudOperations)
	persistence.Clear("")
	t.Run("BeaconsFilePersistence:Get with Filters", fixture.TestGetWithFilters)
}