package main

import (
	"fmt"
	"log"
	"path/filepath"

	homedir "github.com/mitchellh/go-homedir"
)

var dbName = "tasks.db"

func init() {
	homeDir, _ := homedir.Dir()
	dbPath := filepath.Join(homeDir, dbName)
	err := Init(dbPath)
	if err != nil {
		log.Fatalf("Could not initialize database at %s: %v", dbPath, err)
	}
}

func main() {
	fmt.Printf("test\n")
	AddTask("wyrzuuic smieci")
	AddTask("ootworzyc prezenty")

	CompleteTask(5)

	tasks, err := ListTasks()
	if err != nil {
		log.Fatalf("could not list tasks: %v", err)
	}

	for i, task := range tasks {
		fmt.Printf("%d: %s\n", i+1, task.Value)
	}
}
