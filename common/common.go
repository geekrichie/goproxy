package common

import (
	"log"
	"os"
	"strings"
)

func GetConfPath() string{
	path :=  ".."
	return strings.Join([]string{path, "conf"}, string(os.PathSeparator) )
}


func CheckError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}