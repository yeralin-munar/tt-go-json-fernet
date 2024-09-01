package tests

import (
	"context"
	"os"
	"path/filepath"

	"github.com/yeralin-munar/tt-go-json-fernet/internal/constant"
)

// Before running the test, you need to run docker service
func (s *TTGoJsonFernetTestSuite) TestDataCollectorUseCase_Successful_Data_Collect() {
	ctx := context.Background()

	err := s.jsonGeneratorUseCase.GenerateJsonFile(ctx)
	s.Require().NoError(err)

	// Check if the file exists
	c, err := os.ReadDir(filepath.Join(constant.DataFolder))
	s.Require().NoError(err)
	s.Require().NotEmpty(c)

	// Print files in the folder
	s.T().Logf("[BEFORE] Files in the folder %q:", constant.DataFolder)
	for _, entry := range c {
		s.T().Logf("\t- %s", entry.Name())
	}

	// Read the file
	err = s.dataCollectorUseCase.CollectData(ctx)
	s.Require().NoError(err)

	// Check if the file exists
	c, err = os.ReadDir(filepath.Join(constant.DataFolder))
	s.Require().NoError(err)
	s.Require().NotEmpty(c)

	// Print files in the folder
	s.T().Logf("[AFTER] Files in the folder %q:", constant.DataFolder)
	for _, entry := range c {
		if entry.IsDir() {
			childDir, err := os.ReadDir(filepath.Join(constant.DataFolder, entry.Name()))
			s.Require().NoError(err)
			s.Require().NotEmpty(childDir)
			s.T().Logf("\t- %s:", entry.Name())
			for _, entry2 := range childDir {
				s.T().Logf("\t\t- %s", entry2.Name())
			}
			continue
		}
		s.T().Logf("\t- %s", entry.Name())
	}

	// Check if the data is collected
	fileData, _, err := s.fileDataRepo.Query(ctx, nil)
	s.Require().NoError(err)
	s.Require().NotEmpty(fileData)

	// Print data
	s.T().Logf("Data in the database:")
	for _, data := range fileData {
		s.T().Logf("\t- %s: %v", data.Name, data.Data)
	}
}
