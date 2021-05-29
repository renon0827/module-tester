package parameters

import "module-tester/repository"

// NamedPipeServerOption : 名前付きパイプサーバーのオプション
type NamedPipeServerOption struct {
	InputBufferSize  int32       `json:"InputBufferSize"`
	OutputBufferSize int32       `json:"OutputBufferSize"`
	Data             interface{} `json:"Data"`
}

// NamedPipeClientOption : 名前付きパイプクライアントのオプション
type NamedPipeClientOption struct {
	TimeoutSeconds int         `json:"TimeoutSeconds"`
	Data           interface{} `json:"Data"`
	BufferSize     int         `json:"BufferSize"`
}

func (o NamedPipeClientOption) DataToString() string {
	return dataToString(o.Data)
}

func (o *NamedPipeClientOption) DataSetFromRepository(repo repository.Repository) {
	o.Data = dataSetFromRepository(o.Data, repo)
}

func (o NamedPipeServerOption) DataToString() string {
	return dataToString(o.Data)
}

func (o *NamedPipeServerOption) DataSetFromRepository(repo repository.Repository) {
	o.Data = dataSetFromRepository(o.Data, repo)
}
