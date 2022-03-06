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

package gnotes

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
	"golang.org/x/tools/godoc/util"
)

type SelfApp struct {
	// Notes
	Notes *NoteSpec

	// On exit, dont upload if notes did not change
	NotesChanged bool

	// CLI opts
	CliOpts CliOpts

	Config *Config
}

type CliOpts struct {
	SkipDownload bool
	NewNote      bool
}

type NoteSpec struct {
	Notes []NoteInfo `json:"notes"`
}

type NoteInfo struct {
	Dir      string `json:"dir"`
	File     string `json:"file"`
	Created  int64  `json:"created"`
	Modified int64  `json:"modified"`
	Hash     string `json:"hash"`

	// For attachments
	IsAttachment    bool   `json:"attachment"`
	AttachmentTitle string `json:"attachment_title"`
	Size            int64  `json:"size"`
	Type            string `json:"type"` // TODO:
}

func InitApp(configPath string) (*SelfApp, error) {
	app := &SelfApp{}

	app.NotesChanged = false

	app.Config = LoadConfig()

	app.Notes = new(NoteSpec)

	app.CliOpts.SkipDownload = false
	app.CliOpts.NewNote = false

	return app, nil
}

func (n NoteInfo) Title(noteDir string) string {
	if n.IsAttachment {
		return "Attachment: " + n.AttachmentTitle
	}

	notePath := filepath.Join(noteDir, "gnotes", n.File)

	r, err := os.Open(notePath)
	if err != nil {
		return "[error]"
	}
	defer r.Close()

	head := make([]byte, 64)
	_, err = r.Read(head)
	if err != nil {
		return "[error]"
	}

	if string(head[:]) == "" {
		// Note should be removed if its empty
		return "[empty]"
	}

	return strings.ReplaceAll(string(head[:]), "\n", " ")
}

func (a NoteInfo) Info() string {
	if a.IsAttachment {
		c := time.Unix(a.Created, 0)
		return fmt.Sprintf("Type %s, created on %s. %s", a.Type, c.Format("2006-01-02"), formatBytes(a.Size))
	}

	c := time.Unix(a.Created, 0)
	m := time.Unix(a.Modified, 0)

	return fmt.Sprintf("Created on %s. last modified on %s", c.Format("2006-01-02"), m.Format("2006-01-02"))
}

func (self *SelfApp) NewAttachment(path string) error {
	src, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("failed to open file for new attachment: %s", err)
	}
	defer src.Close()

	stat, err := src.Stat()
	if err != nil {
		return fmt.Errorf("failed to stat file: %s", err)
	}

	u := uuid.NewString()
	createdTime := time.Now().Unix()

	// TODO: dont read the whole file into memory
	fileContents, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read file: %s", err)
	}

	if util.IsText(fileContents) {
		// Its a text file, so it needs to be added to notes, not attachments
		// TODO: add flag to disable this
		// TODO: add size limit to disable this
		log.Printf("Adding as not since it seems to be a text file")
		return self.NewNoteWithContentsOfFile(path, nil)
	}

	newAttachment := NoteInfo{
		File:            u,
		IsAttachment:    true,
		AttachmentTitle: filepath.Base(path),
		Created:         createdTime,
		Type:            "unknown file",
		Size:            stat.Size(),
		Hash:            "na",
	}

	err = self.Config.S3.S3UploadFileTo(path, u)
	if err != nil {
		return fmt.Errorf("failed to upload file: %s", err)
	}

	self.Notes.Notes = append(self.Notes.Notes, newAttachment)
	self.NotesChanged = true

	return nil
}

func (self *SelfApp) NewNoteWithContentsOfFile(path string, completion func()) error {
	createdTime := time.Now().Unix()

	uuidP := uuid.NewString()

	newNote := NoteInfo{
		Dir:      "notes/" + uuidP,
		File:     "notes/" + uuidP + "/content",
		Created:  createdTime,
		Modified: createdTime,
		// TODO: Need to generate a hash
		Hash: "na",
	}

	notePath := filepath.Join(self.Config.App.NoteDir, "gnotes", newNote.File)

	err := os.MkdirAll(filepath.Join(self.Config.App.NoteDir, "gnotes", newNote.Dir), 0755)
	if err != nil {
		return err
	}

	// Copy the data (if theres any)
	if path != "" {
		in, err := os.Open(path)
		if err != nil {
			return fmt.Errorf("failed to open contents of file: %s", err)
		}
		defer in.Close()

		out, err := os.Create(notePath)
		if err != nil {
			return fmt.Errorf("failed to create file: %s", err)
		}

		defer out.Close()

		_, err = io.Copy(out, in)
		if err != nil {
			return fmt.Errorf("failed to copy: %s", err)
		}

		err = out.Sync()
		if err != nil {
			return fmt.Errorf("failed to sync: %s", err)
		}
	}

	self.Notes.Notes = append(self.Notes.Notes, newNote)
	self.NotesChanged = true

	// If the ui is not loaded, then just return and dont open the note
	//	if !self.uiLoaded {
	//		log.Printf("not adding item to ui")
	//		return nil
	//	}

	// Append the new note
	//	self.noteList.AddItem("hello", "[new_note]", getShortcutForIndex(len(self.notes.Notes)-1), func() {
	//		err := self.openNote(self.noteList.GetCurrentItem() - 1)
	//		if err != nil {
	//			log.Fatalf("Failed opening note index: %d: %s\n", self.noteList.GetCurrentItem()-1, err)
	//		}
	//	})

	// Open the new note
	if completion != nil {
		completion()
	}
	//return self.openNote(len(self.Notes.Notes) - 1)
	return nil
}

