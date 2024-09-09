package enum

const (
	ErrorConfigFile   = "fatal error config file: %w"
	ErrorConnectDB    = "Unable to connect to database!"
	ErrorLoadEnvFile  = "Load env file failed!"
	ErrorCreateSchema = "Failed to create schema!"

	ErrorReadRequestBody   = "Failed to read request body!"
	ApplicationStartFailed = "application start failed!"
)
