package agent

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/nastts/final-calculator/structs"
)


func Worker(){
	for{
		pendingTask, err := GetTask()
		if err != nil{
			time.Sleep(1 * time.Second)
			continue
		}
		if pendingTask == nil{
			time.Sleep(1 * time.Second)
			continue
		}
		result := CalcTask(pendingTask)
		if err := SendResult(pendingTask.ID, result); err != nil{
			log.Println("Ошибка отправки результата:", err)
		}
	}
}




func GetTask()(*structs.Task, error){
	resp, err := http.Get("http://localhost/internal/task")
	if err != nil{
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK{
		return nil, nil
	}
	var task structs.Task
	err = json.NewDecoder(resp.Body).Decode(&task)
	if err != nil{
		return nil, err
	}
	return &task, nil
}

func SendResult(id string, result float64)error{
	data := map[string]interface{}{
		"id": id,
		"result": result,
	}
	jsondata,_ := json.Marshal(data)
	resp, err := http.Post("http://localhost/internal/task","application/json", bytes.NewBuffer(jsondata))
	if err != nil{
		return err
	}
	defer resp.Body.Close()
	return nil
}


func CalcTask(task *structs.Task)(float64){
	switch task.Operation{
	case "+":
		return task.Arg1 + task.Arg2
	case "-":
		return task.Arg1 - task.Arg2
	case "*":
		return task.Arg1 * task.Arg2
	case "/":
		if task.Arg2 == 0{
			return 0
		}
		return task.Arg1 / task.Arg2
	default:
		return 0
	}
}