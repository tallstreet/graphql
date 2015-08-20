package handler

import (
	"encoding/json"
	"io"
	"log"
	"strings"
	"net/http"

	
	"github.com/tallstreet/graphql/ast"
	"github.com/tallstreet/graphql/executor"
	//"github.com/tallstreet/graphql/executor/tracer"
	"github.com/tallstreet/graphql/parser"
	"sevki.org/lib/prettyprint"
	//"golang.org/x/net/context"
)

// Error represents an error the occured while parsing a graphql query or while generating a response.
type Error struct {
	Message string `json:"message"`
}

// Result represents a relay query.
type Request struct {
	Query      string      `json:"query,omitempty"`
	Variables  interface{} `json:"variables,omitempty"`
}

// Result represents a graphql query result.
type Result struct {
	Trace interface{} `json:"__trace_info,omitempty"`
	Data  interface{} `json:"data,omitempty"`
	Error *Error      `json:"error,omitempty"`
}

// ExecutorHandler makes a executor.Executor querable via HTTP
type ExecutorHandler struct {
	executor *executor.Executor
}

// New constructs a ExecutorHandler from a executor.
func New(executor *executor.Executor) *ExecutorHandler {
	return &ExecutorHandler{executor: executor}
}

func writeErr(w io.Writer, err error) {
	writeJSON(w, Result{Error: &Error{Message: err.Error()}})
}
func writeJSON(w io.Writer, data interface{}) {
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Println("error writing json response:", err)
		// attempt to write error
		writeErr(w, err)
	}
}

func writeJSONIndent(w io.Writer, data interface{}, indentString string) {
	b, err := json.MarshalIndent(data, "", indentString)
	if err != nil {
		log.Println("error encoding json response:", err)
		writeErr(w, err)
	}
	if _, err := w.Write(b); err != nil {
		log.Println("error writing json response:", err)
		writeErr(w, err)
	}
}

// ServeHTTP provides an entrypoint into a graphql executor. It pulls the query from
// the 'q' GET parameter.
func (h *ExecutorHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", r.Header.Get("Access-Control-Request-Headers"))
	if r.Method == "OPTIONS" {
		w.WriteHeader(200)
		return
	}
	
	decoder := json.NewDecoder(r.Body)
  var qreq Request   
  err := decoder.Decode(&qreq)
	if err != nil {
		log.Println("error parsing:", err)
		writeErr(w, err)
		return
	}
	q := qreq.Query
	
	/*
	//TODO(tallstreet): reject non-GET/OPTIONS requests
	q := r.URL.Query().Get("q")
	*/
	var doc ast.Document
	if err := parser.New("graphql", strings.NewReader(q)).Decode(&doc); err != nil {
		
		log.Printf(err.Error())
		
	} else {
		//		fs, _ := os.Open("../tests/complex-as-possible.graphql")
		//		e, _ := ioutil.ReadAll(fs)
		//		log.Printf(string(e))
		log.Printf(prettyprint.AsJSON(doc))
	}
	
	/*
	log.Println("query:", q)
	operation, err := parser.ParseOperation([]byte(q))
	if err != nil {
		log.Println("error parsing:", err)
		writeErr(w, err)
		return
	}
	

	asjson, _ := json.MarshalIndent(operation, "", " ")
	log.Println(string(asjson))
	// if err := h.validator.Validate(operation); err != nil { writeErr(w, err); return }
	ctx := context.Background()
	if r.Header.Get("X-Trace-ID") != "" {
		t, err := tracer.FromRequest(r)
		if err == nil {
			ctx = tracer.NewContext(ctx, t)
		}
	}
	ctx = context.WithValue(ctx, "http_request", r)
	if r.Header.Get("X-GraphQL-Only-Parse") == "1" {
		writeJSONIndent(w, operation, " ")
		return
	}

	data, err := h.executor.HandleOperation(ctx, operation)
	result := Result{Data: data}
	if err != nil {
		w.WriteHeader(400)
		result.Error = &Error{Message: err.Error()}
	}
	if t, ok := tracer.FromContext(ctx); ok {
		t.Done()
		result.Trace = t
	}
	*/

	writeJSONIndent(w, doc, "  ")
}
