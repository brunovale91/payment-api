package api

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/brunovale91/payment-api/services"
	"github.com/brunovale91/payment-api/types"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
)

const paymentIdParam = "paymentID"

var InternalError = &types.HttpError{StatusText: "Internal Error"}
var BadRequest = &types.HttpError{StatusText: "Bad request"}
var NotFound = &types.HttpError{StatusText: "Payment not found"}
var paymentsSelf = "http://localhost:8080/v1/api/payments"

func NewApiRouter(paymentService services.PaymentService) *chi.Mux {
	router := chi.NewRouter()
	router.Use(
		render.SetContentType(render.ContentTypeJSON),
		middleware.RequestID,
		middleware.RealIP,
		middleware.DefaultCompress,
		middleware.RedirectSlashes,
		middleware.Logger,
		middleware.Recoverer,
		middleware.Timeout(60*time.Second))

	router.Route("/v1", func(r chi.Router) {
		r.Mount("/api/payments", addRoutes(paymentService))
	})

	walkFunc := func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		log.Printf("%s %s\n", method, route)
		return nil
	}
	if err := chi.Walk(router, walkFunc); err != nil {
		log.Panicf("Logging err: %s\n", err.Error())
	}

	return router
}

func addRoutes(paymentService services.PaymentService) *chi.Mux {
	router := chi.NewRouter()
	setGetPaymentById(router, paymentService)
	setDeletePayment(router, paymentService)
	setUpdatePayment(router, paymentService)
	setCreatePayment(router, paymentService)
	setGetPayments(router, paymentService)
	return router
}

func setGetPaymentById(router *chi.Mux, paymentService services.PaymentService) {
	router.Get("/{"+paymentIdParam+"}", func(w http.ResponseWriter, r *http.Request) {
		paymentID := chi.URLParam(r, paymentIdParam)
		payment, err := paymentService.GetPayment(paymentID)
		if err != nil {
			renderInternalError(router, w, r)
		} else if payment != nil {
			render.JSON(w, r, payment)
		} else {
			renderNotFound(router, w, r)
		}
	})
}

func setDeletePayment(router *chi.Mux, paymentService services.PaymentService) {
	router.Delete("/{"+paymentIdParam+"}", func(w http.ResponseWriter, r *http.Request) {
		paymentID := chi.URLParam(r, paymentIdParam)
		deleted, err := paymentService.DeletePayment(paymentID)
		if err != nil {
			renderInternalError(router, w, r)
		} else if deleted {
			render.JSON(w, r, &types.PaymentDelete{
				Deleted: deleted,
			})
		} else {
			renderNotFound(router, w, r)
		}
	})
}

func setUpdatePayment(router *chi.Mux, paymentService services.PaymentService) {
	router.Put("/{"+paymentIdParam+"}", func(w http.ResponseWriter, r *http.Request) {
		paymentID := chi.URLParam(r, paymentIdParam)
		var payment types.Payment
		json.NewDecoder(r.Body).Decode(&payment)

		errors := isValidAtrributes(payment.Attributes)
		if errors != nil {
			renderBadRequest(router, w, r, errors)
			return
		}

		updatedPayment, err := paymentService.UpdatePayment(paymentID, payment.Attributes)
		if err != nil {
			renderInternalError(router, w, r)
		} else if updatedPayment != nil {
			render.JSON(w, r, updatedPayment)
		} else {
			renderNotFound(router, w, r)
		}
	})
}

func setCreatePayment(router *chi.Mux, paymentService services.PaymentService) {
	router.Post("/", func(w http.ResponseWriter, r *http.Request) {
		var payment types.Payment
		json.NewDecoder(r.Body).Decode(&payment)

		payment.Version = 0
		errors := isValidPayment(&payment)
		if errors != nil {
			renderBadRequest(router, w, r, errors)
			return
		}

		createdPayment, err := paymentService.CreatePayment(&payment)
		if err != nil {
			renderInternalError(router, w, r)
		} else {
			render.JSON(w, r, createdPayment)
		}
	})
}

func setGetPayments(router *chi.Mux, paymentService services.PaymentService) {
	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		payments, err := paymentService.GetPayments()
		if err != nil {
			renderInternalError(router, w, r)
		} else {
			render.JSON(w, r, &types.Payments{
				Data: payments,
				Links: &types.Links{
					Self: paymentsSelf,
				},
			})
		}
	})
}

func renderNotFound(router *chi.Mux, w http.ResponseWriter, r *http.Request) {
	render.Status(r, 404)
	render.JSON(w, r, NotFound)
}

func renderBadRequest(router *chi.Mux, w http.ResponseWriter, r *http.Request, messages []string) {
	render.Status(r, 400)
	render.JSON(w, r, &types.HttpError{
		StatusText: BadRequest.StatusText,
		Messages:   messages,
	})
}

func renderInternalError(router *chi.Mux, w http.ResponseWriter, r *http.Request) {
	render.Status(r, 500)
	render.JSON(w, r, InternalError)
}
