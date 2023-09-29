// interactive web user interface
package web

import (
	"github.com/gorilla/websocket"
	"net/http"
	"ngrok/client/assets"
	"ngrok/client/mvc"
	"ngrok/log"
	"ngrok/proto"
	"ngrok/util"
	"path"
	"encoding/json"
	"crypto/rand"
	"fmt"
	"net/url"
)
type TunnelListResource struct {
	Tunnels   []tunnel `json:"tunnels"`
	Uri string			`json:"uri"`
}

type Konfig struct {
	Addr string			`json:"addr"`
	Inspect bool		`json:"inspect"`
}


type tunnel struct {
	Name string			`json:"name"`
	Id string			`json:"ID"`
	Uri string			`json:"uri"`
	PublicUrl string	`json:"public_url"`
	Protocol string		`json:"proto"` // http/tcp/https/http+https	
	Config Konfig		`json:"config"`
}
/*

<Name>command_line</Name>
<ID>6a59522eae8051c252b212daf850d7cd</ID>
<URI>/api/tunnels/command_line</URI>
<PublicURL>https://bc9182948fa6.ngrok.app</PublicURL>
<Proto>https</Proto>

tunnel1 := tunnel{"command_line", "6a59522eae8051c252b212daf850d7cd", "/api/tunnels/command_line", "https://bc9182948fa6.ngrok.app", "https"}
tunnels := []tunnel{tunnel1}
fmt.Printf("tunnels is %v\\n", tunnels)

    b, err := json.Marshal(tunnels)
    if err != nil {
        fmt.Println(err)
        return
    }
    fmt.Println(string(b))

	*/
type WebView struct {
	log.Logger

	ctl mvc.Controller

	// messages sent over this broadcast are sent to all websocket connections
	wsMessages *util.Broadcast
}
type Guid [16]byte

// String returns a standard hexadecimal string version of the Guid.
// Lowercase characters are used.
func GuidToString(g *Guid) string {
	return fmt.Sprintf("%x-%x-%x-%x-%x",
		g[0:4], g[4:6], g[6:8], g[8:10], g[10:16])
}

func GuidToStringNoHyphen(g *Guid) string {
	return fmt.Sprintf("%x", g)
}

func NewGuid() *Guid {
	g := new(Guid)
	if _, err := rand.Read(g[:]); err != nil {
		panic(err)
	}
	g[6] = (g[6] & 0x0f) | 0x40 // version = 4
	g[8] = (g[8] & 0x3f) | 0x80 // variant = RFC 4122
	return g
}

func NewWebView(ctl mvc.Controller, addr string) *WebView {
	wv := &WebView{
		Logger:     log.NewPrefixLogger("view", "web"),
		wsMessages: util.NewBroadcast(),
		ctl:        ctl,
	}

	// for now, always redirect to the http view
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/http/in", 302)
	})
	// no support for api/tunnels/tunnelName so maybe should change Uri:"/api/tunnels/" + t.Name to blank?
	http.HandleFunc("/api/tunnels", func(w http.ResponseWriter, r *http.Request) {

		var fs []tunnel
		
		for _, t := range ctl.State().GetTunnels() {
			var addr string
			if t.BindTls {
				addr = "https://"
			} else {
				addr = "http://"
			}
			c := Konfig{Addr:addr + t.LocalAddr, Inspect:true}
			g := NewGuid()
			u, _ := url.Parse(t.PublicUrl)
			tunnel1 := tunnel{Name:t.Name, 
				Id: GuidToString(g), 
				Uri:"/api/tunnels/" + t.Name, 
				PublicUrl:t.PublicUrl, 
				Protocol:u.Scheme, 
				Config:c}
			fs = append(fs, tunnel1)
			wv.Info("/api/tunnels proto %v", t.Protocol.GetName())
		}
		fs1 := TunnelListResource{Tunnels:fs, Uri:"/api/tunnels"}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(fs1)
	})

	// handle web socket connections
	http.HandleFunc("/_ws", func(w http.ResponseWriter, r *http.Request) {
		conn, err := websocket.Upgrade(w, r, nil, 1024, 1024)

		if err != nil {
			http.Error(w, "Failed websocket upgrade", 400)
			wv.Warn("Failed websocket upgrade: %v", err)
			return
		}

		msgs := wv.wsMessages.Reg()
		defer wv.wsMessages.UnReg(msgs)
		for m := range msgs {
			err := conn.WriteMessage(websocket.TextMessage, m.([]byte))
			if err != nil {
				// connection is closed
				break
			}
		}
	})

	// serve static assets
	http.HandleFunc("/static/", func(w http.ResponseWriter, r *http.Request) {
		buf, err := assets.Asset(path.Join("assets", "client", r.URL.Path[1:]))
		if err != nil {
			wv.Warn("Error serving static file: %s", err.Error())
			http.NotFound(w, r)
			return
		}
		w.Write(buf)
	})

	wv.Info("Serving web interface on %s", addr)
	wv.ctl.Go(func() { http.ListenAndServe(addr, nil) })
	return wv
}

func (wv *WebView) NewHttpView(proto *proto.Http) *WebHttpView {
	return newWebHttpView(wv.ctl, wv, proto)
}

func (wv *WebView) Shutdown() {
}
