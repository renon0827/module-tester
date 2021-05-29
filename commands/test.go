package commands

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"module-tester/scenarios"
	"os"

	log "github.com/sirupsen/logrus"
)

// Test : テストコマンド
type Test struct {
	Flags        *flag.FlagSet
	ScenarioFile string
	Out          string
	ScenarioName string
}

func (g Test) showCommandUsage() {
	fmt.Printf("\n使い方: %s\n", g.CommandName())
	fmt.Printf("\n説明: %s\n", g.CommandDescription())
	fmt.Printf("\n実行オプション:\n")
	g.Flags.PrintDefaults()
}

// CommandName : コマンド名を返す
func (t Test) CommandName() string {
	return "Test"
}

// CommandDescription : コマンド説明を返す
func (t Test) CommandDescription() string {
	return "予め作成されたシナリオファイルを元に、テストを実行します。"
}

// Init : コマンドの初期化
func (t *Test) Initialize(args []string) error {
	t.Flags = flag.NewFlagSet(t.CommandName(), flag.ExitOnError)
	t.Flags.Usage = t.showCommandUsage
	t.Flags.StringVar(&t.ScenarioFile, "scenario", "", "テストを実行するシナリオファイルを指定します（必須）")
	t.Flags.StringVar(&t.Out, "out", "", "テストの結果を出力するファイル名を指定します（指定なしの場合、結果ファイルを出力しません）")
	t.Flags.StringVar(&t.ScenarioName, "name", "", "実行するシナリオ名を指定します（指定なしまたは空の場合、すべてのシナリオを実行します）")
	t.Flags.Parse(args)
	return nil
}

// Run : コマンドの実行
func (t Test) Run() int {
	scenarioDataSource, err := ioutil.ReadFile(t.ScenarioFile)
	if err != nil {
		log.Errorln(err)
		return 1
	}

	var scenarioData scenarios.ScenarioTest
	err = json.Unmarshal(scenarioDataSource, &scenarioData)
	if err != nil {
		log.WithError(err).Errorln("シナリオファイルが不正です")
		return 1
	}

	if t.ScenarioName != "" {
		exist := false
		for _, item := range scenarioData.Scenarios {
			if item.Name != t.ScenarioName {
				continue
			}
			scenarioData.Scenarios = []scenarios.Scenario{item}
			exist = true
			break
		}
		if !exist {
			log.WithFields(log.Fields{
				"シナリオ名": t.ScenarioName,
			}).Errorln("シナリオ名が不正です")
		}
	}

	results, correctFlag := scenarios.RunScenarioTest(scenarioData)

	if t.Out != "" {
		resultContent, err := json.Marshal(results)
		if err != nil {
			log.WithError(err).Errorln("テスト結果が不正です")
			return 1
		}

		var buf bytes.Buffer
		err = json.Indent(&buf, resultContent, "", "  ")
		if err != nil {
			log.Errorln(err)
		}

		err = ioutil.WriteFile(t.Out, buf.Bytes(), os.ModePerm)
		if err != nil {
			log.Errorln(err)
			return 1
		}

		log.WithField("出力ファイル", t.Out).Infoln("テスト結果ファイルを出力しました")
	}

	if correctFlag {
		log.Infoln("すべてのテストが正常に終了しました")
		return 0
	}

	log.Errorln("失敗したテストが存在します")
	return 1
}
