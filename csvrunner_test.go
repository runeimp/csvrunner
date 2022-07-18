package csvrunner

import (
	"runtime"
	"testing"
)

const (
	templateNameAge = `echo "${name} is $age years old"`
)

var (
	recordRandolf = []string{"Randolf", "42", "they/them"}
)

func TestCommandLineParsePowerShell(t *testing.T) {
	got := PatternShellVars.ReplaceAllString(templateNameAge, `$$Env:$1`)
	want := `echo "$Env:name is $Env:age years old"`

	if got != want {
		t.Fatalf(`commandLineParse() = '%s', want match for '%s'`, got, want)
	}
}

// func TestCommandLineParse(t *testing.T) {
// 	template := `echo "My name is ${Name} and I prefer '${Pronouns}'." >> output.txt`
// 	headers := []string{"Name", "Age", "Pronouns"}
// 	record := []string{"Lord Henry Hedgefundington VII", "38", "they/them"}

// 	tr, err := NewTemplateRunner(template, headers, record)
// 	if err != nil {
// 		t.Fatalf(`initialization error: %v`, err)
// 	}

// 	got, err := tr.commandLineParse()
// 	if err != nil {
// 		t.Fatalf(`tr.commandLineParse error: %v`, err)
// 	}
// 	want := []string{"echo", `"My name is Lord Henry Hedgefundington VII and I prefer 'they/them'."`, ">>", "output.txt"}
// 	match := true

// 	for i, s := range got {
// 		if i < len(want) && s != want[i] {
// 			match = false
// 		}
// 	}
// 	if match == false {
// 		t.Fatalf(`commandLineParse() = %q, want match for %q`, got, want)
// 	}
// }

func TestTemplateRunnerMapperAge(t *testing.T) {
	template := templateNameAge
	headers := []string{"name", "age", "pronouns"}
	record := recordRandolf

	tr, err := NewTemplateRunner(template, headers, record)
	if err != nil {
		t.Fatalf(`initialization error: %v`, err)
	}

	got := tr.mapper("age")
	want := "42"
	if got != want {
		t.Fatalf(`tr.Validate() = %q, want match for %q`, got, want)
	}
}

func TestTemplateRunnerMapperName(t *testing.T) {
	template := templateNameAge
	headers := []string{"name", "age", "pronouns"}
	record := recordRandolf

	tr, err := NewTemplateRunner(template, headers, record)
	if err != nil {
		t.Fatalf(`initialization error: %v`, err)
	}

	got := tr.mapper("name")
	want := "Randolf"
	if got != want {
		t.Fatalf(`tr.Validate() = %q, want match for %q`, got, want)
	}
}

func TestTemplateRunnerRun(t *testing.T) {
	template := templateNameAge
	headers := []string{"name", "age", "pronouns"}
	record := recordRandolf
	output := true

	tr, err := NewTemplateRunner(template, headers)
	if err != nil {
		t.Fatalf(`initialization error: %v`, err)
	}

	gotCombined, gotErr := tr.Run(record, output)
	wantCombined := "Randolf is 42 years old\n" // Bourne shell echo add \n
	if runtime.GOOS == "windows" {
		wantCombined = "Randolf is 42 years old\r\n" // PowerShell echo adds \r\n
	}
	if gotErr != nil {
		t.Fatalf(`tr.Validate() = %q, want match for nil`, gotErr.Error())
	}
	if gotCombined != wantCombined {
		t.Fatalf(`tr.Validate() = %q, want match for %q`, gotCombined, wantCombined)
	}
}

func TestTemplateRunnerValidationHeader(t *testing.T) {
	template := templateNameAge
	headers := []string{"name", "age", "pronouns"}

	tr, err := NewTemplateRunner(template, headers)
	if err != nil {
		t.Fatalf(`initialization error: %v`, err)
	}

	got := tr.validate(headers)
	if got != nil {
		t.Fatalf(`tr.Validate() = %q, want match for nil`, got)
	}
}

func TestTemplateRunnerValidationHeaderCaseMismatch(t *testing.T) {
	template := "echo ${name}'s' pronouns are $pronouns. At ${age} years old $name is very happy."
	headers := []string{"name", "age", "pronouns"}

	tr, err := NewTemplateRunner(template, headers)
	if err != nil {
		t.Fatalf(`initialization error: %v`, err)
	}

	got := tr.validate([]string{"Name", "Age", "Gender"})
	want := ErrorValidationHeaderCaseMismatch
	if got != want {
		t.Fatalf(`tr.Validate() = %q, want match for ErrorValidationHeaderCaseMismatch`, got)
	}
}

func TestTemplateRunnerValidationHeaderValueMismatch(t *testing.T) {
	template := templateNameAge
	headers := []string{"name", "age", "pronouns"}

	tr, err := NewTemplateRunner(template, headers)
	if err != nil {
		t.Fatalf(`initialization error: %v`, err)
	}

	got := tr.validate([]string{"name", "age", "gender"})
	want := ErrorValidationHeaderValueMismatch
	if got != want {
		t.Fatalf(`tr.Validate() = %q, want match for ErrorValidationHeaderValueMismatch`, got)
	}
}
