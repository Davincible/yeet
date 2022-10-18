package main

import (
	"cli/providers/cli"
	"os"

	"fmt"

	conf "github.com/BoRuDar/configuration/v4"
)

type Config struct {
	Name string `default:"Alex" flag:"name" env:"NAME"`
}

func main() {
	cfg := struct {
		Name   string `flag:"name"`
		TTL    int    `flag:"ttl"`
		Secure bool   `flag:"secure"`
	}{}

	for _, val := range os.Args {
		fmt.Printf("'%s'\n", val)
	}

	configurator := conf.New(
		&cfg,
		// order of execution will be preserved:
		cli.NewFlagProvider(),
		// conf.NewFlagProvider(), // 1st
		// conf.NewEnvProvider(), // 2nd
		// conf.NewJSONFileProvider(fileName), // 3rd
		conf.NewDefaultProvider(), // 4th
	)

	if err := configurator.InitValues(); err != nil {
		panic(err)
	}

	fmt.Printf("Values: %+v\n", cfg)

	cfgg := struct {
		Dicky    string `flag:"dick"`
		LastName string `default:"defaultLastName"`
	}{}

	configurator = conf.New(
		&cfgg,
		// order of execution will be preserved:
		cli.NewFlagProvider(),
		// conf.NewFlagProvider(), // 1st
		conf.NewEnvProvider(), // 2nd
		// conf.NewJSONFileProvider(fileName), // 3rd
		conf.NewDefaultProvider(), // 4th
	)

	if err := configurator.InitValues(); err != nil {
		panic(err)
	}
	fmt.Println(cfgg)
}
