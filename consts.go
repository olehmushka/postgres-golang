package postgres

const (
	DefaultBatchItemsMaxSize = 500
	DefaultSchemaName        = "public"

	TimeLayoutISO = "2006-01-02 15:04:05"
)

type TxIsoLevel string

// Transaction isolation levels
const (
	Serializable    = TxIsoLevel("serializable")
	RepeatableRead  = TxIsoLevel("repeatable read")
	ReadCommitted   = TxIsoLevel("read committed")
	ReadUncommitted = TxIsoLevel("read uncommitted")
)

type TxAccessMode string

// Transaction access modes
const (
	ReadWrite = TxAccessMode("read write")
	ReadOnly  = TxAccessMode("read only")
)

type TxDeferrableMode string

// Transaction deferrable modes
const (
	Deferrable    = TxDeferrableMode("deferrable")
	NotDeferrable = TxDeferrableMode("not deferrable")
)

type TxOptions struct {
	IsoLevel       TxIsoLevel
	AccessMode     TxAccessMode
	DeferrableMode TxDeferrableMode
}
