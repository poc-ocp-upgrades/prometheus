package graphite

import (
	"bytes"
	"fmt"
	"strings"
	"github.com/prometheus/common/model"
)

const (
	symbols		= "(){},=.'\"\\"
	printables	= ("0123456789abcdefghijklmnopqrstuvwxyz" + "ABCDEFGHIJKLMNOPQRSTUVWXYZ" + "!\"#$%&\\'()*+,-./:;<=>?@[\\]^_`{|}~")
)

func escape(tv model.LabelValue) string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	length := len(tv)
	result := bytes.NewBuffer(make([]byte, 0, length))
	for i := 0; i < length; i++ {
		b := tv[i]
		switch {
		case b == '.' || b == '%' || b == '/' || b == '=':
			fmt.Fprintf(result, "%%%X", b)
		case strings.IndexByte(symbols, b) != -1:
			result.WriteString("\\" + string(b))
		case strings.IndexByte(printables, b) != -1:
			result.WriteByte(b)
		default:
			fmt.Fprintf(result, "%%%X", b)
		}
	}
	return result.String()
}
