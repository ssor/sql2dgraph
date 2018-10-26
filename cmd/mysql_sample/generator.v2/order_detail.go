package generator_v2

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
