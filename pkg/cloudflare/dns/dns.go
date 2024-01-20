package dns

type DNSRecord struct {
	Id         string
	Zone_id    string
	Zone_name  string
	Name       string
	Dns_type   string
	Content    string
	Proxiable  bool
	Proxied    bool
	Comment    string
	ModifiedOn string
}
