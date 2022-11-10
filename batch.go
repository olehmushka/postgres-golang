package postgres

type Batch interface {
	SetItem(BatchItem)
	GetItems() []BatchItem
	Len() int
}
