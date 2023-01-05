package config

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gopkg.in/go-playground/validator.v9"
)

type (
	// Endpoints defines multiple API base-urls to fetch the data
	Endpoints struct {
		// EthRPCEndpoint      string `mapstructure:"eth_rpc_endpoint"`
		// BorRPCEndpoint      string `mapstructure:"bor_rpc_end_point"`
		// BorExternalRPC      string `mapstructure:"bor_external_rpc"`
		HeimdallRPCEndpoint string `mapstructure:"heimdall_rpc_endpoint"`
		HeimdallLCDEndpoint string `mapstructure:"heimdall_lcd_endpoint"`
		// HeimdallExternalRPC string `mapstructure:"heimdall_external_rpc"`
	}

	StatsDetails struct {
		SecretKey         string `mapstructure:"secret_key"`
		Node              string `mapstructure:"node"`
		NetStatsIPAddress string `mapstructure:"net_stats_ip"`
		Port              int    `mapstructure:"port"`
		Host              string `mapstructure:"host"`
	}

	// Config defines all the configurations required for the app
	Config struct {
		Endpoints    Endpoints    `mapstructure:"rpc_and_lcd_endpoints"`
		StatsDetails StatsDetails `mapstructure:"stats_details"`
	}
)

// ReadFromFile to read config details using viper
func ReadFromFile() (*Config, error) {

	configPath := `/var/lib/telemetry/config/`
	log.Printf("Config Path : %s", configPath)

	v := viper.New()
	v.AddConfigPath(".")
	v.AddConfigPath(configPath)
	v.SetConfigName("config")
	if err := v.ReadInConfig(); err != nil {
		log.Fatalf("error while reading config.toml: %v", err)
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		log.Fatalf("error unmarshaling config.toml to application config: %v", err)
	}

	if err := cfg.Validate(); err != nil {
		log.Fatalf("error occurred in config validation: %v", err)
	}

	return &cfg, nil
}

// Validate config struct
func (c *Config) Validate(e ...string) error {
	v := validator.New()
	if len(e) == 0 {
		return v.Struct(c)
	}
	return v.StructExcept(c, e...)
}
