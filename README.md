GoKilo
======

GoKilo is an attempt to port Antirez's [kilo](http://antirez.com/news/108) text 
editor to [Go Langugae](https://golang.org/). To build this, I have followed
this [tutorial](https://viewsourcecode.org/snaptoken/kilo/index.html)
originally in C language that breaks down Kilo editor's development into a series
of small steps. Chapters [1 through 5](https://viewsourcecode.org/snaptoken/kilo/index.html) 
of the tutorial are complete and we have a fully functional basic text edtor. 
The next steps will be adding search and syntax highlighting functionality.


Limitations
-----------
- The codebase is currently very close stylistically to the original C code 
  while I'm working through the C tutorial. It may not be fully in idiomatic
  Go style. I plan to refactor after compeltion.
- Currently, the code will compile only in Linux since I'm using Linux specific
  `syscall.IoctlSetTermios` function call in `terminal.go` to enter and leave
  raw mode in the terminal. I will research a windows version in the future but
  pull requests welcome

Building
--------
`go build` 