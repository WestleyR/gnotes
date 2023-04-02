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
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/google/uuid"
)

type SelfApp struct {
	// Notes
	Notes *NoteBook

	// IndexNeedsUpdating indecates if the index file needs to be uploaded
	// like if a new note was created.
	IndexNeedsUpdating bool

	// CLI opts
	CliOpts CliOpts

	Config *Config
}

type CliOpts struct {
	SkipDownload bool
	NewNote      bool
}

// This is a global config
// TODO: use this more instead of passing s3config, noteDir to functions.
var self *SelfApp

// NoteBook is the collection of all sub-categroies.
type NoteBook struct {
	Books []*Book `json:"folders"`
	//LastSelected int     `json:"last_selected"`
}

// Book contains all the notes for the Name of the sub-categroies.
type Book struct {
	// Name is the sub-dir name
	Name string `json:"name"`
	// Notes contains all the notes in the sub-dir
	Notes    []*Note `json:"notes"`
	Modified int64   `json:"modified"`
	Selected bool    `json:"selected"`
}

// Note is all the data for a specific note.
type Note struct {
	// S3Path is the path to the note on the s3 server. Also used for the local
	// path when caching. eg. "catigory/uuid-1/content"
	S3Path   string `json:"path"`
	Created  int64  `json:"created"`
	Modified int64  `json:"modified"`
	Hash     string `json:"hash"`
	Title    string `json:"title"`

	// For attachments
	IsAttachment    bool   `json:"attachment"`
	AttachmentTitle string `json:"attachment_title"`
	Size            int64  `json:"size"`
}

func InitApp(configPath string) (*SelfApp, error) {
	app := &SelfApp{}

	var err error
	app.Config, err = LoadConfig(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed loading config: %w", err)
	}

	app.Notes = &NoteBook{
		Books: []*Book{
			{
				// Default
				Name:     "Notes",
				Notes:    []*Note{},
				Selected: true,
			},
		},
	}

	app.CliOpts.SkipDownload = false
	app.CliOpts.NewNote = false

	self = app

	return app, nil
}

func (b *NoteBook) SetSelected(index int) {
	b.deselectAll()
	b.Books[index].Selected = true
}

func (b *NoteBook) GetSelected() *Book {
	for _, n := range b.Books {
		if n.Selected {
			return n
		}
	}

	// Non selected, could be a bug.
	return b.Books[0]
}

func (b *NoteBook) deselectAll() {
	for _, n := range b.Books {
		n.Selected = false
	}
}

// Download will download the note if needed based on hash.
func (n *Note) Download(noteDir string, s3Config S3Config) error {
	// Skip if theres no hash (like for a newly created note).
	if n.Hash == "" {
		return nil
	}

	noteFile := filepath.Join(noteDir, "notes", n.S3Path)

	currentHash, err := Sha1File(noteFile)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return err
	}

	if n.Hash == currentHash {
		log.Printf("Using cached note\n")
		return nil
	}

	// Download the note

	err = s3Config.DownloadFileFrom(filepath.Join(s3Config.UserID, "notes", n.S3Path), noteFile)
	if err != nil {
		return fmt.Errorf("failed to download file: %w", err)
	}

	// Verify the checksum again

	currentHash, err = Sha1File(noteFile)
	if err != nil {
		return err
	}

	if n.Hash != currentHash {
		fmt.Printf("WARNING!!! hash not the same!\n")
	}

	return nil
}

// SaveNoteIndex does the same thing as Note.Save(), but also updates the
// modified timestamp for the book.
func (b *Book) SaveNoteIndex(noteIndex int) error {
	n := b.Notes[noteIndex]

	noteFile := filepath.Join(self.Config.App.NoteDir, "notes", n.S3Path)

	currentHash, err := Sha1File(noteFile)
	if err != nil {
		return fmt.Errorf("failed to get checksum for local cached file: %w", err)
	}

	// Double check to make sure current hash is not empty
	if currentHash != "" && n.Hash != currentHash {
		// Upload the note that changed
		err := self.Config.S3.UploadFile(
			noteFile,
			filepath.Join(self.Config.S3.UserID, "notes", n.S3Path),
		)
		if err != nil {
			return err
		}

		// After uploading, update the hash tracker
		n.Hash = currentHash
		n.Changed()

		self.IndexNeedsUpdating = true
		b.Changed(noteIndex)

		return nil
	}
	log.Printf("Not uploading note since it has not changed")

	return nil
}

// Save will save a note to the s3 server. Should be called after editing. Will
// NOT update the note index file. Does not update/upload if the checksum did not
// change.
// Depercated: use Book.SaveNoteIndex() instead (MAYBE...)
func (n *Note) Save() error {
	noteFile := filepath.Join(self.Config.App.NoteDir, "notes", n.S3Path)

	currentHash, err := Sha1File(noteFile)
	if err != nil {
		return fmt.Errorf("failed to get checksum for local cached file: %w", err)
	}

	// Double check to make sure current hash is not empty
	if currentHash != "" && n.Hash != currentHash {
		// Upload the note that changed
		err := self.Config.S3.UploadFile(
			noteFile,
			filepath.Join(self.Config.S3.UserID, "notes", n.S3Path),
		)
		if err != nil {
			return err
		}

		// After uploading, update the hash tracker
		n.Hash = currentHash
		n.Changed()

		self.IndexNeedsUpdating = true
		return nil
	}
	log.Printf("Not uploading note since it has not changed")

	return nil
}

