package hw03frequencyanalysis

import (
	"sort"
	"strings"
	"unicode"
)

type word struct {
	word  string
	count int
}

func Top10(text string) []string {
	if text == "" {
		return make([]string, 0)
	}

	text = clearText(text)
	textWords := strings.Split(text, " ")
	counts := make(map[string]word)

	for _, tw := range textWords {
		w := clearWord(tw)
		if _, ok := counts[w]; ok {
			counts[w] = word{
				word:  w,
				count: counts[w].count + 1,
			}
		} else {
			counts[w] = word{
				word:  w,
				count: 1,
			}
		}
	}

	words := make([]word, 0, len(counts))
	for _, w := range counts {
		words = append(words, w)
	}

	sort.Slice(words, func(i, j int) bool {
		if words[i].count == words[j].count {
			return strings.Compare(words[i].word, words[j].word) == -1
		}
		if words[i].count > words[j].count {
			return true
		}
		return false
	})

	boundary := 10
	if len(words) < 10 {
		boundary = len(words)
	}

	words = words[:boundary]
	result := make([]string, 0, boundary)

	for _, w := range words {
		result = append(result, w.word)
	}

	return result
}

func clearText(t string) string {
	var (
		builder         strings.Builder
		lastRuneIsSpace bool
	)

	t = strings.ReplaceAll(t, "\t", " ")
	t = strings.ReplaceAll(t, "\n", " ")
	t = strings.ReplaceAll(t, " - ", " ")

	for _, r := range t {
		if unicode.IsSpace(r) && lastRuneIsSpace {
			continue
		}

		if unicode.IsSpace(r) && !lastRuneIsSpace {
			lastRuneIsSpace = true
		} else {
			lastRuneIsSpace = false
		}

		builder.WriteRune(r)
	}

	return builder.String()
}

func clearWord(w string) string {
	var builder strings.Builder

	for i, r := range w {
		if r == '-' {
			if 0 < i && i < len(w)-1 {
				builder.WriteRune(r)
			}
			continue
		}

		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			builder.WriteRune(r)
			continue
		}
	}

	return strings.ToLower(builder.String())
}
