package postgres

type batchItem struct {
	query     string
	arguments []any
}

func NewBatchItem() BatchItem {
	return &batchItem{}
}

func (item *batchItem) SetQuery(query string) {
	item.query = query
}

func (item *batchItem) SetArgs(args ...any) {
	item.arguments = args
}

func (item *batchItem) GetQuery() string {
	return item.query
}

func (item *batchItem) GetArgs() []any {
	return item.arguments
}
