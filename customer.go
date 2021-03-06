package stripe

import (
	"encoding/json"
)

// CustomerParams is the set of parameters that can be used when creating or updating a customer.
// For more details see https://stripe.com/docs/api#create_customer and https://stripe.com/docs/api#update_customer.
type CustomerParams struct {
	Params         `form:"*"`
	AccountBalance *int64                         `form:"account_balance"`
	BusinessVATID  *string                        `form:"business_vat_id"`
	Coupon         *string                        `form:"coupon"`
	DefaultSource  *string                        `form:"default_source"`
	Description    *string                        `form:"description"`
	Email          *string                        `form:"email"`
	Plan           *string                        `form:"plan"`
	Quantity       *int64                         `form:"quantity"`
	Shipping       *CustomerShippingDetailsParams `form:"shipping"`
	Source         *SourceParams                  `form:"*"` // SourceParams has custom encoding so brought to top level with "*"
	TaxPercent     *float64                       `form:"tax_percent"`
	Token          *string                        `form:"-"` // This doesn't seem to be used?
	TrialEnd       *int64                         `form:"trial_end"`
}

// CustomerShippingDetailsParams is the structure containing shipping information.
type CustomerShippingDetailsParams struct {
	Address *AddressParams `form:"address"`
	Name    *string        `form:"name"`
	Phone   *string        `form:"phone"`
}

// SetSource adds valid sources to a CustomerParams object,
// returning an error for unsupported sources.
func (cp *CustomerParams) SetSource(sp interface{}) error {
	source, err := SourceParamsFor(sp)
	cp.Source = source
	return err
}

// CustomerListParams is the set of parameters that can be used when listing customers.
// For more details see https://stripe.com/docs/api#list_customers.
type CustomerListParams struct {
	ListParams   `form:"*"`
	Created      *int64            `form:"created"`
	CreatedRange *RangeQueryParams `form:"created"`
}

// Customer is the resource representing a Stripe customer.
// For more details see https://stripe.com/docs/api#customers.
type Customer struct {
	AccountBalance int64                    `json:"account_balance"`
	BusinessVATID  string                   `json:"business_vat_id"`
	Created        int64                    `json:"created"`
	Currency       Currency                 `json:"currency"`
	DefaultSource  *PaymentSource           `json:"default_source"`
	Deleted        bool                     `json:"deleted"`
	Delinquent     bool                     `json:"delinquent"`
	Description    string                   `json:"description"`
	Discount       *Discount                `json:"discount"`
	Email          string                   `json:"email"`
	ID             string                   `json:"id"`
	Livemode       bool                     `json:"livemode"`
	Metadata       map[string]string        `json:"metadata"`
	Shipping       *CustomerShippingDetails `json:"shipping"`
	Sources        *SourceList              `json:"sources"`
	Subscriptions  *SubscriptionList        `json:"subscriptions"`
}

// CustomerList is a list of customers as retrieved from a list endpoint.
type CustomerList struct {
	ListMeta
	Data []*Customer `json:"data"`
}

// CustomerShippingDetails is the structure containing shipping information.
type CustomerShippingDetails struct {
	Address Address `json:"address"`
	Name    string  `json:"name"`
	Phone   string  `json:"phone"`
}

// UnmarshalJSON handles deserialization of a Customer.
// This custom unmarshaling is needed because the resulting
// property may be an id or the full struct if it was expanded.
func (c *Customer) UnmarshalJSON(data []byte) error {
	if id, ok := ParseID(data); ok {
		c.ID = id
		return nil
	}

	type customer Customer
	var v customer
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	*c = Customer(v)
	return nil
}
