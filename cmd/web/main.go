package main

import(
	"flag"//for cmd arguements
	"log"
	"net/http"
	"os"
	)

// global application struct that holds application wide dependencies and 
// To Inject dependencies to handler function to use application logger instead of go's logger
type application struct
{
	errorLog *log.Logger
	infoLog *log.Logger
}
func main(){

	// Logs for debugging flags are joined using bitwise OR operator
	// log for info
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	// log for error, log.Lshortfile -> to include relevant file name and line number
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)


	//Initialize a new instance of our application struct, containing the dependencies
	app:= &application{
		errorLog:errorLog,
		infoLog:infoLog,
	}
	// define a new command line flag with name addr variable at runtime
	addr:=flag.String("addr", ":4000", "HTTP network address")
	// flag.Parse() reads the commandline args and assigns it to the vars otherwise it will always contain the address of ":4000"
	// the values returned from flag.String() is a pointer to the flag and not the value itself. so we need to dereference the pointer(prefix with * before using it)
	flag.Parse()
	
	// use the http.NewServeNux() function to initialize a new servemux, then
	// register the home function as the handler for the "/" URL pattern.
	mux:=http.NewServeMux() //will check urlpath & dispatch the HTTP request to matching handles
	// create a new http.fileserver handler to search ui/static directory for corresponding file to send to the user
	fileServer:=http.FileServer(http.Dir("./ui/static/")) 
	// fileServer:=http.FileServer(http.Dir("./demo-ui/oxer-html/")) // note subtree
	//mux.Handle() function to register the file server as the handler all URL paths that start with "/static".
	// For matching paths, we strip the "/static" prefix before the reaches the file server.
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))
	// mux.Handle("/oxer-html/", http.StripPrefix("/oxer-html", fileServer))
	
	
	mux.HandleFunc("/", app.home) // go's servemux treat URL pattern "/" like a catch all. This will handle all http requests to the server regardless of there URL path. example : http://localhost:4000/foo also give same response.
	mux.HandleFunc("/view", app.snippetView)

	//{} http.ListenAndServe() start a new web server
	// we pass 2 parameters : Tcp network address to listen on(:4000)
	// and the servemux we created
	// log.Fatal() function logs the error message and exit
	// any error returned by http.ListenAndServe() is always non-nil
	// log.Println("Starting Server on %s", *addr)}
	// err := http.ListenAndServe(*addr, mux) // pass incomming requests to servemux
	
	// initialize a new http.Server struct. We set address and fields so that the server uses same network address and routes as before 
	// And set Errorlog field so that the server now uses custom errorLog logger in the event of any problems
	srv:=&http.Server{
		Addr: *addr,
		ErrorLog: errorLog,
		Handler:mux,
	}
	infoLog.Printf("Starting Server on %s", *addr)
	err:= srv.ListenAndServe()
	// err.Fatal(err)
	errorLog.Fatal(err)
	// In production stage logs can be directed to disk files/logging service such as Splunk



}

// Sometimes network addresses written using named ports like ":http" or ":http-alt" instead of a number. This named
// port then Go will attempt to look up the relevant port number from your /etc/services file when staring the server/ return an error if match cant't be found.


// Go's serve mux supports 2 different types of URL patterns: fixed paths and subtree paths.

// fixed path - don't end with trailing slash - eg-: www.sippet.view (corresponding handles called when the request url path exactly matches the fixed path.)
// subtree path do end with trailing slash - eg -: www.home.h/ (corresponding handles called when start of the request url path matches the subtree path). 
// 			So need to restrict the "/" pattern.
// 			we want to return 404 page not found else.
// 			it's unable to change go's servemux default behaviour, but by adding a simple check in home handler fn ultimately  has the same effect.
