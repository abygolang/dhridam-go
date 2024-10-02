package main

import (
	// "fmt"
	"html/template" //for safely parsing and rendering HTML tempelates. functions in this package parse the tempelate and execute the templates.
	"log"
	"net/http"
	// "strconv"
)

// Define a home handler function which writes a byte slice containing 
// hello from dhridam as response body
// home handler base signature is changed so that it is defined as a method againt appliication
func  (app *application) home(w http.ResponseWriter, r*http.Request){//(app *application)
	if r.URL.Path!="/"{
		http.NotFound(w,r)
		return
	}
	// use template.ParseFiles() function to read the template file into a template set.
	// if there is an error, we log the detailed error message and use the http.Error() function to send a generic 500 Internal server Error response to the user.
	// As we add more pages to this web application there will e some shared, HTML markup that we want to include on every page  - like header navigation and metadata inside the <head> HTML element.
	// To prevent duplication, its a good idea to create a base(or master) tempelate which contains this shared content , which we can compose with the page specific markup for individual pages. 
	files:=[]string{
		"./ui/html/base.html",
		"./ui/html/pages/home.html",
	}
	//parse single files
	// ts, err:= template.ParseFiles("./ui/html/pages/home.html")
	// parse multi template files
	ts, err:=template.ParseFiles(files...)
	if err!=nil{
		log.Println(err.Error())
		http.Error(w, "internal server Error", 500)
		return
	}
	// We then use Execute() method on the template set to write the template content as response bosy
	// the second parameter of Execute represents ant dynamic data that er want to pass in, if no it's nil
	// err=ts.Execute(w, nil)
	err = ts.ExecuteTemplate(w, "base", nil)
	if err!=nil{
		log.Println(err.Error())
		http.Error(w, "Internal Server Error", 500)
	}
	// w.Write([]byte("Hello from dhridam"))
}

func (app *application) snippetView(w http.ResponseWriter, r*http.Request) {
	w.Write([]byte("Display a specific Snippet"))
}

func (app *application) postreq(w http.ResponseWriter, r*http.Request){
	if r.Method!="POST"{
		//use w.WriteHeader() method to send a 405 status code and
		// w.Write() method to write a "Method Not Allowed" and return from the function so that subsequent code is not executed
		w.Header().Set("Allow", "POST") // to add an Allow post header to the response heaser map. first parameter is the header name and second is the header value.
		w.WriteHeader(405) // if you don't call w.WriteHeader() explicitly, it will auto matically send a 200 ok status code to user. so need to call w.Writeheader() before any call to w.Write().
		w.Write([]byte("Method Not Allowed"))
		return
	}
}
// http.ResponseWriter parameter provides methods for assembling HTTP response and sending it to user
// http.Request parameter is pointer to a struct which holds information about current requests(like HTTp method & URL being requested)