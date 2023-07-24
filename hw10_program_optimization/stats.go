package hw10programoptimization

import (
	"bufio"
	"errors"
	"fmt"
	"io"
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
	return u, nil
}

func getUsers(r io.Reader, domain string) (result DomainStat, err error) {
	if r == nil {
		return nil, ErrNilReader
	}

	result = make(DomainStat, 0)

	sc := bufio.NewScanner(r)

	for sc.Scan() {
		var user User
		if err = easyjson.Unmarshal(sc.Bytes(), &user); err != nil {
			return
		}
		if strings.HasSuffix(user.Email, "."+domain) {
			if i := strings.Index(user.Email, "@"); i != -1 {
				result[strings.ToLower(user.Email[i+1:])]++
			}
		}
	}
	return
}
