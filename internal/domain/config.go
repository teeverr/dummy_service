package domain

import (
	"io/ioutil"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

// local
const configPath string = "/Users/k.sushkov/git/my/dummy_service/cmd/api/k8s/prod/configs/config.yml"

// Inside Docker container
//const configPath string = "/opt/app/config/config.yml"

type Config struct {
	Server struct {
		Port            string        `yaml:"port"`
		Host            string        `yaml:"host"`
		ReadTimeout     time.Duration `yaml:"read_timeout"`
		WriteTimeout    time.Duration `yaml:"write_timeout"`
		GracefulTimeout time.Duration `yaml:"graceful_timeout"`
	} `yaml:"server"`
	Database struct {
		Address  string `yaml:"address"`
		Port     string `yaml:"port"`
		Name     string `yaml:"name"`
		Username string `yaml:"username"`
		Password string `yaml:"password"`
		TimeZone string `yaml:"timezone"`
		Prefix   string `yaml:"prefix"`
	} `yaml:"database"`
	Workload struct {
		Cpu struct {
			Min int `yaml:"min"`
			Max int `yaml:"max"`
		} `yaml:"cpu"`
	} `yaml:"workload"`
}

func (cfg *Config) Parse() error {
	re := regexp.MustCompile(`\$\{.*\}`)
	input, err := ioutil.ReadFile(configPath)
	if err != nil {
		return err
	}
	lines := strings.Split(string(input), "\n")
	for i, line := range lines {
		matches := re.FindStringSubmatch(line)
		if len(matches) == 0 {
			continue
		}
		name := matches[0][2 : len(matches[0])-1]
		value, ok := os.LookupEnv(name)
		if !ok {
			return errors.Errorf("an environment var '%s' not found", name)
		}
		subRe := regexp.MustCompile(regexp.QuoteMeta(matches[0]))
		replacedLine := subRe.ReplaceAllString(line, value)
		lines[i] = replacedLine
	}
	output := strings.Join(lines, "\n")
	err = yaml.Unmarshal([]byte(output), cfg)
	if err != nil {
		return err
	}
	return nil
}
