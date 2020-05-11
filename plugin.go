package uad

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/gobuffalo/flect"
)

type Plugin struct {
	URL         string
	Name        string
	Description string
	Category    string
	Owned       bool
}

func (p Plugin) String() string {
	cl := flect.Pascalize(p.Category)
	bb := &bytes.Buffer{}

	bb.WriteString(fmt.Sprintf("// %s\n", p.Name))
	for _, l := range strings.Split(p.Description, "\n") {
		l = strings.TrimSpace(l)
		if len(l) == 0 {
			continue
		}
		bb.WriteString(fmt.Sprintf("// %s\n", l))
	}
	fn := fmt.Sprintf("func (%s) %s(){}\n\n", cl, flect.Pascalize(p.Name))
	bb.WriteString(fn)

	return bb.String()
}
