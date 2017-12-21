package main

import (
	"net/http"
	"github.com/graphql-go/handler"
	"github.com/graphql-go/graphql"
	"log"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/mattn/go-sqlite3"
	"time"
)

type User struct {
	Id int64 `json:"id"`
	Name  string `json:"name"`
	Favorites []Favorite `json:"favorites"`
}

type Favorite struct {
	Id int64 `json:"id"`
	Uid int64 `json:"uid"`
	Url  string `json:"url"`
	Savetype string `json:"savetype"`
	Time  string `json:"time"`
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {


	db, err := sql.Open("sqlite3", "./favoriteAds.db")
	checkErr(err)


	var sql_table=`CREATE TABLE IF NOT EXISTS favorite (
	id  INTEGER PRIMARY KEY AUTOINCREMENT,
		uid int(11) NOT NULL,
		url varchar(512) NOT NULL,
		savetype varchar(10) NOT NULL,
		time timestamp NOT NULL ,
		UNIQUE (uid, url) ON CONFLICT REPLACE
	); `


	res, err := db.Exec(sql_table)
	if err != nil { panic(err) }
	println(res)






	var fadsType = graphql.NewObject(graphql.ObjectConfig{
		Name: "favorite",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type: graphql.Int,
			},
			"uid": &graphql.Field{
				Type: graphql.Int,
			},
			"url": &graphql.Field{
				Type: graphql.String,
			},
			"savetype": &graphql.Field{
				Type: graphql.String,
			},
			"time": &graphql.Field{
				Type: graphql.String,
			},

		},
	})


	var userType = graphql.NewObject(graphql.ObjectConfig{
		Name: "user",
		Fields: graphql.Fields{
				"id": &graphql.Field{
					Type: graphql.Int,
				},
			"favorites": &graphql.Field{
				Type: graphql.NewList(fadsType),

			},


		},
	})
	

	fields := graphql.Fields{
		"User": &graphql.Field{
			Type: userType,
			Args: graphql.FieldConfigArgument{
				"id": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.Int),
				},
				},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				 var userid = int64(p.Args["id"].(int))
				rows, err := db.Query("SELECT * FROM favorite WHERE uid=1")

				checkErr(err)
				var id int64
				var uid int64
				var url string
				var savetype string
				var time string
				favorites := []Favorite{}

				for rows.Next() {

					err = rows.Scan(&id,&uid, &url,&savetype,&time)
					checkErr(err)
					favorites = append(favorites, Favorite{Id:id,Uid:uid,Url:url,Savetype:savetype,Time:time})
				}

				user := &User{
					Id: userid,
					Favorites:favorites,
				}


				return  user, nil
			},
		},
	}

	
	rootMutation := graphql.NewObject(graphql.ObjectConfig{
		Name: "RootMutation",
		Fields: graphql.Fields{
			"addfads": &graphql.Field{
				Type: fadsType, // the return type for this field
				Args: graphql.FieldConfigArgument{
					"uid": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.Int),
					},
					"url": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
					"savetype": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},

				},
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {

					uid := int64(params.Args["uid"].(int))
					url, _ := params.Args["url"].(string)
					savetype, _ := params.Args["savetype"].(string)

					stmt, err := db.Prepare("INSERT INTO favorite(uid, url, Savetype,time) values(?,?,?,?)")
					checkErr(err)
					res, err := stmt.Exec(uid,url,savetype,time.Now())
					checkErr(err)

					lid,err := res.LastInsertId()

					checkErr(err)
					newFads := &Favorite{
						Id: lid,
						Uid:   uid,
						Url: url,
						Savetype: savetype,
					}
					return newFads,nil
				},
			},
			"deletfads": &graphql.Field{
				Type: fadsType, // the return type for this field
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.Int),
					},
				},
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {

					// marshall and cast the argument value
					id, _ := params.Args["id"].(int)

					print("id to delete ")
					println(id)
					stmt, err := db.Prepare("delete from favorite where id=?")
					checkErr(err)

					res, err := stmt.Exec(id)
					checkErr(err)

					affect, err := res.RowsAffected()
					checkErr(err)

					println("affect ")
					println(affect)

					return affect,nil
				},
			},
		},
	})

	rootQuery := graphql.ObjectConfig{Name: "RootQuery", Fields: fields}
	schemaConfig := graphql.SchemaConfig{Query: graphql.NewObject(rootQuery),Mutation:rootMutation}
	schema, err := graphql.NewSchema(schemaConfig)
	if err != nil {
		log.Fatalf("failed to create new schema, error: %v", err)
	}

	h := handler.New(&handler.Config{
		Schema: &schema,
		//Pretty: true,
		GraphiQL: true,
	})

	// serve HTTP
	http.Handle("/graphql", h)
	http.ListenAndServe(":8080", nil)
}
