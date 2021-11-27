package builder

import "testing"

func TestMain(t *testing.T) {
	b := New()
	b.Add("ffmpeg -re -i      ")
	b.Add("http://11.1.1.1/url")

	if b.String() != "ffmpeg -re -i http://11.1.1.1/url" {
		t.Errorf("unexpected output: %v", b.String())
	}
}
