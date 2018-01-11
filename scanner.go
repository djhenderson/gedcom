/*
This is free and unencumbered software released into the public domain. For more
information, see <http://unlicense.org/> or the accompanying UNLICENSE file.
*/

package gedcom

import (
	"fmt"
	"io"
	"log"
	"strconv"
)

// A scanner is a GEDCOM scanning state machine.
type scanner struct {
	parseState int
	tokenStart int
	level      int    // the level
	tag        []byte // tag
	value      []byte // value
	xref       []byte // xref (between level and tag)
}

const (
	stateBegin = iota
	stateLevel
	stateSeekTagOrXref
	stateSeekTag
	stateTag
	stateXref
	stateSeekValue
	stateValue
	stateEnd
	stateError
)

func (s *scanner) reset() {
	s.parseState = stateBegin
	s.tokenStart = 0
	s.level = 0
	s.xref = make([]byte, 0)
	s.tag = make([]byte, 0)
	s.value = make([]byte, 0)
}

func (s *scanner) nextTag(data []byte) (offset int, err error) {

	for i, c := range data {
		switch s.parseState {
		case stateBegin:
			switch {
			case c >= '0' && c <= '9':
				s.tokenStart = i
				s.parseState = stateLevel
			case isSpace(c):
				continue
			default:
				s.parseState = stateError
				err = fmt.Errorf("Found non-whitespace before level")
				log.Println(err.Error())
				return
			}
		case stateLevel:
			switch {
			case c >= '0' && c <= '9':
				continue
			case c == ' ':
				parsedLevel, perr := strconv.ParseInt(string(data[s.tokenStart:i]), 10, 64)
				if perr != nil {
					err = perr
					log.Println(err.Error())
					return
				}
				s.level = int(parsedLevel)
				s.parseState = stateSeekTagOrXref
			default:
				s.parseState = stateError
				err = fmt.Errorf("Level contained non-numerics")
				log.Println(err.Error())
				return
			}

		case stateSeekTag:
			switch {
			case isAlphaNumeric(c):
				s.tokenStart = i
				s.parseState = stateTag
			case c == ' ':
				continue
			default:
				s.parseState = stateError
				err = fmt.Errorf("Tag \"%s\" contained non-alphanumeric", string(data[s.tokenStart:i]))
				log.Println(err.Error())
				return
			}
		case stateSeekTagOrXref:
			switch {
			case isAlphaNumeric(c):
				s.tokenStart = i
				s.parseState = stateTag
			case c == '@':
				s.tokenStart = i
				s.parseState = stateXref
			case c == ' ':
				continue
			default:
				s.parseState = stateError
				err = fmt.Errorf("Xref \"%s\" contained non-alphanumeric", string(data[s.tokenStart:i]))
				log.Println(err.Error())
				return
			}

		case stateTag:
			switch {
			case isAlphaNumeric(c):
				continue
			case c == '\n' || c == '\r':
				s.tag = data[s.tokenStart:i]
				s.parseState = stateEnd
				offset = i
				return
			case (c == ' ') || (c == '\t'):
				s.tag = data[s.tokenStart:i]
				s.parseState = stateSeekValue
			default:
				s.parseState = stateError
				err = fmt.Errorf("Tag contained non-alphanumeric")
				log.Println(err.Error())
				return
			}

		case stateXref:
			switch {
			case isAlphaNumeric(c) || c == '@':
				continue
			case c == ' ':
				s.xref = data[s.tokenStart+1 : i-1]
				s.parseState = stateSeekTag
			default:
				s.parseState = stateError
				err = fmt.Errorf("Xref contained non-alphanumeric")
				log.Println(err.Error())
				return
			}
		case stateSeekValue:
			switch {
			case c == '\n' || c == '\r':
				s.parseState = stateEnd
				offset = i
				return
			case c == ' ':
				continue
			default:
				s.tokenStart = i
				s.parseState = stateValue
			}

		case stateValue:
			switch {
			case c == '\n' || c == '\r':
				s.value = data[s.tokenStart:i]
				s.parseState = stateEnd
				offset = i
				return
			default:
				continue
			}

		case stateEnd:
			break

		default:
			panic("what state are we in?")
		}
	}

	return 0, io.EOF
}

func isSpace(c byte) bool {
	return c == ' ' || c == '\t' || c == '\r' || c == '\n'
}

func isAlphaNumeric(c byte) bool {
	return (c >= '0' && c <= '9') || (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || c == '_'
}
