package service

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/yeralin-munar/tt-go-json-fernet/config"
)

const (
	EncryptionMode = "ENCRYPTION"
	DecryptionMode = "DECRYPTION"
)

type DataCollectorUseCase interface {
	CollectData(ctx context.Context) error
}

type JsonGeneratorUseCase interface {
	GenerateJsonFile(ctx context.Context) error
}

type TTGoJsonFernetService struct {
	log  *log.Helper
	conf *config.Server

	dataCollectorUseCase DataCollectorUseCase
	jsonGeneratorUseCase JsonGeneratorUseCase
}

func NewTTGoJsonFernetService(
	logger log.Logger,
	conf *config.Server,
	dataCollectorUseCase DataCollectorUseCase,
	jsonGeneratorUseCase JsonGeneratorUseCase,
) *TTGoJsonFernetService {
	return &TTGoJsonFernetService{
		log:                  log.NewHelper(logger),
		conf:                 conf,
		dataCollectorUseCase: dataCollectorUseCase,
		jsonGeneratorUseCase: jsonGeneratorUseCase,
	}
}

func (s *TTGoJsonFernetService) Run(ctx context.Context) error {
	// s.log.Infof("server listening on %s", s.conf.Addr)
	s.log.Infof("Current mode: %s", s.conf.Mode)

	// http server
	http.HandleFunc("/run", s.handle(ctx))

	return http.ListenAndServe(s.conf.HTTP.Addr, nil)
}

func (s *TTGoJsonFernetService) handle(
	ctx context.Context,
) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if strings.EqualFold(s.conf.Mode, EncryptionMode) {
			err := s.jsonGeneratorUseCase.GenerateJsonFile(ctx)
			if err != nil {
				// return fmt.Errorf("Failed to work in encryption mode: %v", err)
				w.WriteHeader(http.StatusInternalServerError)
				_, err = w.Write([]byte(fmt.Sprintf("Failed to work in encryption mode: %v", err)))
				if err != nil {
					s.log.Errorf("Failed to write response: %v", err)
				}

				return
			}

			w.WriteHeader(http.StatusOK)
			_, err = w.Write([]byte("Successfully generated JSON file"))
			if err != nil {
				s.log.Errorf("Failed to write response: %v", err)
			}
		} else if strings.EqualFold(s.conf.Mode, DecryptionMode) {
			err := s.dataCollectorUseCase.CollectData(ctx)
			if err != nil {
				// return fmt.Errorf("Failed to work in decryption mode: %v", err)
				w.WriteHeader(http.StatusInternalServerError)
				_, err = w.Write([]byte(fmt.Sprintf("Failed to work in decryption mode: %v", err)))
				if err != nil {
					s.log.Errorf("Failed to write response: %v", err)
				}

				return
			}

			w.WriteHeader(http.StatusOK)
			_, err = w.Write([]byte("Successfully collected data"))
			if err != nil {
				s.log.Errorf("Failed to write response: %v", err)
			}
		} else {
			s.log.Errorf("Invalid mode: %s", s.conf.Mode)
		}
	}
}
