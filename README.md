GoKilo
======

GoKilo is an attempt to port Antirez's [kilo](http://antirez.com/news/108) text 
editor from C to [Go Langugae](https://golang.org/). To build this, I have followed
this wonderful [tutorial](https://viewsourcecode.org/snaptoken/kilo/index.html)
originally in C language that breaks down Kilo editor's development into a series
of small steps. At present in GoKilo, chapters 
[1 through 6](https://viewsourcecode.org/snaptoken/kilo/index.html) 
of the tutorial are complete and we have a fully functional basic text edtor. 

### Update

GoKilo has been substantially re-written to a much more idiomatic Go style
from the original version that closely mirrored C while retaining
much of the logic. However, two substantial changes were made during the rewrite:
1. A small State Machine added to generalize handling of state changes between 
   editing, save prompt for new files, quit prompt in case of unsaved changes
   and find interaction vs. using embedded screen & keyboard logic
2. A line editor derived from the full Editor used to handle user prompts 
   vs. custom key-handling in find function in the original
3. I am starting development of a tutorial to build a simpler version 
   of GoKilo at github.com/gokilo modeled on the original one

Extra Functionality in GoKilo
-----------------------------
1. Search functionality has been made case-insensitive by default and works slightly differently
2. Native Windows console support - you can use GoKilo in PowerShell or Cmd

Roadmap
-------

1. **Porting over tutorial to Go:** I plan to then re-write and release
   the tutorial to work with Go version of kilo. Given how close Go is 
   to C, it should largely port over in sequence but will probalby be simpler
   with fewer steps due to garbage collection and type safety in Go

2. **Syntax Highlighting**: I haven't yet decided whether to add syntax highlighting to GoKilo


Dog-fooding
-----------
I'm using `gokilo` as my [default git editor](https://stackoverflow.com/questions/2596805/how-do-i-make-git-use-the-editor-of-my-choice-for-commits).

Building
--------
```
go get github.com/srinathh/gokilo
go build
``` 
