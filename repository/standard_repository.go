package repository

import (
	"errors"

	log "github.com/sirupsen/logrus"
)

// standardRepository : Mapによって値を管理するオブジェクト
type standardRepository struct {
	Instance repositoryInstance
}

func (r standardRepository) Get(key string) (interface{}, error) {
	if val, exist := r.Instance[key]; exist {
		return val, nil
	} else {
		return nil, errors.New("変数 '" + key + "' は存在しません")
	}
}

func (r standardRepository) Set(key string, val interface{}) {
	log := log.WithFields(log.Fields{
		"変数名": key,
		"値":   val,
	})
	if v, exist := r.Instance[key]; exist {
		log.WithField("変更前の値", v).Infoln("変数の値を更新しました")
		r.Instance[key] = val
	} else {
		log.Infoln("新しい変数を作成しました")
		r.Instance[key] = val
	}
}

func (r standardRepository) Clear() {
	r.Instance = repositoryInstance{}
}

// CreateStandardRepository : Mapによって値を管理するオブジェクトを生成する
func CreateStandardRepository() Repository {
	return standardRepository{
		Instance: repositoryInstance{},
	}
}
