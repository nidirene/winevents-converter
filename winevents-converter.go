package main

import (
	"bufio"
	"encoding/xml"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"path"
	"strings"
	"syscall"
)

var (
	opts   *Options
	logger *log.Logger
)

func main() {

	var (
		signalCh = make(chan os.Signal, 1)
	)
	signal.Notify(signalCh, syscall.SIGINT, syscall.SIGTERM)

	opts = GetOptions()
	logger = opts.Logger
	flag.Parse()
	//

	file, err := os.Open(path.Join(opts.EventFilePath, "lst.txt"))
	if err != nil {
		logger.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		events := strings.Fields(scanner.Text())
		if len(events) <= 2 {
			opts.Logger.Println(events)
			break
		}
		var (
			event Event
			jb    []byte
		)

		b, err := ioutil.ReadFile(path.Join(opts.EventFilePath, events[0]))
		if err != nil {
			opts.Logger.Println(err)
			return
		}
		err = xml.Unmarshal(b, &event)
		if err != nil {
			opts.Logger.Println(err)
		}

		jb, err = convertToWinC(&event)
		if err != nil {
			opts.Logger.Println(err)
			return
		}
		fmt.Println("-------------------------------------------------------------------------------")
		fmt.Println(events[1:])
		fmt.Println(string(jb))
		fmt.Println()
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

}
