package jsongenerator

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"time"

	"github.com/brianvoe/gofakeit"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/yeralin-munar/tt-go-json-fernet/config"
	"github.com/yeralin-munar/tt-go-json-fernet/internal/biz/usecase"
	"github.com/yeralin-munar/tt-go-json-fernet/internal/constant"
)

const (
	AddNumDepth       = 2
	AddAdditionalKeys = 4
)

// Struct for JsonGeneratorUseCase
// Methods:
// - GenerateJsonFile
// - generateData map[string]interface{}
// - encryptData
// - createFile

// JsonGeneratorUseCase struct
type JsonGeneratorUseCase struct {
	conf *config.Data
	log  *log.Helper

	cipherator      usecase.Cipherator
	scrapingKeyRepo usecase.ScrapingKeyRepo
}

func NewJsonGeneratorUseCase(
	conf *config.Data,
	logger log.Logger,
	cipherator usecase.Cipherator,
	scrapingKeyRepo usecase.ScrapingKeyRepo,
) *JsonGeneratorUseCase {
	return &JsonGeneratorUseCase{
		conf:            conf,
		log:             log.NewHelper(logger),
		scrapingKeyRepo: scrapingKeyRepo,
		cipherator:      cipherator,
	}
}

// GenerateJsonFile generates a JSON file
func (uc *JsonGeneratorUseCase) GenerateJsonFile(ctx context.Context) error {
	uc.log.Infof("Generating JSON file...")
	// Implementation goes here
	// [x] Generate data
	generatedData, err := uc.generateData(ctx)
	if err != nil {
		return fmt.Errorf("failed to generate data: %w", err)
	}
	uc.log.Debugf("Generated data: %+v", generatedData)
	// [x] Marshal data
	jsonBytes, err := json.Marshal(generatedData)
	if err != nil {
		return fmt.Errorf("failed to marshal data: %w", err)
	}
	// [x] Encrypt data
	encryptedBytes, err := uc.encryptData(ctx, jsonBytes)
	if err != nil {
		return fmt.Errorf("failed to encrypt data: %w", err)
	}
	// [x] Create file
	err = uc.createFile(ctx, encryptedBytes)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}

	return nil
}

// generateData generates data for the JSON file
func (uc *JsonGeneratorUseCase) generateData(ctx context.Context) (
	map[string]interface{},
	error,
) {
	// Load scraping keys
	scrapingKeys, _, err := uc.scrapingKeyRepo.Query(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to load scraping keys: %w", err)
	}
	// Generate random map[string]interface{} data
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	// [x] Generate random maxDepth
	maxDepth := r.Intn(5) + AddNumDepth
	if maxDepth == 0 {
		maxDepth = 1
	}
	// [x] Generate random additional keys
	numberOfAdditionalKeys := r.Intn(maxDepth*2) + AddAdditionalKeys
	data := make(map[string]interface{})

	//  totalNumberOfKeys := len(scrapingKeys) + numberOfAdditionalKeys

	// Add scraping keys
	for _, key := range scrapingKeys {
		data[key.Name] = gofakeit.HackerAbbreviation()
	}

	// Add additional keys
	for i := 0; i < numberOfAdditionalKeys; i++ {
		key := gofakeit.HackerAdjective()
		depth := r.Intn(maxDepth)

		if depth == 1 {
			depth = 0
		}
		data[key] = generateDataInDepth(depth)
	}

	return data, nil
}

func generateDataInDepth(depth int) interface{} {
	if depth == 0 {
		return gofakeit.HackerIngverb()
	}

	numberOfKeys := rand.Intn(5) + AddAdditionalKeys
	data := make(map[string]interface{})

	for i := 0; i < numberOfKeys; i++ {
		key := gofakeit.HackerNoun()

		data[key] = generateDataInDepth(depth - 1)
	}

	return data
}

// encryptData encrypts the data
func (uc *JsonGeneratorUseCase) encryptData(ctx context.Context, data []byte) ([]byte, error) {
	encryptedBytes, err := uc.cipherator.Encrypt(data, uc.conf.Crypto.Key)
	if err != nil {
		return nil, err
	}

	return encryptedBytes, nil
}

// createFile creates a file with the given data
func (uc *JsonGeneratorUseCase) createFile(ctx context.Context, data []byte) error {
	// Check if the folder exists
	err := os.MkdirAll(filepath.Join(constant.DataFolder), 0755)
	if err != nil {
		return fmt.Errorf("failed to create folder: %w", err)
	}
	// if _, err := os.Stat(uc.conf.Folder); os.IsNotExist(err) {
	// 	err := os.Mkdir(uc.conf.Folder, 0755)
	// 	if err != nil {
	// 		return fmt.Errorf("failed to create folder: %w", err)
	// 	}
	// }
	filename := fmt.Sprint("data-", time.Now().Unix(), ".json")
	filename = fmt.Sprintf("%s/%s", constant.DataFolder, filename)

	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	_, err = file.Write(data)
	if err != nil {
		return fmt.Errorf("failed to write data to file: %w", err)
	}

	return nil
}
