package api

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"github.com/sukhajata/devicetwin.git/internal/core"
	"github.com/sukhajata/devicetwin.git/pkg/authhelper"
	pb "github.com/sukhajata/ppconfig"
	"github.com/urfave/negroni"
	"net/http"
)

// HTTPServer - provides an HTTP server
type HTTPServer struct {
	Ready         bool
	Live          bool
	configService core.ConfigHandler
	allowedRoles  []string
}

type desiredConfigRequest struct {
	DeviceEUI  string `json:"deviceEUI"`
	FieldName  string `json:"fieldName"`
	FieldValue string `json:"fieldValue"`
	Slot       int32  `json:"slot"`
}

type configField struct {
	Name     string `json:"name"`
	Desired  string `json:"desired"`
	Reported string `json:"reported"`
}

type updateFirmwareRequest struct {
	Firmware string `json:"firmware"`
}

func (s *HTTPServer) readinessHandler(w http.ResponseWriter, r *http.Request) {
	if !s.Ready {
		http.Error(w, "Not ready", http.StatusInternalServerError)
		return
	}
	_, err := fmt.Fprintf(w, "Ready")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

}

func (s *HTTPServer) livenessHandler(w http.ResponseWriter, r *http.Request) {
	if !s.Live {
		http.Error(w, "Not live", http.StatusInternalServerError)
		return
	}
	_, err := fmt.Fprintf(w, "Live")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (s *HTTPServer) postSetDesiredHandler(w http.ResponseWriter, r *http.Request) {
	token, err := authhelper.GetTokenFromHeader(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	decoder := json.NewDecoder(r.Body)
	var content desiredConfigRequest
	err = decoder.Decode(&content)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	req := &pb.SetDesiredRequest{
		Identifier: content.DeviceEUI,
		FieldName:  content.FieldName,
		FieldValue: content.FieldValue,
		Slot:       content.Slot,
	}
	response, err := s.configService.SetDesired(token, req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = w.Write([]byte(response.GetReply()))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (s *HTTPServer) getConfigByNameHandler(w http.ResponseWriter, r *http.Request) {
	token, err := authhelper.GetTokenFromHeader(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	vars := mux.Vars(r)
	deviceeui, ok := vars["deviceeui"]
	if !ok {
		http.Error(w, "missing parameter deviceeui", http.StatusBadRequest)
		return
	}
	name, ok := vars["name"]
	if !ok {
		http.Error(w, "missing parameter name", http.StatusBadRequest)
		return
	}

	req := &pb.GetConfigByNameRequest{
		Identifier: deviceeui,
		FieldName:  name,
	}

	response, err := s.configService.GetConfigByName(token, req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	field := &configField{
		Name:     response.GetName(),
		Desired:  response.GetDesired(),
		Reported: response.GetReported(),
	}
	b, err := json.Marshal(field)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = w.Write(b)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

}

func (s *HTTPServer) getAssignRoffsetHandler(w http.ResponseWriter, r *http.Request) {
	token, err := authhelper.GetTokenFromHeader(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	vars := mux.Vars(r)
	deviceeui, ok := vars["deviceeui"]
	if !ok {
		http.Error(w, "missing parameter deviceeui", http.StatusBadRequest)
		return
	}

	req := &pb.Identifier{
		Identifier: deviceeui,
	}

	response, err := s.configService.AssignRadioOffset(token, req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	b, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = w.Write(b)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

}

func (s *HTTPServer) postUpdateFirmwareHandler(w http.ResponseWriter, r *http.Request) {
	token, err := authhelper.GetTokenFromHeader(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	decoder := json.NewDecoder(r.Body)
	var content updateFirmwareRequest
	err = decoder.Decode(&content)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = s.configService.UpdateFirmwareAllDevices(token)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err = w.Write([]byte("Started updating"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func NewHTTPServer(configService core.ConfigHandler) *HTTPServer {
	s := &HTTPServer{
		configService: configService,
		Ready:         true,
		Live:          true,
	}

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowCredentials: true,
		AllowedHeaders:   []string{"X-Requested-With", "Content-Type", "Authorization"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
		//Debug: true,
	})

	router := mux.NewRouter()

	// default
	router.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		_, err := fmt.Fprintf(w, "Welcome to the home page!")
		if err != nil {
			fmt.Println(err)
		}
	})

	router.HandleFunc("/health/ready", s.readinessHandler).Methods("GET")
	router.HandleFunc("/health/live", s.livenessHandler).Methods("GET")
	router.HandleFunc("/set", s.postSetDesiredHandler).Methods("POST")
	router.HandleFunc("/get/{deviceeui}/{name}", s.getConfigByNameHandler).Methods("GET")
	router.HandleFunc("/roffset/{deviceeui}", s.getAssignRoffsetHandler).Methods("GET")
	router.HandleFunc("/update-firmware", s.postUpdateFirmwareHandler).Methods("POST")

	n := negroni.New()
	n.Use(negroni.NewRecovery())

	n.UseHandler(router)
	n.Use(c)

	//start on new goroutine
	go func() {
		//log.Fatal(http.ListenAndServe(":80", nil))
		n.Run(":80")
	}()

	return s
}
