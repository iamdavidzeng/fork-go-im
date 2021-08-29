package utils

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

type Header struct {
	Authorizatoin string
	Token         string
}

func PostFile(filename string, tragetUrl string, headers *Header) (*http.Response, error) {
	bodyBuf := bytes.NewBufferString("")
	bodyWriter := multipart.NewWriter(bodyBuf)

	_, err := bodyWriter.CreateFormFile("smfile", filename)
	if err != nil {
		fmt.Println("error writing to buffer")
		return nil, err
	}

	fh, err := os.Open(filename)
	if err != nil {
		fmt.Println("error opening file")
		return nil, err
	}
	boundary := bodyWriter.Boundary()
	closeBuf := bytes.NewBufferString(fmt.Sprintf("\r\n--%s--\r\n", boundary))

	requestReader := io.MultiReader(bodyBuf, fh, closeBuf)
	fi, err := fh.Stat()
	if err != nil {
		fmt.Printf("Error Stating file: %s", filename)
		return nil, err
	}
	req, err := http.NewRequest("POST", tragetUrl, requestReader)
	if err != nil {
		return nil, err
	}

	req.Header.Add(headers.Authorizatoin, headers.Token)
	req.Header.Add("Content-Type", "multipart/form-data; boundary="+boundary)
	req.ContentLength = int64(closeBuf.Len()) + int64(bodyBuf.Len()) + fi.Size()
	return http.DefaultClient.Do(req)
}

func GetCurrentDirectory() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		fmt.Println(err)
	}
	return dir
}
