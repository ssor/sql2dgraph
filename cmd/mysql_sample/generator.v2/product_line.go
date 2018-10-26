package generator_v2

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
