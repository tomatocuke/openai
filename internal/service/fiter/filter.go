package fiter

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"sync"

	"github.com/tomatocuke/sieve"
)

var (
	instance = sieve.New()
)

var (
	dir       = "./keyword/"
	filenames = []string{"色情辱骂.txt", "政治.txt", "违法广告.txt"}
)

func init() {
	go func() {
		var ok bool = true
		wg := &sync.WaitGroup{}
		wg.Add(len(filenames))
		for _, name := range filenames {
			go func(name string, wg *sync.WaitGroup) {
				defer wg.Done()
				name = dir + name

				f, err := os.OpenFile(name, os.O_RDONLY, 0755)
				if err != nil {
					ok = false
					return
				}

				var arr []string
				br := bufio.NewReader(f)
				for {
					a, _, c := br.ReadLine()
					if c == io.EOF {
						break
					}
					arr = append(arr, string(a))
				}

				instance.Add(arr)

			}(name, wg)
		}
		wg.Wait()

		if ok {
			fmt.Println("敏感词加载完成")
		} else {
			fmt.Println("加载词典失败，不影响使用。")
		}
	}()
}

func Check(text string) string {
	s, _ := instance.Search(text)
	if s != "" {
		log.Println("敏感词:", s, text)
	}
	return s
}
