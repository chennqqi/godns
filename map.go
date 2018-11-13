package main

import (
	"golang.org/x/net/publicsuffix"
	"strings"
)

type IPList struct {
	addrs []string
	index int
}

func (t *IPList) ToList() []string {
	l := len(t.addrs)
	if l <= 1 {
		return t.addrs
	}

	index := t.index % l
	ret := make([]string, len(t.addrs))
	ret = append(ret, t.addrs[index:]...)
	ret = append(ret, t.addrs[:index]...)
	t.index++
	return ret
}

func (t *IPList) Add(ip string) {
	t.addrs = append(t.addrs, ip)
}

type AddrMap map[string]*IPList

func (t AddrMap) Get(domain string) ([]string, bool) {
	domain = strings.ToLower(domain)

	iplist, ok := t[domain]
	if ok {
		return iplist.ToList(), true
	}

	sld, err := publicsuffix.EffectiveTLDPlusOne(domain)
	if err != nil {
		return nil, false
	}

	for host, iplist := range t {
		if strings.HasPrefix(host, "*.") {
			old, err := publicsuffix.EffectiveTLDPlusOne(host)
			if err != nil {
				continue
			}
			if sld == old {
				return iplist.ToList(), true
			}
		}
	}
	return []string{}, false
}

func (t AddrMap) Add(domain, ip string) {
	domain = strings.ToLower(domain)
	iplist, ok := t[domain]
	if ok {
		iplist.Add(ip)
	} else {
		t[domain] = &IPList{
			[]string{ip}, 0,
		}
	}
}
