#!/usr/bin/env bash
#
# This script builds the application from source for multiple platforms.
set -e

GO_CMD=${GO_CMD:-go}

# Get the parent directory of where this script is.
SOURCE="${BASH_SOURCE[0]}"
while [ -h "$SOURCE" ] ; do SOURCE="$(readlink "$SOURCE")"; done
DIR="$( cd -P "$( dirname "$SOURCE" )/.." && pwd )"

# Change into that directory
cd "$DIR"

# Set build tags
BUILD_TAGS="${BUILD_TAGS:-"fspop"}"
LD_FLAGS="-s -w "

# Get the git commit
GIT_COMMIT="$(git rev-parse HEAD)"
#GIT_DIRTY=""
GIT_DIRTY="$(test -n "`git status --porcelain`" && echo "+CHANGES" || true)"

# Determine the arch/os combos we're building for
XC_ARCH=${XC_ARCH:-"386 amd64"}
XC_OS=${XC_OS:-linux windows}
XC_OSARCH=${XC_OSARCH:-"darwin/amd64 linux/386 linux/amd64 linux/arm linux/arm64 windows/386 windows/amd64"}

# Delete the old dir
echo "==> Removing old directory..."
rm -f bin/*
rm -rf pkg/*
mkdir -p bin/

# Build!
# If GOX_PARALLEL_BUILDS is set, it will be used to add a "-parallel=${GOX_PARALLEL_BUILDS}" gox parameter
echo "==> Building..."
gox \
    -osarch="${XC_OSARCH}" \
    -gcflags "${GCFLAGS}" \
    -ldflags "${LD_FLAGS}-X gitlab.com/merrittcorp/fspop/version.GitCommit='${GIT_COMMIT}${GIT_DIRTY}'" \
    -output "pkg/{{.OS}}_{{.Arch}}/fspop" \
    ${GOX_PARALLEL_BUILDS+-parallel="${GOX_PARALLEL_BUILDS}"} \
    -tags="${BUILD_TAGS}" \
    -gocmd="${GO_CMD}" \
    .

# Move all the compiled things to the $GOPATH/bin
OLDIFS=$IFS
IFS=: MAIN_GOPATH=($GOPATH)
IFS=$OLDIFS

# Copy our OS/Arch to the bin/ directory
#DEV_PLATFORM=${DEV_PLATFORM:-"./pkg/$(${GO_CMD} env GOOS)_$(${GO_CMD} env GOARCH)"}
#for F in $(find ${DEV_PLATFORM} -mindepth 1 -maxdepth 1 -type f); do
#    cp ${F} bin/
#    cp ${F} ${MAIN_GOPATH}/bin/
#done


# Zip and copy to the dist dir
echo "==> Packaging..."
for PLATFORM in $(find ./pkg -mindepth 1 -maxdepth 1 -type d); do
    OSARCH=$(basename ${PLATFORM})
    echo "--> ${OSARCH}"

    pushd $PLATFORM >/dev/null 2>&1
    zip ${OSARCH}.zip *
    mv ${OSARCH}.zip ../../bin/${OSARCH}.zip
    popd >/dev/null 2>&1
done

# Done!
echo
echo "==> Results bin/:"
ls -R bin/
