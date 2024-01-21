package msgs

const (
	// main
	ServerStarting    = "Server Starting"
	ConfigInitialized = "Config Initialized"

	// application
	ApplicationStart                = "Start Application"
	ApplicationStarted              = "Application start is over, waiting when ctx done"
	TelegramAdapterInitialized      = "TelegramAdapterInitialized"
	VkontakteAdapterInitialized     = "VkontakteAdapterInitialized"
	DatabaseInitialized             = "DatabaseInitialized"
	MessagesMapCreated              = "MessagesMapCreated"
	InitGrpcMergerClientInitialized = "InitGrpcMergerClientInitialized"
	ApplicationReceiveCtxDone       = "Application receive ctx.Done signal"
	ApplicationReceiveInternalError = "ApplicationReceiveInternalError"

	//  Runnable
	RunRunnable     = "Run Runnable"
	StoppedRunnable = "Stopped Runnable"

	// db
	RunGracefulShutdownDb        = "RunGracefulShutdownDb"
	StoppedDatabaseWithoutErrors = "StoppedDatabaseWithoutErrors"
	StoppedDatabaseError         = "StoppedDatabaseError"
)
