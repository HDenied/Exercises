# How to run the program

```
#build the executable
$ make

#build and run unit tests
$ make test

#clean the environment
$ make clean
```

# Running the program

After the program has been built, it will be located in the same directory of the make file, the commands to provide are:
```
$ ./shredder -h

Usage of ./shredder:
  -b int
        The block size in bytes to perform the overwriting (default 4096)
  -f string
        The file path of the file to process
```

# Test Cases

Test case coverage is:
```
ok  	shredcmd/shred	0.354s	coverage: 83.9% of statements
```

Test cases covered are:
```
go test -v ./shred
=== RUN   TestErrorFileNotExistent
--- PASS: TestErrorFileNotExistent (0.00s)
=== RUN   TestErrorNotAfile
--- PASS: TestErrorNotAfile (0.00s)
=== RUN   TestErrorZeroSizeFile
--- PASS: TestErrorZeroSizeFile (0.00s)
=== RUN   TestErrorReadOnlyFile
--- PASS: TestErrorReadOnlyFile (0.00s)
=== RUN   TestErrorNegativeBuffer
--- PASS: TestErrorNegativeBuffer (0.00s)
=== RUN   TestGoodPathDefaultBuffer
--- PASS: TestGoodPathDefaultBuffer (0.01s)
=== RUN   TestGoodPathSmallBuffer
--- PASS: TestGoodPathSmallBuffer (0.01s)
=== RUN   TestGoodPathGTMaxBuff
Length of the buffer limited to 1048576 bytes
--- PASS: TestGoodPathGTMaxBuff (0.02s)
PASS
ok      shredcmd/shred  0.034s
```

I decided to cover tests which are related to errors that can be triggered by human mistakes (e.g. file doesn't exist, file doesn't have the right permission, wrong buffer length,...). 
Some of the system errors (e.g. writing to disk and write error) would have been more complex to test in a unit test and, since they are pretty standard, I assumed that an error will be caught correctly in case they occurred, given the simple logic behind it, therefore I omitted them. 
Probably I could have created a file system on a file, mounted it, corrupted the partition table before the writes and see if this would have triggered the system errors. I think the interest in these last test cases would have been important more at system level to make sure that, for instance, other functionalities in the system wouldn't be affected by this kind of errors. I doubt this would occurr anyway given that I actually had these cases on real hw for similar applications and the error is handled properly with a return code like in this case.
Some performance tests could have been added to find the ideal max buffer size for different page cache but this would have been relevant for production environment.

# What this function is used for
This is a function used when we want make sure that some confidential documents or data in general are deleted permanetly and not recoverable. Possible use cases are: parsonal information, business data and sensitive data (e.g. emails, cache files, browsing history and so on).I believe the big disadvantages for the function are that is just tested on a Linux system and the execution time for big files: more tests should be run for big files in order to find a thread off between iterations, generation of random chunk data and buffer size. If someoene would like to run it on other kind of system, more tests should be executed on the specific platforms. Besides that, it seems to work pretty well.

# Additional observation

The program should work for file of any size, as the file is overwritten in chunks.  By default I decided to make the buffer big as a standard Linux page cache to limit page thrashing, neverthless I left the possibility to select a custom value for the buffer if operation seems to be too slow, specially for the biggest file where too many iterations could create further delay. The max buffer size is limited to 1MB. 
I also made sure that each writing is flushed on the disk each time using the Sync command.