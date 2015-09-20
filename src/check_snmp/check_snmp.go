package main

import (
	"flag"
	"fmt"
	"github.com/alouca/gosnmp"
	"os"
	"utils"
)

func main() {

	//    for i := 0; i < 100000 ; i++ {
	//    fmt.Print("weee\a")
	//   }
	ipPtr := flag.String("ip", " ", "ip or hostname")
	commPtr := flag.String("community", "nagios", "SNMP community string")
	flag.Parse()

	if *ipPtr == " " {
		fmt.Printf("-ip=<ip or hostname>")
		os.Exit(utils.UNKNOWN)
	}

	//else {
	//fmt.Printf("wordPtr is nil")
	//	}

	ip := *ipPtr
	comm := *commPtr
	timeout := 5

	s, err := gosnmp.NewGoSNMP(ip, comm, gosnmp.Version2c, timeout)

	if err != nil {
		fmt.Printf("UNKNOWN: Could not create snmp object for %s", ip)
		os.Exit(utils.UNKNOWN)
	}

	resp, err := s.Get(".1.3.6.1.2.1.1.1.0")

	if err != nil {
		fmt.Printf("UNKNOWN: Could not connect to %s", ip)
		os.Exit(utils.UNKNOWN)
	}

	if err == nil {
		for _, v := range resp.Variables {
			switch v.Type {
			case gosnmp.OctetString:
				fmt.Printf("Response: %s : %s : %s \n", v.Name, v.Value.(string), v.Type.String())
			}
		}
	}

}
