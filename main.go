package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"

	core "WebAppAnalyzer/core"
)

// Data structure to hold technology information
type Technology struct {
	Name     string `json:"name"`
	Version  string `json:"version"`
	Category string `json:"category"`
	Icon     string `json:"icon"`
}

func AnalyzeWebApp() (result interface{}) {

	var url, appsJSONPath, scraper, userAgent string
	var help bool
	var timeoutSeconds, loadingTimeoutSeconds, maxDepth, maxVisitedLinks, msDelayBetweenRequests int
	flag.StringVar(&appsJSONPath, "file", "", "Path to override default technologies.json file")
	flag.StringVar(&scraper, "scraper", "rod", "Choose scraper between rod (default) and colly")
	flag.StringVar(&userAgent, "useragent", "", "Override the user-agent string")
	flag.IntVar(&timeoutSeconds, "timeout", 3, "Timeout in seconds for fetching the url")
	flag.IntVar(&loadingTimeoutSeconds, "loadtimeout", 3, "Timeout in seconds for loading the page")
	flag.IntVar(&maxDepth, "depth", 0, "Don't analyze page when depth superior to this number. Default (0) means no recursivity (only first page will be analyzed)")
	flag.IntVar(&maxVisitedLinks, "maxlinks", 5, "Max number of pages to visit. Exit when reached")
	flag.IntVar(&msDelayBetweenRequests, "delay", 100, "Delay in ms between requests")
	flag.BoolVar(&help, "h", false, "Help")
	flag.Parse()

	var Usage = func() {
		fmt.Fprintln(os.Stderr, "Usage : go run main.go [options] <url>")
		flag.PrintDefaults()
	}

	if help {
		Usage()
		os.Exit(1)
	}
	if flag.NArg() == 0 {
		fmt.Fprintln(os.Stderr, "You must specify a url to analyse")
		Usage()
		os.Exit(1)
	} else if flag.NArg() > 1 {
		fmt.Fprintf(os.Stderr, "Too many arguments %s", flag.Args())
		Usage()
		os.Exit(1)
	} else {
		url = flag.Arg(0)
	}
	if scraper != "rod" && scraper != "colly" {
		fmt.Fprintf(os.Stderr, "Unknown scraper %s : only supporting rod and colly", scraper)
		Usage()
		os.Exit(1)
	}

	config := core.NewConfig()
	config.AppsJSONPath = appsJSONPath
	config.TimeoutSeconds = timeoutSeconds
	config.LoadingTimeoutSeconds = loadingTimeoutSeconds
	config.MaxDepth = maxDepth
	config.MaxVisitedLinks = maxVisitedLinks
	config.MsDelayBetweenRequests = msDelayBetweenRequests
	config.Scraper = scraper
	if userAgent != "" {
		config.UserAgent = userAgent
	}

	wapp, err := core.Init(config)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	res, err := wapp.Analyze(url)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	return (res)
}

func main() {
	e := echo.New()

	// Serve static files
	e.Static("/", "public")

	data := AnalyzeWebApp()
	e.POST("/api/analyze", func(c echo.Context) error {
		return c.JSON(http.StatusOK, data)
	})

	// Start the server on port 8080
	e.Logger.Fatal(e.Start(":8080"))
}
