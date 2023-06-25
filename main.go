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

type Response struct {
	Message string `json:"message"`
	Status  int    `json:"status"`
}

func main() {
	/// Webauthn configuration
	requireResidentKey := true
	webAuthnConfig := &webauthn.Config{
		RPDisplayName: "Cashfree",
		RPID:          "localhost",
		RPOrigins:     []string{"http://localhost:8080"},
		AuthenticatorSelection: protocol.AuthenticatorSelection{
			AuthenticatorAttachment: protocol.Platform,
			RequireResidentKey:      &requireResidentKey,
			UserVerification:        protocol.VerificationRequired,
		},
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
	http.HandleFunc("/beginLogin", beginLogin)
	http.HandleFunc("/finishLogin", finishLogin)

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
		resp.Write(getJSON(-1, "registeration failed"))
		return
	}
	fmt.Println(response.Raw.AttestationResponse.ClientDataJSON)
	user, err := store.GetUser(UserId("101")) // Get the user
	if err != nil {
		fmt.Println("cannot find user:", err)
		setCORSHeader(resp)
		resp.Write(getJSON(-1, "registeration failed"))
		return
	}
	session, err := store.GetSession(UserId("101"))
	if err != nil {
		fmt.Println("cannot find session:", err)
		setCORSHeader(resp)
		resp.Write(getJSON(-1, "registeration failed"))
		return
	}
	credential, err := w.CreateCredential(user, *session, response)
	if err != nil {
		fmt.Println("cannot create credential:", err)
		setCORSHeader(resp)
		resp.Write(getJSON(-1, "registeration failed"))
		return
	}
	user.AddCredential(credential)
	store.SaveUser(*user)
	fmt.Println("credential created successfully: ", credential)
	// If creation was successful, store the credential object
	setCORSHeader(resp)
	resp.WriteHeader(http.StatusOK)

	resp.Write(getJSON(100, "registeration successful"))

}

func beginLogin(resp http.ResponseWriter, req *http.Request) {
	user, err := store.GetUser(UserId("101"))
	fmt.Println("Found user: ", user)
	fmt.Println("With Credentials: ", user.Credentials)
	if err != nil {
		resp.Write(getJSON(-1, "failed to find user"))
	}
	options, session, err := w.BeginLogin(user)
	fmt.Println("session ---> ", session)
	if err != nil {
		fmt.Println("Found error: ", err)
		resp.Write(getJSON(-1, "failed to initiate login"))
		return
	}
	//store.StoreSession(UserId(user.Id), session)
	setCORSHeader(resp)
	resp.WriteHeader(http.StatusOK)

	fmt.Println(session)
	json.NewEncoder(resp).Encode(options)
}

func finishLogin(resp http.ResponseWriter, req *http.Request) {
}

func setCORSHeader(resp http.ResponseWriter) {
	resp.Header().Set("Content-Type", "application/json")
	resp.Header().Set("Access-Control-Allow-Headers", "CSRF-Token, X-Requested-By, Authorization, Content-Type")
	resp.Header().Set("Access-Control-Allow-Origin", "*")
}

func getJSON(status int, msg string) []byte {
	r := &Response{
		Message: msg,
		Status:  status,
	}
	b, _ := json.Marshal(r)
	return b
}
