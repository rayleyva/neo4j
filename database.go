// Copyright (c) 2012-2013 Jason McVetta.  This is Free Software, released under
// the terms of the GPL v3.  See http://www.gnu.org/copyleft/gpl.html for details.
// Resist intellectual serfdom - the ownership of ideas is akin to slavery.

// Package neo4j provides a client for the Neo4j graph database.
package neo4j

import (
	"github.com/jmcvetta/restclient"
	"log"
	"net/url"
	"strconv"
)

// A Database is a REST client connected to a Neo4j database.
type Database struct {
	url *url.URL // Root URL for REST API
	rc  *restclient.Client
	//
	HrefNode      string
	HrefRefNode   string
	HrefNodeIndex string
	HrefRelIndex  string
	HrefExtInfo   string
	HrefRelTypes  string
	HrefBatch     string
	HrefCypher    string
	Version       string
}

// A serviceRoot describes services available on the Neo4j server
type serviceRoot struct {
	Extensions    interface{} `json:"extensions"`
	HrefNode      string      `json:"node"`
	HrefRefNode   string      `json:"reference_node"`
	HrefNodeIndex string      `json:"node_index"`
	HrefRelIndex  string      `json:"relationship_index"`
	HrefExtInfo   string      `json:"extensions_info"`
	HrefRelTypes  string      `json:"relationship_types"`
	HrefBatch     string      `json:"batch"`
	HrefCypher    string      `json:"cypher"`
	Version       string      `json:"neo4j_version"`
}

// Connect establishes a connection to the Neo4j server.
func Connect(uri string) (*Database, error) {
	var sr serviceRoot
	var e neoError
	db := &Database{
		rc: restclient.New(),
	}
	u, err := url.Parse(uri)
	if err != nil {
		return db, err
	}
	db.url = u
	req := restclient.RequestResponse{
		Url:    u.String(),
		Method: "GET",
		Result: &sr,
		Error:  &e,
	}
	status, err := db.rc.Do(&req)
	if err != nil {
		logPretty(req)
		return db, err
	}
	switch {
	case status == 404:
		return db, InvalidDatabase
	case status != 200 || sr.Version == "":
		log.Println("Status " + strconv.Itoa(status) + " trying to cconnect to " + u.String())
		logPretty(req)
		return db, BadResponse
	}
	// Populate Database struct
	db.HrefNode = sr.HrefNode
	db.HrefRefNode = sr.HrefRefNode
	db.HrefNodeIndex = sr.HrefNodeIndex
	db.HrefRelIndex = sr.HrefRelIndex
	db.HrefExtInfo = sr.HrefExtInfo
	db.HrefRelTypes = sr.HrefRelTypes
	db.HrefBatch = sr.HrefBatch
	db.HrefCypher = sr.HrefCypher
	db.Version = sr.Version
	// Success!
	return db, nil
}
