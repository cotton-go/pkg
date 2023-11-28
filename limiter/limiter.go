package limiter

import (
	"sync/atomic"
)

const (
	DefaultLimit = 100
)

// Limiter 是一个限制器类，用于控制并发执行的任务数量。
type Limiter struct {
	limit         int
	numInProgress int32
	tickets       chan int32
}

// NewLimiter 使用给定的限制创建一个新的 Limiter 实例。
//
// 如果限制小于等于0，则使用默认限制。
//
// 函数返回新创建的 Limiter 实例的指针。
func NewLimiter(limit int) *Limiter {
	// 如果提供的限制小于或等于0，则设置默认限制
	if limit <= 0 {
		limit = DefaultLimit
	}

	// 使用提供的限制创建一个新的 Limiter 实例
	c := &Limiter{
		limit:   limit,
		tickets: make(chan int32, limit),
	}

	// 用初始票数填充票证通道
	for i := 0; i < c.limit; i++ {
		c.tickets <- int32(i) + 1
	}

	// 返回新创建的 Limiter 实例
	return c
}

// Execute 方法在从限制器获取票据后运行 job 函数。
// job 函数完成后释放票据。
// job 函数是一个不带参数也不返回值的函数。
func (c *Limiter) Execute(job func()) {
	// 从限制器获取票据
	ticket := c.acquire()

	// job 函数完成后，释放票据
	defer c.release(ticket)

	// 执行 job 函数
	job()
}

// ExecuteWithTicket 使用从限制器获取的票据执行给定的作业函数。
// 获取的票据在作业函数执行后被释放。
// 该函数接受一个作业函数作为输入，该函数接受一个表示获取的票据的整数参数。
func (c *Limiter) ExecuteWithTicket(job func(ticket int)) {
	// 从限制器获取一个票据
	ticket := c.acquire()

	// 确保无论作业函数是否发生恐慌，都释放票据
	defer c.release(ticket)

	// 使用获取的票据执行作业函数
	job(int(ticket))
}

// Wait 阻塞直到限制器有可用的票据。
func (c *Limiter) Wait() {
	// 遍历限制数并从通道中消耗票据。
	for i := 0; i < c.limit; i++ {
		<-c.tickets
	}
}

// GetNumInProgress 返回当前正在进行中的任务数量。
// 它使用原子操作来加载 numInProgress 的值。
func (c *Limiter) GetNumInProgress() int32 {
	return atomic.LoadInt32(&c.numInProgress)
}

// acquire 是 Limiter 类的一个方法。它用于获取一张票，
// 原子性地增加正在进行的任务数量，并且如果没有可用的票则等待。
// 此函数通过使用原子操作来增加 numInProgress 以确保线程安全。
// 它返回获取的票。
func (c *Limiter) acquire() int32 {
	// 原子性地增加正在进行的任务数量
	atomic.AddInt32(&c.numInProgress, 1)

	// 从 tickets 通道获取一张票
	return <-c.tickets
}

// release 释放一个票据。
// 它减少正在进行的任务数量并向票据通道发送一个空结构体。
func (c *Limiter) release(ticket int32) {
	atomic.AddInt32(&c.numInProgress, -1)
	c.tickets <- ticket
}
