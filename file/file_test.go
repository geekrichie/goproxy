package file

import (
	"testing"
)

func TestTask_SaveToFile(t *testing.T) {
	task := new(Task)
	task.Port = 9999
	task.TargetAddrs = []string{"127.0.0.1:3306"}
	task.SaveToFile()
	(&Task{}).SaveToFile()

}