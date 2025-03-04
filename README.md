# go-check-http-methods

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
