package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"strconv"
	"text/template"

	"github.com/spf13/viper"
)

var templates = make(map[string]*template.Template)
var host string
var conn net.Conn
var err error

const ColorGrading = "ColorGrading"

func main() {
	viper.SetConfigName("setting")
	viper.SetConfigType("yaml")
	viper.SetDefault("host.ip", "localhost")
	viper.SetDefault("host.port", 8080)
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}

	host = viper.GetString("host.ip") + ":" + strconv.Itoa(viper.GetInt("host.port"))
	fmt.Println("Hello, This is RemoteControlViaGo from " + host)

	conn, err = net.Dial("udp", "172.29.49.211:8081")
	if err != nil {
		panic(fmt.Errorf("fatal error dial udp: %w", err))
	}
	defer conn.Close()

	templates[ColorGrading] = loadTemplate(ColorGrading)
	http.HandleFunc("/", handleIndex)
	http.HandleFunc("/pp", handleIsUseTemperatureType)
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(viper.GetInt("host.port")), nil))
}

func handleIndex(w http.ResponseWriter, r *http.Request) {
	if err := templates[ColorGrading].Execute(w, struct {
		Title string
		Host  string
	}{
		Host:  host,
		Title: ColorGrading,
	}); err != nil {
		log.Printf("failed to execute template: %v", err)
	}
}

func handleIsUseTemperatureType(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	bytearry, _ := ioutil.ReadAll(r.Body)
	conn.Write(bytearry)
	fmt.Println(string(bytearry))
}

func loadTemplate(name string) *template.Template {
	t, err := template.ParseFiles(
		"templates/"+name+".template",
		"templates/_header.template",
		"templates/_footer.template",
	)
	if err != nil {
		log.Fatalf("template error: %v", err)
	}
	return t
}
