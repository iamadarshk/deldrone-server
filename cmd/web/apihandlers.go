// this file contains the handlers required for serving the api
// TODO: Add Authentication
package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/iamadarshk/deldrone-server/pkg/models"
)

// Vendor -----------------------------------------------------------------------------------------

// Method: GET, Path: "/api/vendor/{vendorID}"
// Feteches a particular vendor by their vendorid
func (app *application) apiGetVendorByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	vendorid, err := strconv.Atoi(vars["vendorID"])
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	vendor, err := app.vendors.Get(vendorid)
	if err != nil {
		app.serverError(w, err)
		return
	}
	jsonData, err := json.Marshal(vendor)
	if err != nil {
		app.serverError(w, err)
		return
	}
	app.writeJSON(w, r, jsonData)
}

// Method: GET, PATH: "/api/vendors/{pincode}/{pincoderange}"
// Fetches all vendors whith pincode = {pincode} +- {pincoderange}
func (app *application) apiGetVendorByPincode(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	pincode, err := strconv.Atoi(vars["pincode"])
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	pincodeRange, err := strconv.Atoi(vars["pincoderange"])
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	vendors, err := app.vendors.GetByPincode(pincode, pincodeRange)
	if err != nil {
		app.serverError(w, err)
		return
	}
	jsonData, err := json.Marshal(vendors)
	if err != nil {
		app.serverError(w, err)
		return
	}
	app.writeJSON(w, r, jsonData)
}

// Method: GET, Path: "/api/vendor/{vendorID}/listings"
// Fetches all listings by particular vendor using vendorID
func (app *application) apiGetVendorListings(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	vendorID, err := strconv.Atoi(vars["vendorID"])
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	listings, err := app.listings.All(vendorID)
	if err != nil {
		app.serverError(w, err)
		return
	}
	jsonData, err := json.Marshal(listings)
	if err != nil {
		app.serverError(w, err)
		return
	}
	app.writeJSON(w, r, jsonData)
}

// Method: PUT, Path: "/api/vendor/{vendorID}/listing/{listingID}"
// Updates particular listing using listingID and vendorID

// Method: DELETE, Path: "/api/vendor/{vendorID}/listing/{listingID}"
// Deletes particular listing using vendorID and listingID

// Method: POST, Path: "/api/vendor/{vendorID}/listing/new"
// Creates a new listing for vendor with their vendorID

// Method: GET, Path: "/api/vendor/{vendorID}/deliveries"
// fetches all deliveries from vendor with their vendorID
func (app *application) apiGetVendorDeliveries(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	vendorID, err := strconv.Atoi(vars["vendorID"])
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	deliveries, err := app.deliveries.GetAllByVendorID(vendorID)
	if err != nil {
		app.serverError(w, err)
		return
	}
	jsonData, err := json.Marshal(deliveries)
	if err != nil {
		app.serverError(w, err)
		return
	}
	app.writeJSON(w, r, jsonData)
}

// Method: GET, Path: "/api/vendor/{vendorID}/activedeliveries"
// fetches active deliveries for vendor using their vendorID
func (app *application) apiGetVendorActiveDeliveries(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	vendorID, err := strconv.Atoi(vars["vendorID"])
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	deliveries, err := app.deliveries.GetAllByVendorIDStatus(vendorID, "placed")
	if err != nil {
		app.serverError(w, err)
		return
	}
	jsonData, err := json.Marshal(deliveries)
	if err != nil {
		app.serverError(w, err)
		return
	}
	app.writeJSON(w, r, jsonData)
}

// Customer ---------------------------------------------------------------------------------------

// Method: GET, Path: "api/customer/{customerID}"
// fetches details for particular customer using their customerID
func (app *application) apiGetCustomer(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	customerID, err := strconv.Atoi(vars["customerID"])
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	customer, err := app.customers.Get(customerID)
	if err != nil {
		app.serverError(w, err)
		return
	}
	jsonData, err := json.Marshal(customer)
	if err != nil {
		app.serverError(w, err)
		return
	}
	app.writeJSON(w, r, jsonData)
}

// Method: GET, Path: "api/customer/{customerID}/getcart"
// fetches customer's cart using their customerID
func (app *application) apiGetCustomerCart(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	customerID, err := strconv.Atoi(vars["customerID"])
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	cart := app.carts[customerID]
	jsonData, err := json.Marshal(cart)
	if err != nil {
		app.serverError(w, err)
		return
	}
	app.writeJSON(w, r, jsonData)
}

