package ecodes

type ECode int // @name ECode

const (
	// Common error
	BadRequest        = 40000
	Unauthorized      = 40001
	PermissionDenied  = 40003
	NotExisted        = 40004
	Existed           = 40005
	QueryParamInvalid = 40006
	InvalidSignature  = 40007
	AddressIsEmpty    = 40008

	// Project error

	// IOT error
	IOTNotAllowed      = 41000
	IOTInvalidNonce    = 41001
	IOTInvalidMintSign = 41002

	// Sensor error
	SensorNotAllowed      = 41100
	SensorInvalidNonce    = 41101
	SensorInvalidMintSign = 41102
	SensorInvalidMetric   = 41103
	SensorInvalidType     = 41104
	SensorHasNoAddress    = 41105
	SensorHasAddress      = 41106
)

const (
	// Runtime group
	Internal            = 50000
	ServiceNotAvailable = 50001

	// Dev group
	NotImplement    = 50020
	NotRegisterAuth = 50021
)
