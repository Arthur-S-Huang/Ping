# CloudFlare-ping
A Ping program that accepts a hostname or an IP address as its argument, then send ICMP echo requests in a loop to the target while receiving echo reply messages.

## How to run the ping CLI app
Clone the repository. Enter the correct project folder. Then enter `sudo -E go run ping.go [target address]`

## Demo
* Pinging cloudfront.com
![ping_cloud_front]("https://github.com/Arthur-S-Huang/CloudFlare-ping/blob/master/ping_examples/ping_cf.png")
