package protocol

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
)

type RESPValue struct {
	Type   byte
	Str    string
	Num    int64
	Array  []RESPValue
	IsNull bool
}

type RESPError struct {
	Message string
}

const (
	SimpleString = '+'
	Error        = '-'
	Integer      = ':'
	BulkString   = '$'
	Array        = '*'
	CRLF         = "\r\n"
)

func (v *RESPValue) Serialize() string {
	switch v.Type {
	case SimpleString:
		return fmt.Sprintf("+%s\r\n", v.Str)
	case Error:
		return "-Error message\r\n"
	case Integer:
		return fmt.Sprintf(":%d\r\n", v.Num)
	case BulkString:
		if v.IsNull {
			return "$-1\r\n"
		}
		return fmt.Sprintf("$%d\r\n%s\r\n", len(v.Str), v.Str)
	case Array:
		if v.IsNull {
			return "*-1\r\n"
		}
		var builder strings.Builder
		builder.WriteString(fmt.Sprintf("*%d\r\n", len(v.Array)))
		for _, item := range v.Array {
			builder.WriteString(item.Serialize())
		}
		return builder.String()
	default:
		return ""
	}
}

func Deserialize(reader *bufio.Reader) (*RESPValue, error) {
	typ, err := reader.ReadByte()

	if err != nil {
		return nil, err
	}

	switch typ {
	case SimpleString, Error:
		str, err := readLine(reader)
		if err != nil {
			return nil, err
		}
		return &RESPValue{Type: typ, Str: str}, nil
	case Integer:
		line, err := readLine(reader)
		if err != nil {
			return nil, err
		}
		num, err := strconv.ParseInt(line, 10, 64)
		if err != nil {
			return nil, err
		}
		return &RESPValue{Type: typ, Num: num}, nil
	case BulkString:
		line, err := readLine(reader)
		if err != nil {
			return nil, err
		}
		len, err := strconv.Atoi(line)
		if err != nil {
			return nil, err
		}
		if len == -1 {
			return &RESPValue{Type: typ, IsNull: true}, nil
		}
		str, err := readBulkString(reader, len)
		return &RESPValue{Type: typ, Str: str}, nil
	default:
		return nil, errors.New("invalid resp type found")
	}
}

func readLine(reader *bufio.Reader) (string, error) {
	line, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	return strings.TrimRight(line, CRLF), nil
}

func readBulkString(reader *bufio.Reader, length int) (string, error) {
	buf := make([]byte, length+2)
	_, err := io.ReadFull(reader, buf)
	if err != nil {
		return "", err
	}
	return string(buf[:length]), nil
}
