package controllers

import (
	
	models "assignment-4/Models"
	"errors"
	"time"

	"github.com/go-pg/pg/v10"
)

func GetBanks(db *pg.DB) ([]*models.Bank, error) {

	var banks []*models.Bank

	if selectErr := db.Model(&banks).
		Relation("Branch").
		Order("id ASC").
		Select(); selectErr != nil {
		return banks, selectErr
	}
	return banks, nil
}

func GetBankById(db *pg.DB, bankId uint) (models.Bank, error) {
	var bank models.Bank

	if selectErr := db.Model(&bank).
		Relation("Branch").
		Where("bank.Id = ?", bankId).
		Select(); selectErr != nil {
		return bank, selectErr
	}

	return bank, nil
}

func UpdateBank(db *pg.DB, bank *models.Bank, bankID uint) (uint, error) {

	if _, updateErr := db.Model(bank).Where("id = ?", bankID).UpdateNotZero(); updateErr != nil {
		return bankID, updateErr
	}

	return bankID, nil
}

func CreateNewBank(db *pg.DB, newBank *models.Bank) (uint, error) {

	if _, err := db.Model(newBank).Insert(); err != nil {
		return 0, err
	}
	return newBank.ID, nil
}

func DeleteExistingBank(tx *pg.Tx, bankId uint) error {

	// Delete the bank record
	if _, deleteErr := tx.Model((*models.Bank)(nil)).Where("id = ?", bankId).Delete(); deleteErr != nil {
		return deleteErr
	}
	return nil
}

func GetBranches(db *pg.DB) ([]*models.Branch, error) {

	var branches []*models.Branch
	if selectErr := db.Model(&branches).Relation("Bank").Select(); selectErr != nil {
		return branches, selectErr
	}
	return branches, nil
}

func GetBranchById(db *pg.DB, branchId uint) (*models.Branch, error) {
	branch := &models.Branch{}
	if selectErr := db.Model(branch).
		Relation("Bank").
		Where("branch.id=?", branchId). // Specify 'branch.id'
		Select(); selectErr != nil {
		return branch, selectErr
	}

	return branch, nil
}

func UpdateBranchRec(tx *pg.Tx, branchId uint, branch *models.Branch) (uint, error) {

	res, updateErr := tx.Model(branch).Where("id = ?", branchId).UpdateNotZero(branch)

	if updateErr != nil {
		tx.Rollback()
		return branchId, updateErr
	}

	if res.RowsAffected() == 0 {
		tx.Rollback()
		return branchId, errors.New("no record updated")
	}
	return branchId, nil
}

func DeleteBranchById(tx *pg.Tx, branchId uint) error {
	res, deleteErr := tx.Model((*models.Branch)(nil)).
		Where("id = ?", branchId).
		Delete()

	if deleteErr != nil {
		tx.Rollback()
		return deleteErr
	}
	if res.RowsAffected() == 0 {
		tx.Rollback()
		return errors.New("no record deleted")
	}

	return nil
}

func GetCustDetails(db *pg.DB, custId uint) (models.Customer, error) {

	var customer models.Customer
	if selErr := db.Model(&customer).Where("cust_id = ?", custId).Relation("Account").Select(); selErr != nil {
		return models.Customer{}, selErr
	}
	return customer, nil

}

func updateCustDetails(tx *pg.Tx, customer models.Customer)(error){

	res, updateErr := tx.Model(&customer).WherePK().UpdateNotZero()

	if updateErr != nil {
		tx.Rollback()
		return updateErr
	}

	if res.RowsAffected() == 0 {
		tx.Rollback()
		return errors.New("no record updated")
	}
	return nil

}
func ViewTransDetails(db *pg.DB, transId uint) (models.Transaction, error) {

	//retrieve the transaction details
	var fetchedTrans models.Transaction

	if selErr := db.Model(&fetchedTrans).Relation("Account").Where("id=?", transId).Select(); selErr != nil {
		return models.Transaction{}, selErr
	}
	return fetchedTrans, nil
}

func SearchTrans(db *pg.DB, startDate time.Time) ([]models.Transaction, error) {
	var transactions []models.Transaction

	if selErr := db.Model(&transactions).Where("tr_date >= ? AND tr_date < ?", startDate, startDate.AddDate(0, 0, 1)).Select(); selErr != nil {
		return nil, selErr
	}
	return transactions, nil
}

