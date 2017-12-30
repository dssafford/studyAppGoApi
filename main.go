package main

import (
	// "bytes"
	"database/sql"
	"fmt"
	"net/http"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"os"
	"log"
	"net/url"
	"time"
	
)

func main() {
	//var db *sql.DB
	//var err error

	databaseUrl := os.Getenv("DATABASE_URL")

	if databaseUrl == "" {
		fmt.Println("*** Using local database ***")
		databaseUrl = "mysql://root:mypassword@0.0.0.0:9015/studyApp"
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
	type Entry struct {
		Id int
		Date_Added time.Time
		Project string
		File_Directory string
		Machine string
		Technology	string
		Version int
		Comments string
		Is_Active bool
	}
	router := gin.Default()
	// - No origin allowed by default
	// - GET,POST, PUT, HEAD methods
	// - Credentials share disabled
	// - Preflight requests cached for 12 hours
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:4500"}
	// config.AddAllowOrigins("http://facebook.com")
	// config.AllowOrigins == []string{"http://google.com", "http://facebook.com"}

	router.Use(cors.New(config))
	//router.Run()

	// router.GET("/data", func(c *gin.Context) {
	// 	var result gin.H

	// 	result = gin.H{"databaseUrl:": databaseUrl}

	// 	c.JSON(http.StatusOK, result )

	// })

	// GET a person detail
	// router.GET("/entry/:id", func(c *gin.Context) {
	// 	var (
	// 		entry Entry
	// 		result gin.H
	// 	)
	// 	id := c.Param("id")
	// 	row := db.QueryRow("select id, project, file_directory, machine, technology, version, comments from journal where id = ?;", id)
	// 	// err = row.Scan(&entry.id, &entry.date_added, &entry.project, &entry.file_directory, &entry.machine, &entry.technology, &entry.version, &entry.comments, &entry.is_active)
	// 	err = row.Scan(&entry.id, &entry.project, &entry.file_directory, &entry.machine, &entry.technology, &entry.version, &entry.comments)
	// 	if err != nil {
	// 		// If no results send null
	// 		result = gin.H{
	// 			"result": nil,
	// 			"count":  0,
	// 		}
	// 	} else {
	// 		result = gin.H{
	// 			"result": entry,
	// 			"count":  1,
	// 		}
	// 	}
	// 	c.JSON(http.StatusOK, result)
	// })

	// GET all persons
	router.GET("/entries", func(c *gin.Context) {
		var (
			entry  Entry
			entrys []Entry
		)
		rows, err := db.Query("select Id, Date_Added, Project, File_Directory, Machine, Technology, Version, Comments, Is_Active from journal;")
		if err != nil {
			fmt.Print(err.Error())
		}
		for rows.Next() {
			// err = rows.Scan(&entry.id, &entry.date_added, &entry.project, &entry.file_directory, &entry.machine, &entry.technology, &entry.version, &entry.comments, &entry.is_active)
			err = rows.Scan(&entry.Id, &entry.Date_Added, &entry.Project, &entry.File_Directory, &entry.Machine, &entry.Technology, &entry.Version, &entry.Comments, &entry.Is_Active)
			entrys = append(entrys, entry)
			
			if err != nil {
				fmt.Print(err.Error())
			}
		}
		defer rows.Close()
		c.JSON(http.StatusOK, gin.H{
			"result": entrys,
			"count":  len(entrys),
		})
	})

	// // POST new person details
	// router.POST("/entry", func(c *gin.Context) {
	// 	var buffer bytes.Buffer
	// 	project := c.PostForm("project")
	// 	date_added := c.PostForm("date_added")
	// 	file_directory := c.PostForm("file_directory")
	// 	machine := c.PostForm("machine")
	// 	technology := c.PostForm("technology")
	// 	version := c.PostForm("version")
	// 	comments := c.PostForm("comments")
	// 	is_active := c.PostForm("is_active")
	// 	stmt, err := db.Prepare("insert into journal (project, date_added, file_directory, machine, technology, version, comments, is_active) values(?,?,?,?,?,?,?,?) where id= ?;")
	// 	if err != nil {
	// 		fmt.Print(err.Error())
	// 	}
	// 	_, err = stmt.Exec(project, date_added, file_directory, machine, technology, version, comments, is_active)

	// 	if err != nil {
	// 		fmt.Print(err.Error())
	// 	}

	// 	// Fastest way to append strings
	// 	buffer.WriteString(project)
	// 	buffer.WriteString(" ")
	// 	buffer.WriteString(machine)
	// 	defer stmt.Close()
	// 	name := buffer.String()
	// 	c.JSON(http.StatusOK, gin.H{
	// 		"message": fmt.Sprintf(" %s successfully created", name),
	// 	})
	// })

	// // PUT - update a person details
	// router.PUT("/entry", func(c *gin.Context) {
	// 	var buffer bytes.Buffer
	// 	id := c.Query("id")
	// 	project := c.PostForm("project")
	// 	date_added := c.PostForm("date_added")
	// 	file_directory := c.PostForm("file_directory")
	// 	machine := c.PostForm("machine")
	// 	technology := c.PostForm("technology")
	// 	version := c.PostForm("version")
	// 	comments := c.PostForm("comments")
	// 	is_active := c.PostForm("is_active")
	// 	stmt, err := db.Prepare("update journal set project= ?, date_added= ?, file_directory= ?, machine=?, technology=?,version=?, comments=?, is_active=? where id= ?;")
	// 	if err != nil {
	// 		fmt.Print(err.Error())
	// 	}
	// 	_, err = stmt.Exec(project, date_added, file_directory, machine, technology, version, comments, is_active, id)
	// 	if err != nil {
	// 		fmt.Print(err.Error())
	// 	}

	// 	// Fastest way to append strings
	// 	buffer.WriteString(project)
	// 	buffer.WriteString(" ")
	// 	buffer.WriteString(machine)
	// 	defer stmt.Close()
	// 	updatedEntry := buffer.String()
	// 	c.JSON(http.StatusOK, gin.H{
	// 		"message": fmt.Sprintf("Successfully updated to %s", updatedEntry),
	// 	})
	// })

	// // Delete resources
	// router.DELETE("/entry", func(c *gin.Context) {
	// 	id := c.Query("id")
	// 	stmt, err := db.Prepare("delete from journal where id= ?;")
	// 	if err != nil {
	// 		fmt.Print(err.Error())
	// 	}
	// 	_, err = stmt.Exec(id)
	// 	if err != nil {
	// 		fmt.Print(err.Error())
	// 	}
	// 	c.JSON(http.StatusOK, gin.H{
	// 		"message": fmt.Sprintf("Successfully deleted entry: %s", id),
	// 	})
	// })


	if err := http.ListenAndServe(fmt.Sprintf(":%v", getPort()), router); err != nil {
		log.Fatalln(err)
	}

}

func getPort() string {
	if configuredPort := os.Getenv("PORT"); configuredPort == "" {
		fmt.Println("usimg port = 3000")
		return "3000"
	} else {
		fmt.Println("using port = " + configuredPort)
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