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
	"fmt"
	"log"
	"path/filepath"

	"github.com/WestleyR/gnotes"
)

//export GnotesTest
func GnotesTest() *C.char {
	fmt.Println("Hello world from golang!")

	return C.CString("Hello world from golang!")
}

//export Download
// Download will download the note index file to the specified location in the
// config. The config should have all the needed s3 access tokens and user id.
func Download(config_file *C.char) {
	fmt.Println("Download called")

	configFile := C.GoString(config_file)

	app, err := gnotes.InitApp(configFile)
	if err != nil {
		log.Fatalf("Failed to init app: %s\n", err)
	}

	err = app.LoadNotes()
	if err != nil {
		log.Fatalf("Failed to load notes: %s\n", err)
	}

	fmt.Println("Download finished")
}

//export DownloadNote
// DownloadNote will download a specific note from the json index file.
func DownloadNote(config_file *C.char, json_note_path *C.char) {
	fmt.Println("DownloadNote called")

	configFile := C.GoString(config_file)

	app, err := gnotes.InitApp(configFile)
	if err != nil {
		log.Fatalf("Failed to init app: %s\n", err)
	}

	notePath := filepath.Join(app.Config.App.NoteDir, "notes", C.GoString(json_note_path))

	err = app.Config.S3.DownloadFileFrom(filepath.Join(app.Config.S3.UserID, "notes", C.GoString(json_note_path)), notePath)
	if err != nil {
		log.Fatalf("Failed to download file: %s", err)
	}
}

//export Save
// Save will save the whole all and any note if needed, also downloads/uploads
// the new index.
func Save(config_file *C.char) {
	fmt.Println("Save called")

	configFile := C.GoString(config_file)

	app, err := gnotes.InitApp(configFile)
	if err != nil {
		log.Fatalf("Failed to init app: %s\n", err)
	}

	err = app.LoadNotes()
	if err != nil {
		log.Fatalf("Failed to load notes: %s\n", err)
	}

	for _, note := range app.Notes.Books[0].Notes {
		err := note.Save()
		if err != nil {
			log.Fatalf("Failed to save note: %s", err)
		}
	}

	err = app.SaveIndexFile()
	if err != nil {
		log.Fatalf("Failed to upload notes: %s", err)
	}
}
