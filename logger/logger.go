package logger

import "log"

type Level int

const (
	DEBUG Level = iota + 1
	INFO
	WARN
	SEVERE
)

type Logger struct {
	OutputLevel Level
}

var DefaultLogger *Logger = &Logger{OutputLevel: WARN}

func GetLevel() Level {
	return DefaultLogger.GetLevel()
}

func SetLevel(Level Level) {
	DefaultLogger.SetLevel(Level)
}

func Info(message string) {
	DefaultLogger.Info(message)
}

func Debug(message string) {
	DefaultLogger.Debug(message)
}

func Warn(message string) {
	DefaultLogger.Warn(message)
}

func Severe(message string) {
	DefaultLogger.Severe(message)
}

func Infof(format string, v ...interface{}) {
	DefaultLogger.Infof(format, v...)
}

func Debugf(format string, v ...interface{}) {
	DefaultLogger.Debugf(format, v...)
}
func Warnf(format string, v ...interface{}) {
	DefaultLogger.Warnf(format, v...)
}
func Severef(format string, v ...interface{}) {
	DefaultLogger.Severef(format, v...)
}
func (logger *Logger) Debug(message string) {
	if logger.OutputLevel <= DEBUG {
		log.Println(message)
	}
}

func (logger *Logger) Info(message string) {
	if logger.OutputLevel <= INFO {
		log.Println(message)
	}
}

func (logger *Logger) Warn(message string) {
	if logger.OutputLevel <= WARN {
		log.Println(message)
	}
}

func (logger *Logger) Severe(message string) {
	if logger.OutputLevel <= SEVERE {
		log.Println(message)
	}
}

func (logger *Logger) Debugf(format string, v ...interface{}) {
	if logger.OutputLevel <= DEBUG {
		log.Printf(format, v...)
	}
}

func (logger *Logger) Infof(format string, v ...interface{}) {
	if logger.OutputLevel <= INFO {
		log.Printf(format, v...)
	}
}

func (logger *Logger) Warnf(format string, v ...interface{}) {
	if logger.OutputLevel <= WARN {
		log.Printf(format, v...)
	}
}

func (logger *Logger) Severef(format string, v ...interface{}) {
	if logger.OutputLevel <= SEVERE {
		log.Printf(format, v...)
	}
}

func (logger *Logger) GetLevel() Level {
	return logger.OutputLevel
}

func (logger *Logger) SetLevel(Level Level) {
	if Level >= DEBUG && Level <= SEVERE {
		logger.OutputLevel = Level
	}
}
