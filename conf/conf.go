package conf

import (
	"io/ioutil"
	"time"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

type HTTP struct {
	GinMode        string        `yaml:"gin_mode"`
	Addr           string        `yaml:"addr"`
	ReadTimeout    time.Duration `yaml:"read_timeout"`
	WriteTimeout   time.Duration `yaml:"wirte_timeout"`
	MaxHeaderBytes int           `yaml:"max_header_bytes"`
}

type Config struct {
	HTTP HTTP `yaml:"http"`
}

var GlobalConfig Config

func InitConfig(path string) {
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		panic(errors.WithMessage(err, " [PANIC] read config file failed"))
	}

	err = yaml.Unmarshal(bytes, &GlobalConfig)
	if err != nil {
		panic(errors.WithMessage(err, " [PANIC] unmarshal yaml failed"))
	}
}
