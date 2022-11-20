package main

import (
	"fmt"
	pb "golang.conradwood.net/apis/certmanager"
	"golang.conradwood.net/go-easyops/prometheus"
	"golang.conradwood.net/go-easyops/tokens"
	"golang.conradwood.net/go-easyops/utils"
	"sort"
	"time"
)

const (
	minDays = 30
)

var (
	expiryGauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "certmanager_cert_age",
			Help: "V=1 UNIT=s DESC=age of certificate",
		},
		[]string{"host"},
	)
)

func init() {
	prometheus.MustRegister(expiryGauge)
}

func refresher() {
	for {
		ctx := tokens.ContextWithToken()
		certs, err := certStore.All(ctx)
		if err != nil {
			fmt.Printf("[refresher] failed to get certs from db:%s\n", err)
			time.Sleep(time.Duration(5) * time.Minute)
			continue
		}
		sort.Slice(certs, func(i, j int) bool {
			return certs[i].LastAttempt < certs[j].LastAttempt
		})
		cutoff := time.Now().Add(time.Duration(minDays*24) * time.Hour)
		dorand := time.Now().Add(time.Duration(minDays*12) * time.Hour)
		for _, c := range certs {
			if *debug {
				fmt.Printf("Certificate %s: Expiry: %s, LastAttempt: %s\n",
					c.Host,
					utils.TimestampString(c.Expiry),
					utils.TimestampString(c.LastAttempt),
				)
			}
			t := time.Unix(int64(c.Expiry), 0)
			l := prometheus.Labels{"host": c.Host}
			expiryInSecs := int64(c.Expiry) - time.Now().Unix()
			expiryGauge.With(l).Set(float64(expiryInSecs))
			// t == expiry date
			if t.After(cutoff) {
				// expiry is AFTER now+24 days
				continue
			}
			if t.After(dorand) {
				// expiry is AFTER now+12 days
				// space them out a bit, so to avoid letsencrypt limits
				if utils.RandomInt(12) < 11 {
					continue
				}
			}
			r := &requestCertificate{
				req: &pb.PublicCertRequest{
					Hostname:   c.Host,
					VerifyType: pb.VerifyType_HTTP,
				},
				serviceid: c.CreatorService,
				userid:    c.CreatorUser,
			}
			reqChannel <- r

		}
		time.Sleep(time.Duration(4) * time.Hour)
	}
}
