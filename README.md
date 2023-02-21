# ldap_search


A tool I built to assist me with finding GPO locations based off their name. It is not perfect.

```
go run main.go -u maple -p Welcome9 -d netspi.local -t 192.168.1.1 -P 'Configure Local' -h
usage: LDAP_GPO_Seacher [-h|--help] -u|--user "<value>" -p|--password "<value>"
                        -d|--domain "<value>" -t|--ldap-server "<value>"
                        -P|--policy "<value>"

                        A small program to locate the GPO from from its name

Arguments:

  -h  --help         Print help information
  -u  --user         Valid AD User
  -p  --password     Valid AD Password
  -d  --domain       DNS Domain Name
  -t  --ldap-server  IP of LDAP Server
  -P  --policy       Partial of the Policy Name
  ```
  
 ```
 go run main.go -u maple -p Welcome9 -d netspi.local -t 192.168.1.1 -P 'Configure Local'
  
Results found
Policy Name: Configure Local Accounts
Path: \\192.168.1.1\sysvol\netspi.local\Policies\{EDA4933B-A8FD-472C-AF02-340F93BEAE76}
```
