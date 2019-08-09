package main 

import (
	"os"
	"fmt"

	// golang syscall main package is deprecated and
	// points to sys/<os> packages to be used instead
	syscall "golang.org/x/sys/unix"
)

var origTermios syscall.Termios

// enableRawMode switches from cooked or canonical mode to raw mode
// by using syscalls. Currently this is the implrementation for Unix only
func enableRawMode() error {

	// Gets TermIOS data structure. From glibc, we find the cmd should be TCGETS
	// https://code.woboq.org/userspace/glibc/sysdeps/unix/sysv/linux/tcgetattr.c.html
	termios, err := syscall.IoctlGetTermios(syscall.Stdin, syscall.TCGETS)
	if err != nil{
		return err
	} 

	origTermios = *termios

	// turn off echo & canonical mode by using a bitwise clear operator &^
	termios.Lflag = termios.Lflag &^ (syscall.ECHO|syscall.ICANON)

	// We from the code of tcsetattr in glibc, we find that for TCSAFLUSH, 
	// the corresponding command is TCSETSF 
	// https://code.woboq.org/userspace/glibc/sysdeps/unix/sysv/linux/tcsetattr.c.html
	if err := syscall.IoctlSetTermios(syscall.Stdin, syscall.TCSETSF, termios); err != nil{
		return err
	}

	return nil
}

func disableRawMode() error{
	if err := syscall.IoctlSetTermios(syscall.Stdin, syscall.TCSETSF, &origTermios); err != nil{
		return err
	}
	return nil
}

func safeExit(err error){
	if err1 := disableRawMode(); err1 != nil{
		fmt.Fprintf(os.Stderr, "Error: diabling raw mode: %s\n\r", err)
	}
	
	if err == nil{
		os.Exit(0)
	}

	fmt.Fprintf(os.Stderr, "Error: %s\n\r", err)
	os.Exit(1)
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
			safeExit(nil)
		}
	}
}