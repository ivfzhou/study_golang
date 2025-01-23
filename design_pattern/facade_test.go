// 外观模式是一种结构型设计模式，能为程序库、框架或其他复杂类提供一个简单的接口。
// 如果你需要一个指向复杂子系统的直接接口，且该接口的功能有限，则可以使用外观模式。
// 如果需要将子系统组织为多层结构，可以使用外观。
package design_pattern_test

import (
	"fmt"
	"testing"
)

/*
钱包功能，结合了账户等其他功能。
*/

func TestFacade(t *testing.T) {
	walletFacade := NewWalletFacade("abc", 1234)

	err := walletFacade.AddMoneyToWallet("abc", 1234, 10)
	if err != nil {
		t.Fatalf("Error: %s\n", err.Error())
	}

	err = walletFacade.DeductMoneyFromWallet("abc", 1234, 5)
	if err != nil {
		t.Fatalf("Error: %s\n", err.Error())
	}
}

// ===

type WalletFacade struct {
	account      *Account
	wallet       *Wallet
	securityCode *SecurityCode
	notification *Notification
	ledger       *Ledger
}

func NewWalletFacade(accountID string, code int) *WalletFacade {
	fmt.Println("Starting create account")
	walletFacade := &WalletFacade{
		account:      NewAccount(accountID),
		securityCode: NewSecurityCode(code),
		wallet:       NewWallet(),
		notification: &Notification{},
		ledger:       &Ledger{},
	}
	fmt.Println("Account created")
	return walletFacade
}

func (w *WalletFacade) AddMoneyToWallet(accountID string, securityCode int, amount int) error {
	fmt.Println("Starting Add money to wallet")
	err := w.account.CheckAccount(accountID)
	if err != nil {
		return err
	}
	err = w.securityCode.CheckCode(securityCode)
	if err != nil {
		return err
	}
	w.wallet.CreditBalance(amount)
	w.notification.SendWalletCreditNotification()
	w.ledger.MakeEntry(accountID, "credit", amount)
	return nil
}

func (w *WalletFacade) DeductMoneyFromWallet(accountID string, securityCode int, amount int) error {
	fmt.Println("Starting debit money from wallet")
	err := w.account.CheckAccount(accountID)
	if err != nil {
		return err
	}

	err = w.securityCode.CheckCode(securityCode)
	if err != nil {
		return err
	}
	err = w.wallet.DebitBalance(amount)
	if err != nil {
		return err
	}
	w.notification.SendWalletDebitNotification()
	w.ledger.MakeEntry(accountID, "credit", amount)
	return nil
}

// =

type Account struct {
	name string
}

func NewAccount(accountName string) *Account {
	return &Account{
		name: accountName,
	}
}

func (a *Account) CheckAccount(accountName string) error {
	if a.name != accountName {
		return fmt.Errorf("Account Name is incorrect")
	}
	fmt.Println("Account Verified")
	return nil
}

// =

type SecurityCode struct {
	code int
}

func NewSecurityCode(code int) *SecurityCode {
	return &SecurityCode{
		code: code,
	}
}

func (s *SecurityCode) CheckCode(incomingCode int) error {
	if s.code != incomingCode {
		return fmt.Errorf("Security Code is incorrect")
	}
	fmt.Println("SecurityCode Verified")
	return nil
}

// =

type Wallet struct {
	balance int
}

func NewWallet() *Wallet {
	return &Wallet{
		balance: 0,
	}
}

func (w *Wallet) CreditBalance(amount int) {
	w.balance += amount
	fmt.Println("Wallet balance added successfully")
	return
}

func (w *Wallet) DebitBalance(amount int) error {
	if w.balance < amount {
		return fmt.Errorf("Balance is not sufficient")
	}
	fmt.Println("Wallet balance is Sufficient")
	w.balance = w.balance - amount
	return nil
}

// =

type Ledger struct{}

func (s *Ledger) MakeEntry(accountID, txnType string, amount int) {
	fmt.Printf("Make ledger entry for accountId %s with txnType %s for amount %d\n", accountID, txnType, amount)
	return
}

// =

type Notification struct{}

func (n *Notification) SendWalletCreditNotification() {
	fmt.Println("Sending wallet credit notification")
}

func (n *Notification) SendWalletDebitNotification() {
	fmt.Println("Sending wallet debit notification")
}
