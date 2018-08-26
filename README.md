# crawler
A simple web crawler in Golang

## Introduction
This web crawler provides a web api retries web pages concurrently using go routine and channel without mutex.
The number of worker goroutine is 5 times number of cpu. This ensures it doesn't exhaust resources e.g. file handle limit or memory.
It also supports limit by depth so it retrieves the pages that are close to the target.

## API endpoint

to start the server:

$ go build -v
$ ./crawler


$ curl http://localhost:8080/v1/crawl -d '{"url": "https://google.com.au", "depth": 2}' > 2.json
  % Total    % Received % Xferd  Average Speed   Time    Time     Time  Current
                                 Dload  Upload   Total   Spent    Left  Speed
100  378k    0  378k  100    44  96973     11  0:00:04  0:00:04 --:--:--  100k

$ ll
-rw-r--r-- 1 Administrator 197121  387892 Feb 13 01:40  2.json

I cannot paste 2.json here as it is big, but it retries 21 pages in 4 seconds.

## possible improvements
1. use a second channel to stop the process when a given timeout reached.
2. HTTPS and JWT to secure the endpoint
3. cache result to make it even faster
