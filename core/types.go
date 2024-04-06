package core

import (
	"errors"
	"regexp"
)

const (
	ColorInfo    = 0x003DFF
	ColorResult  = 0xAE00FF
	ColorSuccess = 0x23DB08
	ColorWarning = 0xFFD700
	ColorError   = 0xFF003D
)

var (
	ErrInvalidIdentifier = errors.New("invalid identifier")
	NamespaceRegex       = regexp.MustCompile(`[a-z0-9-]{1,32}`)
	IdRegex              = regexp.MustCompile(`[a-z0-9-]{1,32}(\/[a-z0-9-]{1,32})*`)
)

type Identifier struct {
	namespace string
	id        string
}

func NewIdentifier(namespace, id string) *Identifier {
	if !IsValidNamespace(namespace) || !IsValidId(id) {
		panic(ErrInvalidIdentifier)
	}

	return &Identifier{
		namespace: namespace,
		id:        id,
	}
}

func (i *Identifier) Namespace() string {
	return i.namespace
}

func (i *Identifier) Id() string {
	return i.id
}

func (i *Identifier) String() string {
	return i.namespace + ":" + i.id
}

func IsValidNamespace(namespace string) bool {
	return NamespaceRegex.MatchString(namespace)
}

func IsValidId(id string) bool {
	return IdRegex.MatchString(id)
}
