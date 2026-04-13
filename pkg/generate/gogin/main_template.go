//ff:func feature=gen-gogin type=generator control=sequence topic=output
//ff:what domain 모드 main.go 생성 템플릿 문자열을 반환한다

package gogin

import "fmt"

// mainWithDomainsTemplate returns the fmt.Sprintf template for domain-mode cmd/main.go.
// dbName: 기본 DB 이름 (DATABASE_URL env 미지정 시 fallback 에 삽입).
func mainWithDomainsTemplate(osImport, importBlock, queueImport, builtinImport, jwtFlagLine, authzBlock, queueInitBlock, builtinInitBlock, initBlock, queueSubscribeBlock, dbName string) string {
	return fmt.Sprintf(`package main

import (
	"database/sql"
	"flag"
	"log"%s

	_ "github.com/lib/pq"
%s%s%s
)

func main() {
	addr := flag.String("addr", ":8080", "listen address")
	dsnDefault := os.Getenv("DATABASE_URL")
	if dsnDefault == "" {
		dsnDefault = "postgres://localhost:5432/%s?sslmode=disable"
	}
	dsn := flag.String("dsn", dsnDefault, "database connection string (or DATABASE_URL env)")
	dbDriver := flag.String("db", "postgres", "database driver (postgres, mysql)")%s
	flag.Parse()

	conn, err := sql.Open(*dbDriver, *dsn)
	if err != nil {
		log.Fatalf("database connection failed: %%v", err)
	}
	defer conn.Close()

	if err := conn.Ping(); err != nil {
		log.Fatalf("database ping failed: %%v", err)
	}
%s%s%s
	server := &service.Server{
%s
	}
%s
	r := service.SetupRouter(server)
	log.Printf("server listening on %%s", *addr)
	log.Fatal(r.Run(*addr))
}
`, osImport, importBlock, queueImport, builtinImport, dbName, jwtFlagLine, authzBlock, queueInitBlock, builtinInitBlock, initBlock, queueSubscribeBlock)
}
