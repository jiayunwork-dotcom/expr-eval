package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"expr-eval/internal/compile"
	"expr-eval/internal/eval"
	"expr-eval/internal/parser"
	"expr-eval/internal/validate"
	"expr-eval/internal/vm"
)

type Config struct {
	Addr string
}

func New(cfg Config) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/health", handleHealth)
	mux.HandleFunc("/api/eval", handleEval)
	mux.HandleFunc("/api/validate", handleValidate)
	mux.HandleFunc("/api/compile", handleCompile)
	mux.HandleFunc("/api/batch", handleBatch)
	return mux
}

func ListenAndServe(cfg Config) error {
	mux := New(cfg)
	return http.ListenAndServe(cfg.Addr, mux)
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

type evalRequest struct {
	Expr string             `json:"expr"`
	Vars map[string]float64 `json:"vars"`
}

type evalResponse struct {
	Result float64 `json:"result"`
}

func handleEval(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		httpError(w, http.StatusMethodNotAllowed, "POST required")
		return
	}
	var req evalRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpError(w, http.StatusBadRequest, "invalid JSON: "+err.Error())
		return
	}
	if req.Expr == "" {
		httpError(w, http.StatusBadRequest, "expr is empty")
		return
	}
	node, err := parser.Parse(req.Expr)
	if err != nil {
		httpError(w, http.StatusBadRequest, "parse error: "+err.Error())
		return
	}
	result, err := eval.Eval(node, req.Vars)
	if err != nil {
		httpError(w, http.StatusBadRequest, "eval error: "+err.Error())
		return
	}
	writeJSON(w, http.StatusOK, evalResponse{Result: result})
}

type validateRequest struct {
	Expr string   `json:"expr"`
	Vars []string `json:"vars"`
}

type validateResponse struct {
	OK     bool     `json:"ok"`
	Errors []string `json:"errors,omitempty"`
}

func handleValidate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		httpError(w, http.StatusMethodNotAllowed, "POST required")
		return
	}
	var req validateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpError(w, http.StatusBadRequest, "invalid JSON: "+err.Error())
		return
	}
	if req.Expr == "" {
		httpError(w, http.StatusBadRequest, "expr is empty")
		return
	}
	node, err := parser.Parse(req.Expr)
	if err != nil {
		httpError(w, http.StatusBadRequest, "parse error: "+err.Error())
		return
	}
	funcs := validate.DefaultFuncs()
	result := validate.Validate(node, req.Vars, funcs)
	resp := validateResponse{OK: result.OK()}
	for _, e := range result.Errors() {
		resp.Errors = append(resp.Errors, e.Message)
	}
	writeJSON(w, http.StatusOK, resp)
}

type compileRequest struct {
	Expr string             `json:"expr"`
	Vars map[string]float64 `json:"vars"`
}

type compileResponse struct {
	Instructions int     `json:"instructions"`
	Result       float64 `json:"result"`
}

func handleCompile(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		httpError(w, http.StatusMethodNotAllowed, "POST required")
		return
	}
	var req compileRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpError(w, http.StatusBadRequest, "invalid JSON: "+err.Error())
		return
	}
	if req.Expr == "" {
		httpError(w, http.StatusBadRequest, "expr is empty")
		return
	}
	node, err := parser.Parse(req.Expr)
	if err != nil {
		httpError(w, http.StatusBadRequest, "parse error: "+err.Error())
		return
	}
	prog, err := compile.Compile(node)
	if err != nil {
		httpError(w, http.StatusBadRequest, "compile error: "+err.Error())
		return
	}
	v := vm.New()
	v.RegisterDefaults()
	result, err := v.Run(prog, req.Vars)
	if err != nil {
		httpError(w, http.StatusBadRequest, "vm error: "+err.Error())
		return
	}
	writeJSON(w, http.StatusOK, compileResponse{
		Instructions: len(prog.Instructions),
		Result:       result,
	})
}

type batchRequest struct {
	Expressions []string           `json:"expressions"`
	Vars        map[string]float64 `json:"vars"`
}

type batchResult struct {
	Expr   string   `json:"expr"`
	Result *float64 `json:"result,omitempty"`
	Error  string   `json:"error,omitempty"`
}

func handleBatch(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		httpError(w, http.StatusMethodNotAllowed, "POST required")
		return
	}
	var req batchRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpError(w, http.StatusBadRequest, "invalid JSON: "+err.Error())
		return
	}
	if len(req.Expressions) == 0 {
		httpError(w, http.StatusBadRequest, "expressions array is empty")
		return
	}
	results := make([]batchResult, len(req.Expressions))
	for i, expr := range req.Expressions {
		br := batchResult{Expr: expr}
		node, err := parser.Parse(expr)
		if err != nil {
			br.Error = "parse: " + err.Error()
			results[i] = br
			continue
		}
		val, err := eval.Eval(node, req.Vars)
		if err != nil {
			br.Error = "eval: " + err.Error()
			results[i] = br
			continue
		}
		br.Result = &val
		results[i] = br
	}
	writeJSON(w, http.StatusOK, map[string]interface{}{"results": results})
}

func httpError(w http.ResponseWriter, code int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]string{"error": msg})
}

func writeJSON(w http.ResponseWriter, code int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	enc.Encode(v)
}

func ParsePort(addr string) int {
	parts := strings.Split(addr, ":")
	if len(parts) < 2 {
		return 0
	}
	p, _ := strconv.Atoi(parts[len(parts)-1])
	return p
}

func FormatAddr(addr string) string {
	port := ParsePort(addr)
	if port == 0 {
		return addr
	}
	return fmt.Sprintf("http://0.0.0.0:%d", port)
}
