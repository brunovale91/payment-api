package api

import (
	"github.com/brunovale91/payment-api/types"
	"github.com/xeipuuv/gojsonschema"
)

// Validate payment against payment json schema
func isValidPayment(payment *types.Payment) []string {
	return isValid(getPaymentSchema(), payment)
}

// Validate attributes against attributes json schema
func isValidAtrributes(attributes *types.PaymentAttributes) []string {
	return isValid(getPaymentAttributesSchema(), attributes)
}

func isValid(schema map[string]interface{}, value interface{}) []string {
	schemaLoader := gojsonschema.NewGoLoader(schema)
	valueLoader := gojsonschema.NewGoLoader(value)
	result, err := gojsonschema.Validate(schemaLoader, valueLoader)
	messages := make([]string, 0)
	if err != nil {
		messages = append(messages, "Failed to validate")
	}
	if !result.Valid() {
		for _, desc := range result.Errors() {
			messages = append(messages, desc.Description())
		}
	}
	return returnMessages(messages)
}

func returnMessages(messages []string) []string {
	if len(messages) > 0 {
		return messages
	}
	return nil
}

func getPaymentSchema() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"type": map[string]interface{}{
				"type": "string",
				"enum": [...]string{"Payment"},
			},
			"id": map[string]interface{}{
				"type": "string",
			},
			"organisation_id": map[string]interface{}{
				"type": "string",
			},
			"attributes": getPaymentAttributesSchema(),
		},
		"required": [...]string{"type", "organisation_id"},
	}
}

func getPaymentAttributesSchema() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"amount": map[string]interface{}{
				"type":             "number",
				"exclusiveMinimum": 0,
			},
			"beneficiary_party": getPaymentPartySchema(),
			"debtor_party":      getPaymentPartySchema(),
			"end_to_end_reference": map[string]interface{}{
				"type": "string",
			},
		},
		"required": [...]string{"amount", "beneficiary_party", "debtor_party", "end_to_end_reference"},
	}
}

func getPaymentPartySchema() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"bank_id": map[string]interface{}{
				"type": "string",
			},
			"bank_id_code": map[string]interface{}{
				"type": "string",
			},
			"name": map[string]interface{}{
				"type": "string",
			},
		},
		"required": [...]string{"bank_id", "bank_id_code", "name"},
	}
}
