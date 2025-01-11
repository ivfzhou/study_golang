// 责任链是一种行为设计模式，允许你将请求沿着处理者链进行发送，直至其中一个处理者对其进行处理。
// 当程序需要使用不同方式处理不同种类请求，而且请求类型和顺序预先未知时，可以使用责任链模式。
// 当必须按顺序执行多个处理者时，可以使用该模式。
// 如果所需处理者及其顺序必须在运行时进行改变，可以使用责任链模式。
package design_pattern_test

import (
	"fmt"
	"testing"
)

func TestChainOfResponsibility(t *testing.T) {
	var cashier Department = &Cashier{}

	// Set next for medical department
	var medical Department = &Medical{}
	medical.SetNext(cashier)

	// Set next for doctor department
	var doctor Department = &Doctor{}
	doctor.SetNext(medical)

	// Set next for reception department
	var reception Department = &Reception{}
	reception.SetNext(doctor)

	patient := &Patient{Name: "abc"}
	// Patient visiting
	reception.Execute(patient)
}

// ===

type Department interface {
	Execute(*Patient)
	SetNext(Department)
}

type Patient struct {
	Name              string
	RegistrationDone  bool
	DoctorCheckUpDone bool
	MedicineDone      bool
	PaymentDone       bool
}

// =

type Reception struct {
	next Department
}

func (r *Reception) Execute(p *Patient) {
	if p.RegistrationDone {
		fmt.Println("Patient registration already done")
		r.next.Execute(p)
		return
	}
	fmt.Println("Reception registering patient")
	p.RegistrationDone = true
	r.next.Execute(p)
}

func (r *Reception) SetNext(next Department) {
	r.next = next
}

// =

type Doctor struct {
	next Department
}

func (d *Doctor) Execute(p *Patient) {
	if p.DoctorCheckUpDone {
		fmt.Println("Doctor checkup already done")
		d.next.Execute(p)
		return
	}
	fmt.Println("Doctor checking patient")
	p.DoctorCheckUpDone = true
	d.next.Execute(p)
}

func (d *Doctor) SetNext(next Department) {
	d.next = next
}

// =

type Medical struct {
	next Department
}

func (m *Medical) Execute(p *Patient) {
	if p.MedicineDone {
		fmt.Println("Medicine already given to patient")
		m.next.Execute(p)
		return
	}
	fmt.Println("Medical giving medicine to patient")
	p.MedicineDone = true
	m.next.Execute(p)
}

func (m *Medical) SetNext(next Department) {
	m.next = next
}

// =

type Cashier struct {
	next Department
}

func (c *Cashier) Execute(p *Patient) {
	if p.PaymentDone {
		fmt.Println("Payment Done")
	}
	fmt.Println("Cashier getting money from patient patient")
}

func (c *Cashier) SetNext(next Department) {
	c.next = next
}
