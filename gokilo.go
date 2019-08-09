package main 

import (
	"os"
	"fmt"
)

func main(){
	b := []byte{0}

	for{
		_, err := os.Stdin.Read(b)		
		if err != nil{
			fmt.Fprintf(os.Stderr,"Error: %s\n", err)
			os.Exit(1)
		}
	}
}