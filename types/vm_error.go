package types

type VMStatusCode uint64

const (
	// The status of a transaction as determined by the prologue.
	// Validation Errors: 0-999
	// We don't want the default value to be valid
	UNKNOWN_VALIDATION_STATUS VMStatusCode = 0
	// The transaction has a bad signature
	INVALID_SIGNATURE VMStatusCode = 1
	// Bad account authentication key
	INVALID_AUTH_KEY VMStatusCode = 2
	// Sequence number is too old
	SEQUENCE_NUMBER_TOO_OLD VMStatusCode = 3
	// Sequence number is too new
	SEQUENCE_NUMBER_TOO_NEW VMStatusCode = 4
	// Insufficient balance to pay minimum transaction fee
	INSUFFICIENT_BALANCE_FOR_TRANSACTION_FEE VMStatusCode = 5
	// The transaction has expired
	TRANSACTION_EXPIRED VMStatusCode = 6
	// The sending account does not exist
	SENDING_ACCOUNT_DOES_NOT_EXIST VMStatusCode = 7
	// This write set transaction was rejected because it did not meet the
	// requirements for one.
	REJECTED_WRITE_SET VMStatusCode = 8
	// This write set transaction cannot be applied to the current state.
	INVALID_WRITE_SET VMStatusCode = 9
	// Length of program field in raw transaction exceeded max length
	EXCEEDED_MAX_TRANSACTION_SIZE VMStatusCode = 10
	// This script is not on our whitelist of script.
	UNKNOWN_SCRIPT VMStatusCode = 11
	// Transaction is trying to publish a new module.
	UNKNOWN_MODULE VMStatusCode = 12
	// Max gas units submitted with transaction exceeds max gas units bound
	// in VM
	MAX_GAS_UNITS_EXCEEDS_MAX_GAS_UNITS_BOUND VMStatusCode = 13
	// Max gas units submitted with transaction not enough to cover the
	// intrinsic cost of the transaction.
	MAX_GAS_UNITS_BELOW_MIN_TRANSACTION_GAS_UNITS VMStatusCode = 14
	// Gas unit price submitted with transaction is below minimum gas price
	// set in the VM.
	GAS_UNIT_PRICE_BELOW_MIN_BOUND VMStatusCode = 15
	// Gas unit price submitted with the transaction is above the maximum
	// gas price set in the VM.
	GAS_UNIT_PRICE_ABOVE_MAX_BOUND VMStatusCode = 16

	// When a code module/script is published it is verified. These are the
	// possible errors that can arise from the verification process.
	// Verification Errors: 1000-1999
	UNKNOWN_VERIFICATION_ERROR              VMStatusCode = 1000
	INDEX_OUT_OF_BOUNDS                     VMStatusCode = 1001
	RANGE_OUT_OF_BOUNDS                     VMStatusCode = 1002
	INVALID_SIGNATURE_TOKEN                 VMStatusCode = 1003
	INVALID_FIELD_DEF_REFERENCE             VMStatusCode = 1004
	RECURSIVE_STRUCT_DEFINITION             VMStatusCode = 1005
	INVALID_RESOURCE_FIELD                  VMStatusCode = 1006
	INVALID_FALL_THROUGH                    VMStatusCode = 1007
	JOIN_FAILURE                            VMStatusCode = 1008
	NEGATIVE_STACK_SIZE_WITHIN_BLOCK        VMStatusCode = 1009
	UNBALANCED_STACK                        VMStatusCode = 1010
	INVALID_MAIN_FUNCTION_SIGNATURE         VMStatusCode = 1011
	DUPLICATE_ELEMENT                       VMStatusCode = 1012
	INVALID_MODULE_HANDLE                   VMStatusCode = 1013
	UNIMPLEMENTED_HANDLE                    VMStatusCode = 1014
	INCONSISTENT_FIELDS                     VMStatusCode = 1015
	UNUSED_FIELDS                           VMStatusCode = 1016
	LOOKUP_FAILED                           VMStatusCode = 1017
	VISIBILITY_MISMATCH                     VMStatusCode = 1018
	TYPE_RESOLUTION_FAILURE                 VMStatusCode = 1019
	TYPE_MISMATCH                           VMStatusCode = 1020
	MISSING_DEPENDENCY                      VMStatusCode = 1021
	POP_REFERENCE_ERROR                     VMStatusCode = 1022
	POP_RESOURCE_ERROR                      VMStatusCode = 1023
	RELEASEREF_TYPE_MISMATCH_ERROR          VMStatusCode = 1024
	BR_TYPE_MISMATCH_ERROR                  VMStatusCode = 1025
	ABORT_TYPE_MISMATCH_ERROR               VMStatusCode = 1026
	STLOC_TYPE_MISMATCH_ERROR               VMStatusCode = 1027
	STLOC_UNSAFE_TO_DESTROY_ERROR           VMStatusCode = 1028
	RET_UNSAFE_TO_DESTROY_ERROR             VMStatusCode = 1029
	RET_TYPE_MISMATCH_ERROR                 VMStatusCode = 1030
	FREEZEREF_TYPE_MISMATCH_ERROR           VMStatusCode = 1031
	FREEZEREF_EXISTS_MUTABLE_BORROW_ERROR   VMStatusCode = 1032
	BORROWFIELD_TYPE_MISMATCH_ERROR         VMStatusCode = 1033
	BORROWFIELD_BAD_FIELD_ERROR             VMStatusCode = 1034
	BORROWFIELD_EXISTS_MUTABLE_BORROW_ERROR VMStatusCode = 1035
	COPYLOC_UNAVAILABLE_ERROR               VMStatusCode = 1036
	COPYLOC_RESOURCE_ERROR                  VMStatusCode = 1037
	COPYLOC_EXISTS_BORROW_ERROR             VMStatusCode = 1038
	MOVELOC_UNAVAILABLE_ERROR               VMStatusCode = 1039
	MOVELOC_EXISTS_BORROW_ERROR             VMStatusCode = 1040
	BORROWLOC_REFERENCE_ERROR               VMStatusCode = 1041
	BORROWLOC_UNAVAILABLE_ERROR             VMStatusCode = 1042
	BORROWLOC_EXISTS_BORROW_ERROR           VMStatusCode = 1043
	CALL_TYPE_MISMATCH_ERROR                VMStatusCode = 1044
	CALL_BORROWED_MUTABLE_REFERENCE_ERROR   VMStatusCode = 1045
	PACK_TYPE_MISMATCH_ERROR                VMStatusCode = 1046
	UNPACK_TYPE_MISMATCH_ERROR              VMStatusCode = 1047
	READREF_TYPE_MISMATCH_ERROR             VMStatusCode = 1048
	READREF_RESOURCE_ERROR                  VMStatusCode = 1049
	READREF_EXISTS_MUTABLE_BORROW_ERROR     VMStatusCode = 1050
	WRITEREF_TYPE_MISMATCH_ERROR            VMStatusCode = 1051
	WRITEREF_RESOURCE_ERROR                 VMStatusCode = 1052
	WRITEREF_EXISTS_BORROW_ERROR            VMStatusCode = 1053
	WRITEREF_NO_MUTABLE_REFERENCE_ERROR     VMStatusCode = 1054
	INTEGER_OP_TYPE_MISMATCH_ERROR          VMStatusCode = 1055
	BOOLEAN_OP_TYPE_MISMATCH_ERROR          VMStatusCode = 1056
	EQUALITY_OP_TYPE_MISMATCH_ERROR         VMStatusCode = 1057
	EXISTS_RESOURCE_TYPE_MISMATCH_ERROR     VMStatusCode = 1058
	BORROWGLOBAL_TYPE_MISMATCH_ERROR        VMStatusCode = 1059
	BORROWGLOBAL_NO_RESOURCE_ERROR          VMStatusCode = 1060
	MOVEFROM_TYPE_MISMATCH_ERROR            VMStatusCode = 1061
	MOVEFROM_NO_RESOURCE_ERROR              VMStatusCode = 1062
	MOVETOSENDER_TYPE_MISMATCH_ERROR        VMStatusCode = 1063
	MOVETOSENDER_NO_RESOURCE_ERROR          VMStatusCode = 1064
	CREATEACCOUNT_TYPE_MISMATCH_ERROR       VMStatusCode = 1065
	// The self address of a module the transaction is publishing is not the sender address
	MODULE_ADDRESS_DOES_NOT_MATCH_SENDER VMStatusCode = 1066
	// The module does not have any module handles. Each module or script must have at least one
	// module handle.
	NO_MODULE_HANDLES                             VMStatusCode = 1067
	POSITIVE_STACK_SIZE_AT_BLOCK_END              VMStatusCode = 1068
	MISSING_ACQUIRES_RESOURCE_ANNOTATION_ERROR    VMStatusCode = 1069
	EXTRANEOUS_ACQUIRES_RESOURCE_ANNOTATION_ERROR VMStatusCode = 1070
	DUPLICATE_ACQUIRES_RESOURCE_ANNOTATION_ERROR  VMStatusCode = 1071
	INVALID_ACQUIRES_RESOURCE_ANNOTATION_ERROR    VMStatusCode = 1072
	GLOBAL_REFERENCE_ERROR                        VMStatusCode = 1073
	CONTRAINT_KIND_MISMATCH                       VMStatusCode = 1074
	NUMBER_OF_TYPE_ACTUALS_MISMATCH               VMStatusCode = 1075

	// These are errors that the VM might raise if a violation of internal
	// invariants takes place.
	// Invariant Violation Errors: 2000-2999
	UNKNOWN_INVARIANT_VIOLATION_ERROR VMStatusCode = 2000
	OUT_OF_BOUNDS_INDEX               VMStatusCode = 2001
	OUT_OF_BOUNDS_RANGE               VMStatusCode = 2002
	EMPTY_VALUE_STACK                 VMStatusCode = 2003
	EMPTY_CALL_STACK                  VMStatusCode = 2004
	PC_OVERFLOW                       VMStatusCode = 2005
	LINKER_ERROR                      VMStatusCode = 2006
	LOCAL_REFERENCE_ERROR             VMStatusCode = 2007
	STORAGE_ERROR                     VMStatusCode = 2008
	INTERNAL_TYPE_ERROR               VMStatusCode = 2009
	EVENT_KEY_MISMATCH                VMStatusCode = 2010

	// Errors that can arise from binary decoding (deserialization)
	// Deserializtion Errors: 3000-3999
	UNKNOWN_BINARY_ERROR         VMStatusCode = 3000
	MALFORMED                    VMStatusCode = 3001
	BAD_MAGIC                    VMStatusCode = 3002
	UNKNOWN_VERSION              VMStatusCode = 3003
	UNKNOWN_TABLE_TYPE           VMStatusCode = 3004
	UNKNOWN_SIGNATURE_TYPE       VMStatusCode = 3005
	UNKNOWN_SERIALIZED_TYPE      VMStatusCode = 3006
	UNKNOWN_OPCODE               VMStatusCode = 3007
	BAD_HEADER_TABLE             VMStatusCode = 3008
	UNEXPECTED_SIGNATURE_TYPE    VMStatusCode = 3009
	DUPLICATE_TABLE              VMStatusCode = 3010
	VERIFIER_INVARIANT_VIOLATION VMStatusCode = 3011

	// Errors that can arise at runtime
	// Runtime Errors: 4000-4999
	UNKNOWN_RUNTIME_STATUS VMStatusCode = 4000
	EXECUTED               VMStatusCode = 4001
	OUT_OF_GAS             VMStatusCode = 4002
	// We tried to access a resource that does not exist under the account.
	RESOURCE_DOES_NOT_EXIST VMStatusCode = 4003
	// We tried to create a resource under an account where that resource
	// already exists.
	RESOURCE_ALREADY_EXISTS VMStatusCode = 4004
	// We accessed an account that is evicted.
	EVICTED_ACCOUNT_ACCESS VMStatusCode = 4005
	// We tried to create an account at an address where an account already exists.
	ACCOUNT_ADDRESS_ALREADY_EXISTS VMStatusCode = 4006
	TYPE_ERROR                     VMStatusCode = 4007
	MISSING_DATA                   VMStatusCode = 4008
	DATA_FORMAT_ERROR              VMStatusCode = 4009
	INVALID_DATA                   VMStatusCode = 4010
	REMOTE_DATA_ERROR              VMStatusCode = 4011
	CANNOT_WRITE_EXISTING_RESOURCE VMStatusCode = 4012
	VALUE_SERIALIZATION_ERROR      VMStatusCode = 4013
	VALUE_DESERIALIZATION_ERROR    VMStatusCode = 4014
	// The sender is trying to publish a module named `M`, but the sender's account already
	// contains a module with this name.
	DUPLICATE_MODULE_NAME      VMStatusCode = 4015
	ABORTED                    VMStatusCode = 4016
	ARITHMETIC_ERROR           VMStatusCode = 4017
	DYNAMIC_REFERENCE_ERROR    VMStatusCode = 4018
	CODE_DESERIALIZATION_ERROR VMStatusCode = 4019
	EXECUTION_STACK_OVERFLOW   VMStatusCode = 4020
	CALL_STACK_OVERFLOW        VMStatusCode = 4021

	// A reserved status to represent an unknown vm status.
	UNKNOWN_STATUS VMStatusCode = 0xFFFFFFFFFFFFFFFF
)
