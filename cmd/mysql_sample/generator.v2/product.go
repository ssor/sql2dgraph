package generator_v2

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
