//
//  notes.go - https://github.com/WestleyR/gnotes
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
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/rivo/tview"
)

type appConfigs struct {
	editor      string
	s3Active    bool
	s3Bucket    string
	s3Endpoint  string
	s3Region    string
	s3SaveFile  string
	s3AccessKey string
	s3SecretKey string
}

type selfApp struct {
	notes    []*note
	noteList *tview.List

	// On exit, dont upload if notes did not change
	notesChanged bool

	app *tview.Application

	// TODO: is it okay to store all the configs in the running app?
	// changes wont apply until the app is restarted.
	config appConfigs
}

type note struct {
	DateCreated int64
	DateMod     int64

	Content string
}

func initApp(configPath string) (selfApp, error) {
	var app selfApp

	app.notesChanged = false

	// TODO: dont open/close the config file for every settings,
	// need to open the file once.
	app.config.editor = getEditor()
	app.config.s3Active = getUseS3()
	app.config.s3Bucket = getS3Bucket()
	app.config.s3Endpoint = getS3Endpoint()
	app.config.s3Region = getS3Region()
	app.config.s3SaveFile = getS3FileName()
	app.config.s3AccessKey = getS3AccessKey()
	app.config.s3SecretKey = getS3SecretKey()

	return app, nil
}

// getTitleForNote returns the first line of the note
func getTitleForNote(content string) string {
	return strings.Split(content, "\n")[0]
}

func getSubContentForNote(n *note) string {
	c := time.Unix(n.DateCreated, 0)
	m := time.Unix(n.DateMod, 0)

	return fmt.Sprintf("Created on %s. last modified on %s", c.Format("2006-01-02"), m.Format("2006-01-02"))
}

func (self *selfApp) newNote() error {
	createdTime := time.Now().Unix()

	newNote := &note{
		DateCreated: createdTime,
		DateMod:     createdTime,
		Content:     "",
	}

	self.notes = append(self.notes, newNote)

	// Append the new note
	self.noteList.AddItem(self.notes[len(self.notes)-1].Content, "[new_note]", getShortcutForIndex(len(self.notes)-1), func() {
		err := self.openNote(self.noteList.GetCurrentItem() - 1)
		if err != nil {
			log.Fatalf("Failed opening note index: %d: %s\n", self.noteList.GetCurrentItem()-1, err)
		}
	})

	// Open the new note
	return self.openNote(len(self.notes) - 1)
}

func (self *selfApp) openNote(index int) error {
	// Quit the app before opening the text editor
	self.app.Stop()

	// First, write the note to a file so the editor can open it
	tmpFile := "/tmp/wst.gnotes.current-note"
	err := os.WriteFile(tmpFile, []byte(self.notes[index].Content), 0600)
	if err != nil {
		panic(err)
	}

	// Run the command to open the text file with the specified editor
	cmd := exec.Command(self.config.editor, tmpFile)
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr

	err = cmd.Run()
	if err != nil {
		return err
	}

	// Now read the note back into the struct
	n, err := os.ReadFile(tmpFile)
	if err != nil {
		return err
	}

	// Update the mod time only if the content has changed
	if self.notes[index].Content != string(n) {
		// User edited a note, so make sure to upload on exit
		self.notesChanged = true
		self.notes[index].DateMod = time.Now().Unix()
	}

	self.notes[index].Content = string(n)

	if self.notes[index].Content == "" {
		// User wants to delete this note
		self.notes = append(self.notes[:index], self.notes[index+1:]...)
	}

	// Resort the notes
	sortByModTime(self.notes)

	self.loadUI()

	return nil
}

func sortByModTime(notes []*note) {
	sorting := true

	for sorting {
		sorting = false
		for i := 0; i < len(notes)-1; i++ {
			if notes[i].DateMod < notes[i+1].DateMod {
				tmp := notes[i]
				notes[i] = notes[i+1]
				notes[i+1] = tmp
				sorting = true
			}
		}
	}
}

func (self *selfApp) loadNotes() error {
	// Always use the config dir for save files as backups.
	noteFile := getLocalSaveFile()

	// Download the a tmp file first, just in case the download fails,
	// and we lose all our notes
	downloadNote := "/tmp/gnotes.download.notes.json"

	// Only download from s3 if active = true
	if self.config.s3Active {
		err := s3DownloadFile(self.config.s3SaveFile, downloadNote)
		if err != nil {
			return fmt.Errorf("error downloading the save file: %s: %s\n", self.config.s3SaveFile, err)
		}
	}

	// If not using s3, then only use the local saved notes
	if !self.config.s3Active {
		downloadNote = noteFile
	}

	fmt.Printf("DOWNLOADED: %s\n", downloadNote)
	fmt.Printf("NOTE_FILE: %s\n", noteFile)

	// Now read the downloaded file
	downloadedJson, err := os.ReadFile(downloadNote)
	if err != nil {
		return err
	}

	if self.config.s3Active {
		if string(downloadedJson) == "" {
			return fmt.Errorf("downloaded notes are empty; aborting")
		}
	}

	fmt.Printf("FROM FILE: %s\n", string(downloadedJson))

	if string(downloadedJson) != "" {
		err = json.Unmarshal(downloadedJson, &self.notes)
		if err != nil {
			return err
		}
	}

	// Now sort the notes by mod time
	sortByModTime(self.notes)

	return nil
}

func (self *selfApp) saveNotes() error {
	var jsonData []byte
	jsonData, err := json.Marshal(self.notes)
	if err != nil {
		return err
	}

	// Always use the config dir for save files as backups.
	saveFile := getLocalSaveFile()

	err = os.WriteFile(saveFile, jsonData, 0600)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Errored uploading notes. Just in case, you can recover your notes from this string: %v\n", jsonData)
		return err
	}

	if self.config.s3Active {
		err = s3UploadFile(saveFile, self.config.s3SaveFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Errored uploading notes. Just in case, you can recover your notes from this string: %v\n", jsonData)
			return err
		}
	}

	return nil
}
