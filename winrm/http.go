package winrm

import (
	"fmt"
	"net/http"
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
