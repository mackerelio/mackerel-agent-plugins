package main

import (
	"os"
    "net/http"
    "io/ioutil"
    "strconv"
    "log"
    "split"
	"github.com/codegangsta/cli"
)


// General error handling
type GeneralError struct {
    Error string
}


// Get metrics main function
func getMetrics( c *cli.Context ) {
    status, err := getApache2Metrics(
        c.String( "http_host" ),
        uint16( c.Int( "http_port" ) ),
        c.String( "status_page" ) )
    if err != nil {
        log.Fatal( err.Error )
    }
    status = status
}


// parsing metrics from server-status?auto
func parseApache2Status( str string )( string, error ) {
    const Params := map[string]string{
        "Total Accesses": "requests",
        "Total kBytes": "bytes_sent",
        "CPULoad": "cpu_load",
        "BusyWorkers": "busy_workers",
        "IdleWorkers": "idle_workers"
        }
    datas = strings.Sprit( str, "\n" )

}


// Getting apache2 status from server-status module data.
func getApache2Metrics( host string, port uint16, path string )( string, error ){
    uri := "http://" + host + ":" + strconv.FormatUint( uint64( port ), 10 ) + path
    resp, err := http.Get( uri )
    if err != nil {
        return "", err
    }
    status, err := ioutil.ReadAll( resp.Body )
    resp.Body.Close()
    if err != nil {
        return "", err
    }
    return string( status[:] ), nil
}


// main
func main() {
	app := cli.NewApp()
	app.Name = "apache2_metrics"
	app.Version = Version
	app.Usage = "Get metrics from apache2."
	app.Author = "Yuichiro Saito"
	app.Email = "saito@heartbeats.jp"
	app.Flags = Flags
    app.Action = getMetrics

	app.Run(os.Args)
}
