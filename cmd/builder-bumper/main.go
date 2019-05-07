// Copyright 2019 The Prometheus Authors
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Utility program to bump the Go versions in https://github.com/prometheus/golang-builder.
package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

var (
	help      bool
	versionRe *regexp.Regexp
)

func init() {
	flag.BoolVar(&help, "help", false, "Help message")
	versionRe = regexp.MustCompile(`^(?:1\.(\d+))(?:\.(\d+))?$`)
}

type goVersion struct {
	major int
	minor int
}

func newGoVersion(v string) *goVersion {
	m := versionRe.FindSubmatch([]byte(v))
	if len(m) != 3 {
		return nil
	}
	major, err := strconv.Atoi(string(m[1]))
	if err != nil {
		log.Fatal(err)
	}
	minor, err := strconv.Atoi(string(m[2]))
	if err != nil {
		log.Fatal(err)
	}
	return &goVersion{
		major: major,
		minor: minor,
	}
}

// major returns the version string without the minor version.
func (g *goVersion) Major() string {
	return fmt.Sprintf("1.%d", g.major)
}

// golangVersion returns the full version string but without the leading '.0'
// for the initial revision of a major release.
func (g *goVersion) golangVersion() string {
	if g.minor == 0 {
		return g.Major()
	}
	return g.String()
}

// String returns the full version string.
func (g *goVersion) String() string {
	return fmt.Sprintf("1.%d.%d", g.major, g.minor)
}

func (g *goVersion) less(o *goVersion) bool {
	if g.major == o.major {
		return g.minor < o.minor
	}
	return g.major < o.major
}

func (g *goVersion) equal(o *goVersion) bool {
	return g.major == o.major && g.minor == o.minor
}

// url returns the URL of the Go archive.
func (g *goVersion) url() string {
	return fmt.Sprintf("https://dl.google.com/go/go%s.linux-amd64.tar.gz", g.golangVersion())
}

// getSHA256 returns the SHA256 of the Go archive.
func (g *goVersion) getSHA256() (string, error) {
	resp, err := http.Get(g.url() + ".sha256")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(b)), nil
}

// getLastMinorVersion returns the last minor version for a given Go version.
func (g *goVersion) getLastMinorVersion() (*goVersion, error) {
	last := *g
	for {
		next := last
		next.minor++
		resp, err := http.Head(next.url())
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		if resp.StatusCode/100 != 2 {
			return &last, nil
		}
		last = next
	}
}

// getNextMajor returns the next Go major version for a given Go version.
// It returns nil if the current version is already the latest.
func (g *goVersion) getNextMajor() *goVersion {
	version := newGoVersion(g.Major() + ".0")
	version.major++

	resp, err := http.Head(version.url())
	if err != nil {
		return nil
	}
	defer resp.Body.Close()
	if resp.StatusCode/100 != 2 {
		return nil
	}

	return version
}

// getExactVersionFromDir reads the current Go version from a directory.
func getExactVersionFromDir(d string) (*goVersion, error) {
	re := regexp.MustCompile(fmt.Sprintf(`^\s*VERSION\s*:=\s*(%s.\d+)`, d))
	f, err := os.Open(filepath.Join(d, "Makefile.COMMON"))
	if err != nil {
		return nil, err
	}
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		m := re.FindSubmatch(scanner.Bytes())
		if m != nil {
			return newGoVersion(string(m[1])), nil
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return nil, errors.Errorf("couldn't get exact version for %s", d)
}

func replace(filename string, replacers []func(string) (string, error)) error {
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	out := string(b)
	for _, fn := range replacers {
		out, err = fn(out)
		if err != nil {
			return err
		}
	}
	return ioutil.WriteFile(filename, []byte(out), 0644)
}

func shaReplacer(old, new *goVersion) func(string) (string, error) {
	oldSHA, err := old.getSHA256()
	if err != nil {
		return func(string) (string, error) { return "", err }
	}
	nextSHA, err := new.getSHA256()
	if err != nil {
		return func(string) (string, error) { return "", err }
	}

	return func(out string) (string, error) {
		return strings.ReplaceAll(out, oldSHA, nextSHA), nil
	}
}

func majorVersionReplacer(old, new *goVersion) func(string) (string, error) {
	return func(out string) (string, error) {
		return strings.ReplaceAll(out, old.Major(), new.Major()), nil
	}
}

func golangVersionReplacer(old, new *goVersion) func(string) (string, error) {
	return func(out string) (string, error) {
		return strings.ReplaceAll(out, old.golangVersion(), new.golangVersion()), nil
	}
}

func fullVersionReplacer(old, new *goVersion) func(string) (string, error) {
	return func(out string) (string, error) {
		return strings.ReplaceAll(out, old.String(), new.String()), nil
	}
}

// replaceMajor switches the versions from [1.(N-1), 1.N] to [1.N, 1.(N+1)].
func replaceMajor(old, current, next *goVersion) error {
	// Replace the old version by the next one.
	err := filepath.Walk(old.Major(), func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if info.Name() == "Makefile.COMMON" {
			return replace(path,
				[]func(string) (string, error){
					fullVersionReplacer(old, next),
				},
			)
		}
		return replace(path,
			[]func(string) (string, error){
				golangVersionReplacer(old, next),
				majorVersionReplacer(old, next),
				shaReplacer(old, next),
			},
		)
	})
	if err != nil {
		return err
	}
	if err := os.Rename(old.Major(), next.Major()); err != nil {
		return errors.Wrap(err, "failed to create new version directory")
	}

	// Update Makefile.
	err = replace("Makefile",
		[]func(string) (string, error){
			majorVersionReplacer(current, next),
			majorVersionReplacer(old, current),
		},
	)
	if err != nil {
		return err
	}

	// Update README.md.
	return replace("README.md",
		[]func(string) (string, error){
			fullVersionReplacer(current, next),
			majorVersionReplacer(current, next),
			majorVersionReplacer(old, current),
			fullVersionReplacer(old, current),
		},
	)
}

