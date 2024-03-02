package handlers

import (
	"encoding/json"
	"net"
	"net/http"

	"github.com/kupriyanovkk/shortener/internal/config"
)

// isIPInTrustedSubnet checks if IP is in trusted subnet
func isIPInTrustedSubnet(ip, subnet string) bool {
	_, ipNet, err := net.ParseCIDR(subnet)
	if err != nil {
		return false
	}
	return ipNet.Contains(net.ParseIP(ip))
}

// GetInternalStats process request for getting internal statistics
func GetInternalStats(w http.ResponseWriter, r *http.Request, app *config.App) {
	if app.Flags.TrustedSubnet == "" {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	xRealIP := r.Header.Get("X-Real-Ip")

	if !isIPInTrustedSubnet(xRealIP, app.Flags.TrustedSubnet) {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	stats, err := app.Store.GetInternalStats(r.Context())

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	enc := json.NewEncoder(w)
	if err := enc.Encode(stats); err != nil {
		return
	}
}
