package handlers

import (
	"backend/internal/entities"
	"backend/internal/usecases/storage/product"
	"backend/logger"
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

const (
	productID = "product_id"
)

func getAllProducts(app product.Controller) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.Info.Println("Got new request to get all products")
		errorMessage := "Error getting all products"
		result, err := app.GetProducts()
		if err != nil {
			logger.Error.Println(err.Error())
			http.Error(w, errorMessage, http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(result)
	})
}

func getProduct(app product.Controller) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.Info.Printf("Got new request to get product with id %s", mux.Vars(r)[productID])

		errorMessage := "Error getting product"
		productID, err := strconv.Atoi(mux.Vars(r)[productID])
		if err != nil {
			logger.Error.Printf("Failed to parse product ID. Got %v", err)
			http.Error(w, errorMessage, http.StatusBadRequest)
			return
		}
		result, err := app.GetProduct(productID)
		if err != nil {
			logger.Error.Printf(err.Error())
			http.Error(w, errorMessage, http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(result)
	})
}

func addProduct(app product.Controller) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.Info.Println("Got new request to add product")

		errorMessage := "Error adding new product"

		newProduct := new(entities.Product)
		err := json.NewDecoder(r.Body).Decode(&newProduct)
		if err != nil {
			logger.Error.Println("Failed to decode product. Got %v", err)
			http.Error(w, errorMessage, http.StatusBadRequest)
			return
		}
		err = app.AddProduct(*newProduct)
		if err != nil {
			logger.Error.Println(err.Error())
			http.Error(w, errorMessage, http.StatusInternalServerError)
			return
		}
	})
}

func Make(r *mux.Router, app product.Controller) {
	apiURI := "/api"
	serviceRouter := r.PathPrefix(apiURI).Subrouter()
	serviceRouter.Handle("/products", getAllProducts(app)).Methods("GET")
	serviceRouter.Handle("/products/{product_id}", getProduct(app)).Methods("GET")
	serviceRouter.Handle("/products", addProduct(app)).Methods("POST")
}
