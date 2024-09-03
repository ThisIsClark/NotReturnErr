package main

import (
	"bytes"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"log"
	"os"
	"path/filepath"
)

func formatNode(node ast.Node) string {
	var buf bytes.Buffer
	err := format.Node(&buf, token.NewFileSet(), node)
	if err != nil {
		panic(err)
	}
	return buf.String()
}

func isErrorType(expr ast.Expr) bool {
	ident, ok := expr.(*ast.Ident)
	return ok && ident.Name == "error"
}

func main() {
	fset := token.NewFileSet() // 文件集，用于解析

	err := filepath.Walk("E:\\code\\liuyu\\data\\ingestor", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			if info.Name() == "vendor" {
				return filepath.SkipDir
			}
			return nil
		}
		if filepath.Ext(path) != ".go" {
			return nil
		}

		// 解析文件，得到AST
		f, err := parser.ParseFile(fset, path, nil, parser.AllErrors)
		if err != nil {
			log.Println("Error parsing file:", err)
			return nil
		}

		// 遍历AST
		ast.Inspect(f, func(n ast.Node) bool {
			if n == nil {
				return true
			}
			switch stmt := n.(type) {
			case *ast.FuncDecl:
				if stmt.Type.Results != nil {
					// 函数没有返回值，不继续解析
					if len(stmt.Type.Results.List) == 0 {
						return false
					}
					// 函数返回值不包含error，不继续解析
					for _, result := range stmt.Type.Results.List {
						if !isErrorType(result.Type) {
							return false
						}
					}
				} else {
					return false
				}
			case *ast.IfStmt:
				if binExpr, ok := stmt.Cond.(*ast.BinaryExpr); ok {
					if ident, ok := binExpr.X.(*ast.Ident); ok {
						if ident.Name != "err" {
							break
						}
					} else {
						break
					}
					if binExpr.Op != token.NEQ {
						break
					}
					if ident, ok := binExpr.Y.(*ast.Ident); ok {
						if ident.Name != "nil" {
							break
						}
					} else {
						break
					}
				} else {
					break
				}
				// 判断是否有return err
			LOOP:
				for _, stmt := range stmt.Body.List {
					if returnStmt, ok := stmt.(*ast.ReturnStmt); ok {
						for _, expr := range returnStmt.Results {
							if ident, ok := expr.(*ast.Ident); ok && ident.Name == "err" {
								break LOOP
							}
							// 判断是否有return fmt.Errorf("%s", err)和return fmt.Errorf("%s", err.Error())
							if callExpr, ok := expr.(*ast.CallExpr); ok {
								if selExpr, ok := callExpr.Fun.(*ast.SelectorExpr); ok {
									if selExpr.Sel.Name == "Errorf" {
										for _, arg := range callExpr.Args {
											if ident, ok := arg.(*ast.Ident); ok && ident.Name == "err" {
												break LOOP
											}
											if call, ok := arg.(*ast.CallExpr); ok {
												if sel, ok := call.Fun.(*ast.SelectorExpr); ok {
													if x, ok := sel.X.(*ast.Ident); ok && x.Name == "err" && sel.Sel.Name == "Error" {
														break LOOP
													}
												}
											}
										}
									}
								}
							}
						}
						startPos := fset.Position(returnStmt.Pos())
						endPos := fset.Position(returnStmt.End())
						log.Printf("找到了if err != nil中未返回error的return 语句: %s，开始位置: %s，结束位置: %s\n", formatNode(returnStmt), startPos, endPos)
					}
				}
			}
			return true
		})

		return nil
	})

	if err != nil {
		log.Fatal(err)
	}
}