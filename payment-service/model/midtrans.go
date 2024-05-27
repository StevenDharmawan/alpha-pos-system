package model

import "github.com/midtrans/midtrans-go"

type PaymentRequest struct {
	UserEmail   string                 `json:"user_email"`
	OrderID     string                 `json:"order_id"`
	GrossAmt    int64                  `json:"gross_amt"`
	ItemDetails []midtrans.ItemDetails `json:"item_details"`
}
