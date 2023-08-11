package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/require"
	mockdb "github.com/zcubbs/tlz/db/mock"
	db "github.com/zcubbs/tlz/db/sqlc"
	"github.com/zcubbs/tlz/util"
	"go.uber.org/mock/gomock"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestGetDomains(t *testing.T) {
	domain := randomDomain()

	testCases := []struct {
		name          string
		domainName    string
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, response *http.Response)
	}{
		{
			name:       "OK",
			domainName: domain.Name,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetDomain(gomock.Any(), gomock.Eq(domain.Name)).
					Times(1).
					Return(domain, nil)
			},
			checkResponse: func(t *testing.T, response *http.Response) {
				require.Equal(t, http.StatusOK, response.StatusCode)
				requireBodyMatchDomain(t, response.Body, domain)
			},
		},
		{
			name:       "NotFound",
			domainName: domain.Name,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetDomain(gomock.Any(), gomock.Eq(domain.Name)).
					Times(1).
					Return(db.Domain{}, sql.ErrNoRows)
			},
			checkResponse: func(t *testing.T, response *http.Response) {
				require.Equal(t, http.StatusNotFound, response.StatusCode)
			},
		},
		{
			name:       "InternalError",
			domainName: domain.Name,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetDomain(gomock.Any(), gomock.Eq(domain.Name)).
					Times(1).
					Return(db.Domain{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, response *http.Response) {
				require.Equal(t, http.StatusInternalServerError, response.StatusCode)
			},
		},
		{
			name:       "InvalidDomainName",
			domainName: "invalid",
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetDomain(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, response *http.Response) {
				require.Equal(t, http.StatusBadRequest, response.StatusCode)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			tc.buildStubs(store)

			server := newTestServer(t, store)
			app := fiber.New()
			path := "/api/domains/:name"
			app.Get(path, server.GetDomain)

			pathWithParams := fmt.Sprintf("/api/domains/%s", tc.domainName)
			eq := httptest.NewRequest(fiber.MethodGet, pathWithParams, nil)
			resp, err := app.Test(eq, -1)
			require.NoError(t, err)
			tc.checkResponse(t, resp)
		})
	}
}

func TestCreateDomain(t *testing.T) {
	domain := randomDomain()
	domainRequest := CreateDomainRequest{
		Name: domain.Name,
	}

	testCases := []struct {
		name          string
		request       CreateDomainRequest
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, response *http.Response)
	}{
		{
			name:    "OK",
			request: domainRequest,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					InsertDomain(gomock.Any(), gomock.Eq(domainRequest.Name)).
					Times(1).
					Return(domain, nil)
			},
			checkResponse: func(t *testing.T, response *http.Response) {
				require.Equal(t, http.StatusOK, response.StatusCode)
				requireBodyMatchDomain(t, response.Body, domain)
			},
		},
		{
			name:    "InternalError",
			request: domainRequest,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					InsertDomain(gomock.Any(), gomock.Eq(domainRequest.Name)).
					Times(1).
					Return(db.Domain{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, response *http.Response) {
				require.Equal(t, http.StatusInternalServerError, response.StatusCode)
			},
		},
		{
			name:    "InvalidDomainName",
			request: CreateDomainRequest{Name: "invalid"},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					InsertDomain(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, response *http.Response) {
				require.Equal(t, http.StatusBadRequest, response.StatusCode)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			tc.buildStubs(store)

			server := newTestServer(t, store)
			app := fiber.New()
			path := "/api/domains"
			app.Post(path, server.CreateDomain)

			body, err := json.Marshal(tc.request)
			require.NoError(t, err)

			eq := httptest.NewRequest(fiber.MethodPost, path, bytes.NewReader(body))
			resp, err := app.Test(eq, -1)
			require.NoError(t, err)
			tc.checkResponse(t, resp)
		})
	}
}

func randomDomain() db.Domain {
	return db.Domain{
		Name:              util.RandomDomainName(),
		CertificateExpiry: sql.NullTime{},
		Status:            sql.NullString{},
		Issuer:            sql.NullString{},
		Owner:             util.RandomOwner(),
		CreatedAt:         time.Time{},
	}
}

func requireBodyMatchDomain(t *testing.T, body io.ReadCloser, domain db.Domain) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var got db.Domain
	err = json.Unmarshal(data, &got)
	require.NoError(t, err)

	require.Equal(t, got.Name, domain.Name)
	require.Equal(t, got.CertificateExpiry, domain.CertificateExpiry)
	require.Equal(t, got.Status, domain.Status)
	require.Equal(t, got.Issuer, domain.Issuer)
	require.Equal(t, got.Owner, domain.Owner)
	require.WithinDuration(t, got.CreatedAt, domain.CreatedAt, time.Second)
}
