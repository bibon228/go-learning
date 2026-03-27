package payments

type PaymentMethod interface {
	Pay(usd int) int
	Cancel(id int)
	Info(id int) string
	InfoAll() string
}
type PaymentModule struct {
	info          map[int]PaymentInfo
	paymentMethod PaymentMethod
}

func NewPaymentModule(paymentMethod PaymentMethod) *PaymentModule {
	return &PaymentModule{
		paymentMethod: paymentMethod,
	}
}

/*func (p PaymentModule) Pay(description string, usd int) int {
	id := p.paymentMethod.Pay(usd)
	info := PaymentInfo{
		Description: description,
		Usd:         usd,
		Canceled:    false,
	}

}
func (p PaymentModule) Cancel(id int)      {}
func (p PaymentModule) Info(id int) string {}
func (p PaymentModule) InfoAll() string    {}
*/
