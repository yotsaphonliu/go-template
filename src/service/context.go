package service

import (
	"bytes"
	"io"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"go-template/src/core/azure_ad"
	"go-template/src/core/db"
	"go-template/src/core/log"
	"go-template/src/core/minio"
	"go-template/src/core/smtp_service"
	"go-template/src/puppeteer"
)

const (
	// UserKey user key
	UserKey = "user"
	// LangKey lang key
	LangKey = "lang"
	// DatabaseKey database key
	DatabaseKey = "database"
	// ParametersKey parameters key
	ParametersKey = "parameters"
	// SessionTokenKey session token key
	SessionTokenKey = "sessionToken"
)

// Context context
type Context struct {
	*fiber.Ctx
	Config       *Config
	Logger       log.Logger
	DB           db.DB
	DBLOS        db_los.DB
	AzureAD      azure_ad.AzureADService
	UserID       int64
	AzureUserID  string
	Role         []string
	EmailAddress string
	ProfilePic   string
	SpanID       string
	TraceID      string
	SmtpService  *smtp_service.SmtpServiceClient
	MinIO        minio.MinIO
}

// New new custom fiber context
func (service *Service) NewContext(c *fiber.Ctx) *Context {
	return &Context{
		Ctx:          c,
		Config:       service.Config,
		Logger:       service.Logger,
		DB:           service.DB,
		DBLOS:        service.DB_LOS,
		AzureAD:      service.AzureAD,
		UserID:       service.UserID,
		AzureUserID:  service.AzureUserID,
		Role:         service.Role,
		EmailAddress: service.EmailAddress,
		ProfilePic:   service.ProfilePic,
		SpanID:       service.SpanID,
		TraceID:      service.TraceID,
		DpisService:  service.DpisService,
		MinIO:        service.Minio,
		SmtpService:  service.SmtpService,
	}
}

func (ctx *Context) getLogger(funcName string) log.Logger {
	return ctx.Logger.WithFields(log.Fields{
		"func":     funcName,
		"span_id":  ctx.SpanID,
		"trace_id": ctx.TraceID,
	})
}

func (ctx *Context) CreateActivityLog(uniqueNo string, reqBody, resBody []byte) error {
	err := ctx.DB.CreateActivityLog(ctx.GetServiceCode(), uniqueNo, reqBody, resBody)
	if err != nil {
		return err
	}

	return nil
}

func (ctx *Context) GetServiceCode() string {
	name := ctx.Route().Name
	if len(name) == 0 {
		name = ctx.OriginalURL() // default value
	}

	return name
}

func detectContentType(r io.Reader) (string, error) {
	buf := make([]byte, 512)
	n, err := r.Read(buf)
	if err != nil && err != io.EOF {
		return "", err
	}

	// Restore reader
	r = io.MultiReader(bytes.NewReader(buf[:n]), r)

	// Detect MIME type
	return http.DetectContentType(buf[:n]), nil
}
