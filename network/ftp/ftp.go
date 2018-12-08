package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/golang/glog"
	expect "github.com/google/goexpect"
	"github.com/google/goterm/term"
)

const (
	timeout = 10 * time.Minute
)

//package google.golang.org/grpc/codes: unrecognized import path "google.golang.org/grpc/codes" (https fetch: Get https://google.golang.org/grpc/codes?go-get=1: dial tcp 216.58.195.81:443: i/o timeout)

var (
	addr = flag.String("address", "", "address of telnet server")
	user = flag.String("user", "", "username to use")
	pass = flag.String("pass", "", "password to use")
)

func main() {
	flag.Parse()
	fmt.Println(term.Bluef("Ftp 1 example"))

	e, _, err := expect.Spawn(fmt.Sprintf("ftp %s", *addr), -1)
	if err != nil {
		glog.Exit(err)
	}
	defer e.Close()

	e.ExpectBatch([]expect.Batcher{
		&expect.BExp{R: "username:"},
		&expect.BSnd{S: *user + "\n"},
		&expect.BExp{R: "password:"},
		&expect.BSnd{S: *pass + "\n"},
		&expect.BExp{R: "ftp>"},
		&expect.BSnd{S: "bin\n"},
		&expect.BExp{R: "ftp>"},
		&expect.BSnd{S: "prompt\n"},
		&expect.BExp{R: "ftp>"},
		&expect.BSnd{S: "mget *\n"},
		&expect.BExp{R: "ftp>"},
		&expect.BSnd{S: "bye\n"},
	}, timeout)

	fmt.Println(term.Greenf("All done"))
}
