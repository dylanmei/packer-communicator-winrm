package winrm

import (
	"bytes"
	"encoding/base64"
	"encoding/xml"
	"fmt"
	"github.com/dylanmei/packer-communicator-winrm/winrm/envelope"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

type HttpError struct {
	StatusCode int
	message    string
}

func (he *HttpError) Error() string {
	return fmt.Sprintf("[%d] %s", he.StatusCode, he.message)
}

var ErrHttpAuthenticate = &HttpError{401, "Failed to authenticate"}

func NewHttpError(r *http.Response) *HttpError {
	if r.StatusCode == 401 {
		return ErrHttpAuthenticate
	}
	// fmt.Println("Unexpected HTTP Status", response.Status)
	// for key, value := range response.Header {
	// 	fmt.Println(" ", key, ":", value)
	// }
	return &HttpError{r.StatusCode, r.Status}
}

func sendEnvelope(user, pass string, env *envelope.Envelope) (io.Reader, error) {
	xmlEnvelope, err := xml.MarshalIndent(env, " ", "	")
	if err != nil {
		return nil, err
	}

	if os.Getenv("WINRM_DEBUG") != "" {
		log.Println("sending", string(xmlEnvelope))
	}

	request, _ := http.NewRequest("POST",
		"http://localhost:5985/wsman", bytes.NewReader(xmlEnvelope))
	request.Header.Add("Content-Type", "application/soap+xml;charset=UTF-8")
	request.Header.Add("Authorization", "Basic "+
		base64.StdEncoding.EncodeToString([]byte(user+":"+pass)))

	client := &http.Client{}
	response, err := client.Do(request)

	if err != nil {
		return nil, err
	}
	if response.StatusCode != 200 {
		return nil, NewHttpError(response)
	}

	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	if os.Getenv("WINRM_DEBUG") != "" {
		log.Println("receiving", string(body))
	}

	return bytes.NewReader(body), nil
}
