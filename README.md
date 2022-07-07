# packetmaker
A low allocation, lightweight way to generate packets in the Go programming language. This exists because doing standard byte array slicing is complex and very error-prone, and I wanted a simple way to write packets efficiently (the other way I was using was append, but this sucks for other memory management reasons).

## API usage
The API is documented on GoDoc [here](https://pkg.go.dev/github.com/jakemakesstuff/packetmaker) and should contain everything you need to make packets (if there is anything missing you need, feel free to make an issue/PR!). These functions are designed to be chained, for example `packetmaker.New().String("abc").Uint64(12, true).Make()`.

## Note on Big Endian systems
This library is heavily unit tested to ensure its reliability on little endian systems (most systems you would use come under this, like x86_64/arm64), although the Int functions (not the unsigned ones) are unknown if they work on big endian since they use unsafe. Feel free to make an issue if you get weird results with this (or if it works, I'd love to know!), you will not ever get unsafe results back, but there is a slight chance they might be incorrect.
