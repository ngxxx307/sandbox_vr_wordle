package config

import (
	"fmt"
	"os"
	"reflect"
	"regexp"
	"strings"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type ValidationError struct {
	Fields         []string
	ValidationTags []string
}

func (v ValidationError) Error() string {
	errStr := "Failed validation on fields: \n"
	for i, field := range v.Fields {
		errStr += fmt.Sprintf("Field %v failed on condition %v\n", field, v.ValidationTags[i])
	}
	return errStr
}

type Regex struct {
	Name string
	Exp  *regexp.Regexp
}

type Config struct {
	ServerPort      string `mapstructure:"SERVER_PORT" validate:"required"`
	PingIntervalSec int    `mapstructure:"PING_INTERVAL_SEC" validate:"required"`
	PingTimeoutSec  int    `mapstructure:"PING_TIMEOUT_SEC" validate:"required"`
	PongWaitSec     int    `mapstructure:"PONG_WAIT_SEC" validate:"required"`
	MaxMessageSize  int64  `mapstructure:"MAX_MESSAGE_SIZE" validate:"required"`
}

func NewConfig() (*Config, error) {
	config := &Config{}
	if err := LoadConfig(config); err != nil {
		return nil, err
	}
	return config, nil
}

func LoadConfig(config interface{}) error {
	env := os.Getenv("DEPLOY_ENV")
	if env == "" {
		env = "local"
	}
	envPath := "./env/.env." + strings.ToLower(env)
	if testing.Testing() {
		envPath = "." + envPath
	}
	if err := godotenv.Load(envPath); err != nil {
		return err
	}
	bindEnvs(config)
	if err := viper.Unmarshal(&config); err != nil {
		return err
	}
	return Validate(config)
}

func bindEnvs(iface interface{}, parts ...string) {
	ifv := reflect.ValueOf(iface)
	if ifv.Kind() == reflect.Ptr {
		ifv = ifv.Elem()
	}
	ift := reflect.TypeOf(iface)
	if ift.Kind() == reflect.Ptr {
		ift = ift.Elem()
	}
	for i := 0; i < ift.NumField(); i++ {
		v := ifv.Field(i)
		t := ift.Field(i)
		tv, ok := t.Tag.Lookup("mapstructure")
		if !ok {
			continue
		}
		switch v.Kind() {
		case reflect.Struct:
			bindEnvs(v.Interface(), append(parts, tv)...)
		default:
			viper.BindEnv(strings.Join(append(parts, tv), "."))
		}
	}
}

func Validate(data interface{}) error {
	v := validator.New()

	err := v.Struct(data)
	if err == nil {
		return nil
	}

	validationErrors, ok := err.(validator.ValidationErrors)
	if !ok {
		return err
	}

	var accumulatedErrors = ValidationError{}
	for _, validationErr := range validationErrors {
		accumulatedErrors.Fields = append(accumulatedErrors.Fields, validationErr.Field())
		accumulatedErrors.ValidationTags = append(accumulatedErrors.ValidationTags, validationErr.Tag())
	}
	return accumulatedErrors
}
