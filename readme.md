# DNS Resolver

This is an example of a DNS resolver written in Go.

The purpose of a DNS resolver is to figure out what address is associated with what IP address. 

GOAL: This project is a fun way to learn:

How to parse a binary network protocol like DNS

How DNS works under the hood (what’s happening behind the scenes when you make a DNS query?) 

Other goal: writing under 200 lines of code for this, compared to the original Python project by Julia Evans (I know, I am probably delusional.)

## Writing down the steps 
1- write the header_to_bytes function and the question_to_bytes. header takes a header returnes bytes, question takes a question returns bytes. 

2- write a build_query function that takes domaineName and recordType as arguments, I think those would be strings 

The struct DNS header has multiple fields: queryID, which I think might be a string and 4 ints numQuestions, numAnswers, numAdditionals, and numAuthorities. 

# dnsresolver
