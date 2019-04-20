package types

type HttpError struct {
	StatusText string   `json:"status"`
	Messages   []string `json:"messages"`
}

type Payments struct {
	Data  []*Payment `json:"data"`
	Links *Links     `json:"links"`
}

type Links struct {
	Self string `json:"self"`
}

type PaymentDelete struct {
	Deleted bool `json:"deleted,omitempty"`
}

type Payment struct {
	Type           string             `json:"type,omitempty"`
	Id             string             `json:"id,omitempty"`
	Version        int64              `json:"version"`
	OrganisationId string             `json:"organisation_id,omitempty"`
	Attributes     *PaymentAttributes `json:"attributes,omitempty"`
}

type PaymentAttributes struct {
	Amount            float64       `json:"amount,omitempty"`
	BeneficiaryParty  *PaymentParty `json:"beneficiary_party,omitempty"`
	DebtorParty       *PaymentParty `json:"debtor_party,omitempty"`
	EndToEndReference string        `json:"end_to_end_reference,omitempty"`
}

type PaymentParty struct {
	BankId     string `json:"bank_id,omitempty"`
	BankIdCode string `json:"bank_id_code,omitempty"`
	Name       string `json:"name,omitempty"`
}
