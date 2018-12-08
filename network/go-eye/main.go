package main

import (
	"flag"
	"fmt"
)

func main() {
	act := flag.String("act", "aaa", "selection: aaa;alarmtrans")
	flag.Parse()

	fmt.Printf("Welcome to go-eye : %s\n\n", *act)

	switch *act {
	case "aaa":
		GoPacketAAA()
	case "alarmtrans":
		fmt.Printf("====== this feature is still under development ======\n")
	case "test":
		GoPacketAll()
	}
}
