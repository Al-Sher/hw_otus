package hw10programoptimization

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"regexp"
	"strings"

	"github.com/mailru/easyjson"
)

type DomainStat map[string]int

var ErrNilReader = errors.New("invalid reader")

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	u, err := getUsers(r, domain)
	if err != nil {
		return nil, fmt.Errorf("get users error: %w", err)
	}
	return countDomains(u)
}

type users []User

func getUsers(r io.Reader, domain string) (result users, err error) {
	if r == nil {
		return nil, ErrNilReader
	}

	result = make(users, 0)

	exp, err := regexp.Compile("\\." + domain)
	if err != nil {
		return
	}

	sc := bufio.NewScanner(r)

	for sc.Scan() {
		var user User
		if err = easyjson.Unmarshal(sc.Bytes(), &user); err != nil {
			return
		}
		if exp.MatchString(user.Email) {
			result = append(result, user)
		}
	}
	return
}

func countDomains(u users) (DomainStat, error) {
	result := make(DomainStat, len(u))

	for _, user := range u {
		if i := strings.Index(user.Email, "@"); i != -1 {
			result[strings.ToLower(user.Email[i+1:])]++
		}
	}

	return result, nil
}
