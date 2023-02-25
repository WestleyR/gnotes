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
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

func (self *SelfApp) downloadIndexIfNeeded() error {
	noteSha := filepath.Join(self.Config.App.NoteDir, "notes", "index.json.sha256")

	err := self.Config.S3.DownloadFileFrom(
		filepath.Join(self.Config.S3.UserID, "notes", "index.json.sha256"),
		noteSha,
	)
	if err != nil {
		return fmt.Errorf("failed to download file to: %s: %w", noteSha, err)
	}

	oldSha, err := Sha1File(filepath.Join(self.Config.App.NoteDir, "notes", "index.json"))
	if err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("failed to read index.json file: %w", err)
		}
	}

	newSha, err := os.ReadFile(noteSha)
	if err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("failed to read note sha file: %w", err)
		}
	}

	if string(newSha) != oldSha || oldSha == "" {
		log.Printf("Downloading note index...\n")
		noteIndex := filepath.Join(self.Config.App.NoteDir, "notes", "index.json")

		err := self.Config.S3.DownloadFileFrom(
			filepath.Join(self.Config.S3.UserID, "notes", "index.json"),
			noteIndex,
		)
		if err != nil {
			return err
		}

		// Verify sha after download
		currentSha, err := Sha1File(noteIndex)
		if err != nil {
			return fmt.Errorf("failed to get current sha: %w", err)
		}
		if currentSha != string(newSha) {
			fmt.Printf("WARNING!!! Sha does not match!\n")
			// TODO: probaly abort...
		}
	}

	return nil
}

func (self *SelfApp) LoadNotes() error {
	if self.CliOpts.NewNote {
		log.Printf("skipping json reading since new note is specified")
		return nil
	}

	err := self.downloadIndexIfNeeded()
	if err != nil {
		return err
	}

	// Now read the downloaded file
	downloadedJson, err := os.ReadFile(filepath.Join(self.Config.App.NoteDir, "notes", "index.json"))
	if err != nil {
		return fmt.Errorf("failed to read json: %s", err)
	}

	err = json.Unmarshal(downloadedJson, &self.Notes)
	if err != nil {
		return fmt.Errorf("failed to unmarshal json into notes: %w", err)
	}

	// Now sort the notes by mod time
	self.Notes.Sort()

	return nil
}

func (self *SelfApp) SaveIndexFile() error {
	if !self.IndexNeedsUpdating {
		log.Printf("Not uploading any changes\n")
		return nil
	}

	// Create, and upload the index.json.sha256 and index.json

	noteIndex := filepath.Join(self.Config.App.NoteDir, "notes", "index.json")

	b, err := json.Marshal(self.Notes)
	if err != nil {
		return err
	}
	err = os.WriteFile(noteIndex, b, 0664)
	if err != nil {
		return err
	}

	err = self.Config.S3.UploadFile(
		noteIndex,
		filepath.Join(self.Config.S3.UserID, "notes", "index.json"),
	)
	if err != nil {
		return err
	}

	// Rewrite the sha file
	sha, err := Sha1File(noteIndex)
	if err != nil {
		return err
	}

	err = os.WriteFile(noteIndex+".sha256", []byte(sha), 0664)
	if err != nil {
		return err
	}

	err = self.Config.S3.UploadFile(
		noteIndex+".sha256",
		filepath.Join(self.Config.S3.UserID, "notes", "index.json.sha256"),
	)
	if err != nil {
		return err
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
