package main

import (
	"flag"
	"fmt"
	"os"
	"time"
)

var (
	from, to      string
	limit, offset int64
)

func init() {
	flag.StringVar(&from, "from", "", "file to read from")
	flag.StringVar(&to, "to", "", "file to write to")
	flag.Int64Var(&limit, "limit", 0, "limit of bytes to copy")
	flag.Int64Var(&offset, "offset", 0, "offset in input file")
}

func main() {
	flag.Parse()
	done := make(chan struct{})
	defer close(done)

	go checkProgress(from, to, offset, limit, done)

	if err := Copy(from, to, offset, limit); err != nil {
		panic(err)
	}
}

func checkProgress(fromPath, toPath string, offset, limit int64, done chan struct{}) {
	fromFile, _ := os.Stat(fromPath)

	max := fromFile.Size()
	if limit == 0 || limit > fromFile.Size()-offset {
		max = fromFile.Size() - offset
	}

	for {
		select {
		case <-done:
			return
		default:
			toFile, _ := os.Stat(toPath)
			if toFile == nil {
				continue
			}
			d := toFile.Size() * 100 / max
			fmt.Printf("\033c %d%%", d)
			time.Sleep(10 * time.Millisecond)
		}
	}
}
