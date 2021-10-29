package main

import (
	"flag"
	"fmt"
	"time"
	"github.com/mit-pdos/gokv/grove_ffi"
	"github.com/mit-pdos/gokv/connman"
	"github.com/mit-pdos/gokv/lockservice"
	"github.com/mit-pdos/gokv/memkv"
	"os"
	"strconv"
	"math/rand"
)

var done bool

func client(r int64, d time.Duration, coord memkv.HostName, c chan uint64) {
	ck := lockservice.MakeLockClerk(coord, connman.MakeConnMan())
	var x uint64 = 0
	for !done {
		k := uint64(rand.Int63n(r))
		ck.Lock(k)
		time.Sleep(d * time.Millisecond)
		ck.Unlock(k)
		x++
	}
	fmt.Printf("Done!\n")
	c <- x
}

func main() {
	var coordStr string
	flag.StringVar(&coordStr, "coord", "", "address of coordinator")
	flag.Parse()

	usage_assert := func(b bool) {
		if !b {
			flag.PrintDefaults()
			fmt.Println("Must provide command in form:")
			fmt.Println(" range RANGE")
			fmt.Println(" duration DURATION")
			os.Exit(1)
		}
	}

	usage_assert(coordStr != "")

	coord := grove_ffi.MakeAddress(coordStr)
	ck := lockservice.MakeLockClerk(coord, connman.MakeConnMan())
	var _ = ck

	n := 4
	var r int64 = 10
	var d time.Duration = 0
	a := flag.Args()
	usage_assert(len(a) > 0)
	if a[0] == "range" {
		usage_assert(len(a) >= 2)
		v, err := strconv.ParseInt(a[1], 10, 64)
		usage_assert(err == nil)
		r = v
	}
	if a[2] == "duration" {
		usage_assert(len(a) >= 4)
		v, err := strconv.ParseInt(a[3], 10, 64)
		usage_assert(err == nil)
		d = time.Duration(v)
	}

	fmt.Printf("RANGE = %d; DURATION = %d\n", r, d)

	done = false
	ch := make(chan uint64, n)
	for i := 0; i < n; i++ {
		go client(r, d, coord, ch)
	}
	time.Sleep(3 * time.Second)
	done = true

	for i := 0; i < n; i++ {
		x := <- ch
		fmt.Printf("x = %d.\n", x)
	}
}

