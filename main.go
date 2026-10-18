package main

import (
	"fmt"
	"os"
	"runtime"
)

func main() {

	fmt.Println("================================")
	fmt.Println("   Secure System Inspector")
	fmt.Println("================================")
	fmt.Println("Operating System Security Lab")

	hostname, err := os.Hostname()

	if err != nil {
		fmt.Println("error getting the hostname :", err)
		return
	}

	fmt.Println("Hostname:", hostname)
	fmt.Println("Operating system: ", runtime.GOOS)
	fmt.Println("Architeture:", runtime.GOARCH)
	fmt.Println("CPU Cores:", runtime.NumCPU())
	fmt.Println("Go Version:", runtime.Version())
}
