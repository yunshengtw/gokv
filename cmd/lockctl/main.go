package main

import (
	"flag"
	"fmt"
	"github.com/mit-pdos/gokv/grove_ffi"
	"github.com/mit-pdos/gokv/connman"
	"github.com/mit-pdos/gokv/lockservice"
	"os"
	"strconv"
)

func main() {
	var coordStr string
	flag.StringVar(&coordStr, "coord", "", "address of coordinator")
	flag.Parse()

	usage_assert := func(b bool) {
		if !b {
			flag.PrintDefaults()
			fmt.Println("Must provide command in form:")
			fmt.Println(" lock KEY")
			fmt.Println(" unlock KEY")
			os.Exit(1)
		}
	}

	usage_assert(coordStr != "")

	coord := grove_ffi.MakeAddress(coordStr)
	ck := lockservice.MakeLockClerk(coord, connman.MakeConnMan())

	a := flag.Args()
	usage_assert(len(a) > 0)
	if a[0] == "lock" {
		usage_assert(len(a) == 2)
		k, err := strconv.ParseUint(a[1], 10, 64)
		usage_assert(err == nil)
		ck.Lock(k)
		fmt.Printf("LOCK %d\n", k)
	} else if a[0] == "unlock" {
		usage_assert(len(a) == 2)
		k, err := strconv.ParseUint(a[1], 10, 64)
		usage_assert(err == nil)
		ck.Unlock(k)
		fmt.Printf("UNLOCK %d\n", k)
	}
}
