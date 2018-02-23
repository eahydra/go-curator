package curator

import "log"

var Log Logger = &stdLogger{}

type Logger interface {
	Infof(format string, v ...interface{})
	Infoln(v ...interface{})
	Errorf(format string, v ...interface{})
	Errorln(v ...interface{})
	Debugf(format string, v ...interface{})
	Debugln(v ...interface{})
	Warnf(format string, v ...interface{})
	Warnln(v ...interface{})
}

type stdLogger struct{}

func (*stdLogger) Infof(format string, v ...interface{}) {
	log.Printf(format, v...)
}

func (*stdLogger) Infoln(v ...interface{}) {
	log.Println(v...)
}

func (*stdLogger) Errorf(format string, v ...interface{}) {
	log.Printf(format, v...)
}

func (*stdLogger) Errorln(v ...interface{}) {
	log.Println(v...)
}

func (*stdLogger) Debugf(format string, v ...interface{}) {
	log.Printf(format, v...)
}

func (*stdLogger) Debugln(v ...interface{}) {
	log.Println(v...)
}

func (*stdLogger) Warnf(format string, v ...interface{}) {
	log.Printf(format, v...)
}

func (*stdLogger) Warnln(v ...interface{}) {
	log.Println(v...)
}
