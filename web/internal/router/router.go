package router

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	_ "subsaggregator/docs"
	"subsaggregator/internal/repository"
	"subsaggregator/internal/service"
	"subsaggregator/internal/utils"

	"github.com/go-chi/chi/v5"
	"github.com/swaggo/http-swagger"
)

func NewRouter() *chi.Mux {
	r := chi.NewRouter()

	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"),
	))

	r.Get("/ping", pong)

	r.Post("/subscription", createSubscription)

	r.Post("/subscription/list", listSubscription)

	r.Post("/subscription/sum-price", sumSubscriptionPrices)

	r.Get("/subscription/{subscriptionId}", getOneSubscription)

	r.Post("/subscription/{subscriptionId}", updateSubscription)

	r.Delete("/subscription/{subscriptionId}", deleteSubscription)

	return r
}

// @Summary Возвращает строку "pong"
// @Description Простой эхо-метод, который возвращает фиксированное сообщение "pong"
// @Tags HealthCheck
// @Accept plain
// @Produce plain
// @Success 200 {string} pong
// @Router /ping [get]
func pong(response http.ResponseWriter, request *http.Request) {
	response.Write([]byte("pong"))
}

// createSubscription создаёт запись о подписке
// @Summary Создаёт запись о подписке
// @Description Создаёт запись о подписке
// @Tags Subscriptions
// @Accept json
// @Produce json
// @Param subscription body service.CreateSubscriptionRequest true "Параметры запроса для создания записи о подписке"
// @Success 200 {object} model.Subscription "Запись о подписке"
// @Failure 400
// @Router /subscription [post]
func createSubscription(w http.ResponseWriter, r *http.Request) {
	var req service.CreateSubscriptionRequest

	err := json.NewDecoder(r.Body).Decode(&req)

	if err != nil {
		http.Error(w, "Bad request: "+err.Error(), http.StatusBadRequest)
		return
	}

	sub, err := service.CreateSubscription(req, &repository.SubscriptionRepo{})

	if err != nil {
		http.Error(w, "Bad request: "+err.Error(), http.StatusBadRequest)
		return
	}

	utils.RespondJSON(w, sub, http.StatusOK)
}

// listSubscription получает список записей о подписках за выбранный период с фильтрацией по ИД пользователя и названию сервиса
// @Summary Получает список записей о подписках
// @Description Получает список записей о подписках за выбранный период с фильтрацией по ИД пользователя и названию сервиса
// @Tags Subscriptions
// @Accept json
// @Produce json
// @Param subscription body service.ListSubscriptionsRequest true "Параметры запроса для получения списка записей о подписках"
// @Success 200 {array} model.Subscription "Запись о подписке"
// @Failure 400
// @Failure 404
// @Router /subscription/list [post]
func listSubscription(w http.ResponseWriter, r *http.Request) {
	var req service.ListSubscriptionsRequest

	err := json.NewDecoder(r.Body).Decode(&req)

	if err != nil {
		http.Error(w, "Bad request: "+err.Error(), http.StatusBadRequest)
		return
	}

	subs, err := service.ListSubscriptions(req, &repository.SubscriptionRepo{})

	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			http.Error(w, "Not found: "+err.Error(), http.StatusNotFound)
			return
		}

		http.Error(w, "Bad request: "+err.Error(), http.StatusBadRequest)
		return
	}

	utils.RespondJSON(w, subs, http.StatusOK)
}

// sumSubscriptionPrices получает суммарную стоимость подписок за выбранный период с фильтрацией по ИД пользователя и названию сервиса
// @Summary Получает суммарную стоимость подписок
// @Description Получает суммарную стоимость подписок за выбранный период с фильтрацией по ИД пользователя и названию сервиса
// @Tags Subscriptions
// @Accept json
// @Produce json
// @Param subscription body service.SumSubscriptionsPricesRequest true "Параметры запроса для получения суммарной стоимости подписок"
// @Success 200 {integer} 100
// @Failure 400
// @Router /subscription/sum-price [post]
func sumSubscriptionPrices(w http.ResponseWriter, r *http.Request) {
	var req service.SumSubscriptionsPricesRequest

	err := json.NewDecoder(r.Body).Decode(&req)

	if err != nil {
		http.Error(w, "Bad request: "+err.Error(), http.StatusBadRequest)
		return
	}

	sumPrice, err := service.SumSubscriptionsPrices(req, &repository.SubscriptionRepo{})

	if err != nil {
		http.Error(w, "Bad request: "+err.Error(), http.StatusNotFound)
		return
	}

	utils.RespondJSON(w, &sumPrice, http.StatusOK)
}

