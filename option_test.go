package messagequeue

import "testing"

func TestConfigMap(t *testing.T) {
	c := &ConfigMap{}
	err := c.Set("test", "testvalue")
	if err != nil {
		t.Fatal(err)
	}
	var result = ""
	err = c.Get("test", &result)
	if err != nil {
		t.Fatal(err)
	}
	if result != "testvalue" {
		t.Fatal(result)
	}
}
