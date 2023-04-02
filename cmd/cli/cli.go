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
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	_ "embed"

	"github.com/WestleyR/gnotes"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"golang.org/x/term"
)

//go:embed usage.md
var helpText []byte

const (
	pageNotes = iota
	pageFolders
)

type gui struct {
	// ui is the main gui application.
	ui *tview.Application

	// pages is used for showing other pages like form/text inputs.
	pages *tview.Pages

	currentPage int

	// noteList is used for both the list of notes, and list of folders.
	noteList *tview.List

	// app is the internal gnotes app
	app *gnotes.SelfApp
}

func newGUI() *gui {
	return &gui{
		ui:       tview.NewApplication(),
		noteList: tview.NewList(),
	}
}

const helpKeyBindings = `EXPIRIMENTAL KEYS (some do not work): F1 = Open help with less(1) command    Ctrl+F = Find/Search all notes    F2 = Back to note folder (TODO)    F3 = Search attachment names    Ctrl+D = delete note folder if empty`

func newPrimitive(text string) tview.Primitive {
	return tview.NewTextView().
		SetTextAlign(tview.AlignLeft).
		SetText(text)
}

var uilog = NewLogger()

func (self *gui) loadUI() {
	self.ui = tview.NewApplication()
	self.pages = tview.NewPages()
	self.noteList = tview.NewList()

	self.reloadNoteList()

	self.noteList.Box.SetBorder(false)
	//self.noteList.Box.SetTitle(" Your notes ")

	width, height, err := term.GetSize(0)
	if err != nil {
		panic(err)
	}

	// heightModifier is the hight modifier for the windows
	heightModifier := 4

	// Only check once if window is too small, then make sure theres space to wrap the footer.
	if len(helpKeyBindings) > width {
		height -= 1
	}

	grid := tview.NewGrid().
		SetRows(height-heightModifier). // -4 for the single line footer
		SetColumns(0).
		SetBorders(true).
		AddItem(self.noteList, 0, 0, 1, 3, 0, 0, true).
		AddItem(newPrimitive(helpKeyBindings), 1, 0, 1, 3, 0, 0, false)

	// Key mappings
	keyMapping := map[tcell.Key]func(){
		tcell.KeyCtrlF: func() {
			index := self.noteList.GetCurrentItem()
			uilog.Log("Ctrl F for current item -> %v", index)
		},
		tcell.KeyF1: func() {
			// With builtin tview window, will try with terminal editor.
			//helpPrimitive := tview.NewTextView().SetText(string(helpText))
			//helpPrimitive.SetBorder(true)

			//helpPrimitive.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
			//	self.ui.SetRoot(grid, true)
			//	return event
			//})

			//self.ui.SetRoot(helpPrimitive, true)

			// With opening configured editor
			//file, err := os.CreateTemp("", "gnotes-help-file.md")
			//if err != nil {
			//	log.Fatal(err)
			//}
			//defer os.Remove(file.Name())
			//file.Write(helpText)
			//os.Chmod(file.Name(), 0o444)

			//self.ui.Stop()
			//gnotes.OpenFileWithEditor(file.Name())
			//self.loadUI()

			// With less
			self.ui.Stop()
			openWithLess(helpText)
			self.loadUI()
		},
		tcell.KeyF2: func() {
			uilog.Log("F2")
			self.reloadNoteFolders()
		},
		tcell.KeyF3: func() {
			uilog.Log("F3")

			attachmentIndex := 0
			// Notes are always sorted, so we can rely on the first attachment since they
			// are all at the bottom
			for _, n := range self.app.Notes.GetSelected().Notes {
				if n.IsAttachment {
					break
				}
				attachmentIndex++
			}
			// Add one more since we have a top level menu item
			attachmentIndex++

			uilog.Log("Found first attachment list index at: %v", attachmentIndex)
			self.noteList.SetCurrentItem(attachmentIndex)
		},
		tcell.KeyCtrlD: func() {
			selectedIndex := self.noteList.GetCurrentItem()

			if self.currentPage != pageFolders {
				self.showWarning("Not in folder view, cannot delete notes like this.")
			}

			// Check to make sure its deletable (ie. not a menu item)
			if selectedIndex < 1 || selectedIndex >= self.noteList.GetItemCount()-1 {
				uilog.Log("Invalid index to delete: %d", selectedIndex)
				return
			}

			// TODO: confirm prompt
			// TODO: only if note folder is empty
			// TODO: only delete if in book view
			self.app.Notes.DeleteBook(selectedIndex - 1)

			self.reloadNoteFolders()
			self.noteList.SetCurrentItem(selectedIndex)
		},
	}

	keyCapture := func(event *tcell.EventKey) *tcell.EventKey {
		for k, v := range keyMapping {
			if event.Key() == k {
				v()
				return nil
			}
		}

		return event
	}

	self.pages.AddAndSwitchToPage("main_view", grid, true)

	self.ui.SetInputCapture(keyCapture)
	self.ui.SetRoot(self.pages, true)
	self.ui.EnableMouse(false)
	if err := self.ui.Run(); err != nil {
		panic(err)
	}
}

func (self *gui) showWarning(text string) {
	view := tview.NewForm().
		AddTextView("Warning", text, 40, 2, true, false).
		AddButton("OK", func() {
			self.pages.RemovePage("warning_view")
		})

	self.pages.AddAndSwitchToPage("warning_view", view, true)
}

