package main

import (
	"bmining/bminer"
	"fmt"
	"math/rand"
	"time"
)

func main() {
	fmt.Println("Here we go")
	rand.Seed(42)
	rand.Seed(time.Now().UnixNano())
	m := bminer.New(bminer.Charge)
	for id := 1; id < 1000; id++ {
		m.Start(id)
	}
	fmt.Println("All done")
}
