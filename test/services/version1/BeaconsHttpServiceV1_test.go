package test_services

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"testing"

	bdata "github.com/pip-services-samples/pip-data-microservice-go/data/version1"
	blogic "github.com/pip-services-samples/pip-data-microservice-go/logic"
	bpersist "github.com/pip-services-samples/pip-data-microservice-go/persistence"
	bservices "github.com/pip-services-samples/pip-data-microservice-go/services/version1"
	cconf "github.com/pip-services3-go/pip-services3-commons-go/config"
	cref "github.com/pip-services3-go/pip-services3-commons-go/refer"
	"github.com/stretchr/testify/assert"
)

var Beacon1 bdata.BeaconV1 = bdata.BeaconV1{
	Id:      "1",
	Udi:     "00001",
	Type:    bdata.BeaconTypeV1.AltBeacon,
	Site_id: "1",
	Label:   "TestBeacon1",
	Center:  bdata.GeoPointV1{Type: "Point", Coordinates: [][]float32{{0.0, 0.0}}},
	Radius:  50,
}

var Beacon2 bdata.BeaconV1 = bdata.BeaconV1{
	Id:      "2",
	Udi:     "00002",
	Type:    bdata.BeaconTypeV1.IBeacon,
	Site_id: "1",
	Label:   "TestBeacon2",
	Center:  bdata.GeoPointV1{Type: "Point", Coordinates: [][]float32{{2.0, 2.0}}},
	Radius:  70,
}

