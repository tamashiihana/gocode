// Copyright 2012 Andreas Louca. All rights reserved.
// Use of this source code is goverend by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
)

func main() {

	myfloat := (10.0 / 100)
	fmt.Printf("%T\n", myfloat)
	fmt.Printf("%.2f\n", myfloat)

	fmt.Printf("%v\n", 1.0/2) // 0.5
	i := 2
	fmt.Printf("%v\n", 1.0/i) // 0

}
