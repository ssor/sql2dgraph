package generator_v2

import "time"

/*
CREATE TABLE customers (
  customerNumber int(11) NOT NULL,
  customerName varchar(50) NOT NULL,
  contactLastName varchar(50) NOT NULL,
  contactFirstName varchar(50) NOT NULL,
  phone varchar(50) NOT NULL,
  addressLine1 varchar(50) NOT NULL,
  addressLine2 varchar(50) DEFAULT NULL,
  city varchar(50) NOT NULL,
  state varchar(50) DEFAULT NULL,
  postalCode varchar(15) DEFAULT NULL,
  country varchar(50) NOT NULL,
  salesRepEmployeeNumber int(11) DEFAULT NULL,
  creditLimit decimal(10,2) DEFAULT NULL,
  PRIMARY KEY (customerNumber),
  KEY salesRepEmployeeNumber (salesRepEmployeeNumber)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;
*/
type Customer struct {
    Uid              string    `json:"uid,omitempty"`
    CustomerNumber   int       `json:"customer_number,omitempty"`
    CustomerName     string    `json:"customer_name,omitempty"`
    ContactLastName  string    `json:"contact_last_name,omitempty"`
    ContactFirstName string    `json:"contact_first_name,omitempty"`
    Phone            string    `json:"phone,omitempty"`
    AddressLine1     string    `json:"address_line1,omitempty"`
    AddressLine2     string    `json:"address_line2,omitempty"`
    City             string    `json:"city,omitempty"`
    State            string    `json:"state,omitempty"`
    PostalCode       string    `json:"postal_code,omitempty"`
    Country          string    `json:"country,omitempty"`
    CreditLimit      float64   `json:"credit_limit,omitempty"`
    Employee         *Employee `json:"service_by_employee,omitempty"`
}

func newCustomer(number int) *Customer {
    return &Customer{
        CustomerNumber: number,
    }
}

func (customer *Customer) DependentObjectHasUid() bool {
    if customer.Employee != nil {
        if len(customer.Employee.Uid) <= 0 {
            return false
        }
    }
    return true
}

/*
CREATE TABLE employees (
  employeeNumber int(11) NOT NULL,
  lastName varchar(50) NOT NULL,
  firstName varchar(50) NOT NULL,
  extension varchar(10) NOT NULL,
  email varchar(100) NOT NULL,
  officeCode varchar(10) NOT NULL,
  reportsTo int(11) DEFAULT 0,
  jobTitle varchar(50) NOT NULL,
  PRIMARY KEY (employeeNumber),
  KEY reportsTo (reportsTo),
  KEY officeCode (officeCode)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;
*/
type Employee struct {
    Uid            string    `json:"uid,omitempty"`
    EmployeeNumber int       `json:"employee_number,omitempty"`
    LastName       string    `json:"last_name,omitempty"`
    FirstName      string    `json:"first_name,omitempty"`
    Extension      string    `json:"extension,omitempty"`
    Email          string    `json:"email,omitempty"`
    Office         *Office   `json:"office_work_in,omitempty"`
    ReportsTo      *Employee `json:"report_to,omitempty"`
    JobTitle       string    `json:"job_title,omitempty"`
}

func newEmployee(number int, reportsTo *Employee) *Employee {
    return &Employee{
        EmployeeNumber: number,
        ReportsTo:      reportsTo,
    }
}

func (employee *Employee) DependentObjectHasUid() bool {
    if employee.ReportsTo != nil {
        if len(employee.ReportsTo.Uid) <= 0 {
            return false
        }
    }
    return true
}

/*
CREATE TABLE offices (
  officeCode varchar(10) NOT NULL,
  city varchar(50) NOT NULL,
  phone varchar(50) NOT NULL,
  addressLine1 varchar(50) NOT NULL,
  addressLine2 varchar(50) DEFAULT NULL,
  state varchar(50) DEFAULT NULL,
  country varchar(50) NOT NULL,
  postalCode varchar(15) NOT NULL,
  territory varchar(10) NOT NULL,
  PRIMARY KEY (officeCode)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;
*/
type Office struct {
    Uid          string `json:"uid,omitempty"`
    OfficeCode   string `json:"office_code,omitempty"`
    City         string `json:"city,omitempty"`
    Phone        string `json:"phone,omitempty"`
    AddressLine1 string `json:"address_line1,omitempty"`
    AddressLine2 string `json:"address_line2,omitempty"`
    State        string `json:"state,omitempty"`
    PostalCode   string `json:"postal_code,omitempty"`
    Country      string `json:"country,omitempty"`
    Territory    string `json:"territory,omitempty"`
}

