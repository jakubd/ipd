package ipd

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func basicTestDb(t *testing.T, givenDbName string) {
	db, err := OpenMaxmindDb(givenDbName)
	assert.NoError(t, err)
	db.Close()
}

func TestOpenMaxmindDb(t *testing.T) {
	basicTestDb(t, "GeoLite2-ASN.mmdb")
	basicTestDb(t, "GeoLite2-Country.mmdb")

	_, err := OpenMaxmindDb("df")
	assert.Error(t, err)
}

func TestResolve(t *testing.T) {
	TestDomain := "one.one.one.one"
	good, err := SimpleResolveDomain(TestDomain)
	assert.NoError(t, err)
	assert.Contains(t, []string{"1.1.1.1", "1.0.0.1"}, good)

	JunkDomain := "asdfasdfasdfasdf4f3asdfasdf4fasdfasd"
	junk, err := SimpleResolveDomain(JunkDomain)
	assert.Error(t, err)
	assert.Equal(t, "", junk)
}

func TestLookup(t *testing.T) {
	testIp := "81.2.69.142"
	info, err := Lookup(testIp)
	assert.NoError(t, err)
	assert.Equal(t, info.Input, testIp)
	assert.Equal(t, info.ASNum, 20712)
	assert.Equal(t, info.ASName, "Andrews & Arnold Ltd")
	assert.Equal(t, info.CountryCode, "GB")
	assert.Equal(t, info.CountryName, "United Kingdom")

	info, err = Lookup(testIp, true)
	assert.NoError(t, err)
	assert.Equal(t, info.Input, testIp)
	assert.Equal(t, info.ASNum, 20712)
	assert.Equal(t, info.ASName, "Andrews & Arnold Ltd")
	assert.Equal(t, info.CountryCode, "GB")
	assert.Equal(t, info.CountryName, "United Kingdom")

	badIp := "hamsammich"
	badInfo, err := Lookup(badIp)
	assert.Error(t, err)
	assert.Equal(t, badInfo.Input, badIp)
	assert.Equal(t, badInfo.ASNum, -1)
	assert.Equal(t, badInfo.ASName, "")
	assert.Equal(t, badInfo.CountryCode, "")
	assert.Equal(t, badInfo.CountryName, "")

	badIp2 := "81.2.69as.142"
	badInfo2, err := Lookup(badIp2)
	assert.Error(t, err)
	assert.Equal(t, badInfo2.Input, badIp2)
	assert.Equal(t, badInfo2.ASNum, -1)
	assert.Equal(t, badInfo2.ASName, "")
	assert.Equal(t, badInfo2.CountryCode, "")
	assert.Equal(t, badInfo2.CountryName, "")

	urlInput := "https://one.one.one.one/whatever"
	urlInfo, err := Lookup(urlInput, true)
	assert.NoError(t, err)
	assert.Contains(t, []string{"1.1.1.1", "1.0.0.1"}, urlInfo.IP.String())
	assert.Equal(t, urlInfo.ASNum, 13335)
	assert.Equal(t, urlInfo.ASName, "CLOUDFLARENET")

	domainInput := "one.one.one.one"
	domainInfo, err := Lookup(domainInput, true)
	assert.NoError(t, err)
	fmt.Println(domainInfo)

}
