$version: "2"

namespace io.baraiq.antifraud

/// The core transaction data structure.
structure Transaction {
    @required
    @jsonName("id")
    id: String,

    @required
    @jsonName("type")
    type: TransactionType,

    @required
    @jsonName("date")
    @documentation("Must be in RFC3339Nano format (e.g., 2019-12-12T12:00:00.000000000Z)")
    @pattern("^\\d{4}-\\d{2}-\\d{2}T\\d{2}:\\d{2}:\\d{2}(\\.\\d+)?(Z|[+-]\\d{2}:\\d{2})$")
    date: String,

    @required
    @jsonName("amount")
    amount: String,

    @required
    @jsonName("currency")
    currency: String,

    @jsonName("description")
    description: String,

    @jsonName("channel")
    channel: ChannelType,

    @jsonName("location_ip")
    locationIp: String,

    @jsonName("location_country")
    locationCountry: String,

    // Sender fields
    @jsonName("sender_id")
    senderId: String,

    @jsonName("sender_type")
    senderType: ClientType,

    @jsonName("sender_name")
    senderName: String,

    @jsonName("sender_pan")
    senderPan: String,

    @jsonName("sender_card_cvv")
    senderCvv: String,

    @jsonName("sender_card_holder")
    senderCardHolder: String,

    @jsonName("sender_card_year_month")
    senderCardYearMonth: String,

    @jsonName("sender_contract_number")
    senderContractNumber: String,

    @jsonName("sender_country")
    senderCountry: String,

    @jsonName("sender_phone")
    senderPhone: String,

    // Recipient fields
    @jsonName("recipient_id")
    recipientId: String,

    @jsonName("recipient_type")
    recipientType: ClientType,

    @jsonName("recipient_name")
    recipientName: String,

    @jsonName("recipient_pan")
    recipientPan: String,

    @jsonName("recipient_card_cvv")
    recipientCvv: String,

    @jsonName("recipient_card_holder")
    recipientCardHolder: String,

    @jsonName("recipient_card_year_month")
    recipientCardYearMonth: String,

    @jsonName("recipient_contract_number")
    recipientContractNumber: String,

    @jsonName("recipient_country")
    recipientCountry: String,

    @jsonName("recipient_phone")
    recipientPhone: String
}

enum TransactionType {
    DEPOSIT = "deposit"
    WITHDRAW = "withdraw"
}

enum ChannelType {
    BINANCE = "binance"
    MOBILE = "mobile"
    E_COM = "e-com"
}

enum ClientType {
    PERSON = "person"
    ORGANIZATION = "organization"
}
