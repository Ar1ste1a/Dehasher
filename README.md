# Dehasher
## A cli tool built for interaction with the Dehash API

<img src="https://img.wanman.io/fUSu0/tIFUJOMu64.png/raw" alt="Ar1ste1a" title="Ar1ste1a Offensive Security">

# Features
- Output Format Control
- Request Limiting
- Record Limiting
- Regular Expression Handling
- Exact Match Handling
- Error Handling
- Credential Dumping
- Intelligent Token Usage
# Options
```bash-session
usage: Dehasher [-h --help] {-k --key}  {-a --authorized-email} [-h --help]  [-m --max-records]  [-r --max-requests]  [-B --print-balance]  [-X --exact-match]  [-R --regex-match]  [-t --list-tokens]  [-o --output-file-name]  [-T --output-txt]  [-J --output-json]  [-Y --output-yaml]  [-x --output-xml]  [-A --query-all-fields]  [-u --query-username]  [-U --username-query]  [-e --query-email]  [-E --email-query]  [-i --query-ip-address]  [-I --ip-address-query]  [-p --query-password]  [-P --password-query]  [-q --query-hashed-password]  [-Q --hashed-password-query]  [-n --query-name]  [-N --name-query]  [-C --creds-only]

Dehashed Tool

options:
   -h --help                                      show this help message and exit
   -m --max-records                           Maximum amount of records to return
   -r --max-requests                           Maximum number of requests to make
   -B --print-balance                      Print remaining balance after requests
   -X --exact-match                                  Use Exact Matching on fields
   -R --regex-match                                  Use Regex Matching on fields
   -t --list-tokens                           List the number of tokens remaining
   -o --output-file-name                                File to output results to
   -T --output-txt                                            Output to text file
   -J --output-json                                           Output to JSON file
   -Y --output-yaml                                           Output to YAML file
   -x --output-xml                                             Output to XML file
   -A --query-all-fields                                        Return All Fields
   -u --query-username                                           Return Usernames
   -U --username-query                                             Username Query
   -e --query-email                                                 Return Emails
   -E --email-query                                                   Email Query
   -i --query-ip-address                                      Return IP Addresses
   -I --ip-address-query                                         IP Address Query
   -p --query-password                                           Return Passwords
   -P --password-query                                             Password Query
   -q --query-hashed-password                             Return Hashed Passwords
   -Q --hashed-password-query                               Hashed Password Query
   -n --query-name                                                   Return Names
   -N --name-query                                                     Name Query
   -C --creds-only                                        Return Credentials Only
   -k --key                                                               API Key
   -a --authorized-email                Email to pair with key for authentication


v1.0
```

# Getting Started

To begin, clone the repository
``` bash-session
git clone https://github.com/Ar1ste1a/Dehasher.git
cd Dehasher
go build dehasher.go
```

# Crafting a query

## Simple Query
``` go
# Provide credentials for emails matching @target.com
dehasher -k ddq<redacted> -a ar1ste1a@domain.tld -E @target.com
```

## Simple Credentials Query
``` go
# Provide credentials for emails matching @target.com
dehasher -k ddq<redacted> -a ar1ste1a@domain.tld -E @target.com -C
```

## Simple Query Returning Balance
``` go
# Provide credentials for emails matching @target.com
dehasher -k ddq<redacted> -a ar1ste1a@domain.tld -E @target.com -C -B
```

## Regex Query
``` go
# Return matches for emails matching this given regex query
# -R e: Specify the '-E' field as a regex entry
dehasher -k ddq<redacted> -a ar1ste1a@domain.tld -E '[a-zA-Z0-9]+(?:\.[a-zA-Z0-9]+)?@target.com' -C -B -R e
```

## Exact Match Query
``` go
# Return matches for usernames exactly matching "admin"
# -X u: Specify the '-U' field as an exact match entry
dehasher -k ddq<redacted> -a ar1ste1a@domain.tld -C -B -U admin -X u
```

## Output Text (default JSON)
``` go
# Return matches for usernames exactly matching "admin" and write to text file 'admins_file.txt'
dehasher -k ddq<redacted> -a ar1ste1a@domain.tld -C -B -U admin -X u -T -o admins_file
```

## Output YAML
``` go
# Return matches for usernames exactly matching "admin" and write to yaml file 'admins_file.yaml'
dehasher -k ddq<redacted> -a ar1ste1a@domain.tld -C -B -U admin -X u -Y -o admins_file
```

## Output XML
``` go
# Return matches for usernames exactly matching "admin" and write to xml file 'admins_file.xml'
dehasher -k ddq<redacted> -a ar1ste1a@domain.tld -C -B -U admin -X u -x -o admins_file
```

