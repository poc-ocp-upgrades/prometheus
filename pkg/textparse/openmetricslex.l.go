package textparse

import (
	"fmt"
)

func (l *openMetricsLexer) Lex() token {
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
		goto yystart5
	case 2:
		goto yystart25
	case 3:
		goto yystart27
	case 4:
		goto yystart30
	case 5:
		goto yystart35
	case 6:
		goto yystart39
	case 7:
		goto yystart43
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
		goto yystate2
	case c == ':' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate4
	}
yystate2:
	c = l.next()
	switch {
	default:
		goto yyabort
	case c == ' ':
		goto yystate3
	}
yystate3:
	c = l.next()
	goto yyrule1
yystate4:
	c = l.next()
	switch {
	default:
		goto yyrule8
	case c >= '0' && c <= ':' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate4
	}
	goto yystate5
yystate5:
	c = l.next()
yystart5:
	switch {
	default:
		goto yyabort
	case c == 'E':
		goto yystate6
	case c == 'H':
		goto yystate10
	case c == 'T':
		goto yystate15
	case c == 'U':
		goto yystate20
	}
yystate6:
	c = l.next()
	switch {
	default:
		goto yyabort
	case c == 'O':
		goto yystate7
	}
yystate7:
	c = l.next()
	switch {
	default:
		goto yyabort
	case c == 'F':
		goto yystate8
	}
yystate8:
	c = l.next()
	switch {
	default:
		goto yyrule5
	case c == '\n':
		goto yystate9
	}
yystate9:
	c = l.next()
	goto yyrule5
yystate10:
	c = l.next()
	switch {
	default:
		goto yyabort
	case c == 'E':
		goto yystate11
	}
yystate11:
	c = l.next()
	switch {
	default:
		goto yyabort
	case c == 'L':
		goto yystate12
	}
yystate12:
	c = l.next()
	switch {
	default:
		goto yyabort
	case c == 'P':
		goto yystate13
	}
yystate13:
	c = l.next()
	switch {
	default:
		goto yyabort
	case c == ' ':
		goto yystate14
	}
yystate14:
	c = l.next()
	goto yyrule2
yystate15:
	c = l.next()
	switch {
	default:
		goto yyabort
	case c == 'Y':
		goto yystate16
	}
yystate16:
	c = l.next()
	switch {
	default:
		goto yyabort
	case c == 'P':
		goto yystate17
	}
yystate17:
	c = l.next()
	switch {
	default:
		goto yyabort
	case c == 'E':
		goto yystate18
	}
yystate18:
	c = l.next()
	switch {
	default:
		goto yyabort
	case c == ' ':
		goto yystate19
	}
yystate19:
	c = l.next()
	goto yyrule3
yystate20:
	c = l.next()
	switch {
	default:
		goto yyabort
	case c == 'N':
		goto yystate21
	}
yystate21:
	c = l.next()
	switch {
	default:
		goto yyabort
	case c == 'I':
		goto yystate22
	}
yystate22:
	c = l.next()
	switch {
	default:
		goto yyabort
	case c == 'T':
		goto yystate23
	}
yystate23:
	c = l.next()
	switch {
	default:
		goto yyabort
	case c == ' ':
		goto yystate24
	}
yystate24:
	c = l.next()
	goto yyrule4
	goto yystate25
yystate25:
	c = l.next()
yystart25:
	switch {
	default:
		goto yyabort
	case c == ':' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate26
	}
yystate26:
	c = l.next()
	switch {
	default:
		goto yyrule6
	case c >= '0' && c <= ':' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate26
	}
	goto yystate27
yystate27:
	c = l.next()
yystart27:
	switch {
	default:
		goto yyabort
	case c == ' ':
		goto yystate28
	}
yystate28:
	c = l.next()
	switch {
	default:
		goto yyabort
	case c == '\n':
		goto yystate29
	case c >= '\x01' && c <= '\t' || c >= '\v' && c <= 'ÿ':
		goto yystate28
	}
yystate29:
	c = l.next()
	goto yyrule7
	goto yystate30
yystate30:
	c = l.next()
yystart30:
	switch {
	default:
		goto yyabort
	case c == ',':
		goto yystate31
	case c == '=':
		goto yystate32
	case c == '}':
		goto yystate34
	case c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate33
	}
yystate31:
	c = l.next()
	goto yyrule13
yystate32:
	c = l.next()
	goto yyrule12
yystate33:
	c = l.next()
	switch {
	default:
		goto yyrule10
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate33
	}
yystate34:
	c = l.next()
	goto yyrule11
	goto yystate35
