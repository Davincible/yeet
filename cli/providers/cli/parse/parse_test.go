package configuration

import (
	"strconv"
	"testing"
)

type Test struct {
	Args   []string
	Expect map[string]string
}

var (
	testData = []Test{
		{
			Args:   []string{"-name", "hi", "-two", "two", "--three", "three"},
			Expect: map[string]string{"name": "hi", "two": "two", "three": "three"},
		},
		{
			Args:   []string{"-this-is-a-key", "hi", "--this-key", "hi", "-two", "two", "--three", "three"},
			Expect: map[string]string{"this-is-a-key": "hi", "this-key": "hi", "two": "two", "three": "three"},
		},
		{
			Args:   []string{"--name-one", "Diana", "--name=Alex", "-some-bool", "--another-bool"},
			Expect: map[string]string{"name-one": "Diana", "name": "Alex", "some-bool": "true", "another-bool": "true"},
		},
		{
			Args:   []string{"--name-one", "\"Diana\"", "--name=\"Alex Hormozi\""},
			Expect: map[string]string{"name-one": "Diana", "name": "Alex Hormozi"},
		},
		{
			Args:   []string{"-config", "~/.config/sample.json"},
			Expect: map[string]string{"config": "~/.config/sample.json"},
		},
		{
			Args:   []string{"-bool", "--int", "5", "--test", "val=ue", "--another=\"hi=there\""},
			Expect: map[string]string{"bool": "true", "int": "5", "test": "val=ue", "another": "hi=there"},
		},
		{
			Args:   []string{"bool", "abc"},
			Expect: map[string]string{},
		},
	}
)

func TestParse(t *testing.T) {
	for i, data := range testData {
		t.Run("Test "+strconv.Itoa(i+1), func(t *testing.T) {
			runTest(t, i+1, data)
		})
	}
}

func runTest(t *testing.T, i int, data Test) {
	keys := Parse(data.Args)

	for key, value := range data.Expect {
		if keys[key] != value {
			t.Fatalf("[Test %d] key '%s' failed: expected '%s', but got '%s'", i, key, value, keys[key])
		}
	}

	if len(data.Expect) != len(keys) {
		t.Errorf("[Test %d] Expected: '%+v', got: '%+v'", i, data.Expect, keys)
		t.Fatalf("[Test %d] keys map has unexpected length. Expected %d, got %d", i, len(data.Args), len(keys))
	}
}
