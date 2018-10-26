package generator_v2

import "time"

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
