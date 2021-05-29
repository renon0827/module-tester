package repository

// Repository : 変数を一元的に管理するオブジェクト
type Repository interface {
	Get(key string) (interface{}, error)
	Set(key string, val interface{})
	Clear()
}

type repositoryInstance map[string]interface{}
