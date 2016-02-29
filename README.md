=== Introduction
Go-RTBH is a daemon which hooks into an existing ELK stack using an AMQP
fanout exchange. It will process incoming events based on a field-based, 
regexp matched ruleset, extract the source ip address, and use BIRD to
distribute prefixes containing s/RTBH specific communities which your
network routers can then use to perform null-routing.
