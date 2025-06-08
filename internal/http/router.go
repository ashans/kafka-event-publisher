package http

import (
	"encoding/json"
	"ev_pub/internal/config"
	"ev_pub/internal/errors"
	"ev_pub/internal/publisher"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log"
	"net/http"
)

type Router struct {
	config    config.HttpConfig
	r         chi.Router
	publisher *publisher.Publisher
}

func NewRouter(config config.AppConfig, publisher *publisher.Publisher) *Router {
	r := chi.NewRouter()
	r.Use(middleware.Recoverer)

	router := &Router{
		config:    config.Http,
		r:         r,
		publisher: publisher,
	}

	r.Get("/", infoHandler(config.AppInfo))
	r.Post(fmt.Sprintf("/%s", router.config.Path.Api), router.apiHandlerWithError)

	return router
}

func (r *Router) Start() error {
	log.Default().Println("Starting HTTP server on", r.config.Port)
	return http.ListenAndServe(fmt.Sprintf(":%d", r.config.Port), r.r)
}

func (r *Router) apiHandlerWithError(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	res, err := r.apiHandler(req)
	if err != nil {
		log.Default().Println(err.Error())
		w.WriteHeader(http.StatusBadRequest)
		errorMessage, ok := errors.AppErrMsg(err)
		if !ok {
			errorMessage = "unknown error"
		}
		json.NewEncoder(w).Encode(ErrorResponse{Error: errorMessage})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}

func (r *Router) apiHandler(req *http.Request) (PublisherResponse, error) {
	ctx := req.Context()
	var reqBody PublisherRequest
	err := json.NewDecoder(req.Body).Decode(&reqBody)
	if err != nil {
		return PublisherResponse{}, errors.Wrap(err, `error decoding request body`)
	}

	partition, offset, err := r.publisher.Publish(ctx, reqBody.Topic, reqBody.Key, reqBody.Value, reqBody.Partitioner, reqBody.Headers)
	if err != nil {
		return PublisherResponse{}, errors.Wrap(err, `error publishing request`)
	}

	return PublisherResponse{
		Partition: partition,
		Offset:    offset,
	}, nil
}

func infoHandler(info map[string]interface{}) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(info)
	}
}
