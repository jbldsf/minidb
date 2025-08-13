package main

import "backend/server"

func main() {
	server.Start()
	// b, _ := base64.RawURLEncoding.DecodeString("eyJpZCI6MX0")
	// println(string(b))
}

//{"continents":[{"id":1,"area":1.23,"name":"Africa","population":123},{"id":2,"area":1.23,"name":"America","population":123},{"id":3,"area":1.23,"name":"Asia","population":123}]} eyJpZCI6MX0
