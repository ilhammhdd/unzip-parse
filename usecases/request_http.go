package usecases

import (
	"io"
	"log"
	"net/http"
)

var (
	ApiURL    string
	ImagerURL string
)

func CallAPI(method, url string, body io.Reader, header *http.Header) {
	client := http.Client{}

	req, err := http.NewRequest(method, url, body)
	CheckError(err)

	req.Header = *header

	res, err := client.Do(req)
	CheckError(err)
	log.Println(url, "response : ", res)
}
