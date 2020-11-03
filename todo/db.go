package main

import (
	"encoding/binary"
	"errors"
	"fmt"
	"log"

	"github.com/boltdb/bolt"
)

var bucketName = []byte("tasks")
var db *bolt.DB

type Task struct {
	Key   int
	Value string
}

func Init(dbPath string) error {
	var err error
	db, err = bolt.Open(dbPath, 0600, nil)
	if err != nil {
		return err
	}
	err = db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(bucketName)
		return err
	})
	return err
}

func AddTask(task string) error {
	err := db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucketName)
		seq, err := b.NextSequence()
		if err != nil {
			log.Fatalf("couldnt get next sequence: %v", err)
			return err
		}
		return b.Put(itob(int(seq)), []byte(task))
	})
	if err != nil {
		return err
	}
	return nil
}

func ListTasks() ([]Task, error) {
	tasks := make([]Task, 0)
	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucketName)
		cursor := b.Cursor()
		for k, v := cursor.First(); k != nil; k, v = cursor.Next() {
			tasks = append(tasks, Task{btoi(k), string(v)})
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return tasks, nil
}

func deleteTask(key int) error {
	err := db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucketName)
		return b.Delete(itob(key))
	})
	return err
}

func CompleteTask(pos int) error {
	tasks, err := ListTasks()
	if err != nil {
		return errors.New(fmt.Sprintf("could not list tasks: %v", err))
	}

	if pos < 1 || pos > len(tasks) {
		return errors.New(fmt.Sprintf("no such task: %d", pos))
	}

	for i, task := range tasks {
		if i+1 == pos {
			err := deleteTask(task.Key)
			if err != nil {
				return errors.New(fmt.Sprintf("could not delete task: %v", err))
			}
			break
		}
	}
	return nil
}

func itob(v int) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(v))
	return b
}

func btoi(val []byte) int {
	return int(binary.BigEndian.Uint64(val))
}
