package clients

import (
	bdata "github.com/pip-services-samples/pip-data-microservice-go/data/version1"
	cdata "github.com/pip-services3-go/pip-services3-commons-go/data"
)

type IBeaconsClientV1 interface {
	GetBeacons(correlationId string, filter *cdata.FilterParams, paging *cdata.PagingParams) (page *bdata.BeaconV1DataPage, err error)

	GetBeaconById(correlationId string, beaconId string) (beacon *bdata.BeaconV1, err error)

	GetBeaconByUdi(correlationId string, udi string) (beacon *bdata.BeaconV1, err error)

	CalculatePosition(correlationId string, siteId string, udis []string) (position *bdata.GeoPointV1, err error)

	CreateBeacon(correlationId string, beacon bdata.BeaconV1) (res *bdata.BeaconV1, err error)

	UpdateBeacon(correlationId string, beacon bdata.BeaconV1) (res *bdata.BeaconV1, err error)

	DeleteBeaconById(correlationId string, beaconId string) (beacon *bdata.BeaconV1, err error)
}
