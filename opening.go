//
//  opening.go - https://github.com/WestleyR/gnotes
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
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/h2non/filetype"
)

func (self *SelfApp) LoadNotes() error {
	downloadNote := filepath.Join(self.Config.App.NoteDir, "gnotes.tar.gz")

	// Only download from s3 if active = true
	if self.Config.S3.Active {
		if !self.CliOpts.SkipDownload {
			err := self.Config.S3.S3DownloadFile(downloadNote)
			if err != nil {
				return fmt.Errorf("error downloading the save file: %s: %s\n", self.Config.S3.File, err)
			}
		}
	}

	// New scope so we dont keep it in memory for long
	{
		// Open the blob into memory
		downloadedBytes, err := os.ReadFile(downloadNote)
		if err != nil {
			if !self.CliOpts.NewNote {
				return fmt.Errorf("failed to open downloaded file: %s", err)
			}
			self.Notes.Sort()
			//sortByModTime(self.Notes.Notes)
			return nil
		}

		// Check if the file is a tar, this may be existing notes and encryption is addon
		kind, _ := filetype.Match(downloadedBytes)
		if kind != filetype.NewType("gz", "application/gzip") {
			// Decrypt if its enabled
			downloadedBytes, err = self.Config.Crypt.DecryptIfEnabled(downloadedBytes)
			if err != nil {
				return fmt.Errorf("failed to decrypt data: %s", err)
			}
		} else {
			log.Printf("WARNING: Notes are not encrypted")
		}

		// Now untar the notes
		err = untar(downloadedBytes, filepath.Join(self.Config.App.NoteDir, "gnotes"))
		if err != nil {
			if self.CliOpts.NewNote {
				return fmt.Errorf("failed to untar: %s", err)
			}
		}
	}

	// Now read the downloaded file
	downloadedJson, err := os.ReadFile(filepath.Join(self.Config.App.NoteDir, "gnotes", "notes/gnotes.json"))
	if err != nil {
		return fmt.Errorf("failed to read json: %s", err)
	}

	err = json.Unmarshal(downloadedJson, &self.Notes)
	if err != nil {
		panic(err)
	}

	// Now sort the notes by mod time
	self.Notes.Sort()
	//sortByModTime(self.Notes.Notes)

	return nil
}

func (self *SelfApp) SaveNotes() error {
	var jsonData []byte
	jsonData, err := json.Marshal(self.Notes)
	if err != nil {
		return err
	}

	err = os.WriteFile(filepath.Join(self.Config.App.NoteDir, "gnotes", "notes/gnotes.json"), jsonData, 0644)
	if err != nil {
		panic(err)
	}

	// New scope to minimize time of large memory blobs
	{
		tarData, err := tarCompress(filepath.Join(self.Config.App.NoteDir, "gnotes"))
		if err != nil {
			return fmt.Errorf("failed to tar gzip compress: %s", err)
		}

		// Write the tar to the config dir as backup
		err = writeNewFile(filepath.Join(self.Config.App.NoteDir, "gnotes.backup.tar.gz"), tarData)
		if err != nil {
			return fmt.Errorf("failed to write backup file: %s", err)
		}

		tarData, err = self.Config.Crypt.EncryptIfEnabled(tarData)
		if err != nil {
			return fmt.Errorf("failed to encrypt file: %s", err)
		}

		err = writeNewFile(filepath.Join(self.Config.App.NoteDir, "gnotes.tar.gz"), tarData)
		if err != nil {
			return fmt.Errorf("failed to write note file: %s", err)
		}
	}

	if self.NotesChanged {
		if self.Config.S3.Active {
			err = self.Config.S3.S3UploadFile(filepath.Join(self.Config.App.NoteDir, "gnotes.tar.gz"))
			if err != nil {
				return fmt.Errorf("failed to upload file: %s", err)
			}
		}
	}

	// Remove the open note directory to avoid note conflicts
	err = os.RemoveAll(filepath.Join(self.Config.App.NoteDir, "gnotes"))
	if err != nil {
		return fmt.Errorf("failed to remove dest dir: %s", err)
	}

	return nil
}

func writeNewFile(path string, data []byte) error {
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		return fmt.Errorf("failed to open file for writing: %s", err)
	}

	_, err = f.Write(data)
	if err != nil {
		return fmt.Errorf("failed to write data: %s", err)
	}

	return nil
}
