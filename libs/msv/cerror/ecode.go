package cerror

type ECode int

const (
	// Common error
	ECodeBadRequest        ECode = 40000
	ECodeUnauthorized      ECode = 40001
	ECodePermissionDenied  ECode = 40003
	ECodeNotExisted        ECode = 40004
	ECodeExisted           ECode = 40005
	ECodeQueryParamInvalid ECode = 40006
	ECodeInvalidSignature  ECode = 40007
	ECodeAddressIsEmpty    ECode = 40008

	// Project error

	// IOT error
	ECodeIOTNotAllowed      ECode = 41000
	ECodeIOTInvalidNonce    ECode = 41001
	ECodeIOTInvalidMintSign ECode = 41002

	// Sensor error
	ECodeSensorNotAllowed      ECode = 41100
	ECodeSensorInvalidNonce    ECode = 41101
	ECodeSensorInvalidMintSign ECode = 41102
	ECodeSensorInvalidMetric   ECode = 41103
	ECodeSensorInvalidType     ECode = 41104
	ECodeSensorHasNoAddress    ECode = 41105
	ECodeSensorHasAddress      ECode = 41106
)

const (
	ECodeInternal     = 50000
	ECodeNotImplement = 50001
)
