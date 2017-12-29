package warpcache

import "testing"

func TestParseInputFormat(t *testing.T) {
	message := "1514474334340612// asset{name=l,.app=sw} 42\n"

	var classname string
	var labels map[string]string
	var value float64
	var err error
	classname, labels, value, err = parseInputFormat(message)
	if err != nil {
		t.Error("error")
	}
	if classname != "asset" {
		t.Error("bad classname")
	}
	if labels["name"] != "l" {
		t.Error("bad label")
	}
	if value != 42.0 {
		t.Error("bad value")
	}
}
