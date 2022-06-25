package storageconsul

import (
	"time"
)

// StorageData describes the data that is saved to KV
type StorageData struct {
	Value    []byte    `json:"value"`
	Modified time.Time `json:"modified"`
}
