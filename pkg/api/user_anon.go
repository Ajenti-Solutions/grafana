package api

import (
	"fmt"
	"net/http"
	"time"

	"github.com/grafana/grafana/pkg/api/response"
	contextmodel "github.com/grafana/grafana/pkg/services/contexthandler/model"
)

// swagger:route GET /api/admin/anonstats admin adminGetStats
//
// Get (get all anon users)
//
// Responses:
// 200: userResponse
// 401: unauthorisedError
// 403: forbiddenError
// 404: notFoundError
// 500: internalServerError
func (hs *HTTPServer) GetAnonUsers(c *contextmodel.ReqContext) response.Response {
	fmt.Printf("hit")
	// from Time.Time
	// to Time.Time
	from := time.Now().Add(-time.Hour * 24 * 30)
	to := time.Now().Add(time.Hour * 24 * 30)
	devices, err := hs.anonService.ListDevices(c.Req.Context(), &from, &to)
	if err != nil {
		return response.Error(http.StatusInternalServerError, "Failed to list devices", err)
	}
	return response.JSON(http.StatusOK, devices)
}

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
func (hs *HTTPServer) GetAnonDeviceCount(c *contextmodel.ReqContext) response.Response {
	from := c.Query("from")
	if from == "" {
		return response.Error(http.StatusBadRequest, "from is required", nil)
	}
	to := c.Query("to")
	if to == "" {
		return response.Error(http.StatusBadRequest, "to is required", nil)
	}
	fromTime, err := time.Parse(time.RFC3339, from)
	if err != nil {
		return response.Error(http.StatusBadRequest, "from is not a valid RFC3339 time", nil)
	}
	toTime, err := time.Parse(time.RFC3339, to)
	if err != nil {
		return response.Error(http.StatusBadRequest, "to is not a valid RFC3339 time", nil)
	}

	devices, err := hs.anonService.CountDevices(c.Req.Context(), fromTime, toTime)
	if err != nil {
		return response.Error(http.StatusInternalServerError, "Failed to list devices", err)
	}
	return response.JSON(http.StatusOK, devices)
}
