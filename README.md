GoKilo
======

GoKilo is an attempt to port Antirez's [kilo](http://antirez.com/news/108) text 
editor from C to [Go Langugae](https://golang.org/). To build this, I have followed
this wonderful [tutorial](https://viewsourcecode.org/snaptoken/kilo/index.html)
originally in C language that breaks down Kilo editor's development into a series
of small steps. At present, chapters [1 through 5](https://viewsourcecode.org/snaptoken/kilo/index.html) 
of the tutorial are complete and we have a fully functional basic text edtor. 

Roadmap
-------
1. **Add Search Functionality:** I plan to first add Search functionality
   from [chapter 6](https://viewsourcecode.org/snaptoken/kilo/06.search.html) 
   since search is a core functionality in a text editor.
   
2. **Refactoring to Idiomatic Go:** I then plan to refactor the code from a 
   C style to a more idiomatic Go style using Go paradigms

3. **Porting over tutorial to Go:** I plan to then re-write and release
   the tutorial to work with Go version of kilo. Given how close Go is 
   to C, it should largely port over in sequence but will probalby be simpler
   with fewer steps due to garbage collection and type safety in Go

4. **Syntax Highlighting**: I haven't yet decided whether to add syntax highlighting to GoKilo

Limitations
-----------
The code currently compiles in Linux only as I'm using a Linux specific
System Call `syscall.IoctlSetTermios` in `terminal.go` to enter and leave
terminal raw mode. I haven't yet researched the Windows equivalent.

Pull requests to add Windows support welcome. Essentially only `enableRawMode()`
and `disableRawMode()` need to be ported over.

Dog-fooding
-----------
I'm using `gokilo` as my [default git editor](https://stackoverflow.com/questions/2596805/how-do-i-make-git-use-the-editor-of-my-choice-for-commits).

Building
--------
```
go get github.com/srinathh/gokilo
go build
``` 