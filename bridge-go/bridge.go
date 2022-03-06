//
// Created by WestleyR <westleyr@nym.hush.com> on 2022-03-06
// Source code: https://github.com/WestleyR/gnotes
//
// Copyright (c) 2022 WestleyR. All rights reserved.
// This software is licensed under a BSD 3-Clause Clear License.
// Consult the LICENSE file that came with this software regarding
// your rights to distribute this software.
//

package main

import "C"

import (
	"log"
	"strings"

	"github.com/WestleyR/gnotes"
)

func parseInput(input, question string) bool {
	if input != "" {
		for _, op := range strings.Split(input, " ") {
			opts := strings.Split(op, "=")

			if len(opts) != 2 {
				log.Fatalf("invalid input: %q", op)
			}

			switch opts[0] {
			case question:
				return opts[1] == "true"
			}
		}
	}

	return false
}

func getApp(input string) *gnotes.SelfApp {
	// Parse the input
	skipDownload := false
	newNote := false

	// Input example: `skip_download=true new_note=true`
	if input != "" {
		for _, op := range strings.Split(input, " ") {
			opts := strings.Split(op, "=")

			if len(opts) != 2 {
				log.Fatalf("invalid input: %q", op)
			}

			switch opts[0] {
			case "skip_download":
			case "new_note":
			default:
				log.Printf("unknown input: %s", opts[0])
			}
		}
	}

	// Init the self app
	app, err := gnotes.InitApp(gnotes.GetFileFromConfig("config.ini"))
	if err != nil {
		log.Fatalf("Failed to init app: %s\n", err)
	}

	// Set the cli opts
	app.CliOpts.SkipDownload = skipDownload
	app.CliOpts.NewNote = newNote

	return app
}

//export Download
func Download(input *C.char) *C.char {
	app := getApp(C.GoString(input))

	err := app.LoadNotes()
	if err != nil {
		log.Fatalf("Failed to load notes: %s\n", err)
	}

	return C.CString("notes downloaded")
}

//export Save
func Save(input *C.char) *C.char {
	app := getApp(C.GoString(input))

	app.NotesChanged = parseInput(C.GoString(input), "notes_changed")

	err := app.SaveNotes()
	if err != nil {
		log.Fatalf("Failed to upload notes: %s", err)
	}

	return C.CString("notes saved")
}
