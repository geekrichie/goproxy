package common

import (
	"gopkg.in/ini.v1"
	"os"
	"runtime"
	"strings"
)

func GetConfPath() string{
	_, fullFilename, _, _ := runtime.Caller(0)
	lastIndex := strings.LastIndexByte(fullFilename, '/')
	if runtime.GOOS == "windows" {
		fullFilename = strings.ReplaceAll(fullFilename[:lastIndex],"/","\\")
	}
	return strings.Join([]string{fullFilename,"..", "conf"}, string(os.PathSeparator))
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

func GetServerPort() int{
	port ,_ := serverConf.Section("server").Key("server_port").Int()
	return port
}


func GetSecretKey() string{
	return serverConf.Section("connect").Key("secretkey").String()
}

