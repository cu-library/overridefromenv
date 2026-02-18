// Copyright 2026 Carleton University Library
// All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

// Package overridefromenv is a library which sets unset flags from environment variables.
package overridefromenv

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

// Override sets unset flags using environment variables.
// It finds unset flags in fs, then sets those flags using the value of an
// environment variable.
// To help group environment variables used by an application, an optional
// prefix can be provided. Prefixes that do not end in a '_' character have
// a '_' character appended.
// The name of the environment variable to lookup is the name of the unset
// flag, with the optional prefix prepended, converted to uppercase, and
// with any '-' characters replaced with '_' characters.
func Override(fs *flag.FlagSet, prefix string) error {
	// Add a trailing '_' character to the prefix if it is needed.
	if prefix != "" && !strings.HasSuffix(prefix, "_") {
		prefix += "_"
	}

	// The set of unset flag names.
	unset := make(map[string]struct{})

	// Visit calls a function on "only those flags that have been set."
	// VisitAll calls a function on "all flags, even those not set."
	// No way to ask for "only unset flags". So, we add all, then
	// delete the set flags.

	// First, visit all the flags, and add them to the set.
	fs.VisitAll(func(f *flag.Flag) { unset[f.Name] = struct{}{} })

	// Then delete the set flags.
	fs.Visit(func(f *flag.Flag) { delete(unset, f.Name) })

	for name := range unset {
		// Build the corresponding environment variable name for each flag.
		key := prefix + name
		key = strings.ReplaceAll(key, "-", "_")
		key = strings.ToUpper(key)

		// Look for the environment variable.
		// If found, set the flag to that variable's value.
		// If there's a problem setting the value, return an error.
		value, found := os.LookupEnv(key)
		if found {
			err := fs.Set(name, value)
			if err != nil {
				return fmt.Errorf("unable to set flag %q from environment variable %q: %w",
					name, key, err)
			}
		}
	}
	return nil

}

// Override the flag package's default set of command-line flags.
// flag.Parse() should be called first.
func OverrideCommandLine(prefix string) error {
	return Override(flag.CommandLine, prefix)
}
