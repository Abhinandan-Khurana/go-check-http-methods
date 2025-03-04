package main

import (
	"bufio"
	"crypto/tls"
	"encoding/json"
	"encoding/xml"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/fatih/color"
)

const (
	appName    = "go-check-http-methods"
	appVersion = "1.0.0"
	appAuthor  = "Abhinandan Khurana aka @l0u51f3r007"
)

var defaultMethods = []string{
	"GET", "POST", "PUT", "DELETE", "HEAD", "OPTIONS", "PATCH",
	"TRACE", "CONNECT", "PROPFIND", "PROPPATCH", "MKCOL", "COPY",
	"MOVE", "LOCK", "UNLOCK", "PURGE", "LINK", "UNLINK",
}

var dangerousMethods = map[string]bool{
	"PUT":       true,
	"DELETE":    true,
	"TRACE":     true,
	"PROPFIND":  true,
	"PROPPATCH": true,
	"MKCOL":     true,
	"COPY":      true,
	"MOVE":      true,
	"LOCK":      true,
	"UNLOCK":    true,
}

type options struct {
	url             string
	urlFile         string
	outputFile      string
	outputFormat    string
	methodsFile     string
	concurrent      int
	timeout         int
	followRedirects bool
	insecure        bool
	verbose         bool
	quiet           bool
	silent          bool
	noColor         bool
	viewMode        string
	proxy           string
	userAgent       string
	headers         stringSlice
	cookies         stringSlice
	auth            string
}

// Custom type for collecting multiple flags of the same name
type stringSlice []string

func (s *stringSlice) String() string {
	return strings.Join(*s, ", ")
}

func (s *stringSlice) Set(value string) error {
	*s = append(*s, value)
	return nil
}

// Result for a single HTTP method test
type MethodResult struct {
	Method          string `json:"method" xml:"method"`
	StatusCode      int    `json:"status_code" xml:"status_code"`
	Status          string `json:"status" xml:"status"`
	ResponseTime    int64  `json:"response_time_ms" xml:"response_time_ms"`
	ContentLength   int    `json:"content_length" xml:"content_length"`
	IsDangerous     bool   `json:"is_dangerous" xml:"is_dangerous"`
	IsVulnerable    bool   `json:"is_vulnerable" xml:"is_vulnerable"`
	VulnDescription string `json:"vulnerability_description,omitempty" xml:"vulnerability_description,omitempty"`
}

// Result for a single URL test
type URLResult struct {
	URL     string         `json:"url" xml:"url"`
	Results []MethodResult `json:"results" xml:"results>method_result"`
}

// Overall results structure
type Results struct {
	ToolName    string      `json:"tool_name" xml:"tool_name"`
	ToolVersion string      `json:"tool_version" xml:"tool_version"`
	ToolAuthor  string      `json:"tool_author" xml:"tool_author"`
	Timestamp   string      `json:"timestamp" xml:"timestamp"`
	Results     []URLResult `json:"results" xml:"url_results>url_result"`
}

type XMLResults struct {
	XMLName xml.Name `xml:"http_method_test_results"`
	Results Results  `xml:",innerxml"`
}

func main() {
	opts := parseOptions()

	if opts.noColor {
		color.NoColor = true
	}

	if !opts.quiet && !opts.silent {
		printBanner()
	}

	var urls []string
	var methods []string
	var err error

	if opts.url != "" {
		urls = append(urls, opts.url)
	} else if opts.urlFile != "" {
		urls, err = readLinesFromFile(opts.urlFile)
		if err != nil {
			color.New(color.FgRed, color.Bold).Fprintf(os.Stderr, "Error reading URL file: %v\n", err)
			os.Exit(1)
		}
	} else {
		color.New(color.FgRed, color.Bold).Fprintf(os.Stderr, "Either URL or URL file must be specified\n")
		os.Exit(1)
	}

	// Normalize URLs
	for i, url := range urls {
		if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
			urls[i] = "https://" + url
		}
	}

	if opts.methodsFile != "" {
		methods, err = readLinesFromFile(opts.methodsFile)
		if err != nil {
			color.New(color.FgRed, color.Bold).Fprintf(os.Stderr, "Error reading methods file: %v\n", err)
			os.Exit(1)
		}
	} else {
		methods = defaultMethods
	}

	client := createHTTPClient(opts)

	results := runTests(urls, methods, client, opts)

	outputResults(results, opts)
}

