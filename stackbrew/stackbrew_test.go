package stackbrew

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseReader(t *testing.T) {
	t.Run("golang", func(t *testing.T) {
		f, err := os.Open("testdata/golang")
		if err != nil {
			t.Fatalf("open file: %v", err)
		}
		defer f.Close()

		s := ParseReader(f)

		assert.Equal(t, "https://github.com/docker-library/golang.git", s.GitRepo)
		assert.ElementsMatch(t, []string{"1.17.0-bullseye", "1.17-bullseye", "1-bullseye", "bullseye"}, s.Stacks[0].Tags)
		assert.ElementsMatch(t, []string{"1.17.0", "1.17", "1", "latest"}, s.Stacks[0].SharedTags)
		assert.Equal(t, "48a7371ed6055a97a10adb0b75756192ad5f1c97", s.Stacks[0].GitCommit)
	})

	t.Run("adoptopenjdk", func(t *testing.T) {
		f, err := os.Open("testdata/adoptopenjdk")
		if err != nil {
			t.Fatalf("open file: %v", err)
		}
		defer f.Close()

		s := ParseReader(f)
		assert.Equal(t, "https://github.com/AdoptOpenJDK/openjdk-docker.git", s.GitRepo)
	})
}
