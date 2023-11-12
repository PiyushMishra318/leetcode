package main

import (
	"bufio"
	"fmt"
	"github.com/ktr0731/go-fuzzyfinder"
	"github.com/spf13/cobra"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func main() {
	var rootCmd = &cobra.Command{Use: "leetcode-cli"}
	rootCmd.AddCommand(addCmd, testCmd, historyCmd, genTestCmd)
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

var addCmd = &cobra.Command{
	Use:   "add [problem name] [problem number]",
	Short: "Add a new LeetCode problem",
	Args:  cobra.ExactArgs(2), // Ensures exactly two arguments are passed
	Run: func(cmd *cobra.Command, args []string) {
		problemName := args[0]
		problemNumber := args[1]

		fileName := fmt.Sprintf("%s_%s.go", problemNumber, problemName)
		filePath := filepath.Join("problems", fileName)

		// Check if file exists
		if _, err := os.Stat(filePath); err == nil {
			fmt.Printf("File already exists: %s\n", filePath)
		} else if os.IsNotExist(err) {
			// Create file if it does not exist
			file, err := os.Create(filePath)
			if err != nil {
				fmt.Printf("Error creating file: %s\n", err)
				return
			}
			file.Close()
			fmt.Printf("Created file: %s\n", filePath)
		} else {
			fmt.Printf("Error checking file: %s\n", err)
			return
		}

		// Open file in Neovim
		openTestFileInEditor(filePath)
	},
}

var testCmd = &cobra.Command{
	Use:   "test [problem name]",
	Short: "Test a LeetCode problem",
	Run: func(cmd *cobra.Command, args []string) {
		err := runTestCommand()
		if err != nil {
			fmt.Println("Error:", err)
		}
	},
}

func runTestCommand() error {
	files, err := filepath.Glob("problems/*.go")
	if err != nil {
		fmt.Printf("Error listing problem files: %s\n", err)
		return err
	}

	idx, err := fuzzyfinder.Find(
		files,
		func(i int) string {
			return files[i]
		},
		fuzzyfinder.WithPromptString("Select a problem file:"),
	)
	if err != nil {
		fmt.Printf("Error selecting file: %s\n", err)
		return err
	}

	selectedFile := files[idx]

	return testProblem(selectedFile)
}

func testProblem(problemName string) error {
	_, fileName := filepath.Split(problemName)
	testDir := "tests" // Creates a 'test' subdirectory in the same directory as the original file
	testFileName := strings.TrimSuffix(fileName, ".go") + "_test.go"
	testFileName = filepath.Join(testDir, testFileName)

	cmd := exec.Command("go", "test", "-v", testFileName)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

var historyCmd = &cobra.Command{
	Use:   "history [problem name]",
	Short: "View submission history for a problem",
	Run: func(cmd *cobra.Command, args []string) {
		// Implementation for viewing history
	},
}

var genTestCmd = &cobra.Command{
	Use:   "gentest",
	Short: "Generate tests for a LeetCode problem",
	Run: func(cmd *cobra.Command, args []string) {
		// 1. Implement Fuzzy Finder to choose a problem file
		files, err := filepath.Glob("problems/*.go")
		if err != nil {
			fmt.Printf("Error listing problem files: %s\n", err)
			return
		}

		idx, err := fuzzyfinder.Find(
			files,
			func(i int) string {
				return files[i]
			},
			fuzzyfinder.WithPromptString("Select a problem file:"),
		)
		if err != nil {
			fmt.Printf("Error selecting file: %s\n", err)
			return
		}

		selectedFile := files[idx]

		// 2. Analyze the Go file
		fset := token.NewFileSet()
		node, err := parser.ParseFile(fset, selectedFile, nil, parser.ParseComments)
		if err != nil {
			fmt.Printf("Error parsing file: %s\n", err)
			return
		}

		var testFunctions []string

		for _, f := range node.Decls {
			fn, ok := f.(*ast.FuncDecl)
			if ok && fn.Name.IsExported() {
				if ok && fn.Name.IsExported() {
					// Generate a basic test function template
					testFunc := generateTestFunctionTemplate(fn)
					testFunctions = append(testFunctions, testFunc)
				}
			}
		}

		// Generate the content for the test file
		testFileContent := strings.Join(testFunctions, "\n\n")

		// Create the test file
		_, fileName := filepath.Split(selectedFile)
		testDir := "tests" // Creates a 'test' subdirectory in the same directory as the original file
		testFileName := strings.TrimSuffix(fileName, ".go") + "_test.go"
		testFileName = filepath.Join(testDir, testFileName)
		file, err := os.Create(testFileName)
		if err != nil {
			fmt.Printf("Error creating test file: %s\n", err)
			return
		}
		defer file.Close()

		writer := bufio.NewWriter(file)
		_, err = writer.WriteString(testFileContent)
		if err != nil {
			fmt.Printf("Error writing to test file: %s\n", err)
			return
		}
		writer.Flush()

		fmt.Printf("Test file created: %s\n", testFileName)
		// Open the test file in neovim
		openTestFileInEditor(testFileName)
	},
}

func openTestFileInEditor(filePath string) {
	cmd := exec.Command("nvim", filePath)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		fmt.Printf("Error opening file in editor: %s\n", err)
	}
}

func generateTestFunctionTemplate(fn *ast.FuncDecl) string {
	var params []string
	for _, p := range fn.Type.Params.List {
		for _, n := range p.Names {
			params = append(params, n.Name+": "+exprToString(p.Type))
		}
	}

	var results []string
	if fn.Type.Results != nil {
		for _, r := range fn.Type.Results.List {
			results = append(results, exprToString(r.Type))
		}
	}

	testFuncName := "Test" + fn.Name.Name
	importStatement := "package main\n\nimport (\n\t\"testing\"\n\t\"reflect\"\n\t\"leetcode/problems\"\n)\n\n"
	testBody := fmt.Sprintf(`func %s(t *testing.T) {
    // Define test cases here
    var testCases []struct {
        inputs []interface{} // Modify this as per your function's input types
        want interface{}     // Modify this as per your function's return type
    }

    for _, tc := range testCases {
        got := %s(tc.inputs...) // Assuming the function can handle variadic inputs
        if !reflect.DeepEqual(got, tc.want) {
            t.Errorf("For inputs %%v, got %%v, want %%v", tc.inputs, got, tc.want)
        }
    }
}`, testFuncName, fn.Name.Name)

	return importStatement + testBody
}

func exprToString(expr ast.Expr) string {
	switch e := expr.(type) {
	case *ast.Ident:
		return e.Name
	case *ast.ArrayType:
		return "[]" + exprToString(e.Elt)
	case *ast.StarExpr:
		return "*" + exprToString(e.X)
	// ... add cases for other types as needed
	default:
		return "interface{}" // a fallback type
	}
}

// Add more commands and implementation as needed
