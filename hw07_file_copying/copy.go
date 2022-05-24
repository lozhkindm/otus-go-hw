package main

import (
	"errors"
	"io"
	"os"
	"time"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
)

func Copy(fromPath, toPath string, offset, limit int64) (err error) {
	var (
		srcFileInfo      os.FileInfo
		srcFile, dstFile *os.File
		progressBar      *ProgressBar
		bytesToCopy      int64
		bytesCopied      int64
		writer           io.Writer
	)
	if srcFileInfo, err = os.Stat(fromPath); err != nil {
		return
	}
	if offset == 0 {
		offset = 0
	}
	if limit == 0 {
		limit = srcFileInfo.Size()
	}
	if srcFileInfo.Size() == 0 {
		return ErrUnsupportedFile
	}
	if offset > srcFileInfo.Size() {
		return ErrOffsetExceedsFileSize
	}

	if srcFile, err = os.Open(fromPath); err != nil {
		return
	}
	defer srcFile.Close()

	if dstFile, err = os.Create(toPath); err != nil {
		return
	}
	defer dstFile.Close()

	if _, err = srcFile.Seek(offset, io.SeekStart); err != nil && !errors.Is(err, io.EOF) {
		return
	}

	bytesToCopy = getBytesToCopy(srcFileInfo.Size(), offset, limit)
	progressBar = NewProgressBar(bytesToCopy, 50)
	writer = io.MultiWriter(dstFile, progressBar)

	for !errors.Is(err, io.EOF) {
		if _, err = io.CopyN(writer, srcFile, 1); err != nil && !errors.Is(err, io.EOF) {
			return
		}
		bytesCopied++
		if bytesCopied == limit {
			break
		}
		time.Sleep(100 * time.Microsecond) // delay to see how the progressbar works
	}
	if errors.Is(err, io.EOF) {
		err = nil
	}
	return
}

func getBytesToCopy(filesize, offset, limit int64) int64 {
	if filesize-offset > limit {
		return limit
	}
	return filesize - offset
}
