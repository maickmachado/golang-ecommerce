package database

import "errors"

var (
	ErrCantFindProduct    = errors.New("can't find product")
	ErrCantDecodeProduct  = errors.New("product not available")
	ErrUserIdIsNotValid   = errors.New("this user is not a valid user")
	ErrCantUpdateUser     = errors.New("can't add product to cart")
	ErrCantRemoveItemCart = errors.New("can't remove item from cart")
	ErrCantGetItemCart    = errors.New("can't get item from cart")
	ErrCantBuyItemCart    = errors.New("can't buy item from cart")
)

func AddProductToCart() {}

func RemoveProductFromCart() {}

func BuyFromCart() {}

func InstantBuyFromCart() {}
