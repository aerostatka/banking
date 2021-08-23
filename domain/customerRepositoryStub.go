package domain

type CustomerRepositoryStub struct {
	customers []Customer
}

func (s CustomerRepositoryStub) FindAll() ([]Customer, error) {
	return s.customers, nil
}

func GenerateCustomerRepositoryStub() CustomerRepositoryStub {
	customersList := []Customer{
		{Id: "uuid1111", Name: "Alex Johnson", City: "New York", ZipCode: "10050", DOB: "05/25/1985", Status: true},
		{Id: "uuid1112", Name: "Dakota Fanning", City: "Jersey City", ZipCode: "07306", DOB: "12/31/1990", Status: true},
		{Id: "uuid1113", Name: "Alexis Castle", City: "Bensalem", ZipCode: "19020", DOB: "01/09/1991", Status: true},
	}

	return CustomerRepositoryStub{customers: customersList}
}
