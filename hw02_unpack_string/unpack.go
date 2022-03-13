package hw02unpackstring

import (
	"errors"
	"strings"
	"unicode"
	"unicode/utf8"
)

var ErrInvalidString = errors.New("invalid string")

type complexRune struct {
	rune       rune
	repeatable bool
	repeats    int
	escaped    bool
}

func Unpack(str string) (string, error) {
	var (
		builder  strings.Builder
		lastRune *complexRune
	)

	runes := make([]*complexRune, 0, utf8.RuneCountInString(str))

	for _, char := range str {
		if unicode.IsDigit(char) && runeIsNotRepeatable(lastRune) {
			return "", ErrInvalidString
		}

		if runeIsNotEscaped(lastRune) && isEscape(lastRune.rune) && isNotValidForEscaping(char) {
			return "", ErrInvalidString
		}

		if runeIsNotEscaped(lastRune) && isEscape(lastRune.rune) {
			r := &complexRune{rune: char, repeatable: true, repeats: 1, escaped: true}
			lastRune = r
			runes = append(runes, r)
			continue
		} else {
			if isEscape(char) {
				lastRune = &complexRune{rune: char, repeatable: true}
				continue
			}

			if unicode.IsLetter(char) || unicode.IsControl(char) {
				r := &complexRune{rune: char, repeatable: true, repeats: 1}
				lastRune = r
				runes = append(runes, r)
				continue
			}

			if unicode.IsDigit(char) {
				repeats := int(char - '0')
				lastRune.repeats = repeats
				lastRune = &complexRune{rune: char}
				continue
			}
		}
	}

	for _, r := range runes {
		if r.repeats == 0 {
			continue
		}
		if r.repeats > 1 {
			repeatedRune := strings.Repeat(string(r.rune), r.repeats)
			builder.WriteString(repeatedRune)
		} else {
			builder.WriteRune(r.rune)
		}
	}

	return builder.String(), nil
}

func isEscape(r rune) bool {
	return r == '\\'
}

func isNotValidForEscaping(r rune) bool {
	return !unicode.IsDigit(r) && !isEscape(r) && !unicode.IsControl(r)
}

func runeIsNotRepeatable(r *complexRune) bool {
	return r == nil || !r.repeatable
}

func runeIsNotEscaped(r *complexRune) bool {
	return r != nil && !r.escaped
}
