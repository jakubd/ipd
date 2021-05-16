package ipd

import (
	"fmt"
	"github.com/asaskevich/govalidator"
	"net"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"

	"github.com/oschwald/geoip2-golang"
)

func OutputLookup(givenInput string, intel bool, resolve ...bool) {
	ipinfo, err := Lookup(givenInput, resolve[0])
	status := "good_ip"
	if err != nil {
		status = "bad_ip"
	}
	var record []string

	var ipstr = ""
	if ipinfo.IP == nil {
		ipstr = ""
	} else {
		ipstr = ipinfo.IP.String()
	}

	record = []string{ipinfo.Input, ipstr, ipinfo.CountryCode, ipinfo.ASName, ipinfo.ASNumStr, status}

	if intel {

		intelrecord := []string{
			" https://censys.io/ipv4/" + ipstr + " ",
			" https://www.shodan.io/host/" + ipstr + " ",
			" https://bgp.he.net/" + ipinfo.ASNumStr + " ",
		}
		record = append(record, intelrecord...)
	}
	fmt.Println(strings.Join(record, ","))
}

// IPInfo is the struct of enriched geoip info
type IPInfo struct {
	Input       string // given input string for a lookup
	IP          net.IP // net.IP representation of the IP string or input
	ASNum       int    // Autonomous system number as int
	ASNumStr    string // Autonomous system number as string prefixed with "AS"
	ASName      string // Autonomous system name
	CountryCode string // ISO Country Code
	CountryName string // Country name
}

// IsFileInMaxmindDir will check if the givenFile is in the Maxmind dir and report back.
// If false will output to errs tream
func IsFileInMaxmindDir(givenFile string) bool {
	if _, err := os.Stat(filepath.Join(GetMaxmindDirectory(), givenFile)); os.IsNotExist(err) {
		_ = fmt.Errorf("can not find neccesary file '%s' in dir %s", givenFile, GetMaxmindDirectory())
		return false
	}
	return true
}

// CheckMaxmindEnvironment will check all neccesary files in the environment needed to function.
func CheckMaxmindEnvironment() bool {
	if runtime.GOOS != "linux" {
		_ = fmt.Errorf("unsupported OS: %s", runtime.GOOS)
		return false
	}

	if _, err := os.Stat(GetMaxmindDirectory()); os.IsNotExist(err) {
		_ = fmt.Errorf("can not find maxmind directory: %s", GetMaxmindDirectory())
		return false
	}

	if !IsFileInMaxmindDir("GeoLite2-ASN.mmdb") {
		return false
	}

	if !IsFileInMaxmindDir("GeoLite2-ASN.mmdb") {
		return false
	}

	return true
}

// GetMaxmindDirectory will return the expected directory for the maxmind db files according to OS
func GetMaxmindDirectory() string {
	switch runtime.GOOS {
	case "darwin":
		panic("MacOS is not supported")
	case "windows":
		panic("Windows is not supported")
	case "linux":
		return "/var/lib/GeoIP/"
	default:
		return "/var/lib/GeoIP/"
	}
}

// OpenMaxmindDb will open the givenDbName from the default or givenDirectory and return the Reader object
func OpenMaxmindDb(givenDbName string, givenDirectory ...string) (*geoip2.Reader, error) {
	var maxmindDirectory string
	if len(givenDirectory) == 0 {
		maxmindDirectory = GetMaxmindDirectory()
	} else {
		maxmindDirectory = givenDirectory[0]
	}

	maxmindDb, err := geoip2.Open(path.Join(maxmindDirectory, givenDbName))
	if err != nil {
		return nil, err
	}

	return maxmindDb, nil
}

// SimpleResolveDomain will lookup a domain and return an IP if possible
// TODO: cleanup
func SimpleResolveDomain(givenInput string) (string, error) {
	ips, err := net.LookupIP(givenInput)
	if err != nil || len(ips) < 1 {
		return "", err
	}
	for _, ip := range ips {
		return ip.String(), nil
	}
	return "", nil
}

// CleanupInput does some light sanitization of the givenInput in the Lookup Func.
func CleanupInput(givenInput string) string {
	return strings.ToLower(strings.TrimSpace(givenInput))
}

// Lookup will look up the givenIpStr string and return a fully parsed IPInfo struct
// if resolve is set to true then input can be domain or url
func Lookup(givenInput string, resolve ...bool) (IPInfo, error) {
	givenInput = CleanupInput(givenInput)

	parseFailed := IPInfo{
		Input:       givenInput,
		IP:          nil,
		ASNum:       -1,
		ASNumStr:    "AS0",
		ASName:      "",
		CountryCode: "",
		CountryName: "",
	}

	DoResolutions := false
	if len(resolve) > 0 && resolve[0] {
		DoResolutions = true
	}

	asnDb, err := OpenMaxmindDb("GeoLite2-ASN.mmdb")
	if err != nil {
		panic(fmt.Sprintf("No maxmind db found in: %s.  "+
			"Please download from https://dev.maxmind.com/geoip/geoip2/geolite2/ and place in dir.",
			GetMaxmindDirectory()))
	}

	countryDb, err := OpenMaxmindDb("GeoLite2-Country.mmdb")
	if err != nil {
		panic(fmt.Sprintf("No maxmind db found in: %s.  "+
			"Please download from https://dev.maxmind.com/geoip/geoip2/geolite2/ and place in dir.",
			GetMaxmindDirectory()))
	}

	defer func(asnDb *geoip2.Reader) {
		err := asnDb.Close()
		if err != nil {
			panic("Could not close ASN db.")
		}
	}(asnDb)
	defer func(countryDb *geoip2.Reader) {
		err := countryDb.Close()
		if err != nil {
			panic("Could not close country db.")
		}
	}(countryDb)

	var ip net.IP
	if DoResolutions {
		var answer = ""
		if govalidator.IsIP(givenInput) {
			answer = givenInput
		} else if govalidator.IsDNSName(givenInput) {
			answer, err = SimpleResolveDomain(givenInput)
		} else {
			// assuming it is a URL and parse that out
			parsed, err := url.Parse(givenInput)
			if err != nil {
				return parseFailed, err
			}
			answer, err = SimpleResolveDomain(parsed.Host)
		}

		if err != nil || len(answer) < 1 {
			return parseFailed, err
		}

		ip = net.ParseIP(answer)

	} else {
		ip = net.ParseIP(givenInput)
	}

	asnRecord, err := asnDb.ASN(ip)
	if err != nil {
		return parseFailed, err
	}

	countryRecord, err := countryDb.Country(ip)
	if err != nil {
		return parseFailed, err
	}

	return IPInfo{
		Input:       givenInput,
		IP:          ip,
		ASNum:       int(asnRecord.AutonomousSystemNumber),
		ASNumStr:    "AS" + strconv.Itoa(int(asnRecord.AutonomousSystemNumber)),
		ASName:      asnRecord.AutonomousSystemOrganization,
		CountryCode: countryRecord.Country.IsoCode,
		CountryName: countryRecord.Country.Names["en"],
	}, nil

}
