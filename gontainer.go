package main

import (
	"fmt"
	"os"
)

// Custom tools
func must(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	if 3 != len(os.Args) {
		fmt.Printf("Usage: %s [start|stop|run] cmd\n", os.Args[0])
		return
	}

	switch os.Args[1] {
	case "start":
		bridgeCreate()
		makeFsPrivate()
	case "stop":
		bridgeDestroy()
	case "run":
		parent()
	case "child":
		child()
	default:
		panic("What should I do ?")
	}
}
