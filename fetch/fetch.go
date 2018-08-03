// Fetch exibe o conteúdo encontrado em cada URL especificada.
package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

func main() {
	start := time.Now()
	for _, url := range os.Args[1:] {
		newURL := checkURLPrefix(url)

		resp, err := http.Get(newURL)
		if err != nil {
			fmt.Fprintf(os.Stderr, "fetch: %v\n", err)
			continue
		}

		fmt.Fprintf(os.Stdout, "fetch: request status: %s\n", resp.Status)

		_, err = io.Copy(os.Stdout, resp.Body)
		resp.Body.Close()
		if err != nil {
			fmt.Fprintf(os.Stderr, "fetch: reading %s: %v\n", url, err)
			os.Exit(1)
		}
	}
	end := time.Since(start).Seconds()
	fmt.Printf("\nTempo total: %.2f\n", end)
}

func checkURLPrefix(url string) string {
	if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		return "http://" + url
	}
	return url
}