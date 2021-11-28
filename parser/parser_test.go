package parser

import "testing"

func TestMain(t *testing.T) {

	var firstCalled bool
	var secondCalled string
	var thirdCalled1 string
	var thirdCalled2 string

	p := New()
	p.Register("first-func", func() {
		firstCalled = true
	})
	p.Register("second-func", func(args string) {
		secondCalled = args
	})
	p.Register("third-func", func(cmd, args string) {
		thirdCalled1 = cmd
		thirdCalled2 = args
	})

	p.Call("first-func")

	if firstCalled == false {
		t.Errorf("Expceted fist function to be called")
	}

	p.Call("second-func with-param")

	if secondCalled != "with-param" {
		t.Errorf("Expceted second function to have param")
	}

	p.Call("third-func with-param and all other")

	if thirdCalled1 != "with-param" && thirdCalled2 != "and all other" {
		t.Error("Unexpected parameters")
	}
}
