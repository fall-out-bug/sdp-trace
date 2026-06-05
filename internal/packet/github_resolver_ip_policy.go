package packet

import "net"

// IP resolver policy mirrors the host policy for literal addresses. Loopback,
// private, link-local, multicast link-local, and unspecified addresses are not
// portable evidence references. The rules are table-driven so adding a future
// blocked class does not change the command behavior around relative resolver
// strings or public HTTPS URLs.

func unsafeResolverIP(ip net.IP) bool {
	if ip == nil {
		return false
	}
	return anyResolverIPRuleMatches(ip, []func(net.IP) bool{
		loopbackResolverIP,
		privateResolverIP,
		linkLocalResolverIP,
		unspecifiedResolverIP,
	})
}

func anyResolverIPRuleMatches(ip net.IP, rules []func(net.IP) bool) bool {
	for _, rule := range rules {
		if rule(ip) {
			return true
		}
	}
	return false
}

func loopbackResolverIP(ip net.IP) bool {
	return ip.IsLoopback()
}

func privateResolverIP(ip net.IP) bool {
	return ip.IsPrivate()
}

func linkLocalResolverIP(ip net.IP) bool {
	return ip.IsLinkLocalUnicast() || ip.IsLinkLocalMulticast()
}

func unspecifiedResolverIP(ip net.IP) bool {
	return ip.IsUnspecified()
}
