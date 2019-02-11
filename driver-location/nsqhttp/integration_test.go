package nsqhttp_test

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"os"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	drvloc "github.com/ifraixedes/go-ms-http-nsq-example/driver-location"
	"github.com/ifraixedes/go-ms-http-nsq-example/driver-location/mock"
	"github.com/ifraixedes/go-ms-http-nsq-example/driver-location/nsqhttp"
	"github.com/ifraixedes/go-ms-http-nsq-example/driver-location/nsqhttp/internal"
	"github.com/mmcloughlin/spherand"
	nsq "github.com/nsqio/go-nsq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestServerClientIntegration(t *testing.T) {
	var (
		ctx   = context.Background()
		topic = fmt.Sprintf("drvloc-%d", rand.Int())
		nsqc  = nsq.NewConfig()
	)

	var cli, err = nsqhttp.NewClient(testNSQdAddr, nsqc, topic, testHttpAddr)
	require.NoError(t, err)

	var (
		dID      = uint64(rand.Int63n(100) + 400)
		lat, lng = spherand.Geographical()
		loc      = drvloc.Location{
			Lat: lat,
			Lng: lng,
			At:  time.Now().Round(0),
		}
		ctrl = gomock.NewController(t)
		svc  = mock.NewService(ctrl)
	)
	defer ctrl.Finish()

	svc.EXPECT().SetLocation(gomock.Any(), dID, loc).Return(nil)

	// Make the NSQ client calls
	err = cli.SetLocation(ctx, dID, loc)
	require.NoError(t, err)

	// Run the server
	var (
		nsqSts = internal.NSQSettings{
			Topic:        topic,
			Channel:      fmt.Sprintf("drvloc-%d", rand.Int()),
			Cfg:          nsqc,
			LookupdAddrs: []string{testNSQLookupdAddr},
		}
		httpConfig = internal.HTTPConfig{
			Addr: testHttpAddr,
		}
		lg = log.New(os.Stderr, "", log.Ldate)
	)
	sdown, err := internal.ServerUp(svc, nsqSts, httpConfig, lg)
	require.NoError(t, err)
	defer func() {
		require.NoError(t, sdown(ctx))
	}()

	// Make the HTTP client call
	var mins = uint16(rand.Intn(50) + 1)
	svc.EXPECT().LocationsForLastMinutes(gomock.Any(), dID, mins).Return([]drvloc.Location{loc}, nil)
	locs, err := cli.LocationsForLastMinutes(ctx, dID, mins)
	require.NoError(t, err)
	assert.Equal(t, []drvloc.Location{loc}, locs)

	// TODO: IMPROVEMENT
	// Rather than sleeping, we should have a wrapper around svc for detecting
	// when it's called and after finish the test. This change will allow to have
	// false negative due bigger time delays than the one used here.
	time.Sleep(1 * time.Second)
}
