package main 

import (
	"io"
	"os"
	"errors"
	"fmt"

	// golang syscall main package is deprecated and
	// points to sys/<os> packages to be used instead
	syscall "golang.org/x/sys/unix"
)
/*** defines ***/
func ctrlKey(b byte) byte{
	return b&0x1f
}

/*** data ***/

type editorConfig struct{
	screenRows int
	screenCols int
	origTermios syscall.Termios
}

var cfg editorConfig

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

	cfg.origTermios = *termios

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
	if err := syscall.IoctlSetTermios(syscall.Stdin, syscall.TCSETSF, &cfg.origTermios); err != nil{
		return err
	}
	return nil
}

func safeExit(err error){
	fmt.Fprint(os.Stdout, "\x1b[2J")
	fmt.Fprint(os.Stdout, "\x1b[H")

	if err1 := disableRawMode(); err1 != nil{
		fmt.Fprintf(os.Stderr, "Error: diabling raw mode: %s\r\n", err)
	}
	
	if err == nil{
		os.Exit(0)
	}

	fmt.Fprintf(os.Stderr, "Error: %s\r\n", err)
	os.Exit(1)
}

// single space buffer to reduce allocations
var keyBuf = []byte{0}
func editorReadKey() (byte, error){

	for {
		n, err := os.Stdin.Read(keyBuf)		
		switch{
		case err == io.EOF:
			continue
		case err != nil:
			return 0, err
		case n==0:
			continue
		default:
			return keyBuf[0], nil
		}
	}
}

func getWindowSize() (int, int, error){

	ws, err := syscall.IoctlGetWinsize(syscall.Stdout, syscall.TIOCGWINSZ)
	if err != nil {
		return 0, 0, err
	}
	if ws.Row == 0 || ws.Col == 0{
		return 0, 0, errors.New("got non zero column or row")
	}

	return int(ws.Row), int(ws.Col), nil;

}

/*** output ***/
func editorRefreshScreen(){
	// clear screen
	fmt.Fprint(os.Stdout, "\x1b[2J")
	fmt.Fprint(os.Stdout, "\x1b[H")

	editorDrawRows()

	// reposition cursor
	fmt.Fprint(os.Stdout, "\x1b[H")


}

func editorDrawRows(){
	for j := 0; j < cfg.screenRows; j++{
		fmt.Fprint(os.Stdout,"~\r\n")
	}
}

/*** Input ***/

func editorProcessKeypress()error{

	b, err := editorReadKey()
	if err != nil{
		return err
	}

	switch(b){
	case ctrlKey('q'):
		safeExit(nil)
	}
	return nil
}

/*** init ***/

func initEditor() error {
	rows, cols, err := getWindowSize(); 
	if err != nil{
		return err
	}
	cfg.screenRows = rows
	cfg.screenCols = cols
	return nil
}

func main(){

	if err := enableRawMode(); err != nil{
		safeExit(err)
	}

	if err := initEditor(); err != nil{
		safeExit(err)
	}

	for{
		editorRefreshScreen()
		if err := editorProcessKeypress(); err != nil{
			safeExit(err)
		}
	}
}