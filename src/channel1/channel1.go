package main

import (
	"fmt"
	"github.com/davecheney/profile"
)

func main() {

	defer profile.Start(profile.CPUProfile).Stop()

	c := make(chan int, 1)

	c <- 1
	fmt.Println(<-c)
	c <- 2
	fmt.Println(<-c)

}
