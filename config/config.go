package config

import (
	"fmt"
	"os"
	"reflect"

	"github.com/spf13/viper"
)

type Config struct {
	Name                       string `mapstructure:"NAME"`
	Environment                string `mapstructure:"ENVIRONMENT"`
	Port                       int    `mapstructure:"PORT"`
	Version                    string `mapstructure:"VERSION"`
	AppUrl                     string `mapstructure:"APP_URL"`
	AccessTokenSigningKey      string `mapstructure:"ACCESS_TOKEN_SIGNING_KEY"`
	AccessTokenTokenExpiration int    `mapstructure:"ACCESS_TOKEN_EXPIRATION"`
	RefreshTokenSigningKey     string `mapstructure:"REFRESH_TOKEN_SIGNING_KEY"`
	RefreshTokenExpiration     int    `mapstructure:"REFRESH_TOKEN_EXPIRATION"`
	DatabasePort               int    `mapstructure:"DATABASE_PORT"`
	DatabaseName               string `mapstructure:"DATABASE_NAME"`
	DatabaseHost               string `mapstructure:"DATABASE_HOST"`
	DatabaseUsername           string `mapstructure:"DATABASE_USERNAME"`
	DatabasePassword           string `mapstructure:"DATABASE_PASSWORD"`
	DatabaseSslMode            string `mapstructure:"DATABASE_SSL_MODE"`
	SMTPHost                   string `mapstructure:"SMTP_HOST"`
	SMTPPort                   int    `mapstructure:"SMTP_PORT"`
	SMTPUsername               string `mapstructure:"SMTP_USERNAME"`
	SMTPPassword               string `mapstructure:"SMTP_PASSWORD"`
	SMTPFrom                   string `mapstructure:"SMTP_FROM"`
	LogEnable                  bool   `mapstructure:"LOG_ENABLE"`
	MailEnable                 bool   `mapstructure:"MAIL_ENABLE"`
	SmsApiKey                  string `mapstructure:"SMS_API_KEY"`
	MaxLoginDeviceCount        int    `mapstructure:"MAX_LOGIN_DEVICE_COUNT"`
	AutoDeleteDevice           bool   `mapstructure:"AUTO_DELETE_DEVICE"`
	SendSMS                    bool   `mapstructure:"SEND_SMS"`
	SendEmail                  bool   `mapstructure:"SEND_EMAIL"`
	RoutePrefix                string
}

var AppConfig *Config = &Config{}

func Load() {
	viper.AddConfigPath(".")
	viper.SetConfigType("env")
	viper.SetConfigName(".env")

	viper.AutomaticEnv() // read in environment variables that match

	// read config from system environment
	elem := reflect.TypeOf(AppConfig).Elem()
	for i := 0; i < elem.NumField(); i++ {
		key := elem.Field(i).Tag.Get("mapstructure")
		value := os.Getenv(key)
		if value != "" {
			viper.Set(key, value)
		}
	}

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}

	err := viper.Unmarshal(&AppConfig)

	if err != nil {
		fmt.Fprintln(os.Stderr, "Config unmarshal error: ", err)
	}
}