func (n *NoteBook) DeleteBook(index int) error {
	n.Books = append(n.Books[:index], n.Books[index+1:]...)
	self.IndexNeedsUpdating = true

	return nil
}

// DeleteNote will delete a specific note. Will delete the note from s3 imetitly,
// and reupload the index files.
func (b *Book) DeleteNote(noteIndex int) error {
	notePath := filepath.Dir(filepath.Join(self.Config.App.NoteDir, "notes", b.Notes[noteIndex].S3Path))

	log.Printf("Removing/deleting note: %s", notePath)

	err := os.RemoveAll(notePath)
	if err != nil {
		return fmt.Errorf("failed to delete note: %w", err)
	}

	// Delete it from s3
	err = self.Config.S3.Delete(filepath.Join(self.Config.S3.UserID, "notes", b.Notes[noteIndex].S3Path))
	if err != nil {
		return fmt.Errorf("failed to delete note from s3: %w", err)
	}

	b.Notes = append(b.Notes[:noteIndex], b.Notes[noteIndex+1:]...)

	self.IndexNeedsUpdating = true
	b.Changed(-1)

	return nil
}

// DeleteNote will delete a specific note. Will delete the note from s3 imetitly,
// and reupload the index files.
// Depercated: use Book.DeleteNote()
func (self *SelfApp) DeleteNote(bookIndex, noteIndex int) error {
	notePath := filepath.Dir(filepath.Join(self.Config.App.NoteDir, "notes", self.Notes.Books[bookIndex].Notes[noteIndex].S3Path))

	log.Printf("Removing/deleting note: %s", notePath)

	err := os.RemoveAll(notePath)
	if err != nil {
		return fmt.Errorf("failed to delete note: %w", err)
	}

	// Delete it from s3
	err = self.Config.S3.Delete(filepath.Join(self.Config.S3.UserID, "notes", self.Notes.Books[bookIndex].Notes[noteIndex].S3Path))
	if err != nil {
		return fmt.Errorf("failed to delete note from s3: %w", err)
	}

	self.Notes.Books[bookIndex].Notes = append(self.Notes.Books[bookIndex].Notes[:noteIndex], self.Notes.Books[bookIndex].Notes[noteIndex+1:]...)

	self.IndexNeedsUpdating = true

	return nil
}

// GetTitle returns a title for a note. Requires the local cache path.
// Should be ~/.cache/gnotes/notes
func (n *Note) GetTitle(noteDir string) string {
	if n.IsAttachment {
		return "Attachment: " + n.AttachmentTitle
	}

	notePath := filepath.Join(noteDir, n.S3Path)

	r, err := os.Open(notePath)
	if err != nil {
		return n.Title
	}
	defer r.Close()

	head := make([]byte, 64)
	l, err := r.Read(head)
	if err != nil {
		return "error: " + err.Error()
	}

	if string(head[:]) == "" {
		// Note should be removed if its empty
		return "empty"
	}

	title := strings.ReplaceAll(string(head[:l]), "\n", " ")

	if title == "" {
		return n.Title
	}

	n.Title = title

	return n.Title
}

func (b *Book) HRModifiedTime() string {
	c := time.Unix(b.Modified, 0)
	return c.Format("2006-01-02 07:05:45PM")
}

func (a *Note) Info() string {
	if a.IsAttachment {
		c := time.Unix(a.Created, 0)
		return fmt.Sprintf("Created on %s. %s", c.Format("2006-01-02"), formatBytes(a.Size))
	}

	c := time.Unix(a.Created, 0)
	m := time.Unix(a.Modified, 0)

	return fmt.Sprintf("Created on %s. last modified on %s", c.Format("2006-01-02"), m.Format("2006-01-02"))
}

// Changed updates the modified timestamp for the book and note.
func (b *Book) Changed(noteIndex int) {
	b.Modified = time.Now().Unix()

	if noteIndex != -1 {
		b.Notes[noteIndex].Modified = b.Modified
	}
}

// Changed will update the modified date for a note.
// Depercated: use Book.Changed()
func (n *Note) Changed() {
	n.Modified = time.Now().Unix()
}

