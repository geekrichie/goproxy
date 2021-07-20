package file

import (
	"bufio"
	"encoding/json"
	log "github.com/amoghe/distillog"
	"goproxy/common"
	"io"
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

func (t *Task) Save() error{
	path := common.GetConfPath()
	filename := strings.Join([]string{path, TASK_FILE_NAME}, string(os.PathSeparator))
	fhandle, err := os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer fhandle.Close()
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
	return err
}

func LoadTask() *Db{
	var taskdb Db
	path := common.GetConfPath()
	filename := strings.Join([]string{path, TASK_FILE_NAME}, string(os.PathSeparator))
	sr,_ := os.Open(filename)
	defer sr.Close()
	reader := bufio.NewReader(sr)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Errorln(err.Error())
			return nil
		}
		line = strings.TrimSpace(line)
		if line[0] == '#' {
			continue
		}
		var task Task
		json.Unmarshal([]byte(line), &task)
		taskdb.AddTask(task)
	}
	return &taskdb

}