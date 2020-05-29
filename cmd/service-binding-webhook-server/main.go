package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"

	"k8s.io/api/admission/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
)

const (
	tlsDir          = `/run/secrets/tls`
	tlsCertFile     = `tls.crt`
	tlsKeyFile      = `tls.key`
	jsonContentType = `application/json`
)

var (
	serviceBindingResource = metav1.GroupVersionResource{Group: "apps.openshift.io", Version: "v1alpha1", Resource: "servicebindingrequests"}
	universalDeserializer  = serializer.NewCodecFactory(runtime.NewScheme()).UniversalDeserializer()
)

func validateServiceBindingRequest() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		log.Println("Handling ServiceBinding request")

		var admissionReviewReq v1beta1.AdmissionReview

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Printf("Could not read body: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if _, _, err := universalDeserializer.Decode(body, nil, &admissionReviewReq); err != nil {
			log.Printf("Could not deserialize: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			return
			//return nil, fmt.Errorf("could not deserialize request: %v", err)
		} else if admissionReviewReq.Request == nil {
			log.Printf("Could not find request: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		log.Printf("ServiceBinding is being created by user %s", admissionReviewReq.Request.UserInfo.Username)

		/*
			// ADD CODE.

			* Start SubjectAccessReview ( SAR ) checks

			* SAR checks for resources the user can view.
			* SAR checks for resources the user can edit.

			// Set ALLOWED:TRUE if SAR checks pass.
			// Else, ALLOWED:FALSE if SAR checks fail.

		*/

		admissionReviewResponse := v1beta1.AdmissionReview{
			Response: &v1beta1.AdmissionResponse{
				UID:     admissionReviewReq.Request.UID,
				Allowed: false, // REJECTING EVERYTHING NOW!
				Result: &metav1.Status{
					Message: fmt.Sprintf("ServiceBinding is being created by user %s", admissionReviewReq.Request.UserInfo.Username),
				},
			},
		}
		// Return the AdmissionReview with a response as JSON.
		bytes, err := json.Marshal(&admissionReviewResponse)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			log.Printf("marshaling response: %v", err)
		}
		_, writeErr := w.Write(bytes)
		if writeErr != nil {
			log.Printf("Could not write response: %v", writeErr)
		}
	})
}

func main() {
	certPath := filepath.Join(tlsDir, tlsCertFile)
	keyPath := filepath.Join(tlsDir, tlsKeyFile)

	mux := http.NewServeMux()
	mux.Handle("/validate", validateServiceBindingRequest())
	server := &http.Server{
		Addr:    ":8443",
		Handler: mux,
	}
	log.Fatal(server.ListenAndServeTLS(certPath, keyPath))
}
