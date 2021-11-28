package builder

import "strings"

func New() *Builder {
	return &Builder{}
}

type Builder struct {
	parts []string
}

func (b *Builder) Reset() {
	b.parts = []string{}
}

func (b *Builder) Add(part string) {
	parts := strings.Split(part, " ")
	for _, part := range parts {
		if strings.Trim(part, " ") != "" {
			b.parts = append(b.parts, part)
		}
	}
}

func (b *Builder) Get() []string {
	return b.parts
}

func (b *Builder) String() string {
	return strings.Join(b.parts, " ")
}
