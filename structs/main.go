package main

import "fmt"

type contactInfo struct {
	email   string
	zipCode int
}

type person struct {
	firstname string
	lastname  string
	contact   contactInfo
}

func main() {
	// akshat := person{"Akshat", "Shetty"}
	// akshat := person{firstname: "Akshat", lastname: "Shetty"}
	// var akshat person
	// akshat.firstname = "Akshat"
	// fmt.Println(akshat)
	akshat := person{
		firstname: "Akshat",
		lastname:  "Shetty",
		contact: contactInfo{
			email:   "akshatshetty2908@gmail.com",
			zipCode: 400064,
		},
	}
	akshat.print()
	akshat.updateName("ak")
	akshat.print()
	(&akshat).updateName("Ak")
	akshat.print()
}

func (p *person) updateName(newFirstName string) {
	(*p).firstname = newFirstName
}

func (p person) print() {
	fmt.Printf("%+v\n", p)
}
