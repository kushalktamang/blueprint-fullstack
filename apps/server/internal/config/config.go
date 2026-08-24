package config

import (
	"os"
	"strings"

	"github.com/go-playground/validator/v10"
	// auto-loads .env into the process environment
	_ "github.com/joho/godotenv/autoload"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/v2"
	"github.com/rs/zerolog"
)

// root/parent application configuration built from the BLUEPRINT_* environment variables
type Config struct {
	Primary       Primary              `koanf:"primary" validate:"required"`
	Server        ServerConfig         `koanf:"server" validate:"required"`
	Database      DatabaseConfig       `koanf:"database" validate:"required"`
	Auth          AuthConfig           `koanf:"auth" validate:"required"`
	Redis         RedisConfig          `koanf:"redis" validate:"required"`
	Integration   IntegrationConfig    `koanf:"integration" validate:"required"`
	Observability *ObservabilityConfig `koanf:"observability"`
}

type Primary struct {
	Env string `koanf:"env" validate:"required"`
}

type ServerConfig struct {
	Port               string   `koanf:"port" validate:"required"`
	ReadTimeout        int      `koanf:"read_timeout" validate:"required"`
	WriteTimeout       int      `koanf:"write_timeout" validate:"required"`
	IdleTimeout        int      `koanf:"idle_timeout" validate:"required"`
	CORSAllowedOrigins []string `koanf:"cors_allowed_origins" validate:"required"`
}

type DatabaseConfig struct {
	// url if using a full DSN (e.g. a Neon connection string). If set, it takes
	// priority and Host/Port/User/Name/SSLMode become optional.
	URL             string `koanf:"url"`
	Host            string `koanf:"host" validate:"required_without=URL"`
	Port            int    `koanf:"port" validate:"required_without=URL"`
	User            string `koanf:"user" validate:"required_without=URL"`
	Password        string `koanf:"password"`
	Name            string `koanf:"name" validate:"required_without=URL"`
	SSLMode         string `koanf:"ssl_mode" validate:"required_without=URL"`
	MaxOpenConns    int    `koanf:"max_open_conns" validate:"required"`
	MaxIdleConns    int    `koanf:"max_idle_conns" validate:"required"`
	ConnMaxLifetime int    `koanf:"conn_max_lifetime" validate:"required"`
	ConnMaxIdleTime int    `koanf:"conn_max_idle_time" validate:"required"`
}

type RedisConfig struct {
	Address string `koanf:"address" validate:"required"`
}

type IntegrationConfig struct {
	ResendAPIKey     string `koanf:"resend_api_key" validate:"required"`
	EmailFromName    string `koanf:"email_from_name" validate:"required"`
	EmailFromAddress string `koanf:"email_from_address" validate:"required,email"`
}

type AuthConfig struct {
	SecretKey string `koanf:"secret_key" validate:"required"`
}

// this function reads BLUEPRINT_* environment variable into Config,
// validates it, and exits the process on any failure (via logger.fatal)
func MustLoad() *Config {
	logger := zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr}).With().Timestamp().Logger()

	// koanf instance
	k := koanf.New(".")

	// Load environment variables with prefix "BLUEPRINT_" and merge into config.
	// Example: BLUEPRINT_PARENT1_CHILD1_NAME becomes "parent1.child1.name"
	err := k.Load(env.ProviderWithValue("BLUEPRINT_", ".", func(key, value string) (string, any) {
		key = strings.ToLower(strings.TrimPrefix(key, "BLUEPRINT_"))
		if strings.Contains(value, " ") {
			return key, strings.Split(value, " ")
		}
		return key, value
	}), nil)
	if err != nil {
		logger.Fatal().Err(err).Msg("could not load initial env variables")
	}

	mainConfig := &Config{}

	err = k.Unmarshal("", mainConfig)
	if err != nil {
		logger.Fatal().Err(err).Msg("could not unmarshal main config")
	}

	validate := validator.New()

	err = validate.Struct(mainConfig)
	if err != nil {
		logger.Fatal().Err(err).Msg("config validation failed")
	}

	// Set default observability config if not provided
	if mainConfig.Observability == nil {
		mainConfig.Observability = DefaultObservabilityConfig()
	}

	// Override service name and environment from primary config
	mainConfig.Observability.ServiceName = "blueprint"
	mainConfig.Observability.Environment = mainConfig.Primary.Env

	// Validate observability config
	if err := mainConfig.Observability.Validate(); err != nil {
		logger.Fatal().Err(err).Msg("invalid observability config")
	}

	return mainConfig
}
