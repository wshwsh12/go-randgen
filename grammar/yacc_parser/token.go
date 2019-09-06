package yacc_parser

import (
	"fmt"
	"io"
	"strings"
	"unicode"
)

type token interface {
	toString() string
}
type eof struct{}

func (*eof) toString() string {
	return "EOF"
}

// ':' or '|'
type operator struct {
	val string
}

func (op *operator) toString() string {
	return op.val
}

type keyword struct {
	val string
}

func (kw *keyword) toString() string {
	return kw.val
}

type nonTerminal struct {
	val string
}

func (nt *nonTerminal) toString() string {
	return nt.val
}

type terminal struct {
	val string
}

func (t *terminal) toString() string {
	return t.val
}

type comment struct {
	val string
}

func (c *comment) toString() string {
	return c.val
}

const (
	inSingQuoteStr = iota + 1
	inDoubleQuoteStr
	inOneLineComment
	inComment
)

func getByQuote(r rune) int {
	if r == '"' {
		return inDoubleQuoteStr
	}
	return inSingQuoteStr
}

type quote struct {
	c int
}

func (q *quote) isInsideStr() bool {
	return q.c == inSingQuoteStr || q.c == inDoubleQuoteStr
}

func (q *quote) isInComment() bool {
	return q.c == inOneLineComment || q.c == inComment
}

func (q *quote) isInOneLineComment() bool {
	return q.c == inOneLineComment
}

func (q *quote) isInSome() bool {
	return q.c != 0
}

func (q *quote) tryToggle(other int) bool {
	if q.c == 0 {
		q.c = other
		return true
	} else if q.c == other {
		q.c = 0
		return true
	}
	return false
}

func skipSpace(reader io.RuneScanner) (r rune, err error) {
	for {
		r, _, err = reader.ReadRune()
		if err != nil {
			return 0, err
		}

		if !unicode.IsSpace(r) {
			return r, nil
		}
	}
}

// Tokenize is used to wrap a reader into a token producer.
// simple lexer not look back, have some problem when quote not pair
func Tokenize(reader io.RuneScanner) func() (token, error) {
	q := quote{0}
	return func() (token, error) {
		var r rune
		var err error
		// Skip spaces.
		r, err = skipSpace(reader)
		if err == io.EOF {
			return &eof{}, nil
		} else if err != nil {
			return nil, err
		}

		// Handle delimiter.
		if r == ':' || r == '|' {
			return &operator{string(r)}, nil
		}

		// Toggle isInsideStr.
		if r == '\'' || r == '"' {
			q.tryToggle(getByQuote(r))
		}

		// handle one line comment
		if r == '#' {
			q.tryToggle(inOneLineComment)
		}

		// Handle stringLiteral or identifier
		var last rune
		var stringBuf string
		stringBuf = string(r)

		for {
			last = r
			r, _, err = reader.ReadRune()
			if err == io.EOF {
				break
			} else if err != nil {
				return nil, err
			}

			// enter comment
			if !q.isInComment() {
				if last == '/' && r == '*' {
					q.tryToggle(inComment)
				}
			}

			if (unicode.IsSpace(r) || isDelimiter(r)) && !q.isInSome() {
				if err := reader.UnreadRune(); err != nil {
					panic(fmt.Sprintf("Unable to unread rune: %s.", string(r)))
				}
				break
			}

			stringBuf += string(r)
			if !q.isInComment() {
				// Handle end str.
				if r == '\'' || r == '"' {
					// identifier can not have ' or "
					if !q.isInsideStr() {
						return nil, fmt.Errorf("unexpected character: `%s` in `%s`", string(r), stringBuf)
					}
					if q.tryToggle(getByQuote(r)) {
						break
					}
				}
			} else {
				// in comment
				if r == '\n' && q.isInOneLineComment() {
					q.tryToggle(inOneLineComment)
					return &comment{stringBuf}, nil
				}
				if last == '*' && r == '/' && q.isInComment() {
					q.tryToggle(inComment)
					return &comment{stringBuf}, nil
				}
			}
		}

		// stringLiteral
		if strings.HasPrefix(stringBuf, "'") || strings.HasPrefix(stringBuf, "\"") {
			return &terminal{stringBuf[1 : len(stringBuf)-1]}, nil
		}

		// keyword
		if strings.HasPrefix(stringBuf, "_") {
			return &keyword{stringBuf}, nil
		}

		// nonTerminal
		if isNonTerminal(stringBuf) {
			return &nonTerminal{stringBuf}, nil
		}

		// terminal
		return &terminal{stringBuf}, nil
	}
}

func isDelimiter(r rune) bool {
	return r == '|' || r == ':'
}

func isNonTerminal(token string) bool {
	for _, c := range token {
		if unicode.IsUpper(c) ||
			(!unicode.IsLetter(c) && c != '_') {
			return false
		}
	}
	return true
}

func isEOF(tkn token) bool {
	_, ok := tkn.(*eof)
	return ok
}

func isComment(tkn token) bool {
	_, ok := tkn.(*comment)
	return ok
}

func isTknNonTerminal(tkn token) bool {
	_, ok := tkn.(*nonTerminal)
	return ok
}
