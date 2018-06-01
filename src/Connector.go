package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"net/url"
	"os/exec"
	"strings"
	"time"

	"github.com/beevik/etree"
	"golang.org/x/net/proxy"
)

type Port struct {
	Protocol       string
	PortID         string
	State          string
	Reason         string
	ServiceName    string
	ProductName    string
	ProductVersion string
	ExtraInfo      string
	OSType         string
}
type Host struct {
	TorResponse string
	Address     string
	HostName    string
	Ports       []Port
}

var (
	connectortemplate *template.Template
	oHost             Host
	fName             string
)

func portscan(target string, finflag chan string) {
	nmapPath, _ := exec.LookPath("nmap")
	torproxy := "socks4://127.0.0.1:9050"
	scannerargs := []string{"--proxy", torproxy, "--dns-servers", "8.8.8.8", "-T4", "-sV", "-Pn", "-A", "--reason", "-v", target, "-oX", "anonscanres.xml"}
	//args2 := []string{"-Pn", "-sT", "-sV", "-O", target, "-oX", fName}
	scanner := exec.Command(nmapPath, scannerargs...)
	//scanner.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	//scanner.Stdout = os.Stdout
	//scanner.Stderr = os.Stderr
	//scanner.Stdin = os.Stdin
	err := scanner.Start()
	//fmt.Println(string(output))
	fmt.Println("scan started")
	if err != nil {
		fmt.Println(err)
		return
	}
	scanner.Wait()

	/*fi, err := os.Open("anonscanres2.xml")
	if err != nil {
		fmt.Println(err)
	}
	defer fi.Close()
	content, err := ioutil.ReadAll(fi)*/
	//fullpath := path.Join(path.Dir(fi.Name()), fi.Name())
	//return string(content)
	finflag <- "Scanning is over"

}

func connectTor(targetURL string) string {

	//targeturl := ""
	torproxy := "socks5://127.0.0.1:9050"
	torproxyurl, err := url.Parse(torproxy)
	if err != nil {
		fmt.Println(err)
	}

	torDialer, err := proxy.FromURL(torproxyurl, proxy.Direct)
	if err != nil {
		fmt.Println(err)
	}
	torTransport := &http.Transport{Dial: torDialer.Dial}
	client := &http.Client{
		Transport: torTransport,
		Timeout:   time.Second * 5,
	}

	req, err := http.NewRequest("GET", targetURL, nil)
	if err != nil {
		fmt.Println(err)
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/51.0.2704.103 Safari/537.36")
	//response, err := client.Get(targeturl)
	response, err := client.Do(req)
	if err != nil {
		fmt.Println("No response")
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println("error in body")
	}

	return (string(body))

}
func init() {
	connectortemplate = template.Must(template.ParseFiles("template/torconnector.html"))
	oHost = Host{}
	fName = "anonscanres2.xml"

	//oscanResult = scanResult{}
}

func parseScanResult(filename string, finflag chan string) {
	doc := etree.NewDocument()
	err := doc.ReadFromFile(filename)
	if err != nil {
		fmt.Println(err)
	}

	root := doc.SelectElement("nmaprun")
	fmt.Println(root.Tag)
	host := root.SelectElement("host")
	if host != nil {
		address := host.SelectElement("address")
		if address != nil && address.SelectAttr("addr") != nil {
			oHost.Address = address.SelectAttr("addr").Value
		}

		hostnames := host.SelectElement("hostnames")
		if len(hostnames.SelectElements("hostname")) > 0 {
			oHost.HostName = hostnames.SelectElements("hostname")[0].SelectAttr("name").Value

		}
		ports := host.SelectElement("ports")

		oHost.Ports = []Port{}
		for _, port := range ports.SelectElements("port") {
			oPort := Port{}
			if port.SelectAttr("protocol") != nil {
				oPort.Protocol = (port.SelectAttr("protocol").Value)
			}
			if port.SelectAttr("portid") != nil {
				oPort.PortID = (port.SelectAttr("portid").Value)
			}
			state := port.SelectElement("state")
			if state != nil {
				oPort.State = state.SelectAttr("state").Value
				oPort.Reason = state.SelectAttr("reason").Value
			}
			service := port.SelectElement("service")
			if service != nil {
				if service.SelectAttr("name") != nil {
					oPort.ServiceName = (service.SelectAttr("name").Value)
				}
				if service.SelectAttr("extrainfo") != nil {
					oPort.ExtraInfo = (service.SelectAttr("extrainfo").Value)
				}
				if service.SelectAttr("version") != nil {
					oPort.ProductVersion = (service.SelectAttr("version").Value)
				}
				if service.SelectAttr("product") != nil {
					oPort.ProductName = (service.SelectAttr("product").Value)
				}
			}

			oHost.Ports = append(oHost.Ports, oPort)

		}
	}

	finflag <- "Parsing over"
}

func scan(httpw http.ResponseWriter, req *http.Request) {
	if req.Method == "GET" {
		err := connectortemplate.Execute(httpw, nil)
		if err != nil {
			fmt.Println(err)
		}
	} else {
		err := req.ParseForm()
		if err != nil {
			fmt.Println(err)
		}
		target := req.Form.Get("target")
		fmt.Println(target)
		//oscanResult.Result = portscan(target)
		//fmt.Println(oscanResult.Result)
		//otorResponse.ResposeText = connectTor(urltocheck)
		//err = connectortemplate.Execute(httpw, oscanResult)
		//if err != nil {
		//fmt.Println(err)
		//}
	}
}

func index(httpw http.ResponseWriter, req *http.Request) {
	if req.Method == "GET" {
		err := connectortemplate.Execute(httpw, nil)
		if err != nil {
			fmt.Println(err)
		}
	} else {
		err := req.ParseForm()
		if err != nil {
			fmt.Println(err)
		}
		target := req.Form.Get("target")
		if strings.TrimSpace(target) != "" {
			fmt.Println(target)
			finflag := make(chan string)
			go portscan(target, finflag)
			<-finflag
			go parseScanResult(fName, finflag)
			<-finflag

			//fmt.Println(otorResponse.Result)
			//otorResponse.ResposeText = connectTor(urltocheck)
			err = connectortemplate.Execute(httpw, oHost)
			if err != nil {
				fmt.Println(err)
			}
		}

		urltocheck := req.Form.Get("url")
		if strings.TrimSpace(urltocheck) != "" {
			fmt.Println(urltocheck)
			oHost.TorResponse = connectTor(urltocheck)
			err = connectortemplate.Execute(httpw, oHost)
			if err != nil {
				fmt.Println(err)
			}
		}
	}
}

func main() {
	fmt.Println("App is ready : http://0.0.0.0:7777")
	http.HandleFunc("/", index)
	//http.HandleFunc("/scan", scan)
	http.ListenAndServe(":7777", nil)
	//connectTor("https://www.whatismyip.com")
	//connectTor("https://ifconfig.co/")
}
