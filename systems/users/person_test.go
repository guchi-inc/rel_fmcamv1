// Copyright (c) 2024 xueyingdu zco-inc All Right reserved. privacy_zco@outlook.com
package users

import (
	"database/sql"
	"fmcam/common/databases"
	"log"

	"github.com/jmoiron/sqlx"
)

/*
sqlx 使用范例
*/
var schema = `
DROP TABLE IF EXISTS person;
CREATE TABLE person (
    first_name text,
    last_name text,
    email text
);

DROP TABLE IF EXISTS place;
CREATE TABLE place (
    country text,
    city text NULL,
    telcode integer
)
`

type Person struct {
	FirstName string `db:"first_name"`
	LastName  string `db:"last_name"`
	Email     string
}

type Place struct {
	Country string
	City    sql.NullString
	TelCode int
}

func (u UserRepository) SchemalList() (user *[]Person, err error) {

	var DBConn = databases.DBMysql
	// this Pings the database trying to connect
	// use sqlx.Open() for sql.Open() semantics
	db, err := sqlx.Connect("postgres", "user=foo dbname=bar sslmode=disable")
	if err != nil {
		log.Fatalln(err)
	}

	// exec the schema or fail; multi-statement Exec behavior varies between
	// database drivers;  pq will exec them all, sqlite3 won't, ymmv
	db.MustExec(schema)

	tx := DBConn.MustBegin()

	tx.MustExec("INSERT INTO person (first_name, last_name, email) VALUES ($1, $2, $3)", "Jason", "Moiron", "jmoiron@jmoiron.net")
	tx.MustExec("INSERT INTO person (first_name, last_name, email) VALUES ($1, $2, $3)", "John", "Doe", "johndoeDNE@gmail.net")
	tx.MustExec("INSERT INTO place (country, city, telcode) VALUES ($1, $2, $3)", "United States", "New York", "1")
	tx.MustExec("INSERT INTO place (country, telcode) VALUES ($1, $2)", "Hong Kong", "852")
	tx.MustExec("INSERT INTO place (country, telcode) VALUES ($1, $2)", "Singapore", "65")
	// Named queries can use structs, so if you have an existing struct (i.e. person := &Person{}) that you have populated, you can pass it in as &person
	tx.NamedExec("INSERT INTO person (first_name, last_name, email) VALUES (:first_name, :last_name, :email)", &Person{"Jane", "Citizen", "jane.citzen@example.com"})
	tx.Commit()

	// 批量查询数据库 Query the database, storing results in a []Person (wrapped in []interface{})
	people := []Person{}
	db.Select(&people, "SELECT * FROM person ORDER BY first_name ASC")
	jason, john := people[0], people[1]

	logger.Printf("%#v\n%#v", jason, john)
	// Person{FirstName:"Jason", LastName:"Moiron", Email:"jmoiron@jmoiron.net"}
	// Person{FirstName:"John", LastName:"Doe", Email:"johndoeDNE@gmail.net"}

	// 查询单个数据 You can also get a single result, a la QueryRow
	jason = Person{}
	err = db.Get(&jason, "SELECT * FROM person WHERE first_name=$1", "Jason")
	logger.Printf("%#v\n", jason)
	// Person{FirstName:"Jason", LastName:"Moiron", Email:"jmoiron@jmoiron.net"}

	//如果有null字段，请使用sql.NULL if you have null fields and use SELECT *, you must use sql.Null* in your struct
	places := []Place{}
	err = db.Select(&places, "SELECT * FROM place ORDER BY telcode ASC")
	if err != nil {
		logger.Println(err)
		return
	}
	usa, singsing, honkers := places[0], places[1], places[2]

	logger.Printf("%#v\n%#v\n%#v\n", usa, singsing, honkers)
	// Place{Country:"United States", City:sql.NullString{String:"New York", Valid:true}, TelCode:1}
	// Place{Country:"Singapore", City:sql.NullString{String:"", Valid:false}, TelCode:65}
	// Place{Country:"Hong Kong", City:sql.NullString{String:"", Valid:false}, TelCode:852}

	// 只使用一个结构体时，循环输出 Loop through rows using only one struct
	place := Place{}
	rows, err := db.Queryx("SELECT * FROM place")
	for rows.Next() {
		err := rows.StructScan(&place)
		if err != nil {
			log.Fatalln(err)
		}
		logger.Printf("place:%#v\n", place)
	}
	// Place{Country:"United States", City:sql.NullString{String:"New York", Valid:true}, TelCode:1}
	// Place{Country:"Hong Kong", City:sql.NullString{String:"", Valid:false}, TelCode:852}
	// Place{Country:"Singapore", City:sql.NullString{String:"", Valid:false}, TelCode:65}

	// Named queries, using `:name` as the bindvar.  Automatic bindvar support
	// which takes into account the dbtype based on the driverName on sqlx.Open/Connect
	_, err = db.NamedExec(`INSERT INTO person (first_name,last_name,email) VALUES (:first,:last,:email)`,
		map[string]interface{}{
			"first": "Bin",
			"last":  "Smuth",
			"email": "bensmith@allblacks.nz",
		})

	// 从数据库查询smith先生 Selects Mr. Smith from the database
	rows, err = db.NamedQuery(`SELECT * FROM person WHERE first_name=:fn`,
		map[string]interface{}{"fn": "Bin"})

	// Named queries can also use structs.  Their bind names follow the same rules
	// as the name -> db mapping, so struct fields are lowercased and the `db` tag
	// is taken into consideration.
	//命名查询也可以使用结构。 它们的绑定名称遵循相同的规则
	//作为名称 -> db 映射，因此结构字段是小写的，并且 'db' 标签被考虑在内。
	rows, err = db.NamedQuery(`SELECT * FROM person WHERE first_name=:first_name`, jason)

	//批量插入结构体格式的 batch insert
	// batch insert with structs
	personStructs := []Person{
		{FirstName: "Ardie", LastName: "Savea", Email: "asavea@ab.co.nz"},
		{FirstName: "Sonny Bill", LastName: "Williams", Email: "sbw@ab.co.nz"},
		{FirstName: "Ngani", LastName: "Laumape", Email: "nlaumape@ab.co.nz"},
	}

	_, err = db.NamedExec(`INSERT INTO person (first_name, last_name, email)
        VALUES (:first_name, :last_name, :email)`, personStructs)

	//批量插入字典格式的 batch insert with maps
	personMaps := []map[string]interface{}{
		{"first_name": "Ardie", "last_name": "Savea", "email": "asavea@ab.co.nz"},
		{"first_name": "Sonny Bill", "last_name": "Williams", "email": "sbw@ab.co.nz"},
		{"first_name": "Ngani", "last_name": "Laumape", "email": "nlaumape@ab.co.nz"},
	}

	_, err = db.NamedExec(`INSERT INTO person (first_name, last_name, email)
        VALUES (:first_name, :last_name, :email)`, personMaps)

	return &people, nil

}
