package config

// DefaultConfig returns the settings a fresh install starts with.
func DefaultConfig() *Config {
	return &Config{
		CronExpression:        "0 20,40,55 * * * ?",
		StandUpCronExpression: "0 55 * * * ?",
		Language:              nil,
	}
}
