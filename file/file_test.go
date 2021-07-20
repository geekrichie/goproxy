package file

import (
	"fmt"
	"testing"
)

func TestTask_SaveToFile(t *testing.T) {
	task := new(Task)
	task.Port = 9999
	task.TargetAddrs = []string{"127.0.0.1:3306"}
	task.Save()
	(&Task{}).Save()

}

func TestLoadTask(t *testing.T) {
	db := LoadTask()
	for _,task := range db.Tasks {
		fmt.Println(task)
	}
}