func (book *Book) NewAttachment(noteDir, path string) error {
	src, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("failed to open file for new attachment: %s", err)
	}
	defer src.Close()

	stat, err := src.Stat()
	if err != nil {
		return fmt.Errorf("failed to stat file: %s", err)
	}

	uuidP := uuid.NewString()
	createdTime := time.Now().Unix()

	// TODO: dont read the whole file into memory
	//fileContents, err := os.ReadFile(path)
	//if err != nil {
	//	return fmt.Errorf("failed to read file: %s", err)
	//}

	// TODO: add back later...
	//	if util.IsText(fileContents) {
	//		// Its a text file, so it needs to be added to notes, not attachments
	//		// TODO: add flag to disable this
	//		// TODO: add size limit to disable this
	//		log.Printf("Adding as not since it seems to be a text file")
	//		return book.NewNoteWithContentsOfFile(noteDir, path, nil)
	//	}

	newAttachment := &Note{
		S3Path:          filepath.Join(book.Name, uuidP, filepath.Base(path)),
		IsAttachment:    true,
		AttachmentTitle: filepath.Base(path),
		Created:         createdTime,
		Size:            stat.Size(),
		Hash:            "",
	}

	notePath := filepath.Join(noteDir, "notes", newAttachment.S3Path)

	err = os.MkdirAll(filepath.Dir(notePath), 0755)
	if err != nil {
		return fmt.Errorf("failed to create new note dir: %w", err)
	}

	err = copyFileContents(path, notePath)
	if err != nil {
		return fmt.Errorf("failed to copy file: %w", err)
	}

	book.Notes = append(book.Notes, newAttachment)

	self.IndexNeedsUpdating = true

	err = newAttachment.Save()
	if err != nil {
		return fmt.Errorf("failed to upload new attachment: %w", err)
	}

	// Manually update the modified timestamp for the book
	book.Changed(-1)

	return nil
}

func (book *Book) NewNoteWithContentsOfFile(noteDir, path string, completion func()) error {
	createdTime := time.Now().Unix()

	uuidP := uuid.NewString()

	newNote := &Note{
		S3Path:   filepath.Join(book.Name, uuidP, "content"),
		Created:  createdTime,
		Modified: createdTime,
		Hash:     "",
	}

	notePath := filepath.Join(noteDir, "notes", newNote.S3Path)

	err := os.MkdirAll(filepath.Dir(notePath), 0755)
	if err != nil {
		return fmt.Errorf("failed to create new note dir: %w", err)
	}

	err = os.WriteFile(notePath, []byte{}, 0664)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}

	// Copy the data (if theres any)
	// TODO: use copy func
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

	book.Notes = append(book.Notes, newNote)

	fmt.Printf("NEW NOTES: %+v\n", book)

	self.IndexNeedsUpdating = true

	// Open the new note
	if completion != nil {
		completion()
	}

	return nil
}

func (book *Book) NewNote(noteDir string, completion func()) error {
	return book.NewNoteWithContentsOfFile(noteDir, "", completion)
}

func (noteBook *NoteBook) NewBook(name string) error {
	if name == "" {
		return fmt.Errorf("name cannot be empty")
	}

	newBook := &Book{
		Name:     name,
		Notes:    []*Note{},
		Selected: true,
	}
	newBook.Changed(-1)

	noteBook.deselectAll()

	noteBook.Books = append(noteBook.Books, newBook)

	return nil
}

func (n *NoteBook) Sort() {
	//	sort.Slice(n.Books, func(i, j int) bool {
	//		return n.Books[i].Modified > n.Books[j].Modified
	//	})

	sorting := true
	for sorting {
		sorting = false
		for i := 0; i < len(n.Books)-1; i++ {
			if n.Books[i].Modified < n.Books[i+1].Modified {
				n.Books[i], n.Books[i+1] = n.Books[i+1], n.Books[i]
				sorting = true
			}
		}
	}

	for _, b := range n.Books {
		b.Sort()
	}
}

func (n *Book) Sort() {
	// First, put all attachments at the bottom,
	// then, sort all notes by last modified
	// finally, sort all attachments by date created,

	numNotes := 0

	for _, a := range n.Notes {
		if !a.IsAttachment {
			numNotes++
		}
	}

	// Sort the note type
	sort.Slice(n.Notes, func(i, j int) bool {
		if boolInt(n.Notes[i].IsAttachment) < boolInt(n.Notes[j].IsAttachment) {
			return true
		}
		return false
	})

	// Sort the notes
	sorting := true
	for sorting {
		sorting = false
		for i := 0; i < numNotes-1; i++ {
			if n.Notes[i].Modified < n.Notes[i+1].Modified {
				n.Notes[i], n.Notes[i+1] = n.Notes[i+1], n.Notes[i]
				sorting = true
			}
		}
	}

	// Sort the attachments
	sorting = true
	for sorting {
		sorting = false
		for i := numNotes; i < len(n.Notes)-1; i++ {
			// By name
			// if !sortorder.NaturalLess(filepath.Base(n.Notes[i].S3Path), filepath.Base(n.Notes[i+1].S3Path)) {
			if n.Notes[i].Created < n.Notes[i+1].Created {
				n.Notes[i], n.Notes[i+1] = n.Notes[i+1], n.Notes[i]
				sorting = true
			}
		}
	}
}

// boolInt will convert a bool to int. Used for sorting.
func boolInt(b bool) int {
	if b {
		return 1
	}
	return 0
}

func copyFileContents(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	if _, err = io.Copy(out, in); err != nil {
		return err
	}

	if err := out.Sync(); err != nil {
		return err
	}

	return out.Close()
}
