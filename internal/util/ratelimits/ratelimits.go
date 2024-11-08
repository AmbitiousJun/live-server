package ratelimits

// Bucket 令牌桶接口
type Bucket interface {

	// Consume 消耗一定数量的令牌
	// 如果桶中没有足量的令牌, 该方法会阻塞线程直至成功消耗相应数量的令牌
	Consume(token uint)

	// TryConsume 非阻塞消耗一定数量的令牌
	// 返回是否消耗成功
	TryConsume(token uint) bool

	// Destroy 销毁令牌桶
	Destroy()
}
