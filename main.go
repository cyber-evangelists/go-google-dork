package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"github.com/gocolly/colly"
)

const (
	userAgent   = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.3"
	googleURL   = "https://www.google.com/search?q="
	bingURL     = "https://www.bing.com/search?q="
	baiduURL    = "https://www.baidu.com/s?wd="
	yahooURL    = "https://search.yahoo.com/search?p="
	dorksDomain = "dorks-domain.txt"
	dorksTarget = "dorks-target.txt"
	rateLimit   = 5 * time.Second
)

func main() {
	domainFlag := flag.String("d", "", "Set the target domain")
	subdomainFlag := flag.Bool("s", false, "Search also for subdomains")
	targetFlag := flag.String("t", "", "Set the target (general)")
	versionFlag := flag.Bool("v", false, "Show the version of evildork")

	flag.Parse()

	if *versionFlag {
		version()
		return
	}

	if *domainFlag != "" {
		startDorking(*domainFlag, *subdomainFlag, "domain")
	} else if *targetFlag != "" {
		startDorking(*targetFlag, false, "target")
	} else {
		flag.PrintDefaults()
	}
}

// version prints the current version of the program
func version() {
	fmt.Println("v0.1")
}

// startDorking initiates the dorking process
func startDorking(target string, subdomain bool, targetType string) {
	directory := "output-dorking"
	if subdomain {
		target = "*." + target
	}
	filename := createOutputFile(directory, target)

	addHTMLBanner(filename)
	addDorks(target, filename, targetType)
	addHTMLFooter(filename)

	openInBrowser(filename)
	performWebScraping(target, targetType)
}

// createOutputFolder creates the output directory if it doesn't exist
func createOutputFolder(directory string) {
	if _, err := os.Stat(directory); os.IsNotExist(err) {
		err := os.MkdirAll(directory, os.ModePerm)
		if err != nil {
			fmt.Println("Error creating directory:", err)
			os.Exit(1)
		}
	}
}

// createOutputFile creates an output file and handles potential overwriting
func createOutputFile(directory, target string) string {
	createOutputFolder(directory)
	target = strings.ReplaceAll(target, " ", "-")
	filename := fmt.Sprintf("%s/dorking-%s.html", directory, target)

	_, err := os.Stat(filename)
	if err == nil {
		var choice string
		fmt.Printf("[!] %s already exists. Do you want to overwrite? (y/n): ", filename)
		fmt.Scanln(&choice)
		if strings.ToLower(choice) != "y" {
			os.Exit(1)
		}
	}

	file, err := os.Create(filename)
	if err != nil {
		fmt.Println("Error creating file:", err)
		os.Exit(1)
	}
	defer file.Close()

	return filename
}

// addHTMLBanner adds a banner to the HTML output file
func addHTMLBanner(filename string) {
	htmlBanner := `
    <html>
    <head>
    <link rel="stylesheet" href="../style.css">
    <title> WEBDORKING </title>
    </head>
    <body>
    <div class="topnav">
    <a href="https://github.com/HamzaSajjad141"> Hamza Team</a>
    <a href="https://github.com/HamzaSajjad141">Contribute to Webdorking</a>
    </div>
    <h1 class="evildork">WEBDORKING</h1>
    <h3 class="fricciolosa">by Hamza Sajjad </h3>
    <ul>
    `

	err := ioutil.WriteFile(filename, []byte(htmlBanner), 0644)
	if err != nil {
		fmt.Println("Error writing to file:", err)
		os.Exit(1)
	}
}

// addHTMLFooter adds a footer to the HTML output file
func addHTMLFooter(filename string) {
	htmlFooter := `
    </ul>
    <div class="footer">
        <p>webdorking by <a href='https://github.com/HamzaSajjad141'>Hamza Sajjad Team</a></p>
    </div>
    <br><br><br><br><br><br><br><br>
    </body>
    </html>
    `

	f, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Error opening file:", err)
		os.Exit(1)
	}
	defer f.Close()

	if _, err := f.WriteString(htmlFooter); err != nil {
		fmt.Println("Error writing to file:", err)
		os.Exit(1)
	}
}

// addDorks adds dorks to the HTML output file
func addDorks(target, filename, targetType string) {
	var dorksFile string
	if targetType == "domain" {
		dorksFile = dorksDomain
	} else {
		dorksFile = dorksTarget
	}

	dorksContent, err := ioutil.ReadFile(dorksFile)
	if err != nil {
		fmt.Println("Error reading dorks file:", err)
		os.Exit(1)
	}

	dorks := strings.Split(string(dorksContent), "\n")

	f, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Error opening file:", err)
		os.Exit(1)
	}
	defer f.Close()

	for _, dork := range dorks {
		for _, engine := range []string{googleURL, bingURL, baiduURL, yahooURL} {
			query := dork + " " + target
			encodedQuery := engine + url.QueryEscape(query)
			_, err := f.WriteString(fmt.Sprintf("<li><a href='%s' target='_blank'>%s</a></li>\n", encodedQuery, query))
			if err != nil {
				fmt.Println("Error writing to file:", err)
				os.Exit(1)
			}
		}
	}
}

// openInBrowser opens the output file in the default browser
func openInBrowser(filename string) {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("cmd", "/c", "start", filename)
	case "darwin":
		cmd = exec.Command("open", filename)
	case "linux":
		cmd = exec.Command("xdg-open", filename)
	default:
		fmt.Println("Unsupported operating system")
		return
	}

	err := cmd.Run()
	if err != nil {
		fmt.Println("Error opening file in browser:", err)
	}
}

// performWebScraping performs web scraping on multiple search engines
func performWebScraping(target, targetType string) {
	queries := generateQueries(target, targetType)
	c := colly.NewCollector(
		colly.UserAgent(userAgent),
		colly.Async(true),
	)

	c.Limit(&colly.LimitRule{
		DomainGlob:  "*",
		Parallelism: 4,
		Delay:       rateLimit,
	})

	c.OnHTML("a", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		if strings.HasPrefix(link, "/url?q=") {
			link = strings.TrimPrefix(link, "/url?q=")
			link = strings.Split(link, "&")[0]
			fmt.Println(link)
		}
	})

	for _, query := range queries {
		go func(query string) {
			err := c.Visit(query)
			if err != nil {
				fmt.Println("Error visiting URL:", err)
			}
		}(query)
	}

	c.Wait()
}

// generateQueries generates search engine queries for the target
func generateQueries(target, targetType string) []string {
	var queries []string
	searchEngines := []string{googleURL, bingURL, baiduURL, yahooURL}
	for _, engine := range searchEngines {
		query := "site:" + target
		if targetType == "target" {
			query = target
		}
		queries = append(queries, engine+url.QueryEscape(query))
	}
	return queries
}
