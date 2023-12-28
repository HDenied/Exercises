#!/bin/bash

if [[ -e build ]]; then
    echo "Clean up the environment"
    rm -fr build 
fi

export ROOT_DIR=$(pwd)
export TMP_DIR=$ROOT_DIR/tmp
export BUILD_DIR=$ROOT_DIR/build
export INITRAMFS_DIR=$TMP_DIR/initramfs

mkdir $BUILD_DIR $TMP_DIR

echo "Extracting the kernel"
## Pull the kernel deb package and extract it
if [[ ! -e $(find . -iname linux-image* -print -quit) ]]; then
    wget "https://kernel.ubuntu.com/mainline/v6.6.8/amd64/linux-image-unsigned-6.6.8-060608-generic_6.6.8-060608.202312201634_amd64.deb"
    mv linux-image-unsigned-6.6.8-060608-generic_6.6.8-060608.202312201634_amd64.deb $TMP_DIR
fi

pushd $TMP_DIR &> /dev/null
ar x linux-image-unsigned-6.6.8-060608-generic_6.6.8-060608.202312201634_amd64.deb
tar -xf data.tar
mv ./boot/* $BUILD_DIR
popd &> /dev/null

## Create a working initramfs image using prcompiled busybox ##
# For reference I compiled busybox version 1.35.0 x86_64 using the following commands
#  make defconfig
#  make menuconfig
# and in make_menuconfig I selected static binary
# after that I ran: 
#  make
#  make CONFIG_PREFIX=../busybox_build

echo -e "Generating initramfs\n"
mkdir -p $INITRAMFS_DIR/{bin,sbin,etc,proc,sys,usr/bin,usr/sbin}
cp -a busybox_build/* $INITRAMFS_DIR
cp -a init $INITRAMFS_DIR
pushd $INITRAMFS_DIR
find . -print0 | cpio --null -ov --format=newc  | gzip -9 > $BUILD_DIR/initramfs.cpio.gz
popd

## Final clean-up
rm -fr $TMP_DIR

echo -e "\nRun the image..."

qemu-system-x86_64 -kernel $BUILD_DIR/vmlinuz-6.6.8-060608-generic -initrd $BUILD_DIR/initramfs.cpio.gz -nographic -append "console=ttyS0"

