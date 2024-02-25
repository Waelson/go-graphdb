package main

import (
	"context"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

func main() {
	println("Iniciando...")
	//dbUri := "bolt://localhost"
	dbUri := "bolt+ssc://????.9:8182"
	driver, err := neo4j.NewDriverWithContext(dbUri, neo4j.NoAuth())
	if err != nil {
		panic(err)
	}

	ctx := context.Background()

	println("Verificando conectividade")
	err = driver.VerifyConnectivity(ctx)
	if err != nil {
		panic(err)
	}

	session := driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close(ctx)

	defer driver.Close(ctx) // Make sure to handle errors during deferred calls

	cypher := `
		MATCH path = (u:User {name: 'Root'})-[*]-(p:Payment {name: '001_202402'})
		UNWIND relationships(path) AS rel
		RETURN DISTINCT TYPE(rel) AS RelationshipType
		`
	result, err := session.Run(ctx, cypher, nil)

	list := make([]any, 0)
	for result.Next(ctx) {
		record := result.Record()
		if value, ok := record.Get("RelationshipType"); ok {
			list = append(list, value)
		}
	}

	fmt.Printf("%v\n", list)

	/*
		println("Tentando inserir registros")

		greeting, err := session.ExecuteWrite(ctx, func(transaction neo4j.ManagedTransaction) (any, error) {
			cypher := `
					create(u1:User{id:'1',name:'Root'})
					create(u2:User{id:'2',name:'Colaborador 1'})
					create(u3:User{id:'3',name:'Colaborador 2'})

					create(g:Group{id:'1',name:'group_root_1'})

					create(p1:Payment{id:'1',name:'001_202402'})
					create(p2:Payment{id:'2',name:'002_202402'})
					create(p3:Payment{id:'3',name:'003_202402'})


					CREATE (u1)-[:OWNER]->(g)
					CREATE (u2)-[:MEMBER]->(g)
					CREATE (u3)-[:MEMBER]->(g)

					CREATE (g)-[:PARENT]->(p1)
					CREATE (g)-[:PARENT]->(p2)
					CREATE (g)-[:PARENT]->(p3)
		    `
			result, err := transaction.Run(ctx, cypher, map[string]any{"message": "hello, world"})

			if err != nil {
				return nil, err
			}

			list := make([]any, 0)
			if result.Next(ctx) {
				list = append(list, result.Record().Values[0])
			}

			return list, result.Err()

			// 		if result.Next(ctx) {
			// 			return result.Record().Values[0], nil
			// 		}

			// 		return nil, result.Err()
		})
		if err != nil {
			panic(err)
		}

		fmt.Printf("%v\n", greeting)
	*/

}
