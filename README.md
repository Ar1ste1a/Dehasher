# Dehasher
## A cli tool built for interaction with the Dehash API

<div align="center">
    <img src="https://img.wanman.io/fUSu0/SaCUyEMe87.png/raw" style="width: 350px; height: auto" alt="Ar1ste1a" title="Ar1ste1a Offensive Security">
</div>

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
usage: Dehasher [-h --help] {-k --key}  {-a --authorized-email} [-h --help]  [-m --max-records]  [-r --max-requests]  [-B --print-balance]  [-X --exact-match]  [-R --regex-match]  [-t --list-tokens]  [-o --output-file-name]  [-T --output-txt]  [-J --output-json]  [-Y --output-yaml]  [-x --output-xml]  [-U --username-query]  [-E --email-query]  [-I --ip-address-query]  [-P --password-query]  [-Q --hashed-password-query]  [-N --name-query]  [-C --creds-only]

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
   -U --username-query                                             Username Query
   -E --email-query                                                   Email Query
   -I --ip-address-query                                         IP Address Query
   -P --password-query                                             Password Query
   -Q --hashed-password-query                               Hashed Password Query
   -N --name-query                                                     Name Query
   -C --creds-only                                        Return Credentials Only
   -k --key                                                               API Key
   -a --authorized-email                Email to pair with key for authentication


v1.0
```

# Sample Run
```bash-session
-k ddq<redacted> -a ar1ste1a@<redacted> -E @example.com -C -o example_creds
Making 3 Requests for 10000 Records (30000 Total)
        [*] Performing Request...
        [*] Retrieved 60 Records
[-] Not Enough Entries, ending queries
[+] Discovered 60 Records
        [*] Writing entries file: example_creds.json
                [*] Success

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

