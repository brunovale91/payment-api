package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/brunovale91/payment-api/types"
)

var validPayment = &types.Payment{
	OrganisationId: "test",
	Type:           "Payment",
	Attributes: &types.PaymentAttributes{
		Amount: 3,
		BeneficiaryParty: &types.PaymentParty{
			BankId:     "id",
			BankIdCode: "code",
			Name:       "name",
		},
		DebtorParty: &types.PaymentParty{
			BankId:     "id2",
			BankIdCode: "code2",
			Name:       "name2",
		},
		EndToEndReference: "test1",
	},
}

var validPaymentUpdate = &types.Payment{
	Attributes: &types.PaymentAttributes{
		Amount: 5,
		BeneficiaryParty: &types.PaymentParty{
			BankId:     "id",
			BankIdCode: "code",
			Name:       "name",
		},
		DebtorParty: &types.PaymentParty{
			BankId:     "id2",
			BankIdCode: "code2",
			Name:       "name2",
		},
		EndToEndReference: "test1",
	},
}

var invalidPaymentUpdate = &types.Payment{
	Attributes: &types.PaymentAttributes{
		Amount: 5,
		BeneficiaryParty: &types.PaymentParty{
			BankId:     "id",
			BankIdCode: "code",
			Name:       "name",
		},
		DebtorParty: &types.PaymentParty{
			BankId:     "id2",
			BankIdCode: "code2",
		},
		EndToEndReference: "test1",
	},
}

var invalidPaymentType = &types.Payment{
	OrganisationId: "test",
	Type:           "Payment1",
	Attributes: &types.PaymentAttributes{
		Amount: 3,
		BeneficiaryParty: &types.PaymentParty{
			BankId:     "id",
			BankIdCode: "code",
			Name:       "name",
		},
		DebtorParty: &types.PaymentParty{
			BankId:     "id2",
			BankIdCode: "code2",
			Name:       "name2",
		},
		EndToEndReference: "test1",
	},
}

var invalidPaymentAmout = &types.Payment{
	OrganisationId: "test",
	Type:           "Payment1",
	Attributes: &types.PaymentAttributes{
		Amount: -1,
		BeneficiaryParty: &types.PaymentParty{
			BankId:     "id",
			BankIdCode: "code",
			Name:       "name",
		},
		DebtorParty: &types.PaymentParty{
			BankId:     "id2",
			BankIdCode: "code2",
			Name:       "name2",
		},
		EndToEndReference: "test1",
	},
}

func TestCreatePayment(t *testing.T) {
	ts := httptest.NewServer(getPaymentApi(TestConfig))
	defer ts.Close()
	deleteAllPayments(ts, t)

	res := createPayment(ts, t, createPaymentBody(t, invalidPaymentAmout))
	httpError := parseHttpError(res)
	res.Body.Close()
	if res.StatusCode != 400 {
		t.Errorf("Status code should be 400: is %d", res.StatusCode)
	}
	if len(httpError.Messages) == 0 {
		t.Errorf("Http Error should have error messages")
	}

	res = createPayment(ts, t, createPaymentBody(t, invalidPaymentType))
	res.Body.Close()
	if res.StatusCode != 400 {
		t.Errorf("Status code should be 400: is %d", res.StatusCode)
	}

	res = createPayment(ts, t, createPaymentBody(t, validPayment))
	payment := parsePayment(res)
	res.Body.Close()
	if payment.Id == "" {
		t.Error("Created payment should have id")
	}

	res = getPayment(ts, t, payment.Id)
	payment = parsePayment(res)
	if res.StatusCode != 200 {
		t.Errorf("Status code should be 200: is %d", res.StatusCode)
	}
}

func TestUpdatePayment(t *testing.T) {
	ts := httptest.NewServer(getPaymentApi(TestConfig))
	defer ts.Close()
	deleteAllPayments(ts, t)

	res := updatePayment(ts, t, "invalid", createPaymentBody(t, validPaymentUpdate))
	if res.StatusCode != 404 {
		t.Errorf("Status code should be 404: is %d", res.StatusCode)
	}

	res = createPayment(ts, t, createPaymentBody(t, validPayment))
	payment := parsePayment(res)
	res.Body.Close()

	res = updatePayment(ts, t, payment.Id, createPaymentBody(t, validPaymentUpdate))
	payment = parsePayment(res)

	if payment.Attributes.Amount != validPaymentUpdate.Attributes.Amount {
		t.Errorf("Payment amount should be %f: is %f", validPaymentUpdate.Attributes.Amount, payment.Attributes.Amount)
	}

	res = updatePayment(ts, t, payment.Id, createPaymentBody(t, invalidPaymentUpdate))
	if res.StatusCode != 400 {
		t.Errorf("Status code should be 400: is %d", res.StatusCode)
	}

}

func TestDeletePayment(t *testing.T) {
	ts := httptest.NewServer(getPaymentApi(TestConfig))
	defer ts.Close()
	deleteAllPayments(ts, t)

	res := getPayments(ts, t)
	payments := parsePayments(res)
	res.Body.Close()
	if len(payments.Data) > 0 {
		t.Errorf("Payment list size should be 0: is %d", len(payments.Data))
	}

	res = createPayment(ts, t, createPaymentBody(t, validPayment))
	payment := parsePayment(res)
	res.Body.Close()

	res = getPayments(ts, t)
	payments = parsePayments(res)
	res.Body.Close()
	if len(payments.Data) != 1 {
		t.Errorf("Payment list size should be 1: is %d", len(payments.Data))
	}

	res = deletePayment(ts, t, "invalid")
	if res.StatusCode != 404 {
		t.Errorf("Status code should be 404: is %d", res.StatusCode)
	}

	res = deletePayment(ts, t, payment.Id)
	res = getPayments(ts, t)
	payments = parsePayments(res)
	res.Body.Close()
	if len(payments.Data) > 0 {
		t.Errorf("Payment list size should be 0: is %d", len(payments.Data))
	}
}

