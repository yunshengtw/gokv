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
	c <- x
}

func main() {
	var coordStr string
	flag.StringVar(&coordStr, "coord", "", "address of coordinator")
	flag.Parse()

	usage_assert := func(b bool) {
		if !b {
			flag.PrintDefaults()
			fmt.Println("Must provide # keys, duration (ms), and # clients.")
			os.Exit(1)
		}
	}

	usage_assert(coordStr != "")

	coord := grove_ffi.MakeAddress(coordStr)

	a := flag.Args()
	usage_assert(len(a) == 3)
	r, err := strconv.ParseInt(a[0], 10, 64)
	usage_assert(err == nil)
	d, err := strconv.ParseInt(a[1], 10, 64)
	usage_assert(err == nil)
	n, err := strconv.ParseInt(a[2], 10, 0)
	usage_assert(err == nil)

	done = false
	ch := make(chan uint64, n)
	for i := int64(0); i < n; i++ {
		go client(r, time.Duration(d), coord, ch)
	}
	time.Sleep(10 * time.Second)
	done = true

	var x uint64 = 0
	for i := int64(0); i < n; i++ {
		x += <- ch
	}
	fmt.Printf("%d, %d, %d, %d\n", r, d, n, x / 10)
}