func parseOptions() options {
	var opts options

	flag.StringVar(&opts.url, "u", "", "Single URL to test")
	flag.StringVar(&opts.urlFile, "f", "", "File containing URLs to test (one per line)")
	flag.StringVar(&opts.outputFile, "o", "", "Output file for results")
	flag.StringVar(&opts.outputFormat, "format", "txt", "Output format: txt, json, or xml")
	flag.StringVar(&opts.methodsFile, "m", "", "File containing HTTP methods to test")
	flag.IntVar(&opts.concurrent, "c", 10, "Number of concurrent requests")
	flag.IntVar(&opts.timeout, "t", 10, "Request timeout in seconds")
	flag.BoolVar(&opts.followRedirects, "L", false, "Follow redirects")
	flag.BoolVar(&opts.insecure, "k", false, "Allow insecure TLS connections")
	flag.BoolVar(&opts.verbose, "v", false, "Verbose output")
	flag.BoolVar(&opts.quiet, "q", false, "Quiet mode, no output except results")
	flag.BoolVar(&opts.silent, "silent", false, "Silent mode, only output results")
	flag.BoolVar(&opts.noColor, "nc", false, "No color output")
	flag.StringVar(&opts.viewMode, "view", "all", "View mode: all, enabled, vulnerable")
	flag.StringVar(&opts.proxy, "proxy", "", "Use proxy (format: http://host:port)")
	flag.StringVar(&opts.userAgent, "ua", "Go-HTTP-Method-Tester/"+appVersion, "User agent string")
	flag.StringVar(&opts.auth, "auth", "", "Basic authentication (format: username:password)")
	flag.Var(&opts.headers, "H", "Custom header (can be used multiple times, format: 'Name: Value')")
	flag.Var(&opts.cookies, "cookie", "Cookie to include (can be used multiple times, format: 'name=value')")

	flag.Parse()

	if opts.url == "" && opts.urlFile == "" {
		color.New(color.FgRed).Println("Error: Either a URL (-u) or a file containing URLs (-f) must be specified")
		flag.Usage()
		os.Exit(1)
	}

	validViewModes := map[string]bool{"all": true, "enabled": true, "vulnerable": true}
	if !validViewModes[opts.viewMode] {
		color.New(color.FgRed).Println("Error: Invalid view mode. Must be one of: all, enabled, vulnerable")
		flag.Usage()
		os.Exit(1)
	}

	validFormats := map[string]bool{"txt": true, "json": true, "xml": true}
	if !validFormats[opts.outputFormat] {
		color.New(color.FgRed).Println("Error: Invalid output format. Must be one of: txt, json, xml")
		flag.Usage()
		os.Exit(1)
	}

	return opts
}

func printBanner() {
	color.New(color.FgCyan, color.Bold).Printf(`
                         _               _         _     _   _                              _   _               _     
   __ _  ___         ___| |__   ___  ___| | __    | |__ | |_| |_ _ __        _ __ ___   ___| |_| |__   ___   __| |___ 
  / _' |/ _ \ _____ / __| '_ \ / _ \/ __| |/ /____| '_ \| __| __| '_ \ _____| '_ ' _ \ / _ \ __| '_ \ / _ \ / _' / __|
 | (_| | (_) |_____| (__| | | |  __/ (__|   <_____| | | | |_| |_| |_) |_____| | | | | |  __/ |_| | | | (_) | (_| \__ \
  \__, |\___/       \___|_| |_|\___|\___|_|\_\    |_| |_|\__|\__| .__/      |_| |_| |_|\___|\__|_| |_|\___/ \__,_|___/
  |___/                                                         |_|                                                   
`)
	fmt.Println()
}

func readLinesFromFile(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" && !strings.HasPrefix(line, "#") {
			lines = append(lines, line)
		}
	}

	return lines, scanner.Err()
}

func createHTTPClient(opts options) *http.Client {
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: opts.insecure},
		MaxIdleConns:    100,
		IdleConnTimeout: 90 * time.Second,
	}

	if opts.proxy != "" {
		proxyURL, err := http.ProxyFromEnvironment(&http.Request{URL: nil})
		if err != nil {
			color.New(color.FgRed, color.Bold).Fprintf(os.Stderr, "Warning: Invalid proxy URL: %v\n", err)
		} else {
			transport.Proxy = http.ProxyURL(proxyURL)
		}
	}

	client := &http.Client{
		Transport: transport,
		Timeout:   time.Duration(opts.timeout) * time.Second,
	}

	if !opts.followRedirects {
		client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		}
	}

	return client
}

