package mysqlwip

// Framework code is generated by the generator.

import (
	"fmt"

	"github.com/pingcap/tidb/parser/ast"
	"github.com/pkg/errors"

	"github.com/bytebase/bytebase/backend/plugin/advisor"
)

var (
	_ advisor.Advisor = (*ColumnDisallowChangingOrderAdvisor)(nil)
	_ ast.Visitor     = (*columnDisallowChangingOrderChecker)(nil)
)

// ColumnDisallowChangingOrderAdvisor is the advisor checking for disallow changing column order.
type ColumnDisallowChangingOrderAdvisor struct {
}

// Check checks for disallow changing column order.
func (*ColumnDisallowChangingOrderAdvisor) Check(ctx advisor.Context, _ string) ([]advisor.Advice, error) {
	stmtList, ok := ctx.AST.([]ast.StmtNode)
	if !ok {
		return nil, errors.Errorf("failed to convert to StmtNode")
	}

	level, err := advisor.NewStatusBySQLReviewRuleLevel(ctx.Rule.Level)
	if err != nil {
		return nil, err
	}
	checker := &columnDisallowChangingOrderChecker{
		level: level,
		title: string(ctx.Rule.Type),
	}

	for _, stmt := range stmtList {
		checker.text = stmt.Text()
		checker.line = stmt.OriginTextPosition()
		(stmt).Accept(checker)
	}

	if len(checker.adviceList) == 0 {
		checker.adviceList = append(checker.adviceList, advisor.Advice{
			Status:  advisor.Success,
			Code:    advisor.Ok,
			Title:   "OK",
			Content: "",
		})
	}
	return checker.adviceList, nil
}

type columnDisallowChangingOrderChecker struct {
	adviceList []advisor.Advice
	level      advisor.Status
	title      string
	text       string
	line       int
}

// Enter implements the ast.Visitor interface.
func (checker *columnDisallowChangingOrderChecker) Enter(in ast.Node) (ast.Node, bool) {
	if node, ok := in.(*ast.AlterTableStmt); ok {
		for _, spec := range node.Specs {
			if (spec.Tp == ast.AlterTableChangeColumn || spec.Tp == ast.AlterTableModifyColumn) &&
				spec.Position.Tp != ast.ColumnPositionNone {
				checker.adviceList = append(checker.adviceList, advisor.Advice{
					Status:  checker.level,
					Code:    advisor.ChangeColumnOrder,
					Title:   checker.title,
					Content: fmt.Sprintf("\"%s\" changes column order", checker.text),
					Line:    checker.line,
				})
				break
			}
		}
	}

	return in, false
}

// Leave implements the ast.Visitor interface.
func (*columnDisallowChangingOrderChecker) Leave(in ast.Node) (ast.Node, bool) {
	return in, true
}
