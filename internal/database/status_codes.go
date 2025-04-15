package database

const (
	// Class 23 — Integrity Constraint Violation
	UniqueViolation     = "23505"
	ForeignKeyViolation = "23503"
	NotNullViolation    = "23502"
	CheckViolation      = "23514"
	ExclusionViolation  = "23P01"

	// Class 42 — Syntax Error or Access Rule Violation
	SyntaxError     = "42601"
	UndefinedTable  = "42P01"
	UndefinedColumn = "42703"

	// Class 08 — Connection Exception
	ConnectionException = "08000"
	ConnectionFailure   = "08006"

	// Class 53 — Insufficient Resources
	InsufficientResources = "53000"
	DiskFull              = "53100"
	OutOfMemory           = "53200"
	TooManyConnections    = "53300"
)
