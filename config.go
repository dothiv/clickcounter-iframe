package clickcounteriframe

import (
	"fmt"

	"code.google.com/p/gcfg"
)

const VERSION = "v1.10.8"

type Config struct {
	Server struct {
		Port int
	}
	Database struct {
		Host     string
		Name     string
		User     string
		Password string
		Sslmode  string
	}
}

func (c *Config) DSN() (dsn string) {
	dsn = fmt.Sprintf("user=%s dbname=%s sslmode=%s", c.Database.User, c.Database.Name, c.Database.Sslmode)
	if len(c.Database.Host) > 0 {
		dsn = dsn + " host=" + c.Database.Host
	}
	if len(c.Database.Password) > 0 {
		dsn = dsn + " password=" + c.Database.Password
	}
	return
}

func NewDefaultConfig() (c *Config) {
	c = new(Config)
	c.Database.Sslmode = "disable"
	return
}

func NewConfig(fileName string) (c *Config, err error) {
	c = NewDefaultConfig()
	err = gcfg.ReadFileInto(c, fileName)
	return
}
