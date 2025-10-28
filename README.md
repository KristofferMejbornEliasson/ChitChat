# ChitChat

Assignment for BSDISYS1KU 2025/2026 at ITU Copenhagen.

**ChitChat** is a small real-time messaging application, where multiple clients connect to the same
server, and can then exchange plaintext messages (up to 128 characters long) with each other.

## Requirements

The following need to be installed and on the user's PATH:

* [Go](https://go.dev/)
* the [Protobuffer compiler](https://protobuf.dev/installation/)
* the [protobuffer compiler plugin for Go](https://github.com/protocolbuffers/protobuf-go)
* the [GRPC library](https://grpc.io/docs/languages/go/quickstart/)

## Run without compiling

It is not necessary to compile the programmes to native executables.

The server should (and indeed, _can_) only be started once. To do so, execute the following command
in the root project directory:

```
go run server/server.go
```

For each client instance you wish to create, open a new terminal and execute:

```
go run client/main.go
```

## Compile and Run

### Compile natively

Windows-users with `Make` installed can execute `make` to compile both the server and client
programmes
to Windows executables `server.exe` and `client.exe` in the root project directory.

To instead compile them manually (regardless of operating system), execute the following commands.
If a custom output name is _not_ provided with `-o <output-name>`, then the
results will be named `server` and `main` respectively, once again in the main project directory.

```
go build [-o <ouput-name>] server/server.go
```

```
go build [-o <output-name>] client/main.go
```

### Run

Once compiled, running them is as simple as executing the programmes. For example, if named `server`
and `client`:

```
./server
```

```
./client
```

## Usage

Up to 2<sup>32</sup>-1 client instances can be created during a single session (...theoretically).\
The client should automatically connect to the server if the latter is listening.

New messages can be written in the client terminal. Pressing `Return` will send what has been
written.\
Note that messages have a maximum length of 128 characters.\
Messages are relayed to all other active clients.

To gracefully close the server, type `exit` into it and press `Return`. This is the only
command that the server responds to; anything else is entirely ignored.

To gracefully quit a client, type `exit` and press `Return`.

Killing a process with `Ctrl+C` works as expected in both cases, too, but logs won't include the
shutdown event.

The server's logging information is written to the `log.txt` file in the project root.\
Each client's log information is written to a uniquely-named file in the user's temporary directory
with the name `client-*.txt`, where `*` is a short, randomly-generated string.