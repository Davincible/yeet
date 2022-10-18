package configuration

import (
	"strings"
)

// Parse takes the list of arguments os.Args[1:] and parses them into a map.
func Parse(args []string) map[string]string {
	keys := make(map[string]string)

	for i := 0; i < len(args); i++ {
		key := args[i]
		if !strings.HasPrefix(key, "-") {
			continue
		}

		// Match both -abc and --abc. Trim twice seperately as might only contain one
		key = strings.TrimPrefix(key, "-")
		key = strings.TrimPrefix(key, "-")

		var val string

		// Check if value is contained within this arg, with =, or on the next one
		if strings.Contains(key, "=") {
			// Check if current key is set with =
			split := strings.SplitN(key, "=", 2)
			key = split[0]
			val = split[1]
		} else if i+1 == len(args) || strings.HasPrefix(args[i+1], "-") {
			// If next arg starts with -, or last arg, set current key as bool true
			val = "true"
		} else {
			// Set value equal to next arg
			i++
			val = args[i]
		}

		keys[key] = val
	}

	return keys
}
