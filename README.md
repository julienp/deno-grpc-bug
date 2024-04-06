# GRPC issues

## deno <-> deno

Run the deno server

```bash
GRPC_TRACE=all GRPC_VERBOSITY=DEBUG deno run -A server.ts
```

and deno client

```bash
GRPC_TRACE=all GRPC_VERBOSITY=DEBUG deno run -A client.ts
```

Running this using deno 1.42.1 fails with

```
error: Uncaught (in promise) TypeError: serde_v8 error: invalid type; expected: string, got: Number
    at ServerHttp2Stream.sendTrailers (node:http2:542:5)
    at ServerHttp2Stream.<anonymous> (file:///Users/julien/Library/Caches/deno/npm/registry.npmjs.org/@grpc/grpc-js/1.10.6/build/src/server-interceptors.js:665:33)
    at Object.onceWrapper (ext:deno_node/_events.mjs:509:28)
    at ServerHttp2Stream.emit (ext:deno_node/_events.mjs:386:28)
    at ServerHttp2Stream.end (node:http2:964:12)
    at BaseServerInterceptingCall.sendStatus (file:///Users/julien/Library/Caches/deno/npm/registry.npmjs.org/@grpc/grpc-js/1.10.6/build/src/server-interceptors.js:668:29)
    at file:///Users/julien/Library/Caches/deno/npm/registry.npmjs.org/@grpc/grpc-js/1.10.6/build/src/server.js:1236:18
    at file:///Users/julien/Library/Caches/deno/npm/registry.npmjs.org/@grpc/grpc-js/1.10.6/build/src/server-interceptors.js:639:13
    at node:http2:475:17
    at Object.runMicrotasks (ext:core/01_core.js:642:26)
```

grpc-js creates the trailers in `server-interceptors.js` like this:

```javascript
const trailersToSend = Object.assign(
  { [GRPC_STATUS_HEADER]: status.code, [GRPC_MESSAGE_HEADER]: encodeURI(status.details) },
  (_a = status.metadata) === null || _a === void 0 ? void 0 : _a.toHttp2Headers()
);
```

where status.code is a number.

By manually patching the file to pass the status code as string `` `${status.code}` `` we get it to run a little better.

The request and response are sent successfully on the transport level, however the response is never delivered to the application.

## deno <-> go

Run the go server

```bash
cd server
GODEBUG=http2debug=2 go run main.go
```

and deno client

```bash
GRPC_TRACE=all GRPC_VERBOSITY=DEBUG deno run -A client.ts
```

The go server is unhappy and returns an error:

```
2024/04/05 18:22:13 http2: Framer 0x1400023a000: wrote SETTINGS len=0
2024/04/05 18:22:13 http2: Framer 0x1400023a000: read SETTINGS len=0
2024/04/05 18:22:13 http2: Framer 0x1400023a000: wrote SETTINGS flags=ACK len=0
2024/04/05 18:22:13 http2: Framer 0x1400023a000: read HEADERS flags=END_HEADERS stream=1 len=127
2024/04/05 18:22:13 http2: decoded hpack field header field ":method" = "POST"
2024/04/05 18:22:13 http2: decoded hpack field header field ":scheme" = "http"
2024/04/05 18:22:13 http2: decoded hpack field header field ":authority" = "localhost:8888"
2024/04/05 18:22:13 http2: decoded hpack field header field ":path" = "/pulumirpc.ResourceMonitor/SupportsFeature"
2024/04/05 18:22:13 http2: decoded hpack field header field "grpc-accept-encoding" = "identity,deflate,gzip"
2024/04/05 18:22:13 http2: decoded hpack field header field "accept-encoding" = "identity"
2024/04/05 18:22:13 http2: decoded hpack field header field "user-agent" = "grpc-node-js/1.10.6"
2024/04/05 18:22:13 http2: decoded hpack field header field "content-type" = "application/grpc"
2024/04/05 18:22:13 http2: decoded hpack field header field "te" = "trailers"
2024/04/05 18:22:13 http2: Framer 0x1400023a000: read DATA stream=1 len=14 data="\x00\x00\x00\x00\t\n\asecrets"
2024/04/05 18:22:13 http2: Framer 0x1400023a000: read HEADERS flags=END_STREAM|END_HEADERS stream=1 len=0
2024/04/05 18:22:13 http2: Framer 0x1400023a000: read SETTINGS flags=ACK len=0
2024/04/05 18:22:13 http2: Framer 0x1400023a000: wrote WINDOW_UPDATE len=4 (conn) incr=14
2024/04/05 18:22:13 http2: Framer 0x1400023a000: wrote PING len=8 ping="\x02\x04\x10\x10\t\x0e\a\a"
2024/04/05 18:22:13 http2: Framer 0x1400023a000: wrote GOAWAY len=137 LastStreamID=1 ErrCode=PROTOCOL_ERROR Debug="received an illegal stream id: 1. headers frame: [FrameHeader HEADERS flags=END_STREAM|END_HEADERS stream=1 len=0]"
2024/04/05 18:22:13 http2: Framer 0x1400023a000: read PING flags=ACK len=8 ping="\x02\x04\x10\x10\t\x0e\a\a"
```

