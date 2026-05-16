package interpreter

import (
	"bufio"
	"fmt"
	"os"
	"slices"
	"strconv"
	"strings"
	"time"
)

type variable struct {
	val      []int
	is_const bool
}

type Interpreter struct {
	pos          int
	vars         map[string]*variable
	labels       map[string]int //name to token pos
	instructions map[int]func(i *Interpreter) error
	Tokens       []string
	stack        []int
	callstack    []int
}

func (i *Interpreter) push(v []int) {
	i.stack = append(i.stack, v...)
}
func (i *Interpreter) pop() int {
	var stack_len = len(i.stack)
	if stack_len == 0 {
		fmt.Printf(`   _
  ;'> - << stack underflow! >>
."..`)
		return 0
	}
	var x int
	x, i.stack = i.peek(), i.stack[:stack_len-1]
	return x
}
func (i *Interpreter) expect(n int) error {
	popped := i.pop()
	if popped != n {
		return fmt.Errorf("Expected %d got %d", n, popped)
	}
	return nil
}
func (i *Interpreter) peek() int {
	return i.stack[len(i.stack)-1]
}
func (i *Interpreter) dup() {
	i.push([]int{i.peek()})
}

func (i *Interpreter) readString() string {
	var sb strings.Builder

	for i.peek() != 0 { //c style strings with nullbytes
		sb.WriteByte(byte(i.pop()))
	}
	i.pop() //pop the zero

	return sb.String()

}
func (i *Interpreter) readBuffer() ([]int, error) {

	b := []int{}

	open_bracket := i.vars["{"].val[0]
	closing_bracket := i.vars["}"].val[0]

	if r := i.expect(closing_bracket); r != nil {
		return []int{}, fmt.Errorf("Expected closing bracket, got %d", r)
	}
	for i.peek() != open_bracket {
		b = append(b, i.pop())
	}
	i.pop()

	slices.Reverse(b)
	return b, nil

}

func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}

