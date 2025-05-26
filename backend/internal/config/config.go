package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type MySQLConfig struct {
	Host      string
	Port      int
	User      string
	Password  string
	DBName    string `mapstructure:"dbname"`
	Charset   string
	ParseTime bool `mapstructure:"parse_time"`
	Loc       string
}

type DatabaseConfig struct {
	MySQL MySQLConfig
}

type ServerConfig struct {
	Public struct {
		Port      int
		UploadDir string `mapstructure:"upload_dir"`
	}
	Admin struct {
		Port      int
		JWTSecret string `mapstructure:"jwt_secret"`
		Password  string `mapstructure:"password"`
	}
}

type ConfigStruct struct {
	Server   ServerConfig
	Database DatabaseConfig
}

var Config ConfigStruct

func InitConfig() error {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config")
	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("读取配置文件失败: %w", err)
	}
	if err := viper.Unmarshal(&Config); err != nil {
		return fmt.Errorf("解析配置文件失败: %w", err)
	}
	return nil
}
