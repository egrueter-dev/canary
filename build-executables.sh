#!/usr/bin/env bash

# To use this file:
# chmod +x go-executable-build.bash to build executableA
# ./go-executable-build.bash

package=$1

if [[ -z "$package" ]]; then
  echo "usage: $0 canary"
  exit 1
fi

package_name="canary"

platforms=(
    "linux/amd64"
    "linux/arm"
    "linux/arm64"
    "darwin/amd64"
    "darwin/arm64"
    "windows/amd64" 
    "windows/arm"
    "windows/arm64"
)

for platform in "${platforms[@]}"
do
    platform_split=(${platform//\// })

    GOOS=${platform_split[0]}

    GOARCH=${platform_split[1]}

    output_name=$package_name'-'$GOOS'-'$GOARCH

    if [ $GOOS = "windows" ]; then
        output_name+='.exe'
    fi

    env GOOS=$GOOS GOARCH=$GOARCH go build -o $output_name $package

    mv $output_name ./binaries

    if [ $? -ne 0 ]; then
        echo 'An error has occurred! Aborting the script execution...'
        exit 1
    fi
done