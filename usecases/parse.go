package usecases

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"strings"
)

type jsonStructure map[string]interface{}
type ParseResult struct {
	msjId   uint
	txtJson jsonStructure
}

func Parse() *[]*ParseResult {
	msjs := GetMeasurementSyncJob("progress")

	var parseResults []*ParseResult

	for i := 0; i < len(*msjs); i++ {
		txtFileName := strings.Replace((*msjs)[i].FileName, ".zip", ".txt", 1)

		txtFile, err := ioutil.ReadFile(filepath.Join(UnzipDestDir, txtFileName))
		if CheckError(err) {
			UpdateNoteAndStatus(ById((*msjs)[i].Id, err.Error(), "fail"))
			continue
		}

		syncEndPoint := ApiURL + "/api/v1/measurements/sync/" + fmt.Sprint((*msjs)[i].Id) + "/" + txtFileName
		syncBody := strings.NewReader(string(txtFile))
		syncHeader := make(http.Header)
		syncHeader.Add("Content-Type", "application/json")
		CallAPI(http.MethodPut, syncEndPoint, syncBody, &syncHeader)

		txtJson := new(jsonStructure)
		err = json.Unmarshal(txtFile, txtJson)
		if CheckError(err) {
			UpdateNoteAndStatus(ById((*msjs)[i].Id, err.Error(), "fail"))
			continue
		}

		parseResults = append(parseResults, &ParseResult{(*msjs)[i].Id, *txtJson})
	}

	return &parseResults
}
