package main

import (
	"flag"
)

var (
	dir           = flag.String("dir", "tls", "Directory to store certificates")
	host          = flag.String("host", "localhost,127.0.0.1,122.51.240.88", "Comma-separated hostnames and IPs to generate a certificate for")
	rsaBits       = flag.Int("rsa-bits", 2048, "Size of RSA key to generate. Ignored if `ecdsa-curve` is set")
	ecdsaCurve    = flag.String("ecdsa-curve", "P256", "ECDSA curve to use to generate a key. Valid values are P224, P256 (recommended), P384, P521")
	useEd25519Key = flag.Bool("use-ed25519", false, "Generate an Ed25519 key or not")

	// mode     = flag.String("mode", "flag", "cmd flag or yaml")
	// confFile = flag.String("conf-file", "ca.yml", "path of ca information file")

	rootOrg         = flag.String("root-org", "root", "root organization name")
	serverOrg       = flag.String("server-org", "server", "server organization name")
	clientOrg       = flag.String("client-org", "client", "client organization name")
	rootStartDate   = flag.String("root-start-date", "2022-01-01", "Creation date of root certificate formatted as 2020-07-22")
	rootExpire      = flag.Int("root-expire", 3650, "Days of duration that root certificate is valid")
	serverStartDate = flag.String("server-start-date", "2022-01-01", "Creation date server certificate formatted as 2020-07-22")
	serverExpire    = flag.Int("server-expire", 3650, "Day of duration that server certificate is valid")
	clientStartDate = flag.String("client-start-date", "2022-01-01", "Creation date client certificate formatted as 2020-07-22")
	clientExpire    = flag.Int("client-expire", 3650, "Day of duration that client certificate is valid")
)

func main() {

	// c, err := ioutil.ReadFile(*confFile)
	// if err != nil {
	// 	log.Fatalln(err)
	// }

	// var ci CaInfo
	// err = yaml.UnmarshalStrict(c, &ci)
	// if err != nil {
	// 	log.Fatalln(err)
	// }

	flag.Parse()

	GenRootCA()

	GenServerCA()

	GenClientCA()

}
