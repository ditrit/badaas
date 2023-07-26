package defaults

// Database
const (
	DatabaseRetryTimes             = uint(10)
	DatabaseRetryDuration          = uint(5)
	ServerTimeout                  = 15
	ServerHost                     = "0.0.0.0"
	ServerPort                     = 8000
	ServerPaginationMaxElemPerPage = uint(100)
	SessionDuration                = uint(3600 * 4) // 4 hours
	SessionPullInterval            = uint(30)       // 30 seconds
	SessionRollInterval            = uint(3600)     // 1 hour
	LoggerSlowQueryThreshold       = 200
	LoggerSlowTransactionThreshold = 200
)