func newOffice(code string) *Office {
    return &Office{
        OfficeCode: code,
    }
}

func (office *Office) DependentObjectHasUid() bool {
    return true
}

/*

CREATE TABLE orderdetails (
  orderNumber int(11) NOT NULL,
  productCode varchar(15) NOT NULL,
  quantityOrdered int(11) NOT NULL,
  priceEach decimal(10,2) NOT NULL,
  orderLineNumber smallint(6) NOT NULL,
  PRIMARY KEY (orderNumber,productCode),
  KEY productCode (productCode)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

*/
type OrderDetail struct {
    Uid             string   `json:"uid,omitempty"`
    Order           *Order   `json:"detail_of_order,omitempty"`
    Product         *Product `json:"product_ordered,omitempty"`
    QuantityOrdered int      `json:"quantity_ordered,omitempty"`
    PriceEach       float64  `json:"price_each,omitempty"`
    OrderLineNumber int      `json:"order_line_number,omitempty"`
}

/*
CREATE TABLE orders (
  orderNumber int(11) NOT NULL,
  orderDate date NOT NULL,
  requiredDate date NOT NULL,
  shippedDate date DEFAULT NULL,
  orderStatus varchar(15) NOT NULL,
  comments text,
  customerNumber int(11) NOT NULL,
  PRIMARY KEY (orderNumber),
  KEY customerNumber (customerNumber)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

*/

type Order struct {
    Uid          string    `json:"uid,omitempty"`
    OrderNumber  int       `json:"order_number,omitempty"`
    OrderDate    time.Time `json:"order_date,omitempty"`
    RequiredDate time.Time `json:"required_date,omitempty"`
    ShippedDate  time.Time `json:"shipped_date,omitempty"`
    OrderStatus  string    `json:"order_status,omitempty"`
    Comments     string    `json:"comments,omitempty"`
    Customer     *Customer `json:"ordered_by_customer,omitempty"`
}

/*

CREATE TABLE payments (
  customerNumber int(11) NOT NULL,
  checkNumber varchar(50) NOT NULL,
  paymentDate date NOT NULL,
  amount decimal(10,2) NOT NULL,
  PRIMARY KEY (customerNumber,checkNumber)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

*/

type Payment struct {
    Uid            string    `json:"uid,omitempty"`
    CustomerNumber *Customer `json:"payment_by_customer,omitempty"`
    CheckNumber    string    `json:"check_number,omitempty"`
    PaymentDate    time.Time `json:"payment_date,omitempty"`
    Amount         float64   `json:"amount,omitempty"`
}

/*
CREATE TABLE productlines (
  productLine varchar(50) NOT NULL,
  textDescription varchar(4000) DEFAULT NULL,
  PRIMARY KEY (productLine)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;
*/

type ProductLine struct {
    Uid             string `json:"uid,omitempty"`
    ProductLine     string `json:"product_line,omitempty"`
    TextDescription string `json:"text_description,omitempty"`
}

/*

CREATE TABLE products (
  productCode varchar(15) NOT NULL,
  productName varchar(70) NOT NULL,
  productLine varchar(50) NOT NULL,
  productScale varchar(10) NOT NULL,
  productVendor varchar(50) NOT NULL,
  productDescription text NOT NULL,
  quantityInStock smallint(6) NOT NULL,
  buyPrice decimal(10,2) NOT NULL,
  MSRP decimal(10,2) NOT NULL,
  PRIMARY KEY (productCode),
  KEY productLine (productLine)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

*/

type Product struct {
    Uid                string       `json:"uid,omitempty"`
    ProductCode        string       `json:"product_code,omitempty"`
    ProductName        string       `json:"product_name,omitempty"`
    ProductLine        *ProductLine `json:"produced_by_line,omitempty"`
    ProductScale       string       `json:"product_scale,omitempty"`
    ProductVendor      string       `json:"product_vendor,omitempty"`
    ProductDescription string       `json:"product_description,omitempty"`
    QuantityInStock    int          `json:"quantity_in_stock,omitempty"`
    BuyPrice           float64      `json:"buy_price,omitempty"`
    MSRP               float64      `json:"msrp,omitempty"`
}
