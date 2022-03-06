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
	"os"
	"path/filepath"

	"github.com/WestleyR/gnotes"
)

func getApp(input string) *gnotes.SelfApp {
	// Input example: `skip_download=true new_note=true`
	opts := parseInput(input)

	// Init the self app
	app, err := gnotes.InitApp(gnotes.GetFileFromConfig("config.ini"))
	if err != nil {
		log.Fatalf("Failed to init app: %s\n", err)
	}

	// Set the cli opts
	app.CliOpts.SkipDownload = isBool(opts["skip_download"])
	app.CliOpts.NewNote = isBool(opts["new_note"])

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

	opts := parseInput(C.GoString(input))

	app.NotesChanged = isBool(opts["notes_changed"])

	err := app.SaveNotes()
	if err != nil {
		log.Fatalf("Failed to upload notes: %s", err)
	}

	return C.CString("notes saved")
}

//export List
func List(input *C.char) *C.char {
	app := getApp(C.GoString(input))

	// Now read the downloaded file
	noteJson, err := os.ReadFile(filepath.Join(app.Config.App.NoteDir, "gnotes", "notes/gnotes.json"))
	if err != nil {
		log.Fatalf("failed to read json: %s", err)
	}

	return C.CString(string(noteJson))
}
