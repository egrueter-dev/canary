# Canary - Endpoint Activity Generator

Contents:
1. [About Canary](#about)
2. [Getting Started](#startup)
3. [Commands](#commands)
4. [Canary Platform Support](#platforms)

## About

Canary is an endpoint activity generator. Canary can be used to modify, create, and delete files. The application also comes with a logging capability. Logs can be exported to a remote server via a network call. 

## Startup

Canary comes with pre-built binaries located in `/binaries`.
You can download the binaries for your specific architecture directly from this repository, or you can clone the repository, `cd` into `/binaries` and run them
from there.

## Commands

Below are a list of the available commands:

| Command         | Parameters          | Description                                     |
| --------------- | ------------------- | ------------------------------------------------|
| -list           |                     | List all available commands                     |
| -setup          |                     | Generate log & example files                    |
| -start-process  | filepath, arguments | Execute binary at path                          |
| -create         | filepath            | Create specific file at path                    |
| -delete         | filepath            | Delete specific file at path                    |
| -send-data      | destination         | Execute binary at provided path                 |
| -modify         | filepath text       | Modify file at path with text (text files only) |

These commands are available at the command line by running the `-list` command.

## Platforms 

The Canary application comes with pre built binaries for a number of different architectures.
The `build-executables.sh` shell script is used to build these binaries from the source code. 

 The Following Have been manually tested: 
  - [x]  Darwin/Arm64: 
  - [x] Linux/Amd64
  - []  Windows/Amd64

Binaries for several platforms have been generated. All binaries are located /binaries. These include:
   * Darwin/Arm64
   * linux/amd64
   * linux/arm
   * linux/arm64
   * darwin/amd64
   * darwin/arm64
   * windows/amd64
   * windows/arm
   * windows/arm64

To re-build new binaries after code changes:

 1. Navigate to the root directory
 2. Run the build-executables script and specify the absolute path to the Canary directory Ex: `./build-executables.sh ~/go/src/github.com/canary`
 3. Compiled binaries should appear in `./binaries`
 4. `./build-executables.sh` can be modified to add other architectures. 


 

  
