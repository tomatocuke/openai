package fiter

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/tomatocuke/sieve"
)

var (
	instance = sieve.New()
)

func init() {
	f, err := os.Open("./keyword.txt")
	if err != nil {
		fmt.Println("加载词典失败")
		return
	}
	var arr []string
	var builder strings.Builder
	builder.Grow(10)
	br := bufio.NewReader(f)
	for {
		a, _, c := br.ReadLine()
		if c == io.EOF {
			break
		}

		runes := strings.Split(string(a), " ")

		for _, w := range runes {
			i, _ := strconv.Atoi(w)
			builder.WriteRune(rune(i))
		}
		arr = append(arr, builder.String())
		builder.Reset()
	}

	instance.Add(arr)

}

func Check(text string) bool {
	s, _ := instance.Search(text)
	if s != "" {
		log.Println("敏感词:", s, text)
	}
	return s == ""
}