func (self *gui) reloadNoteFolders() {
	self.currentPage = pageFolders
	self.noteList.Clear()

	self.noteList.AddItem("Create new folder", "", 'n', func() {
		uilog.Log("Creating new folder")

		form := tview.NewForm().
			AddInputField("Note folder name", "", 80, nil, nil)

		form.AddButton("Create", func() {
			nameField := form.GetFormItemByLabel("Note folder name").(*tview.InputField)
			noteBookName := nameField.GetText()

			uilog.Log("Got folder name: %s", noteBookName)

			self.app.Notes.NewBook(noteBookName)

			self.reloadNoteFolders()
			self.app.IndexNeedsUpdating = true
			self.pages.RemovePage("request_name_form")
		}).
			AddButton("Cancel", func() {
				self.ui.Stop()
				self.loadUI()
			})

		self.pages.AddAndSwitchToPage("request_name_form", form, true)
	})

	for i, book := range self.app.Notes.Books {
		self.noteList.AddItem(book.Name, fmt.Sprintf("%d notes, last modified %s", len(book.Notes), book.HRModifiedTime()), getShortcutForIndex(i), func() {
			self.app.Notes.SetSelected(self.noteList.GetCurrentItem() - 1)
			self.reloadNoteList()
		})
	}

	self.noteList.AddItem("Quit", "", 'q', func() {
		self.ui.Stop()
	})
}

func (self *gui) reloadNoteList() {
	self.currentPage = pageNotes
	self.noteList.Clear()

	self.noteList.AddItem("Create new note", "", 'n', func() {
		err := self.app.Notes.GetSelected().NewNote(
			self.app.Config.App.NoteDir,
			func() {
				err := self.openNote(len(self.app.Notes.GetSelected().Notes) - 1)
				if err != nil {
					log.Printf("Error opening new note: %s", err)
				}
			})
		if err != nil {
			log.Printf("Failed creating new note: %s", err)
		}
	})

	for i, n := range self.app.Notes.GetSelected().Notes {
		self.noteList.AddItem(n.GetTitle(self.app.Config.App.NoteDir+"/notes"), n.Info(), getShortcutForIndex(i), func() {
			err := self.openNote(self.noteList.GetCurrentItem() - 1)
			if err != nil {
				log.Printf("Failed to open note at index: %d: %s", self.noteList.GetCurrentItem()-1, err)
			}
		})
	}

	self.noteList.AddItem("Quit", "Press to exit", 'q', func() {
		// Stop the gui on quit. Will save the notes at main exit.
		self.ui.Stop()
	})
}

func (self *gui) openNote(index int) error {
	// Quit the app before opening the text editor
	self.ui.Stop()

	//if self.app.Notes.Books[self.app.Notes.LastSelected].Notes[index].IsAttachment {
	if self.app.Notes.GetSelected().Notes[index].IsAttachment {
		for attempts := 0; attempts < 10; attempts++ {
			succeed := true
			action := ""

			fmt.Printf(`What do you want to do with: %s:
	  d      - Download
	  e      - Edit name
	  delete - Delete the attachment
	  b      - Back
	: `, filepath.Base(self.app.Notes.GetSelected().Notes[index].S3Path))
			fmt.Scanln(&action)

			switch action {
			case "d":
				downloadTo := self.app.Notes.GetSelected().Notes[index].AttachmentTitle
				fmt.Printf("Downloading %s to %s...\n", downloadTo, downloadTo)

				err := self.app.Config.S3.DownloadFileFrom(
					filepath.Join(self.app.Config.S3.UserID, "notes", self.app.Notes.GetSelected().Notes[index].S3Path),
					downloadTo,
				)
				if err != nil {
					return fmt.Errorf("failed to download attachment: %s", err)
				}

				fmt.Printf("Downloaded attachment (%s) to: %s\n", self.app.Notes.GetSelected().Notes[index].Title, downloadTo)
				self.app.Notes.Sort()
				self.loadUI()
			case "e":
				return fmt.Errorf("not impmented")
			case "delete":
				fmt.Printf("Deleting %s...\n", self.app.Notes.GetSelected().Notes[index].Title)

				err := self.app.Notes.GetSelected().DeleteNote(index)
				if err != nil {
					return fmt.Errorf("failed to delete attachment: %w", err)
				}

				self.app.Notes.Sort()
				self.loadUI()
			case "b":
				self.loadUI()
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
	err := self.app.Notes.GetSelected().Notes[index].Download(self.app.Config.App.NoteDir, self.app.Config.S3)
	if err != nil {
		return err
	}

	// Run the command to open the text file with the specified editor
	editFile := filepath.Join(self.app.Config.App.NoteDir, "notes", self.app.Notes.GetSelected().Notes[index].S3Path)
	cmd := exec.Command(self.app.Config.App.Editor, editFile)
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
		err := self.app.Notes.GetSelected().DeleteNote(index)
		if err != nil {
			return fmt.Errorf("failed to delete empty note: %s", err)
		}

		self.app.Notes.Sort()
		self.loadUI()

		return nil
	}

	// Save the note if needed
	err = self.app.Notes.GetSelected().SaveNoteIndex(index)
	if err != nil {
		return fmt.Errorf("failed to save the note: %w", err)
	}

	// Resort the notes
	self.app.Notes.Sort()

	self.loadUI()

	return nil
}

func getShortcutForIndex(index int) rune {
	var s = []rune{'1', '2', '3', '4', '5', '6', '7', '8', '9'}

	if index >= len(s) {
		return 0
	}

	return s[index]
}

func openWithLess(input []byte) error {
	cmd := exec.Command("less")
	cmd.Stdin = bytes.NewReader(input)
	cmd.Stdout = os.Stdout

	err := cmd.Run()
	if err != nil {
		return err
	}

	return nil
}
