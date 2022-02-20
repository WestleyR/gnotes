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

package gnotes

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

// Default upload to save the notes.
func (c S3Config) S3UploadFile(filepath string) error {
	return c.S3UploadFileTo(filepath, c.File)
}

func (c S3Config) S3UploadFileTo(filepath, to string) error {
	s3Config := &aws.Config{
		Credentials:      credentials.NewStaticCredentials(c.AccessKey, c.SecretKey, ""),
		Endpoint:         aws.String(c.Endpoint),
		Region:           aws.String(c.Region),
		DisableSSL:       aws.Bool(true),
		S3ForcePathStyle: aws.Bool(true),
	}

	newSession, err := session.NewSession(s3Config)
	if err != nil {
		return fmt.Errorf("error creating session: %s", err)
	}

	s3Client := s3.New(newSession)

	bucket := aws.String(c.Bucket)
	uploadFile := aws.String(to)

	cparams := &s3.CreateBucketInput{
		Bucket: bucket,
	}

	_, err = s3Client.CreateBucket(cparams)
	if err != nil {
		return err
	}

	f, err := os.ReadFile(filepath)
	if err != nil {
		return fmt.Errorf("failed to read file to upload: %s", err)
	}

	_, err = s3Client.PutObject(&s3.PutObjectInput{
		Body:   strings.NewReader(string(f)),
		Bucket: bucket,
		Key:    uploadFile,
	})
	if err != nil {
		return fmt.Errorf("failed to upload data to %s/%s, %s", c.Bucket, to, err)
	}

	fmt.Printf("Successfully uploaded data to: %s\n", c.Bucket)

	return nil
}

func (c S3Config) S3DownloadFile(endPath string) error {
	return c.S3DownloadFileFrom(c.File, endPath)
}

func (c S3Config) S3DownloadFileFrom(s3File, endPath string) error {
	bucket := aws.String(c.Bucket)

	s3Config := &aws.Config{
		Credentials:      credentials.NewStaticCredentials(c.AccessKey, c.SecretKey, ""),
		Endpoint:         aws.String(c.Endpoint),
		Region:           aws.String(c.Region),
		DisableSSL:       aws.Bool(true),
		S3ForcePathStyle: aws.Bool(true),
	}
	newSession, err := session.NewSession(s3Config)
	if err != nil {
		return fmt.Errorf("error creating session: %s", err)
	}

	// Create the base dir if it does not exist
	baseDir := filepath.Dir(endPath)
	if _, err := os.Stat(baseDir); os.IsNotExist(err) {
		err := os.MkdirAll(baseDir, 0766)
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
			Key:    aws.String(s3File),
		})
	if err != nil {
		return fmt.Errorf("failed to download file: %s", err)
	}

	log.Printf("Downloaded file: %s (%v bytes)\n", file.Name(), numBytes)

	return nil
}

func (c S3Config) Delete(s3File string) error {
	s3Config := &aws.Config{
		Credentials:      credentials.NewStaticCredentials(c.AccessKey, c.SecretKey, ""),
		Endpoint:         aws.String(c.Endpoint),
		Region:           aws.String(c.Region),
		DisableSSL:       aws.Bool(true),
		S3ForcePathStyle: aws.Bool(true),
	}

	newSession, err := session.NewSession(s3Config)
	if err != nil {
		return fmt.Errorf("error creating session: %s", err)
	}

	svc := s3.New(newSession)

	bucket := aws.String(c.Bucket)
	deleteFile := aws.String(s3File)

	_, err = svc.DeleteObject(&s3.DeleteObjectInput{
		Bucket: bucket,
		Key:    deleteFile,
	})
	if err != nil {
		return fmt.Errorf("failed to call delete object: %s", err)
	}

	err = svc.WaitUntilObjectNotExists(&s3.HeadObjectInput{
		Bucket: bucket,
		Key:    deleteFile,
	})
	if err != nil {
		return fmt.Errorf("failed to wait for delete")
	}

	fmt.Printf("Successfully deleted file: %s\n", s3File)

	return nil
}
