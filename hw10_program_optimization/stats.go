package hw10programoptimization

import (
	"bufio"
	"errors"
	"io"
	"strings"

	jsoniter "github.com/json-iterator/go"
)

type User struct {
	ID       int
	Name     string
	Username string
	Email    string
	Phone    string
	Password string
	Address  string
}

// UserEmail we don't need every field for parsing email addresses.
type UserEmail struct {
	Email string
}

type DomainStat map[string]int

var errInvalidEmail = errors.New("invalid email")

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	var (
		scanner = bufio.NewScanner(r)
		json    = jsoniter.ConfigFastest
		stat    = make(DomainStat)
		user    UserEmail
		dmn     string
		err     error
	)
	for scanner.Scan() {
		user.Email = ""
		if err := json.Unmarshal(scanner.Bytes(), &user); err != nil {
			return nil, err
		}
		if strings.HasSuffix(user.Email, domain) {
			if dmn, err = extractDomainFromEmail(user.Email); err != nil {
				return nil, err
			}
			stat[dmn]++
		}
	}
	return stat, nil
}

func extractDomainFromEmail(email string) (string, error) {
	var (
		builder  strings.Builder
		inDomain = false
	)
	for _, r := range email {
		if r == '@' && inDomain {
			return "", errInvalidEmail
		}
		if r == '@' {
			inDomain = true
			continue
		}
		if inDomain {
			builder.WriteRune(r)
		}
	}
	return strings.ToLower(builder.String()), nil
}
