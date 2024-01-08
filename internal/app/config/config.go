package config

import (
	"flag"
	"os"
)

type ServerConfig struct {
	startAddress string
	shortBaseURL string
	filePath     string
}

func (sa *ServerConfig) GetStartAddress() string {
	return sa.startAddress
}

func (sa *ServerConfig) GetShortBaseURL() string {
	return sa.shortBaseURL
}

func (sa *ServerConfig) GetFilePath() string {
	return sa.filePath
}

func (sa *ServerConfig) SetStartAddress(value string) {
	sa.startAddress = value
}

func (sa *ServerConfig) SetShortBaseURL(value string) {
	sa.shortBaseURL = value
}

func (sa *ServerConfig) SetFilePath(value string) {
	sa.filePath = value
}

func (sa *ServerConfig) ParseFlags() {
	flag.StringVar(&sa.startAddress, "a", "localhost:8080", "address and port to run shortener")
	flag.StringVar(&sa.shortBaseURL, "b", "http://localhost:8080", "address and port for base short URL")
	flag.StringVar(&sa.filePath, "f", "/tmp/short-url-db.json", "storage file path")

	flag.Parse()

	if envStartAddress := os.Getenv("SERVER_ADDRESS"); envStartAddress != "" {
		sa.SetStartAddress(envStartAddress)
	}

	if envShortBaseURL := os.Getenv("BASE_URL"); envShortBaseURL != "" {
		sa.SetShortBaseURL(envShortBaseURL)
	}

	if envFilePath := os.Getenv("FILE_STORAGE_PATH"); envFilePath != "" {
		sa.SetFilePath(envFilePath)
	}
}
