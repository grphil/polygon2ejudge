package console

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func YesOrNo(message string) bool {
	fmt.Printf("%s (y/n)\n", message)
	reader := bufio.NewReader(os.Stdin)
	text, _ := reader.ReadString('\n')
	return strings.ToLower(strings.TrimSpace(text)) == "y"
}
