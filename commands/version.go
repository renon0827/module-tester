package commands

import (
	"fmt"
	"module-tester/constants"
)

// ShowVersion : バージョン情報を表示する
func ShowVersion() {
	fmt.Printf("version %s\n", constants.Version)
}
