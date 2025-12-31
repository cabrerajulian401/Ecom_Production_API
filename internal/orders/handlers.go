package orders

import (
	"log"
	"net/http"

	"github.com/cabrerajulian401/ecom/internal/writer"
)

type handler struct {
	service Service
}

func NewHandler(service Service) *handler {
	return &handler{
		service: service,
	}
}

// HTTP Handlers do not have return values so to exit in an error just write an empty `return`
func (h *handler) PlaceOrder(w http.ResponseWriter, r *http.Request) {

	var tempOrder createOrderParams
	if err := writer.Read(r, &tempOrder); err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	createdOrder, err := h.service.PlaceOrder(r.Context(), tempOrder)
	if err != nil {
		log.Println(err)

		if err == ErrProductNotFound {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writer.Write(w, http.StatusCreated, createdOrder)
}
