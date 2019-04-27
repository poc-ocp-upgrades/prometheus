package runtime

func charsToString(ca []int8) string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	s := make([]byte, 0, len(ca))
	for _, c := range ca {
		if byte(c) == 0 {
			break
		}
		s = append(s, byte(c))
	}
	return string(s)
}
