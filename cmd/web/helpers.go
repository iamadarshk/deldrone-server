package main

import (
	"bytes"
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/iamadarshk/deldrone-server/pkg/models"

	"github.com/justinas/nosurf"
)

type cartRow struct {
	ListingID int
	Name      string
	Price     int
	Quantity  int
	Amount    int
}

func (app *application) serverError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.errorLog.Println(trace)

	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (app *application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

func (app *application) addDefaultData(td *templateData, w http.ResponseWriter, r *http.Request) (*templateData, error) {
	// get session
	session, err := app.sessionStore.Get(r, "session-name")
	if err != nil {
		return nil, err // TODO: handle cookie error
	}

	if td == nil {
		td = &templateData{}
	}

	// add flash message
	flashes := session.Flashes()
	if len(flashes) > 0 {
		td.Flash = flashes[0].(string)
	}

	td.CSRFToken = nosurf.Token(r)
	td.AuthenticatedCustomer = app.authenticatedCustomer(r)
	td.AuthenticatedVendor = app.authenticatedVendor(r)

	// save session
	err = session.Save(r, w)
	if err != nil {
		return nil, err
	}

	return td, nil
}

// Sets Header for JSON and writes bytes to w
func (app *application) writeJSON(w http.ResponseWriter, r *http.Request, jsonData []byte) {
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}

// renders the web page
func (app *application) render(w http.ResponseWriter, r *http.Request, name string, td *templateData) {
	// creates a template set
	ts, ok := app.templateCache[name]
	if !ok {
		app.serverError(w, fmt.Errorf("Template not found: %s", name))
		return
	}

	buf := new(bytes.Buffer)
	td, err := app.addDefaultData(td, w, r)
	if err != nil {
		app.serverError(w, err)
		return
	}
	err = ts.Execute(buf, td)
	if err != nil {
		app.serverError(w, err)
		return
	}
	buf.WriteTo(w)
}

func (app *application) authenticatedVendor(r *http.Request) int {
	session, err := app.sessionStore.Get(r, "session-name")
	if err != nil {
		app.errorLog.Println(err.Error())
		return 0
	}
	if session.Values["vendorID"] == nil {
		return 0
	}
	return session.Values["vendorID"].(int)
}

func (app *application) authenticatedCustomer(r *http.Request) int {
	session, err := app.sessionStore.Get(r, "session-name")
	if err != nil {
		app.errorLog.Println(err.Error())
		return 0
	}
	if session.Values["customerID"] == nil {
		return 0
	}
	return session.Values["customerID"].(int)
}

func (app *application) order2CartRow(orders []*models.Order) ([]*cartRow, error) {
	cartRows := make([]*cartRow, 0)
	for _, order := range orders {
		listing, err := app.listings.Get(order.ListingID)
		if err != nil {
			return nil, err
		}
		row := &cartRow{
			ListingID: order.ListingID,
			Amount:    order.Amount,
			Name:      listing.Name,
			Price:     listing.Price,
			Quantity:  order.Quant,
		}
		cartRows = append(cartRows, row)
	}
	return cartRows, nil
}
