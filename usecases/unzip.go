package usecases

import (
	"archive/zip"
	"io"
	"os"
	"path/filepath"
	"time"
)

var (
	UnzipSrcDir  string
	UnzipDestDir string
)

func Extractor() {
	var extractOk bool

	msjs := GetMeasurementSyncJob("pending")

	for i := 0; i < len(*msjs); i++ {
		filePath := filepath.Join(UnzipSrcDir, (*msjs)[i].FileName)
		_, err := os.Stat(filePath)
		checkCounter := 1
		for os.IsNotExist(err) {
			if checkCounter == 10 {
				break
			}
			time.Sleep(5 * time.Second)
			_, err = os.Stat(filePath)
			checkCounter++

		}
		UpdateNoteAndStatus(ById((*msjs)[i].Id, (*msjs)[i].FileName+" extracting", "progress"))
		ziprc, err := zip.OpenReader(filePath)
		if CheckError(err) {
			UpdateNoteAndStatus(ById((*msjs)[i].Id, err.Error(), "fail"))
			continue
		}
		defer func() {
			if ziprc != nil {
				CheckError(ziprc.Close())
			}
		}()

		for _, f := range ziprc.File {
			writeFile(f, &extractOk)
		}

		if extractOk {
			UpdateNoteAndStatus(ById((*msjs)[i].Id, (*msjs)[i].FileName+" extracted", "progress"))
		}
	}
}

func writeFile(zipFile *zip.File, extractOk *bool) {
	iorc, err := zipFile.Open()
	if CheckError(err) {
		UpdateNoteAndStatus(ByName(zipFile.Name, err.Error(), "fail"))
		return
	}
	defer func() {
		if iorc != nil {
			CheckError(iorc.Close())
		}
	}()

	path := filepath.Join(UnzipDestDir, zipFile.Name)

	if zipFile.FileInfo().IsDir() {
		os.MkdirAll(path, zipFile.Mode())
	} else {
		os.MkdirAll(filepath.Dir(path), zipFile.Mode())

		file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, zipFile.Mode())
		if CheckError(err) {
			UpdateNoteAndStatus(ByName(zipFile.Name, err.Error(), "fail"))
			return
		}

		_, err = io.Copy(file, iorc)
		if CheckError(err) {
			UpdateNoteAndStatus(ByName(zipFile.Name, err.Error(), "fail"))
			return
		}

		CheckError(file.Close())
	}

	*extractOk = true
}
