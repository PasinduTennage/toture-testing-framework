package cmd

/*
	This file defines the client struct and the new method that is invoked when creating a new client by the main

*/

type Toture struct {
}

func New() *Toture {
	to := Toture{}
	return &to

}
