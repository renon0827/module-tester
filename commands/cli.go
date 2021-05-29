package commands

import (
	"flag"
	"fmt"
	"module-tester/commons"
	"module-tester/constants"
	"os"

	"github.com/mattn/go-colorable"
	log "github.com/sirupsen/logrus"
)

const commandDescription = "モジュールテスター"

// Command : コマンドのインターフェース
type Command interface {
	CommandName() string
	CommandDescription() string
	Initialize([]string) error
	Run() int
}

// Options : ModuleTesterのオプション
type Options struct {
	Commands []Command
	Version  bool
}

var options Options

func init() {
	log.SetFormatter(&log.TextFormatter{
		ForceColors:   true,
		FullTimestamp: true,
	})
	log.SetOutput(colorable.NewColorableStdout())

	flag.CommandLine.Init(constants.AppName, flag.ExitOnError)
	flag.CommandLine.Usage = showCommandUsage

	options.Commands = append(options.Commands, &Test{})

	flag.BoolVar(&options.Version, "v", false, "バージョン情報を表示します")
}

func showCommandUsage() {
	fmt.Printf("\n使い方: %s\n", constants.AppName)
	fmt.Printf("\n説明: %s\n", commandDescription)
	fmt.Printf("\n実行オプション:\n")
	flag.PrintDefaults()
	for _, command := range options.Commands {
		fmt.Printf("\nコマンド:\n")
		fmt.Printf("  %s\t%s", command.CommandName(), command.CommandDescription())
	}
	fmt.Printf("\n")
}

// CommandsEntryPoint : このアプリケーションのエントリーポイント
func CommandsEntryPoint() {
	flag.Parse()

	if options.Version {
		ShowVersion()
		os.Exit(0)
	}

	if flag.NArg() > 0 {
		args := flag.Args()
		for _, command := range options.Commands {
			if command.CommandName() != args[0] {
				continue
			}
			if err := command.Initialize(args[1:]); err != nil {
				commons.FailedToPanic(err)
			}
			os.Exit(command.Run())
			break
		}
		fmt.Printf(args[0] + " は存在しないコマンドです\n")
	}
	showCommandUsage()
	os.Exit(1)
}
