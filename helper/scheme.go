package helper

import "fmt"

func newSchemeStringExactIndex(name string) *Scheme {
    return newScheme(name, "string", "@index(exact)")
}

func newSchemeString(name string) *Scheme {
    return newScheme(name, "string", "")
}

func newScheme(name, schemeType, others string) *Scheme {
    return &Scheme{
        Name:   name,
        Type:   schemeType,
        Others: others,
    }
}

type Scheme struct {
    Name   string
    Type   string
    Others string
}

func (scheme *Scheme) String() string {
    return fmt.Sprintf("%s: %s %s . \r\n", scheme.Name, scheme.Type, scheme.Others)
}

func (scheme *Scheme) Equal(src *Scheme) bool {
    if scheme.Name == src.Name &&
        scheme.Type == src.Type &&
        scheme.Others == src.Others {
        return true
    }
    return false
}

type Schemes []*Scheme

func (schemes Schemes) Add(newSchemes ...*Scheme) Schemes {
    list := schemes
    for _, ns := range newSchemes {
        if schemes.Exists(ns) {
            continue
        }
        list = append(list, ns)
    }
    return list
}

func (schemes Schemes) Exists(s *Scheme) bool {
    for _, scheme := range schemes {
        if s.Equal(scheme) {
            return true
        }
    }
    return false
}

func (schemes Schemes) String() string {
    all := ""
    for _, s := range schemes {
        all = all + s.String()
    }
    return all
}
