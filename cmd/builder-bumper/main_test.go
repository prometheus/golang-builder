package main

import (
	"reflect"
	"strings"
	"testing"
)

func TestNewGoVersion(t *testing.T) {
	tests := []struct {
		input string
		want  *goVersion
	}{
		{"1.22.1", &goVersion{22, 1}},
		{"1.20.5", &goVersion{20, 5}},
		{"1.19.0", &goVersion{19, 0}},
		{"1.25.3", &goVersion{25, 3}},
		{"1.18.10", &goVersion{18, 10}},
	}

	for _, tt := range tests {
		t.Run("Parsing "+tt.input, func(t *testing.T) {
			got := newGoVersion(tt.input)

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newGoVersion(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestAvailableVersions(t *testing.T) {
	t.Run("Test that available versions scan correctly removes 'go' prefix", func(t *testing.T) {
		got := getAvailableVersions()

		if len(got) == 0 {
			t.Fatalf("Expected some versions, got none")
		}

		for _, v := range got {
			if strings.Contains(v.String(), "go") {
				t.Errorf("Version still contains 'go' prefix: %v", v)
			}
		}
	})
}
