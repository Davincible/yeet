package cli

import (
	"fmt"
	"os"
	"reflect"
	"strings"

	cli "cli/providers/cli/parse"

	conf "github.com/BoRuDar/configuration/v4"
)

const (
	FlagProviderName = `FlagProvider`
	flagSeparator    = "|"
)

type FlagProviderOption func(*flagProvider)

// FlagSet is the part of flag.FlagSet that NewFlagProvider uses
type FlagSet interface {
	Parse([]string) error
	String(string, string, string) *string
}

type flagProvider struct {
	args  []string
	flags map[string]string
}

type flagData struct {
	key        string
	defaultVal string
	usage      string
}

// WithArgs specifies the custom args to parse instead of os.Args[1:].
// Mostly used for testing. There would be no need to use this yourself.
func WithArgs(args []string) FlagProviderOption {
	return func(c *flagProvider) {
		c.args = args
	}
}

// NewFlagProvider creates a new provider to fetch data from flags like: --flag_name some_value
func NewFlagProvider(opts ...FlagProviderOption) *flagProvider {
	fp := flagProvider{
		args:  os.Args[1:],
		flags: map[string]string{},
	}

	for _, f := range opts {
		f(&fp)
	}

	return &fp
}

func (*flagProvider) Name() string {
	return FlagProviderName
}

func (fp *flagProvider) Init(ptr any) (err error) {
	fp.flags = cli.Parse(fp.args)

	return nil
}

func (fp *flagProvider) Provide(field reflect.StructField, v reflect.Value) error {
	fd, err := fp.getFlagData(field)
	if err != nil {
		fmt.Println(err)
		return err
	}

	val, ok := fp.flags[fd.key]
	if !ok {
		fmt.Println("not foundddd", fp.flags)
		return conf.ErrEmptyValue
	}

	if err := SetField(field, v, val); err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

// getFlagData retrieves the flag tag.
func (fp *flagProvider) getFlagData(field reflect.StructField) (*flagData, error) {
	key := field.Tag.Get("flag")
	if len(key) == 0 {
		return nil, conf.ErrNoTag
	}

	flagInfo := strings.Split(key, flagSeparator)
	switch len(flagInfo) {
	case 3:
		return &flagData{
			key:        strings.TrimSpace(flagInfo[0]),
			defaultVal: strings.TrimSpace(flagInfo[1]),
			usage:      flagInfo[2],
		}, nil

	case 2:
		return &flagData{
			key:        strings.TrimSpace(flagInfo[0]),
			defaultVal: strings.TrimSpace(flagInfo[1]),
		}, nil

	case 1:
		return &flagData{
			key: strings.TrimSpace(flagInfo[0]),
		}, nil

	default:
		return nil, fmt.Errorf("wrong flag definition [%s]", key)
	}
}
