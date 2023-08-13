package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/require"
	mockdb "github.com/zcubbs/tlz/db/mock"
	db "github.com/zcubbs/tlz/db/sqlc"
	"github.com/zcubbs/tlz/pkg/util"
	"go.uber.org/mock/gomock"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestGetDomain(t *testing.T) {
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
					Return(db.Domain{}, pgx.ErrNoRows)
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
					GetDomain(gomock.Any(), gomock.Eq(domainRequest.Name)).
					Times(1).
					Return(db.Domain{}, pgx.ErrNoRows)
				store.EXPECT().
					InsertDomain(gomock.Any(), gomock.Any()).
					Times(1).
					Return(domain, nil)
			},
			checkResponse: func(t *testing.T, response *http.Response) {
				require.Equal(t, http.StatusCreated, response.StatusCode)
				requireBodyMatchDomain(t, response.Body, domain)
			},
		},
		{
			name:    "InternalError",
			request: domainRequest,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetDomain(gomock.Any(), gomock.Eq(domainRequest.Name)).
					Times(1).
					Return(db.Domain{}, pgx.ErrNoRows)
				store.EXPECT().
					InsertDomain(gomock.Any(), gomock.Any()).
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

			req := httptest.NewRequest(fiber.MethodPost, path, bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			resp, err := app.Test(req, -1)
			require.NoError(t, err)
			tc.checkResponse(t, resp)
		})
	}
}

func TestGetDomains(t *testing.T) {
	domains := []db.Domain{
		randomDomain(),
		randomDomain(),
		randomDomain(),
	}

	testCases := []struct {
		name          string
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, response *http.Response)
	}{
		{
			name: "OK",
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAllDomains(gomock.Any()).
					Times(1).
					Return(domains, nil)
			},
			checkResponse: func(t *testing.T, response *http.Response) {
				require.Equal(t, http.StatusOK, response.StatusCode)
				var got []db.Domain
				err := json.NewDecoder(response.Body).Decode(&got)
				require.NoError(t, err)
				requireBodyMatchDomains(t, domains, got)
			},
		},
		{
			name: "InternalError",
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAllDomains(gomock.Any()).
					Times(1).
					Return([]db.Domain{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, response *http.Response) {
				require.Equal(t, http.StatusInternalServerError, response.StatusCode)
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
			app.Get(path, server.GetDomains)

			req := httptest.NewRequest(fiber.MethodGet, path, nil)
			resp, err := app.Test(req, -1)
			require.NoError(t, err)
			tc.checkResponse(t, resp)
		})
	}
}

func randomDomain() db.Domain {
	return db.Domain{
		Name:              util.RandomDomainName(),
		CertificateExpiry: pgtype.Timestamp{},
		Status:            pgtype.Text{},
		Issuer:            pgtype.Text{},
		Owner:             uuid.UUID{},
		CreatedAt:         time.Time{},
	}
}

func requireBodyMatchDomain(t *testing.T, body io.ReadCloser, domain db.Domain) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var got db.Domain
	err = json.Unmarshal(data, &got)
	require.NoError(t, err)

	requireDomainMatch(t, domain, got)
}

func requireBodyMatchDomains(t *testing.T, want, got []db.Domain) {
	require.Equal(t, len(want), len(got))
	for i := range want {
		requireDomainMatch(t, want[i], got[i])
	}
}

func requireDomainMatch(t *testing.T, want, got db.Domain) {
	require.Equal(t, want.Name, got.Name)
	require.Equal(t, want.Status.String, got.Status.String)
	require.Equal(t, want.Owner, got.Owner)
}
