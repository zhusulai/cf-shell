#!/bin/sh

OUTPUTDIR=$1
FILE=$2
VERSION=$3

if [ "$OUTPUTDIR" = "" -o "$FILE" = "" -o "$VERSION" = "" ]; then
#    echo "Requires output directory, output filename, and version as input."
#    exit 1
    OUTPUTDIR="builds"
    FILE="cf-shell"
    VERSION="current"
fi

mkdir -p $OUTPUTDIR/$VERSION

GOOS=linux GOARCH=amd64 go build -o $OUTPUTDIR/$VERSION/${FILE}.linux64
rc=`echo $?`
if [ $rc -ne 0 ]; then
   echo "linux64 build failed"
   exit 1
fi

GOOS=linux GOARCH=386 go build -o $OUTPUTDIR/$VERSION/${FILE}.linux32
if [ $rc -ne 0 ]; then
   echo "linux32 build failed"
   exit 1
fi

GOOS=darwin GOARCH=amd64 go build -o $OUTPUTDIR/$VERSION/${FILE}.osx
if [ $rc -ne 0 ]; then
   echo "darwin build failed"
   exit 1
fi

GOOS=windows GOARCH=amd64 go build -o $OUTPUTDIR/$VERSION/${FILE}.win64.exe
if [ $rc -ne 0 ]; then
   echo "windows x64 build failed"
   exit 1
fi

GOOS=windows GOARCH=386 go build -o $OUTPUTDIR/$VERSION/${FILE}.win32.exe
if [ $rc -ne 0 ]; then
   echo "windows x32 build failed"
   exit 1
fi