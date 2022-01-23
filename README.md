# Canary - Endpoint Activity Generator


 Contents
1. [Canary Platform Support](#platforms)
2. [Example2](#example2)
3. [Third Example](#third-example)
4. [Fourth Example](#fourth-examplehttpwwwfourthexamplecom)

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


 

  
