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
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"golang.org/x/mod/semver"
)

const updatesURL = "https://go.dev/dl/?mode=json"

var (
	help      bool
	versionRe *regexp.Regexp

	logger = slog.New(slog.NewTextHandler(os.Stdout, nil))
)

func init() {
	flag.BoolVar(&help, "help", false, "Help message")
	versionRe = regexp.MustCompile(`^(?:1\.(\d+))(?:\.(\d+))?$`)
}

type goVersion struct {
	major int
	minor int
}

type VersionInfo struct {
	Version string `json:"version"`
}

func newGoVersion(v string) *goVersion {
	c := semver.Canonical("v" + v)
	if c == "" {
		logger.Error("couldn't parse semver", "version", v)
		os.Exit(1)
	}
	m := strings.Split(c, ".")
	major, err := strconv.Atoi(string(m[1]))
	if err != nil {
		logger.Error("error parsing major verison", "error", err)
		os.Exit(1)
	}
	minor, err := strconv.Atoi(string(m[2]))
	if err != nil {
		logger.Error("error parsing minor verison", "error", err)
		os.Exit(1)
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
	if g.major < 21 && g.minor == 0 {
		return g.Major()
	}
	return fmt.Sprintf("1.%d.%d", g.major, g.minor)
}

// String returns the full version string.
func (g *goVersion) String() string {
	return g.golangVersion()
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

func fetchJSON(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func getAvailableVersions() []goVersion {
	var availableVersions []goVersion

	jsonData, err := fetchJSON(updatesURL)
	if err != nil {
		logger.Error("Error fetching JSON", "error", err)
		return availableVersions
	}

	var availableVersionsJSON []VersionInfo
	if err := json.Unmarshal(jsonData, &availableVersionsJSON); err != nil {
		logger.Error("Error parsing JSON", "error", err)
		return availableVersions
	}

	for i := range availableVersionsJSON {
		// remove "go" from a string like "go1.22.0"
		newGoVersion := newGoVersion(strings.TrimLeft(availableVersionsJSON[i].Version, "go"))
		availableVersions = append(availableVersions, *newGoVersion)
	}
	logger.Info("found available versions", "num", len(availableVersions), "versions", availableVersions)
	return availableVersions
}

// getSHA256 returns the SHA256 of the Go archive.
func (g *goVersion) getSHA256() (string, error) {
	resp, err := http.Get(g.url() + ".sha256")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(b)), nil
}

// getLastMinorVersion returns the last minor version for a given Go version.
// if no new minor version available, it will return the given Go version back.
func (g *goVersion) getLastMinorVersion(availableVersions []goVersion) *goVersion {
	sort.Slice(availableVersions, func(i, j int) bool {
		if availableVersions[i].major == availableVersions[j].major {
			return availableVersions[i].minor < availableVersions[j].minor
		}
		return availableVersions[i].major < availableVersions[j].major
	})

	for _, availableVersion := range availableVersions {
		if availableVersion.major == g.major && availableVersion.minor > g.minor {
			return &availableVersion
		}
	}

	return g
}

// getNextMajor returns the next Go major version for a given Go version.
// It returns nil if the current version is already the latest.
func (g *goVersion) getNextMajor(availableVersions []goVersion) *goVersion {
	version := newGoVersion(g.Major() + ".0")
	version.major++

	for _, availableVersion := range availableVersions {
		if version.major == availableVersion.major {
			return version
		}
	}
	return nil
}

// getExactVersionFromDir reads the current Go version from a directory.
func getExactVersionFromDir(d string) (*goVersion, error) {
	re := regexp.MustCompile(fmt.Sprintf(`^\s*VERSION\s*:=\s*(%s(.\d+)?)`, d))
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
	return nil, fmt.Errorf("couldn't get exact version for %s", d)
}

func replace(filename string, replacers []func(string) (string, error)) error {
	b, err := os.ReadFile(filename)
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
	return os.WriteFile(filename, []byte(out), 0644)
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

func majorVersionReplacer(prefix string, old, new *goVersion) func(string) (string, error) {
	return func(out string) (string, error) {
		return strings.ReplaceAll(out, prefix+old.Major(), prefix+new.Major()), nil
	}
}

func golangVersionReplacer(prefix string, old, new *goVersion) func(string) (string, error) {
	return func(out string) (string, error) {
		return strings.ReplaceAll(out, prefix+old.golangVersion(), prefix+new.golangVersion()), nil
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
		if info.Name() == "Dockerfile" {
			return replace(path,
				[]func(string) (string, error){
					golangVersionReplacer("GOLANG_VERSION ", old, next),
					majorVersionReplacer("quay.io/prometheus/golang-builder:", old, next),
					shaReplacer(old, next),
				},
			)
		}
		return replace(path,
			[]func(string) (string, error){
				golangVersionReplacer("", old, next),
				majorVersionReplacer("", old, next),
			},
		)
	})
	if err != nil {
		return err
	}
	if err := os.Rename(old.Major(), next.Major()); err != nil {
		return fmt.Errorf("failed to create new version directory: %w", err)
	}

	// Update CircleCI.
	err = replace(".circleci/config.yml",
		[]func(string) (string, error){
			majorVersionReplacer("", current, next),
			majorVersionReplacer("", old, current),
		},
	)
	if err != nil {
		return err
	}

	// Update Makefile.
	err = replace("Makefile",
		[]func(string) (string, error){
			majorVersionReplacer("", current, next),
			majorVersionReplacer("", old, current),
		},
	)
	if err != nil {
		return err
	}

	// Update README.md.
	return replace("README.md",
		[]func(string) (string, error){
			fullVersionReplacer(current, next),
			majorVersionReplacer("", current, next),
			fullVersionReplacer(old, current),
			majorVersionReplacer("", old, current),
		},
	)
}

// updateNextMinor bumps the given directory to the next minor version.
// It returns nil if no new version exists.
func updateNextMinor(dir string, availableVersions []goVersion) (*goVersion, error) {
	current, err := getExactVersionFromDir(dir)
	if err != nil {
		return nil, fmt.Errorf("failed to detect current version of %s: %w", dir, err)
	}

	next := current.getLastMinorVersion(availableVersions)

	if next.equal(current) {
		logger.Info("no version change for Go", "version", next.golangVersion())
		return nil, nil
	}

	err = replace(filepath.Join(current.Major(), "base/Dockerfile"),
		[]func(string) (string, error){
			golangVersionReplacer("GOLANG_VERSION ", current, next),
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

	logger.Info("updated version", "current", current, "next", next)
	return next, nil
}

func main() {
	flag.Parse()
	if help {
		logger.Info("Bump Go versions in github.com/prometheus/golang-builder.")
		flag.PrintDefaults()
		os.Exit(0)
	}
	if err := run(); err != nil {
		logger.Error("update run failed", "error", err)
		os.Exit(1)
	}
}

func run() error {
	dirs := make([]string, 0)
	files, err := os.ReadDir(".")
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
		return fmt.Errorf("expected 2 versions of Go but got %d", len(dirs))
	}

	// Get list of available versions
	availableVersions := getAvailableVersions()
	if len(availableVersions) == 0 {
		logger.Error("failed to fetch avilable versions from update URL", "url", updatesURL)
		return errors.New("failed to fetch available versions")
	}

	// Check if a new major Go version exists.
	nexts := make([]*goVersion, 0)
	if next := newGoVersion(dirs[1] + ".0").getNextMajor(availableVersions); next != nil {
		logger.Info("found a new major version of Go", "version", next)
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
			logger.Info("processing version dir", "dir", d)
			next, err := updateNextMinor(d, availableVersions)
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
	logger.Info("Run the following command to commit the changes:")
	logger.Info(fmt.Sprintf("git checkout -b golang-%s", strings.Join(vs, "-")))
	logger.Info(fmt.Sprintf("git commit . --no-edit --message \"Bump to Go %s\"", strings.Join(vs, " and ")))

	return nil
}
