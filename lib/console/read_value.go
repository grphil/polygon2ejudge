package console

import (
	"bufio"
	"fmt"
	"golang.org/x/term"
	"os"
	"strconv"
	"strings"
)

func ReadValue(message string) string {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println(message)
	res, err := reader.ReadString('\n')
	if err != nil {
		panic(err)
	}
	return strings.TrimSpace(res)
}

func ReadSecret(message string) string {
	fmt.Println(message)
	secret, err := term.ReadPassword(int(os.Stdin.Fd()))
	if err != nil {
		panic(err)
	}
	return strings.TrimSpace(string(secret))
}

func ReadInts(message string) []int {
	ids := ReadValue(message)

	var res []int
	for _, idStr := range strings.Fields(ids) {
		id, err := strconv.Atoi(idStr)
		if err != nil {
			panic(err)
		}
		res = append(res, id)
	}
	return res
}
