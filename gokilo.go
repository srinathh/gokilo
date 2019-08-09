package main 

import (
	"os"
	"fmt"
	syscall "golang.org/x/sys/unix"
)

func enableRawMode() error {

	termios, err := syscall.IoctlGetTermios(syscall.Stdin, syscall.TCGETS)
	if err != nil{
		return err
	} 

	termios.Lflag = termios.Lflag &^syscall.ECHO

	if err := syscall.IoctlSetTermios(syscall.Stdin, syscall.TCSETS, termios); err != nil{
		return err
	}

	return nil
}

func main(){

	if err := enableRawMode(); err != nil{
		fmt.Fprintln(os.Stderr, err)
	}
	

	b := []byte{0}
	for{
		_, err := os.Stdin.Read(b)		
		if err != nil{
			fmt.Fprintf(os.Stderr,"Error: %s\n", err)
			os.Exit(1)
		}

		if b[0]=='q'{
			break
		}
	}
}