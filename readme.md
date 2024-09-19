# DNS Resolver

This is an example of a DNS resolver written in Go.

The purpose of a DNS resolver is to figure out what address is associated with what IP address. 

GOAL: This project is a fun way to learn:

How to parse a binary network protocol like DNS

How DNS works under the hood (what’s happening behind the scenes when you make a DNS query?) 

Other goal: writing under 200 lines of code for this, compared to the original Python project by Julia Evans. 

	// Example usage
	q := DNSQuery{
		Name:  domain,
		Type_: 1, // A record
		Class: 1, // IN class
	}
    
	// fix my test 
	// TODO: send my query to an actual server 
	//TODO: Read again what DNS is supposed to do