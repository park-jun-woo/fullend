//ff:func feature=gen-gogin type=generator control=sequence
//ff:what main.go 생성 템플릿 문자열을 반환한다

package gogin

import "fmt"

// mainTemplate returns the fmt.Sprintf template for flat-mode cmd/main.go.
func mainTemplate(modulePath, authzImport, queueImport, authzInitBlock, queueInitBlock, initBlock, queueSubscribeBlock string) string {
	return fmt.Sprintf(`package main

import (
	"database/sql"
	"flag"
	"log"
	"net/http"

	_ "github.com/lib/pq"

	"%s/internal/model"
	"%s/internal/service"%s%s
)

func main() {
	addr := flag.String("addr", ":8080", "listen address")
	dsn := flag.String("dsn", "postgres://localhost:5432/app?sslmode=disable", "database connection string")
	dbDriver := flag.String("db", "postgres", "database driver (postgres, mysql)")
	flag.Parse()

	conn, err := sql.Open(*dbDriver, *dsn)
	if err != nil {
		log.Fatalf("database connection failed: %%v", err)
	}
	defer conn.Close()

	if err := conn.Ping(); err != nil {
		log.Fatalf("database ping failed: %%v", err)
	}
%s%s
	server := &service.Server{
%s
	}
%s
	handler := service.Handler(server)
	log.Printf("server listening on %%s", *addr)
	log.Fatal(http.ListenAndServe(*addr, handler))
}
`, modulePath, modulePath, authzImport, queueImport, authzInitBlock, queueInitBlock, initBlock, queueSubscribeBlock)
}
