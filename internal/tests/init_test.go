package tests

import (
	"context"
	"os"
	"testing"

	"github.com/fernet/fernet-go"
	"github.com/yeralin-munar/tt-go-json-fernet/config"
	"github.com/yeralin-munar/tt-go-json-fernet/internal/constant"
	innerpg "github.com/yeralin-munar/tt-go-json-fernet/internal/data/postgres"
	innerfernet "github.com/yeralin-munar/tt-go-json-fernet/internal/encryption/fernet"
	"github.com/yeralin-munar/tt-go-json-fernet/internal/service"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/yeralin-munar/tt-go-json-fernet/internal/biz/usecase"
	"github.com/yeralin-munar/tt-go-json-fernet/internal/biz/usecase/datacollector"
	"github.com/yeralin-munar/tt-go-json-fernet/internal/biz/usecase/jsongenerator"
)

type TTGoJsonFernetTestSuite struct {
	suite.Suite

	logger      log.Logger
	pgContainer *postgres.PostgresContainer
	cfg         *config.Data
	db          *innerpg.DB

	// Repositories
	transactionManager usecase.TransactionManager
	scrapingKeyRepo    usecase.ScrapingKeyRepo
	fileDataRepo       usecase.FileDataRepo

	cipherator usecase.Cipherator

	//UseCases
	dataCollectorUseCase service.DataCollectorUseCase
	jsonGeneratorUseCase service.JsonGeneratorUseCase
}

// Executes before test suite begins execution
func (s *TTGoJsonFernetTestSuite) SetupSuite() {
	s.logger = log.NewStdLogger(os.Stdout)

	s.cfg, s.pgContainer = GetTestPostgresConnection(s.T(), "../../migrations")
}

// Executes after all tests executed
func (s *TTGoJsonFernetTestSuite) TearDownSuite() {
	// Close GRPC connection
	// Delete test folder
	os.RemoveAll(constant.DataFolder)
}

// Executes before each test
func (s *TTGoJsonFernetTestSuite) SetupTest() {
	var err error
	// Connect to the database
	s.db, err = innerpg.NewDB(s.cfg)
	if err != nil {
		s.T().Fatalf("failed to connect to database: %s", err)
	}

	// s.cfg.Folder = "./test_data"
	var key fernet.Key
	err = key.Generate()
	if err != nil {
		s.T().Fatalf("failed to generate key: %s", err)
	}
	s.cfg.Crypto.Key = string(key.Encode())

	s.transactionManager = innerpg.NewTransactionManager(s.db, s.logger)

	// Init repositories
	s.scrapingKeyRepo = innerpg.NewScrapingKeyRepo(s.db, s.logger)
	s.fileDataRepo = innerpg.NewFileDataRepo(s.db, s.logger)

	s.cipherator = innerfernet.NewFernetCipherator()

	// Init usecases
	s.dataCollectorUseCase = datacollector.NewDataCollectorUseCase(
		s.cfg,
		s.logger,
		s.transactionManager,
		s.cipherator,
		s.scrapingKeyRepo,
		s.fileDataRepo,
	)

	s.jsonGeneratorUseCase = jsongenerator.NewJsonGeneratorUseCase(
		s.cfg,
		s.logger,
		s.cipherator,
		s.scrapingKeyRepo,
	)
}

// Executes after each test
func (s *TTGoJsonFernetTestSuite) TearDownTest() {
	ctx := context.Background()

	// Restore the database to the initial state
	err := s.pgContainer.Restore(ctx)
	if err != nil {
		s.T().Fatalf("failed to restore container: %s", err)
	}
}

func TestHelloworldTestSuite(t *testing.T) {
	suite.Run(t, new(TTGoJsonFernetTestSuite))
}