// updateNextMinor bumps the given directory to the next minor version.
// It returns nil if no new version exists.
func updateNextMinor(dir string) (*goVersion, error) {
	current, err := getExactVersionFromDir(dir)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to detect current version of %s", dir)
	}

	next, err := current.getLastMinorVersion()
	if err != nil {
		return nil, err
	}
	if next.equal(current) {
		log.Printf("no version change for Go %s", next.golangVersion())
		return nil, nil
	}

	err = replace(filepath.Join(current.Major(), "base/Dockerfile"),
		[]func(string) (string, error){
			golangVersionReplacer(current, next),
			shaReplacer(current, next),
		},
	)
	if err != nil {
		return nil, err
	}

	err = replace(filepath.Join(current.Major(), "Makefile.COMMON"),
		[]func(string) (string, error){
			fullVersionReplacer(current, next),
		},
	)
	if err != nil {
		return nil, err
	}

	err = replace(filepath.Join("README.md"),
		[]func(string) (string, error){
			fullVersionReplacer(current, next),
		},
	)
	if err != nil {
		return nil, err
	}

	log.Printf("updated from %s to %s", current, next)
	return next, nil
}

func main() {
	flag.Parse()
	if help {
		log.Print("Bump Go versions in github.com/prometheus/golang-builder.")
		flag.PrintDefaults()
		os.Exit(0)
	}
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	dirs := make([]string, 0)
	files, err := ioutil.ReadDir(".")
	if err != nil {
		return err
	}
	for _, f := range files {
		if !f.IsDir() {
			continue
		}
		if !versionRe.Match([]byte(f.Name())) {
			continue
		}
		dirs = append(dirs, f.Name())
	}

	if len(dirs) != 2 {
		return errors.Errorf("Expected 2 versions of Go but got %d\n", len(dirs))
	}

	// Check if a new major Go version exists.
	nexts := make([]*goVersion, 0)
	if next := newGoVersion(dirs[1] + ".0").getNextMajor(); next != nil {
		log.Printf("found a new major version of Go: %s", next)
		old, err := getExactVersionFromDir(dirs[0])
		if err != nil {
			return err
		}
		current, err := getExactVersionFromDir(dirs[1])
		if err != nil {
			return err
		}
		if err = replaceMajor(old, current, next); err != nil {
			return err
		}
		nexts = append(nexts, next)
	} else {
		// Otherwise check for new minor versions.
		for _, d := range dirs {
			log.Printf("processing %s", d)
			next, err := updateNextMinor(d)
			if err != nil {
				return err
			}
			if next != nil {
				nexts = append(nexts, next)
			}
		}
	}

	if len(nexts) == 0 {
		return nil
	}

	sort.SliceStable(nexts, func(i, j int) bool {
		return nexts[i].less(nexts[j])
	})
	vs := make([]string, 0)
	for _, v := range nexts {
		vs = append(vs, v.String())
	}
	log.Print("Run the following command to commit the changes:")
	log.Printf("git checkout -b golang-%s", strings.Join(vs, "-"))
	log.Printf("git commit . --no-edit --message \"Bump to Go %s\"", strings.Join(vs, " and "))

	return nil
}
