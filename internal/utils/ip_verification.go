package utils

import (
	"fmt"
	"net"
	"net/http"
)

func isTrustedIP(IP string, trustedSubnet string) (bool, error) {
	parsedIP := net.ParseIP(IP)

	if parsedIP == nil {
		return false, fmt.Errorf("неверный IP-адрес: %s", IP)
	}

	_, trustedNet, err := net.ParseCIDR(trustedSubnet)

	if err != nil {
		return false, fmt.Errorf("неверная подсеть: %s", trustedSubnet)
	}

	return trustedNet.Contains(parsedIP), nil
}

func GetLocalIP() (string, error) {
	addrs, err := net.InterfaceAddrs()

	if err != nil {
		return "", err
	}

	for _, addr := range addrs {
		if ipNet, ok := addr.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
			if ipNet.IP.To4() != nil {
				return ipNet.IP.String(), nil
			}
		}
	}

	return "", fmt.Errorf("не удалось получить локальный IP-адрес")
}

func VerifyIPMiddleware(trustedSubnet string) func(next http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			if trustedSubnet == "" {
				next.ServeHTTP(w, r)
				return
			}

			IP := r.Header.Get("X-Real-IP")

			if IP == "" {
				http.Error(w, "X-Real-IP заголовок не найден", http.StatusForbidden)
				return
			}

			isTrusted, err := isTrustedIP(IP, trustedSubnet)

			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			if !isTrusted {
				http.Error(w, "IP-адрес не является доверенным", http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		}
	}
}
