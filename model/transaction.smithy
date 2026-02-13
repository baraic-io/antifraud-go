$version: "2"

namespace io.baraiq.antifraud

/// The core transaction data structure.
structure Transaction {
    @required
    @jsonName("id")
    @documentation("Client side merchant_order_id")
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

    @jsonName("client_id")
    @documentation("Client side client_id or client_login")
    clientId: String,

    @jsonName("client_name")
    clientName: String,

    @required
    @jsonName("client_pan")
    clientPan: String,

    @jsonName("client_card_cvv")
    clientCvv: String,

    @jsonName("client_card_holder")
    clientCardHolder: String,

    @jsonName("client_country")
    clientCountry: String,

    @jsonName("client_city")
    clientCity: String,

    @jsonName("client_phone")
    clientPhone: String,

    @jsonName("merchant_id")
    merchantId: String,

    @jsonName("merchant_descriptor")
    merchantDescriptor: String,

    @required
    @jsonName("merchant_terminal_id")
    merchantTerminalId: String,

    @jsonName("merchant_country")
    merchantCountry: String,

    @jsonName("channel")
    channel: ChannelType,

    @jsonName("location_ip")
    locationIp: String,

    @jsonName("location_country")
    locationCountry: String
}

enum TransactionType {
    DEPOSIT = "deposit"
    WITHDRAW = "withdraw"
}

enum ChannelType {
    E_COM = "E-com"
}