func TestGetPayment(t *testing.T) {
	ts := httptest.NewServer(getPaymentApi(TestConfig))
	defer ts.Close()

	deleteAllPayments(ts, t)
	res := createPayment(ts, t, createPaymentBody(t, validPayment))
	payment := parsePayment(res)
	res.Body.Close()

	res = getPayment(ts, t, payment.Id)
	payment = parsePayment(res)

	if payment.OrganisationId != validPayment.OrganisationId {
		t.Errorf("Payment organisation id should be %s: is %s", validPayment.OrganisationId, payment.OrganisationId)
	}

	res = getPayment(ts, t, "invalid")
	if res.StatusCode != 404 {
		t.Errorf("Status code should be 404: is %d", res.StatusCode)
	}

}

func TestGetPayments(t *testing.T) {
	ts := httptest.NewServer(getPaymentApi(TestConfig))
	defer ts.Close()

	deleteAllPayments(ts, t)

	res := getPayments(ts, t)
	if res.StatusCode != 200 {
		t.Errorf("Status should be 200: is %d", res.StatusCode)
	}
	payments := parsePayments(res)
	res.Body.Close()

	if len(payments.Data) > 0 {
		t.Errorf("Payment list size should be 0: is %d", len(payments.Data))
	}

	res = createPayment(ts, t, createPaymentBody(t, validPayment))
	res.Body.Close()

	res = getPayments(ts, t)
	if res.StatusCode != 200 {
		t.Errorf("Status should be 200: is %d", res.StatusCode)
	}
	payments = parsePayments(res)
	res.Body.Close()

	if len(payments.Data) != 1 {
		t.Errorf("Payment list size should be 1: is %d", len(payments.Data))
	}
}

func getPayments(ts *httptest.Server, t *testing.T) *http.Response {
	res, err := http.Get(ts.URL + "/v1/api/payments")
	if err != nil {
		log.Fatal(err)
		t.Errorf("Failed to get payments: %s", err.Error())
	}
	return res
}

func getPayment(ts *httptest.Server, t *testing.T, id string) *http.Response {
	res, err := http.Get(ts.URL + "/v1/api/payments/" + id)
	if err != nil {
		log.Fatal(err)
		t.Errorf("Failed to get payment with id %s: %s", id, err.Error())
	}
	return res
}

func createPayment(ts *httptest.Server, t *testing.T, reqBody []byte) *http.Response {
	res, err := http.Post(ts.URL+"/v1/api/payments", "application/json", bytes.NewReader(reqBody))
	if err != nil {
		log.Fatal(err)
		t.Errorf("Failed to create payment: %s", err.Error())
	}
	return res
}

func updatePayment(ts *httptest.Server, t *testing.T, id string, reqBody []byte) *http.Response {
	req, err := http.NewRequest("PUT", ts.URL+"/v1/api/payments/"+id, bytes.NewReader(reqBody))
	if err != nil {
		log.Fatal(err)
		t.Errorf("Failed to update payment with id %s: %s", id, err.Error())
	}
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
		t.Errorf("Failed to update payment with id %s: %s", id, err.Error())
	}
	return res
}

func deletePayment(ts *httptest.Server, t *testing.T, id string) *http.Response {
	req, err := http.NewRequest("DELETE", ts.URL+"/v1/api/payments/"+id, nil)
	if err != nil {
		log.Fatal(err)
		t.Errorf("Failed to delete payment with id %s: %s", id, err.Error())
	}
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
		t.Errorf("Failed to delete payment with id %s: %s", id, err.Error())
	}
	return res
}

func deleteAllPayments(ts *httptest.Server, t *testing.T) {
	res := getPayments(ts, t)
	if res.StatusCode != 200 {
		t.Errorf("Status should be 200: is %d", res.StatusCode)
	}
	payments := parsePayments(res)
	res.Body.Close()

	for _, payment := range payments.Data {
		deletePayment(ts, t, payment.Id)
	}
}

func createPaymentBody(t *testing.T, payment *types.Payment) []byte {
	req, err := json.Marshal(payment)
	if err != nil {
		log.Fatal(err)
		t.Errorf("Failed to encode payment: %s", err.Error())
	}
	return req
}

func parsePayments(res *http.Response) *types.Payments {
	var payments types.Payments
	json.NewDecoder(res.Body).Decode(&payments)
	return &payments
}

func parsePayment(res *http.Response) *types.Payment {
	var payment types.Payment
	json.NewDecoder(res.Body).Decode(&payment)
	return &payment
}

func parsePaymentDelete(res *http.Response) *types.PaymentDelete {
	var paymentDelete types.PaymentDelete
	json.NewDecoder(res.Body).Decode(&paymentDelete)
	return &paymentDelete
}

func parseHttpError(res *http.Response) *types.HttpError {
	var httpError types.HttpError
	json.NewDecoder(res.Body).Decode(&httpError)
	return &httpError
}
