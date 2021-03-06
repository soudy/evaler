// Copyright 2016 Steven Oud. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be found
// in the LICENSE file.

package main

import (
	"flag"
	"fmt"
	"math/big"
	"os"
	"runtime"

	"github.com/chzyer/readline"
	"github.com/soudy/mathcat"
)

var (
	precision   = flag.Uint("precision", 64, "bits of precision used in decimal float results")
	literalMode = flag.String("mode", "decimal", "type of literal used as result. can be decimal (default), hex, binary or octal")
)

func getHomeDir() string {
	if runtime.GOOS == "windows" {
		home := os.Getenv("HOMEDRIVE") + os.Getenv("HOMEPATH")
		if home == "" {
			home = os.Getenv("USERPROFILE")
		}

		return home
	}

	return os.Getenv("HOME")
}

func repl(mode Mode) {
	p := mathcat.New()
	rl, err := readline.NewEx(&readline.Config{
		Prompt:      "mc> ",
		HistoryFile: getHomeDir() + "/.mathcat_history",
	})

	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(-1)
	}
	defer rl.Close()

	for {
		line, err := rl.Readline()
		if err != nil {
			break
		}

		res, err := p.Run(line)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			continue
		}

		switch mode {
		case Decimal:
			if res.IsInt() {
				fmt.Println(res.Num())
			} else {
				stringResult := new(big.Float).
					SetPrec(*precision).
					SetRat(res).
					Text('f', -1)
				fmt.Println(stringResult)
			}
		case Hex, Binary, Octal:
			formats := map[Mode]string{
				Hex:    "%#x",
				Binary: "%b",
				Octal:  "%#o",
			}
			integer := mathcat.RationalToInteger(res)
			fmt.Printf(formats[mode]+"\n", integer)
		}
	}
}

func main() {
	var mode Mode
	var ok bool

	flag.Parse()

	if mode, ok = modes[*literalMode]; !ok {
		fmt.Fprintf(os.Stderr, "Invalid mode type ‘%s’\n", *literalMode)
		os.Exit(-1)
	}

	repl(mode)
}
