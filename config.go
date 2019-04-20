package main

type ConfigProperties struct {
	MongoURL   string
	Database   string
	Collection string
	Port       string
}

var Config = &ConfigProperties{
	MongoURL:   "mongodb://mongodb:27017",
	Database:   "paymentsDev",
	Collection: "payments",
	Port:       "8080",
}

var TestConfig = &ConfigProperties{
	MongoURL:   "mongodb://mongodb:27017",
	Database:   "paymentsTest",
	Collection: "payments",
	Port:       "8080",
}
