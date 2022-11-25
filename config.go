//
//  config.go - https://github.com/WestleyR/gnotes
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
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/wildwest-productions/goini"
)

const appID = "wst.gnotes"

type Config struct {
	Account accountConfig `ini:"account"`
	App     appSettings   `ini:"settings"`
	S3      S3Config      `ini:"s3"`
	Crypt   cryptConfig   `ini:"encrypt"`
}

// TODO: impment me
type accountConfig struct {
	User  string `ini:"user"`
	Token string `ini:"token"`
}

type cryptConfig struct {
	Enable bool   `ini:"enable"`
	Key    string `ini:"key"`
}

type appSettings struct {
	Editor  string `ini:"editor"`
	NoteDir string `ini:"notes_dir"`
}

type S3Config struct {
	Active    bool   `ini:"active"`
	Bucket    string `ini:"bucket"`
	Endpoint  string `ini:"endpoint"`
	Region    string `ini:"region"`
	File      string `ini:"file"`
	AccessKey string `ini:"accesskey"`
	SecretKey string `ini:"secretkey"`
}

func LoadConfig(configFile string) *Config {
	conf := &Config{
		Account: accountConfig{
			User:  "",
			Token: "",
		},
		App: appSettings{
			Editor:  "vim",
			NoteDir: "/tmp",
		},
		S3: S3Config{
			Active:    false,
			Bucket:    "",
			Endpoint:  "",
			Region:    "",
			File:      "change-me",
			AccessKey: "",
			SecretKey: "",
		},
		Crypt: cryptConfig{
			Enable: false,
			Key:    "",
		},
	}

	iniBytes, err := ioutil.ReadFile(configFile)
	if err != nil {
		log.Printf("Error opening: %s", configFile)
		return conf
	}

	err = goini.Unmarshal(iniBytes, &conf)
	if err != nil {
		log.Printf("Error unmarshaling file: %s", err)
	}

	// Replace the ${HOME} string in the NoteDir if needed
	home, err := os.UserHomeDir()
	if err != nil {
		log.Printf("Error getting home dir: %s", err)
		return conf
	}
	conf.App.NoteDir = strings.ReplaceAll(conf.App.NoteDir, "${HOME}", home)

	return conf
}

func GetFileFromConfig(file string) string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}

	configFile := filepath.Join(home, ".config", appID)

	if _, err = os.Stat(configFile); os.IsNotExist(err) {
		err = os.MkdirAll(configFile, 0700)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to create: %s: %s\n", configFile, err)
			return ""
		}
	}

	configFile = filepath.Join(configFile, file)

	f, err := os.OpenFile(configFile, os.O_CREATE, 0644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open: %s: %s\n", configFile, err)
		return ""
	}
	f.Close()

	return configFile
}
