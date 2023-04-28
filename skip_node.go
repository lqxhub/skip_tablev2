package skiptablev2

type SkipListLevel[K comparable, V SkipListItem[K]] struct {
	//指向下一个结点
	forward *SkipListNode[K, V]

	/*
	 * 到下一个node的距离;
	 * 思考,为啥是记录到下一个node, 而不是记录上一个node到这的距离
	 */
	span int64
}

type SkipListNode[K comparable, V SkipListItem[K]] struct {
	//指向上一个结点
	backward *SkipListNode[K, V]
	//索引用的层
	level []SkipListLevel[K, V]
	//存储的值
	value V
	//排名用的分数
	score float64
}

func NewSkipListNode[K comparable, V SkipListItem[K]](level int, score float64, value V) *SkipListNode[K, V] {
	return &SkipListNode[K, V]{
		backward: nil,
		level:    make([]SkipListLevel[K, V], level),
		value:    value,
		score:    score,
	}
}

// Next 第i层的下一个元素
func (node *SkipListNode[K, V]) Next(i int) *SkipListNode[K, V] {
	return node.level[i].forward
}

// SetNext 设置第i层的下一个元素
func (node *SkipListNode[K, V]) SetNext(i int, next *SkipListNode[K, V]) {
	node.level[i].forward = next
}

// Span 第i层的span值
func (node *SkipListNode[K, V]) Span(i int) int64 {
	return node.level[i].span
}

// SetSpan 设置第i层的span值
func (node *SkipListNode[K, V]) SetSpan(i int, span int64) {
	node.level[i].span = span
}

// Pre 上一个元素    想一下,为啥指向上一个的元素不需要i呢???
func (node *SkipListNode[K, V]) Pre() *SkipListNode[K, V] {
	return node.backward
}
