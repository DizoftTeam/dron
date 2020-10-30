package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"

	"gopkg.in/yaml.v2"
)

var (
	debugMode = true
)

const (
	input  = "$input"
	env    = "env"
	argSym = "$"
)

type iCommand struct {
	Name     string            `yaml:"name"`
	Args     map[string]string `yaml:"args"`
	Commands []string          `yaml:"commands"`
}

type config struct {
	Commands []iCommand `yaml:"commands"`
}

// Печать отладочной информации
func debug(data ...interface{}) {
	if debugMode == true {
		fmt.Println("[debug]", data)
	}
}

// Проверка что в строке еще есть аргументы
func checkHasArgPointer(command string) int {
	for i, v := range command {
		if (string(v)) == argSym {
			return i
		}
	}

	return -1
}

func parseArgs(args map[string]interface{}, command string) string {
	result := command

	for argPos := checkHasArgPointer(result); argPos > -1; argPos = checkHasArgPointer(result) {
		length := len(result)

		argName := ""
		pos := -1

		for j := argPos + 1; j < length; j++ {
			ch := string(result[j])
			argName += ch

			if ch == " " || j == length-1 {
				pos = j
				break
			}
		}

		argName = strings.Trim(argName, "\n\t ")

		if argName != "" && args[argName] != nil && pos != -1 {
			argParam := args[argName].(string)

			debug("ARG_NAME", argName, argParam)

			end := ""

			if pos == length-1 {
				end = ""
			} else {
				end = result[pos:]
			}

			result = fmt.Sprintf("%s%s%s", result[:argPos], argParam, end)
		} else {
			debug("[error]", argName, args[argName], pos)

			log.Fatal(fmt.Sprintf("Argument $%s not found in `args` block", argName))
		}
	}

	return result
}

func main() {
	parsed := parseArgs(map[string]interface{}{
		"world": "WiRight",
		"padla": "Loh",
	}, "hello $world suka $padla")

	println("\nPARSED:", parsed)

	return

	// Параметры командной строки
	isDebug := flag.Bool("debug", false, "Print debug info")

	flag.Parse()

	debugMode = *isDebug
	// -------------------------

	c := config{}

	data, err := ioutil.ReadFile("./dron.yaml")

	if err != nil {
		log.Fatal("File 'dron.yaml' can not be located in current folder")
	}

	if err := yaml.Unmarshal(data, &c); err != nil {
		log.Fatal("Cant read config!\n", err)
	}

	// Проверка что есть первый аргумент
	var fArg string

	for _, v := range os.Args[1:] {
		// Собираем только аргументы без "-"
		if !strings.Contains(v, "-") {
			fArg = v
		}
	}

	if fArg == "" {
		log.Fatal("Command name not specified")
	}

	commandExist := false
	var command iCommand

	for _, k := range c.Commands {
		if k.Name == fArg {
			commandExist = true
			command = k

			break
		}
	}

	if commandExist == false {
		log.Fatal("Unknown command \"", fArg, "\"")
	}

	for i, k := range command.Commands {
		debug("Run [", i, "] command")
		debug("--> ", k)

		cmd := exec.Command("bash", "-c", k)
		cmd.Stdout = os.Stdout

		if err := cmd.Run(); err != nil {
			log.Fatal("Cant run command!\n", err)
		}

		//if err := cmd.Wait(); err != nil {
		//	log.Fatal("Cant Wait end of command!", err)
		//}
	}

	println(">>> Done!")
}
