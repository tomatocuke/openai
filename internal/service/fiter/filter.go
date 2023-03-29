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

	go func() {
		// 加载预定义词典
		arr := strings.Split(keywords, "\n")
		var builder strings.Builder
		builder.Grow(20)

		for _, v := range arr {
			runes := strings.Split(v, " ")
			for _, w := range runes {
				i, _ := strconv.Atoi(w)
				builder.WriteRune(rune(i))
			}
			arr = append(arr, builder.String())
			builder.Reset()
		}

		instance.Add(arr)

		// 加载你定义的词典
		f, err := os.Open("./keyword.txt")
		if err != nil {
			return
		}
		arr = arr[:0]
		br := bufio.NewReader(f)
		for {
			a, _, c := br.ReadLine()
			if c == io.EOF {
				break
			}
			arr = append(arr, string(a))
		}

		instance.Add(arr)
		fmt.Println("敏感词词典加载完成")
	}()

}

func Check(text string) bool {
	s, _ := instance.Search(text)
	if s != "" {
		log.Println("检测到敏感词:", s)
	}
	return s == ""
}
