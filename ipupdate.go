package main

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	dns "google.golang.org/api/dns/v1"

	"golang.org/x/oauth2/google"
)

type (
	data struct {
		IP string
	}
)

func getIP() string {
	resp, getErr := http.Get("https://api.ipify.org?format=json")
	if getErr != nil {
		log.Fatal(getErr)
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	d := data{}
	if err := json.Unmarshal(body, &d); err != nil {
		log.Println("Error unmarshalling data from ipify")
	}
	return d.IP
}

func dnsClientAuth() *dns.Service {
	// Use oauth2.NoContext if there isn't a good context to pass in.
	ctx := context.TODO()

	client, err := google.DefaultClient(ctx, dns.NdevClouddnsReadwriteScope)
	if err != nil {
		log.Fatal("Error getting DNS scope")
	}
	dnsService, err := dns.New(client)
	if err != nil {
		log.Fatal("Error creating DNS client")
	}
	return dnsService
}

func getCurrentSet(service *dns.Service) *dns.ResourceRecordSet {
	var currentSet *dns.ResourceRecordSet
	rrs := dns.NewResourceRecordSetsService(service)
	rrslc := rrs.List(os.Getenv("PROJECT_ID"), os.Getenv("ZONE"))
	response, err := rrslc.Do()
	if err != nil {
		log.Fatal("Error getting ManagedZone")
	}
	sets := response.Rrsets
	for _, set := range sets {
		if set.Name == os.Getenv("ADDR") {
			currentSet = set
		}
	}
	return currentSet
}

func getCurrentHome(service *dns.Service) string {
	set := getCurrentSet(service)
	ip := set.Rrdatas[0]
	return string(ip)
}

func updateIP(service *dns.Service, currentSet *dns.ResourceRecordSet, actualIP string) {
	changesService := dns.NewChangesService(service)
	change := &dns.Change{}
	newSet := &dns.ResourceRecordSet{}
	newSet.Kind = "dns#resourceRecordSet"
	newSet.Rrdatas = []string{actualIP}
	newSet.Name = os.Getenv("ADDr")
	newSet.Ttl = 300
	newSet.Type = "A"

	change.Deletions = []*dns.ResourceRecordSet{currentSet}
	change.Additions = []*dns.ResourceRecordSet{newSet}

	changeCreateCall := changesService.Create(os.Getenv("PROJECT_ID"), os.Getenv("ZONE"), change)
	ch, err := changeCreateCall.Do()
	if err != nil {
		log.Println("Error creating change")
		log.Fatal(err)
	}
	log.Println(ch)
}

func ipUpdate() {
	ip := getIP()
	log.Println("Actual home: " + ip)
	service := dnsClientAuth()
	currentHome := getCurrentHome(service)
	if currentHome != ip {
		log.Println("Home IP address needs updated")
		updateIP(service, getCurrentSet(service), ip)
	} else {
		log.Println("Current home IP address is correct")
	}
}

func doEvery(d time.Duration, f func()) {
	f()
	for x := range time.Tick(d) {
		log.Println(x)
		f()
	}
}

func checkConfig() {
	for _, envVar := range []string{"PROJECT_ID", "ZONE", "ADDR"} {
		val := os.Getenv(envVar)
		if val == "" {
			log.Fatal("Environment variable " + envVar + " needs to be set")
		}
	}
	//Address needs to end with '.'
	addr := os.Getenv("ADDR")
	if !strings.HasSuffix(".", addr) {
		addr = addr + "."
	}
	os.Setenv("ADDR", addr)

	//Check for the Google Credential file (default docker location is /data/creds.json)
	if os.Getenv("GOOGLE_APPLICATION_CREDENTIALS") == "" {
		os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", "/data/googlecreds.json")
	}
	_, err := os.Stat(os.Getenv("GOOGLE_APPLICATION_CREDENTIALS"))
	if err != nil {
		log.Fatal("GOOGLE_APPLICATION_CREDENTIALS file missing (defauls to /data/googlecreds.json)")
	}
}

func main() {
	checkConfig()
	doEvery(5*time.Hour, ipUpdate)
}
