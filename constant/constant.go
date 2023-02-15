package constant

type PaymentMethod struct {
	PaymentMethodIdCard     string
	PaymentMethodFlowDirect string
}

var PaymentMethodConstants PaymentMethod = PaymentMethod{
	PaymentMethodIdCard:     "CARD",
	PaymentMethodFlowDirect: "DIRECT",
}
