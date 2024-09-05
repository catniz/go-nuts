package config

import (
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"github.com/sirupsen/logrus"
	"time"
)

type Config struct {
	PodName      string `envconfig:"POD_NAME" default:""`
	PodNamespace string `envconfig:"POD_NAMESPACE" default:""`
	ServerEnv    string `envconfig:"SERVER_ENV" default:"dev"`
	HttpPort     string `envconfig:"HTTP_PORT" default:"80"`
	LogLevel     string `envconfig:"LOG_LEVEL" default:"debug"`
	MySql        MysqlConfig
}

func (c Config) LogrusLevel() logrus.Level {
	if level, err := logrus.ParseLevel(c.LogLevel); err != nil {
		return logrus.InfoLevel
	} else {
		return level
	}
}

type MysqlConfig struct {
	Host         string `envconfig:"MYSQL_HOST" default:"localhost"`
	Port         string `envconfig:"MYSQL_PORT" default:"3306"`
	User         string `envconfig:"MYSQL_USER" default:""`
	Password     string `envconfig:"MYSQL_PASSWORD" default:""`
	Database     string `envconfig:"MYSQL_DATABASE" default:""`
	ConnMax      int    `envconfig:"MYSQL_CONNECTION_MAX" default:"4"`
	ConnTimeout  int    `envconfig:"MYSQL_CONNECTION_TIMEOUT" default:"1"`    // minutes, connection max lifetime
	ConnIdleTime int    `envconfig:"MYSQL_CONNECTION_IDLE_TIME" default:"10"` // minutes, idle connection close time
	LogDebugSQL  bool   `envconfig:"MYSQL_LOG_DEBUG_SQL" default:"false"`
}

func (m MysqlConfig) Dsn() string {
	// "<user>:<pass>@tcp(<host>:<port>)/<database>?parseTime=True"
	return m.User + ":" + m.Password + "@tcp(" + m.Host + ":" + m.Port + ")/" + m.Database + "?parseTime=True"
}

func (m MysqlConfig) ConnMaxTtl() time.Duration {
	return time.Duration(m.ConnTimeout) * time.Minute
}

func (m MysqlConfig) ConnIdleTtl() time.Duration {
	return time.Duration(m.ConnIdleTime) * time.Minute
}

// LoadConfig loads the configuration from the environment variables.
func LoadConfig(root string) (Config, error) {
	cfg := Config{}
	//Load the .env file
	_ = godotenv.Load(root + "/.env")

	if err := envconfig.Process("", &cfg); err != nil {
		return cfg, err
	}
	return cfg, nil
}
