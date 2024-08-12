package configuration

import (
	"os"

	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

type Configuration struct {
	TinkoffEndpoint string `yaml:"tinkoff_endpoint"`
}

var instance *Configuration

func (c *Configuration) Load(path string) (*Configuration, error) {
	if instance != nil {
		return instance, nil
	}
	yamlFile, err := os.ReadFile(path)
	if err != nil {
		log.Printf("yamlFile.Get err #%v ", err)
		return c, nil
	}
	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		log.Fatalf("Unmarshalling config: %v", err)
		return nil, err
	}
	instance = c
	return c, nil
}
func (c *Configuration) Get() *Configuration {
	if instance == nil {
		return c
	}
	return instance
}