package main

import "encoding/json"
import "fmt"
import "flag"
import "log"
import "os"
import "net/http"
import "bytes"
import "io/ioutil"

var rpc = flag.String("rpc", "http://127.0.0.1:6800/jsonrpc", "Aria2 RPC server address")
var cookie = flag.String("cookie", "", "Cookies")
var dir = flag.String("dir", "", "Saved dest directory (server side)")
var out = flag.String("out", "", "Saved output file name")
var split = flag.Int("split", 15, "One file N connections")
var server = flag.Int("server", 15, "One server N connections")
var referer = flag.String("referer", "", "Set referer")
var secret = flag.String("secret", "", "Set token for authorization")
var user = flag.String("user", "", "[DEPRECATED]Set user name")
var pw = flag.String("passwd", "", "[DEPRECATED]Set password")
var ua = flag.String("ua", "Mozilla/5.0 (X11; Linux; rv:5.0) Gecko/5.0 Firefox/5.0", "Set user agent")

// var session = flag.String("session-dir", "", "Directory for session file (server side)")

func main() {
	flag.Parse()
	URIs := flag.Args()
	if len(URIs) == 0 {
		flag.PrintDefaults()
		os.Exit(0)
	}
	params := makeParamsArry(URIs)

	jsonreq, err := makeJsonStruct(params)
	if err != nil {
		log.Fatal(err)
	}
	// fmt.Println(string(jsonreq))
	os.Stdout.Write(addTask(*rpc, jsonreq))
	os.Stdout.Write([]byte{'\n'})
}

func makeParamsArry(uris []string) []interface{} {

	output := make([]interface{}, 0, 2)

	if *secret != "" {
		output = append(output, "token:"+*secret)
	}

	output = append(output, uris)

	opts := make(map[string]interface{}, 11)
	if *dir != "" {
		opts["dir"] = *dir
	}
	if *out != "" {
		opts["out"] = *out
		// if *session != "" {
		// 	opts["save-session"] = filepath.Join(*session, *out+".session")
		// }
	}

	if *cookie != "" {
		opts["header"] = []string{fmt.Sprintf("Cookie: %s", *cookie)}
	}

	if *referer != "" {
		opts["referer"] = *referer
	}

	if *secret != "" {
		opts["rpc-secret"] = secret
	} else if *user != "" && *pw != "" {
		opts["rpc-user"] = user
		opts["rpc-passwd"] = pw
	}

	opts["continue"] = "true"
	opts["max-connection-per-server"] = *server
	opts["split"] = *split
	opts["min-split-size"] = "5M"
	opts["user-agent"] = *ua

	output = append(output, opts)

	return output
}

func makeJsonStruct(params []interface{}) ([]byte, error) {
	output := make(map[string]interface{}, 4)
	output["jsonrpc"] = "2.0"
	output["id"] = "meoow/aria2rpc"
	output["method"] = "aria2.addUri"
	output["params"] = params
	return json.Marshal(output)
}

func addTask(url string, json []byte) []byte {
	body := bytes.NewReader(json)
	resp, err := http.Post(*rpc, "text/plain", body)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	result, _ := ioutil.ReadAll(resp.Body)

	return result
}
