package http_test

import (
	"context"
	"math/rand"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/ifraixedes/go-ms-http-nsq-example/zombie-driver/http"
	"github.com/ifraixedes/go-ms-http-nsq-example/zombie-driver/http/internal"
	"github.com/ifraixedes/go-ms-http-nsq-example/zombie-driver/mock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestServerClientIntegration(t *testing.T) {
	var (
		dID  = uint64(rand.Int63n(100) + 400)
		ctrl = gomock.NewController(t)
		svc  = mock.NewService(ctrl)
	)
	defer ctrl.Finish()

	svc.EXPECT().IsZombie(gomock.Any(), dID).Return(true, nil)

	var (
		httpConfig = internal.HTTPConfig{
			Addr: testHttpAddr,
		}
		ctx = context.Background()
	)

	var sdown, err = internal.ServerUp(svc, httpConfig)
	require.NoError(t, err)
	defer func() {
		require.NoError(t, sdown(ctx))
	}()

	cli, err := http.NewClient(httpConfig.Addr)
	require.NoError(t, err)

	isz, err := cli.IsZombie(ctx, dID)
	require.NoError(t, err)
	assert.True(t, isz)

	// TODO: IMPROVEMENT
	// Rather than sleeping, we should have a wrapper around svc for detecting
	// when it's called and after finish the test. This change will allow to have
	// false negative due bigger time delays than the one used here.
	time.Sleep(1 * time.Second)
}
