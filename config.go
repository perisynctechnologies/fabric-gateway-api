package assetTransfer

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	Server        Server `yaml:"server"`
	ChaincodeName string `yaml:"chaincodeName"`
	ChannelName   string `yaml:"channelName"`
	MspId         string `yaml:"mspId"`
	CryptoPath    string `yaml:"cryptoPath"`
	CertPath      string `yaml:"certPath"`
	KeyPath       string `yaml:"keyPath"`
	TlsCertPath   string `yaml:"tlsCertPath"`
	PeerEndpoint  string `yaml:"peerEndpoint"`
	GatewayPeer   string `yaml:"gatewayPeer"`
}

type Server struct {
	Port           int      `yaml:"port"`
	AllowedHeaders []string `yaml:"allowedHeaders"`
	AllowedMethods []string `yaml:"allowedMethods"`
	AllowedOrigins []string `yaml:"allowedOrigins"`
	ExposedHeaders []string `yaml:"exposedHeaders"`
	JwtKey         string   `yaml:"jwtKey"`
}

func LoadConfig() *Config {
	// set the directory path pointing to the config.yaml file
	fp := "/home/azureuser/fabric-gateway-api/"

	viper.SetConfigType("yaml")
	viper.AddConfigPath(fp)
	viper.SetConfigName("config")

	if err := viper.ReadInConfig(); err != nil {
		log.Panic(err)
	}

	config := Config{}
	if err := viper.Unmarshal(&config); err != nil {
		log.Panic(err)
	}

	return &config
}
