package postgres

type batch struct {
	items []BatchItem
}

func NewBatch() Batch {
	return &batch{}
}

func (b *batch) SetItem(item BatchItem) {
	b.items = append(b.items, item)
}

func (b *batch) GetItems() []BatchItem {
	return b.items
}

func (b *batch) Len() int {
	return len(b.items)
}