func NewInterpreter() *Interpreter {
	var result Interpreter = Interpreter{
		vars: map[string]*variable{
			"nop":  {[]int{0}, true},
			"dup":  {[]int{1}, true},
			"pop":  {[]int{2}, true},
			"swap": {[]int{3}, true},
			"over": {[]int{4}, true},
			"len":  {[]int{5}, true},
			"{":    {[]int{6}, true},
			"}":    {[]int{7}, true},

			"+": {[]int{10}, true},
			"-": {[]int{11}, true},
			"*": {[]int{12}, true},
			"/": {[]int{13}, true},
			"%": {[]int{14}, true},
			"=": {[]int{15}, true},
			">": {[]int{16}, true},
			"<": {[]int{17}, true},

			"set":      {[]int{20}, true},
			"var":      {[]int{21}, true},
			"alias":    {[]int{22}, true},
			"buf":      {[]int{23}, true},
			"bufalias": {[]int{24}, true},
			"setbuf":   {[]int{25}, true},
			"del":      {[]int{26}, true},
			"get":      {[]int{27}, true},

			"goto":  {[]int{40}, true},
			"gosub": {[]int{41}, true},
			"ret":   {[]int{42}, true},
			"jz":    {[]int{43}, true},
			"jnz":   {[]int{44}, true},

			"outn":   {[]int{50}, true},
			"outc":   {[]int{51}, true},
			"outs":   {[]int{52}, true},
			"outb":   {[]int{53}, true},
			"inputs": {[]int{54}, true},
			"inputn": {[]int{55}, true},
			"inputc": {[]int{56}, true},

			"sleep": {[]int{60}, true},
			"time":  {[]int{61}, true},
		},
		labels: map[string]int{},
		instructions: map[int]func(i *Interpreter) error{
			1: func(i *Interpreter) error { i.dup(); return nil },
			2: func(i *Interpreter) error { i.pop(); return nil },
			3: func(i *Interpreter) error { a, b := i.pop(), i.pop(); i.push([]int{a, b}); return nil },
			4: func(i *Interpreter) error { i.push([]int{i.stack[len(i.stack)-2]}); return nil },
			5: func(i *Interpreter) error { i.push([]int{len(i.stack)}); return nil },

			10: func(i *Interpreter) error { b, a := i.pop(), i.pop(); i.push([]int{a + b}); return nil },
			11: func(i *Interpreter) error { b, a := i.pop(), i.pop(); i.push([]int{a - b}); return nil },
			12: func(i *Interpreter) error { b, a := i.pop(), i.pop(); i.push([]int{a * b}); return nil },
			13: func(i *Interpreter) error { b, a := i.pop(), i.pop(); i.push([]int{a / b}); return nil },
			14: func(i *Interpreter) error { b, a := i.pop(), i.pop(); i.push([]int{a % b}); return nil },
			15: func(i *Interpreter) error { b, a := i.pop(), i.pop(); i.push([]int{boolToInt(a == b)}); return nil },
			16: func(i *Interpreter) error { b, a := i.pop(), i.pop(); i.push([]int{boolToInt(a > b)}); return nil },
			17: func(i *Interpreter) error { b, a := i.pop(), i.pop(); i.push([]int{boolToInt(a < b)}); return nil },
			20: func(i *Interpreter) error {
				name := strings.ToLower(i.readString())
				v, ok := i.vars[name]
				if !ok {
					return fmt.Errorf("Variable %s does not exist", name)
				}
				if v.is_const {
					return fmt.Errorf("Cannot assign to constant variable %s", name)
				}
				val := i.pop()

				v.val = []int{val}

				return nil
			},
			21: func(i *Interpreter) error {
				name := strings.ToLower(i.readString())
				_, ok := i.vars[name]
				if ok {
					return fmt.Errorf("Variable %s already exists", name)
				}
				val := i.pop()
				i.vars[name] = &variable{val: []int{val}, is_const: false}

				return nil
			},
			22: func(i *Interpreter) error {
				name := strings.ToLower(i.readString())
				_, ok := i.vars[name]
				if ok {
					return fmt.Errorf("Variable %s already exists", name)
				}
				val := i.pop()
				i.vars[name] = &variable{val: []int{val}, is_const: true}

				return nil
			},
			23: func(i *Interpreter) error {
				name := strings.ToLower(i.readString())
				_, ok := i.vars[name]
				if ok {
					return fmt.Errorf("Variable %s already exists", name)
				}
				val, err := i.readBuffer()
				if err != nil {
					return err
				}

				i.vars[name] = &variable{val: val, is_const: false}

				return nil
			},
			24: func(i *Interpreter) error {
				name := strings.ToLower(i.readString())
				_, ok := i.vars[name]
				if ok {
					return fmt.Errorf("Variable %s already exists", name)
				}
				val, err := i.readBuffer()
				if err != nil {
					return err
				}
				i.vars[name] = &variable{val: val, is_const: true}

				return nil
			},
			25: func(i *Interpreter) error {
				name := strings.ToLower(i.readString())
				v, ok := i.vars[name]
				if !ok {
					return fmt.Errorf("Variable %s does not exist", name)
				}
				if v.is_const {
					return fmt.Errorf("Cannot assign to constant variable %s", name)
				}
				val, err := i.readBuffer()
				if err != nil {
					return err
				}

				v.val = val

				return nil
			},
			26: func(i *Interpreter) error {
				name := strings.ToLower(i.readString())
				_, ok := i.vars[name]
				if !ok {
					return fmt.Errorf("Variable %s does not exist", name)
				}
				delete(i.vars, name)

				return nil
			},
			27: func(i *Interpreter) error {
				name := strings.ToLower(i.readString())
				v, ok := i.vars[name]
				if !ok {
					return fmt.Errorf("Variable %s does not exist", name)
				}
				i.push(v.val)
				return nil
			},
			40: func(i *Interpreter) error {
				name := strings.ToLower(i.readString())
				v, ok := i.labels[name]
				if !ok {
					return fmt.Errorf("Label %s does not exist", name)
				}
				i.pos = v
				return nil
			},
			41: func(i *Interpreter) error {
				name := strings.ToLower(i.readString())
				v, ok := i.labels[name]
				if !ok {
					return fmt.Errorf("Label %s does not exist", name)
				}
				i.callstack = append(i.callstack, i.pos)
				i.pos = v
				return nil
			},
			42: func(i *Interpreter) error {
				i.pos = i.callstack[len(i.callstack)-1]
				i.callstack = i.callstack[:len(i.callstack)-1]
				return nil
			},
			43: func(i *Interpreter) error {
				name := strings.ToLower(i.readString())
				v, ok := i.labels[name]
				if !ok {
					return fmt.Errorf("Label %s does not exist", name)
				}
				cond := i.pop()
				if cond == 0 {
					i.pos = v
				}
				return nil
			},
			44: func(i *Interpreter) error {
				name := strings.ToLower(i.readString())
				v, ok := i.labels[name]
				if !ok {
					return fmt.Errorf("Label %s does not exist", name)
				}
				cond := i.pop()
				if cond != 0 {
					i.pos = v
				}
				return nil
			},
			50: func(i *Interpreter) error { val := i.pop(); fmt.Printf("%d", val); return nil },
			51: func(i *Interpreter) error { val := i.pop(); fmt.Printf("%c", val); return nil },
			52: func(i *Interpreter) error { val := i.readString(); fmt.Printf("%s", val); return nil },
			53: func(i *Interpreter) error {
				val, err := i.readBuffer()
				if err != nil {
					return err
				}
				fmt.Printf("{ ")
				for _, v := range val {
					fmt.Printf("%d, ", v)
				}
				fmt.Printf(" }")
				return nil
			},
			54: func(i *Interpreter) error {
				var input string
				scanner := bufio.NewScanner(os.Stdin)
				if scanner.Scan() {
					input = scanner.Text()
				}
				i.push([]int{0})
				content := []rune(input)
				slices.Reverse(content)
				i.push(runesToInts(content))
				return nil
			},
			55: func(i *Interpreter) error {
				var input int
				scanner := bufio.NewScanner(os.Stdin)
				if scanner.Scan() {
					var err error
					input, err = strconv.Atoi(scanner.Text())
					if err != nil {
						return err
					}
				}
				i.push([]int{input})
				return nil
			},
			56: func(i *Interpreter) error {
				var input rune
				reader := bufio.NewReader(os.Stdin)
				var err error
				input, _, err = reader.ReadRune()
				if err != nil {
					return err
				}
				i.push([]int{int(input)})
				return nil
			},

			60: func(i *Interpreter) error {
				t := i.pop()
				time.Sleep(time.Duration(t) * time.Millisecond)
				return nil
			},
			61: func(i *Interpreter) error {
				i.push([]int{int(time.Now().UnixMilli())})
				return nil
			},
		},
		Tokens: []string{},
		stack:  []int{},
	}

	return &result
}

