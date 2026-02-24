# GoCrawler

A concurrent, depth-limited web crawler written in Go.

GoCrawler focuses on correctness, bounded concurrency, URL canonicalization, and production-style crawl controls rather than just recursively fetching links.

---

## Overview

GoCrawler starts from a seed URL and crawls pages within the same host. It:

- Normalizes URLs to prevent duplicate crawling
- Uses a bounded worker pool for concurrency
- Applies per-host rate limiting
- Filters non-HTML responses
- Handles redirects safely
- Produces a structured JSON crawl report

This project demonstrates practical backend and networking concepts including HTTP handling, graph traversal, synchronization, and crawl policy design.

---

## Features

### Core Crawling

- Same-host restriction
- Depth-limited traversal
- Max pages limit
- Queue-based crawl frontier

### URL Canonicalization

All URLs are normalized before deduplication and storage:

- Fragment stripping (`#section`)
- Host and scheme normalization
- Default port removal (80/443)
- Trailing slash normalization
- Configurable query string policy

### Concurrency

- Fixed-size worker pool
- Shared crawl frontier queue
- Safe visited-set tracking
- Controlled shutdown when limits are reached

No unbounded goroutines.

### Politeness & Safety

- Custom User-Agent
- Configurable request timeout
- Per-host rate limiting
- Hard limits (depth, max pages)

### HTTP Handling

- Redirect tracking
- Content-Type filtering (HTML-only parsing)
- Response time measurement
- Error categorization (timeout, DNS, parse, etc.)

### Reporting

Structured JSON report including:

- Normalized URL
- Depth
- Status code
- Final URL (after redirects)
- Content-Type
- Response time
- Bytes downloaded
- Number of outlinks
- Error category (if any)

Summary statistics are printed at completion.

---

## Architecture

High-level pipeline:

```
Seed URL
   ↓
Frontier Queue
   ↓
Worker Pool
   ↓
HTTP Fetch
   ↓
Redirect Handling
   ↓
Content-Type Filter
   ↓
HTML Parse
   ↓
Link Extraction
   ↓
URL Normalization
   ↓
Deduplication
   ↓
Enqueue New URLs
   ↓
Structured Report
```

Design goals:

- Deterministic termination
- Bounded memory usage
- No duplicate fetches
- Production-style crawl controls

---

## CLI Usage

Example:

```bash
go run main.go https://example.com \
  -depth 3 \
  -max-pages 500 \
  -concurrency 10 \
  -rate 2 \
  -timeout 5s \
  -output report.json
```

### Flags

| Flag | Description |
|------|------------|
| `-depth` | Maximum crawl depth |
| `-max-pages` | Maximum number of pages to fetch |
| `-concurrency` | Number of worker goroutines |
| `-rate` | Requests per second per host |
| `-timeout` | HTTP request timeout |
| `-output` | Path to JSON report file |

Defaults are conservative to avoid aggressive crawling.

---

## Example Summary Output

```
Crawl completed in 4.2s

Discovered URLs: 148
Fetched: 92
Parsed (HTML): 87
Unique URLs: 92
Errors: 5
Total bytes: 3.1 MB
Average response time: 120ms
```

---

## Example Report Snippet

```json
{
  "url": "https://example.com/about",
  "depth": 1,
  "status": 200,
  "final_url": "https://example.com/about",
  "content_type": "text/html",
  "response_time_ms": 84,
  "bytes": 14231,
  "outlinks": 12,
  "error": null
}
```

---

## Design Tradeoffs

- Same-host restriction simplifies scope and prevents uncontrolled expansion.
- Depth and page limits guarantee bounded execution.
- HTML-only parsing avoids large binary downloads.
- Query string policy is intentionally conservative to reduce infinite URL spaces.
- Concurrency is bounded to prevent resource exhaustion.

---

## Future Improvements

- robots.txt support
- Persistent crawl frontier
- Distributed crawling
- Graph visualization output
- Pluggable storage backend
