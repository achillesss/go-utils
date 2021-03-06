package rate

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"sort"
	"sync"
	"sync/atomic"
	"time"
)

type durationCount map[time.Duration]int64

type RateLimitType int

var rateLimitTypeName = make(map[RateLimitType]string)

func registerRateLimitTypeName(typ int, name string) {
	rateLimitTypeName[RateLimitType(typ)] = name
}

func (t RateLimitType) String() string {
	return rateLimitTypeName[t]
}

type rate struct {
	l sync.RWMutex
	m durationCount
}

func (r *rate) countDuration(duration time.Duration, cnt int64) {
	r.l.Lock()
	var count = r.m[duration]
	r.m[duration] = atomic.AddInt64(&count, cnt)
	r.l.Unlock()
}

func (r *rate) getCount(duration time.Duration) int64 {
	r.l.RLock()
	var count = r.m[duration]
	r.l.RUnlock()
	return count
}

func newRate(durations ...time.Duration) *rate {
	var r rate
	r.m = make(durationCount)
	for _, duration := range durations {
		r.m[duration] = 0
	}
	return &r
}

type tries struct {
	l   sync.RWMutex
	cnt map[string]durationCount
}

func (t *tries) getCount(id string, duration time.Duration) int64 {
	t.l.RLock()
	var count = t.cnt[id][duration]
	t.l.RUnlock()
	return count
}

func (t *tries) count(id string, duration time.Duration) {
	t.l.Lock()
	var dc, ok = t.cnt[id]
	if !ok {
		dc = make(durationCount)
		t.cnt[id] = dc
	}

	count, ok := dc[duration]
	if !ok {
		go func(t *tries) {
			var tm = time.NewTimer(duration)
			<-tm.C
			t.l.Lock()
			delete(dc, duration)
			t.l.Unlock()
		}(t)
	}

	dc[duration] = atomic.AddInt64(&count, 1)
	t.l.Unlock()
}

type CountInfo struct {
	ID       string
	Name     string
	Max      int64
	Tries    int64
	Duration time.Duration
	Typ      RateLimitType
}

func (f *CountInfo) String() string {
	return fmt.Sprintf("ID[%s]METHOD[%s]TYP[%s]CNT[%d/%d]IN[%s]", f.ID, f.Name, f.Typ, f.Tries, f.Max, f.Duration)
}

type rule struct {
	name      string
	typ       RateLimitType
	max       rate
	tries     tries
	durations []time.Duration
}

func (r *rule) count(id string) (*CountInfo, bool) {
	for _, duration := range r.durations {
		var max = r.max.getCount(duration)
		var try = r.tries.getCount(id, duration)
		if try >= max {
			return &CountInfo{
				ID:       id,
				Name:     r.name,
				Max:      max,
				Tries:    try,
				Duration: duration,
				Typ:      r.typ,
			}, false
		}
	}

	for _, duration := range r.durations {
		r.tries.count(id, duration)
	}

	return nil
}

type Rules struct {
	l     sync.RWMutex
	rules map[string]*rule
}

func NewRules() *Rules {
	var r Rules
	r.rules = make(map[string]*rule)
	return &r
}

func nameTypeID(name string, typ int) string {
	var h = md5.New()
	h.Write([]byte(fmt.Sprintf("%s*|*%d", name, typ)))
	return hex.EncodeToString(h.Sum(nil))
}

func (r *Rules) AddRule(name string, duration time.Duration, max int64, typ int) {
	var ruleID = nameTypeID(name, typ)
	r.l.RLock()
	var rl, ok = r.rules[ruleID]
	r.l.RUnlock()

	if !ok {
		rl = new(rule)
		rl.name = name
		rl.typ = RateLimitType(typ)
		rl.max.m = make(durationCount)
		rl.tries.cnt = make(map[string]durationCount)
	}

	r.l.Lock()
	r.rules[ruleID] = rl
	rl.max.m[duration] = max
	rl.durations = nil
	for k := range rl.max.m {
		rl.durations = append(rl.durations, k)
	}
	sort.Slice(rl.durations, func(i, j int) bool {
		return rl.durations[i] < rl.durations[j]
	})
	r.l.Unlock()
}

func (r *Rules) Call(id string, name string, typ int) (*CountInfo, bool) {
	var ruleID = nameTypeID(name, typ)
	r.l.RLock()
	var rl, ok = r.rules[ruleID]
	r.l.RUnlock()
	if !ok {
		return nil
	}

	return rl.count(id)
}

func (r *Rules) SetRateLimitTypeName(typ int, name string) {
	registerRateLimitTypeName(typ, name)
}
