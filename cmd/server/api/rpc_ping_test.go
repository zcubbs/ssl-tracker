package api

import (
	"context"
	"github.com/stretchr/testify/require"
	mockdb "github.com/zcubbs/ssl-tracker/cmd/server/db/mock"
	mockwk "github.com/zcubbs/ssl-tracker/cmd/server/worker/mock"
	pb "github.com/zcubbs/ssl-tracker/pb"
	"go.uber.org/mock/gomock"
	"testing"
)

func TestPing(t *testing.T) {
	testCases := []struct {
		name          string
		req           *pb.Empty
		checkResponse func(t *testing.T, res *pb.PingResponse, err error)
	}{
		{
			name: "OK",
			req:  &pb.Empty{},
			checkResponse: func(t *testing.T, res *pb.PingResponse, err error) {
				require.NoError(t, err)
				require.Equal(t, "Pong", res.Message)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			storeCtrl := gomock.NewController(t)
			defer storeCtrl.Finish()

			store := mockdb.NewMockStore(storeCtrl)
			worker := mockwk.NewMockTaskDistributor(storeCtrl)

			s := newTestServer(t, store, worker)
			res, err := s.Ping(context.Background(), tc.req)
			tc.checkResponse(t, res, err)
		})
	}
}
