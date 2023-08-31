// client create: CertManagerClient
/*
  Created by /home/cnw/devel/go/yatools/src/golang.yacloud.eu/yatools/protoc-gen-cnw/protoc-gen-cnw.go
*/

/* geninfo:
   filename  : protos/golang.conradwood.net/apis/certmanager/certmanager.proto
   gopackage : golang.conradwood.net/apis/certmanager
   importname: ai_0
   clientfunc: GetCertManager
   serverfunc: NewCertManager
   lookupfunc: CertManagerLookupID
   varname   : client_CertManagerClient_0
   clientname: CertManagerClient
   servername: CertManagerServer
   gsvcname  : certmanager.CertManager
   lockname  : lock_CertManagerClient_0
   activename: active_CertManagerClient_0
*/

package certmanager

import (
   "sync"
   "golang.conradwood.net/go-easyops/client"
)
var (
  lock_CertManagerClient_0 sync.Mutex
  client_CertManagerClient_0 CertManagerClient
)

func GetCertManagerClient() CertManagerClient { 
    if client_CertManagerClient_0 != nil {
        return client_CertManagerClient_0
    }

    lock_CertManagerClient_0.Lock() 
    if client_CertManagerClient_0 != nil {
       lock_CertManagerClient_0.Unlock()
       return client_CertManagerClient_0
    }

    client_CertManagerClient_0 = NewCertManagerClient(client.Connect(CertManagerLookupID()))
    lock_CertManagerClient_0.Unlock()
    return client_CertManagerClient_0
}

func CertManagerLookupID() string { return "certmanager.CertManager" } // returns the ID suitable for lookup in the registry. treat as opaque, subject to change.

func init() {
   client.RegisterDependency("certmanager.CertManager")
}
