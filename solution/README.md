
## Minerva
Gathers all the wisdom about your machine's condition (all that you need at least)

### Ecosystem
Although Golang is neither a language that I feel proficient in, nor a language that I used professionally, I considered it a fine choice for this particular task as being robust, well documented and having access to some of the low-level methods. The project is split into four source files, a **Makefile**, a **yaml** configuration file and HTML page whose purpose I'll explain later.

The project is compatible with go >=1.11

### Getting started
Just enter to the project's directory and run `make`. The Makefile will automatically get the dependencies and build the binary for either Linux or OS X (you won't have to fetch Windows machine :)). Make sure you have the Go environment [configured properly](https://golang.org/doc/install).

Afterwards run the minerva **executable** from within the project root directory. It is important that the yaml and html files are inside your [current working directory](https://linux.die.net/man/3/cwd).

Minerva accepts two optional flags `-v` to enable debug messages and `-vv ` to enable verbose logging. `-vv` encapsulates `-v`.

### Configuration

Minerva is configurable from within attached **yaml** file. Description of each key is to be found within the configuration file itself.

### How does it work
Minerva is polling specified TCP server with given intervals and considers given timeouts. It performs either Layer 7 (HTTP) or Layer 4 (TCP) checks where L7 check naturally does L4 check as well. This option is configurable in the yaml file.

Optionally minerva is able to validate HTTP response's status code as well as response body. This allows user to be able to determine if HTTP server responds with the correct web page or whether it is responsive at all.

Http checks have configurable maximum number of bytes that minerva is going to pull from the server. This is to prevent stalling minerva should the response be very long.

Minerva reports four different server health states:
* Up
* Down
* Down, transitionally up
* Up, transitionally down

While the first two are quite obvious, the latter vary per scenario. Down, transitionally up means that the server is Down but the last health check succeeded. Respectively Up, transitionally down means that the server is Up but the last health check has failed.

How long the server is considered  to be "in transition" is configurable from within the yaml file in `Tries` section. It is important to **note** that the server will be back in the current state if the consecutive check are disrupted (ie. if the configuration specifies that the server is considered Alive after 5 consecutive checks, then if you eg. have 3 successful health checks and the 4th one fails then server will be back in Dead state and you must once again collect 5 consecutive successful checks in order for minerva to consider it Alive).

### Reporting/Monitoring
Although the assignment description **suggested** that log files could be used to expose magnificent's health, I considered flat files a very oldschool way of piping data without any service operating on a socket level, thus I decided to go with WebSocket protocol

If Monitoring is enabled in yaml file, minerva is going to expose a WebSocket endpoint `/ws` bound to ip:port specified in the configuration file, eg `ws://localhost:8080/ws`. The WS port is for reading only. All writes to it are discarded.

 Additionally there's a mini-client-webpage (the previously mentioned HTML file) available at `http://localhost:8080/`which connects to the websocket and prints the results. Please don't judge my UI skills -- I really like them as-is :)

Minerva does not report reason why magnificent was considered dead apart from that it was. It reports back  `state` which is one of the four states described before (this should conform the requirement *"The service has to indicate how healthy the Magnificent service has been over the last little while"*) and latest check status which can be either `OK` or `ERR`.  Additionally it passes metadata such as latest check HTTP status code (if `Tcponly` is set to `false`), duration of the last health check in ms, the type of check (`L4` or `L7`) and a timestamp.

If check interval is low this can generate quite a large amount of information. You might consider enabling `TransitionOnly` flag to reduce the amount of data pushed through the socket.

### Demo
I modified magnificent's code a little bit so that it looks like so
```
def render_GET(self, request):
        if choice([True, True, True, False]):
            if choice([True, True, False]):
                print('Magnificent')
                return "Magnificent!".encode('utf-8')  # Twisted expects bytes.
            else:
                print('Something else')
                return "Something else".encode('utf-8')  # Twisted expects bytes.
        else:
            raise Uninspiring()
```

Also I have monitoring enabled with `Transitiononly` set to `false`. The check interval is `700ms`, `Tcponly` set to `false` and I need 2 consecutive checks to get up as well as 3 consecutive to go down.

You can see it in action here:
https://www.youtube.com/watch?v=2kjrFhK3w9A

### Known limitations
* Lack of support for HTTP/2.0
* Lacks proper SSL support
* IPv6 was not tested
* HTTP verification is limited to only one status code per application run
* Mini monitor webpage requires Java -version "script" (sorry, *this is the world we live in*)

### Further work
The project does not ultimately exhaust the assignment. Future improvements include, but are not limited to:
* Support multiple ways of exposing magnificent health status, preferably using a common interface
* Improve efficiency in Layer 4 checks by terminating TCP handshake (RST after SYN-ACK) which requires custom TCP stack implementation that is not trivial
* The websocket being duplex protocol could be used for changing minerva configuration at runtime
* Verbose configuration parsing
* Rewrite it [Dyalog APL](https://github.com/jayfoad/aoc2018apl) and [sell it to the client without support](https://www.youtube.com/watch?v=BdvUR67nZs0)