func (self *SelfApp) NewNote(completion func()) error {
	return self.NewNoteWithContentsOfFile("", completion)
}

//func (self *SelfApp) openNote(index int) error {
//	// Quit the app before opening the text editor
//	//self.app.Stop()
//
//	if self.notes.Notes[index].IsAttachment {
//		action := ""
//
//		fmt.Printf(`Do you want to:
//  d      - Download
//  e      - Edit name
//  delete - Delete the attachment
//  b      - Back
//: `)
//		fmt.Scanln(&action)
//
//		switch action {
//		case "d":
//			downloadTo := self.notes.Notes[index].AttachmentTitle
//			fmt.Printf("Downloading %s to %s...\n", downloadTo, downloadTo)
//
//			err := self.config.S3.s3DownloadFileFrom(self.notes.Notes[index].File, downloadTo)
//			if err != nil {
//				return fmt.Errorf("failed to download attachment: %s", err)
//			}
//
//			fmt.Printf("Downloaded attachment (%s) to: %s\n", self.notes.Notes[index].Title(""), downloadTo)
//		case "e":
//			return fmt.Errorf("not impmented")
//		case "delete":
//			fmt.Printf("Deleting %s...\n", self.notes.Notes[index].Title(""))
//			err := self.config.S3.Delete(self.notes.Notes[index].File)
//			if err != nil {
//				return fmt.Errorf("failed to delete file from s3: %s", err)
//			}
//
//			self.Notes.Notes = append(self.notes.Notes[:index], self.notes.Notes[index+1:]...)
//			self.NotesChanged = true
//
//			sortByModTime(self.notes.Notes)
//			//self.loadUI()
//		case "b":
//			//self.loadUI()
//		default:
//			return fmt.Errorf("unknown input: %s", action)
//		}
//
//		return nil
//	}
//
//	// Run the command to open the text file with the specified editor
//	editFile := filepath.Join(self.config.App.NoteDir, "gnotes", self.notes.Notes[index].File)
//	cmd := exec.Command(self.config.App.Editor, editFile)
//	cmd.Stdout = os.Stdout
//	cmd.Stdin = os.Stdin
//	cmd.Stderr = os.Stderr
//
//	err := cmd.Run()
//	if err != nil {
//		return err
//	}
//
//	// Check if the file is empty
//	r, err := os.Open(editFile)
//	if err != nil {
//		return fmt.Errorf("failed to open file: %s", err)
//	}
//	defer r.Close()
//
//	b, err := os.ReadFile(editFile)
//	if err != nil {
//		return fmt.Errorf("failed to read file: %s", err)
//	}
//
//	if string(b) == "" {
//		// Note is empty, so delete it
//		err := os.RemoveAll(filepath.Join(self.config.App.NoteDir, "gnotes", self.notes.Notes[index].Dir))
//		if err != nil {
//			return fmt.Errorf("failed to delete empty note: %s", err)
//		}
//
//		self.notes.Notes = append(self.notes.Notes[:index], self.notes.Notes[index+1:]...)
//
//		self.notesChanged = true
//
//		sortByModTime(self.notes.Notes)
//		//self.loadUI()
//
//		return nil
//	}
//
//	newSha, err := Sha1File(editFile)
//	if err != nil {
//		return fmt.Errorf("failed to sha file: %s", err)
//	}
//
//	if self.notes.Notes[index].Hash != newSha {
//		// Note changed
//		self.notes.Notes[index].Hash = newSha
//		self.notes.Notes[index].Modified = time.Now().Unix()
//		self.notesChanged = true
//	}
//
//	// Resort the notes
//	sortByModTime(self.notes.Notes)
//
//	//self.loadUI()
//
//	return nil
//}

func (n *NoteSpec) Sort() {
	sorting := true

	for sorting {
		sorting = false
		for i := 0; i < len(n.Notes)-1; i++ {
			if n.Notes[i].Modified < n.Notes[i+1].Modified {
				tmp := n.Notes[i]
				n.Notes[i] = n.Notes[i+1]
				n.Notes[i+1] = tmp
				sorting = true
			}
		}
	}
}

//func sortByModTime(notes []NoteInfo) {
//	sorting := true
//
//	for sorting {
//		sorting = false
//		for i := 0; i < len(notes)-1; i++ {
//			if notes[i].Modified < notes[i+1].Modified {
//				tmp := notes[i]
//				notes[i] = notes[i+1]
//				notes[i+1] = tmp
//				sorting = true
//			}
//		}
//	}
//}
