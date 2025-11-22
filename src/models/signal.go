package models

type Signal struct {
	Type    string `cbor:"Type"`
	Message string `cbor:"Message"`
}
