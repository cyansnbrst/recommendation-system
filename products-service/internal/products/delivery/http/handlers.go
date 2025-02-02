package http

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"

	"cyansnbrst/products-service/config"
	"cyansnbrst/products-service/internal/models"
	"cyansnbrst/products-service/internal/products"
	"cyansnbrst/products-service/pkg/db"
	erp "cyansnbrst/products-service/pkg/error_responses"
	kf "cyansnbrst/products-service/pkg/kafka"
	"cyansnbrst/products-service/pkg/utils"
)

// Validation errors
var (
	errNameRequired = errors.New("name is required")
	errTagsRequired = errors.New("tags are required")
)

// Products handlers
type productsHandlers struct {
	cfg                *config.Config
	productsUC         products.UseCase
	logger             *zap.Logger
	kafkaUserWriter    *kafka.Writer
	kafkaProductWriter *kafka.Writer
}

// Products handlers constructor
func NewProductsHandlers(cfg *config.Config, productsUC products.UseCase, logger *zap.Logger, kafkaUserWriter *kafka.Writer, kafkaProductWriter *kafka.Writer) products.Handlers {
	return &productsHandlers{
		cfg:                cfg,
		productsUC:         productsUC,
		logger:             logger,
		kafkaUserWriter:    kafkaUserWriter,
		kafkaProductWriter: kafkaProductWriter,
	}
}

//	@Summary		Get products's info
//	@Description	Retrieves info about the product.
//	@Tags			products
//	@Produce		json
//	@Security		cookieAuth
//	@Param			id	path		int						true	"Product ID"
//	@Success		200	{object}	models.ProductResponse	"success response with product"
//	@Failure		404	{object}	models.ErrorResponse	"not found error"
//	@Failure		500	{object}	models.ErrorResponse	"internal server error"
//	@Router			/view/{id} [get]
func (h *productsHandlers) Get() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := utils.ReadIDParam(r)
		if err != nil {
			erp.NotFoundResponse(w, r, h.logger)
			return
		}

		product, err := h.productsUC.Get(id)
		if err != nil {
			erp.NotFoundResponse(w, r, h.logger)
			return
		}

		messagePayload := kf.KafkaMessage{
			Action: "view_products",
			Time:   time.Now().Format(time.RFC3339),
			Tags:   nil,
		}

		err = h.productsUC.SendToKafka(r.Context(), strconv.Itoa(int(id)), messagePayload, h.kafkaProductWriter)
		if err != nil {
			erp.ServerErrorResponse(w, r, h.logger, err)
			return
		}
		err = utils.WriteJSON(w, http.StatusOK, utils.Envelope{
			"product": product,
		}, nil)
		if err != nil {
			erp.ServerErrorResponse(w, r, h.logger, err)
		}
	}
}

//	@Summary		Create a new product
//	@Description	Creates a new product (admin-only).
//	@Tags			products
//	@Accept			json
//	@Produce		json
//	@Security		cookieAuth
//	@Success		200	{object}	models.SuccessResponse	"success response with product"
//	@Failure		400	{object}	models.ErrorResponse	"bad request error"
//	@Failure		500	{object}	models.ErrorResponse	"internal server error"
//	@Router			/create [post]
func (h *productsHandlers) Create() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var requestData models.CreateProductDTO

		if err := utils.ReadJSON(w, r, &requestData); err != nil {
			erp.BadRequestResponse(w, r, h.logger, err)
			return
		}

		if requestData.Name == "" {
			erp.BadRequestResponse(w, r, h.logger, errNameRequired)
			return
		}

		if requestData.Tags == nil {
			erp.BadRequestResponse(w, r, h.logger, errTagsRequired)
			return
		}

		id, err := h.productsUC.Create(requestData.Name, requestData.Tags)
		if err != nil {
			erp.ServerErrorResponse(w, r, h.logger, err)
			return
		}

		messagePayload := kf.KafkaMessage{
			Action: "product_create",
			Time:   time.Now().Format(time.RFC3339),
			Tags:   requestData.Tags,
		}

		err = h.productsUC.SendToKafka(r.Context(), strconv.Itoa(int(id)), messagePayload, h.kafkaProductWriter)
		if err != nil {
			erp.ServerErrorResponse(w, r, h.logger, err)
			return
		}

		err = utils.WriteJSON(w, http.StatusCreated, utils.Envelope{
			"message": "product created successfully",
		}, nil)
		if err != nil {
			erp.ServerErrorResponse(w, r, h.logger, err)
		}
	}
}

