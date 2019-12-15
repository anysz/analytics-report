package storage

import (
	"os"
	"time"
	"strconv"
	"sync"
	"encoding/json"
)

type fileStorageConfig struct{
	items []int64 `json:"items"`
}

type fileStorage struct{
	mx sync.Mutex
	metaFilename	string
	config 			*fileStorageConfig
}
func NewFileStorage(metaFilename string) {
	Driver = &fileStorage{metaFilename: metaFilename, config: new(fileStorageConfig)}
}
func (f *fileStorage) loadFileStorageMeta() bool {
	file, err := os.Open(f.metaFilename)
	defer file.Close()
	if os.IsNotExist(err) {
		f.saveFileStorageMeta()
		return true
	} else if err != nil {
		panic(err)
	} else {
		decoder := json.NewDecoder(file)
		decoder.Decode(f.config)
	}
	return false
}
func (f fileStorage) saveFileStorageMeta() bool {
	file, err := os.Create(f.metaFilename)
	defer file.Close()
	if err != nil && !os.IsExist(err) {
		return false
	}
	encoder := json.NewEncoder(file)
	encoder.SetIndent("","\t")
	encoder.Encode(f.config)
	return true
}
func (f fileStorage) Commit() bool {
	return f.saveFileStorageMeta()
}
func (f *fileStorage) InitConfig() *ConfigState {
	return &ConfigState{
		IsNew: f.loadFileStorageMeta(),
	}
}


func (f *fileStorage) GetAllLogDataIds() ([]int64, bool) {
	return f.config.items, true
}

func (f *fileStorage) CreateLogData(d LogData) bool {
	d.Id = time.Now().UnixNano()
	f.mx.Lock()
	file, err := os.Open("data/" + strconv.Itoa(int(d.Id)) + ".json")
	defer file.Close()
	defer f.mx.Unlock()
	if err != nil && !os.IsExist(err) {
		return false
	}
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "\t")
	encoder.Encode(&d)
	f.config.items = append(f.config.items, d.Id)
	if errx := f.Commit(); !errx {
		return false
	}
	return true
}

func (f *fileStorage) InputLogData(pageId int64, data LogItem) bool {
	filename := "data/" + strconv.Itoa(int(pageId)) + ".json"
	f.mx.Lock()
	file, err := os.Open(filename)
	defer file.Close()
	defer f.mx.Unlock()
	if os.IsNotExist(err) || err != nil {
		return false
	}
	var ld LogData
	decoder := json.NewDecoder(file)
	decoder.Decode(&ld)
	ld.Items = append(ld.Items, data)
	encoder := json.NewEncoder(file)
	encoder.Encode(&ld)
	return true
}

func (f *fileStorage) GetLogData(pageId int64) (*LogData, bool) {
	filename := "data/" + strconv.Itoa(int(pageId)) + ".json"
	f.mx.Lock()
	file, err := os.Open(filename)
	defer file.Close()
	defer f.mx.Unlock()
	if os.IsNotExist(err) || err != nil {
		return nil, false
	}
	var ld LogData
	decoder := json.NewDecoder(file)
	decoder.Decode(&ld)
	return &ld, true
}