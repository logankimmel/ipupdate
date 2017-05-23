package main

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"

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
	rrslc := rrs.List("lkimmel-1069", "lkimmel")
	response, err := rrslc.Do()
	if err != nil {
		log.Fatal("Error getting ManagedZone")
	}
	sets := response.Rrsets
	for _, set := range sets {
		if set.Name == "home.lkimmel.com." {
			currentSet = set
			// set = set.Rrdatas[0]
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
	newSet.Name = "home.lkimmel.com."
	newSet.Ttl = 300
	newSet.Type = "A"

	change.Deletions = []*dns.ResourceRecordSet{currentSet}
	change.Additions = []*dns.ResourceRecordSet{newSet}

	changeCreateCall := changesService.Create("lkimmel-1069", "lkimmel", change)
	ch, err := changeCreateCall.Do()
	if err != nil {
		log.Println("Error creating change")
		log.Fatal(err)
	}
	log.Println(ch)
}

func main() {
	os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", "lkimmel-f071bfdb867f.json")

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
