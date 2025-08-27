package config

import (
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/go-jose/go-jose/v4"
	"gopkg.in/yaml.v3"
)

const (
	defaultServicePort uint16 = 5000
	defaultMaxConns    int32  = 10
	defaultLogFile     string = "/var/log/saintjacques/logs"
	defaultImageSize   int64  = 5
	defaultUploadDir   string = "download/noco/Patrimoine Jacquaire/TMedias/CheminMedia"
)

type (
	Config struct {
		Port       uint16     `yaml:"-"`
		Logs       Logs       `yaml:"logs"`
		Database   *Database  `yaml:"database"`
		FileSystem FileSystem `yaml:"filesystem"`
		SMTP       *SMTP      `yaml:"smtp"`
		Session    Session    `yaml:"session"`
	}

	Logs struct {
		Level string `yaml:"level"`
		File  string `yaml:"file"`
	}

	Database struct {
		DSN      string `yaml:"dsn"`
		MaxConns int32  `yaml:"maxConns"`
	}

	FileSystem struct {
		UploadDir string `yaml:"uploadDir"`
		MaxSize   int64  `yaml:"maxSize"`
	}

	SMTP struct {
		Host   string `yaml:"host"`
		APIKey string `yaml:"apikey"`
		From   From   `yaml:"from"`
	}

	From struct {
		Name  string `yaml:"name"`
		Email string `yaml:"email"`
	}

	Session struct {
		Duration     time.Duration `yaml:"duration"`
		JWTSecret    []byte        `yaml:"-"`
		JWTSecretRaw string        `yaml:"jwtSecret"`
		Signer       jose.Signer   `yaml:"-"`
	}
)

func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg Config

	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	err = cfg.checkDatabase()
	if err != nil {
		return nil, err
	}

	cfg.checkFileSystem()

	err = cfg.checkSMTP()
	if err != nil {
		return nil, err
	}

	err = cfg.checkSession()
	if err != nil {
		return nil, err
	}

	if cfg.Logs.Level == "" {
		cfg.Logs.Level = "info"
	}

	if cfg.Logs.File == "" {
		cfg.Logs.File = defaultLogFile
	}

	cfg.Port = defaultServicePort

	return &cfg, nil
}

func (cfg *Config) checkDatabase() error {
	if cfg.Database == nil {
		return errors.New("database section is required")
	}

	if cfg.Database.DSN == "" {
		return errors.New("DSN is required")
	}

	if cfg.Database.MaxConns == 0 {
		cfg.Database.MaxConns = defaultMaxConns
	}

	return nil
}

func (cfg *Config) checkFileSystem() {
	if cfg.FileSystem.UploadDir == "" {
		cfg.FileSystem.UploadDir = defaultUploadDir
	}

	if cfg.FileSystem.MaxSize == 0 {
		cfg.FileSystem.MaxSize = defaultImageSize
	}
}

func (cfg *Config) checkSMTP() error {
	if cfg.SMTP == nil {
		return errors.New("smtp section is required")
	}

	if cfg.SMTP.Host == "" {
		return errors.New("smtp host is required")
	}

	if cfg.SMTP.APIKey == "" {
		return errors.New("smtp api key is required")
	}

	if cfg.SMTP.From.Name == "" || cfg.SMTP.From.Email == "" {
		return errors.New("smtp FROM is required")
	}

	return nil
}

func (cfg *Config) checkSession() error {
	if cfg.Session.Duration == 0 {
		cfg.Session.Duration = 30 * 24 * time.Hour
	}

	if cfg.Session.JWTSecretRaw == "" {
		return errors.New("jwt secret is required")
	}

	var err error

	cfg.Session.JWTSecret, err = hex.DecodeString(cfg.Session.JWTSecretRaw)
	if err != nil {
		return err
	}

	signerOpt := jose.SignerOptions{}

	cfg.Session.Signer, err = jose.NewSigner(jose.SigningKey{
		Algorithm: jose.HS256,
		Key:       cfg.Session.JWTSecret,
	}, signerOpt.WithType("JWT"))
	if err != nil {
		panic(fmt.Sprintf(`unable to create signer: %s`, err))
	}

	return nil
}
