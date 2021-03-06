package main

import (
	"crypto/tls"
	"net"

	mgo "gopkg.in/mgo.v2"
)

//DBHandler holds a session information for accessing MongoDB instance using mgo
type DBHandler struct {
	DB *mgo.Session
}

func (h *DBHandler) initDB() error {
	tlsConfig := &tls.Config{}

	dialInfo := &mgo.DialInfo{
		Addrs: []string{"healthyrepo-shard-00-00-7p8dp.mongodb.net:27017",
			"healthyrepo-shard-00-01-7p8dp.mongodb.net:27017",
			"healthyrepo-shard-00-02-7p8dp.mongodb.net:27017"},
		Database: "admin",
		Username: getEnv("ATLAS_USERNAME"),
		Password: getEnv("ATLAS_PASSWORD"),
	}
	dialInfo.DialServer = func(addr *mgo.ServerAddr) (net.Conn, error) {
		conn, err := tls.Dial("tcp", addr.String(), tlsConfig)
		return conn, err
	}
	session, err := mgo.DialWithInfo(dialInfo)
	if err != nil {
		return err
	}

	h.DB = session

	return nil
}
