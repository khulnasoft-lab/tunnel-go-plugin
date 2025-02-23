package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/khulnasoft-lab/tunnel-go-plugin/pkg/common"
)

var (
	tempJsonFileName = "tunnel-go-plugin-temp.json"
	version          = "dev"
)

func main() {
	if common.IsHelp() {
		helpMessage()
		return
	}

	pluginArgs, tunnelCmd := common.RetrievePluginArguments([]string{"--plugin-output", "--output"})

	pluginOutput := pluginArgs["--plugin-output"]
	if pluginOutput == "" {
		log.Fatal("flag --plugin-output is required")
	}

	tunnelOutputFileName := pluginArgs["--output"]
	if tunnelOutputFileName == "" {
		tunnelOutputFileName = filepath.Join(os.TempDir(), tempJsonFileName)
		defer removeFile(tunnelOutputFileName)
	}

	if err := common.MakeTunnelJsonReport(tunnelCmd, tunnelOutputFileName); err != nil {
		log.Fatalf("failed to make tunnel report: %v", err)
	}
	_, err := common.ReadReport(tunnelOutputFileName)
	if err != nil {
		log.Fatalf("failed to get report from json: %v", err)
	}

	if err := saveResult(pluginOutput, []byte{}); err != nil {
		log.Fatalf("failed to save result: %v", err)
	}
}

func removeFile(file string) {
	if err := os.Remove(file); err != nil {
		log.Fatalf("failed to remove file %v", err)
	}
}

func closeFile(file *os.File) {
	if err := file.Close(); err != nil {
		log.Fatalf("failed to remove file %v", err)
	}
}

func saveResult(filename string, result []byte) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	_, err = file.Write(result)
	if err != nil {
		return err
	}
	defer closeFile(file)
	return nil
}

func helpMessage() {
	_, err := fmt.Printf(`
tunnel-go-plugin v%s
Usage: tunnel tunnel-go-plugin [-h,--help] command target filename
 A Tunnel common plugin.
Options:
  -h, --help    Show usage.
Examples:
  # example
  tunnel tunnel-go-plugin
`, version)
	if err != nil {
		log.Fatalf("Failed to display help message %v", err)
	}
	os.Exit(0)
}
