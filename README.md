# Goroutine intro
This repository exists to 
* keep the presentation about a goroutine intro
* be a home to the sample code

The sample analyses Readmes of Novatec's Github projects.
There are 3 approaches:
* without goroutines (execute `sample.sh`)
* goroutines synchronized by a WaitGroup (execute `sample.sh --waitgroup`)
* goroutines synchronized via Channels (execute `sample.sh --channel`)

For the ease of runtime comparision I suggest to use `time sample.sh
[--waitgroup|--channel]`

If you want to analyse additional projects simply edit `sample.sh` and append 
the project's name to the parameter list.

Please note, that just because we are using goroutines we are not necessarily
faster.
