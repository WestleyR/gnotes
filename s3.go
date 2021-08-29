//
//  s3.go - https://github.com/WestleyR/gnotes
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
	"os"
	"path/filepath"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

func s3UploadFile(filepath, name string) error {
	bucket := aws.String(getS3Bucket())
	uploadName := aws.String(name)

	s3Config := &aws.Config{
		//Credentials: credentials.NewStaticCredentials(s3AccessKey, s3SecretKey, ""),
		Credentials: credentials.NewStaticCredentials(getS3AccessKey(), getS3SecretKey(), ""),
		//Endpoint:         aws.String("https://objects-us-east-1.dream.io"),
		Endpoint: aws.String(getS3Endpoint()),
		//Region:           aws.String("us-east-1"),
		Region:           aws.String(getS3Region()),
		DisableSSL:       aws.Bool(true),
		S3ForcePathStyle: aws.Bool(true),
	}
	// newSession := session.New(s3Config) // deprecated, trying NewSession
	newSession, err := session.NewSession(s3Config)
	if err != nil {
		return fmt.Errorf("error creating session: %s", err)
	}

	s3Client := s3.New(newSession)

	cparams := &s3.CreateBucketInput{
		Bucket: bucket, // Required
	}

	// Create a new bucket using the CreateBucket call.
	_, err = s3Client.CreateBucket(cparams)
	if err != nil {
		return err
	}

	// Read the file to upload
	f, err := os.ReadFile(filepath)
	if err != nil {
		return fmt.Errorf("failed to read file to upload: %s", err)
	}

	// Upload a new object "testobject" with the string "Hello World!" to our "newbucket".
	_, err = s3Client.PutObject(&s3.PutObjectInput{
		Body:   strings.NewReader(string(f)),
		Bucket: bucket,
		Key:    uploadName,
	})
	if err != nil {
		return fmt.Errorf("failed to upload data to %s/%s, %s", *bucket, *uploadName, err.Error())
	}
	fmt.Printf("Successfully created bucket %s and uploaded data with key %s\n", *bucket, *uploadName)

	return nil
}

func s3DownloadFile(s3Path, endPath string) error {
	bucket := aws.String(getS3Bucket())

	s3Config := &aws.Config{
		Credentials:      credentials.NewStaticCredentials(getS3AccessKey(), getS3SecretKey(), ""),
		Endpoint:         aws.String(getS3Endpoint()),
		Region:           aws.String(getS3Region()),
		DisableSSL:       aws.Bool(true),
		S3ForcePathStyle: aws.Bool(true),
	}
	// newSession := session.New(s3Config) // deprecated, trying NewSession
	newSession, err := session.NewSession(s3Config)
	if err != nil {
		return fmt.Errorf("error creating session: %s", err)
	}

	//	s3Client := s3.New(newSession)
	//	cparams := &s3.CreateBucketInput{
	//		Bucket: bucket, // Required
	//	}
	//
	//	// Create a new bucket using the CreateBucket call.
	//	_, err := s3Client.CreateBucket(cparams)
	//	if err != nil {
	//		return err
	//	}

	// Create the base dir if it does not exist
	baseDir := filepath.Dir(endPath)
	if _, err := os.Stat(baseDir); os.IsNotExist(err) {
		err := os.MkdirAll(baseDir, 0700)
		if err != nil {
			return err
		}
	}

	file, err := os.Create(endPath)
	if err != nil {
		return fmt.Errorf("failed to create output file: %s", err)
	}
	defer file.Close()

	downloader := s3manager.NewDownloader(newSession)
	numBytes, err := downloader.Download(file,
		&s3.GetObjectInput{
			Bucket: bucket,
			Key:    aws.String(s3Path),
		})
	if err != nil {
		return fmt.Errorf("failed to download file: %s", err)
	}

	fmt.Println("Downloaded file", file.Name(), numBytes, "bytes")

	return nil
}
