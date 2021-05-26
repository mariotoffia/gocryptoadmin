package parsers

import (
	"fmt"
	"strings"

	"github.com/mariotoffia/gocryptoadmin/common"
)

// ResolverParser parses a resolver expression on form
// _btx:USDT = btx,all:USDT-USD -> ofx,all:USD-EUR_.
//
// All patterns are terminated with a new-line.
//
// When prefix is omitted, the _all_ prefix is appended.
type ResolverParser struct {
	expressions []ResolverExpression
}

type ResolverExpressionPathItem struct {
	AssetPrefixes []string
	AssetPair     common.AssetPair
}

type ResolverExpression struct {
	Asset         common.AssetType
	AssetPrefixes []string
	Path          []ResolverExpressionPathItem
}

func NewResolverParser() *ResolverParser {
	return &ResolverParser{
		expressions: []ResolverExpression{},
	}
}

func (parser *ResolverParser) Parse(expr string) *ResolverParser {

	lines := strings.Split(strings.ReplaceAll(expr, "\r\n", "\n"), "\n")

	for _, line := range lines {

		expr := ResolverExpression{
			Path: []ResolverExpressionPathItem{},
		}

		eq := strings.Split(line, "=")

		if len(eq) != 2 {

			panic(
				fmt.Sprintf("expr: %s is not valid", line),
			)

		}

		assetPrefixes, assetExpr := parser.getPrefixes(parser.cleanString(eq[0]))

		asset := common.AssetType(assetExpr)
		expr.Asset = asset
		expr.AssetPrefixes = assetPrefixes

		paths := strings.Split(parser.cleanString(eq[1]), "->")
		for _, path := range paths {

			assetPairPrefixes, costUnitExpr := parser.getPrefixes(path)

			costUnit := common.AssetType(costUnitExpr)
			expr.Path = append(expr.Path, ResolverExpressionPathItem{
				AssetPrefixes: assetPairPrefixes,
				AssetPair: common.AssetPair{
					Asset:    asset,
					CostUnit: costUnit,
				},
			})

			asset = costUnit

		}

		parser.expressions = append(parser.expressions, expr)
	}

	return parser
}

func (parser *ResolverParser) GetExpressions() []ResolverExpression {
	return parser.expressions
}

func (parser *ResolverParser) cleanString(expr string) string {

	expr = strings.ReplaceAll(expr, " ", "")
	return strings.ReplaceAll(expr, "\t", "")

}

func (parser *ResolverParser) getPrefixes(expr string) ([]string, string) {

	c := strings.Split(expr, ":")

	if len(c) <= 1 {
		return []string{"all"}, expr
	}

	if len(c) != 2 {

		panic(
			fmt.Sprintf("expr: %s is not valid", expr),
		)

	}

	return strings.Split(c[0], ","), c[1]

}
