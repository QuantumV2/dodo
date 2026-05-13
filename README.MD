# DODO

DODO is a stack based retro-style programming language inspired by FORTH, BASIC, etc. \
(Almost) everything in this language is a number on the stack, except minor stuff like labels, variable storage, and DO

DO is used to execute every command because everything just pushes its opcode.

## Contributing
We are always free for contributions, suggestions, and feedback! 
Submit issues and pull requests to our codeberg repo instead of github. https://codeberg.org/QuantumV/dodo

## Installation 
1. Run `go install codeberg.org/QuantumV/dodo/cmd/dodo`
2. Ensure your go binaries path is in your PATH 
3. Run `dodo help` for info
4. Optionally, run `dodo init-lib` and move the `util` directory into the directory it gives you, so you can run code that depends on it!
5. Try writing some code!

## Notes
### Strings
Strings work by first pushing 0 and then pushing themselves in reverse for convenient printing and such. So like reversed C Strings.
### Buffers
Buffers are multiple ints enclosed by `{` (80) and `}` (81). When they are pushed by a variable, the brackets are excluded.

## Instructions 
### Misc
| Instruction | Opcode | Description |
|:------------|:------:|:------------|
| DO  | N/A | Executes the top instruction on the stack. Technically is not an instruction |
| @NAME: | N/A | Creates a label with a name, no name is also valid. This has a separate pass so you can safely jump over it. |
| IMPORT | N/A | Replaces itself and the preceding string by the contents of the file pointed to by the string, relative to the location where the command is run. Falls back to the DODO lib folder if it cant find it. |
| ; | N/A | Initializes a comment from this token to the end of the line, the comment is not tokenized | 
| VARIABLENAME | ~ | Pushes the value of the variable named, can push mutliple values.
| NOP | 0 | Does nothing. At all. |
### Stack Ops
| Instruction | Opcode | Description |
|:------------|:------:|:------------|
| DUP | 1 | Duplicate top value on the stack |
| POP | 2 | Discard top value on the stack |
| SWAP| 3 | Swap two top values on the stack |
| OVER| 4 | Push the value before the top value on the stack
### Math Ops
B is the top of the stack, A is the element before the top
| Instruction | Opcode | Description |
|:------------|:------:|:------------|
| + | 10 | Performs A + B
| - | 11 | Performs A - B
| * | 12 | Performs A * B
| / | 13 | Performs A / B
| % | 14 | Performs A % B
| = | 15 | Checks if A is equal to B and pushes 1 if it is, and 0 if it isnt.
| > | 16 | Checks if A is more than B and pushes 1 if it is, and 0 if it isnt.
| < | 17 | Checks if A is less than B and pushes 1 if it is, and 0 if it isnt.
### Var Ops
| Instruction | Opcode | Description |
|:------------|:------:|:------------|
| SET | 20 | Sets a variable pointed to by the preceding string to a value. E.g. `10 "HI" SET DO` |
| VAR | 21 | Initialize a numeric var. E.g. `10 "HI" VAR DO`   |
| ALIAS | 22 | Initialize a constant numeric var. E.g. `10 "HI" ALIAS DO`   |
| BUF | 23 | Initialize a buffer var that will push the contents of the buffer on use. E.g. `{ 1 2 3 4 } "HI" BUF DO`   |
| BUFALIAS | 24 | Initialize a constant buffer var that will push the contents of the buffer on use. E.g. `{ 1 2 3 4 }  "HI" BUFALIAS DO`   |
| SETBUF | 25 | Sets a buffer variable. E.g. `{ 1 2 3 4 } "HI" SETBUF DO` |
| DEL | 26 | Permanently delete a variable or alias. This has no restrictions. E.g. `"HI" DEL DO`|
| GET | 27 | Pushes a variable value to the stack by name, useful for getting variables generated at runtime |


### Jump Ops
| Instruction | Opcode | Description |
|:------------|:------:|:------------|
| GOTO | 40 | Jumps to the label pointed to by the preceding string. E.g. `"MYLABEL" GOTO` | 
| GOSUB | 41 | Jumps to a label, but also pushes the position to the call stack, so it can be returned to by...  |
| RET | 42 | Jumps back to the caller |
| JZ | 43 |  Jumps to a label, but only if the top of the stack is 0, also consumes it.
| JNZ | 44 | Jumps to a label, but only if the top of the stack is 1, also consumes it.
### I/O Ops
| Instruction | Opcode | Description |
|:------------|:------:|:------------|
| OUTN | 50 | Output top value as a number |
| OUTC | 51 | Output top value as a character |
| OUTS | 52 | Output top value as a string |
| OUTB | 53 | Output top value as a buffer (`{ 1, 2, 3, }`) |
| INPUTS | 54 | Input value as a string  |
| INPUTN | 55 | Input value as a number  |
| INPUTC | 56 | Input value as a character  |
### Time Ops 
| Instruction | Opcode | Description |
|:------------|:------:|:------------|
| SLEEP | 60 | Sleep for an amount of ms |
| TIME | 61 | Get unix time milliseconds |

### Default Consts
| Const | Value | Description |
|:------------|:------:|:------------|
|{|80|Used for defining buffers |
|}|81| Ditto |



```
  _    
<=' < Hello, I'm DODO! >   
  /__  
  \__".
   //  
```