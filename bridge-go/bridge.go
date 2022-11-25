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
	"encoding/json"
	"log"

	"github.com/WestleyR/gnotes"
)

var app *gnotes.SelfApp

//export InitApp
func InitApp(input *C.char) *C.char {
	// Input example: `config=/path/to/config skip_download=true new_note=true`
	opts := parseInput(C.GoString(input))

	var err error

	// Init the self app
	app, err = gnotes.InitApp(opts["config"])
	if err != nil {
		log.Fatalf("Failed to init app: %s", err)
	}

	// Set the cli opts
	app.CliOpts.SkipDownload = isBool(opts["skip_download"])
	app.CliOpts.NewNote = isBool(opts["new_note"])

	return nil
}

//export NewNote
func NewNote(input *C.char) *C.char {
	err := app.NewNote(nil)
	if err != nil {
		log.Fatalf("failed to create new note: %s", err)
	}

	return C.CString("new note created")
}

//export Download
func Download(input *C.char) *C.char {
	err := app.LoadNotes()
	if err != nil {
		log.Fatalf("Failed to load notes: %s", err)
	}

	return C.CString("notes downloaded")
}

//export Save
func Save(input *C.char) *C.char {
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
	jsonData, err := json.Marshal(app.Notes)
	if err != nil {
		log.Fatalf("failed to marshal json: %s", err)
	}

	return C.CString(string(jsonData))
}
