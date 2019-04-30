package textparse

import (
	"fmt"
)

const (
	sInit	= iota
	sComment
	sMeta1
	sMeta2
	sLabels
	sLValue
	sValue
	sTimestamp
)

func (l *promlexer) Lex() token {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if l.i >= len(l.b) {
		return tEOF
	}
	c := l.b[l.i]
	l.start = l.i
yystate0:
	switch yyt := l.state; yyt {
	default:
		panic(fmt.Errorf(`invalid start condition %d`, yyt))
	case 0:
		goto yystart1
	case 1:
		goto yystart8
	case 2:
		goto yystart19
	case 3:
		goto yystart21
	case 4:
		goto yystart24
	case 5:
		goto yystart29
	case 6:
		goto yystart33
	case 7:
		goto yystart36
	}
	goto yystate0
	goto yystate1
yystate1:
	c = l.next()
yystart1:
	switch {
	default:
		goto yyabort
	case c == '#':
		goto yystate5
	case c == ':' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate7
	case c == '\n':
		goto yystate4
	case c == '\t' || c == ' ':
		goto yystate3
	case c == '\x00':
		goto yystate2
	}
yystate2:
	c = l.next()
	goto yyrule1
yystate3:
	c = l.next()
	switch {
	default:
		goto yyrule3
	case c == '\t' || c == ' ':
		goto yystate3
	}
yystate4:
	c = l.next()
	goto yyrule2
yystate5:
	c = l.next()
	switch {
	default:
		goto yyrule5
	case c == '\t' || c == ' ':
		goto yystate6
	}
yystate6:
	c = l.next()
	switch {
	default:
		goto yyrule4
	case c == '\t' || c == ' ':
		goto yystate6
	}
yystate7:
	c = l.next()
	switch {
	default:
		goto yyrule10
	case c >= '0' && c <= ':' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate7
	}
	goto yystate8
yystate8:
	c = l.next()
yystart8:
	switch {
	default:
		goto yyabort
	case c == 'H':
		goto yystate9
	case c == 'T':
		goto yystate14
	case c == '\t' || c == ' ':
		goto yystate3
	}
yystate9:
	c = l.next()
	switch {
	default:
		goto yyabort
	case c == 'E':
		goto yystate10
	}
yystate10:
	c = l.next()
	switch {
	default:
		goto yyabort
	case c == 'L':
		goto yystate11
	}
yystate11:
	c = l.next()
	switch {
	default:
		goto yyabort
	case c == 'P':
		goto yystate12
	}
yystate12:
	c = l.next()
	switch {
	default:
		goto yyabort
	case c == '\t' || c == ' ':
		goto yystate13
	}
yystate13:
	c = l.next()
	switch {
	default:
		goto yyrule6
	case c == '\t' || c == ' ':
		goto yystate13
	}
yystate14:
	c = l.next()
	switch {
	default:
		goto yyabort
	case c == 'Y':
		goto yystate15
	}
yystate15:
	c = l.next()
	switch {
	default:
		goto yyabort
	case c == 'P':
		goto yystate16
	}
yystate16:
	c = l.next()
	switch {
	default:
		goto yyabort
	case c == 'E':
		goto yystate17
	}
yystate17:
	c = l.next()
	switch {
	default:
		goto yyabort
	case c == '\t' || c == ' ':
		goto yystate18
	}
yystate18:
	c = l.next()
	switch {
	default:
		goto yyrule7
	case c == '\t' || c == ' ':
		goto yystate18
	}
	goto yystate19
yystate19:
	c = l.next()
yystart19:
	switch {
	default:
		goto yyabort
	case c == ':' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate20
	case c == '\t' || c == ' ':
		goto yystate3
	}
yystate20:
	c = l.next()
	switch {
	default:
		goto yyrule8
	case c >= '0' && c <= ':' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate20
	}
	goto yystate21
yystate21:
	c = l.next()
yystart21:
	switch {
	default:
		goto yyrule9
	case c == '\t' || c == ' ':
		goto yystate23
	case c >= '\x01' && c <= '\b' || c >= '\v' && c <= '\x1f' || c >= '!' && c <= 'ÿ':
		goto yystate22
	}
yystate22:
	c = l.next()
	switch {
	default:
		goto yyrule9
	case c >= '\x01' && c <= '\t' || c >= '\v' && c <= 'ÿ':
		goto yystate22
	}
yystate23:
	c = l.next()
	switch {
	default:
		goto yyrule3
	case c == '\t' || c == ' ':
		goto yystate23
	case c >= '\x01' && c <= '\b' || c >= '\v' && c <= '\x1f' || c >= '!' && c <= 'ÿ':
		goto yystate22
	}
	goto yystate24
yystate24:
	c = l.next()
yystart24:
	switch {
	default:
		goto yyabort
	case c == ',':
		goto yystate25
	case c == '=':
		goto yystate26
	case c == '\t' || c == ' ':
		goto yystate3
	case c == '}':
		goto yystate28
	case c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate27
	}
yystate25:
	c = l.next()
	goto yyrule15
yystate26:
	c = l.next()
	goto yyrule14
yystate27:
	c = l.next()
	switch {
	default:
		goto yyrule12
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate27
	}
yystate28:
	c = l.next()
	goto yyrule13
	goto yystate29
yystate29:
	c = l.next()
yystart29:
	switch {
	default:
		goto yyabort
	case c == '"':
		goto yystate30
	case c == '\t' || c == ' ':
		goto yystate3
	}
yystate30:
	c = l.next()
	switch {
	default:
		goto yyabort
	case c == '"':
		goto yystate31
	case c == '\\':
		goto yystate32
	case c >= '\x01' && c <= '!' || c >= '#' && c <= '[' || c >= ']' && c <= 'ÿ':
		goto yystate30
	}
yystate31:
	c = l.next()
	goto yyrule16
yystate32:
	c = l.next()
	switch {
	default:
		goto yyabort
	case c >= '\x01' && c <= '\t' || c >= '\v' && c <= 'ÿ':
		goto yystate30
	}
	goto yystate33
yystate33:
	c = l.next()
yystart33:
	switch {
	default:
		goto yyabort
	case c == '\t' || c == ' ':
		goto yystate3
	case c == '{':
		goto yystate35
	case c >= '\x01' && c <= '\b' || c >= '\v' && c <= '\x1f' || c >= '!' && c <= 'z' || c >= '|' && c <= 'ÿ':
		goto yystate34
	}
yystate34:
	c = l.next()
	switch {
	default:
		goto yyrule17
	case c >= '\x01' && c <= '\b' || c >= '\v' && c <= '\x1f' || c >= '!' && c <= 'z' || c >= '|' && c <= 'ÿ':
		goto yystate34
	}
yystate35:
	c = l.next()
	goto yyrule11
	goto yystate36
yystate36:
	c = l.next()
yystart36:
	switch {
	default:
		goto yyabort
	case c == '\n':
		goto yystate37
	case c == '\t' || c == ' ':
		goto yystate3
	case c >= '0' && c <= '9':
		goto yystate38
	}
yystate37:
	c = l.next()
	goto yyrule19
yystate38:
	c = l.next()
	switch {
	default:
		goto yyrule18
	case c >= '0' && c <= '9':
		goto yystate38
	}
yyrule1:
	{
		return tEOF
	}
yyrule2:
	{
		l.state = sInit
		return tLinebreak
		goto yystate0
	}
yyrule3:
	{
		return tWhitespace
	}
yyrule4:
	{
		l.state = sComment
		goto yystate0
	}
yyrule5:
	{
		return l.consumeComment()
	}
yyrule6:
	{
		l.state = sMeta1
		return tHelp
		goto yystate0
	}
yyrule7:
	{
		l.state = sMeta1
		return tType
		goto yystate0
	}
yyrule8:
	{
		l.state = sMeta2
		return tMName
		goto yystate0
	}
yyrule9:
	{
		l.state = sInit
		return tText
		goto yystate0
	}
yyrule10:
	{
		l.state = sValue
		return tMName
		goto yystate0
	}
yyrule11:
	{
		l.state = sLabels
		return tBraceOpen
		goto yystate0
	}
yyrule12:
	{
		return tLName
	}
yyrule13:
	{
		l.state = sValue
		return tBraceClose
		goto yystate0
	}
yyrule14:
	{
		l.state = sLValue
		return tEqual
		goto yystate0
	}
yyrule15:
	{
		return tComma
	}
yyrule16:
	{
		l.state = sLabels
		return tLValue
		goto yystate0
	}
yyrule17:
	{
		l.state = sTimestamp
		return tValue
		goto yystate0
	}
yyrule18:
	{
		return tTimestamp
	}
yyrule19:
	{
		l.state = sInit
		return tLinebreak
		goto yystate0
	}
	panic("unreachable")
	goto yyabort
yyabort:
	if l.state == sComment {
		return l.consumeComment()
	}
	return tInvalid
}
func (l *promlexer) consumeComment() token {
	_logClusterCodePath()
	defer _logClusterCodePath()
	for c := l.cur(); ; c = l.next() {
		switch c {
		case 0:
			return tEOF
		case '\n':
			l.state = sInit
			return tComment
		}
	}
}
