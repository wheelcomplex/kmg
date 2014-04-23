package kmgTls

import (
	"crypto/tls"
	"fmt"
	"github.com/bronze1man/kmg/crypto/kmgTls"
	"net"
	"net/http"
)

func SelfCertHttpListenAndServe(addr string, handler http.Handler) error {
	l, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	tlsConfig, err := kmgTls.CreateTlsConfig()
	if err != nil {
		return fmt.Errorf("fail at kmgTls.CreateTlsConfig,error:%s", err.Error())
	}
	l = tls.NewListener(l, tlsConfig)
	return http.Serve(l, handler)
}
