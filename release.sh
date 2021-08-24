#!/bin/bash


TARGETS="windows,zip darwin,tar.gz linux,tar.gz"

read -p "Version: (x.y.z format): " VERSION
echo 

for target in $TARGETS
do
    GOOS=$(echo $target | cut -d , -f 1)
    EXTENSION=$(echo $target | cut -d , -f 2)

    GOOS=$GOOS go build -ldflags="-X main.Version=v$VERSION" -o build/$GOOS/ ./cmd/mpdrp
    # Additional assets can be added below
    case $GOOS in
        windows)
            cp config/mpdrp.bat build/windows/
            ;;
        linux)
            cp config/mpdrp.service build/linux/
            ;;
        darwin)
            cp config/mpdrp.plist build/darwin/
            ;;
    esac
    # Finally, pack it up
    ARCHIVE_FILENAME="mpdrp-$GOOS-$(go env GOARCH).$EXTENSION"
    echo "Creating $GOOS release: $ARCHIVE_FILENAME"
    tar -c --strip-components=2 -a -f build/$ARCHIVE_FILENAME build/$GOOS/*
    
done




