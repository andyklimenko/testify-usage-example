package api

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/andyklimenko/testify-usage-example/api/storage"
	"github.com/andyklimenko/testify-usage-example/api/storage/database"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/suite"
)

type srvSuite struct {
	suite.Suite

	repo    *storage.Storage
	httpCli *http.Client
}

func TestServer(t *testing.T) {
	t.Parallel()
	suite.Run(t, &srvSuite{})
}

func (s *srvSuite) SetupSuite() {
	s.httpCli = &http.Client{Timeout: time.Second}
	db := database.DB()
	s.repo = storage.New(db)
}

func (s *srvSuite) setupServer(changelog userChangelog) (string, func()) {
	srv := &Server{
		repo:          s.repo,
		userChangelog: changelog,
	}
	testSrv := httptest.NewServer(setupRouter(srv))
	srv.httpSrv = testSrv.Config

	return testSrv.URL, func() {
		testSrv.Close()
	}
}

func TestMain(m *testing.M) {
	closer, repoErr := database.InitDockerDB()
	if repoErr != nil {
		os.Exit(-1)
	}
	status := m.Run()
	if closer != nil {
		closer()
	}
	os.Exit(status)
}
