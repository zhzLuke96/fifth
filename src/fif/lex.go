package main

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

func Parse(input []byte) error {
	l := newLex(input)
	_ = FifParse(l)
	return l.err
}

type lex struct {
	input   []byte
	pos     int
	linepos int
	line    int
	lastT   string
	err     error
}

func newLex(input []byte) *lex {
	return &lex{
		input: input,
	}
}

// Lex satisfies yyLexer.
func (l *lex) Lex(lval *FifSymType) int {
	return l.scanNormal(lval)
}

var idRegex = regexp.MustCompile("[\\w\\$_]")
var needLook = regexp.MustCompile("[=+\\-*/<>]")

func (l *lex) scanNormal(lval *FifSymType) int {
	for b := l.next(); b != 0; b = l.next() {
		switch {
		case b == '\n':
			l.line++
			l.linepos = 0
			continue
		case unicode.IsSpace(rune(b)):
			continue
		case b == '"' || b == '\'':
			return l.scanString(lval, b)
		case unicode.IsDigit(rune(b)):
			l.backup()
			return l.scanNum(lval)
		// case unicode.IsLetter(rune(b)):
		case idRegex.Match([]byte{b}):
			l.backup()
			return l.scanIdentifier(lval)
		case needLook.Match([]byte{b}) && !l.nxtIsBreakWord():
			l.backup()
			return l.scanOpt(lval)
		default:
			l.lastT = string(b)
			return int(b)
		}
	}
	return 0
}

var escape = map[byte]byte{
	'"':  '"',
	'\\': '\\',
	'/':  '/',
	'b':  '\b',
	'f':  '\f',
	'n':  '\n',
	'r':  '\r',
	't':  '\t',
}

func (l *lex) scanString(lval *FifSymType, match byte) int {
	buf := bytes.NewBuffer(nil)
	for b := l.next(); b != 0; b = l.next() {
		switch b {
		case '\\':
			// TODO(sougou): handle \uxxxx construct.
			b2 := escape[l.next()]
			if b2 == 0 {
				return LexError
			}
			buf.WriteByte(b2)
		case match:
			lval.val = buf.String()
			l.lastT = buf.String()
			return StringConstant
		default:
			buf.WriteByte(b)
		}
	}
	return LexError
}

func (l *lex) scanNum(lval *FifSymType) int {
	buf := bytes.NewBuffer(nil)
	for {
		b := l.next()
		switch {
		case unicode.IsDigit(rune(b)):
			buf.WriteByte(b)
		case strings.IndexByte(".+-eE", b) != -1:
			buf.WriteByte(b)
		default:
			l.backup()
			l.lastT = buf.String()
			val, err := strconv.ParseFloat(buf.String(), 64)
			if err != nil {
				return LexError
			}
			lval.val = val
			return NumConstant
		}
	}
}

var reserved_words = map[string]int{
	"func":        FuncDefined,
	"return":      FuncReturn,
	"if":          T_IF,
	"else":        T_ELSE,
	"then":        T_THEN,
	"true":        T_TRUE,
	"false":       T_FALSE,
	">=":          T_GE,
	"<=":          T_LE,
	"==":          T_EQ,
	"for":         T_FOR,
	"while":       T_WHILE,
	"__fifcode__": T_FIF,
}

func (l *lex) scanIdentifier(lval *FifSymType) int {
	buf := bytes.NewBuffer(nil)
	for {
		b := l.next()
		switch {
		case idRegex.Match([]byte{b}):
			buf.WriteByte(b)
		// case unicode.IsSpace(rune(b)):
		default:
			l.backup()
			lval.val = buf.String()
			l.lastT = buf.String()
			if ty, ok := reserved_words[lval.val.(string)]; ok {
				// fmt.Printf("\nty:[%v] [%v]\n", ty, lval.val.(string))
				return ty
			}
			return Identifier
		}
	}
}

func (l *lex) scanOpt(lval *FifSymType) int {
	buf := bytes.NewBuffer(nil)
	b := l.next()
	buf.WriteByte(b)
	b = l.next()
	buf.WriteByte(b)

	lval.val = buf.String()
	l.lastT = buf.String()
	if ty, ok := reserved_words[lval.val.(string)]; ok {
		// fmt.Printf("\nty:[%v] [%v]\n", ty, lval.val.(string))
		return ty
	}
	return LexError
}

var breakWords = regexp.MustCompile("[[\\](){}.;\\s]")

func (l *lex) nxtIsBreakWord() bool {
	nxt := l.next()
	l.backup()
	return breakWords.Match([]byte{nxt})
}

func (l *lex) backup() {
	if l.pos == -1 {
		return
	}
	l.pos--
	l.linepos--
}

func (l *lex) next() byte {
	if l.pos >= len(l.input) || l.pos == -1 {
		l.pos = -1
		l.linepos = -1
		return 0
	}
	l.pos++
	l.linepos++
	return l.input[l.pos-1]
}

// Error satisfies yyLexer.
func (l *lex) Error(s string) {
	l.err = errors.New(s)
	fmt.Fprintf(os.Stderr, "[ERROR]:%v pos:[L:%v,P:%v] token:[%v]", s, l.line+1, l.linepos, l.lastT)
}
