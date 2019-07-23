/* Before you execute the program, Launch `cqlsh` and execute:
CREATE KEYSPACE hydro_monitor_data WITH REPLICATION = { 'class' : 'SimpleStrategy', 'replication_factor' : 3 };
CREATE TABLE hydro_monitor_data.users (email string, password string, admin bool, primary key (email));
*/
package main

import (
	"fmt"
	"log"

	"github.com/gocql/gocql"
	// "github.com/hailocab/go-hostpool"
)

func main() {
	// connect to the cluster
	cluster := gocql.NewCluster("192.168.50.41", "192.168.50.42", "192.168.50.43")
	// Select keyspace
	cluster.Keyspace = "usertest"

	// cluster.PoolConfig.HostSelectionPolicy = gocql.DCAwareRoundRobinPolicy("local-datacenter-name")

	// Create host selection policy using a simple host pool
	// cluster.PoolConfig.HostSelectionPolicy = gocql.HostPoolHostPolicy(hostpool.New(nil))
	// Create host selection policy using an epsilon greedy pool
	// cluster.PoolConfig.HostSelectionPolicy = gocql.HostPoolHostPolicy(
		// hostpool.NewEpsilonGreedy(nil, 0, &hostpool.LinearEpsilonValueCalculator{}),)
	
	cluster.PoolConfig.HostSelectionPolicy = gocql.RoundRobinHostPolicy()
	//cluster.PoolConfig.HostSelectionPolicy = gocql.TokenAwareHostPolicy(gocql.RoundRobinHostPolicy())
	//cluster.Keyspace = "example"
	//cluster.Consistency = gocql.Quorum
	session, _ := cluster.CreateSession()
	defer session.Close()

	// insert a user
	if err := session.Query(`INSERT INTO users (email, password, admin) VALUES (?, ?, ?)`,
		"bob@example.com", "secretandencryptedpassword", false).Exec(); err != nil {
		log.Fatal(err)
	}

	var email string

	/* Search for a specific set of records whose 'timeline' column matches
	 * the value 'me'. The secondary index that we created earlier will be
	 * used for optimizing the search */
	// if err := session.Query(`SELECT id, text FROM tweet WHERE timeline = ? LIMIT 1`,
	if err := session.Query(`SELECT email FROM users WHERE admin = ? LIMIT 1 ALLOW FILTERING`,
		false).Consistency(gocql.One).Scan(&email); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Non admin:", email)

	// list all users
	iter := session.Query(`SELECT email FROM users`).Iter()
	for iter.Scan(&email) {
		fmt.Println("User:", email)
	}
	if err := iter.Close(); err != nil {
		log.Fatal(err)
	}
}