/*
    The MIT License (MIT)
    
    Copyright (c) 2015 zhouwench zhouwench@gmail.com
    
    Permission is hereby granted, free of charge, to any person obtaining a copy
    of this software and associated documentation files (the "Software"), to deal
    in the Software without restriction, including without limitation the rights
    to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
    copies of the Software, and to permit persons to whom the Software is
    furnished to do so, subject to the following conditions:
    
    The above copyright notice and this permission notice shall be included in all
    copies or substantial portions of the Software.
    
    THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
    IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
    FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
    AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
    LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
    OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
    SOFTWARE.
*/
package main

import (
		"./cmd"
		"fmt"
		"flag"
		"os"
		"strings"
	   )

func main() {
    servers := flag.String("servers", "", "zk servers")
    chroot := flag.String("chroot", "", "zk chroot")
	recursive := flag.Bool("r", false, "recursive ?")
	help := flag.Bool("h", false, "help info")
	flag.Parse()
	if *help {
		printUsage()
		os.Exit(0)
	}
	if len(*servers) == 0 {
		fmt.Println("ERROR: --servers is required\n")
		printUsage()
		os.Exit(255)
	}
	if flag.NArg() < 2 {
		fmt.Println("ERROR: cmd path is requred\n")
		printUsage()
		os.Exit(255)
	}
    opt := cmd.NewCmdOption()
	opt.Servers = strings.Split(*servers, ",")
	opt.Chroot = *chroot
	opt.Recursive = *recursive
	err := cmd.Call(opt, flag.Args())
	if err != nil {
		fmt.Println("ERROR: fail to zk:", err)
		os.Exit(255)
	}
	os.Exit(0)
}

func printUsage() {
	fmt.Printf("Usage: \n")
	fmt.Printf(" %s [OPTIONS] CMD PATH [CmdArgs]\n", os.Args[0])

	fmt.Printf("\n")
	fmt.Printf("OPTIONS:\n")
	fmt.Printf(" --servers: required, zk servers\n")
	fmt.Printf(" --chroot: optional, zk chroot\n")
	fmt.Printf(" -r: optional, recursive operation ? \n")

	fmt.Printf("\n")
	fmt.Printf("CMD: [ls|children|get|set|create|delete|exists]\n")

	fmt.Printf("\n")
	fmt.Printf("PATH: zookeeper path\n")
}

