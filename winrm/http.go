package winrm

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
	"launchpad.net/xmlpath"
	"log"
	"net/http"
	"os"
	"strings"
)

type HttpError struct {
	StatusCode int
	Status     string
}

func (e *HttpError) Error() string {
	return fmt.Sprintf("[%d] %s", e.StatusCode, e.Status)
}

var ErrHttpAuthenticate = &HttpError{401, "Failed to authenticate"}

func fault(r *http.Response) error {
	if r.StatusCode == 401 {
		return ErrHttpAuthenticate
	}

	if h := r.Header.Get("Content-Type"); strings.HasPrefix(h, "application/soap+xml") {
		body, _ := ioutil.ReadAll(r.Body)
		buffer := bytes.NewBuffer(body)

		f := &HttpError{500, "Unparsable SOAP error"}
		root, err := xmlpath.Parse(buffer)

		if err != nil {
			return f
		}

		path := xmlpath.MustCompile("//Fault/Reason/Text")
		if reason, ok := path.String(root); ok {
			f.Status = "FAULT: " + reason
		}

		return f
	}

	return &HttpError{r.StatusCode, r.Status}
}

func deliver(user, pass, env string) (io.Reader, error) {
	if os.Getenv("WINRM_DEBUG") != "" {
		log.Println("delivering", string(env))
	}

	request, _ := http.NewRequest("POST",
		"http://localhost:5985/wsman", bytes.NewBufferString(env))
	request.Header.Add("Content-Type", "application/soap+xml;charset=UTF-8")
	request.Header.Add("Authorization", "Basic "+
		base64.StdEncoding.EncodeToString([]byte(user+":"+pass)))

	client := &http.Client{}
	response, err := client.Do(request)

	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	if response.StatusCode != 200 {
		return nil, fault(response)
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	if os.Getenv("WINRM_DEBUG") != "" {
		log.Println("receiving", string(body))
	}

	return bytes.NewReader(body), nil
}
