package db

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"

	log "github.com/sirupsen/logrus"
)

// DB Контракт для базы данных. Чтобы скрыть детали реализации (как именно храним инфу)
type DB struct{}

func (d *DB) getStoragePath(storageName []string) (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", errors.New("Not ok getting info about caller")
	}
	// _, base, _, ok := runtime.Caller(0)
	// if !ok {
	// 	return "", errors.New("Not ok getting info about caller")
	// }
	dir := path.Join(path.Dir(wd), "trade-tech")
	rootDir := filepath.Dir(dir)
	paths := append([]string{rootDir, "storage"}, storageName...)
	p := path.Join(paths...)
	p = path.Clean(p)
	return p, nil
}

// Prune Очистить директорию с данными
func (d *DB) Prune(storageName []string) error {
	storageFile, err := d.getStoragePath(storageName)
	if err != nil {
		log.Errorf("Failed to get storage path %v: %v", storageName, err)
		return err
	}

	if _, err := os.Stat(storageFile); os.IsNotExist(err) {
		log.Warnf("Directory %v doesnt exists", storageFile)
		return nil
	}

	err = os.RemoveAll(storageFile)
	return err
}

// Append Добавить данные в конец стореджа
func (d *DB) Append(storageName []string, content []byte) error {
	storageFile, err := d.getStoragePath(storageName)
	if err != nil {
		log.Errorf("Failed to get storage path %v: %v", storageName, err)
		return err
	}

	if _, err := os.Stat(storageFile); err != nil {
		dir, _ := d.getStoragePath(storageName[:len(storageName)-1])
		os.MkdirAll(dir, 0700)
	}

	file, err := os.OpenFile(storageFile, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0660)
	if err != nil {
		log.Errorf("Failed to open file %v: %v", storageName, err)

		return err
	}
	defer file.Close()

	log.Tracef("Appending to %v", storageFile)
	_, err = file.Write(content)

	return err
}

// Get Получить данные из стореджа
func (d *DB) Get(storageName []string) ([]byte, error) {
	var result []byte
	storageFile, err := d.getStoragePath(storageName)
	if err != nil {
		log.Errorf("Failed to get storage path %v: %v", storageName, err)
		return result, err
	}

	if _, err := os.Stat(storageFile); os.IsNotExist(err) {
		log.Warnf("Directory %v doesnt exists", storageFile)
		return result, err
	}

	file, err := os.OpenFile(storageFile, os.O_RDONLY, 0660)
	if err != nil {
		log.Errorf("Failed to open file %v: %v", storageName, err)
		return result, err
	}
	defer file.Close()

	line := getLastLineWithSeek(file)
	result = []byte(line)

	return result, nil
}

func getLastLineWithSeek(fileHandle *os.File) string {
	line := ""
	var cursor int64 = 0
	stat, _ := fileHandle.Stat()
	filesize := stat.Size()
	for {
		cursor--
		fileHandle.Seek(cursor, io.SeekEnd)

		char := make([]byte, 1)
		fileHandle.Read(char)

		if cursor != -1 && (char[0] == 10 || char[0] == 13) && (len(char) > 0) { // stop if we find not empty line
			break
		}

		line = fmt.Sprintf("%s%s", string(char), line) // there is more efficient way

		if cursor == -filesize { // stop if we are at the begining
			break
		}
	}

	return line
}
