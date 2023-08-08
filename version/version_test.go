package version

import (
	"testing"

	"github.com/go-vela/types/version"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	t.Run("Valid Version Output", func(t *testing.T) {
		// to avoid flaky tests, set the values
		Go = "go1.20.7"
		OS = "darwin"
		Compiler = "gc"
		Tag = "v1.1.1"
		Commit = "000"
		Date = "111"
		Arch = "arm64"

		expected := version.Version{
			Canonical: "v1.1.1",
			Major:     1,
			Minor:     1,
			Patch:     1,
			Metadata: version.Metadata{
				Architecture:    "arm64",
				BuildDate:       "111",
				Compiler:        "gc",
				GitCommit:       "000",
				GoVersion:       "go1.20.7",
				OperatingSystem: "darwin",
			},
		}
		assert.Equal(t, expected, *New())
	})
	t.Run("no tag", func(t *testing.T) {
		// to avoid flaky tests, set the values
		Go = "go1.20.7"
		OS = "darwin"
		Compiler = "gc"
		Tag = ""
		Commit = "000"
		Date = "111"
		Arch = "arm64"

		expected := version.Version{
			Canonical: "v0.0.0",
			Metadata: version.Metadata{
				Architecture:    "arm64",
				BuildDate:       "111",
				Compiler:        "gc",
				GitCommit:       "000",
				GoVersion:       "go1.20.7",
				OperatingSystem: "darwin",
			},
		}
		assert.Equal(t, expected, *New())
	})
	t.Run("invalid tag", func(t *testing.T) {
		// to avoid flaky tests, set the values
		Go = "go1.20.7"
		OS = "darwin"
		Compiler = "gc"
		Tag = "something"
		Commit = "000"
		Date = "111"
		Arch = "arm64"

		expected := version.Version{
			Canonical: "something",
			Metadata: version.Metadata{
				Architecture:    "arm64",
				BuildDate:       "111",
				Compiler:        "gc",
				GitCommit:       "000",
				GoVersion:       "go1.20.7",
				OperatingSystem: "darwin",
			},
		}
		assert.Equal(t, expected, *New())
	})
}
