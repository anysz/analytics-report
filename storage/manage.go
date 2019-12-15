package storage

type ConfigState struct {
	IsNew bool
}

type LogData struct {
	Id int64 			`json:"id"` // page tracking id
	Items []LogItem		`json:"items"` // tracking logs
	PageContent string  `json:"page_content"` // page content 
}


type LogType int
var (
	LogType_SUCCESS LogType = 0
	LogType_DENIED  LogType = 1

	LogType_UNAVAILABLE LogType = 2
	LogType_TIMEOUT     LogType = 3

	LogType_UNKNOWN     LogType = 999
)


type LogItem struct {
	Id int64 			`json:"id"` // log id
	Timestamp int64 	`json:"timestamp"` // log timestamp
	ItemType  LogType	`json:"type"` // log message type [success, denied, unavailable, timeout, unknown]

	Latitude  float64 	`json:"latitude"`  // log coord latitude
	Longitude float64 	`json:"longitude"` // log coord longitude

}


type StorageDriver interface{
	InitConfig() *ConfigState

	GetAllLogDataIds() ([]int64, bool)

	CreateLogData(LogData) bool
	InputLogData(int64, LogItem) bool
	GetLogData(int64) (*LogData, bool)

	Commit() bool
}

var Driver StorageDriver