#!/bin/bash

## Initial cleanup
if [[ -e build ]]; then
    echo "Clean up the environment"
    rm -fr build 
fi

echo "Create build and tmp dir"
mkdir build tmp/

## Pull busybox
if test ! -f busybox; then
    echo "Download busybox..."
    wget https://busybox.net/downloads/binaries/1.35.0-x86_64-linux-musl/busybox
fi

chmod a+x busybox
mv busybox build
 
## Pull the kernel deb package and extract it
if [[ ! -e $(find . -iname linux-image* -print -quit) ]]; then
    wget "https://kernel.ubuntu.com/mainline/v6.6.8/amd64/linux-image-unsigned-6.6.8-060608-generic_6.6.8-060608.202312201634_amd64.deb"
    mv linux-image-unsigned-6.6.8-060608-generic_6.6.8-060608.202312201634_amd64.deb ./tmp
fi

pushd tmp &> /dev/null
ar x linux-image-unsigned-6.6.8-060608-generic_6.6.8-060608.202312201634_amd64.deb
tar -xf data.tar
mv ./boot/* ../build
popd &> /dev/null

pushd build &> /dev/null

## Run the image
qemu-system-x86_64 -kernel vmlinuz-6.6.8-060608-generic -serial stdio

## Final clean-up
popd &> /dev/null
#rm -fr tmp
