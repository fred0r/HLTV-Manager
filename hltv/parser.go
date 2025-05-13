package hltv

import (
	"fmt"
	"regexp"
	"strings"
)

type Status int

const (
	HLTV_CHECK_INIT   Status = iota // 0 HLTV инициализируется
	HLTV_FAIL_INIT                  // 1 HLTV иницилизирован
	HLTV_CONNECT                    // 2 HLTV подключается к серверу
	HLTV_FAIL_CONNECT               // 3 У HLTV проблеммы с подключением к серверу
	HLTV_RECORD                     // 4 HLTV Записывает демо
	HLTV_GOOD                       // 5 HLTV Все хорошо
)

type Parser struct {
	Status      Status
	initModules map[string]bool
}

func (hltv *HLTV) TerminalControl() {
	buf := make([]byte, 1024)
	hltv.Parser = newParser()
	for {
		n, err := hltv.Docker.Attach.Reader.Read(buf)
		if err != nil {
			break
		}
		line := string(buf[:n])
		line = strings.TrimSpace(line)
		hltv.ParseHltvOutLines(line)
	}
}

var timeoutPattern = regexp.MustCompile(`^WARNING! Server::Challenge: Timeout after \d+ retries$`)
var rejectedPattern = regexp.MustCompile(`^Connection rejected: No password set.*$`)
var buildPattern = regexp.MustCompile(`^BUILD \d+ SERVER \(\d+ CRC\)$`)
var recordPattern = regexp.MustCompile(`^Start recording to [a-zA-Z0-9]+-\d+-[a-zA-Z0-9_]+\.dem.$`)
var disconnectPattern = regexp.MustCompile(`^Disconnected.*$`)

func (hltv *HLTV) ParseHltvOutLines(input string) {
	lines := strings.Split(input, "\n")
	for _, line := range lines {
		line = strings.ReplaceAll(line, "\r", "")
		//fmt.Printf("[%s]\n", line)
		if line == "" {
			continue
		}
		switch hltv.Parser.Status {
		case HLTV_CHECK_INIT:
			if initialized, exists := hltv.Parser.initModules[line]; exists {
				if !initialized {
					hltv.Parser.initModules[line] = true
				}
			}

			if hltv.allModulesInitialized() {
				hltv.Parser.Status = HLTV_CONNECT
				fmt.Println("ОК! Все модули инициализированы.")
			}
			continue
		case HLTV_CONNECT, HLTV_FAIL_CONNECT:
			{
				if timeoutPattern.MatchString(line) {
					fmt.Println("HLTV Не может подключиться к серверу:", hltv.Settings.Connect)
				} else if rejectedPattern.MatchString(line) {
					fmt.Println("HLTV Не может подключиться к серверу из-за пароля:", hltv.Settings.Connect)
				} else if buildPattern.MatchString(line) {
					fmt.Println("HLTV Подключился к серверу:", hltv.Settings.Connect)
					hltv.Parser.Status = HLTV_RECORD
				}
				continue
			}
		case HLTV_RECORD:
			{
				if recordPattern.MatchString(line) {
					fmt.Println("Началась запись демки:", line)
					hltv.Parser.Status = HLTV_GOOD
					hltv.DemoControl()
				}
				continue
			}
		case HLTV_GOOD:
			{
				if buildPattern.MatchString(line) {
					hltv.Parser.Status = HLTV_RECORD
					fmt.Println("HLTV ПЕРЕПОДКЛЮЧАЕТСЯ")
				} else if disconnectPattern.MatchString(line) {
					hltv.Parser.Status = HLTV_CONNECT
					fmt.Println("HLTV Кикнули?")
				}
				continue
			}
		}

		fmt.Printf("[%s]\n", line)
	}
}

func (hltv *HLTV) allModulesInitialized() bool {
	for _, initialized := range hltv.Parser.initModules {
		if !initialized {
			return false
		}
	}
	return true
}

func newParser() *Parser {
	return &Parser{
		Status: HLTV_CHECK_INIT,
		initModules: map[string]bool{
			"Console initialized.":       false,
			"FileSystem initialized.":    false,
			"Network initialized.":       false,
			"Master module initialized.": false,
			"Server module initialized.": false,
			"World module initialized.":  false,
			"Demo client initialized.":   false,
			"Executing file hltv.cfg.":   false,
			"Proxy module initialized.":  false,
			"Recording initialized.":     false,
		},
	}
}
