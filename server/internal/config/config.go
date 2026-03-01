package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/spf13/viper"
)

type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	Logger   LoggerConfig   `mapstructure:"logger"`
	LLM      LLMConfig      `mapstructure:"llm"`
}

type ServerConfig struct {
	Port int    `mapstructure:"port"`
	Mode string `mapstructure:"mode"` // debug, release
	Host string `mapstructure:"host"`
}

type DatabaseConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	DBName   string `mapstructure:"dbname"`
	SSLMode  string `mapstructure:"sslmode"`
}

type LoggerConfig struct {
	Level  string `mapstructure:"level"`  // debug, info, warn, error
	Format string `mapstructure:"format"` // json, console
}

type LLMConfig struct {
	DefaultProvider string            `mapstructure:"default_provider"`
	Providers      map[string]Provider `mapstructure:"providers"`
}

type Provider struct {
	BaseURL string `mapstructure:"base_url"`
	APIKey  string `mapstructure:"api_key"`
}

func Load(configPath string) (*Config, error) {
	viper.SetConfigFile(configPath)
	viper.SetConfigType("yaml")

	setDefaults()

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		// 配置文件不存在时仅使用默认值和环境变量
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			if _, statErr := os.Stat(configPath); os.IsNotExist(statErr) {
				// 文件不存在，使用默认值
			} else {
				return nil, fmt.Errorf("读取配置文件失败: %w", err)
			}
		}
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("解析配置失败: %w", err)
	}

	// 环境变量覆盖（Docker 场景）
	applyEnvOverrides(&config)

	return &config, nil
}

// applyEnvOverrides 支持 Docker 等场景通过环境变量覆盖配置
func applyEnvOverrides(c *Config) {
	if v := os.Getenv("SERVER_PORT"); v != "" {
		if port, err := strconv.Atoi(v); err == nil {
			c.Server.Port = port
		}
	}
	if v := os.Getenv("SERVER_MODE"); v != "" {
		c.Server.Mode = v
	}
	if v := os.Getenv("DB_HOST"); v != "" {
		c.Database.Host = v
	}
	if v := os.Getenv("DB_PORT"); v != "" {
		if port, err := strconv.Atoi(v); err == nil {
			c.Database.Port = port
		}
	}
	if v := os.Getenv("DB_USER"); v != "" {
		c.Database.User = v
	}
	if v := os.Getenv("DB_PASSWORD"); v != "" {
		c.Database.Password = v
	}
	if v := os.Getenv("DB_NAME"); v != "" {
		c.Database.DBName = v
	}
	if v := os.Getenv("DB_SSLMODE"); v != "" {
		c.Database.SSLMode = v
	}
	if v := os.Getenv("LOG_LEVEL"); v != "" {
		c.Logger.Level = v
	}
	if v := os.Getenv("LOG_FORMAT"); v != "" {
		c.Logger.Format = v
	}
}

func setDefaults() {
	viper.SetDefault("server.port", 8080)
	viper.SetDefault("server.mode", "debug")
	viper.SetDefault("server.host", "0.0.0.0")

	viper.SetDefault("database.host", "localhost")
	viper.SetDefault("database.port", 5432)
	viper.SetDefault("database.user", "postgres")
	viper.SetDefault("database.password", "")
	viper.SetDefault("database.dbname", "x_novel")
	viper.SetDefault("database.sslmode", "disable")

	viper.SetDefault("logger.level", "info")
	viper.SetDefault("logger.format", "console")

	viper.SetDefault("llm.default_provider", "openai")
}

func (c *Config) GetDSN() string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		c.Database.Host,
		c.Database.Port,
		c.Database.User,
		c.Database.Password,
		c.Database.DBName,
		c.Database.SSLMode,
	)
}

func (c *Config) GetServerAddr() string {
	return fmt.Sprintf("%s:%d", c.Server.Host, c.Server.Port)
}