yystate35:
	c = l.next()
yystart35:
	switch {
	default:
		goto yyabort
	case c == '"':
		goto yystate36
	}
yystate36:
	c = l.next()
	switch {
	default:
		goto yyabort
	case c == '"':
		goto yystate37
	case c == '\\':
		goto yystate38
	case c >= '\x01' && c <= '\t' || c >= '\v' && c <= '!' || c >= '#' && c <= '[' || c >= ']' && c <= 'ÿ':
		goto yystate36
	}
yystate37:
	c = l.next()
	goto yyrule14
yystate38:
	c = l.next()
	switch {
	default:
		goto yyabort
	case c >= '\x01' && c <= '\t' || c >= '\v' && c <= 'ÿ':
		goto yystate36
	}
	goto yystate39
yystate39:
	c = l.next()
yystart39:
	switch {
	default:
		goto yyabort
	case c == ' ':
		goto yystate40
	case c == '{':
		goto yystate42
	}
yystate40:
	c = l.next()
	switch {
	default:
		goto yyabort
	case c >= '\x01' && c <= '\t' || c >= '\v' && c <= '\x1f' || c >= '!' && c <= 'ÿ':
		goto yystate41
	}
yystate41:
	c = l.next()
	switch {
	default:
		goto yyrule15
	case c >= '\x01' && c <= '\t' || c >= '\v' && c <= '\x1f' || c >= '!' && c <= 'ÿ':
		goto yystate41
	}
yystate42:
	c = l.next()
	goto yyrule9
	goto yystate43
yystate43:
	c = l.next()
yystart43:
	switch {
	default:
		goto yyabort
	case c == ' ':
		goto yystate45
	case c == '\n':
		goto yystate44
	}
yystate44:
	c = l.next()
	goto yyrule18
yystate45:
	c = l.next()
	switch {
	default:
		goto yyabort
	case c == '#':
		goto yystate47
	case c >= '\x01' && c <= '\t' || c >= '\v' && c <= '\x1f' || c == '!' || c == '"' || c >= '$' && c <= 'ÿ':
		goto yystate46
	}
yystate46:
	c = l.next()
	switch {
	default:
		goto yyrule16
	case c >= '\x01' && c <= '\t' || c >= '\v' && c <= '\x1f' || c >= '!' && c <= 'ÿ':
		goto yystate46
	}
yystate47:
	c = l.next()
	switch {
	default:
		goto yyrule16
	case c == ' ':
		goto yystate48
	case c >= '\x01' && c <= '\t' || c >= '\v' && c <= '\x1f' || c >= '!' && c <= 'ÿ':
		goto yystate46
	}
yystate48:
	c = l.next()
	switch {
	default:
		goto yyabort
	case c == '\n':
		goto yystate49
	case c >= '\x01' && c <= '\t' || c >= '\v' && c <= 'ÿ':
		goto yystate48
	}
yystate49:
	c = l.next()
	goto yyrule17
yyrule1:
	{
		l.state = sComment
		goto yystate0
	}
yyrule2:
	{
		l.state = sMeta1
		return tHelp
		goto yystate0
	}
yyrule3:
	{
		l.state = sMeta1
		return tType
		goto yystate0
	}
yyrule4:
	{
		l.state = sMeta1
		return tUnit
		goto yystate0
	}
yyrule5:
	{
		l.state = sInit
		return tEofWord
		goto yystate0
	}
yyrule6:
	{
		l.state = sMeta2
		return tMName
		goto yystate0
	}
yyrule7:
	{
		l.state = sInit
		return tText
		goto yystate0
	}
yyrule8:
	{
		l.state = sValue
		return tMName
		goto yystate0
	}
yyrule9:
	{
		l.state = sLabels
		return tBraceOpen
		goto yystate0
	}
yyrule10:
	{
		return tLName
	}
yyrule11:
	{
		l.state = sValue
		return tBraceClose
		goto yystate0
	}
yyrule12:
	{
		l.state = sLValue
		return tEqual
		goto yystate0
	}
yyrule13:
	{
		return tComma
	}
yyrule14:
	{
		l.state = sLabels
		return tLValue
		goto yystate0
	}
yyrule15:
	{
		l.state = sTimestamp
		return tValue
		goto yystate0
	}
yyrule16:
	{
		return tTimestamp
	}
yyrule17:
	{
		l.state = sInit
		return tLinebreak
		goto yystate0
	}
yyrule18:
	{
		l.state = sInit
		return tLinebreak
		goto yystate0
	}
	panic("unreachable")
	goto yyabort
yyabort:
	return tInvalid
}
