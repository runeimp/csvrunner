PROJECT_NAME := "CSV Runner"
PROJECT_CLI := "csvrun"

alias arc := archive

set dotenv-load := false
set positional-arguments


@_default: _term-wipe
	just --list


# Archive GoReleaser dist
archive: _term-wipe
	#!/bin/sh
	tag="$(git tag --points-at main)"
	app="{{PROJECT_NAME}}"
	arc="${app}_${tag}"

	echo "app = '${app}'"
	echo "tag = '${tag}'"
	echo "arc = '${arc}'"
	if [ ! -e _dist ]; then
		mkdir _dist
	fi
	if [ -e dist ]; then
		echo "Move dist -> _dist/${arc}"
		# mv dist "_dist/${arc}"

		# echo "cd distro"
		# cd distro
		
		printf "pwd = "
		pwd
		
		ls -Alh
	else
		echo "dist directory not found for archiving"
	fi


# Build app
build target: _term-wipe
	#!/bin/sh

	if [ ]

	@# if [ '{{os()}}' = 'windows' ]; then
	@# 	go build -o bin/windows/{{PROJECT_CLI}}.exe ./cmd/{{PROJECT_CLI}}/{{PROJECT_CLI}}.go
	@# else
	@# 	go build -o bin/{{os()}}/{{PROJECT_CLI}} ./cmd/{{PROJECT_CLI}}/main.go
	ls -al bin/*/*

_build-linux:
	GOOS=linux GOARCH=amd64 go build -o bin/linux/{{PROJECT_CLI}} ./cmd/{{PROJECT_CLI}}/{{PROJECT_CLI}}.go

_build-macos:
	GOOS=darwin GOARCH=amd64 go build -o bin/macos/{{PROJECT_CLI}} ./cmd/{{PROJECT_CLI}}/{{PROJECT_CLI}}.go

_build-windows:
	GOOS=windows GOARCH=amd64 go build -o bin/windows/{{PROJECT_CLI}}.exe ./cmd/{{PROJECT_CLI}}/{{PROJECT_CLI}}.go

# Clean up this place!
clean: _term-wipe
	rm -f c.out


# Build distro
distro: _term-wipe
	#!/bin/sh
	# goreleaser
	just archive


# Build and install the app
install: _term-wipe
	#!/bin/sh
	cd cmd/{{PROJECT_CLI}}
	go install


# Run code
run *args: _term-wipe
	@just _run "$@"

_run *args:
	@go run ./cmd/{{PROJECT_CLI}}/{{PROJECT_CLI}}.go "$@"

# Run a test: cli, coverage, or package
@test cmd="package": _term-wipe
	just test-{{cmd}}

# Quick CLI test
test-cli:
	#!/bin/sh
	TEST_ENV="echo \"My name is \${Name} and I prefer '\${Pronouns}'.\""
	export TEST_ENV

	GOARCH=amd64 go build -o {{PROJECT_CLI}}.exe ./cmd/{{PROJECT_CLI}}/{{PROJECT_CLI}}.go
	# printf "" > output.txt
	echo "line one" > output.txt

	echo

	echo "==> Test with debug enabled (-e)"
	just _run -debug -e TEST_ENV data/test.csv

	echo

	echo "==> Test with debug disabled (-ot)"
	just _run -ot "${TEST_ENV}" data/test.csv

	echo

	echo "==> Test with debug disabled (-ot)"
	echo '$ just _run -ot "${TEST_ENV}" data/test.csv'
	just _run -ot "${TEST_ENV}" data/test.csv

	echo

	echo "==> Test with stdin and debug disabled (-oe)"
	echo '$ cat data/test.csv | {{PROJECT_CLI}}.exe -oe TEST_ENV'
	cat data/test.csv | {{PROJECT_CLI}}.exe -oe TEST_ENV

	echo

	echo "==> Test with file writing and debug enabled (-t)"
	echo '$ cat data/test.csv | {{PROJECT_CLI}}.exe -t '"echo \"My name is \${Name} and I prefer '\${Pronouns}'.\" >> output.txt"
	cat data/test.csv | {{PROJECT_CLI}}.exe -debug -t "echo \"My name is \${Name} and I prefer '\${Pronouns}'.\" >> output.txt"

	echo

	echo "$ cat output.txt"
	cat output.txt

	echo

# Run Go Test Coverage
@test-coverage:
	echo "You need to run:"
	echo "go test -coverprofile=c.out"
	echo "go tool cover -func=c.out"

# Run Go Unit Tests
test-package:
	go test
	@# go test parser/*
	@# cd parser; go test


_term-wipe:
	#!/bin/sh
	set -exo pipefail
	if [ ${#VISUAL_STUDIO_CODE} -gt 0 ]; then
		clear
	elif [ ${KITTY_WINDOW_ID} -gt 0 ] || [ ${#TMUX} -gt 0 ] || [ "${TERM_PROGRAM}" = 'vscode' ]; then
		printf '\033c'
	elif [ "${TERM_PROGRAM}" = 'Apple_Terminal' ] || [ "${TERM_PROGRAM}" = 'iTerm.app' ]; then
		osascript -e 'tell application "System Events" to keystroke "k" using command down'
	elif [ -x "$(which tput)" ]; then
		tput reset
	elif [ -x "$(which tcap)" ]; then
		tcap rs
	elif [ -x "$(which reset)" ]; then
		reset
	else
		clear
	fi

