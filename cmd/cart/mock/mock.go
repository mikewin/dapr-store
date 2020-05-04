package mock

import (
	cartspec "github.com/benc-uk/dapr-store/cmd/cart/spec"
	orderspec "github.com/benc-uk/dapr-store/cmd/orders/spec"
	productspec "github.com/benc-uk/dapr-store/cmd/products/spec"
	"github.com/benc-uk/dapr-store/pkg/problem"
)

// CartService mock
type CartService struct {
}

var mockCart = &cartspec.Cart{
	Products: map[string]int{},
	ForUser:  "demo@example.net",
}

//
// Get fetches saved cart for a given user, if not exists an empty cart is returned
//
func (s CartService) Get(username string) (*cartspec.Cart, error) {
	if username != "demo@example.net" {
		cart := &cartspec.Cart{}
		cart.ForUser = username
		cart.Products = make(map[string]int)
		return cart, nil
	}
	return mockCart, nil
}

//
// Submit a cart and turn into an order
//
func (s CartService) Submit(cart cartspec.Cart) (*orderspec.Order, error) {
	if len(cart.Products) == 0 {
		return nil, problem.New("err://bad", "Cart empty", 400, "Cart empty", "mock-cart")
	}

	o := &orderspec.Order{
		Title:   "Mock Order",
		Amount:  12.34,
		ForUser: cart.ForUser,
		ID:      "order-01",
		Status:  orderspec.OrderNew,
		LineItems: []orderspec.LineItem{
			{
				Count: 1,
				Product: productspec.Product{
					ID:          "4",
					Name:        "foo",
					Cost:        12.34,
					Description: "blah",
					Image:       "blah.jpg",
					OnOffer:     false,
				},
			},
		},
	}
	return o, nil
}

//
// SetProductCount updates the count of a given product in the cart
//
func (s CartService) SetProductCount(cart *cartspec.Cart, productID string, count int) error {
	if count < 0 {
		return problem.New("err://bad", "SetProductCount", 500, "count can not be negative", "mock-cart")
	}
	if count == 0 {
		delete(mockCart.Products, productID)
		return nil
	}
	mockCart.Products[productID] = count
	return nil
}

//
// Clear the cart
//
func (s CartService) Clear(cart *cartspec.Cart) error {
	cart.Products = map[string]int{}
	return nil
}