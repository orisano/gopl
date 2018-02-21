package ftpcodes

// rfc640
// https://tools.ietf.org/rfc/rfc640.txt

const (
	PositivePreliminary         = 100
	PositiveCompletion          = 200
	PositiveIntermediate        = 300
	TransientNegativeCompletion = 400
	PermanentNegativeCompletion = 500
)

const (
	Syntax         = 00
	Information    = 10
	Connections    = 20
	Authentication = 30
	Unspecified    = 40
	FileSystem     = 50
)

const (
	CommandOkay                       = PositiveCompletion + Syntax + 0
	CommandSyntaxError                = PermanentNegativeCompletion + Syntax + 0
	ParamsSyntaxError                 = PermanentNegativeCompletion + Syntax + 1
	CommandNotImplementedSuperfluous  = PositiveCompletion + Syntax + 2
	CommandNotImplemented             = PermanentNegativeCompletion + Syntax + 2
	BadSequenceCommand                = PermanentNegativeCompletion + Syntax + 3
	CommandNotImplementedForParameter = PermanentNegativeCompletion + Syntax + 4
	RestartMarkerReply                = PositivePreliminary + Information + 0
	SystemStatus                      = PositiveCompletion + Information + 1
	DirectoryStatus                   = PositiveCompletion + Information + 2
	FileStatus                        = PositiveCompletion + Information + 3
	HelpMessage                       = PositiveCompletion + Information + 4
	SystemType                        = PositiveCompletion + Information + 5
	ServiceReadyInNMinutes            = PositivePreliminary + Connections + 0
	ServiceReadyForNewUser            = PositiveCompletion + Connections + 0
	ServiceClosingTELNETConnection    = PositiveCompletion + Connections + 1
	ServiceNotAvailable               = TransientNegativeCompletion + Connections + 1
	DataConnectionAlreadyOpen         = PositivePreliminary + Connections + 5
	DataConnectionOpen                = PositiveCompletion + Connections + 5
	CantOpenDataConnection            = TransientNegativeCompletion + Connections + 5
	ClosingDataConnection             = PositiveCompletion + Connections + 6
	ConnectionTrouble                 = TransientNegativeCompletion + Connections + 6
	Entering                          = PositiveCompletion + Connections + 7
	ExtenedPassiveModeEntered         = PositiveCompletion + Connections + 9
	UserLoggedOn                      = PositiveCompletion + Authentication + 0
	NotLoggedIn                       = PermanentNegativeCompletion + Authentication + 0
	UserNameOkay                      = PositiveIntermediate + Authentication + 1
	NeedAccountForLogin               = PositiveIntermediate + Authentication + 2
	NeedAccountForStoringFiles        = PermanentNegativeCompletion + Authentication + 2
	FileStatusOkay                    = PositivePreliminary + FileSystem + 0
	RequestedFileActionOkey           = PositiveCompletion + FileSystem + 0
	RequestedFileActionPending        = PositiveIntermediate + FileSystem + 0
	FileUnavailable                   = TransientNegativeCompletion + FileSystem + 0
	FileUnavailableBusy               = PermanentNegativeCompletion + FileSystem + 0
	LocalErrorInProcessing            = TransientNegativeCompletion + FileSystem + 1
	InsufficientStorageSpace          = TransientNegativeCompletion + FileSystem + 2
	ExceededStorageAllocation         = PermanentNegativeCompletion + FileSystem + 2
	FileNameNotAllowed                = PermanentNegativeCompletion + FileSystem + 3
	StartMailInput                    = PositiveIntermediate + FileSystem + 4
)
