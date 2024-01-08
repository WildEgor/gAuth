package configs

import (
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
)

type Configurator struct{}

func NewConfigurator() *Configurator {
	c := &Configurator{}
	c.Load()
	return c
}

func (c *Configurator) Load() {
	err := godotenv.Load(".env", ".env.local")
	if err != nil {
		log.Fatal("Error loading envs file")
	}
}
