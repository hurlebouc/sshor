package config

import (
	"encoding/json"
	"os"
	"path/filepath"

	"cuelang.org/go/cue/cuecontext"
	"cuelang.org/go/cue/load"
)

func existFile(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return !info.IsDir()
}

func existDir(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.IsDir()
}

func ReadConf() (*Config, error) {
	ctx := cuecontext.New()

	configdir, err := os.UserConfigDir()
	if err != nil {
		panic(err)
	}

	path := "."
	if !existDir(path) {
		path, err = filepath.Abs(configdir + "/sshor")
		if err != nil {
			panic(err)
		}
		if !existDir(path) {
			return nil, nil
		}
	}

	// Load the package "example" from the current directory.
	// We don't need to specify a config in this example.
	insts := load.Instances([]string{"."}, &load.Config{
		Dir: path,
	})

	// The current directory just has one file without any build tags,
	// and that file belongs to the example package.
	// So we get a single instance as a result.
	value := ctx.BuildInstance(insts[0])

	if value.Err() != nil {
		return nil, value.Err()
	}

	jsonBytes, err := value.MarshalJSON()
	if err != nil {
		return nil, err
	}

	config := Config{}
	err = json.Unmarshal(jsonBytes, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
