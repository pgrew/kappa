package skl

import (
	"reflect"
	"strings"
	"testing"

	"github.com/stretchr/testify/suite"
)

type TestCase struct {
	skip bool
	s    string
	stmt Statement
	err  string
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestParserTestSuite(t *testing.T) {
	suite.Run(t, new(ParserTestSuite))
}

// ParserTestSuite executes all the parser tests
type ParserTestSuite struct {
	suite.Suite
}

func (suite *ParserTestSuite) SetupTest() {
}

func (suite *ParserTestSuite) TearDownTest() {
}

func (suite *ParserTestSuite) validate(tests []TestCase) {
	for i, tt := range tests {
		if tt.skip {
			suite.T().Logf("skipping test of '%s'", tt.s)
			continue
		}
		stmt, err := NewParser(strings.NewReader(tt.s)).ParseStatement()

		if !reflect.DeepEqual(tt.err, errstring(err)) {
			suite.T().Errorf("%d. %q: error mismatch:\n  exp=%s\n  got=%s\n\n", i, tt.s, tt.err, err)
		} else if tt.err == "" && !reflect.DeepEqual(tt.stmt, stmt) {
			suite.T().Errorf("%d. %q\n\nstmt mismatch:\n\nexp=%#v\n\ngot=%#v\n\n", i, tt.s, tt.stmt, stmt)
		}
	}
}

// Ensure the parser will return an error for unknown statements
func (suite *ParserTestSuite) TestInvalidStatement() {
	var tests = []TestCase{

		// Errors
		{s: `a bad statement.`, err: `found a, expected USE, CREATE, SHOW, DROP at line 1, char 1`},
	}

	suite.validate(tests)
}

// Ensure the parser can parse strings into USE statements
func (suite *ParserTestSuite) TestUseStatement() {
	var tests = []TestCase{
		{
			s:    `USE acme.example`,
			stmt: &UseStatement{name: "acme.example"},
		},

		// Errors
		{s: `USE `, err: `found EOF, expected namespace at line 1, char 6`},
		{s: `USE acme.example.`, err: `found EOF, expected identifier at line 1, char 18`},
		{s: `USE acme.example. `, err: `found WS, expected identifier at line 1, char 18`},
		{s: `USE .example`, err: `found ., expected namespace at line 1, char 5`},
	}

	suite.validate(tests)
}

// Ensure the parser can parse strings into CREATE NAMESPACE statements
func (suite *ParserTestSuite) TestCreateNamespace() {
	var tests = []TestCase{
		{
			s:    `CREATE NAMESPACE acme`,
			stmt: &CreateNamespaceStatement{name: "acme"},
		},

		// Errors
		{s: `CREATE `, err: `found EOF, expected NAMESPACE at line 1, char 9`},
		{s: `CREATE NAMESPACE `, err: `found EOF, expected namespace at line 1, char 19`},
		{s: `CREATE NAMESPACE acme.example.`, err: `found EOF, expected identifier at line 1, char 31`},
		{s: `CREATE NAMESPACE acme.example. `, err: `found WS, expected identifier at line 1, char 31`},
		{s: `CREATE NAMESPACE .example`, err: `found ., expected namespace at line 1, char 18`},
	}

	suite.validate(tests)
}

// Ensure the parser can parse strings into DROP NAMESPACE statements
func (suite *ParserTestSuite) TestDropNamespace() {
	var tests = []TestCase{
		{
			s:    `DROP NAMESPACE acme`,
			stmt: &DropNamespaceStatement{name: "acme"},
		},

		// Errors
		{s: `DROP `, err: `found EOF, expected NAMESPACE at line 1, char 7`},
		{s: `DROP NAMESPACE `, err: `found EOF, expected namespace at line 1, char 17`},
		{s: `DROP NAMESPACE acme.example.`, err: `found EOF, expected identifier at line 1, char 29`},
		{s: `DROP NAMESPACE acme.example. `, err: `found WS, expected identifier at line 1, char 29`},
		{s: `DROP NAMESPACE .example`, err: `found ., expected namespace at line 1, char 16`},
	}

	suite.validate(tests)
}

// Ensure the parser can parse strings into SHOW NAMESPACES statements
func (suite *ParserTestSuite) TestShowNamespaces() {
	var tests = []TestCase{
		{
			s:    `SHOW NAMESPACES`,
			stmt: &ShowNamespacesStatement{},
		},

		// Errors
		{s: `SHOW `, err: `found EOF, expected NAMESPACES at line 1, char 7`},
		{s: `SHOW NAMESPACE`, err: `found NAMESPACE, expected NAMESPACES at line 1, char 6`},
	}

	suite.validate(tests)
}

// errstring converts an error to its string representation.
func errstring(err error) string {
	if err != nil {
		return err.Error()
	}
	return ""
}

func BenchmarkBadStatement(b *testing.B) {
	stmt := "a bad statement"
	for i := 0; i < b.N; i++ {
		NewParser(strings.NewReader(stmt)).ParseStatement()
	}
}

func BenchmarkUseStatement(b *testing.B) {
	stmt := "USE acme"
	for i := 0; i < b.N; i++ {
		NewParser(strings.NewReader(stmt)).ParseStatement()
	}
}

func BenchmarkNestedUseStatement(b *testing.B) {
	stmt := "USE acme.example.dynamite"
	for i := 0; i < b.N; i++ {
		NewParser(strings.NewReader(stmt)).ParseStatement()
	}
}

func BenchmarkCreateNamespaceStatement(b *testing.B) {
	stmt := "CREATE NAMESPACE acme"
	for i := 0; i < b.N; i++ {
		NewParser(strings.NewReader(stmt)).ParseStatement()
	}
}

func BenchmarkDropNamespaceStatement(b *testing.B) {
	stmt := "DROP NAMESPACE acme"
	for i := 0; i < b.N; i++ {
		NewParser(strings.NewReader(stmt)).ParseStatement()
	}
}

func BenchmarkShowNamespacesStatement(b *testing.B) {
	stmt := "SHOW NAMESPACES"
	for i := 0; i < b.N; i++ {
		NewParser(strings.NewReader(stmt)).ParseStatement()
	}
}
