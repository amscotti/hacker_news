# hacker_news
A command line tool that shows the top stories from [Hacker News](https://news.ycombinator.com/) using the [Hacker News API from Firebase](https://github.com/HackerNews/API).

This is a port of an dotnet version, [hn-top](https://github.com/amscotti/hn-top) which I used to learn about Go and understand how to keep the order of the stories when using Goroutines with help of Mutex and WaitGroup.

![hacker_news Output](https://github.com/amscotti/hacker_news/blob/main/hacker_news.png?raw=true)

## Building and Running

### With Go
* Build with `go build`
* Then run with `./hacker_news`

### With Docker
* Build with `docker build -t hacker_news . `
* Then run with `docker run hacker_news`

### With pre-build Docker image from [ghcr.io](https://github.com/amscotti/hacker_news/pkgs/container/hacker_news)
* To download and run, use `docker run ghcr.io/amscotti/hacker_news:main`

### Command Line Arguments
```
Usage of ./hacker_news:
  -n int
        Number of top news to show. (default 5)
  -u    Show source urls.
```