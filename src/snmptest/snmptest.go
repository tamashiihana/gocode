// Copyright 2012 Andreas Louca. All rights reserved.
// Use of this source code is goverend by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"fmt"
	"github.com/alouca/gosnmp"
	"os"
	"utils"
)


const totalSwapOid = ".1.3.6.1.4.1.2021.4.3.0"
const availSwapOid = ".1.3.6.1.4.1.2021.4.4.0"

var (
	cmdDebug     bool
	cmdCommunity string
	cmdVerbose   bool
	cmdTarget    string
	cmdTimeout   int64
	availSwap    int
	totalSwap    int
	swapLimit    int
)

func init() {

	flag.StringVar(&cmdTarget, "target", "", "Target SNMP Agent")
	flag.StringVar(&cmdCommunity, "community", "public", "SNNP Community")
	flag.BoolVar(&cmdDebug, "debug", false, "Debugging output")
	flag.Int64Var(&cmdTimeout, "timeout", utils.TIMEOUT, "Set the timeout in seconds")
	flag.IntVar(&swapLimit, "limit", 85, "Swap percentage to alert")
	flag.BoolVar(&cmdVerbose, "verbose", false, "Verbose output")
	flag.Parse()
}

func main() {

	if cmdTarget == "" {
		flag.PrintDefaults()
		return
	}

	s, err := gosnmp.NewGoSNMP(cmdTarget, cmdCommunity, gosnmp.Version2c, cmdTimeout)

	if err != nil {
		fmt.Printf("UNKNOWN: Error creating SNMP instance: %s\n", err.Error())
		os.Exit(utils.UNKNOWN)
	}

	if cmdVerbose == true {
		s.SetVerbose(true)
	}

	if cmdDebug == true {
		s.SetDebug(true)
	}

	s.SetTimeout(cmdTimeout)

	resp, err := s.Get(totalSwapOid)
	if err != nil {
		fmt.Printf("UNKNOWN: Error getting response: %s\n", err.Error())
		os.Exit(utils.UNKNOWN)
	} else if cmdVerbose == true {
		fmt.Printf("total swap space [%v]\n", resp.Variables[0].Value)
		totalSwap = resp.Variables[0].Value.(int)
	}

	availResp, err := s.Get(availSwapOid)
	if err != nil {
		fmt.Printf("UNKNOWN: Error getting response: %s\n", err.Error())
		os.Exit(utils.UNKNOWN)
	} else if cmdVerbose == true {
		fmt.Printf("available swap space [%v]\n", availResp.Variables[0].Value)
		availSwap = availResp.Variables[0].Value.(int)
	}

	var usedSwap =  totalSwap - availSwap
	if cmdVerbose == true {
		fmt.Printf("swapLimit = %d\n", swapLimit)
		fmt.Printf("usedSwap = %d\n", usedSwap)
	}


	if availSwap != totalSwap {

		//usedSwap := int(100 - (100 * (totalSwap / availSwap)))
		usedSwap := int(100 - int(100*float64(float64(availSwap)/float64(totalSwap))))

		if availSwap >= swapLimit {
			fmt.Printf("CRITICAL - Swap at [%d%]", usedSwap)
			os.Exit(utils.CRITICAL)
		}
	} else {

		fmt.Printf("OK - Swap at [%d%] Available Swap [%d] Used [%d]", usedSwap, availSwap, totalSwap)
		os.Exit(utils.OK)

	}

}
