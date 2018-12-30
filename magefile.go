// +build mage

package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/magefile/mage/mg" // mg contains helpful utility functions, like Deps
	"github.com/magefile/mage/sh"
)

func init() {
	// ProjectName
	os.Setenv("MATTEL_REPO", "go-chilopod")
	os.Setenv("TAG", "v0.1.2")
}

// Default target to run when none is specified
// If not set, running mage will list available targets
// var Default = Build

// A build step that requires additional params, or platform specific steps for example
func Build_v1() error {
	mg.Deps(InstallDeps)
	fmt.Println("Building...")
	cmd := exec.Command("go", "build", "-o", "MyApp", ".")
	return cmd.Run()
}

func Build() error {
	mg.Deps(Dep)
	return sh.RunV("go", "build", "-o", "$MATTEL_REPO", "-ldflags="+ldflags(), "github.com/riboseyim/$MATTEL_REPO")
}

func Dep() error {
	fmt.Println("Installing Deps...")
	cmd := exec.Command("dep", "ensure")
	return cmd.Run()
}

// A custom install step if you need your bin someplace other than go/bin
func Install() error {
	mg.Deps(Build)
	fmt.Println("Installing...")
	return os.Rename("./MyApp", "/usr/bin/MyApp")
}

// Manage your deps, or running package managers.
func InstallDeps() error {
	fmt.Println("Installing Deps...")
	cmd := exec.Command("go", "get", "github.com/stretchr/piglatin")
	return cmd.Run()
}

// Clean up after yourself
func Clean() {
	fmt.Println("Cleaning...")
	os.RemoveAll("MyApp")
}

func Tools() error {
	//mg.Deps(Protoc)

	update := envBool("UPDATE")

	retool := "github.com/twitchtv/retool"

	args := []string{"get", retool}
	if update {
		args = []string{"get", "-u", retool}
	}

	if err := sh.Run("go", args...); err != nil {
		return err
	}

	return sh.Run("retool", "sync")
}

// retool runs a command using a retool-cached binary.
func retool(cmd string, args ...string) error {
	return sh.Run("retool", append([]string{"do", cmd}, args...)...)
}

func Release() (err error) {
	if os.Getenv("TAG") == "" {
		return errors.New("TAG environment variable is required")
	}
	if err := sh.RunV("git", "tag", "-a", "$TAG"); err != nil {
		return err
	}
	if err := sh.RunV("git", "push", "origin", "$TAG"); err != nil {
		return err
	}

	defer func() {
		if err != nil {
			sh.RunV("git", "tag", "--delete", "$TAG")
			sh.RunV("git", "push", "--delete", "origin", "$TAG")
		}
	}()
	return sh.RunV("goreleaser", "release", "--skip-publish")
	//return retool("goreleaser")
}

func ldflags() string {
	timestamp := time.Now().Format(time.RFC3339)
	hash := hash()
	tag := tag()
	if tag == "" {
		tag = "dev"
	}

	return fmt.Sprintf(`-X "github.com/riboseyim/project/proj.timestamp=%s" `+
		`-X "github.com/riboseyim/project/proj.commitHash=%s" `+
		`-X "github.com/riboseyim/project/proj.gitTag=%s"`, timestamp, hash, tag)
}

// tag returns the git tag for the current branch or "" if none.
func tag() string {
	s, _ := sh.Output("git", "describe", "--tags")
	return s
}

// hash returns the git hash for the current repo or "" if none.
func hash() string {
	hash, _ := sh.Output("git", "rev-parse", "--short", "HEAD")
	return hash
}

func envBool(key string) bool {
	value := os.Getenv(key)
	if value != "" {
		return true
	} else {
		return false
	}
}
