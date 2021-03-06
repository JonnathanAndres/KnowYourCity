// The account handles all the user related tasks for the administrator
package main

import (
	"log"
	"net/http"
	"strings"
	"time"
	"net/rpc"
	"strconv"

	"github.com/gorilla/mux"

	"github.com/Code4SierraLeone/KnowYourCity/base"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// Register handles the registration of a new admin user account
func Register(rw http.ResponseWriter, r *http.Request) {
	// Get request variables
	dbSession := GetVar(r, "db").(*mgo.Session)
	client := GetVar(r, "rpc").(*rpc.Client)
	defer RemoveVars(r)
	email := r.FormValue("email")
	uname := r.FormValue("userName")
	password := r.FormValue("password")

	if email == "" || password == "" {
		http.Error(rw, "Invalid parameters", http.StatusBadRequest)
		return
	}

	// generate  uuid
	// filename := "placeholder"
	var uid string
	if err := client.Call("UUIDGenerator.GenerateAdmin", 0, &uid); err != nil {
		ReportFatal(rw, r, err)
		return
	}

	admin := new(base.Admin)
	admin.Setup(dbSession)
	admin.Id = bson.NewObjectId()
	admin.Email = email
	admin.UID = uid
	admin.UserName = uname
	admin.Designation = "0"
	if err := admin.Create(); err != nil {
		if mgo.IsDup(err) {
			base.RespondError(rw, r, "Email already exists")
			return
		}
		base.RespondError(rw, r, "Server Error")
		return
	}
	LogInfo("Registered admin details, creating auth details")

	auth := new(base.Authentication)
	auth.Setup(dbSession)
	auth.Id = bson.NewObjectId()
	auth.UID = uid
	auth.Email = email
	auth.Password = base.GetSha(password)
	auth.Code = 0
	if err := auth.Create(); err != nil {
		LogError(err)
		ReportFatal(rw, r, err)
		return
	}
	LogInfo("Created auth details")
	// feedchan <- time.Now().String() + "::" + admin.UserName + " registered for a Admin Account"
	base.RespondSuccess(rw, r, "success")
}

// Login handles admin users login and creates new sessions
func Login(rw http.ResponseWriter, r *http.Request) {
	dbSession := GetVar(r, "db").(*mgo.Session)
	client := GetVar(r, "rpc").(*rpc.Client)
	// uid := GetVar(r, "uid").(string)
	defer RemoveVars(r)

	email := r.FormValue("email")
	password := r.FormValue("password")
	if email == "" || password == "" {
		http.Error(rw, "Incomplete parameters", http.StatusBadRequest)
		return
	}
	auth := new(base.Authentication)
	auth.Setup(dbSession)
	if auth.AuthenticateAdmin(email, base.GetSha(password), client) != nil {
		base.RespondError(rw, r, "Invalid Email/Password Combination")
		return
	}
	feedchan <- Payload{"e", time.Now().String() + "::" + auth.Email + " logged in to the Admin panel"}
	base.RespondSuccess(rw, r, auth.UID)
}


// Logout handles admin user logout and stopping of their current sessions
func Logout(rw http.ResponseWriter, r *http.Request) {
	uid := GetVar(r, "uid").(string)
	client := GetVar(r, "rpc").(*rpc.Client)

	defer RemoveVars(r)

	var deleted bool
	if err := client.Call("Sessions.Delete", uid, &deleted); err != nil {
		ReportFatal(rw, r, err)
		return
	}
	DestroySession(uid)
	LogInfo("Deleted Active session")
	feedchan <- Payload{uid, time.Now().String() + "::" + uid + " just logged out."}
	base.RespondSuccess(rw, r, "success")
}

// ForgotPass handles Forgotten password reset
func ForgotPass(rw http.ResponseWriter, r *http.Request) {
	dbSession := base.GetVar(r, "db").(*mgo.Session)
	defer base.RemoveVars(r)

	email := r.FormValue("email")

	userModel := new(base.User)
	userModel.Setup(dbSession)
	if err := userModel.Get("", email); err != nil {
		base.RespondError(rw, r, "No such email")
		return
	}

	userModel.ForgotPassKey = base.GetBase64(userModel.UID.String() + ":" + time.Now().String())
	userModel.Setup(dbSession)
	if err := userModel.Update(); err != nil {
		log.Println(err)
		base.RespondError(rw, r, "Something's broken")
		return
	}
	// TODO:send password reset link
	base.RespondSuccess(rw, r, "success")
}

// Rendering of the Reset Password page
func ResetPassPage(rw http.ResponseWriter, r *http.Request) {
	dbSession := base.GetVar(r, "db").(*mgo.Session)
	defer base.RemoveVars(r)

	vars := mux.Vars(r)
	key := vars["key"]

	userModel := new(base.User)
	userModel.Setup(dbSession)
	if err := userModel.ForgotPassLinkValid(key, dbSession); err != nil {
		// redirect to 404 page
		return
	}
	decodedKey, err := base.DecodeBase64(key)
	if err != nil {
		// redirect to 404 page
		return
	}
	uid := strings.SplitAfter(decodedKey, ":")[0]

	vals := map[string]string{
		"uid": uid,
	}
	if encoded, err := sc.Encode("knowyourcity", vals); err == nil {
		cookie := &http.Cookie{
			Name:  "knowyourcity",
			Value: encoded,
			Path:  "/",
		}
		http.SetCookie(rw, cookie)
	} else {
		// render 404 page
		return
	}
	// set secure cookie
	// render page
}

// ResetPass handles the reset password request
func ResetPass(rw http.ResponseWriter, r *http.Request) {
	dbSession := base.GetVar(r, "db").(*mgo.Session)
	defer base.RemoveVars(r)

	// fetch secure cookie
	value := make(map[string]string)
	if cookie, err := r.Cookie("knowyourcity"); err == nil {
		if err := sc.Decode("knowyourcity", cookie.Value, &value); err != nil {
			log.Println(err)
			base.RespondError(rw, r, "Server Error")
			return
		}
	}

	userModel := new(base.User)
	userModel.Setup(dbSession)
	if err := userModel.Get(value["uid"], ""); err != nil {
		base.RespondError(rw, r, "Invalid UID")
		return
	}

	pass := r.FormValue("pass")
	pass2 := r.FormValue("pass2")

	if pass != pass2 {
		base.RespondError(rw, r, "Passwords dont match")
		return
	}

	userModel.Password = base.GetSha(pass)
	if err := userModel.Update(); err != nil {
		base.RespondError(rw, r, "Something's broken")
		return
	}
	base.RespondSuccess(rw, r, "Success")
}
