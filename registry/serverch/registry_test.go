package registry

import (
	"testing"
	"time"
)

func TestRegistry(t *testing.T) {
	djb33Seed = 1
	r := New("http")
	r.SetTTL(1)

	r.Ping("channel", "host", "80")
	r.Ping("channel", "host1", "80")

	// file, err := r.GetURL("uniq1", "channel", "path/to/file")
	file, err := r.GetURL("2", "channel", "path/to/file")

	if file != "http://host1:80/channel/path/to/file" {
		t.Errorf("unexpected output: %v", file)
	}

	file, err = r.GetURL("7oqwie0", "channel", "path/to/file")

	// r.Debug()
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

	file, err = r.GetURL("uniq", "channel_notexist", "/path/to/file")

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

	file, err = r.GetURL("uniq", "channel", "/path/to/file")

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
