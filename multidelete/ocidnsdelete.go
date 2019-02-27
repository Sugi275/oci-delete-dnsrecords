package main

import (
	"context"
	"fmt"
	"time"

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

	response, err := client.PatchDomainRecords(ctx, patchRequest)
	if err != nil {
		panic(err)
	}

	// 削除されたことを確認。PatchDomainRecordsは非同期に削除されるため
	for yes := true; yes; {
		yes, err = existRecord(ctx, client, getRequest, *deletehash)

		if err != nil {
			panic(err)
		}

		time.Sleep(3 * time.Second)
	}

	fmt.Println(response)
}

func existRecord(ctx context.Context, client dns.DnsClient, getRequest dns.GetDomainRecordsRequest, recordHash string) (bool, error) {
	domainRecords, err := client.GetDomainRecords(ctx, getRequest)
	if err != nil {
		panic(err)
	}

	for _, record := range domainRecords.RecordCollection.Items {
		if *record.RecordHash == recordHash {
			return true, nil
		}
	}

	return false, nil
}