// Method: POST, Path: "api/customer/{customerID}/cart/{listingID}"
// adds listing with listingID in customer's cart using their customerID
func (app *application) apiCustomerAddToCart(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	customerID, err := strconv.Atoi(vars["customerID"])
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	listingID, err := strconv.Atoi(vars["listingID"])
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	var quant int
	err = json.NewDecoder(r.Body).Decode(&quant)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	if quant <= 0 || listingID <= 0 {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	cart := app.carts[customerID]
	if cart == nil {
		cart = make(models.Cart)
		app.carts[customerID] = cart
	}
	cart.Add(listingID, quant)
	url := fmt.Sprintf("/api/customer/%d/cart", customerID)
	http.Redirect(w, r, url, http.StatusSeeOther)
}

// Method: PUT, Path: "api/customer/{customerID}/cart/{listingID}"
// updates listing with listingID in customer's cart using their customerID
func (app *application) apiCustomerUpdateCart(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	customerID, err := strconv.Atoi(vars["customerID"])
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	listingID, err := strconv.Atoi(vars["listingID"])
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	var quant int
	err = json.NewDecoder(r.Body).Decode(&quant)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	if quant <= 0 || listingID <= 0 {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	cart := app.carts[customerID]
	if cart == nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	cart.Add(listingID, quant)
	url := fmt.Sprintf("/api/customer/%d/cart", customerID)
	http.Redirect(w, r, url, http.StatusSeeOther)
}

// Method: DELETE, Path: "api/customer/{customerID}/cart/{listingID}"
// deletes listing with listingID in customer's cart using their customerID
func (app *application) apiCustomerDeleteFromCart(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	customerID, err := strconv.Atoi(vars["customerID"])
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	listingID, err := strconv.Atoi(vars["listingID"])
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	cart := app.carts[customerID]
	delete(cart, listingID)
	url := fmt.Sprintf("/api/customer/%d/cart", customerID)
	http.Redirect(w, r, url, http.StatusSeeOther)
}

// Method: POST, Path: "api/customer/{customerID}/checkout"
// checks out customer using their customerID
func (app *application) apiCheckout(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	customerID, err := strconv.Atoi(vars["customerID"])
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	cart := app.carts[customerID]
	var vendorID int
	for listingID := range cart {
		listing, err := app.listings.Get(listingID)
		if err != nil {
			app.serverError(w, err)
			return
		}
		if vendorID == 0 {
			vendorID = listing.VendorID
		} else if vendorID != listing.VendorID {
			app.clientError(w, http.StatusConflict)
			return
		}
	}
	t := struct {
		DropLat  float64 `json:"dropLatitude"`
		DropLong float64 `json:"dropLongitude"`
	}{}
	err = json.NewDecoder(r.Body).Decode(&t)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	deliveryID, err := app.deliveries.Insert(customerID, vendorID, time.Now(), t.DropLat, t.DropLong)
	if err != nil {
		app.serverError(w, err)
		return
	}
	for listingID, quant := range cart {
		listing, err := app.listings.Get(listingID)
		if err != nil {
			app.serverError(w, err)
			return
		}
		err = app.orders.Insert(deliveryID, listingID, quant, listing.Price*quant)
		if err != nil {
			app.serverError(w, err)
			return
		}
	}
	url := fmt.Sprintf("/api/customer/%d/activeorders", customerID)
	http.Redirect(w, r, url, http.StatusSeeOther)
}

// Method: GET, Path: "api/customer/{customerID}/deliveries"
// fetches all deliveries for customer using their customerID
func (app *application) apiGetCustomerDeliveries(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	customerID, err := strconv.Atoi(vars["customerID"])
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	deliveries, err := app.deliveries.GetAllByCustomerID(customerID)
	if err != nil {
		app.serverError(w, err)
		return
	}
	jsonData, err := json.Marshal(deliveries)
	if err != nil {
		app.serverError(w, err)
		return
	}
	app.writeJSON(w, r, jsonData)
}

// Method: GET, Path: "api/customer/{customerID}/activedeliveries"
// fetches active deliveries for customer using their customerID
func (app *application) apiGetCustomerActiveDeliveries(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	customerID, err := strconv.Atoi(vars["customerID"])
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	deliveries, err := app.deliveries.GetAllByCustomerIDStatus(customerID, "placed")
	if err != nil {
		app.serverError(w, err)
		return
	}
	jsonData, err := json.Marshal(deliveries)
	if err != nil {
		app.serverError(w, err)
		return
	}
	app.writeJSON(w, r, jsonData)
}

// Listings ---------------------------------------------------------------------------------------

// Method: GET, Path: "api/listing/{listingID}"
// Fetches a listing by it's listingID
func (app *application) apiGetListingByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	listingID, err := strconv.Atoi(vars["listingID"])
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	listing, err := app.listings.Get(listingID)
	if err != nil {
		app.serverError(w, err)
		return
	}

	jsonData, err := json.Marshal(listing)
	if err != nil {
		app.serverError(w, err)
		return
	}
	app.writeJSON(w, r, jsonData)
}

// Deliveries -------------------------------------------------------------------------------------

// Method: GET, Path: "api/delivery/{deliveryID}"
// fetch details of a delivery using deliveryID (NOT LIST OF ITEMS ASSOCIATED WITH THE DELIVERYID)
func (app *application) apiGetDelivery(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	deliveryID, err := strconv.Atoi(vars["deliveryID"])
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	delivery, err := app.deliveries.Get(deliveryID)
	if err != nil {
		app.serverError(w, err)
		return
	}

	jsonData, err := json.Marshal(delivery)
	if err != nil {
		app.serverError(w, err)
		return
	}
	app.writeJSON(w, r, jsonData)
}

// Method: GET, Path: "api/delivery/{deliveryID}/orders"
// fetches all orders associated with a deliveryID
func (app *application) apiGetOrdersFromDelivery(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	deliveryID, err := strconv.Atoi(vars["deliveryID"])
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	orders, err := app.orders.AllFromDeliveryID(deliveryID)
	if err != nil {
		app.serverError(w, err)
		return
	}

	jsonData, err := json.Marshal(orders)
	if err != nil {
		app.serverError(w, err)
		return
	}
	app.writeJSON(w, r, jsonData)
}

// Orders -----------------------------------------------------------------------------------------

// Method: GET, Path: "api/order/{orderID}"
// fetches an order using orderID
func (app *application) apiGetOrder(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	orderID, err := strconv.Atoi(vars["orderID"])
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	order, err := app.orders.Get(orderID)
	if err != nil {
		app.serverError(w, err)
		return
	}

	jsonData, err := json.Marshal(order)
	if err != nil {
		app.serverError(w, err)
		return
	}
	app.writeJSON(w, r, jsonData)
}
