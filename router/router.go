package router

import (
	"encoding/json"
	"github.com/r2binx/heimboard-wakeup-go/config"
	"github.com/r2binx/heimboard-wakeup-go/schedule"
	"github.com/r2binx/heimboard-wakeup-go/util"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/auth0/go-jwt-middleware/v2"
	"github.com/auth0/go-jwt-middleware/v2/validator"
	"github.com/gorilla/mux"

	"github.com/r2binx/heimboard-wakeup-go/middleware"
)

var w *util.Wakeup
var conf config.Config

func New(config config.Config, wakeup *util.Wakeup) *mux.Router {
	conf = config
	w = wakeup
	router := mux.NewRouter()
	router.Handle("/ping", StatusHandler).Methods("GET")
	router.Handle("/boot-schedule", middleware.EnsureValidToken(conf.Auth0Url, conf.Auth0Audience)(ScheduleHandler)).Methods("GET")
	router.Handle("/boot-schedule", middleware.EnsureValidToken(conf.Auth0Url, conf.Auth0Audience)(SetScheduleHandler)).Methods("POST")
	router.Handle("/wakeup", middleware.EnsureValidToken(conf.Auth0Url, conf.Auth0Audience)(WakeupHandler)).Methods("GET")

	return router
}

var StatusHandler = http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
	res.WriteHeader(http.StatusOK)
	res.Write([]byte(time.Now().String()))
})

var ScheduleHandler = http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
	permitted := checkPermission(res, req, "guest")
	if permitted {
		payload, _ := json.Marshal(w.GetSchedule())
		res.WriteHeader(http.StatusOK)
		res.Write([]byte(payload))
	}
})

var SetScheduleHandler = http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
	permitted := checkPermission(res, req, "admin")
	if permitted {
		body, err := ioutil.ReadAll(req.Body)
		if err != nil {
			returnFailure(res, http.StatusBadRequest, err)
			return
		}
		var s schedule.Schedule
		err = json.Unmarshal(body, &s)
		if err != nil {
			returnFailure(res, http.StatusBadRequest, err)
			return
		}
		log.Println("writing schedule", s)
		err = w.SetSchedule(s)
		if err != nil {
			log.Println("Failed writing schedule:", err)
			returnFailure(res, http.StatusInternalServerError, err)
			return
		} else {
			returnSuccess(res)
		}
	}
})

var WakeupHandler = http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
	permitted := checkPermission(res, req, "guest")
	if permitted {
		awake, err := w.Wake(conf.WolMac)
		if !awake {
			returnFailure(res, http.StatusInternalServerError, err)
			return
		} else {
			returnSuccess(res)
		}
	}
})

func checkPermission(res http.ResponseWriter, req *http.Request, permission string) bool {
	token := req.Context().Value(jwtmiddleware.ContextKey{}).(*validator.ValidatedClaims)
	claims := token.CustomClaims.(*middleware.CustomClaims)
	if !claims.HasPermission(permission) {
		res.Header().Set("Content-Type", "application/json")
		res.WriteHeader(http.StatusForbidden)
		res.Write([]byte(`{"message":"Insufficient permission."}`))
		return false
	}
	return true
}

func returnSuccess(res http.ResponseWriter) {
	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)
	res.Write([]byte(`{"success": true}`))
}

func returnFailure(res http.ResponseWriter, status int, err error) {
	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(status)
	res.Write([]byte(`{"success": false, "message": "` + err.Error() + `"}`))
}
