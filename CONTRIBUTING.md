# Documento de Apoio aos Desenvolvedores

Este documento trata de algumas práticas comuns no desenvolvimento desse projeto, e deve servir de guia inicial para começar o desenvolvimento de novas funcionalidades ou manutenção das já existentes.

## Como está estruturado o projeto?

O projeto está estruturado da seguinte forma:

```bash
.
├── bin
├── cmd
│   └── api
│       └── main.go
├── docs
│   ├── docs.go
│   ├── swagger.json
│   └── swagger.yaml
├── go.mod
├── go.sum
├── internal
│   ├── database
│   │   ├── database.go
│   │   └── queries
│   │       └── relatorio_fip215.sql
│   ├── model
│   │   ├── conta.go
│   │   └── relatorio.go
│   └── server
│       ├── conta_handler.go
│       ├── relatorio_fip215_handler.go
│       ├── routes.go
│       ├── server.go
│       └── util.go
├── Makefile
├── README.md
├── test
└── tmp
```

As pastas têm os seguintes objetivos:

- `bin`: ondem ficam todos os executáveis do projeto (e. g. main);
- `cmd`: contém um código mínimo, responsável por iniciar o servidor. Essa pasta não deve sofrer muitas modificações;
- `docs`: contém a documentação autogerada pelo [swag](https://github.com/swaggo/swag). Essa pasta só deve ser alterada pela execução do comando `make docs`;
- `internal`: contém a maior parte do código-fonte. Essa é a pasta na qual você mais vai mexer como desenvolvedor;
    - `internal/database`: contém o código-fonte de conexão com o banco de dados;
    - `internal/database/queries`: contém os *scripts* SQL que serão usados para comunicação com o banco de dados. Note que todos os *scripts* se tratam de uma única *query*, composta por um ou mais `SELECT`s. Não deve haver nenhum *script* DDL, 
    - `internal/model`: contém os modelos dos dados extraídos do banco de dados, em forma de *structs*;
    - `internal/server`: contém o código-fonte relativo ao servidor, como as suas rotas (`routes.go`), os controladores das rotas (terminados em `_handler.go`), seus *middlewares* (`server.go`);
- `test`: contém os testes do projeto. Essa é a pasta na qual você mais vai mexer como QA;
- `tmp`: contém os arquivos temporários criados pelo [air](https://github.com/cosmtrek/air) (responsável pelo *live reload*). Essa pasta pode ser ignorada;

Dada essa estrutura, é importante seguí-la para facilitar o desenvolvimento e manutenção futuros.

## [EM CONSTRUÇÃO] Como criar um novo teste?

## [EM CONSTRUÇÃO] Como criar um novo *endpoint*?

## [EM CONSTRUÇÃO] Como adicionar suporte a um relatório?
