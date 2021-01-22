package main

import api "github.com/bastianhussi/todos-api"

func main() {
	c := new(api.Config)
	if err := c.ReadConfig(); err != nil {
		panic(err)
	}
	rdb := api.NewRedisClient()
	defer rdb.Close()

	db, err := api.NewDB()
	if err != nil {
		panic(err)
	}

	defer db.Close()
}
