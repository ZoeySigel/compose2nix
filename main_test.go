package main

import (
	"reflect"
	"strings"
	"testing"
)

func TestParseRuntimeEnvFiles(t *testing.T) {
	for _, tt := range []struct {
		name          string
		value         string
		only          bool
		ignoreMissing bool
		want          []string
		wantError     string
	}{
		{name: "omitted"},
		{name: "paths", value: "/run/secrets/app.env,relative.env", want: []string{"/run/secrets/app.env", "relative.env"}},
		{name: "boolean filename", value: "./true", want: []string{"./true"}},
		{name: "spaces in filename", value: "my file.env", want: []string{"my file.env"}},
		{name: "legacy true", value: "true", wantError: "not a boolean"},
		{name: "legacy false", value: "false", wantError: "not a boolean"},
		{name: "legacy uppercase", value: "TRUE", wantError: "not a boolean"},
		{name: "legacy numeric", value: "1", wantError: "not a boolean"},
		{name: "legacy env files only", only: true, wantError: "deprecated"},
		{name: "legacy ignore missing", ignoreMissing: true, wantError: "deprecated"},
		{name: "empty entry", value: "a.env,,b.env", wantError: "empty file path"},
		{name: "trailing comma", value: "a.env,", wantError: "empty file path"},
		{name: "blank", value: " ", wantError: "empty file path"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseRuntimeEnvFiles(tt.value, tt.only, tt.ignoreMissing)
			if tt.wantError != "" {
				if err == nil || !strings.Contains(err.Error(), tt.wantError) {
					t.Fatalf("got error %v, want %q", err, tt.wantError)
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("got %v, want %v", got, tt.want)
			}
		})
	}
}
