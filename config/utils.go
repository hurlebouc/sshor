package config

import (
	"encoding/json"
	"os"
	"path/filepath"

	"cuelang.org/go/cue/cuecontext"
	"cuelang.org/go/cue/load"
)

func existCuePackage(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	if !info.IsDir() {
		return false
	}
	files, err := os.ReadDir(path)
	if err != nil {
		panic(err)
	}
	for _, file := range files {
		if filepath.Ext(file.Name()) == ".cue" {
			return true
		}
	}
	return false
}

func ReadConf() (*Config, error) {
	ctx := cuecontext.New()

	configdir, err := os.UserConfigDir()
	if err != nil {
		panic(err)
	}

	path := "."
	if !existCuePackage(path) {
		path, err = filepath.Abs(configdir + "/sshor")
		if err != nil {
			panic(err)
		}
		if !existCuePackage(path) {
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
