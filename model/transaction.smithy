$version: "2"

namespace io.baraiq.antifraud

/// The core transaction data structure.
structure Transaction {
    @required
    @jsonName("id")
    id: String,

    @jsonName("request_id")
    requestId: String,

    @required
    @jsonName("type")
    type: TransactionType,

    @required
    @jsonName("date")
    @documentation("Must be in RFC3339Nano format (e.g., 2019-12-12T12:00:00.000000000Z)")
    @pattern("^\\d{4}-\\d{2}-\\d{2}T\\d{2}:\\d{2}:\\d{2}(\\.\\d+)?(Z|[+-]\\d{2}:\\d{2})$")
    date: String,
    
    @jsonName("creation_date")
    creationDate: String,

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

    @jsonName("product_id")
    productId: String,

    @jsonName("knp_code")
    knpCode: String,

    @jsonName("oper_type")
    operType: Integer,

    @jsonName("location_ip")
    locationIp: String,

    @jsonName("location_country")
    locationCountry: String,

    // Sender fields
    @jsonName("sender_id")
    senderId: String,

    @jsonName("sender_iinbin")
    senderIINBIN: String,

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

    @jsonName("sender_card_exp_date")
    senderCardExpDate: String,

    @jsonName("sender_card_open_date")
    senderCardOpenDate: String,

    @jsonName("sender_contract_number")
    senderContractNumber: String,

    @jsonName("sender_country")
    senderCountry: String,

    @jsonName("sender_phone")
    senderPhone: String,

    @jsonName("sender_reg_date")
    senderRegDate: String,

    @jsonName("sender_bank_bic")
    senderBankBic: String,

    @jsonName("sender_bank_name")
    senderBankName: String,

    @jsonName("sender_is_client")
    senderIsClient: Boolean,

    // Recipient fields
    @jsonName("recipient_id")
    recipientId: String,

    @jsonName("recipient_iinbin")
    recipientIINBIN: String,

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

    @jsonName("recipient_card_exp_date")
    recipientCardExpDate: String,

    @jsonName("recipient_card_open_date")
    recipientCardOpenDate: String,

    @jsonName("recipient_contract_number")
    recipientContractNumber: String,

    @jsonName("recipient_country")
    recipientCountry: String,

    @jsonName("recipient_phone")
    recipientPhone: String,
    
    @jsonName("recipient_reg_date")
    recipientRegDate: String,

    @jsonName("recipient_bank_bic")
    recipientBankBic: String,

    @jsonName("recipient_bank_name")
    recipientBankName: String,

    @jsonName("recipient_is_client")
    recipientIsClient: Boolean
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
