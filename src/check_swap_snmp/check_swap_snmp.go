// Copyright 2012 Andreas Louca. All rights reserved.
// Use of this source code is goverend by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"fmt"
	"github.com/alouca/gosnmp"
	"github.com/davecheney/profile"
	"os"
	"utils"
)

var (
	cmdCommunity string
	cmdTarget    string
	cmdTimeout   int64
	availSwap    int
	totalSwap    int
	usedSwap     int
	swapLimit    int
)

func init() {

	flag.StringVar(&cmdTarget, "target", "", "Target SNMP Agent")
	flag.StringVar(&cmdCommunity, "community", "public", "SNNP Community")
	flag.Int64Var(&cmdTimeout, "timeout", utils.TIMEOUT, "Set the timeout in seconds")
	flag.IntVar(&swapLimit, "limit", 85, "Swap percentage to alert")
	flag.Parse()
}

func main() {
	defer profile.Start(profile.CPUProfile).Stop()

	if cmdTarget == "" {
		flag.PrintDefaults()
		return
	}

	s, err := gosnmp.NewGoSNMP(cmdTarget, cmdCommunity, gosnmp.Version2c, cmdTimeout)
	s.SetVerbose(true)
	if err != nil {
		fmt.Printf("UNKNOWN: Error creating SNMP instance: %s\n", err.Error())
		os.Exit(utils.UNKNOWN)
	}

	const totalSwapOid string = ".1.3.6.1.4.1.2021.4.3.0"
	const availSwapOid string = ".1.3.6.1.4.1.2021.4.4.0"

	resp, err := s.Get(totalSwapOid)
	if err != nil {
		fmt.Printf("UNKNOWN: Error getting response: %s\n", err.Error())
		os.Exit(utils.UNKNOWN)
	} else {
		totalSwap = resp.Variables[0].Value.(int)
		fmt.Printf("total swap [%d]\n", totalSwap)
	}

	availResp, err := s.Get(availSwapOid)
	if err != nil {
		fmt.Printf("UNKNOWN: Error getting response: %s\n", err.Error())
		os.Exit(utils.UNKNOWN)
	} else {
		availSwap = availResp.Variables[0].Value.(int)
		fmt.Printf("availSwap [%d]\n", availSwap)
	}

	if availSwap != totalSwap {

		usedSwap = int(100 - int(100*float64(float64(availSwap)/float64(totalSwap))))
		fmt.Printf("Used swap = %d\n", usedSwap)

		if usedSwap >= swapLimit {

			fmt.Printf("CRITICAL - SWAP Available [%d] Total [%d] In Use [%d%%]", availSwap, totalSwap, usedSwap)
			os.Exit(utils.CRITICAL)
		}

	}

	fmt.Printf("OK - SWAP Available [%d] Total [%d] In Use [%d%%]", availSwap, totalSwap, usedSwap)
	os.Exit(utils.OK)

}
