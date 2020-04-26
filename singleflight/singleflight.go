package singleflight

import "sync"

// 代表正在进行中或者已经结束的请求
type call struct {
	wg  sync.WaitGroup
	val interface{}
	err error
}

// 管理不同key的请求
type Group struct {
	mu sync.Mutex
	m  map[string]*call
}

// 核心函数，保证一个key只调用一次fn
func (g *Group) Do(key string, fn func() (interface{}, error)) (interface{}, error) {
	g.mu.Lock() // 保证只有一个协程调用fn
	if g.m == nil {
		g.m = make(map[string]*call)
	}
	// 判断有没有进行中的请求
	if c, ok := g.m[key]; ok {
		g.mu.Unlock()
		c.wg.Wait()         // 等待正在进行中的请求结束
		return c.val, c.err // 返回结果
	}
	c := new(call)
	c.wg.Add(1) // 请求前加锁
	g.m[key] = c
	g.mu.Unlock() // 获取调用fn资格后释放调用锁

	c.val, c.err = fn() // 调用fn，获取返回值
	c.wg.Done()         // 释放请求锁

	g.mu.Lock()
	delete(g.m, key) // 更新 g.m
	g.mu.Unlock()

	return c.val, c.err // 返回结果
}
