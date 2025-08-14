package service

import (
	"reflect"
	"strings"

	"go-template/src/core/db_los"
	"go-template/src/core/minio"
	"go-template/src/core/smtp_service"
	"go-template/src/puppeteer"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
	"go-template/src/core/azure_ad"
	"go-template/src/core/db"
	"go-template/src/core/dpis_service"
	"go-template/src/core/log"
	"go-template/src/custom_error"
)

var (
	uni      *ut.UniversalTranslator
	trans    ut.Translator
	validate *validator.Validate
)

func init() {
	en := en.New()
	uni = ut.New(en, en)

	// this is usually known or extracted from http 'Accept-Language' header
	// also see uni.FindTranslator(...)
	translator, found := uni.GetTranslator("en")
	if !found {
		panic("translator not found")
	}

	validate = validator.New()

	if err := en_translations.RegisterDefaultTranslations(validate, translator); err != nil {
		panic(err)
	}

	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	trans = translator
}

type Service struct {
	Config       *Config
	Logger       log.Logger
	DB           db.DB
	DB_LOS       db_los.DB
	AzureAD      azure_ad.AzureADService
	UserID       int64
	AzureUserID  string
	Role         []string
	EmailAddress string
	ProfilePic   string
	SpanID       string
	TraceID      string
	DpisService  dpis_service.DpisService
	Minio        minio.MinIO
	SmtpService  *smtp_service.SmtpServiceClient
	Puppeteer    puppeteer.Puppeteer
}

func NewService(logger log.Logger) (service *Service, err error) {
	service = &Service{
		Logger: logger,
	}
	service.Config, err = InitConfig()
	if err != nil {
		return nil, err
	}

	dbConfig, err := db.InitConfig()
	if err != nil {
		return nil, err
	}

	service.DB, err = db.New(dbConfig, logger)
	if err != nil {
		return nil, err
	}

	dbLOSConfig, err := db_los.InitConfig()
	if err != nil {
		return nil, err
	}

	service.DB_LOS, err = db_los.New(dbLOSConfig, logger)
	if err != nil {
		return nil, err
	}

	service.Minio, err = minio.New(logger)
	if err != nil {
		return nil, err
	}

	service.Puppeteer, err = puppeteer.New(logger)
	if err != nil {
		return nil, err
	}

	err = service.Minio.CreateDefaultBucket()
	if err != nil {
		return nil, err
	}

	return service, nil
}

func ValidateInput(input interface{}) *custom_error.ValidationError {
	err := validate.Struct(input)
	if err != nil {
		messages := make([]string, 0)
		for _, e := range err.(validator.ValidationErrors) {
			messages = append(messages, e.Translate(trans))
		}
		errMessage := strings.Join(messages, ", ")
		return &custom_error.ValidationError{
			Code:    custom_error.InputValidationError,
			Message: errMessage,
		}
	}
	return nil
}