func runTests(urls []string, methods []string, client *http.Client, opts options) Results {
	results := Results{
		ToolName:    appName,
		ToolVersion: appVersion,
		ToolAuthor:  appAuthor,
		Timestamp:   time.Now().Format(time.RFC3339),
		Results:     make([]URLResult, 0),
	}

	for _, url := range urls {
		if opts.verbose && !opts.quiet && !opts.silent {
			// color.New(color.FgHiYellow).Printf("Testing URL: %s\n", url)
			urlColor := color.New(color.FgGreen, color.Bold).SprintFunc()
			fmt.Printf("Testing URL: %s\n", urlColor(url))
		}

		urlResult := URLResult{
			URL:     url,
			Results: make([]MethodResult, 0),
		}

		sem := make(chan bool, opts.concurrent)
		var wg sync.WaitGroup
		var mu sync.Mutex

		for _, method := range methods {
			sem <- true
			wg.Add(1)

			go func(method string) {
				defer func() {
					<-sem
					wg.Done()
				}()

				result := testMethod(client, url, method, opts)

				shouldAdd := false
				switch opts.viewMode {
				case "all":
					shouldAdd = true
				case "enabled":
					shouldAdd = result.StatusCode != 405 && result.StatusCode != 501
				case "vulnerable":
					shouldAdd = result.IsVulnerable
				}

				if shouldAdd {
					mu.Lock()
					urlResult.Results = append(urlResult.Results, result)
					mu.Unlock()
				}

				if opts.verbose && !opts.quiet && !opts.silent {
					methodColor := color.New(color.FgYellow).SprintFunc()
					statusColor := color.New(color.FgWhite).SprintFunc()

					if result.StatusCode >= 200 && result.StatusCode < 300 {
						statusColor = color.New(color.FgGreen).SprintFunc()
					} else if result.StatusCode >= 400 {
						statusColor = color.New(color.FgRed).SprintFunc()
					} else if result.StatusCode >= 300 {
						statusColor = color.New(color.FgBlue).SprintFunc()
					}

					fmt.Printf("%s %s - %s %s\n",
						methodColor(method),
						url,
						statusColor(fmt.Sprintf("%d", result.StatusCode)),
						statusColor(result.Status))
				}
			}(method)

		}

		wg.Wait()

		results.Results = append(results.Results, urlResult)
	}

	return results
}

func testMethod(client *http.Client, url, method string, opts options) MethodResult {
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return MethodResult{
			Method:          method,
			StatusCode:      0,
			Status:          "Error: " + err.Error(),
			ResponseTime:    0,
			IsDangerous:     dangerousMethods[method],
			IsVulnerable:    false,
			VulnDescription: "",
		}
	}

	for _, header := range opts.headers {
		parts := strings.SplitN(header, ":", 2)
		if len(parts) == 2 {
			req.Header.Add(strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1]))
		}
	}

	req.Header.Set("User-Agent", opts.userAgent)

	for _, cookie := range opts.cookies {
		parts := strings.SplitN(cookie, "=", 2)
		if len(parts) == 2 {
			req.AddCookie(&http.Cookie{
				Name:  parts[0],
				Value: parts[1],
			})
		}
	}

	// Set authentication if provided
	if opts.auth != "" {
		parts := strings.SplitN(opts.auth, ":", 2)
		if len(parts) == 2 {
			req.SetBasicAuth(parts[0], parts[1])
		}
	}

	startTime := time.Now()
	resp, err := client.Do(req)
	endTime := time.Now()
	responseTime := endTime.Sub(startTime).Milliseconds()

	if err != nil {
		return MethodResult{
			Method:          method,
			StatusCode:      0,
			Status:          "Error: " + err.Error(),
			ResponseTime:    responseTime,
			IsDangerous:     dangerousMethods[method],
			IsVulnerable:    false,
			VulnDescription: "",
		}
	}
	defer resp.Body.Close()

	bodyBytes, _ := io.ReadAll(resp.Body)
	bodyStr := string(bodyBytes)

	contentLength := len(bodyBytes)

	isVulnerable := false
	vulnDescription := ""

	if dangerousMethods[method] && (resp.StatusCode >= 200 && resp.StatusCode < 400) {
		isVulnerable = true
		vulnDescription = color.New(color.FgHiRed).Sprintf("Potentially dangerous method %s is allowed", method)
	}

	if method == "TRACE" && resp.StatusCode == 200 && strings.Contains(bodyStr, req.Header.Get("User-Agent")) {
		isVulnerable = true
		vulnDescription = "TRACE method enabled and echoing request headers (potential Cross-Site Tracing vulnerability)"
	}

	// Check for HTTP verb tampering
	if method != "GET" && method != "HEAD" && method != "OPTIONS" &&
		(resp.StatusCode == 200 || resp.StatusCode == 204) {
		isVulnerable = true
		vulnDescription = color.New(color.FgHiRed).Sprintf("Potential HTTP verb tampering vulnerability with method %s", method)
	}

	return MethodResult{
		Method:          method,
		StatusCode:      resp.StatusCode,
		Status:          resp.Status,
		ResponseTime:    responseTime,
		ContentLength:   contentLength,
		IsDangerous:     dangerousMethods[method],
		IsVulnerable:    isVulnerable,
		VulnDescription: vulnDescription,
	}
}

