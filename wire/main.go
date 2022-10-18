package main

import "fmt"

type Options struct {
	mapy map[string]string
}

type OptionsTwo struct {
	mapy map[string]string
}

type Option func(o *Options)
type OptionTwo func(o *OptionsTwo)

func WithString(k, v string) Option {
	return func(o *Options) {
		fmt.Println("Setting options:", k, v)
		o.mapy[k] = v
	}
}

func WithStringTwo(k, v string) OptionTwo {
	return func(o *OptionsTwo) {
		fmt.Println("Setting options:", k, v)
		o.mapy[k] = v
	}
}

func ProvideOptions() []Option {
	return []Option{
		WithString("one", "hi"),
		WithString("two", "hi"),
		WithString("three", "hi"),
	}
}

func ProvideOptionsTwo() []OptionTwo {
	return []OptionTwo{
		WithStringTwo("one", "hi"),
		WithStringTwo("two", "hi"),
		WithStringTwo("three", "hi"),
	}
}

func ProvideOptionTwo(o ...OptionTwo) *OptionsTwo {
	opts := OptionsTwo{}

	for _, opt := range o {
		opt(&opts)
	}

	return &opts
}

func ProvideOption(o ...Option) *Options {
	opts := Options{}

	for _, opt := range o {
		opt(&opts)
	}

	return &opts
}

type ReturnType struct{}

func ProvideCollector(_ *Options, _ *OptionsTwo) (*ReturnType, error) {
	return new(ReturnType), nil
}