func runesToInts(runes []rune) []int {
	intSlice := make([]int, len(runes))
	for i, r := range runes {
		intSlice[i] = int(r)
	}
	return intSlice
}

func isInteger(s string) bool {
	_, err := strconv.Atoi(s)
	return err == nil
}
func (i *Interpreter) do() error {
	val := i.pop()
	op, ok := i.instructions[val]
	if !ok {
		return fmt.Errorf("Attempted to execute nonexistent instruction %d", val)
	}
	err := op(i)
	if err != nil {
		return err
	}
	return nil
}

func (i *Interpreter) Interpret(Debug bool, Delay time.Duration) ([]int, error) {
	i.pos = 0
	for i.pos < len(i.Tokens) {
		tok := i.Tokens[i.pos]
		if tok[0] == '@' {
			if tok[len(tok)-1] != ':' {
				return i.stack, fmt.Errorf("Unfinished label declaration: %q. No colon found", tok)
			}

			i.labels[tok[1:len(tok)-1]] = i.pos
		}
		i.pos++
	}
	i.pos = 0
	for i.pos < len(i.Tokens) {
		if Debug {
			fmt.Printf("STACK PRE: %v\n", i.stack)
			time.Sleep(Delay)
		}
		tok := i.Tokens[i.pos]
		switch {
		case tok[0] == '"':
			if tok[len(tok)-1] != '"' {
				return i.stack, fmt.Errorf("Unfinished string declaration: %q. No quote found", tok)
			}
			i.push([]int{0})
			content := []rune(tok[1 : len(tok)-1])
			slices.Reverse(content)
			i.push(runesToInts(content))
		case tok[0] == '@':
		case isInteger(tok):
			v, _ := strconv.Atoi(tok)
			i.stack = append(i.stack, v)
		case tok == "do":
			err := i.do()
			if err != nil {
				return i.stack, err
			}
		default:
			v, ok := i.vars[tok]
			if !ok {
				return i.stack, fmt.Errorf("Found bad token: %q", tok)
			}
			i.push(v.val)
		}

		i.pos++
		if Debug {
			fmt.Printf("STACK POST: %v\n", i.stack)
		}
	}
	return i.stack, nil
}
