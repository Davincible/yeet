package cli

import (
	"reflect"
	"testing"

	"github.com/BoRuDar/configuration/v4"
)

type One struct {
	Name string `default:"Alex" flag:"name" env:"NAME"`
}

type Test struct {
	Args     []string
	Expected any
}

type testOne struct {
	Name   string   `flag:"name"`
	TTL    int      `flag:"ttl"`
	Secure bool     `flag:"secure"`
	Files  []string `flag:"files"`
}

var (
	testData = []Test{
		{
			Args: []string{"-name", "Lia", "--ttl=5", "--secure", "--files=/home/config,/var/log,/home/lia"},
			Expected: testOne{
				Name:   "Lia",
				TTL:    5,
				Secure: true,
				Files:  []string{"/home/config", "/var/log", "/home/lia"},
			},
		},
	}
)

func TestCliProvider(t *testing.T) {
	data := testData[0]

	cfg := new(testOne)
	c := configuration.New(cfg,
		NewFlagProvider(WithArgs(data.Args)),
	)

	if err := c.InitValues(); err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(data.Expected.(testOne), *cfg) {
		t.Fatalf("Structs are not deeply equal: %+v != %+v", cfg, data.Expected)
	}
}
