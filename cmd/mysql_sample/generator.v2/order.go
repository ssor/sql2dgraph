package generator_v2

import "time"

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
