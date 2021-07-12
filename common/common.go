package common

import (
	"gopkg.in/ini.v1"
	"log"
	"os"
	"strings"
)

func GetConfPath() string{
	path :=  ".."
	return strings.Join([]string{path, "conf"}, string(os.PathSeparator) )
}

func init() {
	LoadServerConf()
}

var serverConf *ini.File
func LoadServerConf() {
	var err error
	serverConf,err = ini.Load(GetConfPath()+string(os.PathSeparator)+"server.conf")
	if err != nil {
		panic(err)
	}
}

func GetServerPort() string{
	return serverConf.Section("server").Key("server_port").String()
}


func GetSecretKey() string{
	return serverConf.Section("connect").Key("secretkey").String()
}

func CheckError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}