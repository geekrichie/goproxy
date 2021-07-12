package common

import (
	"fmt"
	log "github.com/amoghe/distillog"
	"testing"
)

func TestGetConfPath(t *testing.T) {
	path := GetConfPath()
	fmt.Println(path)
}

func TestGetServerPort(t *testing.T) {
	port := GetServerPort()
	fmt.Println(port)
}

func TestLog(t *testing.T) {

	// ... later ...

	log.Infoln("Starting program")
	log.Debugln("initializing the frobnicator")
	log.Warningln("frobnicator failure detected, proceeding anyways...")
	log.Infoln("Exiting")
}
