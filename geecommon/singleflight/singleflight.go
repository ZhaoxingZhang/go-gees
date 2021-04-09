package singleflight

import "sync"

/*
call 代表正在进行中，或已经结束的请求。使用 sync.WaitGroup 锁避免重入。
Group 是 singleflight 的主数据结构，管理不同 key 的请求(call)。
*/
type call struct {
	wg  sync.WaitGroup
	val interface{}
	err error
}

type Group struct {
	mu sync.Mutex       // protects m
	m  map[string]*call
}

func (g *Group) Do(key string, fn func() (interface{}, error)) (interface{}, error) {
	g.mu.Lock()
	// 延迟初始化, 目的很简单，提高内存使用效率
	if g.m == nil {
		g.m = make(map[string]*call)
	}
	// 之前有相同请求，等候完成
	if c, ok := g.m[key]; ok {
		g.mu.Unlock()
		c.wg.Wait()
		return c.val, c.err
	}

	c := new(call)
	c.wg.Add(1)
	g.m[key] = c
	g.mu.Unlock()

	c.val, c.err = fn()
	c.wg.Done()

	g.mu.Lock()
	delete(g.m, key)
	g.mu.Unlock()

	return c.val, c.err
}