package db

import "fmt"

type Config struct {
	Host     string `env:"HOST"`
	Port     int    `env:"PORT"`
	Username string `env:"USERNAME"`
	Password string `env:"PASSWORD"`
	DbName   string `env:"DB_NAME"`
}

func (c *Config) GetConnectionString() string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s",
		c.Host,
		c.Port,
		c.Username,
		c.Password,
		c.DbName,
	)
}
