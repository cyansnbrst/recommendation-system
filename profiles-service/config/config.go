package config

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/spf13/viper"
)

// App config struct
type Config struct {
	Port             int
	Env              string
	AuthURL          string
	DefaultLocation  string
	DefaultInterests string
	PostgreSQL       PostgreSQL
	Kafka            Kafka
	Timeout          Timeout
}

// PostgreSQL config struct
type PostgreSQL struct {
	Host         string
	Port         int
	User         string
	Password     string
	DBName       string
	SSLMode      string
	MaxOpenConns int
	MaxIdleConns int
	MaxIdleTime  time.Duration
}

// Kafka config struct
type Kafka struct {
	Brokers     []string
	Topics      map[string]string
	MaxAttempts int
}

// Timeouts config struct
type Timeout struct {
	PostgreSQLConn   time.Duration
	PostgreSQLAction time.Duration
	ServerIdle       time.Duration
	ServerRead       time.Duration
	ServerWrite      time.Duration
	ServerShutdown   time.Duration
}

// Load config file from given path
func LoadConfig(filename string) (*viper.Viper, error) {
	v := viper.New()

	v.SetConfigName(filename)
	v.SetConfigType("env")
	v.AddConfigPath(".")
	v.AutomaticEnv()

	v.SetEnvPrefix("")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if err := v.ReadInConfig(); err != nil {
		var configNotFoundError viper.ConfigFileNotFoundError
		if errors.As(err, &configNotFoundError) {
			return nil, errors.New("config file not found")
		}
		return nil, err
	}

	return v, nil
}

func ParseConfig(v *viper.Viper) (*Config, error) {
	var c Config
	var err error

	// Server config
	c.Port = v.GetInt("port")
	c.Env = v.GetString("env")
	c.AuthURL = v.GetString("auth_url")
	c.DefaultLocation = v.GetString("default_location")
	c.DefaultInterests = v.GetString("default_interests")

	// PostgreSQL config
	c.PostgreSQL.Host = v.GetString("postgresql_host")
	c.PostgreSQL.Port = v.GetInt("postgresql_port")
	c.PostgreSQL.User = v.GetString("postgresql_user")
	c.PostgreSQL.Password = v.GetString("postgresql_password")
	c.PostgreSQL.DBName = v.GetString("postgresql_db")
	c.PostgreSQL.SSLMode = v.GetString("postgresql_sslmode")
	c.PostgreSQL.MaxOpenConns = v.GetInt("postgresql_max_open_conns")
	c.PostgreSQL.MaxIdleConns = v.GetInt("postgresql_max_idle_conns")
	c.PostgreSQL.MaxIdleTime, err = parseTimeout(v, "postgresql_max_idle_time")
	if err != nil {
		return nil, err
	}

	// Kafka config
	if brokers := v.GetString("kafka_brokers"); brokers != "" {
		c.Kafka.Brokers = strings.Split(brokers, ",")
	}
	c.Kafka.Topics = make(map[string]string)
	for _, key := range v.AllKeys() {
		if strings.HasPrefix(key, "kafka_topic_") {
			topicName := strings.TrimPrefix(key, "kafka_topic_")
			c.Kafka.Topics[topicName] = v.GetString(key)
		}
	}
	c.Kafka.MaxAttempts = v.GetInt("kafka_max_attempts")

	// Timeout config
	c.Timeout.PostgreSQLConn, err = parseTimeout(v, "timeout_postgresql_conn")
	if err != nil {
		return nil, err
	}
	c.Timeout.PostgreSQLAction, err = parseTimeout(v, "timeout_postgresql_action")
	if err != nil {
		return nil, err
	}
	c.Timeout.ServerIdle, err = parseTimeout(v, "timeout_server_idle")
	if err != nil {
		return nil, err
	}
	c.Timeout.ServerRead, err = parseTimeout(v, "timeout_server_read")
	if err != nil {
		return nil, err
	}
	c.Timeout.ServerWrite, err = parseTimeout(v, "timeout_server_write")
	if err != nil {
		return nil, err
	}
	c.Timeout.ServerShutdown, err = parseTimeout(v, "timeout_server_shutdown")
	if err != nil {
		return nil, err
	}

	return &c, nil
}

func parseTimeout(v *viper.Viper, key string) (time.Duration, error) {
	durationStr := v.GetString(key)
	duration, err := time.ParseDuration(durationStr)
	if err != nil {
		return 0, fmt.Errorf("invalid duration for %s: %w", key, err)
	}
	return duration, nil
}
