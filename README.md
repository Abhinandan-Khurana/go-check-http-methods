<p align="center">
  <img src="./go-check-http-methods.svg" width="300" height="300">
</p>

<p align="center">
  <a href="https://goreportcard.com/report/github.com/Abhinandan-Khurana/go-check-http-methods"><img src="https://goreportcard.com/badge/github.com/Abhinandan-Khurana/go-check-http-methods" alt="Go Report Card"></a>
  <a href="LICENSE"><img src="https://img.shields.io/badge/license-MIT-blue.svg" alt="License"></a>
  <a href="https://golang.org/doc/devel/release.html"><img src="https://img.shields.io/badge/Go-1.23+-00ADD8.svg" alt="Go Version"></a>
<img src="https://img.shields.io/badge/version-v1.0.0-blue.svg" alt="Version">
</p>

A powerful, fast, and lightweight Go tool for testing HTTP method security across web applications and servers. Detect HTTP Verb Tampering vulnerabilities, dangerous methods, and server misconfigurations with ease.

## 🚀 Features

- 🔥 **Comprehensive Method Testing** - Tests a wide range of HTTP methods including standard and non-standard methods
- 🔍 **Vulnerability Detection** - Identifies potential HTTP verb tampering vulnerabilities and dangerous method configurations
- 🌍 **Multiple URL Support** - Test individual URLs or bulk test from a file
- 📊 **Flexible Output Formats** - Results in plain text, JSON, or XML format
- 🔒 **Custom Authentication** - Support for basic auth, custom headers, and cookies
- 🚦 **Concurrent Testing** - Configurable concurrency for faster scanning
- 🎨 **Colorized Output** - Easy-to-read color-coded results (with no-color option)
- 🔄 **Proxy Support** - Route requests through a proxy for anonymity or testing internal networks
- 📱 **Cross-Platform** - Works on Windows, macOS, and Linux (both amd64 and arm64)

## 📥 Installation

### Pre-compiled Binaries

