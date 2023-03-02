package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type Product struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type productList struct {
	Products []Product `json:"products"`
}

var products = productList{
	Products: []Product{
		{ID: 1, Name: "Product 1", Description: "Description for Product 1"},
		{ID: 2, Name: "Product 2", Description: "Description for Product 2"},
	},
}

func getProducts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(products)
}

func getProduct(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	id := getIdFromRequest(r)
	for _, product := range products.Products {
		if product.ID == id {
			json.NewEncoder(w).Encode(product)
			return
		}
	}
	http.NotFound(w, r)
}

func addProduct(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var product Product
	err := json.NewDecoder(r.Body).Decode(&product)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	product.ID = getNextId()
	products.Products = append(products.Products, product)
	json.NewEncoder(w).Encode(product)
}

func updateProduct(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	id := getIdFromRequest(r)
	for index, product := range products.Products {
		if product.ID == id {
			var updatedProduct Product
			err := json.NewDecoder(r.Body).Decode(&updatedProduct)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			updatedProduct.ID = id
			products.Products[index] = updatedProduct
			json.NewEncoder(w).Encode(updatedProduct)
			return
		}
	}
	http.NotFound(w, r)
}

func deleteProduct(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	id := getIdFromRequest(r)
	for index, product := range products.Products {
		if product.ID == id {
			products.Products = append(products.Products[:index], products.Products[index+1:]...)
			w.WriteHeader(http.StatusNoContent)
			return
		}
	}
	http.NotFound(w, r)
}

func getIdFromRequest(r *http.Request) int {
	id := 0
	if idParam, ok := r.URL.Query()["id"]; ok {
		fmt.Sscanf(idParam[0], "%d", &id)
	}
	return id
}

func getNextId() int {
	return len(products.Products) + 1
}

func main() {
	http.HandleFunc("/products", getProducts)
	http.HandleFunc("/product", getProduct)
	http.HandleFunc("/add-product", addProduct)
	http.HandleFunc("/update-product", updateProduct)
	http.HandleFunc("/delete-product", deleteProduct)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
