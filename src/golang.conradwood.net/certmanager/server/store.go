package main

import (
	"context"
	"fmt"
	cm "golang.conradwood.net/apis/certmanager"
	"golang.conradwood.net/go-easyops/errors"
	"sort"
	"sync"
	"time"
)

var (
	templock sync.Mutex
)

func store(domain string, token string, keyauth string) error {
	templock.Lock()
	defer templock.Unlock()

	ctx := context.Background()
	clearStore(ctx, domain)
	s := &cm.StoreAuth{Domain: domain, Token: token, KeyAuth: keyauth, Created: uint32(time.Now().Unix())}
	_, err := authStore.Save(ctx, s)
	if err != nil {
		fmt.Printf("Failed to store auth: %s\n", err)
		return err
	}
	return nil
}

func getFromStore(ctx context.Context, domain string) (*cm.StoreAuth, error) {
	if domain == "" {
		return nil, errors.InvalidArgs(ctx, "getstore: missing domain", "getstore: missing domain")
	}
	templock.Lock()
	defer templock.Unlock()
	obs, err := authStore.ByDomain(ctx, domain)
	if err != nil {
		fmt.Printf("Failed to get auth from store: %s\n", err)
		return nil, err
	}
	if len(obs) == 0 {
		return nil, errors.NotFound(ctx, "auth not found", "auth for domain \"%s\" not found", domain)
	}
	sort.Slice(obs, func(i, j int) bool {
		return obs[i].Created < obs[i].Created
	})
	return obs[0], nil
}

func getFromStoreByChallenge(ctx context.Context, challenge string) (*cm.StoreAuth, error) {
	if challenge == "" {
		return nil, errors.InvalidArgs(ctx, "getstore: missing challenge", "getstore: missing challenge")
	}
	templock.Lock()
	defer templock.Unlock()
	obs, err := authStore.ByToken(ctx, challenge)
	if err != nil {
		fmt.Printf("Failed to get auth from store by challenge: %s\n", err)
		return nil, err
	}
	if len(obs) == 0 {
		return nil, errors.NotFound(ctx, "auth not found", "auth for challenge \"%s\" not found", challenge)
	}
	sort.Slice(obs, func(i, j int) bool {
		return obs[i].Created < obs[i].Created
	})
	return obs[0], nil
}

func clearStore(ctx context.Context, hostname string) error {
	s := "delete from " + authStore.Tablename() + " where domain = $1"
	_, err := psql.ExecContext(ctx, "clearstore", s, hostname)
	return err
}
