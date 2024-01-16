package main

import (
	"encoding/json"
	"os"
)

type Task struct {
	Question string `json:"Question"`
	Solution string `json:"Solution"`
	Respond  Table
}

func loadTasks(fn string) (map[string]*Task, error) {
	f, err := os.ReadFile(fn)
	if err != nil {
		return nil, err
	}

	var tasks map[string]*Task
	err = json.Unmarshal(f, &tasks)
	if err != nil {
		return nil, err
	}

	for k := range tasks {
		req := &Request{
			Query: tasks[k].Solution,
			Done:  make(chan error),
		}
		querys <- req

		err := <-req.Done
		if err != nil {
			return nil, err
		}

		tasks[k].Respond = req.Respond
	}

	return tasks, nil
}