On the deno side we see a broken pipe:

```
D 2024-04-05T16:19:11.877Z | v1.10.6 18901 | transport | (1) 127.0.0.1:8888 connection closed with error connection error received: unspecific protocol error detected (b"received an illegal stream id: 1. maxStreamID=1, headers frame: [FrameHeader HEADERS flags=END_STREAM|END_HEADERS stream=1 len=0]")
error: Uncaught (in promise) Error: stream closed because of a broken pipe
```

## node <-> go

When running the node version of our client, we get the following output from the go server:

````
2024/04/05 21:40:56 http2: Framer 0x140002d2000: wrote SETTINGS len=0
2024/04/05 21:40:56 http2: Framer 0x140002d2000: read SETTINGS len=0
2024/04/05 21:40:56 http2: Framer 0x140002d2000: read HEADERS flags=END_HEADERS stream=1 len=127
2024/04/05 21:40:56 http2: Framer 0x140002d2000: wrote SETTINGS flags=ACK len=0
2024/04/05 21:40:56 http2: decoded hpack field header field ":authority" = "localhost:8888"
2024/04/05 21:40:56 http2: decoded hpack field header field ":method" = "POST"
2024/04/05 21:40:56 http2: decoded hpack field header field ":path" = "/pulumirpc.ResourceMonitor/SupportsFeature"
2024/04/05 21:40:56 http2: decoded hpack field header field ":scheme" = "http"
2024/04/05 21:40:56 http2: decoded hpack field header field "grpc-accept-encoding" = "identity,deflate,gzip"
2024/04/05 21:40:56 http2: decoded hpack field header field "accept-encoding" = "identity"
2024/04/05 21:40:56 http2: decoded hpack field header field "user-agent" = "grpc-node-js/1.10.6"
2024/04/05 21:40:56 http2: decoded hpack field header field "content-type" = "application/grpc"
2024/04/05 21:40:56 http2: decoded hpack field header field "te" = "trailers"
2024/04/05 21:40:56 http2: Framer 0x140002d2000: read DATA stream=1 len=14 data="\x00\x00\x00\x00\t\n\asecrets"
2024/04/05 21:40:56 http2: Framer 0x140002d2000: read SETTINGS flags=ACK len=0
2024/04/05 21:40:56 http2: Framer 0x140002d2000: wrote WINDOW_UPDATE len=4 (conn) incr=14
2024/04/05 21:40:56 http2: Framer 0x140002d2000: read DATA flags=END_STREAM stream=1 len=0 data=""
2024/04/05 21:40:56 http2: Framer 0x140002d2000: wrote PING len=8 ping="\x02\x04\x10\x10\t\x0e\a\a"
2024/04/05 21:40:56 http2: Framer 0x140002d2000: wrote HEADERS flags=END_HEADERS stream=1 len=14
2024/04/05 21:40:56 http2: Framer 0x140002d2000: wrote DATA stream=1 len=7 data="\x00\x00\x00\x00\x02\b\x01"
2024/04/05 21:40:56 http2: Framer 0x140002d2000: wrote HEADERS flags=END_STREAM|END_HEADERS stream=1 len=24
2024/04/05 21:40:56 http2: Framer 0x140002d2000: read PING flags=ACK len=8 ping="\x02\x04\x10\x10\t\x0e\a\a"
```

Of note is that it looks like node is sending an empty data frame with the end_stream flag

```
2024/04/05 21:40:56 http2: Framer 0x140002d2000: read DATA flags=END_STREAM stream=1 len=0 data=""
```

compared to deno instead sending empty trailers, which is the frame that seems to upset the go server.

```
2024/04/05 18:22:13 http2: Framer 0x1400023a000: read HEADERS flags=END_STREAM|END_HEADERS stream=1 len=0
```
````
