## anonym - An easy to use framework for dark web threat intelligence and scanning.

### Getting Started

##### git clone https://github.com/diljith369/anonym.git
##### go get github.com/gocolly/colly
##### go get github.com/gocolly/colly/proxy
##### go get github.com/beevik/etree
##### go get github.com/gorilla/mux

##### Make sure that TOR is running in your machine

##### Set TOR as Windows Service
###### Download TOR browser bundle , move to the folder \Tor Browser\Browser\TorBrowser\Tor using command prompt as admin
###### Supply the following command 
###### tor.exe -service install 
###### Above command will install TOR as windows service (You may need to change the persmission of the service to LocalSystem to make it run)

####  Debian /Ubuntu 
##### apt-get install tor
##### servie tor start

### Navigate inside /src/bin/linux/
### change executable permissions to anonym
##### chmod a+x anonym
#### Run the anonym binary 
##### ./anonym
#### On your favourite browser go to http://127.0.0.1:7777

### Prerequisites
#### Go 
#### nmap and TOR

### Built With
#### Go 

### Original Author
#### * **Diljith S** - *Initial work* - (https://github.com/diljith369)
