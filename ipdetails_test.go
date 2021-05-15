package ipd

import (
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

func TestLookup(t *testing.T) {
	testIp := "81.2.69.142"
	info, err := Lookup(testIp)
	assert.NoError(t, err)
	assert.Equal(t, info.IPStr, testIp)
	assert.Equal(t, info.ASNum, 20712)
	assert.Equal(t, info.ASName, "Andrews & Arnold Ltd")
	assert.Equal(t, info.CountryCode, "GB")
	assert.Equal(t, info.CountryName, "United Kingdom")

	badIp := "hamsammich"
	badInfo, err := Lookup(badIp)
	assert.Error(t, err)
	assert.Equal(t, badInfo.IPStr, badIp)
	assert.Equal(t, badInfo.ASNum, -1)
	assert.Equal(t, badInfo.ASName, "")
	assert.Equal(t, badInfo.CountryCode, "")
	assert.Equal(t, badInfo.CountryName, "")

	badIp2 := "81.2.69as.142"
	badInfo2, err := Lookup(badIp2)
	assert.Error(t, err)
	assert.Equal(t, badInfo2.IPStr, badIp2)
	assert.Equal(t, badInfo2.ASNum, -1)
	assert.Equal(t, badInfo2.ASName, "")
	assert.Equal(t, badInfo2.CountryCode, "")
	assert.Equal(t, badInfo2.CountryName, "")
}
