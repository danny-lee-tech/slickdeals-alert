package main

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/danny-lee-tech/slickdeals-alert/internal/config"
	"github.com/danny-lee-tech/slickdeals-alert/internal/emailer"
	"github.com/danny-lee-tech/slickdeals-alert/internal/pushbulleter"
	"github.com/danny-lee-tech/slickdeals-alert/internal/scraper"
	"gopkg.in/yaml.v2"
)

var DefaultConfigLocation = "configs/config.yml"

func main() {
	cfg, err := getConfig()
	if err != nil {
		log.Fatal("Error:", err)
		return
	}

	fmt.Println("Starting scraper")

	scraper := scraper.Scraper{
		VoteFilter:        *cfg.VoteFilter,
		NotifyMinimumRank: *cfg.NotifyMinimumRank,
	}

	if cfg.Email != nil {
		scraper.Emailer = &emailer.Emailer{
			SMTP:               cfg.Email.SMTP,
			Port:               *cfg.Email.Port,
			SourceEmailAddress: cfg.Email.SourceEmailAddress,
			TargetEmailAddress: cfg.Email.TargetEmailAddress,
			Subject:            cfg.Email.Subject,
			PasswordFile:       cfg.Email.PasswordFile,
		}
	}

	if cfg.PushBullet != nil {
		scraper.PushBulleter = &pushbulleter.PushBulleter{
			APIKey: cfg.PushBullet.APIKey,
			Tag:    cfg.PushBullet.Tag,
		}
	}

	err = scraper.Execute()
	if err != nil {
		log.Fatal("Error:", err)
		return
	}
}

func getConfig() (config.Config, error) {
	configLocation := getConfigLocation()
	fmt.Println("Setting up configs", configLocation)
	cfgBytes, err := os.ReadFile(configLocation)
	if err != nil {
		return config.Config{}, err
	}

	fmt.Println(string(cfgBytes))

	var cfg config.Config
	err = yaml.Unmarshal(cfgBytes, &cfg)
	if err != nil {
		return config.Config{}, err
	}

	err = validateConfig(&cfg)
	if err != nil {
		return config.Config{}, err
	}

	return cfg, nil
}

func validateConfig(cfg *config.Config) error {
	if cfg.NotifyMinimumRank == nil {
		return errors.New("missing required field: notify_min_rank")
	}

	if cfg.VoteFilter == nil {
		return errors.New("missing required field: vote_filter")
	}

	if cfg.Email != nil {
		if cfg.Email.SMTP == "" {
			return errors.New("missing required field: email.smtp")
		}

		if cfg.Email.Port == nil {
			return errors.New("missing required field: email.port")
		}

		if cfg.Email.SourceEmailAddress == "" {
			return errors.New("missing required field: email.source_email")
		}

		if cfg.Email.TargetEmailAddress == "" {
			return errors.New("missing required field: email.target_email")
		}

		if cfg.Email.Subject == "" {
			return errors.New("missing required field: email.subject")
		}

		if cfg.Email.PasswordFile == "" {
			return errors.New("missing required field: email.password_file")
		}
	}

	return nil
}

func getConfigLocation() string {
	if len(os.Args) > 1 {
		return os.Args[1]
	}

	return DefaultConfigLocation
}
