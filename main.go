package main

import (
	"bytes"
	"database/sql"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"os"
	"log"

	"net/url"
)

func main() {
	//var db *sql.DB
	//var err error

	databaseUrl := os.Getenv("DATABASE_URL")

	if databaseUrl == "" {
		fmt.Println("*** Using local database ***")
		databaseUrl = "mysql://root:mypassword@0.0.0.0:9010/douggo"
		//fmt.Println("*** Using local database *** == " + )
	//} else {
	//
	//	databaseUrl := "b8dd0c2d93dec9:ba1aea79@tcp(us-cdbr-iron-east-05.cleardb.net:3306)/ad_547d2c245fcfb2b"
	//	db, err = sql.Open("mysql", databaseUrl)
	//	if err != nil {
	//		fmt.Print("xxxxxxxxxxxxx xxxxxxxxxxx database error == " + err.Error())
	//	}
	}

	url, err := url.Parse(databaseUrl)

	if err != nil {
		log.Fatalln("Error parsing DATABASE_URL", err)
	}

	fmt.Println("database URL = " + formattedUrl(url))

	db, err := sql.Open("mysql", formattedUrl(url))

	if err != nil {
		log.Fatalln("Failed to establish database connection", err)
	}


	//db, err := sql.Open("mysql", databaseUrl)
	//if err != nil {
	//	fmt.Print("xxxxxxxxxxxxx xxxxxxxxxxx database error == " + err.Error())
	//}
	defer db.Close()

	// make sure connection is available
	err = db.Ping()
	if err != nil {
		fmt.Print("database can't ping == " + err.Error())
	}
	type Person struct {
		Id         int
		First_Name string
		Last_Name  string
	}
	router := gin.Default()


	router.GET("/data", func(c *gin.Context) {
		var result gin.H

		result = gin.H{"databaseUrl:": databaseUrl}

		c.JSON(http.StatusOK, result )

	})

	// GET a person detail
	router.GET("/person/:id", func(c *gin.Context) {
		var (
			person Person
			result gin.H
		)
		id := c.Param("id")
		row := db.QueryRow("select id, first_name, last_name from person where id = ?;", id)
		err = row.Scan(&person.Id, &person.First_Name, &person.Last_Name)
		if err != nil {
			// If no results send null
			result = gin.H{
				"result": nil,
				"count":  0,
			}
		} else {
			result = gin.H{
				"result": person,
				"count":  1,
			}
		}
		c.JSON(http.StatusOK, result)
	})

	// GET all persons
	router.GET("/persons", func(c *gin.Context) {
		var (
			person  Person
			persons []Person
		)
		rows, err := db.Query("select id, first_name, last_name from person;")
		if err != nil {
			fmt.Print(err.Error())
		}
		for rows.Next() {
			err = rows.Scan(&person.Id, &person.First_Name, &person.Last_Name)
			persons = append(persons, person)
			if err != nil {
				fmt.Print(err.Error())
			}
		}
		defer rows.Close()
		c.JSON(http.StatusOK, gin.H{
			"result": persons,
			"count":  len(persons),
		})
	})

	// POST new person details
	router.POST("/person", func(c *gin.Context) {
		var buffer bytes.Buffer
		first_name := c.PostForm("first_name")
		last_name := c.PostForm("last_name")
		stmt, err := db.Prepare("insert into person (first_name, last_name) values(?,?);")
		if err != nil {
			fmt.Print(err.Error())
		}
		_, err = stmt.Exec(first_name, last_name)

		if err != nil {
			fmt.Print(err.Error())
		}

		// Fastest way to append strings
		buffer.WriteString(first_name)
		buffer.WriteString(" ")
		buffer.WriteString(last_name)
		defer stmt.Close()
		name := buffer.String()
		c.JSON(http.StatusOK, gin.H{
			"message": fmt.Sprintf(" %s successfully created", name),
		})
	})

	// PUT - update a person details
	router.PUT("/person", func(c *gin.Context) {
		var buffer bytes.Buffer
		id := c.Query("id")
		first_name := c.PostForm("first_name")
		last_name := c.PostForm("last_name")
		stmt, err := db.Prepare("update person set first_name= ?, last_name= ? where id= ?;")
		if err != nil {
			fmt.Print(err.Error())
		}
		_, err = stmt.Exec(first_name, last_name, id)
		if err != nil {
			fmt.Print(err.Error())
		}

		// Fastest way to append strings
		buffer.WriteString(first_name)
		buffer.WriteString(" ")
		buffer.WriteString(last_name)
		defer stmt.Close()
		name := buffer.String()
		c.JSON(http.StatusOK, gin.H{
			"message": fmt.Sprintf("Successfully updated to %s", name),
		})
	})

	// Delete resources
	router.DELETE("/person", func(c *gin.Context) {
		id := c.Query("id")
		stmt, err := db.Prepare("delete from person where id= ?;")
		if err != nil {
			fmt.Print(err.Error())
		}
		_, err = stmt.Exec(id)
		if err != nil {
			fmt.Print(err.Error())
		}
		c.JSON(http.StatusOK, gin.H{
			"message": fmt.Sprintf("Successfully deleted user: %s", id),
		})
	})
	if err := http.ListenAndServe(fmt.Sprintf(":%v", getPort()), router); err != nil {
		log.Fatalln(err)
	}
	//router.Run(":3000")
}

func getPort() string {
	if configuredPort := os.Getenv("PORT"); configuredPort == "" {
		return "3000"
	} else {
		return configuredPort
	}
}

func formattedUrl(url *url.URL) string {
	return fmt.Sprintf(
		"%v@tcp(%v)%v?parseTime=true",
		url.User,
		url.Host,
		url.Path,
	)
}