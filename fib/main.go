package main

import (
	"flag"
	"fmt"
)

type args struct {
	months    int
	offspring int
}

func parseArgs() args {
	result := args{
		months:    0,
		offspring: 0,
	}
	flag.IntVar(&result.months, "m", 0, "How many months to run")
	flag.IntVar(&result.offspring, "o", 0, "Offspring per pair")
	flag.Parse()

	return result
}

type Fibber struct {
	memo map[int]int
	k    int
}

func NewFibber(k int) *Fibber {
	return &Fibber{
		memo: make(map[int]int),
		k:    k,
	}
}

func (self Fibber) Fib(n int) int {
	if n <= 2 {
		return 1
	} else if rval, ok := self.memo[n]; ok {
		return rval
	} else {
		rval := self.Fib(n-1) + (self.k * (self.Fib(n - 2)))
		self.memo[n] = rval
		return rval
	}
}

func main() {
	args := parseArgs()
	fibber := NewFibber(args.offspring)
	fmt.Printf("%d\n", fibber.Fib(args.months))
}
