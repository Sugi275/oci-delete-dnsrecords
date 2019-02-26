package main

import (
	"context"
	"fmt"

	"github.com/Sugi275/oci-env-configprovider/envprovider"
	"github.com/oracle/oci-go-sdk/dns"
)

func main() {
	zn := "test.enc"
	dn := "_acme-challenge.test.enc"

	client, err := dns.NewDnsClientWithConfigurationProvider(envprovider.GetEnvConfigProvider())
	if err != nil {
		panic(err)
	}

	compartmentid, err := envprovider.GetCompartmentID()
	if err != nil {
		panic(err)
	}

	request := dns.DeleteDomainRecordsRequest{
		ZoneNameOrId:  &zn,
		Domain:        &dn,
		CompartmentId: &compartmentid,
	}

	ctx := context.Background()
	response, err := client.DeleteDomainRecords(ctx, request)
	if err != nil {
		panic(err)
	}
	fmt.Println(response)
}
