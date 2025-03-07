package config

import "github.com/spf13/viper"

type App struct {
	AppPort      string `json:"app_port"`
	AppEnv       string `json:"app_env"`
	JwtSecretKey string `json:"jwt_secret_key"`
	JwtIssuer    string `json:"jwt_issuer"`
}

type PgsqlDB struct {
	DbHost     string `json:"db_host"`
	DbPort     string `json:"db_port"`
	DbUser     string `json:"db_user"`
	DbPassword string `json:"db_password"`
	DbName     string `json:"db_name"`
	DbMaxOpen  int    `json:"db_max_open"`
	DbMaxIdle  int    `json:"db_max_idle"`
}

type CloudflareR2 struct {
	BucketName string `json:"bucket_name"`
	ApiKey     string `json:"api_key"`
	ApiSecret  string `json:"api_secret"`
	Token      string `json:"token"`
	AccountID  string `json:"account_id"`
	PublicUrl  string `json:"public_url"`
}

type Config struct {
	App     App
	PgsqlDB PgsqlDB
	R2      CloudflareR2
}

func NewConfig() *Config {
	return &Config{
		App: App{
			AppPort:      viper.GetString("APP_PORT"),
			AppEnv:       viper.GetString("APP_ENV"),
			JwtSecretKey: viper.GetString("JWT_SECRET_KEY"),
			JwtIssuer:    viper.GetString("JWT_ISSUER"),
		},
		PgsqlDB: PgsqlDB{
			DbHost:     viper.GetString("DATABASE_HOST"),
			DbPort:     viper.GetString("DATABASE_PORT"),
			DbUser:     viper.GetString("DATABASE_USER"),
			DbPassword: viper.GetString("DATABASE_PASSWORD"),
			DbName:     viper.GetString("DATABASE_NAME"),
			DbMaxOpen:  viper.GetInt("DATABASE_MAX_OPEN_CONNECTION"),
			DbMaxIdle:  viper.GetInt("DATABASE_MAX_IDLE_CONNECTION"),
		},
		R2: CloudflareR2{
			BucketName: viper.GetString("CLOUDFLARE_R2_BUCKET_NAME"),
			ApiKey:     viper.GetString("CLOUDFLARE_R2_API_KEY"),
			ApiSecret:  viper.GetString("CLOUDFLARE_R2_API_SECRET"),
			Token:      viper.GetString("CLOUDFLARE_R2_TOKEN"),
			AccountID:  viper.GetString("CLOUDFLARE_R2_ACCOUNT_ID"),
			PublicUrl:  viper.GetString("CLOUDFLARE_R2_PUBLIC_URL"),
		},
	}
}
