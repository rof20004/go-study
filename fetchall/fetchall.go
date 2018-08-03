// Fetchall busca URLs em paralelo e informa os tempos gastos e os tamanhos.
package main

import (
	"fmt"
	"time"
	"io"
	"os"
	"net/http"
	"strings"
)

func main() {
	start := time.Now()
	ch := make(chan string)

	for _, url := range os.Args[1:] {
		go fetch(url, ch)
	}

	for range os.Args[1:] {
		fmt.Println(<-ch)
	}

	fmt.Printf("%.2f elapsed\n", time.Since(start).Seconds())
}

func fetch(url string, ch chan<- string) {
	start := time.Now()
	response, err := http.Get(url)
	if err != nil {
		ch <- fmt.Sprint(err)
		return
	}

	newURL := createFileNameByURL(url)

	file, err := os.OpenFile(newURL + ".txt", os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		ch <- fmt.Sprint(err)
		return
	}

	nbytes, err := io.Copy(file, response.Body)
	file.Close()
	response.Body.Close()
	if err != nil {
		ch <- fmt.Sprintf("while reading %s: %v", url, err)
		return
	}

	secs := time.Since(start).Seconds()
	ch <- fmt.Sprintf("%.2f %7d %s", secs, nbytes, url)
}

func createFileNameByURL(url string) string {
	if strings.HasPrefix(url, "http://") {
		return strings.TrimPrefix(url, "http://")
	} else if strings.HasPrefix(url, "https://") {
		return strings.TrimPrefix(url, "https://")
	}
	return url
}