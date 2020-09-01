package stack

import (
	"encoding/hex"
	"hash/maphash"
	"regexp"
	"runtime/debug"
	"strings"
	"sync"
	"sync/atomic"
)

type Stacker interface {
	StoreStack([]byte) string
	QueryStack(string) (int64, []byte)
	ListStacks() map[string]int64
}

var internalStacker Stacker

func SetStacker(stacker Stacker) {
	internalStacker = stacker
}

type defaultStacker struct {
	stacks     map[string][]byte
	stacksLock sync.RWMutex

	counts     map[string]int64
	countsLock sync.RWMutex

	seed maphash.Seed
}

func init() {
	var s defaultStacker
	s.stacks = make(map[string][]byte)
	s.counts = make(map[string]int64)
	s.seed = maphash.MakeSeed()
	internalStacker = &s

}

var stackHashPrefix string

func SetDefaultStackKeyPrefix(pre string) {
	stackHashPrefix = pre
}

func (s *defaultStacker) StoreStack(stack []byte) string {
	//	var h = md5.New()
	//	h.Write(stack)
	//	var key = stackHashPrefix + hex.EncodeToString(h.Sum(nil))

	var h = maphash.Hash{}
	h.SetSeed(s.seed)
	h.Write(stack)
	var key = stackHashPrefix + hex.EncodeToString(h.Sum(nil))

	s.stacksLock.Lock()
	s.stacks[key] = stack
	s.stacksLock.Unlock()

	s.countsLock.Lock()
	var count = s.counts[key]
	atomic.AddInt64(&count, 1)
	s.counts[key] = count
	s.countsLock.Unlock()

	return key
}

func (s *defaultStacker) QueryStack(key string) (int64, []byte) {
	s.stacksLock.RLock()
	var st, ok = s.stacks[key]
	s.stacksLock.RUnlock()

	if ok {
		s.countsLock.RLock()
		var count = s.counts[key]
		s.countsLock.RUnlock()
		return count, st
	}

	return int64(len(s.stacks)), nil
}

func (s *defaultStacker) ListStacks() map[string]int64 {
	var count = make(map[string]int64)
	s.countsLock.RLock()
	for k, v := range s.counts {
		count[k] = v
	}
	s.countsLock.RUnlock()
	return count
}

var reg = regexp.MustCompile(`\((0x[^)]+)\)`)

func getStackHash(stack []byte) string {
	stack = reg.ReplaceAll(stack, []byte("()"))
	var stacks = strings.Split(string(stack), "\n")
	var newStack = strings.Join(stacks[1:], "\n")
	return internalStacker.StoreStack([]byte(newStack))
}

func GetStackHash() string {
	var stack = debug.Stack()
	return getStackHash(stack)
}

func QueryStack(key string) (int64, []byte) {
	return internalStacker.QueryStack(key)
}

func ListStacksCount() map[string]int64 {
	return internalStacker.ListStacks()
}