//	@Summary		Edit a product
//	@Description	Edits an existing product (admin-only).
//	@Tags			products
//	@Accept			json
//	@Produce		json
//	@Security		cookieAuth
//	@Param			id	path		int						true	"Product ID"
//	@Success		200	{object}	models.SuccessResponse	"success response with product"
//	@Failure		400	{object}	models.ErrorResponse	"bad request error"
//	@Failure		404	{object}	models.ErrorResponse	"not found error"
//	@Failure		500	{object}	models.ErrorResponse	"internal server error"
//	@Router			/update/{id} [put]
func (h *productsHandlers) Update() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := utils.ReadIDParam(r)
		if err != nil {
			erp.NotFoundResponse(w, r, h.logger)
			return
		}

		var requestData models.UpdateProductDTO

		if err = utils.ReadJSON(w, r, &requestData); err != nil {
			erp.BadRequestResponse(w, r, h.logger, err)
			return
		}

		err = h.productsUC.Update(id, requestData.Name, requestData.Tags)
		if err != nil {
			if errors.Is(err, db.ErrRecordNotFound) {
				erp.NotFoundResponse(w, r, h.logger)
			} else {
				erp.ServerErrorResponse(w, r, h.logger, err)
			}
			return
		}

		if requestData.Tags != nil {
			messagePayload := kf.KafkaMessage{
				Action: "product_update",
				Time:   time.Now().Format(time.RFC3339),
				Tags:   requestData.Tags,
			}

			err = h.productsUC.SendToKafka(r.Context(), strconv.Itoa(int(id)), messagePayload, h.kafkaProductWriter)
			if err != nil {
				erp.ServerErrorResponse(w, r, h.logger, err)
				return
			}
		}

		err = utils.WriteJSON(w, http.StatusOK, utils.Envelope{
			"message": "product updated successfully",
		}, nil)
		if err != nil {
			erp.ServerErrorResponse(w, r, h.logger, err)
		}
	}
}

//	@Summary		Delete a product
//	@Description	Deletes an existing product (admin-only).
//	@Tags			products
//	@Produce		json
//	@Security		cookieAuth
//	@Param			id	path		int						true	"Product ID"
//	@Success		200	{object}	models.SuccessResponse	"success response with product"
//	@Failure		400	{object}	models.ErrorResponse	"bad request error"
//	@Failure		404	{object}	models.ErrorResponse	"not found error"
//	@Failure		500	{object}	models.ErrorResponse	"internal server error"
//	@Router			/delete/{id} [delete]
func (h *productsHandlers) Delete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := utils.ReadIDParam(r)
		if err != nil {
			erp.BadRequestResponse(w, r, h.logger, err)
			return
		}

		err = h.productsUC.Delete(id)
		if err != nil {
			if errors.Is(err, db.ErrRecordNotFound) {
				erp.NotFoundResponse(w, r, h.logger)
			} else {
				erp.ServerErrorResponse(w, r, h.logger, err)
			}
			return
		}

		messagePayload := kf.KafkaMessage{
			Action: "product_delete",
			Time:   time.Now().Format(time.RFC3339),
			Tags:   nil,
		}

		err = h.productsUC.SendToKafka(r.Context(), strconv.Itoa(int(id)), messagePayload, h.kafkaProductWriter)
		if err != nil {
			erp.ServerErrorResponse(w, r, h.logger, err)
			return
		}

		err = utils.WriteJSON(w, http.StatusNoContent, utils.Envelope{
			"message": "product successfully deleted",
		}, nil)
		if err != nil {
			erp.ServerErrorResponse(w, r, h.logger, err)
		}
	}
}
