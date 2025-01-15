# Wget Clone in Go

## Project Overview

This project aims to replicate the core functionalities of `wget`, a widely used utility for downloading files from the web. The implementation is done in Go, focusing on efficiency, usability, and adherence to the core principles of `wget`.

### Objectives
The program includes the following functionalities:

1. Downloading a single file from a given URL.
2. Saving the file under a different name.
3. Saving the file to a specific directory.
4. Limiting the download speed using a rate limit.
5. Downloading a file in the background with logs redirected to a file.
6. Downloading multiple files asynchronously by reading links from a file.
7. Mirroring an entire website for offline usage.
8. Read a file and extract download URLs.
9. Convert links for offline viewing.

---

## Features and Usage

### Installation

To use this application, make sure you have Go installed on your system. If you don’t have Go installed, you can download it from the [official Go website](https://go.dev/dl/). Follow the instructions provided there to complete the installation.

Once Go is installed, follow these steps:

1. **Clone the Repository**:
   Clone the project repository to your local machine:
   ```bash
   git clone https://learn.zone01kisumu.ke/git/aaochieng/wget
   cd wget
   ```

2. **Build the Program**:  
   Use the `go build` command to build the application into an executable:  
   ```bash
   go build -o wget .
   ```

3. **Run the Program**:  
   Execute the program using the generated executable. For example:  
   ```bash
   ./wget <flags> <URL>
   ```

4. **Verify Installation**:  
   To confirm that the application is working, try downloading a sample file:  
   ```bash
   ./wget https://pbs.twimg.com/media/EMtmPFLWkAA8CIS.jpg
   ```

---

Let me know if you'd like to expand on this further!

### Basic Usage

To download a file, pass its URL as an argument to the program:

```bash
$ go run . https://pbs.twimg.com/media/EMtmPFLWkAA8CIS.jpg
```

The program will display details such as start time, HTTP status, content size, download progress, and end time.

### Flags and Options

#### Background Download (`-B`)
Downloads a file in the background, logging the output to `wget-log`:

```bash
$ go run . -B https://pbs.twimg.com/media/EMtmPFLWkAA8CIS.jpg
Output will be written to "wget-log".
```

#### Save with a Different Name (`-O`)
Saves the file with a specified name:

```bash
$ go run . -O=newfile.zip <url>
```

#### Save to a Specific Directory (`-P`)
Saves the file to the given directory:

```bash
$ go run . -P=~/Downloads/ <url>
```

#### Rate Limiting (`--rate-limit`)
Limits the download speed:

```bash
$ go run . --rate-limit=500k <url>
```

#### Asynchronous Download (`-i`)
Downloads multiple files asynchronously by reading a file containing URLs:

```bash
$ go run . -i=links.txt
```
The `links.txt` file should contain one URL per line.

#### Website Mirroring (`--mirror`)
Mirrors an entire website:

```bash
$ go run . --mirror https://example.com
```

#### Optional Flags for Mirroring

- **Exclude File Types (`-R`)**: Avoid downloading specified file types:

  ```bash
  $ go run . --mirror -R=jpg,png https://example.com
  ```

- **Exclude Directories (`-X`)**: Avoid specific directories:

  ```bash
  $ go run . --mirror -X=/assets,/images https://example.com
  ```

- **Convert Links (`--convert-links`)**: Converts links for offline viewing.

  ```bash
  $ go run . --mirror --convert-links https://example.com
  ```

**Note:** Prefer to download websites with `--convert-links` for better offline viewing.
---

## Implementation Details

### Program Output
The program provides feedback for each download operation:

- **Start Time**: Displayed in `YYYY-MM-DD HH:MM:SS` format.
- **HTTP Status**: Indicates the response status (e.g., `200 OK` or error messages).
- **Content Size**: Shows the file size in bytes, KiB, MiB, or GiB.
- **File Path**: Displays where the file will be saved.
- **Progress Bar**: Updates in real-time with downloaded size, percentage, and estimated time remaining.
- **End Time**: Displayed in `YYYY-MM-DD HH:MM:SS` format.

### Example Output (disclaimer, your output maybe different!)

```bash
$ go run . https://example.com/file.zip
start at 2025-01-15 12:34:56
sending request, awaiting response... status 200 OK
content size: 20485760 [~20.49MB]
saving file to: ./file.zip
 20.49 MiB / 20.49 MiB [================================================] 100.00% 2.00 MiB/s 0s

Downloaded [https://example.com/file.zip]
finished at 2025-01-15 12:34:59
```

---

## Development

### Technologies Used
- **Language**: Go
- **Concurrency**: For asynchronous downloads
- **Standard Library**: Extensive use of Go’s standard library to avoid third-party dependencies

### Key Concepts
- **Rate Limiting**: Ensures bandwidth control.
- **Progress Display**: Real-time updates for user feedback.
- **Error Handling**: Graceful handling of HTTP errors and invalid inputs.

---

## Future Enhancements
- Add support for FTP protocols.
- Implement resume functionality for interrupted downloads.
- Support custom headers and authentication mechanisms.
