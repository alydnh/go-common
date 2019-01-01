package logs

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
)

func CreateManager(configPath string) *LogManager {
	logger := CreateDefaultLogger("LOG-MANAGER")
	if manager, err := readConfig(configPath); nil != err {
		logger.Error("Initialize log manager failed with error:", err, "use default logger")
		return &LogManager{
			Stdout:   true,
			Syslog:   false,
			Level:    "DEBUG",
			Prefixes: make(map[string]*logConfig),
			logger:   logger,
		}
	} else {
		manager.logger = logger
		return manager
	}
}

type LogManager struct {
	Stdout   bool                  `yaml:"stdout"`
	Syslog   bool                  `yaml:"syslog"`
	Level    string                `yaml:"level"`
	Prefixes map[string]*logConfig `yaml:"prefixes"`
	logger   Logger
}

type logConfig struct {
	Stdout bool   `yaml:"stdout"`
	Syslog bool   `yaml:"syslog"`
	Level  string `yaml:"level"`
}

func readConfig(configPath string) (*LogManager, error) {
	fileInfo, err := os.Stat(configPath)
	if nil != err {
		return nil, err
	}
	if fileInfo.IsDir() {
		return nil, fmt.Errorf("FILE_NOT_EXISTS")
	}

	var file *os.File
	if file, err = os.OpenFile(configPath, os.O_RDONLY, os.FileMode(0660)); nil != err {
		return nil, err
	}
	defer file.Close()
	if bytes, err := ioutil.ReadAll(file); nil != err {
		return nil, err
	} else {
		manager := &LogManager{}
		if err := yaml.Unmarshal(bytes, manager); nil != err {
			return nil, err
		}

		return manager, nil
	}
}

func (m LogManager) InitializeComponents(components ...LogComponent) {
	for _, component := range components {
		if config, ok := m.Prefixes[component.LogPrefix()]; ok {
			if logger, err := NewLoggerWithConfig(config.Stdout, config.Syslog, config.Level, component.LogPrefix()); nil == err {
				m.logger.Info("Set logger for prefix:", component.LogPrefix())
				component.SetLogger(logger, m)
				continue
			} else {
				m.logger.Error("Initialize logger to component with prefix:", component.LogPrefix(), "failed with error:", err)
			}
		}

		m.logger.Warning("Use default debug level log for component with prefix:", component.LogPrefix())
		logger, _ := NewLoggerWithConfig(true, false, "DEBUG", component.LogPrefix())
		component.SetLogger(logger, m)
	}
}
