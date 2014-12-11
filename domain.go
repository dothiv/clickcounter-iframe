package clickcounteriframe

import "time"

type EntityInterface interface {
}

type Domain struct {
	EntityInterface
	Id       int64
	Name     string
	Redirect string
	Created  *time.Time
	Updated  *time.Time
}
