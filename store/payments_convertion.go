package store

import (
	"github.com/brunovale91/payment-api/types"
	"go.mongodb.org/mongo-driver/bson"
)

func docToPayment(payment interface{}) *types.Payment {
	if payment != nil {
		paymentBson := payment.(bson.D).Map()
		return &types.Payment{
			Id:             paymentBson["Id"].(string),
			Version:        paymentBson["Version"].(int64),
			OrganisationId: paymentBson["OrganisationId"].(string),
			Type:           paymentBson["Type"].(string),
			Attributes:     docToAttributes(paymentBson["Attributes"]),
		}
	}
	return nil
}

func docToAttributes(attributes interface{}) *types.PaymentAttributes {
	if attributes != nil {
		attBson := attributes.(bson.D).Map()
		return &types.PaymentAttributes{
			Amount:            attBson["Amount"].(float64),
			EndToEndReference: attBson["EndToEndReference"].(string),
			BeneficiaryParty:  docToParty(attBson["BeneficiaryParty"]),
			DebtorParty:       docToParty(attBson["DebtorParty"]),
		}
	}
	return nil
}

func docToParty(party interface{}) *types.PaymentParty {
	if party != nil {
		partyBson := party.(bson.D).Map()
		return &types.PaymentParty{
			BankId:     partyBson["BankId"].(string),
			BankIdCode: partyBson["BankIdCode"].(string),
			Name:       partyBson["Name"].(string),
		}
	}
	return nil
}

func paymentToDoc(payment *types.Payment) bson.M {
	if payment != nil {
		return bson.M{
			"_id":            payment.Id,
			"Id":             payment.Id,
			"OrganisationId": payment.OrganisationId,
			"Type":           payment.Type,
			"Version":        payment.Version,
			"Attributes":     attributesToDoc(payment.Attributes),
		}
	}
	return nil
}

func attributesToDoc(attributes *types.PaymentAttributes) bson.M {
	if attributes != nil {
		return bson.M{
			"Amount":            attributes.Amount,
			"BeneficiaryParty":  partyToDoc(attributes.BeneficiaryParty),
			"DebtorParty":       partyToDoc(attributes.DebtorParty),
			"EndToEndReference": attributes.EndToEndReference,
		}
	}
	return nil
}

func partyToDoc(party *types.PaymentParty) bson.M {
	if party != nil {
		return bson.M{
			"BankId":     party.BankId,
			"BankIdCode": party.BankIdCode,
			"Name":       party.Name,
		}
	}
	return nil
}
