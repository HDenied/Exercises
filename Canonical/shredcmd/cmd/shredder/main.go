package main

import (
	"flag"
	"fmt"
	"shredcmd/shred"
)

func main() {
	filePath := flag.String("f", "", "The file path of the file to process")
	blockSize := flag.Int("b", shred.DEF_BUFF, "The block size in bytes to perform the overwriting")

	flag.Parse()

	// Shred the file
	res, err := shred.Shred(*filePath, *blockSize)
	if err != nil {
		fmt.Println("Error:", err)
		return
	} else {
		fmt.Printf("Shredding completed successfully: overwritten %d bytes in %d iterations with block size of %d bytes\n", res.BytesWritten, res.Iteration, res.BlockSize)
	}

}
