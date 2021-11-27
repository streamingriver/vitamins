package registry

import (
	"io"
	"net/http"
)

type Fetcher interface {
	Fetch(string) error
	SetAuth(string, string)
}

type Fetch struct {
	auth *auth
}

func (f *Fetch) SetAuth(username, password string) {
	f.auth = &auth{username: username, password: password}
}

func (f *Fetch) Fetch(url string) error {
	HTTPClient := &http.Client{}

	request, _ := http.NewRequest("GET", url, nil)
	if f.auth != nil {
		request.SetBasicAuth(f.auth.username, f.auth.password)
	}

	response, err := HTTPClient.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	io.Copy(io.Discard, response.Body)
	if err != nil {
		return err
	}
	return nil

}
