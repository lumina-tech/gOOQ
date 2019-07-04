package main

import (
	"fmt"
	"log"
	"strings"
	"text/template"

	"github.com/lumina-tech/gooq/common"
	"github.com/lumina-tech/gooq/meta"
)

type PredicateInfo struct {
	Predicate     string
	FieldFunction string
	JoinFunction  string
}

var preds = []PredicateInfo{
	PredicateInfo{Predicate: "EqPredicate", FieldFunction: "Eq", JoinFunction: "IsEq"},
	PredicateInfo{Predicate: "GtPredicate", FieldFunction: "Gt", JoinFunction: "IsGt"},
	PredicateInfo{Predicate: "GtePredicate", FieldFunction: "Gte", JoinFunction: "IsGte"},
	PredicateInfo{Predicate: "ILikePredicate", FieldFunction: "ILike", JoinFunction: "IsILike"},
	PredicateInfo{Predicate: "LikePredicate", FieldFunction: "Like", JoinFunction: "IsLike"},
	PredicateInfo{Predicate: "LtPredicate", FieldFunction: "Lt", JoinFunction: "IsLt"},
	PredicateInfo{Predicate: "LtePredicate", FieldFunction: "Lte", JoinFunction: "IsLte"},
	PredicateInfo{Predicate: "NotEqPredicate", FieldFunction: "NotEq", JoinFunction: "IsNotEq"},
	PredicateInfo{Predicate: "NotILikePredicate", FieldFunction: "NotILike", JoinFunction: "IsNotILike"},
	PredicateInfo{Predicate: "NotLikePredicate", FieldFunction: "NotLike", JoinFunction: "IsNotLike"},
	PredicateInfo{Predicate: "NotSimilarPredicate", FieldFunction: "NotSimilar", JoinFunction: "IsNotSimilar"},
	PredicateInfo{Predicate: "SimilarPredicate", FieldFunction: "Simiar", JoinFunction: "IsSimilar"},
}

func argifier(arg meta.FunctionInfo) string {
	argN := strings.Count(arg.Expr, "%")
	if argN == 1 {
		return ""
	} else {
		argsToSub := argN - 1
		args := make([]string, argsToSub)
		for i := 0; i < argsToSub; i++ {
			args[i] = fmt.Sprintf("_%d", i)
		}
		return strings.Join(args, ",")
	}
}

func signifier(arg meta.FunctionInfo) string {
	args := argifier(arg)
	if args == "" {
		return ""
	} else {
		return fmt.Sprintf("%s interface{}", args)
	}
}

func injectifier(arg meta.FunctionInfo) string {
	args := argifier(arg)
	if args == "" {
		return ""
	} else {
		return fmt.Sprintf(", %s", args)
	}
}

func main() {
	params := make(map[string]interface{})
	params["types"] = meta.Types
	params["predicates"] = preds
	params["functions"] = meta.Funcs
	m := template.FuncMap{
		"toLower":     strings.ToLower,
		"signifier":   signifier,
		"injectifier": injectifier,
	}
	t, err := template.New("fields.go.tmpl").Funcs(m).ParseFiles("./templates/fields.go.tmpl")
	if err != nil {
		log.Fatal(err)
	}
	if err := common.RenderToFile(t, "./gooq/fields.go", params); err != nil {
		log.Fatal(err)
	}
	log.Println("Regenerated fields")
}
