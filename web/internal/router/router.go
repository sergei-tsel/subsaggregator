package router

import (
	"encoding/json"
	"net/http"
	"strconv"
	"subsaggregator/internal/repository"
	"subsaggregator/internal/service"
	"subsaggregator/internal/utils"

	"github.com/go-chi/chi/v5"
)

func NewRouter() *chi.Mux {
	r := chi.NewRouter()

	r.Get("/ping", pong)

	r.Post("/subscription", createSubscription)

	r.Post("/subscription/list", listSubscription)

	r.Post("/subscription/sumPrices", sumSubscriptionPrices)

	r.Get("/subscription/{subscriptionId}", getOneSubscription)

	r.Post("/subscription/{subscriptionId}", updateSubscription)

	r.Post("/subscription/{subscriptionId}/delete", deleteSubscription)

	return r
}

func pong(response http.ResponseWriter, request *http.Request) {
	response.Write([]byte("pong"))
}

func listSubscription(w http.ResponseWriter, r *http.Request) {
	var req service.ListSubscriptionsRequest

	err := json.NewDecoder(r.Body).Decode(&req)

	if err != nil {
		http.Error(w, "Invalid request: "+err.Error(), http.StatusBadRequest)
		return
	}

	subs, err := service.ListSubscriptions(req, &repository.SubscriptionRepo{})

	if err != nil {
		http.Error(w, "Subscriptions not found: "+err.Error(), http.StatusNotFound)
		return
	}

	utils.RespondJSON(w, subs, http.StatusOK)
}

func createSubscription(w http.ResponseWriter, r *http.Request) {
	var req service.CreateSubscriptionRequest

	err := json.NewDecoder(r.Body).Decode(&req)

	if err != nil {
		http.Error(w, "Invalid request: "+err.Error(), http.StatusBadRequest)
		return
	}

	sub, err := service.CreateSubscription(req, &repository.SubscriptionRepo{})

	if err != nil {
		http.Error(w, "Invalid request: "+err.Error(), http.StatusBadRequest)
		return
	}

	utils.RespondJSON(w, sub, http.StatusOK)
}

func sumSubscriptionPrices(w http.ResponseWriter, r *http.Request) {
	var req service.SumSubscriptionsRequest

	err := json.NewDecoder(r.Body).Decode(&req)

	if err != nil {
		http.Error(w, "Invalid request: "+err.Error(), http.StatusBadRequest)
		return
	}

	sumPrice, err := service.SumSubscriptionsPrices(req, &repository.SubscriptionRepo{})

	if err != nil {
		http.Error(w, "Subscriptions not found: "+err.Error(), http.StatusNotFound)
		return
	}

	utils.RespondJSON(w, sumPrice, http.StatusOK)
}

func getOneSubscription(w http.ResponseWriter, r *http.Request) {
	stringSubId := chi.URLParam(r, "subscriptionId")
	subId, err := strconv.Atoi(stringSubId)

	if err != nil {
		http.Error(w, "Invalid request: "+err.Error(), http.StatusBadRequest)
		return
	}

	sub, err := service.GetOneSubscription(&repository.SubscriptionRepo{}, subId)

	if err != nil {
		http.Error(w, "Subscription not found: "+err.Error(), http.StatusNotFound)
		return
	}

	utils.RespondJSON(w, sub, http.StatusOK)
}

func updateSubscription(w http.ResponseWriter, r *http.Request) {
	stringSubsId := chi.URLParam(r, "subscriptionId")
	subsId, err := strconv.Atoi(stringSubsId)

	if err != nil {
		http.Error(w, "Invalid request: "+err.Error(), http.StatusBadRequest)
		return
	}

	var req service.UpdateSubscriptionRequest

	err = json.NewDecoder(r.Body).Decode(&req)

	if err != nil {
		http.Error(w, "Invalid request: "+err.Error(), http.StatusBadRequest)
		return
	}

	sub, err := service.UpdateSubscription(req, &repository.SubscriptionRepo{}, subsId)

	if err != nil {
		http.Error(w, "Subscription not found: "+err.Error(), http.StatusNotFound)
		return
	}

	utils.RespondJSON(w, sub, http.StatusOK)
}

func deleteSubscription(w http.ResponseWriter, r *http.Request) {
	stringSubsId := chi.URLParam(r, "subscriptionId")
	subsId, err := strconv.Atoi(stringSubsId)

	if err != nil {
		http.Error(w, "Invalid request: "+err.Error(), http.StatusBadRequest)
		return
	}

	err = service.DeleteSubscription(&repository.SubscriptionRepo{}, subsId)

	if err != nil {
		http.Error(w, "Subscription not found: "+err.Error(), http.StatusNotFound)
		return
	}

	utils.RespondJSON(w, nil, http.StatusOK)
}
