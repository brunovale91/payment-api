package main

import (
	"log"
	"net/http"

	"github.com/brunovale91/payment-api/api"
	"github.com/brunovale91/payment-api/services"
	"github.com/brunovale91/payment-api/store"
)

func main() {
	api := getPaymentApi(Config)
	if api != nil {
		log.Fatal(http.ListenAndServe(":"+Config.Port, api))
	}
}

func getPaymentApi(config *ConfigProperties) http.Handler {
	paymentStore, err := store.NewPaymentStore(&store.PaymentStoreConfig{
		URL:        config.MongoURL,
		Database:   config.Database,
		Collection: config.Collection,
	})
	if err != nil {
		log.Fatal("Failed to initialize data store")
		return nil
	}
	paymentService := services.NewPaymentService(paymentStore)
	router := api.NewApiRouter(paymentService)
	return router
}
