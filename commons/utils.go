package commons

import (
	"encoding/json"

	log "github.com/sirupsen/logrus"
)

// FailedToPanic : 渡されたerrorがnil以外だった場合、Panicを発生させる
func FailedToPanic(err error) {
	if err != nil {
		log.Panic(err)
	}
}

func ToStruct(m interface{}, val interface{}) error {
	tmp, err := json.Marshal(m)
	if err != nil {
		return err
	}
	err = json.Unmarshal(tmp, val)
	if err != nil {
		return err
	}
	return nil
}

func StringArrayContains(target string, list []string) bool {
	for _, item := range list {
		if item == target {
			return true
		}
	}
	return false
}
