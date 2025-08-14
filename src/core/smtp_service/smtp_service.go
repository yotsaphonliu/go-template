package smtp_service

import (
	"go-template/src/core/log"
)

type SmtpService interface {
}

type SmtpServiceClient struct {
	logger log.Logger
	Config *Config
}

func New(config *Config, logger log.Logger) (serviceClient *SmtpServiceClient, err error) {

	serviceClient = &SmtpServiceClient{
		logger: logger.WithFields(log.Fields{
			"module": "smtp_service",
		}),
		Config: config,
	}

	return serviceClient, nil
}
