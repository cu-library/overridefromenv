// Copyright 2019 Carleton University Library
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
// It finds unset flags in fs, then sets those flags using the value of the
// environment variable with the key strings.ToUpper(prefix+flag.Name).
func Override(fs *flag.FlagSet, prefix string) error {

	// A map of pointers to unset flags.
	listOfUnsetFlags := make(map[*flag.Flag]bool)

	// Visit calls a function on "only those flags that have been set."
	// VisitAll calls a function on "all flags, even those not set."
	// No way to ask for "only unset flags". So, we add all, then
	// delete the set flags.

	// First, visit all the flags, and add them to our map.
	fs.VisitAll(func(f *flag.Flag) { listOfUnsetFlags[f] = true })

	// Then delete the set flags.
	fs.Visit(func(f *flag.Flag) { delete(listOfUnsetFlags, f) })

	// Loop through our list of unset flags.
	// We don't care about the values in our map, only the keys.
	for f := range listOfUnsetFlags {
		// Build the corresponding environment variable name for each flag.
		envVarName := fmt.Sprintf("%v%v", strings.ToUpper(prefix), strings.ToUpper(f.Name))

		// Look for the environment variable name.
		// If found, set the flag to that value.
		// If there's a problem setting the flag value,
		// there's a serious problem we can't recover from.
		envVarValue, found := os.LookupEnv(envVarName)
		if found {
			err := f.Value.Set(envVarValue)
			if err != nil {
				return fmt.Errorf("unable to set flag %v from environment variable %v, "+
					"which has a value of \"%v\": %w",
					f.Name, envVarName, envVarValue, err)
			}
		}
	}
	return nil
}