// getOneSubscription получает запись о подписке
// @Summary Получает запись о подписке
// @Description Получает запись о подписке
// @Tags Subscriptions
// @Accept json
// @Produce json
// @Param subscriptionId path int true "Идентификатор записи о подписке"
// @Success 200 "Список записей о подписках"
// @Failure 400
// @Failure 404
// @Router /subscription/{subscriptionId} [get]
func getOneSubscription(w http.ResponseWriter, r *http.Request) {
	stringSubId := chi.URLParam(r, "subscriptionId")
	subId, err := strconv.Atoi(stringSubId)

	if err != nil {
		http.Error(w, "Bad request: "+err.Error(), http.StatusBadRequest)
		return
	}

	sub, err := service.GetOneSubscription(&repository.SubscriptionRepo{}, subId)

	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			http.Error(w, "Not found: "+err.Error(), http.StatusNotFound)
			return
		}

		http.Error(w, "Bad request: "+err.Error(), http.StatusBadRequest)
		return
	}

	utils.RespondJSON(w, sub, http.StatusOK)
}

// updateSubscription изменяет запись о подписке
// @Summary Изменяет запись о подписке
// @Description Изменяет запись о подписке
// @Tags Subscriptions
// @Accept json
// @Produce json
// @Param subscriptionId path int true "Идентификатор пользователя"
// @Param subscription body service.UpdateSubscriptionRequest true "Параметры запроса для изменения записи о подписке"
// @Success 200 {object} model.Subscription "Запись о подписке"
// @Failure 400
// @Failure 404
// @Router /subscription/{subscriptionId} [post]
func updateSubscription(w http.ResponseWriter, r *http.Request) {
	stringSubId := chi.URLParam(r, "subscriptionId")
	subId, err := strconv.Atoi(stringSubId)

	if err != nil {
		http.Error(w, "Bad request: "+err.Error(), http.StatusBadRequest)
		return
	}

	var req service.UpdateSubscriptionRequest

	err = json.NewDecoder(r.Body).Decode(&req)

	if err != nil {
		http.Error(w, "Bad request: "+err.Error(), http.StatusBadRequest)
		return
	}

	sub, err := service.UpdateSubscription(req, &repository.SubscriptionRepo{}, subId)

	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			http.Error(w, "Not found: "+err.Error(), http.StatusNotFound)
			return
		}

		http.Error(w, "Bad request: "+err.Error(), http.StatusBadRequest)
		return
	}

	utils.RespondJSON(w, sub, http.StatusOK)
}

// deleteSubscription удаляет запись о подписке
// @Summary Удаляет запись о подписке
// @Description Удаляет запись о подписке
// @Tags Subscriptions
// @Accept json
// @Produce json
// @Param subscriptionId path int true "Идентификатор пользователя"
// @Success 204
// @Failure 400
// @Failure 404
// @Router /subscription/{subscriptionId} [delete]
func deleteSubscription(w http.ResponseWriter, r *http.Request) {
	stringSubId := chi.URLParam(r, "subscriptionId")
	subId, err := strconv.Atoi(stringSubId)

	if err != nil {
		http.Error(w, "Bad request: "+err.Error(), http.StatusBadRequest)
		return
	}

	err = service.DeleteSubscription(&repository.SubscriptionRepo{}, subId)

	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			http.Error(w, "Not found: "+err.Error(), http.StatusNotFound)
			return
		}

		http.Error(w, "Bad request: "+err.Error(), http.StatusBadRequest)
		return
	}

	utils.RespondJSON(w, nil, http.StatusNoContent)
}
