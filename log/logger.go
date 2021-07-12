package log

import log "github.com/amoghe/distillog"

func Info(msg string) {
     log.Infoln(msg)
}

func Infof(msg string, v ...interface{}) {
	log.Infof(msg, v...)
}

func Debug(msg string) {
	log.Debugln(msg)
}

func Debugf(msg string, v ...interface{}) {
	log.Debugf(msg, v...)
}

func Warning(msg string) {
	log.Warningln(msg)
}

func Warningf(msg string, v ...interface{}) {
	log.Debugf(msg, v...)
}

func Error(msg string) {
	log.Errorln(msg)
}

func Errorf(msg string, v ...interface{}) {
	log.Errorf(msg, v...)
}

