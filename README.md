# GRPC issues

Run the server

```bash
GRPC_TRACE=all GRPC_VERBOSITY=DEBUG deno run -A server.ts
```

and client

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
