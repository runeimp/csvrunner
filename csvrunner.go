package csvrunner

import (
	// "encoding/csv"
	// "encoding/json"
	"errors"
	"fmt" // For test debugging
	"io"
	// "io/ioutil"
	"bytes"
	"os"
	"os/exec"
	"regexp"
	"runtime"
	"strings"

	// shellquote "github.com/kballard/go-shellquote"
	"github.com/runeimp/csvrunner/app"
)

var (
	ErrorQuoteCountIsNegative          = errors.New("quote counter went negative")
	ErrorValidationHeaderCaseMismatch  = errors.New("header field case mismatch")
	ErrorValidationHeaderValueMismatch = errors.New("header field value mismatch")
	PatternShellVars                   = regexp.MustCompile(`\$\{?(\w+)\}?`)
)

type TemplateRunner struct {
	headersToUpper []string       // headers with fields trimmed and made uppercase for case insensitive matches
	headersTrimmed []string       // headers with fields trimmed of excess space
	indexMap       map[string]int // To track the correct index for a header field
	templateFields map[string]int // To count the occurrences of a fields in a template
	templateString string
	record         []string
}

func (tr *TemplateRunner) init() (err error) {
	headers := make([]string, len(tr.headersTrimmed))
	copy(headers, tr.headersTrimmed)

	for i, f := range headers {
		ht := strings.TrimSpace(f)
		hu := strings.ToUpper(ht)
		tr.headersTrimmed[i] = ht
		tr.headersToUpper[i] = hu
		tr.indexMap[ht] = i
	}

	err = tr.validate(headers)

	matches := PatternShellVars.FindAllStringSubmatch(tr.templateString, -1)
	for _, m := range matches {
		tr.templateFields[m[1]] += 1
	}
	// err = fmt.Errorf("tr.templateString: '%s' | matches: %#v | tr.templateFields: %+v", tr.templateString, matches, tr.templateFields) // For test debugging

	return err
}

func (tr *TemplateRunner) mapper(field string) string {

	if i, ok := tr.indexMap[field]; ok {
		if i < len(tr.record) {
			return tr.record[i]
		}
	}
	return ""
}

// func (tr *TemplateRunner) commandLineParse(input ...string) (result []string, err error) {
// 	shellString := ""
// 	if len(input) > 0 {
// 		shellString = input[0]
// 	}
// 	result, err = shellquote.Split(shellString)
// 	app.PrintDebug("csvrunner.tr.commandLineParse() | result: %#v", result)
// 	return result, err
// }

func (tr *TemplateRunner) isWhiteSpace(char rune) bool {
	whitespace := []rune{0x09, 0x0A, 0x0B, 0x0C, 0x0D, 0x20} // Horizontal Tab, Line Feed/New Line, Vertical Tab, Form Feed, Carriage Return, Space
	for _, r := range whitespace {
		if char == r {
			return true
		}
	}
	return false
}

func (tr *TemplateRunner) Run(record []string, output bool) (combined string, err error) {
	tr.record = record

	var cmd *exec.Cmd

	if runtime.GOOS == "windows" {
		// For Windows we run the command "raw"

		commandLine := PatternShellVars.ReplaceAllString(tr.templateString, `$$Env:$1`)

		args := []string{"-c", commandLine}
		app.PrintDebug("csvrunner.tr.Run() | sh %v", args)
		cmd = exec.Command("powershell", args...)

		// commandLine := os.Expand(tr.templateString, tr.mapper)

		// clArgs, err := tr.commandLineParse(commandLine)
		// if err != nil {
		// 	return combined, err
		// }
		// command := clArgs[0]
		// args := clArgs[1:]

		// app.PrintDebug("csvrunner.tr.Run() | %s %v", command, args)
		// cmd = exec.Command(command, args...)
	} else {
		// For all other platforms we run through Bourne Shell

		commandLine := tr.templateString

		args := []string{"-c", commandLine}
		app.PrintDebug("csvrunner.tr.Run() | sh %v", args)
		cmd = exec.Command("sh", args...)
	}

	newEnv := []string{}
	for i, h := range tr.headersTrimmed {
		env := fmt.Sprintf("%s=%s", h, record[i])
		newEnv = append(newEnv, env)
	}
	cmd.Env = newEnv // the default is os.Environ()

	outputType := ""

	if output {
		outputType = "CombinedOutput"
	}

	switch outputType {
	case "StdOutErr":
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	case "CombinedOutput":
		var combinedBytes []byte
		combinedBytes, err = cmd.CombinedOutput()
		combined = string(combinedBytes)
	case "MultiWriter":
		// Testing io.MultiWriter
		var stdbuf bytes.Buffer
		cmd.Stdout = io.MultiWriter(os.Stdout, &stdbuf)
		cmd.Stderr = io.MultiWriter(os.Stderr, &stdbuf)
		err = cmd.Run()
		if err != nil {
			return combined, err
		}
		combined = string(stdbuf.Bytes())
	default:
		err = cmd.Run()
	}

	// stdout, err := cmd.StdoutPipe()
	// if err != nil {
	// 	return combined, err
	// }

	// if err = cmd.Start(); err != nil {
	// 	return combined, err
	// }

	// stdoutBytes, err := ioutil.ReadAll(stdout)
	// if err != nil {
	// 	return combined, err
	// }

	// combined = string(stdoutBytes)

	// if err = cmd.Wait(); err != nil {
	// 	return combined, err
	// }

	return combined, err
}

func (tr *TemplateRunner) validate(headers []string) (err error) {
	for i, f := range headers {
		ht := strings.TrimSpace(f)
		hu := strings.ToUpper(ht)
		if hu != tr.headersToUpper[i] {
			err = ErrorValidationHeaderValueMismatch
			// err = fmt.Errorf("%s | f: '%s' | hu: '%s' | tr.headersToUpper[%d]: '%s'", ErrorValidationHeaderValueMismatch, f, hu, i, tr.headersToUpper[i]) // For test debugging
			return err
		}
		if ht != tr.headersTrimmed[i] {
			err = ErrorValidationHeaderCaseMismatch
			// err = fmt.Errorf("%s | f: '%s' | ht: '%s' | tr.headersTrimmed[%d]: '%s'", ErrorValidationHeaderCaseMismatch, f, ht, i, tr.headersTrimmed[i]) // For test debugging
			return err
		}
	}
	return err
}

func NewTemplateRunner(templateString string, headers []string, record ...[]string) (tr *TemplateRunner, err error) {
	tr = &TemplateRunner{
		headersToUpper: make([]string, len(headers)),
		headersTrimmed: make([]string, len(headers)),
		indexMap:       make(map[string]int, len(headers)),
		templateFields: make(map[string]int),
		templateString: templateString,
	}
	copy(tr.headersTrimmed, headers)
	// err = fmt.Errorf("headers: %q | tr.headersTrimmed: %q", headers, tr.headersTrimmed)

	if len(record) > 0 {
		tr.record = record[0]
	}

	err = tr.init()

	return tr, err
}
