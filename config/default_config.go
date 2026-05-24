package config

// DefaultConfig returns a new Config with default settings.
func DefaultConfig() *Config {
	return &Config{
		CronExpression:        "0 20,40,55 * * * ?",
		StandUpCronExpression: "0 55 * * * ?",
		Language:              nil,
	}
}
