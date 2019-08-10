package main 

import (
	"io"
	"os"
	"fmt"
	"unicode"

	// golang syscall main package is deprecated and
	// points to sys/<os> packages to be used instead
	syscall "golang.org/x/sys/unix"
)

/*** data ***/
var origTermios syscall.Termios

/*** terminal ***/

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
	termios.Lflag = termios.Lflag &^ (syscall.ECHO|syscall.ICANON|syscall.ISIG|syscall.IEXTEN)
	termios.Iflag = termios.Iflag &^ (syscall.IXON| syscall.ICRNL|syscall.BRKINT|syscall.INPCK|syscall.ISTRIP)
	termios.Oflag = termios.Oflag &^ (syscall.OPOST)
	termios.Cflag = termios.Cflag | syscall.CS8
	termios.Cc[syscall.VMIN]=0
	termios.Cc[syscall.VTIME]=1
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
		fmt.Fprintf(os.Stderr, "Error: diabling raw mode: %s\r\n", err)
	}
	
	if err == nil{
		os.Exit(0)
	}

	fmt.Fprintf(os.Stderr, "Error: %s\r\n", err)
	os.Exit(1)
}

/*** init ***/

func main(){

	if err := enableRawMode(); err != nil{
		safeExit(err)
	}
	

	b := []byte{0}
	for{
		n, err := os.Stdin.Read(b)		

		switch{
		case (err == io.EOF)||(n==0):
			fmt.Print("No input\r\n")
		case err != nil:
			safeExit(err)
		case  b[0]=='q':
			safeExit(nil)
		default:
			if unicode.IsControl(rune(b[0])){
				fmt.Printf("%d\r\n", b[0])
			}else{
			fmt.Printf("%d (%c)\r\n", b[0], b[0])
			}
		}
	}
}