package main

type EsErr struct{}

func (m *EsErr) Error() string {
	return "boom"
}
