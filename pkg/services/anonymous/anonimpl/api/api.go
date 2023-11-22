package api

import (
	"net/http"
	"time"

	"github.com/grafana/grafana/pkg/api/response"
	"github.com/grafana/grafana/pkg/api/routing"
	"github.com/grafana/grafana/pkg/infra/log"
	"github.com/grafana/grafana/pkg/services/accesscontrol"
	"github.com/grafana/grafana/pkg/services/anonymous"
	"github.com/grafana/grafana/pkg/services/anonymous/anonimpl/anonstore"
	contextmodel "github.com/grafana/grafana/pkg/services/contexthandler/model"
	"github.com/grafana/grafana/pkg/setting"
)

type AnonDeviceServiceAPI struct {
	cfg            *setting.Cfg
	store          anonstore.AnonStore
	accesscontrol  accesscontrol.AccessControl
	RouterRegister routing.RouteRegister
	log            log.Logger
}

func NewAnonDeviceServiceAPI(
	cfg *setting.Cfg,
	anonstore anonstore.AnonStore,
	accesscontrol accesscontrol.AccessControl,
	routerRegister routing.RouteRegister,
) *AnonDeviceServiceAPI {
	return &AnonDeviceServiceAPI{
		cfg:            cfg,
		store:          anonstore,
		accesscontrol:  accesscontrol,
		RouterRegister: routerRegister,
		log:            log.New("anon.api"),
	}
}

func (api *AnonDeviceServiceAPI) RegisterAPIEndpoints() {
	auth := accesscontrol.Middleware(api.accesscontrol)
	api.RouterRegister.Group("/api/anonymous", func(anonRoutes routing.RouteRegister) {
		anonRoutes.Get("/anonstats", auth(accesscontrol.EvalPermission(accesscontrol.ActionServerStatsRead)), routing.Wrap(api.CountDevices))
	})
}

// type AnonStore interface {
// 	// ListDevices returns all devices that have been updated between the given times.
// 	ListDevices(ctx context.Context, from *time.Time, to *time.Time) ([]*Device, error)
// 	// CreateOrUpdateDevice creates or updates a device.
// 	CreateOrUpdateDevice(ctx context.Context, device *Device) error
// 	// CountDevices returns the number of devices that have been updated between the given times.
// 	CountDevices(ctx context.Context, from time.Time, to time.Time) (int64, error)
// 	// DeleteDevice deletes a device by its ID.
// 	DeleteDevice(ctx context.Context, deviceID string) error
// 	// DeleteDevicesOlderThan deletes all devices that have no been updated since the given time.
// 	DeleteDevicesOlderThan(ctx context.Context, olderThan time.Time) error
// }

// ListDevices returns all devices that have been updated between the given times.
// func (a *AnonDeviceServiceAPI) ListDevices(ctx context.Context) ([]*anonstore.Device, error) {
// 	//, from *time.Time, to *time.Time
// 	from := time.Now()

// 	return a.service.ListDevices(ctx, from, to)
// }

// swagger:route GET /anonusers signed_in_user getSignedInUser
//
// Get (get all anon users)
//
// Responses:
// 200: userResponse
// 401: unauthorisedError
// 403: forbiddenError
// 404: notFoundError
// 500: internalServerError
func (api *AnonDeviceServiceAPI) CountDevices(c *contextmodel.ReqContext) response.Response {
	fromTime := time.Now().Add(-anonymous.ThirtyDays)
	toTime := time.Now().Add(time.Minute)

	devices, err := api.store.CountDevices(c.Req.Context(), fromTime, toTime)
	if err != nil {
		return response.Error(http.StatusInternalServerError, "Failed to list devices", err)
	}
	return response.JSON(http.StatusOK, devices)
}
