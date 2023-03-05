package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	//cert --domain default.xxxx.cn --default
	//获取命令参数
	//fmt.Println(os.Args)
	//for _, arg := range os.Args {
	//	switch arg {
	//	case "cert":
	//		issue()
	//	}
	//}
	//issue()
	//var domains []string
	//for _, arg := range os.Args {
	//	result, ok := strings.CutPrefix(arg, "--domain=")
	//	if ok {
	//		domains = append(domains, result)
	//	}
	//}
	//fmt.Println(domains)

	//domain := flag.String("domain", "", "")
	var domain string
	var defaultMode bool
	var defaultString string
	if os.Args[1] == "cert" {

	}
	certCmd := flag.NewFlagSet("cert", flag.ExitOnError)

	certCmd.StringVar(&domain, "domain", "", "--domain default.xxxxx.cn")
	certCmd.BoolVar(&defaultMode, "default", false, "")
	certCmd.StringVar(&defaultString, "default", "", "")
	err := certCmd.Parse(os.Args[2:])
	if err != nil {
		fmt.Println(err)
		return
	}

	//flag.Parse()
	fmt.Println(domain)
	fmt.Println(defaultMode)
	fmt.Println(defaultString)

}
