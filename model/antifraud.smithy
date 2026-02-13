$version: "2"

namespace io.baraiq.antifraud

use aws.protocols#restJson1
use smithy.framework#ValidationException

/// The Antifraud Service API for transaction validation and resolution.
@restJson1
service AntifraudService {
    version: "2026-02-13"
    operations: [
        FinalizeTransaction,
        AddTransactionServiceCheck,
        ValidateTransactionByAML,
        ValidateTransactionByFC,
        ValidateTransactionByML,
        ValidateTransactionByLST,
        StoreTransaction,
        StoreServiceResolution,
        StoreFinalResolution
    ]
}

/// Finalize a transaction and get the final resolution.
@http(method: "POST", uri: "/api/fzrsvc/transaction/finalize")
operation FinalizeTransaction {
    input: AF_Transaction
    output: FinalResolution
}

/// Add a service check result to a transaction.
@http(method: "POST", uri: "/api/fzrsvc/transaction/add-service-check")
operation AddTransactionServiceCheck {
    input: ServiceResolution
    output: Unit
}

/// Validate a transaction using the AML service.
@http(method: "POST", uri: "/api/amlsvc/validate")
operation ValidateTransactionByAML {
    input: AF_Transaction
    output: ServiceResolution
}

/// Validate a transaction using the FC service.
@http(method: "POST", uri: "/api/fcsvc/validate")
operation ValidateTransactionByFC {
    input: AF_Transaction
    output: ServiceResolution
}

/// Validate a transaction using the ML service.
@http(method: "POST", uri: "/api/mlsvc/validate")
operation ValidateTransactionByML {
    input: AF_Transaction
    output: ServiceResolution
}

/// Validate a transaction using the LST service.
@http(method: "POST", uri: "/api/lstsvc/validate")
operation ValidateTransactionByLST {
    input: AF_Transaction
    output: ServiceResolution
}

/// Store a transaction.
@http(method: "POST", uri: "/api/storagesvc/store/transaction")
operation StoreTransaction {
    input: AF_Transaction
    output: Unit
}

/// Store a service resolution.
@http(method: "POST", uri: "/api/storagesvc/store/service-resolution")
operation StoreServiceResolution {
    input: ServiceResolution
    output: Unit
}

/// Store a final resolution.
@http(method: "POST", uri: "/api/storagesvc/store/final-resolution")
operation StoreFinalResolution {
    input: FinalResolution
    output: Unit
}




/// Extended definition of a transaction with Antifraud specific metadata.
structure AF_Transaction {
    @jsonName("transaction")
    transaction: Transaction,

    @required
    @jsonName("af_id")
    afId: String,

    @required
    @jsonName("af_add_date")
    afAddDate: String,

    @jsonName("af_retries")
    afRetries: AF_RetryMap
}

map AF_RetryMap {
    key: String,
    value: AF_Retry
}

structure AF_Retry {
    @required
    @jsonName("retry_count")
    retryCount: Integer,

    @required
    @jsonName("retry_max")
    retryMax: Integer
}

/// The result of a single service check.
structure ServiceResolution {
    @required
    @jsonName("af_id")
    afId: String,

    @required
    @jsonName("txn_id")
    txnId: String,

    @required
    @jsonName("id")
    id: String,

    @required
    @jsonName("date")
    @timestampFormat("date-time")
    date: Timestamp,

    @required
    @jsonName("service")
    service: String,

    @jsonName("error")
    error: String,

    @jsonName("details")
    details: StringMap,

    @jsonName("action")
    action: StringNotMapButSet, // Go: map[string]struct{}

    @required
    @jsonName("process_time")
    processTime: Long,

    @required
    @jsonName("retry")
    retry: Integer, // Go uses uint, mapped to Integer

    @required
    @jsonName("fraud")
    fraud: Boolean,

    @required
    @jsonName("validated")
    validated: Boolean,

    @required
    @jsonName("blocked")
    blocked: Boolean,

    @required
    @jsonName("alert")
    alert: Boolean,

    @required
    @jsonName("in_white_list")
    inWhiteList: Boolean
}

list StringNotMapButSet {
    member: String
}

map StringMap {
    key: String,
    value: String
}

/// The final decision for a transaction.
structure FinalResolution {
    @required
    @jsonName("af_id")
    afId: String,

    @required
    @jsonName("af_transaction")
    afTransaction: AF_Transaction,

    @required
    @jsonName("id")
    id: String,

    @jsonName("af_error")
    error: String,

    @jsonName("af_details")
    details: Document, // Go: map[string]interface{}

    @required
    @jsonName("af_add_date")
    @timestampFormat("date-time")
    addDate: Timestamp,

    @required
    @jsonName("af_finalized_date")
    @timestampFormat("date-time")
    finalizedDate: Timestamp,

    @required
    @jsonName("af_finalized_action")
    finalizedAction: String,

    @required
    @jsonName("af_process_time")
    processTime: Long,

    @jsonName("af_validated_services")
    validatedServices: StringList,

    @jsonName("af_unvalidated_services")
    unvalidatedServices: StringList,

    @required
    @jsonName("af_retry_attempt")
    retry: Integer,

    @required
    @jsonName("af_fraud")
    fraud: Boolean,

    @required
    @jsonName("af_validated")
    validated: Boolean,

    @required
    @jsonName("af_blocked")
    blocked: Boolean,

    @required
    @jsonName("af_alert")
    alert: Boolean
}

list StringList {
    member: String
}
