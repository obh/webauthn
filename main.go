package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-webauthn/webauthn/protocol"
	"github.com/go-webauthn/webauthn/webauthn"
)

var (
	w     *webauthn.WebAuthn
	err   error
	store *Datastore
)

func main() {
	/// Webauthn configuration
	webAuthnConfig := &webauthn.Config{
		RPDisplayName: "Cashfree",
		RPID:          "localhost",
		RPOrigins:     []string{"http://localhost"},
	}
	if w, err = webauthn.New(webAuthnConfig); err != nil {
		fmt.Println("cannot initialize webauthn...", err)
	}

	//datastore configuration
	store = InitData()

	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/", fs)
	http.HandleFunc("/home", goHome)
	http.HandleFunc("/options", registerOptions)
	http.HandleFunc("/register", register)

	fmt.Println("starting server...")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Failed in starting server...", err)
	}
}

func goHome(resp http.ResponseWriter, req *http.Request) {
	s := "hello world"
	setCORSHeader(resp)
	resp.Write([]byte(s))
}

func registerOptions(resp http.ResponseWriter, req *http.Request) {
	user, err := store.GetUser(UserId("101"))
	if err != nil {
		resp.Write([]byte("failed to find user"))
	}
	options, session, err := w.BeginRegistration(user)
	store.StoreSession(UserId("101"), session)
	if err != nil {
		resp.Write([]byte("failed in registration"))
	}
	setCORSHeader(resp)
	resp.WriteHeader(http.StatusOK)

	fmt.Println(session)
	json.NewEncoder(resp).Encode(options)
}

func register(resp http.ResponseWriter, req *http.Request) {
	body := req.Body
	fmt.Println("Got body: ", body)
	response, err := protocol.ParseCredentialCreationResponseBody(body)
	if err != nil {
		fmt.Println("failed in parsing data:", err.Error())
		setCORSHeader(resp)
		resp.Write([]byte("registeration failed"))
		return
	}
	fmt.Println(response.Raw.AttestationResponse.ClientDataJSON)
	user, err := store.GetUser(UserId("101")) // Get the user
	if err != nil {
		fmt.Println("cannot find user:", err)
		setCORSHeader(resp)
		resp.Write([]byte("registeration failed"))
		return
	}
	session, err := store.GetSession(UserId("101"))
	if err != nil {
		fmt.Println("cannot find session:", err)
		setCORSHeader(resp)
		resp.Write([]byte("registeration failed"))
		return
	}

	credential, err := w.CreateCredential(user, *session, response)
	if err != nil {
		fmt.Println("cannot create credential:", err)
		setCORSHeader(resp)
		resp.Write([]byte("registeration failed"))
		return
	}
	fmt.Println("credential created successfully: ", credential)
	// If creation was successful, store the credential object
	setCORSHeader(resp)
	resp.WriteHeader(http.StatusOK)

	resp.Write([]byte("registeration successful"))

	// Pseudocode to add the user credential.
	// user.AddCredential(credential)
	// datastore.SaveUser(user)
}

func setCORSHeader(resp http.ResponseWriter) {
	resp.Header().Set("Content-Type", "application/json")
	// resp.Header().Set("Access-Control-Allow-Headers", "CSRF-Token, X-Requested-By, Authorization, Content-Type")
	// resp.Header().Set("Access-Control-Allow-Origin", "*")
}
