package ipd

import (
	"fmt"
	"net"
	"path"
	"strconv"
	"strings"

	"github.com/oschwald/geoip2-golang"
)

func OutputLookup(givenInput string, intel bool) {
	ipinfo, err := Lookup(givenInput)
	status := "good_ip"
	if err != nil {
		status = "bad_ip"
	}
	var record []string
	record = []string{ipinfo.IPStr, ipinfo.CountryCode, ipinfo.ASName, ipinfo.ASNumStr, status}

	if intel {
		intelrecord := []string{
			" https://censys.io/ipv4/" + ipinfo.IPStr + " ",
			" https://www.shodan.io/host/" + ipinfo.IPStr + " ",
			" https://bgp.he.net/" + ipinfo.ASNumStr + " ",
		}
		record = append(record, intelrecord...)
	}
	fmt.Println(strings.Join(record, ","))
}

// IPInfo is the struct of enriched geoip info
type IPInfo struct {
	IPStr       string // given IP String or input
	IP          net.IP // net.IP representation of the IP string or input
	ASNum       int    // Autonomous system number as int
	ASNumStr    string // Autonomous system number as string prefixed with "AS"
	ASName      string // Autonomous system name
	CountryCode string // ISO Country Code
	CountryName string // Country name
}

// OpenMaxmindDb will open the givenDbName from the default or givenDirectory and return the Reader object
func OpenMaxmindDb(givenDbName string, givenDirectory ...string) (*geoip2.Reader, error) {
	var maxmindDirectory string
	if len(givenDirectory) == 0 {
		maxmindDirectory = "/var/lib/GeoIP/"
	} else {
		maxmindDirectory = givenDirectory[0]
	}

	maxmindDb, err := geoip2.Open(path.Join(maxmindDirectory, givenDbName))
	if err != nil {
		return nil, err
	}

	return maxmindDb, nil
}

// Lookup will look up the givenIpStr string and return a fully parsed IPInfo struct
func Lookup(givenIpStr string) (IPInfo, error) {

	parseFailed := IPInfo{
		IPStr:       givenIpStr,
		IP:          nil,
		ASNum:       -1,
		ASNumStr:    "AS0",
		ASName:      "",
		CountryCode: "",
		CountryName: "",
	}

	asnDb, err := OpenMaxmindDb("GeoLite2-ASN.mmdb")
	if err != nil {
		return parseFailed, err
	}

	countryDb, err := OpenMaxmindDb("GeoLite2-Country.mmdb")
	if err != nil {
		return parseFailed, err
	}

	defer asnDb.Close()
	defer countryDb.Close()

	ip := net.ParseIP(givenIpStr)

	asnRecord, err := asnDb.ASN(ip)
	if err != nil {
		return parseFailed, err
	}

	countryRecord, err := countryDb.Country(ip)
	if err != nil {
		return parseFailed, err
	}

	return IPInfo{
		IPStr:       givenIpStr,
		IP:          ip,
		ASNum:       int(asnRecord.AutonomousSystemNumber),
		ASNumStr:    "AS" + strconv.Itoa(int(asnRecord.AutonomousSystemNumber)),
		ASName:      asnRecord.AutonomousSystemOrganization,
		CountryCode: countryRecord.Country.IsoCode,
		CountryName: countryRecord.Country.Names["en"],
	}, nil

}
