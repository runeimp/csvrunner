package cli

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"os"

	"github.com/runeimp/csvrunner"
	"github.com/runeimp/csvrunner/app"
)

const (
	CLIName    = "csvrun"
	CLILabel   = CLIName + " v" + CLIVersion
	CLIVersion = app.Version
	usage      = `
%s

Usage: %s [OPTIONS] [CSV1 [CSV2] [CSV3]...]

OPTIONS
-------
  -e ENV     Specify the environment variable to define the template string
  -f FILE    Specify the file contents to define the template string
  -h         Display this help information
  -o         Display template run output
  -t STRING  Specify the template string on the command line
  -v         Display the app version

Contents for a CSV can also be piped via stdin


`
)

var (
	csvCount       int
	csvFiles       = []string{}
	outputHelp     bool
	outputVersion  bool
	templateEnv    string
	templateFile   string
	templateOutput bool
	templateString string
)

func help() {
	fmt.Printf(usage, CLILabel, CLIName)
}

func parseArgs(args []string) {
	// fmt.Fprintln(os.Stderr, "cli.parseArgs()")
	argSkip := false
	argsOnly := false
	for i, arg := range args {
		if i == 0 || argSkip {
			argSkip = false
			continue
		}
		switch arg {
		case "-debug", "--debug":
			app.Debug = true
			app.PrintDebug("cli.parseArgs() | arg: %q | app.Debug: %t", arg, app.Debug)
		case "-e", "-env", "-environment", "--environment", "-template-environment", "--template-environment":
			templateEnv = args[i+1]
			templateString = os.Getenv(templateEnv)
			app.PrintDebug("cli.parseArgs() | arg: %q | templateEnv: %q | templateString: %q", arg, templateEnv, templateString)
			argSkip = true
		case "-f", "-file", "-template-file", "--template-file":
			templateFile = args[i+1]
			app.PrintDebug("cli.parseArgs() | arg: %q | templateFile: %q", arg, templateFile)
			argSkip = true
		case "-h", "-help", "--help":
			app.PrintDebug("cli.parseArgs() | arg: %q", arg)
			help()
			os.Exit(0)
		case "-o", "-out", "-template-output", "--template-output":
			templateOutput = true
			app.PrintDebug("cli.parseArgs() | arg: %q | templateOutput: %t", arg, templateOutput)
		// case "-ot":
		// 	templateOutput = true
		// 	templateString = args[i+1]
		// 	app.PrintDebug("cli.parseArgs() | arg: %q | templateOutput: %t | templateString: %q", arg, templateOutput, templateString)
		// 	argSkip = true
		case "-t", "-temp", "-template", "--template":
			templateString = args[i+1]
			app.PrintDebug("cli.parseArgs() | arg: %q | templateString: %q", arg, templateString)
			argSkip = true
		case "-v", "-ver", "-version", "--version":
			fmt.Println(CLILabel)
			os.Exit(0)
		default:
			if arg[0] == '-' {
				if argsOnly {
					app.PrintDebug("cli.parseArgs() | arg: %q", arg)
					csvFiles = append(csvFiles, arg)
				} else if arg == "--" {
					argsOnly = true
				} else {
					for _, char := range arg[1:] {
						short := fmt.Sprintf("-%s", string(char))
						switch short {
						case "-e", "-f", "-t":
							parseArgs([]string{"skip", short, args[i+1]})
							argSkip = true
						default:
							parseArgs([]string{"skip", short})
						}
					}
				}
			} else {
				app.PrintDebug("cli.parseArgs() | arg: %q", arg)
				csvFiles = append(csvFiles, arg)
			}
		}
	}
}

func parseCSV(reader io.Reader) (err error) {
	r := csv.NewReader(reader)

	var templateRunner *csvrunner.TemplateRunner

	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			app.PrintDebug("parse csv error: %v\n", err)
			err = fmt.Errorf("parse csv error: %v\n", err)
		}
		if templateRunner == nil {
			templateRunner, err = csvrunner.NewTemplateRunner(templateString, record)
			if err != nil {
				return err
			}
			continue
		}
		app.PrintDebug("cli.parseCSV() | record: %q", record)
		// _, err = templateRunner.Run(record, templateOutput)
		out, err := templateRunner.Run(record, templateOutput)
		if templateOutput {
			fmt.Print(out)
		}
	}

	return err
}

func Run(args []string) {
	if len(args) == 1 {
		help()
		os.Exit(0)
	}

	var err error

	parseArgs(args)

	app.PrintDebug("cli.Run() | csvFiles: %q", csvFiles)

	fi, _ := os.Stdin.Stat()
	if (fi.Mode() & os.ModeCharDevice) == 0 {
		// Data is from a Pipe
		reader := bufio.NewReader(os.Stdin)
		err = parseCSV(reader)
		if err != nil {
			app.PrintError("stdin parsing error: %v", err)
		}
	}

	for _, file := range csvFiles {
		app.PrintDebug("cli.Run() | csvFiles | file: %q", file)
		if _, err := os.Stat(file); os.IsNotExist(err) {
			app.PrintError("could not find %q", file)
			continue
		}

		csvFile, _ := os.Open(file)
		reader := bufio.NewReader(csvFile)
		err = parseCSV(reader)
		if err != nil {
			app.PrintError("csv file parsing error: %v", err)
		}
	}
}
