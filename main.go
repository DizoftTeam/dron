package main

import (
  "bufio"
  "flag"
  "fmt"
  "io/ioutil"
  "log"
  "os"
  "os/exec"
  "strings"

  _ "github.com/joho/godotenv/autoload"
  "gopkg.in/yaml.v2"
)

const (
  input   = "input" // Read info from user input
  env     = "env"   // Get param from env
  version = "1.2.2" // Application version
)

var (
  debugMode = false // Debug mode - print debug info
  argSym    = "$"   // Start of argument TODO: maybe set in config
)

// Command structure
type iCommand struct {
  Name     string            `yaml:"name"`
  Args     map[string]string `yaml:"args"`
  Commands []string          `yaml:"commands"`
}

// Struct of dron.yaml file
type config struct {
  Commands []iCommand `yaml:"commands"`
}

// Print debug into
func debug(data ...interface{}) {
  if debugMode == true {
    fmt.Println("[debug]", data)
  }
}

// Check that command has arguments
func checkHasArgPointer(command string) int {
  for i, v := range command {
    if (string(v)) == argSym {
      return i
    }
  }

  return -1
}

// Parsing and processing ENV argument
func parseEnv(arg string) string {
  // $env(ENV_NAME)
  envName := arg[5 : len(arg)-1]

  envVal, exist := os.LookupEnv(envName)

  if exist == false {
    log.Fatal(fmt.Sprintf("Env param '%s' not exist!", envName))
  }

  return envVal
}

// Parsing and processing INPUT argument
func parseInput(argName string) string {
  scanner := bufio.NewScanner(os.Stdin)

  fmt.Printf("[%s]> ", argName)

  scanner.Scan()

  return scanner.Text()
}

// Parse arguments
func parseArgs(args map[string]string, command string) string {
  result := command

  for argPos := checkHasArgPointer(result); argPos > -1; argPos = checkHasArgPointer(result) {
    length := len(result)

    argName := ""
    pos := -1

    for j := argPos + 1; j < length; j++ {
      ch := string(result[j])
      argName += ch

      if strings.ContainsAny(ch, "\n\t\"' ") || j == length-1 {
        pos = j

        if strings.ContainsAny(ch, `"'(){}[]`) {
          pos = j - 1
        }

        break
      }
    }

    argName = strings.Trim(argName, "\n\t'\" ")

    if argName != "" && args[argName] != "" && pos != -1 {
      argParam := args[argName]

      debug("ARG_NAME", argName, argParam)

      end := ""

      if pos == length-1 {
        end = ""
      } else {
        end = result[pos:]
      }

      // CHECK ENV ---

      if strings.Contains(argParam, env) {
        argParam = parseEnv(argParam)
      }

      // --- CHECK ENV

      // CHECK INPUT ---

      if strings.Contains(argParam, input) {
        argParam = parseInput(argName)
      }

      // --- CHECK INPUT

      result = fmt.Sprintf("%s%s%s", result[:argPos], argParam, end)
    } else {
      debug("[error]", argName, args[argName], pos)

      log.Fatal(fmt.Sprintf("Argument $%s not found in `args` block", argName))
    }
  }

  return result
}

func main() {
  // CLI arguments
  isDebug := flag.Bool("debug", false, "Print debug info")
  showList := flag.Bool("list", false, "Print list commands")
  showVersion := flag.Bool("version", false, "Show version")

  flag.Parse()

  debugMode = *isDebug
  // -------------------------

  c := config{}

  data, err := ioutil.ReadFile("./dron.yaml")

  if err != nil {
    data2, err2 := ioutil.ReadFile("./dron.yml")

    if err2 != nil {
      log.Fatal("File 'dron.yaml' or 'dron.yml' can not be located in current folder")
    }

    data = data2
  }

  if err := yaml.Unmarshal(data, &c); err != nil {
    log.Fatal("Cant read config!\n", err)
  }

  // Check that using -list argument
  if *showList == true {
    for _, v := range c.Commands {
      fmt.Println(v.Name)
    }

    return
  }

  // Check that using -version argument
  if *showVersion == true {
    fmt.Println("Dron version:", version)

    return
  }

  // Check the first argument exist
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
    debug("-->", k)

    parsed := parseArgs(command.Args, k)

    debug("[parsed]", parsed)

    cmd := exec.Command("bash", "-c", parsed)
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
