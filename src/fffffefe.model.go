/**
 *
 * Copyright (c) 2020, Oracle and/or its affiliates. All rights reserved.
 *
 */
package src

import (
	"example.com/fffffefe/lib/util/date"
)

type Bank_details struct {
	AssetType string `json:"AssetType" final:"fffffefe.Bank_details"`

	RawMaterialAvailable int         `json:"RawMaterialAvailable" validate:"int,min=0"`
	License              string      `json:"License" validate:"string,min=2,max=4"`
	ExpiryDate           date.Date   `json:"ExpiryDate" validate:"date,before=2020-06-26"`
	Active               bool        `json:"Active" validate:"bool" default:"true"`
	Metadata             interface{} `json:"Metadata,omitempty"`
}

type Customer struct {
	AssetType string `json:"AssetType" final:"fffffefe.Customer"`

	CustomerId     string       `json:"CustomerId" validate:"string" id:"true" mandatory:"true"`
	Name           string       `json:"Name" validate:"string" mandatory:"true"`
	ProductsBought int          `json:"ProductsBought" validate:"int"`
	OfferApplied   int          `json:"OfferApplied" validate:"int,max=0"`
	PhoneNumber    string       `json:"PhoneNumber" validate:"string,regexp=^\\(?([0-9]{3})\\)?[-. ]?([0-9]{3})[-. ]?([0-9]{4})$"`
	Received       bool         `json:"Received" validate:"bool"`
	Bank_details   Bank_details `json:"Bank_details" validate:""`
	Metadata       interface{}  `json:"Metadata,omitempty"`
}

type Retailer struct {
	AssetType string `json:"AssetType" final:"fffffefe.Retailer"`

	RetailerId        string      `json:"RetailerId" validate:"string" id:"true" mandatory:"true"`
	Customer          Customer    `json:"Customer" validate:""`
	ProductsOrdered   int         `json:"ProductsOrdered" validate:"int" mandatory:"true"`
	ProductsAvailable int         `json:"ProductsAvailable" validate:"int" default:"1"`
	ProductsSold      int         `json:"ProductsSold" validate:"int"`
	Remarks           string      `json:"Remarks" validate:"string" default:"open for business"`
	Items             []int       `json:"Items" validate:"array=int,range=1-5"`
	Domain            string      `json:"Domain" validate:"string,url,min=30,max=50"`
	Metadata          interface{} `json:"Metadata,omitempty"`
}

type Account struct {
	AssetType string `json:"AssetType" final:"fffffefe.Account"`

	RawMaterialAvailable int         `json:"RawMaterialAvailable" validate:"int,min=0"`
	License              string      `json:"License" validate:"string,min=2,max=4"`
	ExpiryDate           date.Date   `json:"ExpiryDate" validate:"date,before=2020-06-26"`
	Active               bool        `json:"Active" validate:"bool" default:"true"`
	Metadata             interface{} `json:"Metadata,omitempty"`
}

type Supplier struct {
	AssetType string `json:"AssetType" final:"fffffefe.Supplier"`

	SupplierId           string      `json:"SupplierId" validate:"string,regexp=^[a-zA-Z]$" id:"true" mandatory:"true"`
	Retailer             Retailer    `json:"Retailer" validate:""`
	RawMaterialAvailable int         `json:"RawMaterialAvailable" validate:"int,min=0"`
	License              string      `json:"License" validate:"string,min=2,max=4"`
	ExpiryDate           date.Date   `json:"ExpiryDate" validate:"date,before=2020-06-26"`
	Active               bool        `json:"Active" validate:"bool" default:"true"`
	Account              Account     `json:"Account" validate:""`
	Metadata             interface{} `json:"Metadata,omitempty"`
}

type Manufacturer struct {
	AssetType string `json:"AssetType" final:"fffffefe.Manufacturer"`

	ManufacturerId       string       `json:"ManufacturerId" validate:"string" id:"true" mandatory:"true"`
	Bank_details         Bank_details `json:"Bank_details" validate:""`
	RawMaterialAvailable int          `json:"RawMaterialAvailable" validate:"int,max=8"`
	ProductsAvailable    int          `json:"ProductsAvailable" validate:"int"`
	CompletionDate       date.Date    `json:"CompletionDate" validate:"date,after=2020-06-26T02:30:55Z,before=2020-06-28T02:30:55Z"`
	Account              Account      `json:"Account" validate:""`
	Metadata             interface{}  `json:"Metadata,omitempty"`
}

type Distributor struct {
	AssetType string `json:"AssetType" final:"fffffefe.Distributor"`

	DistributorId       string      `json:"DistributorId" validate:"string" id:"true" mandatory:"true"`
	ProductsToBeShipped int         `json:"ProductsToBeShipped" validate:"int"`
	ProductsShipped     int         `json:"ProductsShipped" validate:"int,min=3"`
	ProductsReceived    int         `json:"ProductsReceived" validate:"int"`
	MailId              string      `json:"MailId" validate:"string,email"`
	DistributionDate    date.Date   `json:"DistributionDate" validate:"date"`
	Metadata            interface{} `json:"Metadata,omitempty"`
}
