package logs

import (
	"fmt"
	"github.com/op/go-logging"
	"go-common/utils"
	"io"
	"os"
	"time"
)

var (
	// LoggingPattern is the pattern to use for rendering the logs
	loggingPattern = ` %{time:15:04:05.000} [%{module}] %{color}▶ %{level:.10s}%{color:reset} %{message}`
)

type LogComponent interface {
	LogPrefix() string
	SetLogger(logger Logger, manager LogManager)
}

type Logger interface {
	Debug(v ...interface{})
	Info(v ...interface{})
	Warning(v ...interface{})
	Error(v ...interface{})
	Critical(v ...interface{})
	Fatal(v ...interface{})
}

func NewLoggerWithConfig(stdout, syslog bool, level, prefix string, ws ...io.Writer) (Logger, error) {
	return NewLogger(map[string]interface{}{
		"stdout": stdout,
		"syslog": syslog,
		"level":  level,
		"prefix": prefix,
	}, ws...)
}

func NewLogger(cfg map[string]interface{}, ws ...io.Writer) (Logger, error) {
	logConfig := ConfigGetter(cfg)
	moduleName := logConfig.Prefix
	if utils.EmptyOrWhiteSpace(moduleName) {
		moduleName = "DEFAULT"
	}
	logger := logging.MustGetLogger(moduleName)

	backend := make([]logging.Backend, 0, 3)
	var b logging.Backend
	if logConfig.StdOut {
		b = logging.NewLogBackend(os.Stdout, logConfig.Prefix, 0)
		backend = append(backend, b)
	}

	for _, w := range ws {
		b = logging.NewLogBackend(w, logConfig.Prefix, 0)
		backend = append(backend, b)
	}

	if logConfig.Syslog {
		var err error
		b, err = logging.NewSyslogBackend(logConfig.Prefix)
		if err != nil {
			return nil, err
		}
		backend = append(backend, b)
	}

	for i, b := range backend {
		format := logging.MustStringFormatter(loggingPattern)
		backendLeveled := logging.AddModuleLevel(logging.NewBackendFormatter(b, format))
		logLevel, err := logging.LogLevel(logConfig.Level)
		if err != nil {
			return nil, err
		}
		backendLeveled.SetLevel(logLevel, moduleName)
		backend[i] = backendLeveled
	}

	logging.SetBackend(backend...)
	return internalLogger{logger}, nil
}

func ConfigGetter(logConfig map[string]interface{}) config {
	cfg := config{}
	if v, ok := logConfig["stdout"]; ok {
		cfg.StdOut = v.(bool)
	}
	if v, ok := logConfig["syslog"]; ok {
		cfg.Syslog = v.(bool)
	}
	if v, ok := logConfig["level"]; ok {
		cfg.Level = v.(string)
	}
	if v, ok := logConfig["prefix"]; ok {
		cfg.Prefix = v.(string)
	}
	return cfg
}

type config struct {
	Level  string
	StdOut bool
	Syslog bool
	Prefix string
}

type internalLogger struct {
	logger *logging.Logger
}

// Debug implements the internalLogger interface
func (l internalLogger) Debug(v ...interface{}) {
	l.logger.Debug(v...)
}

// Info implements the internalLogger interface
func (l internalLogger) Info(v ...interface{}) {
	l.logger.Info(v...)
}

// Warning implements the internalLogger interface
func (l internalLogger) Warning(v ...interface{}) {
	l.logger.Warning(v...)
}

// Error implements the internalLogger interface
func (l internalLogger) Error(v ...interface{}) {
	l.logger.Error(v...)
}

// Critical implements the internalLogger interface
func (l internalLogger) Critical(v ...interface{}) {
	l.logger.Critical(v...)
}

// Fatal implements the internalLogger interface
func (l internalLogger) Fatal(v ...interface{}) {
	l.logger.Fatal(v...)
}

func CreateDefaultLogger(prefix string) Logger {
	return &defaultLogger{prefix: prefix}
}

type defaultLogger struct {
	prefix string
}

func (d defaultLogger) print(level string, v []interface{}) {
	prefix := fmt.Sprintf("%s [%s]\t", utils.ToDatetimeString(time.Now()), d.prefix)
	fmt.Println(append(append([]interface{}{prefix, level, " "}), v...)...)
}

func (d defaultLogger) Debug(v ...interface{}) {
	d.print("DEBUG", v)
}

func (d defaultLogger) Info(v ...interface{}) {
	d.print("INFO", v)
}

func (d defaultLogger) Warning(v ...interface{}) {
	d.print("WARNING", v)
}

func (d defaultLogger) Error(v ...interface{}) {
	d.print("ERROR", v)
}

func (d defaultLogger) Critical(v ...interface{}) {
	d.print("CRITICAL", v)
}

func (d defaultLogger) Fatal(v ...interface{}) {
	d.print("FATAL", v)
}
