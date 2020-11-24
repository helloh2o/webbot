package core

import (
	"bufio"
	"log"
	"os"
	"strings"
)

var (
	accounts []User
)

func ReadAccounts() {
	f, err := os.OpenFile("./accounts.txt", os.O_APPEND|os.O_CREATE|os.O_RDWR, os.ModePerm)
	if err != nil {
		log.Printf("open accounts file error %s", err)
	}
	br := bufio.NewReader(f)
	for {
		line, _, err := br.ReadLine()
		if err != nil {
			break
		}
		ap := strings.Split(string(line), ":")
		if len(ap) == 2 {
			accounts = append(accounts, User{Username: ap[0], Password: ap[1]})
		}
	}
}
