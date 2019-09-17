package postgres

import (
	"bytes"
	"database/sql"
	"database/sql/driver"
	"fmt"
	"strconv"
	"strings"
)

// Cube is a type that is used to work with the postgresql cube extension
// which enables the use of euclidian operators on a matrix of sorts. As
// of writing this, the only supported use case is a 1xn matrix, since
// that is all we support in faces at the moment. For more information
// on cube see: https://www.postgresql.org/docs/9.6/cube.html.
type Cube []float64

var (
	_ driver.Valuer = Cube{}
	_ sql.Scanner   = (*Cube)(nil)
)

// Value is the cube value to convert the cube into the postgreql wire
// protocol. It is was inspired by the github.com/lib/pq.FloatArray type.
func (c Cube) Value() (driver.Value, error) {
	if c == nil {
		return nil, nil
	}

	if n := len(c); n > 0 {
		// There will be at least two parenthesis, N bytes of values,
		// and N-1 bytes of delimiters.
		b := make([]byte, 0, 2*n)

		b = strconv.AppendFloat(b, c[0], 'f', -1, 64)
		for i := 1; i < n; i++ {
			b = append(b, ',')
			b = strconv.AppendFloat(b, c[i], 'f', -1, 64)
		}

		return fmt.Sprintf("(%s)", string(b)), nil
	}

	return "()", nil
}

// Scan converts the postgesql cube bytes into a valid cube type. This
// again was inspired by github.com/lib/pq.FloatArray type.
func (c *Cube) Scan(src interface{}) error {
	switch src := src.(type) {
	case []byte:
		return c.scanBytes(src)
	case string:
		return c.scanBytes([]byte(src))
	case nil:
		*c = nil
		return nil
	}

	return fmt.Errorf("pq: cannot convert %T to CUBE", src)
}

func (c *Cube) scanBytes(src []byte) error {
	elems, err := scanLinearArray(src, []byte{','})
	if err != nil {
		return err
	}
	if *c != nil && len(elems) == 0 {
		*c = (*c)[:0]
	} else {
		b := make(Cube, len(elems))
		for i, v := range elems {
			if b[i], err = strconv.ParseFloat(strings.TrimSpace(string(v)), 64); err != nil {
				return fmt.Errorf("pq: parsing array element index %d: %v", i, err)
			}
		}
		*c = b
	}
	return nil
}

func scanLinearArray(src, del []byte) (elems [][]byte, err error) {
	dims, elems, err := parseArray(src, del)
	if err != nil {
		return nil, err
	}
	if len(dims) > 1 {
		return nil, fmt.Errorf("pq: cannot convert ARRAY%s to CUBE", strings.Replace(fmt.Sprint(dims), " ", "][", -1))
	}
	return elems, err
}

func parseArray(src, del []byte) (dims []int, elems [][]byte, err error) {
	var depth, i int

	if len(src) < 1 || src[0] != '(' {
		return nil, nil, fmt.Errorf("pq: unable to parse array; expected %q at offset %d", '(', 0)
	}

Open:
	for i < len(src) {
		switch src[i] {
		case '(':
			depth++
			i++
		case ')':
			elems = make([][]byte, 0)
			goto Close
		default:
			break Open
		}
	}
	dims = make([]int, i)

Element:
	for i < len(src) {
		switch src[i] {
		case '(':
			if depth == len(dims) {
				break Element
			}
			depth++
			dims[depth-1] = 0
			i++
		case '"':
			var elem = []byte{}
			var escape bool
			for i++; i < len(src); i++ {
				if escape {
					elem = append(elem, src[i])
					escape = false
				} else {
					switch src[i] {
					default:
						elem = append(elem, src[i])
					case '\\':
						escape = true
					case '"':
						elems = append(elems, elem)
						i++
						break Element
					}
				}
			}
		default:
			for start := i; i < len(src); i++ {
				if bytes.HasPrefix(src[i:], del) || src[i] == ')' {
					elem := src[start:i]
					if len(elem) == 0 {
						return nil, nil, fmt.Errorf("pq: unable to parse array; unexpected %q at offset %d", src[i], i)
					}
					if bytes.Equal(elem, []byte("NULL")) {
						elem = nil
					}
					elems = append(elems, elem)
					break Element
				}
			}
		}
	}

	for i < len(src) {
		if bytes.HasPrefix(src[i:], del) && depth > 0 {
			dims[depth-1]++
			i += len(del)
			goto Element
		} else if src[i] == ')' && depth > 0 {
			dims[depth-1]++
			depth--
			i++
		} else {
			return nil, nil, fmt.Errorf("pq: unable to parse array; unexpected %q at offset %d", src[i], i)
		}
	}

Close:
	for i < len(src) {
		if src[i] == ')' && depth > 0 {
			depth--
			i++
		} else {
			return nil, nil, fmt.Errorf("pq: unable to parse array; unexpected %q at offset %d", src[i], i)
		}
	}
	if depth > 0 {
		err = fmt.Errorf("pq: unable to parse array; expected %q at offset %d", ')', i)
	}
	if err == nil {
		for _, d := range dims {
			if (len(elems) % d) != 0 {
				err = fmt.Errorf("pq: multidimensional arrays must have elements with matching dimensions")
			}
		}
	}
	return
}
