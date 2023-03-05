package main

import (
	"fmt"
	"github.com/miekg/dns"
)

func main() {
	c := new(dns.Client)
	m := new(dns.Msg)
	m.SetQuestion(dns.Fqdn(""), dns.TypeCNAME) // 设置查询的域名和类型
	m.RecursionDesired = true                  // 设置递归查询
	r, _, err := c.Exchange(m, "8.8.8.8:53")   // 设置DNS服务器地址
	if r == nil {
		fmt.Printf("*** error: %s\n", err.Error())
		return
	}
	if r.Rcode != dns.RcodeSuccess {
		fmt.Printf(" *** invalid answer name %s after MX query for %s\n", "www.baidu.com", "www.baidu.com")
		return
	}
	fmt.Println(r.Rcode)
	fmt.Println(dns.RcodeSuccess)
	for _, a := range r.Answer { // 遍历返回的答案部分
		cname := a.(*dns.CNAME)
		fmt.Println(cname.Target)
		fmt.Printf("%v\n", a) // 打印每条记录
	}
}
