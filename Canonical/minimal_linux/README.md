# General info

The script can be run from any directory with the command:

```
./minimal_image.sh
```

It follows the requirements:
- No other external or system directory are used, only the cwd.
- No sudo command is necessary to run the script
- It prints hello world
- It comes with basic functionality of a linux distribution

The only assumptions is that qemu and git are installed on a system.

The "minimal_image.sh" script performs the following steps:

- Extract a precompiled kernel to a local directory from a .deb package
- A compiled busybox and an init file are pulled from my repo
- A file system is generated
- Busybox executables and the init file are copied over the generated file system
- The fs is compressed in a cpio.gz file
- qemu command is executed


I precompiled the BusyBox image to speed up the process of spawning the system, all commands I used
are in the comments of "minimal_image.sh". For the same reason I provided an init file where I also
generate essential system file node.