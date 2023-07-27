package main

import (
	"fmt"
	"os"
	"sync"
	"time"
)

type FileLogWriter struct {
	lock     sync.Mutex
	filename string
	file     *os.File
}

func (w *FileLogWriter) Write(p []byte) (n int, err error) {
	w.lock.Lock()
	defer w.lock.Unlock()
	if err := w.changeFileNameByTime(); err != nil {
		return 0, err
	}
	return w.file.Write(p)
}

func (w *FileLogWriter) changeFileNameByTime() error {
	now := time.Now()
	name := fmt.Sprintf("%s.%04d%02d%02d", gConfig.Log.File, now.Year(), now.Month(), now.Second())
	if name != w.filename {
		if w.file != nil {
			w.file.Close()
			w.file = nil
		}
		tmp, err := os.OpenFile(name, os.O_CREATE|os.O_APPEND, 0644)
		if err != nil {
			return err
		}
		w.file = tmp
		w.filename = name
	}
	return nil
}
