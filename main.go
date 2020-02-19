package main

import (
	"bytes"
	"debug/elf"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/BurntSushi/toml"
)

type Configuration struct {
	Symbols []string
}

var conf *Configuration
var defaultConf = `Symbols = [
	"_ZN7QbModelC1Ev",
	"_ZN7QbModelD1Ev",
]`

func main() {
	var filePath string
	if len(os.Args) > 1 {
		filePath = os.Args[1]
	} else {
		Fatalln("error: missing input file")
	}
	if !isExists(filePath) {
		Fatalln("error: the specified binary doesn't exists")
	}

	if !isExists("conf.toml") {
		fmt.Println("info: couldn't find conf.toml, writing default configuration...")
		if err := ioutil.WriteFile("conf.toml", []byte(defaultConf), 0644); err != nil {
			Fatalf("error: failed to write default configuration\n%s", err.Error())
		}
	}
	conf = &Configuration{}
	if _, err := toml.DecodeFile("conf.toml", conf); err != nil {
		Fatalf("error: failed to parse configuration\n%s", err.Error())
	}

	if len(conf.Symbols) == 0 {
		fmt.Println("aborting because there's nothing to find (Symbols in conf.toml is empty)")
		os.Exit(0)
	}

	fmt.Println("info: opening binary...")
	file, err := elf.Open(filePath)
	if err != nil {
		Fatalf("error: failed to open the specified binary\n%s", err.Error())
	}

	fmt.Println("info: getting dynamic symbols...")
	symbols, err := file.DynamicSymbols()
	if err != nil {
		Fatalf("error: failed to get dynamic symbols\n%s", err.Error())
	}

	fmt.Println("info: getting offsets...")
	var offsets bytes.Buffer
	for _, symbol := range symbols {
		for _, whitelistedSymbolName := range conf.Symbols {
			if symbol.Name == whitelistedSymbolName {
				offsets.Write([]byte(fmt.Sprintf("%s = 0x%x\n", symbol.Name, symbol.Value)))
			}
		}
	}

	if offsets.Len() > 0 {
		fmt.Println("info: writing offsets to a file...")
		if err := ioutil.WriteFile("offsets.txt", offsets.Bytes(), 0644); err != nil {
			Fatalf("error: failed to write offsets to a file\n%s", err.Error())
		}
	} else {
		fmt.Println("warning: len of offsets is zero")
	}
	fmt.Println("finished")
}

func isExists(path string) bool {
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

func Fatalf(format string, a ...interface{}) {
	fmt.Printf(format, a...)
	os.Exit(1)
}

func Fatalln(a ...interface{}) {
	fmt.Println(a...)
	os.Exit(1)
}
