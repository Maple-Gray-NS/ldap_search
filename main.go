package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/akamensky/argparse"
	"github.com/go-ldap/ldap/v3"
)

type Policy struct {
	DisplayName string
	GUID        string
	Path        string
}

type CLI_Arugment struct {
	Domain        string
	User          string
	Pass          string
	Remote_Target string
}

func main() {
	parser := argparse.NewParser("LDAP_GPO_Seacher", "A small program to locate the GPO from from its name")
	user := parser.String("u", "user", &argparse.Options{Required: true, Help: "Valid AD User"})
	password := parser.String("p", "password", &argparse.Options{Required: true, Help: "Valid AD Password"})
	domain := parser.String("d", "domain", &argparse.Options{Required: true, Help: "DNS Domain Name"})
	dc_ip := parser.String("t", "ldap-server", &argparse.Options{Required: true, Help: "IP of LDAP Server"})
	policy := parser.String("P", "policy", &argparse.Options{Required: true, Help: "Partial of the Policy Name"})
	err := parser.Parse(os.Args)
	if err != nil {
		log.Fatal(err)
	}

	base_dn := create_base_dn(*domain)

	conn, err := connect_to_ldap(fmt.Sprintf("%s@%s", *user, *domain), *password, *dc_ip)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	results, err := search_ldap(conn, *policy, base_dn)
	if err != nil {
		log.Fatal(err)
	}
	if len(results) == 0 {
		fmt.Println("results not found")
		return
	} else {
		fmt.Println("Results found")
	}

	policies := parse_entry(results)
	for p := range policies {
		fmt.Printf("Policy Name: %s\n", policies[p].DisplayName)
		fmt.Printf("Path: \\\\%s\\sysvol\\%s\\Policies\\%s", *dc_ip, *domain, policies[p].GUID)
	}

}

func create_base_dn(domain string) string {
	dn_comp := strings.Split(domain, ".")
	var base_dn string
	for d := range dn_comp {
		base_dn += fmt.Sprintf("DC=%s,", dn_comp[d])
	}
	base_dn = strings.Trim(base_dn, ",")
	return base_dn

}

func connect_to_ldap(ldap_username string, ldap_password string, dc_ip string) (*ldap.Conn, error) {
	l, err := ldap.Dial("tcp", fmt.Sprintf("%s:389", dc_ip))
	if err != nil {
		return l, err
	}

	err = l.Bind(ldap_username, ldap_password)

	if err != nil {
		return l, err
	}

	return l, nil

}

func search_ldap(conn *ldap.Conn, policy_name string, base_dn string) ([]*ldap.Entry, error) {
	filter := fmt.Sprintf("(&(objectClass=groupPolicyContainer)(displayName=*%s*))", policy_name)
	//filter := fmt.Sprintf("(&(objectClass=groupPolicyContainer))")

	search_req := ldap.NewSearchRequest(
		base_dn,
		ldap.ScopeWholeSubtree,
		ldap.NeverDerefAliases,
		0,
		0,
		false,
		filter,
		[]string{"displayName", "cn"},
		nil,
	)

	result, err := conn.SearchWithPaging(search_req, 100)
	if err != nil {
		return result.Entries, err
	}

	return result.Entries, err

}

func parse_entry(entries []*ldap.Entry) []Policy {
	var policy_slice []Policy
	for _, e := range entries {
		var new_policy Policy
		for _, attributes := range e.Attributes {
			switch attributes.Name {
			case "cn":
				new_policy.GUID = attributes.Values[0]
			case "displayName":
				new_policy.DisplayName = attributes.Values[0]
			}

		}
		policy_slice = append(policy_slice, new_policy)
	}

	return policy_slice

}

func Parse_Target(arg string) (CLI_Arugment, error) {
	var argument_obj CLI_Arugment
	regex_target := "(?:(?:([^/@:]*)/)?([^@:]*)(?::([^@]*))?@)?(.*)"
	r := regexp.MustCompile(regex_target)
	results := r.FindStringSubmatch(arg)

	if len(results) == 5 {
		argument_obj.Domain = results[1]
		argument_obj.User = results[2]
		argument_obj.Pass = results[3]
		argument_obj.Remote_Target = results[4]
	} else {
		return argument_obj, errors.New("unable to parse Credentials/Target")
	}
	return argument_obj, nil

}
