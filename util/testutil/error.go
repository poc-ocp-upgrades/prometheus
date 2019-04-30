package testutil

func ErrorEqual(left, right error) bool {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if left == right {
		return true
	}
	if left != nil && right != nil {
		return left.Error() == right.Error()
	}
	return false
}
