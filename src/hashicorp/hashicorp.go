package hashicorp

import (
	"context"
	"encoding/base64"
	"fmt"
	"go-template/src/core/log"
	"reflect"
	"strings"

	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/api/auth/approle"
	"github.com/hashicorp/vault/api/auth/userpass"
)

//go:generate mockery --name Interface
type Interface interface {
	GetData(field string) (string, error)
	GetDataFromField(input interface{}) error
}

type Client struct {
	client *api.Client
	log    log.Logger
	Config *Config
}

func New(config *Config, logger log.Logger) (Interface, error) {

	client, err := getClient(config, logger)
	if err != nil {
		return nil, fmt.Errorf("error creating client: %s", err)
	}

	return &Client{
		client: client,
		log: logger.WithFields(log.Fields{
			"module": "hashicorp",
		}),
		Config: config,
	}, nil
}

func getClient(config *Config, logger log.Logger) (*api.Client, error) {
	client, err := api.NewClient(&api.Config{Address: config.URL})
	if err != nil {
		return nil, fmt.Errorf("error creating client: %s", err)
	}

	if config.AuthMethod == Token {
		client.SetToken(config.TokenMethod.Token)

		return client, nil
	}

	var auth api.AuthMethod
	switch config.AuthMethod {
	case UserPassword:
		auth, err = userpass.NewUserpassAuth(config.UserPassMethod.User, &userpass.Password{FromString: config.UserPassMethod.Password})
		if err != nil {
			return nil, fmt.Errorf("unable to initialize Userpass auth method: %w", err)
		}

	case AppRole:
		auth, err = approle.NewAppRoleAuth(config.AppRoleMethod.RoleID, &approle.SecretID{FromString: config.AppRoleMethod.SecretID}, approle.WithMountPath(config.AppRoleMethod.Path))
		if err != nil {
			return nil, fmt.Errorf("unable to initialize AppRole auth method: %w", err)
		}
	default:
		return nil, fmt.Errorf("unknown auth method: %s", config.AuthMethod)
	}

	// Authenticate and obtain a token
	authInfo, err := client.Auth().Login(context.Background(), auth)
	if err != nil {
		return nil, fmt.Errorf("unable to login to Userpass auth method: %w", err)
	}
	if authInfo == nil {
		return nil, fmt.Errorf("no auth info was returned after login")
	}

	// Set the client token
	client.SetToken(authInfo.Auth.ClientToken)

	// renew token
	watcher, err := client.NewLifetimeWatcher(&api.LifetimeWatcherInput{
		Secret:    authInfo,
		Increment: authInfo.Auth.LeaseDuration - 100, //3600-100 second
	})
	if err != nil {
		logger.Errorf("unable to watch auth info: %s", err)
	}

	go func() {
		go watcher.Start()
		defer watcher.Stop()
		for {
			select {
			case err := <-watcher.DoneCh():
				logger.Errorf("unable to renew token: %s", err)

			case renewal := <-watcher.RenewCh():
				sec, _ := renewal.Secret.TokenTTL()
				logger.Infof("auth info watcher renewed : %s", sec)
			}
		}
	}()

	return client, nil
}

// GetData retrieves a value from HashiCorp Vault using a formatted string.
//
// The format of the input must be: 'hashicorp:path:field[:option...]'
//
// - path: the Vault path where the secret is stored
// - field: the key within the secret to retrieve
// - options (optional): post-processing flags such as 'decodeBase64'
//
// Multiple options can be specified, separated by colons.
//
// Supported options:
//   - decodeBase64: decodes the retrieved base64 string
//   - otherOption: (reserved for future use)
//
// If the input string does not follow the 'hashicorp:' format, it is returned as-is.
//
// Example:
//
//	"hashicorp:service-a/api-secrets:username"
//	"hashicorp:service-a/api-secrets:password:decodeBase64"
//	"hashicorp:service-a/api-secrets:password:decodeBase64:otherOption"
func (client *Client) GetData(field string) (string, error) {
	if len(field) == 0 {
		return "", fmt.Errorf("no field specified")
	}
	fieldList := strings.Split(field, ":")
	if len(fieldList) < 3 || fieldList[0] != "hashicorp" {
		return field, nil
	}
	mOption := make(map[option]bool)
	for i, f := range fieldList {
		if i < 3 {
			continue // skip base field
		}
		mOption[option(f)] = true
	}

	path := fieldList[1]
	fieldName := fieldList[2]

	secret, err := client.client.Logical().Read(path)
	if err != nil {
		return "", fmt.Errorf("error reading secret: %s", err)
	}

	// Check if the secret was found
	if secret == nil || secret.Data["data"] == nil {
		return "", fmt.Errorf("no secret found")
	}

	// Access the key-value pairs in the secret data
	data := secret.Data["data"].(map[string]interface{})
	val, ok := data[fieldName]
	if !ok {
		return "", fmt.Errorf("no field %s found", fieldName)
	}

	valStr, ok := val.(string)
	if !ok {
		return "", fmt.Errorf("field %s is not a string", fieldName)
	}

	if mOption[optionDecodeBase64] {
		valByte, err := base64.StdEncoding.DecodeString(valStr)
		if err != nil {
			return "", fmt.Errorf("error decoding secret: %s", err)
		}
		valStr = string(valByte)
	}

	return valStr, nil
}

// GetDataFromField traverses the input data and replaces string values that match the hashicorp format
// with the retrieved secret. It supports structs, pointers, slices, and strings.
func (client *Client) GetDataFromField(input interface{}) error {
	v := reflect.ValueOf(input)
	if !v.IsValid() {
		return nil
	}
	return client.processValue(v)
}

func (client *Client) processValue(v reflect.Value) error {
	if !v.IsValid() {
		return nil
	}

	switch v.Kind() {
	case reflect.Ptr:
		if v.IsNil() {
			return nil
		}
		return client.processValue(v.Elem())

	case reflect.Struct:
		for i := 0; i < v.NumField(); i++ {
			field := v.Field(i)
			if field.CanSet() {
				if err := client.processValue(field); err != nil {
					return err
				}
			}
		}

	case reflect.Slice:
		for i := 0; i < v.Len(); i++ {
			if err := client.processValue(v.Index(i)); err != nil {
				return err
			}
		}

	case reflect.Map:
		// only handle maps with string keys
		if v.Type().Key().Kind() != reflect.String {
			return nil
		}
		for _, key := range v.MapKeys() {
			orig := v.MapIndex(key)
			if !orig.IsValid() {
				continue
			}

			// unwrap interface{} values
			val := orig
			if val.Kind() == reflect.Interface && !val.IsNil() {
				val = val.Elem()
			}

			var newVal reflect.Value
			switch val.Kind() {
			case reflect.String:
				// fetch & replace raw string values
				data, err := client.GetData(val.String())
				if err != nil {
					return err
				}
				newVal = reflect.ValueOf(data)

			case reflect.Ptr, reflect.Struct, reflect.Slice, reflect.Map:
				// deep-copy and recurse
				copy := reflect.New(val.Type()).Elem()
				copy.Set(val)
				if err := client.processValue(copy); err != nil {
					return err
				}
				newVal = copy

			default:
				// leave other primitive kinds (int, bool, etc.) unchanged
				newVal = val
			}

			v.SetMapIndex(key, newVal)
		}

	case reflect.String:
		s := v.String()
		data, err := client.GetData(s)
		if err != nil {
			return err
		}
		v.SetString(data)

	default:
		// nothing to do for other Kinds
	}

	return nil
}
