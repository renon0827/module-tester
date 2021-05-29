package scenarios

import (
	"module-tester/drivers"
	"module-tester/repository"
	"module-tester/stubs"

	log "github.com/sirupsen/logrus"
)

func getTestResultFromDriver(result drivers.DriveResult) TestResult {
	var detail *TestResultDetail
	evaluationResult := result.Evaluator.Evaluation()
	if result.Response != nil {
		detail = &TestResultDetail{
			Begin:    result.Begin,
			End:      result.End,
			Request:  result.Request,
			Response: result.Response,
			Score:    float64(result.End.Sub(result.Begin).Microseconds()) / 1000,
		}
	}
	return TestResult{
		Correct: evaluationResult.Correct,
		Message: evaluationResult.Message,
		Detail:  detail,
	}
}

func getTestResultFromStub(result stubs.ProcessResult) TestResult {
	var detail *TestResultDetail
	evaluationResult := result.Evaluator.Evaluation()
	if result.Request != nil {
		detail = &TestResultDetail{
			Begin:    result.Begin,
			End:      result.End,
			Request:  result.Request,
			Response: result.Response,
			Score:    float64(result.End.Sub(result.Begin).Microseconds()) / 1000,
		}
	}
	return TestResult{
		Correct: evaluationResult.Correct,
		Message: evaluationResult.Message,
		Detail:  detail,
	}
}

func RunScenarioTest(scenario ScenarioTest) ([]ScenarioTestResult, bool) {
	correctFlag := true
	results := []ScenarioTestResult{}
	for _, scenario := range scenario.Scenarios {
		log.WithFields(log.Fields{
			"シナリオ名": scenario.Name,
		}).Infoln("テスト開始")

		repository := repository.CreateStandardRepository()
		moduleTestResult := make([]ModuleTestResult, len(scenario.Sequence))
		for sequenceNo, moduleTest := range scenario.Sequence {
			var stubsChan chan stubs.ProcessResult
			if moduleTest.Stubs != nil {
				stubsChan = make(chan stubs.ProcessResult, len(*moduleTest.Stubs))
				for stubNo, stubOption := range *moduleTest.Stubs {
					l := log.WithFields(log.Fields{
						"スタブ番号":    stubNo,
						"シークエンス番号": sequenceNo,
					})
					if stubOption.TimeoutSeconds == 0 {
						l.Warnln("タイムアウトが設定されていません")
					}
					stub, err := stubs.CreateStub(stubOption, repository)
					if err != nil {
						log.WithError(err).Errorln("Stubの生成に失敗しました")
						return results, correctFlag
					}
					go func(stub stubs.Stub, log *log.Entry) {
						log.Infoln("スタブを開始します")
						stubsChan <- stub.Listen()
						log.Infoln("スタブを終了します")
					}(stub, l)
				}
			}

			var driverResult *TestResult
			if moduleTest.Driver != nil {
				log := log.WithFields(log.Fields{
					"シークエンス番号": sequenceNo,
				})
				driver, err := drivers.CreateDriver(*moduleTest.Driver, repository)
				if err != nil {
					log.WithError(err).Errorln("Driverの生成に失敗しました")
					return results, correctFlag
				}
				log.Infoln("ドライバーを開始します")
				d := getTestResultFromDriver(driver.Drive())
				log.Infoln("ドライバーを終了します")
				driverResult = &d
				if !d.Correct {
					log.Errorln(d.Message)
					correctFlag = false
				} else {
					log.Infoln(d.Message)
				}
			}

			var stubResult *[]TestResult
			if moduleTest.Stubs != nil {
				s := make([]TestResult, len(*moduleTest.Stubs))
				for i := range *moduleTest.Stubs {
					log := log.WithFields(log.Fields{
						"スタブ番号":    i,
						"シークエンス番号": sequenceNo,
					})
					s[i] = getTestResultFromStub(<-stubsChan)
					if !s[i].Correct {
						log.Errorln(s[i].Message)
						correctFlag = false
					} else {
						log.Infoln(s[i].Message)
					}
				}
				stubResult = &s
			}

			moduleTestResult[sequenceNo] = ModuleTestResult{
				Driver: driverResult,
				Stubs:  stubResult,
			}
		}

		results = append(results, ScenarioTestResult{
			ScenarioName: scenario.Name,
			Results:      moduleTestResult,
		})
	}
	return results, correctFlag
}
