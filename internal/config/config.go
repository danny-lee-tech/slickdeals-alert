package config

type EmailConfig struct {
	SMTP               string `yaml:"smtp"`
	Port               *int   `yaml:"port"`
	SourceEmailAddress string `yaml:"source_email"`
	TargetEmailAddress string `yaml:"target_email"`
	Subject            string `yaml:"subject"`
	PasswordFile       string `yaml:"password_file"`
}

type PushBulletConfig struct {
	APIKey string `yaml:"api_key"`
	Tag    string `yaml:"tag"`
}

type Config struct {
	VoteFilter        *int              `yaml:"vote_filter"`     // Search Filter on minimum number of votes. Used to determine the URL to scrape, specifically the vote query parameter
	NotifyMinimumRank *int              `yaml:"notify_min_rank"` // the minimum number of thumbs up x 2 before a notification occurs
	Email             *EmailConfig      `yaml:"email,omitempty"`
	PushBullet        *PushBulletConfig `yaml:"pushbullet"`
}
