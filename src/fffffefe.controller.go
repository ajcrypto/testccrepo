/**
 *
 * Copyright (c) 2020, Oracle and/or its affiliates. All rights reserved.
 *
 */
package src

import (
	"example.com/fffffefe/lib/model"
)

type Controller struct {
}

func (t *Controller) Init() (interface{}, error) {
	return nil, nil
}

//-----------------------------------------------------------------------------
//Bank_details
//-----------------------------------------------------------------------------

func (t *Controller) CreateBank_details(asset Bank_details) (interface{}, error) {
	return model.Save(&asset)
}

func (t *Controller) GetBank_detailsById(id string) (Bank_details, error) {
	var asset Bank_details
	_, err := model.Get(id, &asset)
	return asset, err
}

func (t *Controller) UpdateBank_details(asset Bank_details) (interface{}, error) {
	return model.Update(&asset)
}

func (t *Controller) DeleteBank_details(id string) (interface{}, error) {
	return model.Delete(id)
}

func (t *Controller) GetBank_detailsHistoryById(id string) (interface{}, error) {
	historyArray, err := model.GetHistoryByID(id)
	return historyArray, err
}

func (t *Controller) GetBank_detailsByRange(startkey string, endKey string) ([]Bank_details, error) {
	var assets []Bank_details
	_, err := model.GetByRange(startkey, endKey, &assets)
	return assets, err
}

//-----------------------------------------------------------------------------
//Customer
//-----------------------------------------------------------------------------

func (t *Controller) CreateCustomer(asset Customer) (interface{}, error) {
	return model.Save(&asset)
}

func (t *Controller) GetCustomerById(id string) (Customer, error) {
	var asset Customer
	_, err := model.Get(id, &asset)
	return asset, err
}

//-----------------------------------------------------------------------------
//Retailer
//-----------------------------------------------------------------------------

func (t *Controller) CreateRetailer(asset Retailer) (interface{}, error) {
	return model.Save(&asset)
}

func (t *Controller) GetRetailerById(id string) (Retailer, error) {
	var asset Retailer
	_, err := model.Get(id, &asset)
	return asset, err
}

//-----------------------------------------------------------------------------
//Account
//-----------------------------------------------------------------------------

func (t *Controller) CreateAccount(asset Account) (interface{}, error) {
	return model.Save(&asset)
}

func (t *Controller) GetAccountById(id string) (Account, error) {
	var asset Account
	_, err := model.Get(id, &asset)
	return asset, err
}

func (t *Controller) UpdateAccount(asset Account) (interface{}, error) {
	return model.Update(&asset)
}

func (t *Controller) DeleteAccount(id string) (interface{}, error) {
	return model.Delete(id)
}

func (t *Controller) GetAccountHistoryById(id string) (interface{}, error) {
	historyArray, err := model.GetHistoryByID(id)
	return historyArray, err
}

func (t *Controller) GetAccountByRange(startkey string, endKey string) ([]Account, error) {
	var assets []Account
	_, err := model.GetByRange(startkey, endKey, &assets)
	return assets, err
}

//-----------------------------------------------------------------------------
//Supplier
//-----------------------------------------------------------------------------

func (t *Controller) CreateSupplier(asset Supplier) (interface{}, error) {
	return model.Save(&asset)
}

func (t *Controller) GetSupplierById(id string) (Supplier, error) {
	var asset Supplier
	_, err := model.Get(id, &asset)
	return asset, err
}

func (t *Controller) UpdateSupplier(asset Supplier) (interface{}, error) {
	return model.Update(&asset)
}

func (t *Controller) DeleteSupplier(id string) (interface{}, error) {
	return model.Delete(id)
}

func (t *Controller) GetSupplierHistoryById(id string) (interface{}, error) {
	historyArray, err := model.GetHistoryByID(id)
	return historyArray, err
}

func (t *Controller) GetSupplierByRange(startkey string, endKey string) ([]Supplier, error) {
	var assets []Supplier
	_, err := model.GetByRange(startkey, endKey, &assets)
	return assets, err
}

//-----------------------------------------------------------------------------
//Manufacturer
//-----------------------------------------------------------------------------

func (t *Controller) CreateManufacturer(asset Manufacturer) (interface{}, error) {
	return model.Save(&asset)
}

func (t *Controller) GetManufacturerById(id string) (Manufacturer, error) {
	var asset Manufacturer
	_, err := model.Get(id, &asset)
	return asset, err
}

func (t *Controller) GetManufacturerHistoryById(id string) (interface{}, error) {
	historyArray, err := model.GetHistoryByID(id)
	return historyArray, err
}

//-----------------------------------------------------------------------------
//Distributor
//-----------------------------------------------------------------------------

func (t *Controller) CreateDistributor(asset Distributor) (interface{}, error) {
	return model.Save(&asset)
}

func (t *Controller) GetDistributorById(id string) (Distributor, error) {
	var asset Distributor
	_, err := model.Get(id, &asset)
	return asset, err
}

//-----------------------------------------------------------------------------
//Custom Methods
//-----------------------------------------------------------------------------

/**
 *
 * BDB sql rich queries can be executed in OBP CS/EE.
 * This method can be invoked only when connected to remote OBP CS/EE network.
 *
 */
func (t *Controller) ExecuteQuery(inputQuery string) (interface{}, error) {
	resultArray, err := model.Query(inputQuery)
	return resultArray, err
}

func (t *Controller) fetchRawMaterial(supplierId string, rawMaterialSupply int) (interface{}, error) {

	return nil, nil
}

func (t *Controller) getRawMaterialFromSupplier(manufacturerId string, supplierId string, rawMaterialSupply int) (interface{}, error) {

	return nil, nil
}

func (t *Controller) createProducts(manufacturerId string, rawMaterialConsumed int, productsCreated int) (interface{}, error) {

	return nil, nil
}

func (t *Controller) sendProductsToDistribution() (interface{}, error) {

	return nil, nil
}

func (t *Controller) someFunc() (interface{}, error) {

	return nil, nil
}

func (t *Controller) someFunc2() (interface{}, error) {

	return nil, nil
}
