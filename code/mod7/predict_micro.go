package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"

	"cloud.google.com/go/bigquery"
	"github.com/gorilla/mux"
	"google.golang.org/api/iterator"
)

func defaultHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Serving:", r.URL.Path, "from", r.Host)
	w.WriteHeader(http.StatusOK)
	Body := "Enter ID to predict.\n"
	fmt.Fprintf(w, "%s", Body)
}

func predictQuery(recidId string, w http.ResponseWriter) error {
	projectID := "msds434-mod7"
	ctx := context.Background()
	client, err := bigquery.NewClient(ctx, projectID)
	if err != nil {
		return fmt.Errorf("bigquery.NewClient: %v", err)
	}
	defer client.Close()

	qtxt := fmt.Sprintf(
		"SELECT * FROM ML.PREDICT(MODEL `recidivism.recid_xgb_model`, "+
			"(SELECT * FROM `recidivism.test` "+
			"WHERE id = %s)) ", recidId)

	if recidId == "all" {
		qtxt = fmt.Sprintf(
			"SELECT * FROM ML.PREDICT(MODEL `recidivism.recid_xgb_model`, " +
				"(SELECT * FROM `recidivism.test`))")
	}

	q := client.Query(qtxt)

	// Location must match that of the dataset(s) referenced in the query.
	q.Location = "US"
	// Run the query and print results when the query job is completed.
	job, err := q.Run(ctx)
	if err != nil {
		return err
	}
	status, err := job.Wait(ctx)
	if err != nil {
		return err
	}
	if err := status.Err(); err != nil {
		return err
	}
	it, err := job.Read(ctx)
	for {
		var row []bigquery.Value
		err := it.Next(&row)
		if err == iterator.Done {
			break
		}
		if err != nil {
			return err
		}
		fmt.Fprintln(w, row)
	}
	return nil
}

// GetPredictionHandler is used to get the prediction result for a given ID
func getPredictionHandler(w http.ResponseWriter, r *http.Request) {
	paramStr := strings.Split(r.URL.Path, "/")

	recidID := paramStr[len(paramStr)-1]
	fmt.Fprintf(w, "Prediction results for %s:\n", recidID)
	err := predictQuery(recidID, w)
	if err != nil {
		fmt.Fprintf(w, "%s", err)
	}
}

func main() {
	// Create new Router
	router := mux.NewRouter()

	// route properly to respective handlers
	router.Handle("/", http.HandlerFunc(defaultHandler)).Methods("GET")
	router.Handle("/{id:[0-9]+}", http.HandlerFunc(getPredictionHandler)).Methods("GET")
	router.Handle("/all", http.HandlerFunc(getPredictionHandler)).Methods("GET")

	// Create new server and assign the router
	server := http.Server{
		Addr:    ":8080",
		Handler: router,
	}
	fmt.Println("Staring Recidivism Prediction server on Port 8080")
	// Start Server on defined port/host.
	server.ListenAndServe()
}
