package gateway

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	drvloc "github.com/ifraixedes/go-ms-http-nsq-example/driver-location"
	dmock "github.com/ifraixedes/go-ms-http-nsq-example/driver-location/mock"
	zmock "github.com/ifraixedes/go-ms-http-nsq-example/zombie-driver/mock"
	"github.com/mmcloughlin/spherand"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGateway_Integration_RouteEndpoints(t *testing.T) {
	var (
		dID      = uint64(rand.Int63n(99) + 1)
		lat, lng = spherand.Geographical()
		locrj    = fmt.Sprintf(`{"latitude":%.5f,"longitude":%.5f}`, lat, lng)
		c        = Config{
			setDriverLocation: configEndpoint{
				Method: http.MethodPatch,
				Path:   "/drivers/:id/locations",
			},
			getDriver: configEndpoint{
				Method: http.MethodGet,
				Path:   "/drivers/:id",
			},
		}
	)

	t.Run("OK", func(t *testing.T) {
		var (
			ctrl  = gomock.NewController(t)
			dlsvc = dmock.NewService(ctrl)
			zdsvc = zmock.NewService(ctrl)
		)
		defer ctrl.Finish()

		dlsvc.EXPECT().SetLocation(gomock.Any(), dID, mockFuncMatcher{
			Func: func(v interface{}) bool {
				l := v.(drvloc.Location)

				if !assert.InEpsilon(t, lat, l.Lat, 0.000009) {
					return false
				}
				if !assert.InEpsilon(t, lng, l.Lng, 0.000009) {
					return false
				}

				return assert.WithinDuration(t, time.Now(), l.At, time.Second)
			},
		})
		zdsvc.EXPECT().IsZombie(gomock.Any(), dID).Return(true, nil)

		var svr = httptest.NewServer(routeEndpoints(c, &endpoints{
			dlsvc: dlsvc,
			zdsvc: zdsvc,
		}))
		defer svr.Close()

		var path = strings.Replace(c.setDriverLocation.Path, ":id", strconv.Itoa(int(dID)), 1)
		var req, err = http.NewRequest(
			c.setDriverLocation.Method,
			fmt.Sprintf("%s%s", svr.URL, path),
			bytes.NewReader([]byte(locrj)),
		)
		require.NoError(t, err)

		var hc = &http.Client{}
		res, err := hc.Do(req)
		require.NoError(t, err)

		assert.Equal(t, http.StatusNoContent, res.StatusCode)
		rawb, err := ioutil.ReadAll(res.Body)
		defer res.Body.Close()
		require.NoError(t, err)
		assert.Empty(t, rawb)

		path = strings.Replace(c.getDriver.Path, ":id", strconv.Itoa(int(dID)), 1)
		req, err = http.NewRequest(
			c.getDriver.Method,
			fmt.Sprintf("%s%s", svr.URL, path),
			nil,
		)
		require.NoError(t, err)

		res, err = hc.Do(req)
		require.NoError(t, err)

		assert.Equal(t, http.StatusOK, res.StatusCode)
		rawb, err = ioutil.ReadAll(res.Body)
		defer res.Body.Close()
		require.NoError(t, err)
		assert.Equal(t, fmt.Sprintf(`{"id":%d,"zombie":true}`, dID), string(rawb))
	})

	t.Run("error: set location error", func(t *testing.T) {
		t.Skipf("TODO: TO BE IMPLEMENTED")
	})

	t.Run("error: get driver not found", func(t *testing.T) {
		t.Skipf("TODO: TO BE IMPLEMENTED")
	})

	t.Run("error: get driver error", func(t *testing.T) {
		t.Skipf("TODO: TO BE IMPLEMENTED")
	})
}
