//
//  ui.go - https://github.com/WestleyR/gnotes
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
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/WestleyR/gnotes"
	"github.com/rivo/tview"
)

type gui struct {
	app      *tview.Application
	noteList *tview.List
}

func newGUI() *gui {
	return &gui{
		app:      tview.NewApplication(),
		noteList: tview.NewList(),
	}
}

func (self *gui) loadUI(app *gnotes.SelfApp) {
	self.app = tview.NewApplication()
	self.noteList = tview.NewList()

	self.reloadUI(app)

	self.noteList.Box.SetBorder(true)
	self.noteList.Box.SetTitle(" Your notes ")
	self.app.SetRoot(self.noteList, true)
	self.app.EnableMouse(true)

	if err := self.app.Run(); err != nil {
		panic(err)
	}
}

func (self *gui) reloadUI(app *gnotes.SelfApp) {
	self.noteList.Clear()

	self.noteList.AddItem("Create new note", "", 'n', func() {
		err := app.Notes.Books[0].NewNote(
			app.Config.App.NoteDir,
			func() {
				err := self.openNote(app, len(app.Notes.Books[0].Notes)-1)
				if err != nil {
					log.Printf("Error opening new note: %s", err)
				}
			})
		if err != nil {
			log.Printf("Failed creating new note: %s", err)
		}
	})

	for i, n := range app.Notes.Books[0].Notes {
		self.noteList.AddItem(n.GetTitle(app.Config.App.NoteDir+"/notes"), n.Info(), getShortcutForIndex(i), func() {
			err := self.openNote(app, self.noteList.GetCurrentItem()-1)
			if err != nil {
				log.Printf("Failed to open note at index: %d: %s", self.noteList.GetCurrentItem()-1, err)
			}
		})
	}

	self.noteList.AddItem("Quit", "Press to exit", 'q', func() {
		// Stop the gui on quit. Will save the notes at main exit.
		self.app.Stop()
	})
}

func (self *gui) openNote(app *gnotes.SelfApp, index int) error {
	// Quit the app before opening the text editor
	self.app.Stop()

	if app.Notes.Books[0].Notes[index].IsAttachment {
		for attempts := 0; attempts < 10; attempts++ {
			succeed := true
			action := ""

			fmt.Printf(`What do you want to do with: %s:
	  d      - Download
	  e      - Edit name
	  delete - Delete the attachment
	  b      - Back
	: `, filepath.Base(app.Notes.Books[0].Notes[index].S3Path))
			fmt.Scanln(&action)

			switch action {
			case "d":
				downloadTo := app.Notes.Books[0].Notes[index].AttachmentTitle
				fmt.Printf("Downloading %s to %s...\n", downloadTo, downloadTo)

				err := app.Config.S3.DownloadFileFrom(
					filepath.Join(app.Config.S3.UserID, "notes", app.Notes.Books[0].Notes[index].S3Path),
					downloadTo,
				)
				if err != nil {
					return fmt.Errorf("failed to download attachment: %s", err)
				}

				fmt.Printf("Downloaded attachment (%s) to: %s\n", app.Notes.Books[0].Notes[index].Title, downloadTo)
				app.Notes.Sort()
				self.loadUI(app)
			case "e":
				return fmt.Errorf("not impmented")
			case "delete":
				fmt.Printf("Deleting %s...\n", app.Notes.Books[0].Notes[index].Title)

				err := app.DeleteNote(0, index)
				if err != nil {
					return fmt.Errorf("failed to delete attachment: %w", err)
				}

				app.Notes.Sort()
				self.loadUI(app)
			case "b":
				self.loadUI(app)
			default:
				log.Printf("Unknown input: %s", action)
				time.Sleep(1 * time.Second)
				succeed = false
			}

			if succeed {
				break
			}

			// Yes, I know, loop of 10. But I want to be able to print a message if over.
			if attempts >= 9 {
				return fmt.Errorf("too many attempts")
			}
		}

		return nil
	}

	// Make sure the note is up-to-date
	err := app.Notes.Books[0].Notes[index].Download(app.Config.App.NoteDir, app.Config.S3)
	if err != nil {
		return err
	}

	// Run the command to open the text file with the specified editor
	editFile := filepath.Join(app.Config.App.NoteDir, "notes", app.Notes.Books[0].Notes[index].S3Path)
	cmd := exec.Command(app.Config.App.Editor, editFile)
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr

	err = cmd.Run()
	if err != nil {
		return err
	}

	// Check if the file is empty
	r, err := os.Open(editFile)
	if err != nil {
		return fmt.Errorf("failed to open file: %s", err)
	}
	defer r.Close()

	b, err := os.ReadFile(editFile)
	if err != nil {
		return fmt.Errorf("failed to read file: %s", err)
	}

	if string(b) == "" {
		// Note is empty, so delete it
		err := app.DeleteNote(0, index)
		if err != nil {
			return fmt.Errorf("failed to delete empty note: %s", err)
		}

		//app.NotesChanged = true
		// TODO: delete note

		app.Notes.Sort()
		self.loadUI(app)

		return nil
	}

	// Save the note if needed
	err = app.Notes.Books[0].Notes[index].Save()
	if err != nil {
		return fmt.Errorf("failed to save the note: %w", err)
	}

	// Resort the notes
	app.Notes.Sort()

	self.loadUI(app)

	return nil
}

func getShortcutForIndex(index int) rune {
	var s = []rune{'1', '2', '3', '4', '5', '6', '7', '8', '9'}

	if index >= len(s) {
		return 0
	}

	return s[index]
}
