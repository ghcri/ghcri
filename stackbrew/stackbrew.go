package stackbrew

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"strings"
)

type Stackbrew struct {
	Maintainers []string
	GitRepo     string

	Stacks []Stack
}

type Stack struct {
	Tags          []string
	SharedTags    []string
	Architectures []string
	GitCommit     string
	File          string
	Directory     string
	Constraints   []string
}

func ParseBytes(bs []byte) Stackbrew {
	return ParseReader(bytes.NewReader(bs))
}

func ParseReader(r io.Reader) Stackbrew {
	s := Stackbrew{}

	var cur *Stack
	sc := bufio.NewScanner(r)
	for sc.Scan() {
		content := sc.Text()

		// Ignore empty line.
		if len(content) == 0 {
			continue
		}

		// Ignore all content start with `#`
		if strings.HasPrefix(content, "#") {
			continue
		}

		// Parse Maintainers.
		//
		// The space padding is hardcoded.
		//
		// Example content looks like:
		//
		// Maintainers: Tianon Gravi <admwiggin@gmail.com> (@tianon),
		//             Joseph Ferguson <yosifkit@gmail.com> (@yosifkit),
		//             Johan Euphrosine <proppy@google.com> (@proppy)
		if strings.HasPrefix(content, "Maintainers") || strings.HasPrefix(content, "             ") {
			content = strings.TrimPrefix(content, "Maintainers:")
			content = strings.TrimSuffix(content, ",")
			content = strings.TrimSpace(content)
			s.Maintainers = append(s.Maintainers, content)
			continue
		}

		if strings.HasPrefix(content, "GitRepo") {
			s.GitRepo = parseLine(content, "GitRepo")
			continue
		}

		// Tags is always the first field of Stack, we use it to detect whether the previous has finished.
		if strings.HasPrefix(content, "Tags") {
			// Check the previous stack.
			// If the cur is not nil, we should append it into Stackbrew and reset it.
			if cur != nil {
				s.Stacks = append(s.Stacks, *cur)
				cur = nil
			}
			cur = &Stack{}
			cur.Tags = parseSlice(content, "Tags")
			continue
		}

		if strings.HasPrefix(content, "SharedTags") {
			cur.SharedTags = parseSlice(content, "SharedTags")
			continue
		}

		if strings.HasPrefix(content, "Architectures") {
			cur.Architectures = parseSlice(content, "Architectures")
			continue
		}

		if strings.HasPrefix(content, "GitCommit") {
			cur.GitCommit = parseLine(content, "GitCommit")
			continue
		}

		if strings.HasPrefix(content, "File") {
			cur.File = parseLine(content, "File")
			continue
		}

		if strings.HasPrefix(content, "Directory") {
			cur.Directory = parseLine(content, "Directory")
			continue
		}

		if strings.HasPrefix(content, "Constraints") {
			cur.Constraints = parseSlice(content, "Constraints")
			continue
		}

		panic(fmt.Errorf("line %s is not parsed correctly", content))
	}
	// Handle the last item of stack.
	if cur != nil {
		s.Stacks = append(s.Stacks, *cur)
	}
	return s
}

func parseLine(content, prefix string) string {
	content = strings.TrimPrefix(content, prefix+":")
	return strings.TrimSpace(content)
}

func parseSlice(content, prefix string) []string {
	content = strings.TrimPrefix(content, prefix+":")
	vs := strings.Split(content, ",")
	for i := range vs {
		vs[i] = strings.TrimSpace(vs[i])
	}
	return vs
}
