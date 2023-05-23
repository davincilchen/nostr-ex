package config

import (
	"encoding/json"
	"io/ioutil"
	"nostr-ex/pkg/db"

	"github.com/caarlos0/env"
	_ "github.com/joho/godotenv/autoload" //support .env && autoload
)

type Config struct {
	Server Server    `json:"Server"`
	Nostr  Nostr     `json:"Nostr"`
	Relay  Relay     `json:"Relay"`
	DB     db.Config `json:"DB"`
}

type Server struct {
	Port string `json:"Port" env:"SERVER_PORT" envDefault:"8000"`
}

type Nostr struct {
	PublicKey  string `json:"PublicKey" env:"PUBLIC_KEY" envDefault:"6a11da0c34b6881f61ab5116a52c72e161af9cdca5ca8fd59f296c3d94e8532a"`
	PrivateKey string `json:"PrivateKey" env:"PRIVATE_KEY" envDefault:"073af9c587de8e9ac5afcfc3eea195f21aa9940daa09614c5dd59260c1a77812"`
}

type Relay struct {
	//URL string `json:"URL" env:"RELAY_URL" envDefault:"ws://127.0.0.1:8100/"`
	URL string `json:"URL" env:"RELAY_URL" envDefault:"wss://relay.nekolicio.us/"`
}

var config *Config

func GetConfig() *Config {
	return config
}

func GetRelayUrl() string {
	return config.Relay.URL
}

func New(path string) (*Config, error) {

	cfg, err := new(path)

	if cfg != nil {
		config = cfg
	}

	return cfg, err

}

func new(path string) (*Config, error) {

	cfg, err := NewFromFile(path)
	if err == nil { //priority first
		return cfg, nil
	}

	cfg, _ = NewFromEnv()

	return cfg, nil
}

func NewFromFile(path string) (*Config, error) {
	buf, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var config Config
	err = json.Unmarshal(buf, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

func NewFromEnv() (*Config, error) {

	var config Config

	config.Server = *GetServerConfig()
	config.Nostr = *GetNostrConfig()
	config.Relay = *GetRelayConfig()
	config.DB = *GetDBConfig()
	return &config, nil
}

func GetServerConfig() *Server {
	cfg := &Server{}
	env.Parse(cfg)
	return cfg
}

func GetNostrConfig() *Nostr {
	cfg := &Nostr{}
	env.Parse(cfg)
	return cfg
}

func GetRelayConfig() *Relay {
	cfg := &Relay{}
	env.Parse(cfg)
	return cfg
}

func GetDBConfig() *db.Config {
	cfg := &db.Config{}
	env.Parse(cfg)
	return cfg
}