Download the latest release for your platform from the [releases page](https://github.com/Abhinandan-Khurana/go-check-http-methods/releases).

### From Source

Make sure you have Go 1.23+ installed, then run:

```bash
# Clone the repository
git clone https://github.com/Abhinandan-Khurana/go-check-http-methods.git
cd go-check-http-methods

# Build for your platform
make install

# Or build for all platforms
make all
```

### Using Go Install

```bash
go install github.com/Abhinandan-Khurana/go-check-http-methods@latest
```

## 📚 Usage

### Basic Usage

```bash
# Test a single URL

go-check-http-methods -u <https://example.com>

# Test URLs from a file

go-check-http-methods -f urls.txt

# Output results in JSON format

go-check-http-methods -u <https://example.com> -format json

# Save results to a file

go-check-http-methods -u <https://example.com> -o results.txt
```

### Advanced Usage

```bash
# Test with custom HTTP methods

go-check-http-methods -u <https://example.com> -m methods.txt

# Increase concurrency for faster scanning

go-check-http-methods -u <https://example.com> -c 20

# Use silent mode (no banner, only results)

go-check-http-methods -u <https://example.com> -silent

# Disable color output

go-check-http-methods -u <https://example.com> -nc

# Add custom headers and authentication

go-check-http-methods -u <https://example.com> -H "Authorization: Bearer token" -H "X-Custom: Value"

# Use a proxy

go-check-http-methods -u <https://example.com> -proxy <http://127.0.0.1:8080>
```

### Command-Line Options

```bash
-u string
Single URL to test
-f string
File containing URLs to test (one per line)
-o string
Output file for results
-format string
Output format: txt, json, or xml (default "txt")
-m string
File containing HTTP methods to test
-c int
Number of concurrent requests (default 10)
-t int
Request timeout in seconds (default 10)
-L Follow redirects
-k Allow insecure TLS connections
-v Verbose output
-q Quiet mode, no output except results
-silent
Silent mode, no banner but show results
-nc No color output
-view string
View mode: all, enabled, vulnerable (default "all")
-proxy string
Use proxy (format: <http://host:port>)
-ua string
User agent string
-H string
Custom header (can be used multiple times, format: 'Name: Value')
-cookie string
Cookie to include (can be used multiple times, format: 'name=value')
-auth string
Basic authentication (format: username:password)
```

## 🔍 Examples

### Testing for Verb Tampering Vulnerabilities

```bash
go-check-http-methods -u <https://example.com> -view vulnerable -v
```

This command will:

- Test the URL with all supported HTTP methods
- Show only vulnerable methods in the results
- Display verbose output during scanning

### Sample Output

```bash
                         _               _         _     _   _                              _   _               _
   __ _  ___         ___| |__   ___  ___| | __    | |__ | |_| |_ _ __        _ __ ___   ___| |_| |__   ___   __| |___
  / _' |/ _ \ _____ / __| '_ \ / _ \/ __| |/ /____| '_ \| __| __| '_ \ _____| '_ ' _ \ / _ \ __| '_ \ / _ \ / _' / __|
 | (_| | (_) |_____| (__| | | |  __/ (__|   <_____| | | | |_| |_| |_) |_____| | | | | |  __/ |_| | | | (_) | (_| \__ \
  \__, |\___/       \___|_| |_|\___|\___|_|\_\    |_| |_|\__|\__| .__/      |_| |_| |_|\___|\__|_| |_|\___/ \__,_|___/
  |___/                                                         |_|

Author: Abhinandan Khurana aka @l0u51f3r007
go-check-http-methods v1.0.0 - Results
Timestamp: 2025-03-04T16:41:01+05:30

URL: https://example.com
METHOD     CODE     STATUS                                   RESPONSE_TIME  VULNERABILITY
----------------------------------------------------------------------------------------------------
POST        403       403 Forbidden                             928
----------------------------------------------------------------------------------------------------
PROPFIND    501       501 Not Implemented                       928           DANGEROUS
----------------------------------------------------------------------------------------------------
HEAD        200       200 OK                                    929
----------------------------------------------------------------------------------------------------
PATCH       501       501 Not Implemented                       929
----------------------------------------------------------------------------------------------------
DELETE      501       501 Not Implemented                       939           DANGEROUS
----------------------------------------------------------------------------------------------------
GET         200       200 OK                                    939
----------------------------------------------------------------------------------------------------
CONNECT     400       400 Bad Request                           946
----------------------------------------------------------------------------------------------------
TRACE       403       403 Forbidden                             960           DANGEROUS
----------------------------------------------------------------------------------------------------
PUT         501       501 Not Implemented                       963           DANGEROUS
----------------------------------------------------------------------------------------------------
OPTIONS     501       501 Not Implemented                       1027
----------------------------------------------------------------------------------------------------
MKCOL       501       501 Not Implemented                       227           DANGEROUS
----------------------------------------------------------------------------------------------------
COPY        501       501 Not Implemented                       245           DANGEROUS
----------------------------------------------------------------------------------------------------
MOVE        501       501 Not Implemented                       283           DANGEROUS
----------------------------------------------------------------------------------------------------
PROPPATCH   501       501 Not Implemented                       348           DANGEROUS
----------------------------------------------------------------------------------------------------
LOCK        501       501 Not Implemented                       810           DANGEROUS
----------------------------------------------------------------------------------------------------
UNLINK      400       400 Bad Request                           785
----------------------------------------------------------------------------------------------------
LINK        400       400 Bad Request                           789
----------------------------------------------------------------------------------------------------
PURGE       400       400 Bad Request                           803
----------------------------------------------------------------------------------------------------
UNLOCK      501       501 Not Implemented                       809           DANGEROUS
----------------------------------------------------------------------------------------------------
```

### Sample JSON output

```json
{
  "tool_name": "go-check-http-methods",
  "tool_version": "1.0.0",
  "tool_author": "Abhinandan Khurana aka @l0u51f3r007",
  "timestamp": "2025-03-04T16:41:57+05:30",
  "results": [
    {
      "url": "https://example.com",
      "results": [
        {
          "method": "CONNECT",
          "status_code": 400,
          "status": "400 Bad Request",
          "response_time_ms": 810,
          "content_length": 312,
          "is_dangerous": false,
          "is_vulnerable": false
        },
        {
          "method": "TRACE",
          "status_code": 403,
          "status": "403 Forbidden",
          "response_time_ms": 810,
          "content_length": 359,
          "is_dangerous": true,
          "is_vulnerable": false
        },
        {
          "method": "PATCH",
          "status_code": 501,
          "status": "501 Not Implemented",
          "response_time_ms": 811,
          "content_length": 336,
          "is_dangerous": false,
          "is_vulnerable": false
        },
        {
          "method": "POST",
          "status_code": 403,
          "status": "403 Forbidden",
          "response_time_ms": 813,
          "content_length": 359,
          "is_dangerous": false,
          "is_vulnerable": false
        },
        {
          "method": "OPTIONS",
          "status_code": 501,
          "status": "501 Not Implemented",
          "response_time_ms": 813,
          "content_length": 19,
          "is_dangerous": false,
          "is_vulnerable": false
        },
        {
          "method": "DELETE",
          "status_code": 501,
          "status": "501 Not Implemented",
          "response_time_ms": 810,
          "content_length": 339,
          "is_dangerous": true,
          "is_vulnerable": false
        },
        {
          "method": "PUT",
          "status_code": 501,
          "status": "501 Not Implemented",
          "response_time_ms": 810,
          "content_length": 334,
          "is_dangerous": true,
          "is_vulnerable": false
        },
        {
          "method": "HEAD",
          "status_code": 200,
          "status": "200 OK",
          "response_time_ms": 810,
          "content_length": 0,
          "is_dangerous": false,
          "is_vulnerable": false
        },
        {
          "method": "GET",
          "status_code": 200,
          "status": "200 OK",
          "response_time_ms": 810,
          "content_length": 1256,
          "is_dangerous": false,
          "is_vulnerable": false
        },
        {
          "method": "PROPFIND",
          "status_code": 501,
          "status": "501 Not Implemented",
          "response_time_ms": 810,
          "content_length": 339,
          "is_dangerous": true,
          "is_vulnerable": false
        },
        {
          "method": "LOCK",
          "status_code": 501,
          "status": "501 Not Implemented",
          "response_time_ms": 299,
          "content_length": 337,
          "is_dangerous": true,
          "is_vulnerable": false
        },
        {
          "method": "MKCOL",
          "status_code": 501,
          "status": "501 Not Implemented",
          "response_time_ms": 300,
          "content_length": 336,
          "is_dangerous": true,
          "is_vulnerable": false
        },
        {
          "method": "LINK",
          "status_code": 400,
          "status": "400 Bad Request",
          "response_time_ms": 762,
          "content_length": 312,
          "is_dangerous": false,
          "is_vulnerable": false
        },
        {
          "method": "UNLOCK",
          "status_code": 501,
          "status": "501 Not Implemented",
          "response_time_ms": 764,
          "content_length": 339,
          "is_dangerous": true,
          "is_vulnerable": false
        },
        {
          "method": "MOVE",
          "status_code": 501,
          "status": "501 Not Implemented",
          "response_time_ms": 790,
          "content_length": 335,
          "is_dangerous": true,
          "is_vulnerable": false
        },
        {
          "method": "PURGE",
          "status_code": 400,
          "status": "400 Bad Request",
          "response_time_ms": 789,
          "content_length": 312,
          "is_dangerous": false,
          "is_vulnerable": false
        },
        {
          "method": "PROPPATCH",
          "status_code": 501,
          "status": "501 Not Implemented",
          "response_time_ms": 791,
          "content_length": 342,
          "is_dangerous": true,
          "is_vulnerable": false
        },
        {
          "method": "COPY",
          "status_code": 501,
          "status": "501 Not Implemented",
          "response_time_ms": 793,
          "content_length": 335,
          "is_dangerous": true,
          "is_vulnerable": false
        },
        {
          "method": "UNLINK",
          "status_code": 400,
          "status": "400 Bad Request",
          "response_time_ms": 799,
          "content_length": 312,
          "is_dangerous": false,
          "is_vulnerable": false
        }
      ]
    }
  ]
}
```

## 🏗️ Building from Source

You can build this tool for multiple platforms using the included Makefile:

```bash
# Build for all platforms

make all

# Build for specific platform

make linux-amd64

# Clean build artifacts

make clean

# Package builds into zip files

make package
```

## 📋 HTTP Method Descriptions

| Method    | Description                                      | Potential Risk |
| --------- | ------------------------------------------------ | -------------- |
| GET       | Retrieve a resource                              | Low            |
| POST      | Create a new resource                            | Medium         |
| PUT       | Update a resource                                | High           |
| DELETE    | Delete a resource                                | High           |
| HEAD      | Similar to GET but returns only headers          | Low            |
| OPTIONS   | Returns the HTTP methods supported by the server | Low            |
| PATCH     | Partial update of a resource                     | Medium         |
| TRACE     | Echo the received request                        | High           |
| CONNECT   | Establish a network connection                   | Medium         |
| PROPFIND  | WebDAV method to retrieve properties             | High           |
| PROPPATCH | WebDAV method to change properties               | High           |
| MKCOL     | WebDAV method to create collections              | High           |
| COPY      | WebDAV method to copy a resource                 | High           |
| MOVE      | WebDAV method to move a resource                 | High           |
| LOCK      | WebDAV method to lock a resource                 | High           |
| UNLOCK    | WebDAV method to unlock a resource               | High           |

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📝 License

This project is licensed under the MIT License - see the LICENSE file for details.

## 🙏 Acknowledgments

- [Fatih Color](https://github.com/fatih/color) - For the colorized output

---

Made with <3 by <a href="https://x.com/l0u51f3r007">Abhinandan Khurana</a>
