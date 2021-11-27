package registry

import (
	"testing"
	"time"
)

func TestRegistry(t *testing.T) {
	r := New("http")
	r.SetTTL(1)
	r.Ping("channel", "host", "80")

	file, err := r.GetURL("channel", "/path/to/file")

	if file != "http://host:80/channel/path/to/file" {
		t.Errorf("unexpected output: %v", file)
	}

	if err != nil {
		if err.(*RegistryError).Expired() == true {
			t.Errorf("Expected expired to be false, got true")
		}
		if err.(*RegistryError).NotFound() == true {
			t.Errorf("Why channel found?")
		}
	}

	file, err = r.GetURL("channel_notexist", "/path/to/file")

	if file != "" {
		t.Errorf("unexpected output: %v", file)
	}

	if err != nil {
		if err.(*RegistryError).Expired() != false {
			t.Errorf("Expected expired to be false, got true")
		}
		if err.(*RegistryError).NotFound() != true {
			t.Errorf("not existent channel found?")
		}
	}
	time.Sleep(time.Millisecond * 1200)

	file, err = r.GetURL("channel", "/path/to/file")

	if file != "" {
		t.Errorf("unexpected output: %v", file)
	}

	if err != nil {
		if err.(*RegistryError).Expired() != true {
			t.Errorf("Expected expired to be true, got false")
		}
		if err.(*RegistryError).NotFound() != false {
			t.Errorf("not existent channel found?")
		}
	}
}