func TestBeaconsHttpServiceV1(t *testing.T) {

	var persistence *bpersist.BeaconsMemoryPersistence
	var controller *blogic.BeaconsController
	var service *bservices.BeaconsHttpServiceV1
	var url string = "http://localhost:3000"

	persistence = bpersist.NewBeaconsMemoryPersistence()
	persistence.Configure(cconf.NewEmptyConfigParams())

	controller = blogic.NewBeaconsController()
	controller.Configure(cconf.NewEmptyConfigParams())
	service = bservices.NewBeaconsHttpServiceV1()
	service.Configure(cconf.NewConfigParamsFromTuples(
		"connection.protocol", "http",
		"connection.port", "3000",
		"connection.host", "localhost",
	))

	references := cref.NewReferencesFromTuples(
		cref.NewDescriptor("beacons", "persistence", "memory", "default", "1.0"), persistence,
		cref.NewDescriptor("beacons", "controller", "default", "default", "1.0"), controller,
		cref.NewDescriptor("beacons", "service", "http", "default", "1.0"), service,
	)

	controller.SetReferences(references)
	service.SetReferences(references)

	opnErr := persistence.Open("")
	if opnErr != nil {
		panic("Can't open persistence")
	}
	service.Open("")
	defer service.Close("")
	defer persistence.Close("")

	var beacon1 bdata.BeaconV1
	// Create the first beacon
	bodyMap := make(map[string]interface{})
	bodyMap["beacon"] = Beacon1
	jsonBody, _ := json.Marshal(bodyMap)
	bodyReader := bytes.NewReader(jsonBody)
	postResponse, postErr := http.Post(url+"/v1/beacons/create_beacon", "application/json", bodyReader)
	assert.Nil(t, postErr)
	resBody, bodyErr := ioutil.ReadAll(postResponse.Body)
	assert.Nil(t, bodyErr)
	var beacon bdata.BeaconV1
	jsonErr := json.Unmarshal(resBody, &beacon)
	assert.Nil(t, jsonErr)
	assert.NotNil(t, beacon)
	assert.Equal(t, Beacon1.Udi, beacon.Udi)
	assert.Equal(t, Beacon1.Site_id, beacon.Site_id)
	assert.Equal(t, Beacon1.Type, beacon.Type)
	assert.Equal(t, Beacon1.Label, beacon.Label)
	assert.NotNil(t, beacon.Center)

	// Create the second beacon
	bodyMap = make(map[string]interface{})
	bodyMap["beacon"] = Beacon2
	jsonBody, _ = json.Marshal(bodyMap)
	bodyReader = bytes.NewReader(jsonBody)
	postResponse, postErr = http.Post(url+"/v1/beacons/create_beacon", "application/json", bodyReader)
	assert.Nil(t, postErr)
	resBody, bodyErr = ioutil.ReadAll(postResponse.Body)
	assert.Nil(t, bodyErr)

	jsonErr = json.Unmarshal(resBody, &beacon)
	assert.Nil(t, jsonErr)
	assert.NotNil(t, beacon)
	assert.Equal(t, Beacon2.Udi, beacon.Udi)
	assert.Equal(t, Beacon2.Site_id, beacon.Site_id)
	assert.Equal(t, Beacon2.Type, beacon.Type)
	assert.Equal(t, Beacon2.Label, beacon.Label)
	assert.NotNil(t, beacon.Center)

	// Get all beacons
	postResponse, postErr = http.Post(url+"/v1/beacons/get_beacons", "application/json", nil)
	assert.Nil(t, postErr)
	resBody, bodyErr = ioutil.ReadAll(postResponse.Body)
	assert.Nil(t, bodyErr)
	var page bdata.BeaconV1DataPage
	jsonErr = json.Unmarshal(resBody, &page)
	assert.Nil(t, jsonErr)
	assert.NotNil(t, page)
	assert.Len(t, page.Data, 2)
	beacon1 = *page.Data[0]

	// Update the beacon
	beacon1.Label = "ABC"
	bodyMap = make(map[string]interface{})
	bodyMap["beacon"] = beacon1

	jsonBody, _ = json.Marshal(bodyMap)
	bodyReader = bytes.NewReader(jsonBody)
	postResponse, postErr = http.Post(url+"/v1/beacons/update_beacon", "application/json", bodyReader)
	assert.Nil(t, postErr)
	resBody, bodyErr = ioutil.ReadAll(postResponse.Body)
	assert.Nil(t, bodyErr)

	jsonErr = json.Unmarshal(resBody, &beacon)
	assert.Nil(t, jsonErr)
	assert.NotNil(t, beacon)
	assert.Equal(t, Beacon1.Id, beacon.Id)
	assert.Equal(t, "ABC", beacon.Label)

	// Get beacon by udi
	bodyMap = make(map[string]interface{})
	bodyMap["udi"] = beacon1.Udi

	jsonBody, _ = json.Marshal(bodyMap)
	bodyReader = bytes.NewReader(jsonBody)
	postResponse, postErr = http.Post(url+"/v1/beacons/get_beacon_by_udi", "application/json", bodyReader)
	assert.Nil(t, postErr)
	resBody, bodyErr = ioutil.ReadAll(postResponse.Body)
	assert.Nil(t, bodyErr)

	jsonErr = json.Unmarshal(resBody, &beacon)
	assert.Nil(t, jsonErr)
	assert.NotNil(t, beacon)
	assert.Equal(t, Beacon1.Id, beacon.Id)

	// Calculate position for one beacon
	bodyMap = make(map[string]interface{})
	bodyMap["site_id"] = "1"
	bodyMap["udis"] = []string{"00001"}

	jsonBody, _ = json.Marshal(bodyMap)
	bodyReader = bytes.NewReader(jsonBody)
	postResponse, postErr = http.Post(url+"/v1/beacons/calculate_position", "application/json", bodyReader)
	assert.Nil(t, postErr)
	resBody, bodyErr = ioutil.ReadAll(postResponse.Body)
	assert.Nil(t, bodyErr)

	var position bdata.GeoPointV1
	jsonErr = json.Unmarshal(resBody, &position)
	assert.Nil(t, jsonErr)
	assert.NotNil(t, beacon)

	assert.NotNil(t, position)
	assert.Equal(t, "Point", position.Type)
	assert.Equal(t, (float32)(0.0), position.Coordinates[0][0])
	assert.Equal(t, (float32)(0.0), position.Coordinates[0][1])

	// Delete the beacon
	bodyMap = make(map[string]interface{})
	bodyMap["beacon_id"] = beacon1.Id

	jsonBody, _ = json.Marshal(bodyMap)
	bodyReader = bytes.NewReader(jsonBody)
	postResponse, postErr = http.Post(url+"/v1/beacons/delete_beacon_by_id", "application/json", bodyReader)
	assert.Nil(t, postErr)
	resBody, bodyErr = ioutil.ReadAll(postResponse.Body)
	assert.Nil(t, bodyErr)

	jsonErr = json.Unmarshal(resBody, &beacon)
	assert.Nil(t, jsonErr)
	assert.NotNil(t, beacon)
	assert.Equal(t, Beacon1.Id, beacon.Id)

	// Try to get deleted beacon
	bodyMap = make(map[string]interface{})
	bodyMap["beacon_id"] = beacon1.Id

	jsonBody, _ = json.Marshal(bodyMap)
	bodyReader = bytes.NewReader(jsonBody)
	postResponse, postErr = http.Post(url+"/v1/beacons/get_beacon_by_id", "application/json", bodyReader)
	assert.Nil(t, postErr)
	resBody, bodyErr = ioutil.ReadAll(postResponse.Body)
	assert.Nil(t, bodyErr)

	beacon = bdata.BeaconV1{}

	jsonErr = json.Unmarshal(resBody, &beacon)
	assert.Nil(t, jsonErr)
	assert.NotNil(t, beacon)
	assert.Empty(t, beacon)
}
