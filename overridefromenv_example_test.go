// Copyright 2026 Carleton University Library
// All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package overridefromenv_test

import (
	"flag"
	"fmt"
	"os"
	"testing"

	"github.com/cu-library/overridefromenv/v2"
)

func TestMain(m *testing.M) {

	prefix := "APP_"
	envs := []struct {
		name  string
		value string
	}{
		{"PORT", "9090"},
		{"CONFIG_FILE", "my-config.toml"},
	}

	// Borrowed from the stdlib's cleanup after using testing's t.Setenv().
	for _, env := range envs {
		key := prefix + env.name
		prevValue, ok := os.LookupEnv(key)
		err := os.Setenv(key, env.value)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Could not set environment variable:", err)
			os.Exit(1)
		}
		if ok {
			defer os.Setenv(key, prevValue)
		} else {
			defer os.Unsetenv(key)
		}
	}

	m.Run()
}

func ExampleOverride() {
	// We assume the environment variables APP_PORT
	// and APP_CONFIG_FILE have been set.
	// Normally, we don't need to create a new FlagSet.
	// Instead, we use the flag package's CommandLine, which is
	// the default set of command-line flags, parsed from os.Args.
	fs := flag.NewFlagSet("demo", flag.ContinueOnError)

	prefix := "APP"

	host := fs.String("host", "localhost", "server host")
	port := fs.Int("port", 8080, "server port")
	config := fs.String("config-file", "config.toml", "config file")

	// Simulate user explicitly setting one flag (port) via parsing args.
	_ = fs.Parse([]string{"-port=7777"})

	// Override the defaults of unset flags from the environment.
	// Normally, you would use flag.CommandLine as the first
	// argument, but in this example we use the explicit FlagSet fs.
	err := overridefromenv.Override(fs, prefix)
	if err != nil {
		fmt.Printf("error: %v", err)
		return
	}

	// host should be the default, as it wasn't explicitly set and
	// APP_HOST wasn't an environment variable.
	fmt.Println(*host)
	// port should be the value which was set explicitly.
	fmt.Println(*port)
	// config should be the value of the environment variable APP_CONFIG_FILE.
	fmt.Println(*config)

	// Output:
	// localhost
	// 7777
	// my-config.toml
}
