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
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/WestleyR/gnotes"
	"github.com/google/uuid"
	"github.com/spf13/pflag"
)

var Version string = "v0.2.0"

func main() {
	helpFlag := pflag.BoolP("help", "h", false, "print this help output.")
	versionFlag := pflag.BoolP("version", "V", false, "print srm version.")
	uploadFileFlag := pflag.StringP("add-file", "a", "", "add an attachment file.")
	skipDownloadFlag := pflag.BoolP("skip-download", "s", false, "skips downloading the note file, used for devs, or if starting notes from scratch.")
	newNoteFlag := pflag.BoolP("reset", "R", false, "dont fail if local notes dont exist, DANGER: could delete all existing notes!")
	decryptFlag := pflag.StringP("decrypt", "d", "", "decrypt for devs")
	genCryptKeyFlag := pflag.BoolP("gen-crypt-key", "", false, "generate an 16 bit encryption key (for first initalization)")
	genUUIDFlag := pflag.BoolP("gen-uuid", "", false, "generate a uuid for user id (for first initalization)")

	pflag.Parse()

	switch {
	case *helpFlag:
		fmt.Printf("Copyright (c) 2021-2022 WestleyR. All rights reserved.\n")
		fmt.Printf("This software is licensed under the terms of The Clear BSD License.\n")
		fmt.Printf("Source code: https://github.com/WestleyR/gnotes\n")
		fmt.Printf("\n")
		pflag.Usage()
		return

	case *versionFlag:
		if *versionFlag {
			fmt.Printf("%s\n", Version)
			return
		}

	case *genCryptKeyFlag:
		fmt.Println(ranStr(16))
		return

	case *genUUIDFlag:
		u := uuid.New()
		fmt.Println(u.String())
		return
	}

	// Init the self app
	app, err := gnotes.InitApp(gnotes.GetFileFromConfig("config.ini"))
	if err != nil {
		log.Fatalf("Failed to init app: %s\n", err)
	}

	if *decryptFlag != "" {
		b, err := os.ReadFile(*decryptFlag)
		if err != nil {
			log.Fatal(err)
		}
		b, err = app.Config.S3.Decrypt(b)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(string(b))

		return
	}

	// Set the cli opts
	app.CliOpts.SkipDownload = *skipDownloadFlag
	app.CliOpts.NewNote = *newNoteFlag

	// Setup the gui (cli)
	gui := newGUI()

	// Load the notes either from s3, or local stash
	err = app.LoadNotes()
	if err != nil {
		log.Fatalf("Failed to load notes: %s\n", err)
	}

	// Before starting the ui, see if theres anything to be done first
	if *uploadFileFlag != "" {
		err := app.Notes.Books[0].NewAttachment(app.Config.App.NoteDir, *uploadFileFlag)
		if err != nil {
			log.Fatalf("Failed to add attachment: %s", err)
		}

		err = app.SaveIndexFile()
		if err != nil {
			log.Fatalf("Failed to save notes: %s", err)
		}

		// After uploading the file, just quit
		return
	}

	// Start the app
	gui.loadUI(app)

	// Always save the notes if needed
	err = app.SaveIndexFile()
	if err != nil {
		log.Fatalf("Failed to upload notes: %s", err)
	}

	fmt.Println("END")
}

const validKeyChars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"

func ranStr(n int) string {
	rand.Seed(time.Now().UnixNano())

	b := make([]byte, n)
	for i := range b {
		b[i] = validKeyChars[rand.Intn(len(validKeyChars))]
	}

	return string(b)
}
