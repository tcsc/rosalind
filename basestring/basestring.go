package basestring

import "fmt"

type baseString struct {
	chars  []byte
	length int
}

// ----------------------------------------------------------------------------
//
// ----------------------------------------------------------------------------

type InvalidBaseError rune

func (self InvalidBaseError) Error() string {
	return fmt.Sprintf("Invalid base char: %c", self)
}

// ----------------------------------------------------------------------------
//
// ----------------------------------------------------------------------------

func New() baseString {
	return baseString{
		chars:  make([]byte, 0),
		length: 0,
	}
}

func toBaseChar(c rune) (byte, error) {
	switch c {
	case 'G':
		return 1, nil
	case 'A':
		return 2, nil
	case 'T':
		return 3, nil
	case 'C':
		return 4, nil
	}
	return 0, InvalidBaseError(c)
}

func FromString(s string) (baseString, error) {
	cb := (len(s) + 1) / 2

	str := baseString{
		chars:  make([]byte, cb),
		length: len(s),
	}

	i := 0
	for _, c := range s {
		if err := str.setBase(i, c); err != nil {
			return baseString{nil, 0}, err
		}
		i++
	}
	return str, nil
}

func (self *baseString) Length() int {
	return self.length
}

func (self *baseString) setBase(i int, b rune) error {
	base, err := toBaseChar(b)
	if err != nil {
		return err
	}

	byteOffset := i / 2
	nibbleOffset := uint(i % 2)
	pair := self.chars[byteOffset]

	// surely there's a way we can do this without branching
	if nibbleOffset == 0 {
		pair = (pair & 0xF0) | byte(base)
	} else {
		pair = (pair & 0x0F) | byte(base<<4)
	}

	self.chars[byteOffset] = pair
	return nil
}
