/* Before you execute the program, Launch `cqlsh` and execute:
create keyspace example with replication = { 'class' : 'SimpleStrategy', 'replication_factor' : 1 };
create table example.tweet(timeline text, id UUID, text text, PRIMARY KEY(id));
create index on example.tweet(timeline);
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

	// cluster.PoolConfig.HostSelectionPolicy = gocql.DCAwareRoundRobinPolicy("local-datacenter-name")

	// Create host selection policy using a simple host pool
	// cluster.PoolConfig.HostSelectionPolicy = gocql.HostPoolHostPolicy(hostpool.New(nil))
	// Create host selection policy using an epsilon greedy pool
	// cluster.PoolConfig.HostSelectionPolicy = gocql.HostPoolHostPolicy(
		// hostpool.NewEpsilonGreedy(nil, 0, &hostpool.LinearEpsilonValueCalculator{}),)
	
	// cluster.PoolConfig.HostSelectionPolicy = gocql.RoundRobinHostPolicy()
	cluster.PoolConfig.HostSelectionPolicy = gocql.TokenAwareHostPolicy(gocql.RoundRobinHostPolicy())
	cluster.Keyspace = "example"
	cluster.Consistency = gocql.Quorum
	session, _ := cluster.CreateSession()
	defer session.Close()

	// insert a tweet
	if err := session.Query(`INSERT INTO tweet (timeline, id, text) VALUES (?, ?, ?)`,
		"me", gocql.TimeUUID(), "hello world").Exec(); err != nil {
		log.Fatal(err)
	}

	var id gocql.UUID
	var text string

	/* Search for a specific set of records whose 'timeline' column matches
	 * the value 'me'. The secondary index that we created earlier will be
	 * used for optimizing the search */
	if err := session.Query(`SELECT id, text FROM tweet WHERE timeline = ? LIMIT 1`,
		"me").Consistency(gocql.One).Scan(&id, &text); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Tweet:", id, text)

	// list all tweets
	iter := session.Query(`SELECT id, text FROM tweet WHERE timeline = ?`, "me").Iter()
	for iter.Scan(&id, &text) {
		fmt.Println("Tweet:", id, text)
	}
	if err := iter.Close(); err != nil {
		log.Fatal(err)
	}
}