func outputResults(results Results, opts options) {
	var output []byte
	var err error

	switch opts.outputFormat {
	case "json":
		color.NoColor = true
		output, err = json.MarshalIndent(results, "", "  ")
	case "xml":
		color.NoColor = true
		wrapper := XMLResults{Results: results}
		output, err = xml.MarshalIndent(wrapper, "", "  ")
		output = append([]byte(xml.Header), output...)
	case "txt":
		output = []byte(formatTextResults(results, opts))
	default:
		color.New(color.FgRed).Fprintf(os.Stderr, "Unsupported output format: %s\n", opts.outputFormat)
		return
	}

	if err != nil {
		color.New(color.FgRed).Fprintf(os.Stderr, "Error formatting results: %v\n", err)
		return
	}

	if opts.outputFile != "" {
		err = os.WriteFile(opts.outputFile, output, 0644)
		if err != nil {
			color.New(color.FgRed).Fprintf(os.Stderr, "Error writing to output file: %v\n", err)
		} else if !opts.quiet {
			color.New(color.FgGreen).Printf("Results written to %s\n", opts.outputFile)
		}
	}

	// Output to stdout unless quiet mode
	if !opts.quiet {
		fmt.Println(string(output))
	}
}

func formatTextResults(results Results, opts options) string {
	var sb strings.Builder

	sb.WriteString(color.New(color.FgBlue, color.Bold).Sprintf("Author: %s\n", results.ToolAuthor))
	sb.WriteString(color.New(color.FgBlue, color.Bold).Sprintf("%s v%s - Results\n", results.ToolName, results.ToolVersion))
	sb.WriteString(color.New(color.FgYellow, color.Bold).Sprintf("Timestamp: %s\n\n", results.Timestamp))

	for _, urlResult := range results.Results {
		sb.WriteString(color.New(color.FgYellow).Sprintf("URL: %s\n", urlResult.URL))
		sb.WriteString(color.New(color.FgHiCyan, color.Bold).Sprintf("%-10s %-8s %-40s %-12s %s\n", "METHOD", "CODE", "STATUS", "RESPONSE_TIME", " VULNERABILITY"))
		sb.WriteString(strings.Repeat("-", 100) + "\n")

		for _, methodResult := range urlResult.Results {
			vulnStatus := ""
			if methodResult.IsVulnerable {
				vulnStatus = "VULNERABLE"
			} else if methodResult.IsDangerous {
				vulnStatus = "DANGEROUS"
			}

			// sb.WriteString(fmt.Sprintf("%-10s %-8d %-40s %-12d %s\n",
			//	methodResult.Method,
			//	methodResult.StatusCode,
			//	methodResult.Status,
			//	methodResult.ResponseTime,
			//	vulnStatus))

			sb.WriteString(fmt.Sprintf("%-10s  ", methodResult.Method))
			// sb.WriteString(fmt.Sprintf("%-8d  ", methodResult.StatusCode))
			if methodResult.StatusCode >= 200 && methodResult.StatusCode < 300 {
				sb.WriteString(color.New(color.FgRed).Sprintf("%-8d  ", methodResult.StatusCode))
			} else if methodResult.StatusCode >= 400 && methodResult.StatusCode < 500 {
				sb.WriteString(color.New(color.FgGreen).Sprintf("%-8d  ", methodResult.StatusCode))
			} else if methodResult.StatusCode >= 300 && methodResult.StatusCode < 400 {
				sb.WriteString(color.New(color.FgBlue).Sprintf("%-8d  ", methodResult.StatusCode))
			} else if methodResult.StatusCode >= 500 {
				sb.WriteString(color.New(color.FgHiBlue).Sprintf("%-8d  ", methodResult.StatusCode))
			}

			sb.WriteString(fmt.Sprintf("%-40s  ", methodResult.Status))
			sb.WriteString(color.New(color.FgHiYellow).Sprintf("%-12d  ", methodResult.ResponseTime))
			sb.WriteString(color.New(color.FgRed).Sprintf("%s\n", vulnStatus))

			if methodResult.IsVulnerable && methodResult.VulnDescription != "" {
				sb.WriteString(color.New(color.FgHiRed, color.Bold).Sprintf("  - %s\n", methodResult.VulnDescription))
			}
			sb.WriteString(strings.Repeat("-", 100) + "\n")
		}
		sb.WriteString("\n")
	}

	return sb.String()
}
