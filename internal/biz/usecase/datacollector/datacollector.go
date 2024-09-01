package datacollector

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/yeralin-munar/tt-go-json-fernet/config"
	"github.com/yeralin-munar/tt-go-json-fernet/internal/biz/dto"
	"github.com/yeralin-munar/tt-go-json-fernet/internal/biz/usecase"
	"github.com/yeralin-munar/tt-go-json-fernet/internal/constant"
)

type jsonFile struct {
	Name string
	Data []byte
}

type DataCollectorUseCase struct {
	conf *config.Data
	log  *log.Helper

	transactionManager usecase.TransactionManager

	cipherator      usecase.Cipherator
	scrapingKeyRepo usecase.ScrapingKeyRepo
	fileDataRepo    usecase.FileDataRepo
}

func NewDataCollectorUseCase(
	conf *config.Data,
	logger log.Logger,
	transactionManager usecase.TransactionManager,
	cipherator usecase.Cipherator,
	scrapingKeyRepo usecase.ScrapingKeyRepo,
	fileDataRepo usecase.FileDataRepo,
) *DataCollectorUseCase {
	return &DataCollectorUseCase{
		conf:               conf,
		log:                log.NewHelper(logger),
		transactionManager: transactionManager,
		scrapingKeyRepo:    scrapingKeyRepo,
		cipherator:         cipherator,
		fileDataRepo:       fileDataRepo,
	}
}

func (uc *DataCollectorUseCase) CollectData(ctx context.Context) error {
	// [x] Load scraping keys
	scrapingKeys, err := uc.loadScrapingKeys(ctx)
	if err != nil {
		return fmt.Errorf("failed to load scraping keys: %w", err)
	}
	// [x] Load files from folder
	files, err := uc.loadFilesFromFolder(ctx)
	if err != nil {
		return fmt.Errorf("failed to load files from folder: %w", err)
	}
	// [x] Decrypt files
	files, err = uc.decryptFiles(ctx, files)
	if err != nil {
		return fmt.Errorf("failed to decrypt files: %w", err)
	}
	// [x] Get file data from files using scraping keys
	fileDatas, err := uc.getDataFromFiles(ctx, files, scrapingKeys)
	if err != nil {
		return fmt.Errorf("failed to get file data from files: %w", err)
	}
	// [x] Store file data in database
	err = uc.storeFileData(ctx, fileDatas)
	if err != nil {
		return fmt.Errorf("failed to store file data: %w", err)
	}
	// [ ] Move files from data folder to collected folder
	err = uc.moveFiles(ctx, files)
	if err != nil {
		return fmt.Errorf("failed to move files: %w", err)
	}

	return nil
}

func (uc *DataCollectorUseCase) loadScrapingKeys(ctx context.Context) (
	[]*dto.ScrapingKey,
	error,
) {
	scrappingKeysDao, _, err := uc.scrapingKeyRepo.Query(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to load scraping keys: %w", err)
	}

	return dto.ConvertScrapingKeyDaosToDtos(scrappingKeysDao), nil
}

func (uc *DataCollectorUseCase) loadFilesFromFolder(ctx context.Context) (
	[]*jsonFile,
	error,
) {
	// [x] Read files from folder
	c, err := os.ReadDir(filepath.Join(constant.DataFolder))
	if err != nil {
		return nil, fmt.Errorf("failed to read folder: %w", err)
	}

	// [x] Store files in jsonFile struct
	var files []*jsonFile
	for _, entry := range c {
		if entry.IsDir() {
			continue
		}

		filename := filepath.Join(constant.DataFolder, entry.Name())
		data, err := os.ReadFile(filename)
		if err != nil {
			return nil, fmt.Errorf("failed to read file %q: %w", filename, err)
		}

		files = append(files, &jsonFile{
			Name: filename,
			Data: data,
		})
	}

	return files, nil
}

func (uc *DataCollectorUseCase) decryptFiles(ctx context.Context, files []*jsonFile) (
	[]*jsonFile,
	error,
) {
	// [x] Decrypt files
	for _, file := range files {
		decryptedData, err := uc.cipherator.Decrypt(file.Data, uc.conf.Crypto.Key)
		if err != nil {
			return nil, fmt.Errorf("failed to decrypt file %q: %w", file.Name, err)
		}

		file.Data = decryptedData
	}

	return files, nil
}

func (uc *DataCollectorUseCase) getDataFromFiles(
	ctx context.Context,
	files []*jsonFile,
	scrapingKeys []*dto.ScrapingKey,
) (
	[]*dto.FileData,
	error,
) {
	// [x] Get file data from files using scraping keys
	var fileDatas []*dto.FileData
	for _, file := range files {
		fileData, err := uc.getDataFromFile(ctx, file, scrapingKeys)
		if err != nil {
			return nil, fmt.Errorf("failed to get file data from file %q: %w", file.Name, err)
		}

		fileDatas = append(fileDatas, fileData)
	}

	return fileDatas, nil
}

// Get file data from file using scraping keys
func (uc *DataCollectorUseCase) getDataFromFile(
	ctx context.Context,
	file *jsonFile,
	scrapingKeys []*dto.ScrapingKey,
) (
	*dto.FileData,
	error,
) {
	jsonData := make(map[string]interface{})
	err := json.Unmarshal(file.Data, &jsonData)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal file data: %w", err)
	}

	uc.log.Debugf("File data: %+v", jsonData)

	data := make(map[string]interface{})
	for _, key := range scrapingKeys {
		// [x] Get data from file by key
		value, ok := jsonData[key.Name]
		if !ok {
			continue
		}
		// [x] Add data to map
		data[key.Name] = value
	}

	var fileData *dto.FileData
	if len(data) > 0 {
		fileData = &dto.FileData{
			Name: file.Name,
			Data: data,
		}
	}

	return fileData, nil
}

func (uc *DataCollectorUseCase) storeFileData(ctx context.Context, fileDatas []*dto.FileData) error {
	fileDataDaos := dto.ConvertFileDataDtosToDao(fileDatas)
	// Transaction
	insertErr := uc.transactionManager.InTransaction(ctx, func(txCtx context.Context) error {
		_, err := uc.fileDataRepo.Insert(txCtx, fileDataDaos)
		if err != nil {
			return fmt.Errorf("failed to store file data: %w", err)
		}

		return nil

	})

	return insertErr
}

// Move files from data folder to collected folder (<data folder>/collected)
func (uc *DataCollectorUseCase) moveFiles(ctx context.Context, files []*jsonFile) error {
	// Check if the folder exists
	err := os.MkdirAll(filepath.Join(constant.DataFolder, constant.WorkFinishFolder), 0755)
	if err != nil {
		return fmt.Errorf("failed to create folder: %w", err)
	}

	for _, file := range files {
		_, filename := filepath.Split(file.Name)
		err := os.Rename(
			file.Name,
			filepath.Join(constant.DataFolder, constant.WorkFinishFolder, filename),
		)
		if err != nil {
			return fmt.Errorf("failed to move file %q: %w", file.Name, err)
		}
	}

	return nil
}
