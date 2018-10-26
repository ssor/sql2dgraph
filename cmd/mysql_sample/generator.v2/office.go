package generator_v2

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
