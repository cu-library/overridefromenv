// Copyright 2026 Carleton University Library
// All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package overridefromenv

import (
	"flag"
	"testing"
	"time"
)

func TestOverrideIgnoreSetFlags(t *testing.T) {

	prefix := "OVERRIDEFROMENVTEST_"
	target := prefix + "TEST"

	t.Setenv(target, "override")

	// Setup the test flag.
	fs := flag.NewFlagSet("test", flag.ExitOnError)
	s := fs.String("test", "default", "")
	fs.Set("test", "newvalue")

	Override(fs, prefix)

	if *s != "newvalue" {
		t.Error("An already set flag was overridden.")
	}
}

func TestOverrideError(t *testing.T) {

	prefix := "OVERRIDEFROMENVTEST_"
	target := prefix + "TEST"

	t.Setenv(target, "override")

	// Setup the test flag.
	fs := flag.NewFlagSet("test", flag.ExitOnError)
	fs.Float64("test", 0.1, "")

	err := Override(fs, prefix)

	if err == nil {
		t.Error("Overriding a float flag with a string didn't cause an error.")
	}
}

func TestOverrideUnsetFlags(t *testing.T) {

	prefix := "OVERRIDEFROMENVTEST_"

	fs := flag.NewFlagSet("test", flag.ExitOnError)

	b := fs.Bool("booltest", true, "")
	t.Setenv(prefix+"BOOLTEST", "false")

	defaultduration, _ := time.ParseDuration("1h")
	d := fs.Duration("durationtest", defaultduration, "")
	nd, _ := time.ParseDuration("2h")
	t.Setenv(prefix+"DURATIONTEST", "2h")

	fl := fs.Float64("floattest", 0.1, "")
	t.Setenv(prefix+"FLOATTEST", "0.2")

	i := fs.Int("inttest", 1, "")
	t.Setenv(prefix+"INTTEST", "2")

	i64 := fs.Int64("int64test", 1, "")
	t.Setenv(prefix+"INT64TEST", "2")

	s := fs.String("stringtest", "default", "")
	t.Setenv(prefix+"STRINGTEST", "newvalue")

	u := fs.Uint64("uinttest", 1, "")
	t.Setenv(prefix+"UINTTEST", "2")

	u64 := fs.Uint64("uint64test", 1, "")
	t.Setenv(prefix+"UINT64TEST", "2")

	Override(fs, prefix)

	if *b != false {
		t.Error("bool flag was not overwritten.")
	}
	if *d != nd {
		t.Error("duration flag was not overwritten.")
	}
	if *fl != 0.2 {
		t.Error("float flag was not overwritten.")
	}
	if *i != 2 {
		t.Error("int flag was not overwritten.")
	}
	if *i64 != 2 {
		t.Error("int64 flag was not overwritten.")
	}
	if *s != "newvalue" {
		t.Error("string flag was not overwritten.")
	}
	if *u != 2 {
		t.Error("uint flag was not overwritten.")
	}
	if *u64 != 2 {
		t.Error("uint64 flag was not overwritten.")
	}
}

func TestOverrideUnsetFlagsNormalizeKey(t *testing.T) {

	prefix_without_underscore := "OVERRIDEFROMENVTEST"
	prefix := prefix_without_underscore + "_"

	fs := flag.NewFlagSet("test", flag.ExitOnError)

	b := fs.Bool("bool-test", true, "")
	t.Setenv(prefix+"BOOL_TEST", "false")

	defaultduration, _ := time.ParseDuration("1h")
	d := fs.Duration("duration-test", defaultduration, "")
	nd, _ := time.ParseDuration("2h")
	t.Setenv(prefix+"DURATION_TEST", "2h")

	fl := fs.Float64("float-test", 0.1, "")
	t.Setenv(prefix+"FLOAT_TEST", "0.2")

	i := fs.Int("int-test", 1, "")
	t.Setenv(prefix+"INT_TEST", "2")

	i64 := fs.Int64("int-64-test", 1, "")
	t.Setenv(prefix+"INT_64_TEST", "2")

	s := fs.String("string-test", "default", "")
	t.Setenv(prefix+"STRING_TEST", "newvalue")

	u := fs.Uint64("uint-test", 1, "")
	t.Setenv(prefix+"UINT_TEST", "2")

	u64 := fs.Uint64("uint_64-test", 1, "")
	t.Setenv(prefix+"UINT_64_TEST", "2")

	Override(fs, prefix_without_underscore)

	if *b != false {
		t.Error("bool flag was not overwritten.")
	}
	if *d != nd {
		t.Error("duration flag was not overwritten.")
	}
	if *fl != 0.2 {
		t.Error("float flag was not overwritten.")
	}
	if *i != 2 {
		t.Error("int flag was not overwritten.")
	}
	if *i64 != 2 {
		t.Error("int64 flag was not overwritten.")
	}
	if *s != "newvalue" {
		t.Error("string flag was not overwritten.")
	}
	if *u != 2 {
		t.Error("uint flag was not overwritten.")
	}
	if *u64 != 2 {
		t.Error("uint64 flag was not overwritten.")
	}
}
