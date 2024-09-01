package tests

import (
	"context"
	"os"
	"path/filepath"

	"github.com/yeralin-munar/tt-go-json-fernet/internal/constant"
)

// Before running the test, you need to run docker service
func (s *TTGoJsonFernetTestSuite) TestJSONGeneratorUseCase_Successful_File_Generation() {
	ctx := context.Background()

	err := s.jsonGeneratorUseCase.GenerateJsonFile(ctx)
	s.Require().NoError(err)

	// Check if the file exists
	c, err := os.ReadDir(filepath.Join(constant.DataFolder))
	s.Require().NoError(err)
	s.Require().NotEmpty(c)

	// Print files in the folder
	s.T().Log("Files in the folder:")
	for _, entry := range c {
		s.T().Logf("\t- %s", entry.Name())
	}
}
