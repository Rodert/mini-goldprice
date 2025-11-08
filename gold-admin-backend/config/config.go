package config

import (
	"bytes"
	"embed"
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	JWT      JWTConfig      `mapstructure:"jwt"`
	Log      LogConfig      `mapstructure:"log"`
}

type ServerConfig struct {
	Port int    `mapstructure:"port"`
	Mode string `mapstructure:"mode"`
}

type DatabaseConfig struct {
	Type         string `mapstructure:"type"`
	Path         string `mapstructure:"path"`
	Host         string `mapstructure:"host"`
	Port         int    `mapstructure:"port"`
	Username     string `mapstructure:"username"`
	Password     string `mapstructure:"password"`
	Database     string `mapstructure:"database"`
	Charset      string `mapstructure:"charset"`
	MaxIdleConns int    `mapstructure:"max_idle_conns"`
	MaxOpenConns int    `mapstructure:"max_open_conns"`
}

type JWTConfig struct {
	Secret string `mapstructure:"secret"`
	Expire int    `mapstructure:"expire"`
}

type LogConfig struct {
	Level string `mapstructure:"level"`
	Path  string `mapstructure:"path"`
}

var AppConfig *Config

// LoadConfigFromFile 从文件加载配置（兼容旧方法）
func LoadConfigFromFile(configPath string) error {
	viper.Reset()
	viper.SetConfigFile(configPath)
	viper.SetConfigType("yaml")

	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("读取配置文件失败: %w", err)
	}

	AppConfig = &Config{}
	if err := viper.Unmarshal(AppConfig); err != nil {
		return fmt.Errorf("解析配置文件失败: %w", err)
	}

	return nil
}

// LoadConfigFromEmbed 从 embed.FS 加载配置
func LoadConfigFromEmbed(embedFS embed.FS, configPath string) error {
	viper.Reset()
	viper.SetConfigType("yaml")

	// 从 embed.FS 读取文件内容
	data, err := embedFS.ReadFile(configPath)
	if err != nil {
		return fmt.Errorf("读取嵌入配置文件失败: %w", err)
	}

	// 将内容读取到 viper
	if err := viper.ReadConfig(bytes.NewReader(data)); err != nil {
		return fmt.Errorf("解析嵌入配置失败: %w", err)
	}

	AppConfig = &Config{}
	if err := viper.Unmarshal(AppConfig); err != nil {
		return fmt.Errorf("解析配置结构失败: %w", err)
	}

	return nil
}

// LoadConfig 兼容旧方法（保持向后兼容）
func LoadConfig(configPath string) error {
	return LoadConfigFromFile(configPath)
}
