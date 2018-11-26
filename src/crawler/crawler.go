package main

import (
	"fmt"
	"sync"
)

type Fetcher interface {
	// Fetch returns the body of URL and
	// a slice of URLs found on that page.
	//对单个网页的抓取逻辑
	Fetch(url string) (body string, urls []string, err error)
}
//保存抓取过的url，起到过滤作用
type UrlState struct{
	mutex sync.Mutex
	filter map[string]bool
}

// Crawl uses fetcher to recursively crawl
// pages starting with url, to a maximum of depth.
//递归爬取，与Urlstate一起使用，内部使用sync.WaitGroup来实现多goroutine的等待
func Crawl(url string, depth int, fetcher Fetcher, urlState UrlState) {
	if depth <= 0 {
		return
	}
	urlState.mutex.Lock()
	if urlState.filter[url]{
		defer urlState.mutex.Unlock()
		return 
	}
	
	urlState.filter[url] = true

	urlState.mutex.Unlock()
	body, urls, err := fetcher.Fetch(url)
	if err != nil {
		urlState.mutex.Lock()
		fmt.Println(err)
		delete(urlState.filter, url)
		defer urlState.mutex.Unlock()
		return
	}
	fmt.Printf("found: %s %q\n", url, body)
	var wg sync.WaitGroup
	for _, u := range urls {
		wg.Add(1)
		go func(u string){
			defer wg.Done()
			Crawl(u, depth-1, fetcher, urlState)
		}(u)
	}
	wg.Wait()
	return
}

//使用channel来实现的并行爬虫
func ConcurrentChannel(url string, fetcher Fetcher){
	ch := make(chan []string)
	go func(){
		ch <- []string{url}
	}()
	master(ch, fetcher)
}

func master(ch chan []string, fetcher Fetcher){
	n := 1//起始放入的种子个数
	fetched := make(map[string]bool)//用于过滤网页
	for urls := range ch{
		for _, u := range urls {
			if fetched[u] == false {
				fetched[u] = true
				n += 1//每开一个worker，爬一个url，n++
				go worker(u, ch, fetcher)
			}
		}
		n -= 1//执行完一个url的，n--
		if n == 0 {
			return 
		}
	}
}

func worker(url string, ch chan []string, fetcher Fetcher) {
	body, urls, err := fetcher.Fetch(url)
	if err != nil {
		fmt.Println(err)
		ch <- []string{}
	}else{
		fmt.Printf("found: %s %q\n", url, body)
		ch <- urls
	}
}


func main() {
	urlState := UrlState{sync.Mutex{}, make(map[string]bool)}
	Crawl("https://golang.org/", 4, fetcher, urlState)
	fmt.Println("finish!")
	ConcurrentChannel("https://golang.org/", fetcher)
}




// fakeFetcher is Fetcher that returns canned results.
type fakeFetcher map[string]*fakeResult

type fakeResult struct {
	body string
	urls []string
}

func (f fakeFetcher) Fetch(url string) (string, []string, error) {
	if res, ok := f[url]; ok {
		return res.body, res.urls, nil
	}
	return "", nil, fmt.Errorf("not found: %s", url)
}

// fetcher is a populated fakeFetcher.
var fetcher = fakeFetcher{
	"https://golang.org/": &fakeResult{
		"The Go Programming Language",
		[]string{
			"https://golang.org/pkg/",
			"https://golang.org/cmd/",
		},
	},
	"https://golang.org/pkg/": &fakeResult{
		"Packages",
		[]string{
			"https://golang.org/",
			"https://golang.org/cmd/",
			"https://golang.org/pkg/fmt/",
			"https://golang.org/pkg/os/",
		},
	},
	"https://golang.org/pkg/fmt/": &fakeResult{
		"Package fmt",
		[]string{
			"https://golang.org/",
			"https://golang.org/pkg/",
		},
	},
	"https://golang.org/pkg/os/": &fakeResult{
		"Package os",
		[]string{
			"https://golang.org/",
			"https://golang.org/pkg/",
		},
	},
}
