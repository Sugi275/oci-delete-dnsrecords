package main

import (
	"context"
	"fmt"

	"github.com/Sugi275/oci-env-configprovider/envprovider"
	"github.com/oracle/oci-go-sdk/common"
	"github.com/oracle/oci-go-sdk/dns"
)

func main() {
	zn := "test.enc"
	dn := "multitest.test.enc"

	client, err := dns.NewDnsClientWithConfigurationProvider(envprovider.GetEnvConfigProvider())
	if err != nil {
		panic(err)
	}

	compartmentid, err := envprovider.GetCompartmentID()
	if err != nil {
		panic(err)
	}

	getRequest := dns.GetDomainRecordsRequest{
		ZoneNameOrId:  common.String(zn),
		Domain:        common.String(dn),
		CompartmentId: common.String(compartmentid),
		Rtype:         common.String("TXT"),
	}

	// Domainに設定されているすべてのレコードを取得
	ctx := context.Background()
	domainRecords, err := client.GetDomainRecords(ctx, getRequest)
	if err != nil {
		panic(err)
	}

	if *domainRecords.OpcTotalItems == 0 {
		fmt.Println("no records")
		return
	}

	var deletehash *string
	for _, record := range domainRecords.RecordCollection.Items {
		fmt.Print("Domain: ", *record.Domain, "  ")
		fmt.Print("RecordHash: ", *record.RecordHash, "  ")
		fmt.Println("Rdata: ", *record.Rdata)

		if *record.Rdata == "\"deleteme\"" {
			deletehash = record.RecordHash
			break
		}
	}

	if deletehash == nil {
		fmt.Println("no records")
		return
	}

	// 一部のレコードを削除。rdataの中身がdeletemeのものを削除
	fmt.Println("deletehash: ", *deletehash)
	recordOperation := dns.RecordOperation{
		RecordHash: deletehash,
		Operation:  dns.RecordOperationOperationRemove,
	}

	patchRequest := dns.PatchDomainRecordsRequest{
		ZoneNameOrId: common.String(zn),
		Domain:       common.String(dn),
		PatchDomainRecordsDetails: dns.PatchDomainRecordsDetails{
			Items: []dns.RecordOperation{
				recordOperation,
			},
		},
		CompartmentId: common.String(compartmentid),
	}

	response, err := client.PatchDomainRecords(context.Background(), patchRequest)
	if err != nil {
		panic(err)
	}

	fmt.Println(response)
}
