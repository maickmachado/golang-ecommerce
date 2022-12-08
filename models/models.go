package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserModel struct {
	ID             primitive.ObjectID `json:"_id" bson:"_id"`
	FistName       *string            `json:"fist_name" validate:"required,min=2,max=30"`
	LastName       *string            `json:"last_name" validate:"required,min=2,max=30"`
	Email          *string            `json:"email" validate:"email,required"`
	Password       *string            `json:"password" validate:"required,min=6"`
	PhoneNumber    *string            `json:"phone_number" validate:"required"`
	Token          *string            `json:"token"`
	RefreshToken   *string            `json:"refresh_token"`
	CreatedAt      time.Time          `json:"created_at"`
	UpdatedAt      time.Time          `json:"updated_at"`
	UserID         string             `json:"user_id"`
	UserCart       []UserCart         `json:"user_cart" bson:"user_cart"`
	AddressDetails []AdressDetails    `json:"address" bson:"address"`
	Order          []Order            `json:"order" bson:"order"`
}

type Product struct {
	ProductID       primitive.ObjectID `bson:"_id"`
	ProductName     *string            `json:"product_name"`
	ProductPrice    *uint64            `json:"product_price"`
	ProductRating   *uint8             `json:"product_rating"`
	ProductImage    *string            `json:"product_image"`
	ProductQuantity *uint8             `json:"product_quantity"`
}

type UserCart struct {
	UserCartProductID       primitive.ObjectID `bson:"_id"`
	UserCartProductName     *string            `json:"product_name"`
	UserCartProductPrice    *uint64            `json:"product_price"`
	UserCartProductRating   *uint8             `json:"product_rating"`
	UserCartProductImage    string             `json:"product_image"`
	UserCartProductQuantity *uint8             `json:"product_quantity"`
}

type AdressDetails struct {
	AddressID primitive.ObjectID `bson:"_id"`
	House     *string            `json:"house"`
	Street    *string            `json:"street"`
	City      *string            `json:"city"`
	State     *string            `json:"state"`
	ZipCode   *string            `json:"zip_code"`
	Country   *string            `json:"country"`
}

type Order struct {
	OrderID            primitive.ObjectID `bson:"_id"`
	OrderStatus        string             `json:"order_status"`
	OrderCart          []UserCart         `json:"order_cart" bson:"order_cart"`
	OrderDate          time.Time          `json:"order_date" bson:"order_date"`
	OrderTotal         int                `json:"order_total" bson:"order_total"`
	OrderDiscount      *int               `json:"order_discount" bson:"discount"`
	OrderPaymentMethod OrderPaymentMethod `json:"order_payment_method" bson:"order_payment_method"`
}

type OrderPaymentMethod struct {
	DigitalPaymentMethod bool `json:"digital_payment_method" bson:"digital_payment_method"`
	PixPaymentMethod     bool `json:"pix_payment_method" bson:"pix_payment_method"`
}
