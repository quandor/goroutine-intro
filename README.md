# Analyzes
* Check for README.md, readme.md, README.adoc, readme.adoc
* Report reponame, readme file format, readme size

# Talk
## goroutines What are they used for?
* concurrency made easy
* Lightweight threading
* Ever made yourself something with plain Threads? Don't :-)
## What does that look like?
* simple programm calling a goroutine
* May a goroutine return something?
## How can you communicate with the rest?
* Channels
* Select
* Timers, Tickers
## Sublte Problems
http://www.sarathlakshman.com/2016/06/15/pitfall-of-golang-scheduler
## Differences to coroutines
* Structered concurrency??
* there is no launch and async in go. there is just a launch
## Summary
* Looks easy, but it ain't
* Still a lot easier than standard Java mechanisms like Executors or plain
Threads
## Link list
Intro: https://riteeksrivastava.medium.com/a-complete-journey-with-goroutines-8472630c7f5c
technical details: https://dave.cheney.net/2015/08/08/performance-without-the-event-loop
Nice examples : https://gobyexample.com
