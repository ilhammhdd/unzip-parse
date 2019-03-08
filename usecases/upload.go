package usecases

import (
	"bytes"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

type formDataValues map[string]io.Reader

var PhotoCategories []string

func UploadPhotos(parseResults *[]*ParseResult) {
	for j := 0; j < len(*parseResults); j++ {
		var uploadOk bool = true
	IteratePhotos:
		for i := 0; i < len(PhotoCategories); i++ {
			photoObj, ok := ((*(*parseResults)[j]).txtJson)[(PhotoCategories)[i]]
			if !ok {
				uploadOk = false
				UpdateNoteAndStatus(ById((*(*parseResults)[j]).msjId, (PhotoCategories)[i]+" object is not found", "fail"))
				break
			}
			photoObjSlc, ok := photoObj.([]interface{})
			if !ok {
				uploadOk = false
				UpdateNoteAndStatus(ById((*(*parseResults)[j]).msjId, (PhotoCategories)[i]+" object is not an array", "fail"))
				break
			}
			for i := 0; i < len(photoObjSlc); i++ {
				fileName, originalName, errMessage, ok := getNames(i, &photoObjSlc)
				if !ok {
					uploadOk = false
					UpdateNoteAndStatus(ById((*(*parseResults)[j]).msjId, errMessage, "fail"))
					break IteratePhotos
				}

				contentType, body, err := fillMultiPartForm(i, fileName, originalName)
				if CheckError(err) {
					uploadOk = false
					UpdateNoteAndStatus(ById((*(*parseResults)[j]).msjId, err.Error()+" in : "+(PhotoCategories)[i], "fail"))
					break IteratePhotos
				}

				uploadEndpoint := ImagerURL + "/upload"
				uploadHeader := make(http.Header)
				uploadHeader.Add("Content-Type", contentType)
				CallAPI(http.MethodPost, uploadEndpoint, body, &uploadHeader)
			}
		}

		if uploadOk {
			UpdateNoteAndStatus(ById((*(*parseResults)[j]).msjId, "successfully extracted, parsed, and uploaded", "done"))
		}
	}
}

func getNames(iteration int, photoObjSlc *[]interface{}) (fnm, onm, erm string, ok bool) {
	photoObjEleMap, ok := (*photoObjSlc)[iteration].(map[string]interface{})
	if !ok {
		return "", "", (PhotoCategories)[iteration] + " slice is not a map from json", false
	}

	fileName, ok := photoObjEleMap["file_name"]
	if !ok {
		return "", "", (PhotoCategories)[iteration] + " file_name key/value not exists", false
	}

	originalName, ok := photoObjEleMap["original_name"]
	if !ok {
		return "", "", (PhotoCategories)[iteration] + " original_name key/value not exists", false
	}
	return fileName.(string), originalName.(string), "", true
}

func fillMultiPartForm(iteration int, fileName, originalName string) (string, *bytes.Buffer, error) {
	var body bytes.Buffer
	mpWriter := multipart.NewWriter(&body)

	photoFilePath := filepath.Join(UnzipDestDir, originalName)

	w, err := mpWriter.CreateFormFile("imgFile", photoFilePath)
	if CheckError(err) {
		return "", nil, err
	}

	photoFile, err := os.Open(photoFilePath)
	if CheckError(err) {
		return "", nil, err
	}

	_, err = io.Copy(w, photoFile)
	if CheckError(err) {
		return "", nil, err
	}

	photoFile.Close()

	mpWriter.WriteField("code", fileName)

	err = mpWriter.Close()
	if CheckError(err) {
		return "", nil, err
	}

	return mpWriter.FormDataContentType(), &body, nil
}
