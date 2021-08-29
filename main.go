//
//  main.go - https://github.com/WestleyR/gnotes
//  gnotes - CLI based S3 syncing note app
//
// Created by WestleyR <westleyr@nym.hush.com> on 2021-08-28
// Source code: https://github.com/WestleyR/gnotes
//
// Copyright (c) 2021 WestleyR. All rights reserved.
// This software is licensed under a BSD 3-Clause Clear License.
// Consult the LICENSE file that came with this software regarding
// your rights to distribute this software.
//

package main

import (
	"fmt"
	"github.com/spf13/pflag"
	"log"
)

var Version string = "v0.1.0"

func main() {
	helpFlag := pflag.BoolP("help", "h", false, "print this help output.")
	versionFlag := pflag.BoolP("version", "V", false, "print srm version.")

	pflag.Parse()

	if *helpFlag {
		fmt.Printf("Copyright (c) 2021 WestleyR. All rights reserved.\n")
		fmt.Printf("This software is licensed under the terms of The Clear BSD License.\n")
		fmt.Printf("Source code: https://github.com/.../gnotes\n")
		fmt.Printf("\n")
		pflag.Usage()
		return
	}

	if *versionFlag {
		fmt.Printf("%s\n", Version)
		return
	}

	// Init the self app
	self, err := initApp(getFileFromConfig("config.ini"))
	if err != nil {
		log.Fatalf("Failed to init app: %s\n", err)
	}

	// Load the notes either from s3, or local stash
	err = self.loadNotes()
	if err != nil {
		log.Fatalf("Failed to load notes: %s\n", err)
	}

	// Start the app
	self.loadUI()

	fmt.Println("END")
}