func GetAccDetails(db *pg.DB, accId uint) (models.Account, error) {
	// Retrieve the account details along with its branch and bank
	var account models.Account

	selectErr := db.Model(&account).
		Relation("Branch.Bank").
		Relation("Customers").
		Where("acc_id = ?", accId).
		Select()
	if selectErr != nil {
		return models.Account{}, selectErr
	}
	return account, nil
}

func GetTransDetailsOfAcc(db *pg.DB, accId uint) (models.Account, error) {

	//retrieve all the transactions related to that account
	var account models.Account
	if selErr := db.Model(&account).Relation("Transaction").Where("acc_id=?", accId).Select(); selErr != nil {
		return models.Account{}, selErr
	}
	return account, nil
}

func SaveCustomers(tx *pg.Tx, inputDetails OpenCustAcc) ([]uint, error){

	var custIds []uint
	//insert record in customers table
	for _, customer := range inputDetails.Customers{
		//check if the customer details are already present in customers table, if yes then only retrieve
		//the custId else store the customer details
			existingCustomer := &models.Customer{}
			err := tx.Model(existingCustomer).
				Where("pan = ?", customer.PAN).
				Where("branch_id = ?", inputDetails.BranchID).
				Select()

			if err == nil {
				custIds = append(custIds, existingCustomer.CustId) //customer already exists
			} else if err == pg.ErrNoRows {

				//create customer instance
				cust := &models.Customer{
					Name:     customer.Name,
					PAN:      customer.PAN,
					DOB:      customer.DOB,
					Phone:    customer.Phone,
					Address:  customer.Address,
					BranchID: inputDetails.BranchID,
				}
				//insert record in customers table
				_, custErr := tx.Model(cust).Insert()

				if custErr != nil {	
					tx.Rollback()				
					return []uint{} ,custErr
				}
				//this slice will be used for mapping accounts with customer
				custIds = append(custIds, cust.CustId)
			} else {
				// Handle other errors
				tx.Rollback()
				return []uint{} , errors.New("failed to check existing customer")
			}
	}
	return custIds,nil
}

func SaveAccount(tx *pg.Tx, account *models.Account)(uint,error){

	if _, accErr := tx.Model(account).Insert(); accErr != nil {
		tx.Rollback()
		return 0, errors.New("failed to insert record in accounts table")
	}
	return account.AccID,nil
}

func SaveCustAcc(tx *pg.Tx, custIds []uint, accId uint) error {

	for _, custId := range custIds {

		_, insertErr := tx.Model(&models.CustomerToAccount{
									AccId:  accId,
									CustId: custId,
								}).Insert()

		if insertErr != nil {
			tx.Rollback()
			return errors.New("failed to insert record in accounts table")
		}
	}
	return nil
}

func GetCustAcc(tx *pg.Tx, accId uint)([] models.CustomerToAccount,error){

	var accCust []models.CustomerToAccount
	if selErr := tx.Model(&accCust).Where("acc_id=?", accId).Select(); selErr != nil{
		
		tx.Rollback()
		return []models.CustomerToAccount{},selErr
	}
	return accCust,nil
}

func DeleteCust(tx *pg.Tx, accCust []models.CustomerToAccount)error{

	//iterating over this slice to delete corresponding customer details
	for _, customer := range accCust {

		count, err := tx.Model((*models.CustomerToAccount)(nil)).
			Where("cust_id = ?", customer.CustId).
			Count()
		if err != nil {
			tx.Rollback()
			return errors.New("failed to retrieve count of customer from mapping table")
		}
		if count == 1 {
			//delete the corresponding customer record
			if _, delErr := tx.Model((*models.Customer)(nil)).Where("cust_id=?", customer.CustId).Delete(); delErr != nil {
				tx.Rollback()
				return errors.New("failed to delete customer data from customers table")
			}

		} 
	}

	return nil
}

func DeleteAcc(tx *pg.Tx, accId uint) error {
	if _, delErr := tx.Model((*models.Account)(nil)).Where("acc_id=?", accId).Delete(); delErr != nil{
		tx.Rollback()
		return delErr
	}
	return nil
}

func InsertTransaction(tx *pg.Tx, transaction models.Transaction) error {

	if _, err := tx.Model(&transaction).Insert(); err != nil {
		tx.Rollback()
		return err
	}
	return nil
}

func UpdateAcc(tx *pg.Tx, account models.Account) error {
	
	if _, err := tx.Model(&account).Where("acc_id = ?", account.AccID).Update(); err != nil {
		tx.Rollback()
		return err
	}
	return nil
}


