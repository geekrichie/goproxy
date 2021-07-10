package file

import (
	"encoding/json"
	"goproxy/common"
	"os"
	"strings"
)

type Task struct {
	Port int   //server监听端口号
	TargetAddrs []string //目标地址
}

const (
	TASK_FILE_NAME= "tasks.json"
)

type Db struct {
	Tasks []Task
}

func (d *Db) AddTask(task Task){
	d.Tasks = append(d.Tasks, task)
}

func (t *Task) SaveToFile() error{
	path := common.GetConfPath()
	filename := strings.Join([]string{path, TASK_FILE_NAME}, string(os.PathSeparator))
	fhandle, err := os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	serializedTask, err := json.Marshal(t)
	if err != nil {
		return err
	}
	tmpData := make([]byte, 0,100)
	tmpData = append(tmpData, []byte("#**#\n")...)
	tmpData = append(tmpData,  serializedTask...)
	tmpData = append(tmpData, []byte("\n")...)

	_, err = fhandle.Write(tmpData)
	if err != nil {
		return err
	}
	err  = fhandle.Sync()
	if err != nil {
		return err
	}
	fhandle.Close()
	return err
}