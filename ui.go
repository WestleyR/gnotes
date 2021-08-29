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

	"github.com/rivo/tview"
)

func getShortcutForIndex(index int) rune {
	var s = []rune{'1', '2', '3', '4', '5', '6', '7', '8', '9'}

	if index >= len(s) {
		return 0
	}

	return s[index]
}

func (self *selfApp) loadUI() {
	self.app = tview.NewApplication()
	self.noteList = tview.NewList()

	self.reloadUI()

	self.noteList.Box.SetBorder(true)
	self.noteList.Box.SetTitle(" Your notes ")
	self.app.SetRoot(self.noteList, true)
	self.app.EnableMouse(true)

	if err := self.app.Run(); err != nil {
		panic(err)
	}
}

func (self *selfApp) reloadUI() {
	self.noteList.Clear()

	self.noteList.AddItem("Create new note", "", 'n', func() {
		err := self.newNote()
		if err != nil {
			log.Fatalf("Failed creating new note: %s\n", err)
		}
	})

	for i, n := range self.notes {
		self.noteList.AddItem(getTitleForNote(n.Content), getSubContentForNote(n), getShortcutForIndex(i), func() {
			err := self.openNote(self.noteList.GetCurrentItem() - 1)
			if err != nil {
				log.Fatalf("Failed to open note at index: %d: %s\n", self.noteList.GetCurrentItem()-1, err)
			}
		})
	}

	self.noteList.AddItem("About", "About gnotes", 'a', func() {
		self.app.Stop()
		fmt.Println(`
gnotes - S3 syncing notes



# gnotes - Terminal based S3 syncing note app

## Installation

_todo..._

  $ go install ...

Or

  $ git clone ... \
    cd gnotes/    \
    go build      \
    cp gnotes ~/.local/bin  # or your perfured path

## Setting up gnotes for your S3 server

The user config file for gnotes is located in: _(subject to change)_

  ${HOME}/.config/wst.gnotes/config.ini

_**NOTE:** You will need to restart the app after makeing changes to the config file._

And as an example, for dream host, it should look like this:

  [settings]
  editor = vim
  
  [s3]
  active = true
  bucket = gnotes
  endpoint = https://objects-us-east-1.dream.io
  region = us-east-1
  savefile = gnotes.json
  accesskey = <ACCESS_KEY>
  secretkey = <SECRET_KEY>`)
	})

	self.noteList.AddItem("Quit", "Press to exit", 'q', func() {
		self.app.Stop()
		if self.notesChanged {
			err := self.saveNotes()
			if err != nil {
				log.Fatalf("Failed to upload notes: %s\n", err)
			}
		} else {
			log.Printf("Not uploading notes since it did not change\n")
		}
	})
}
