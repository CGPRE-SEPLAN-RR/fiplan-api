# Documento de Apoio aos Desenvolvedores

Este documento trata de algumas práticas comuns no desenvolvimento desse projeto, e deve servir de guia inicial para começar o desenvolvimento de novas funcionalidades ou manutenção das já existentes. Ele será muito alterado nos estágios iniciais do projeto, onde a arquitetura do sistema e os processos de desenvolvimento padrão não estão precisamente definidos.

## O que é importante que eu saiba?

- Em Go, os tipos, variáveis e constantes definidos no escopo global do arquivo apenas são exportados se iniciarem com letra maíúscula
- É importante se ater ao tipo de dado correto para utilizar o mínimo de memória possível, como por exemplo:
    - Usar `uint16` para anos de exercício, já que o ano tem sempre 4 dígitos e o uint16 permite valores de 0 a 65536

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
│   │   └── database.go
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

- `bin`: ondem ficam todos os executáveis do projeto (`main`);
- `cmd`: contém um código mínimo, responsável por iniciar o servidor. Essa pasta não deve sofrer muitas modificações;
- `docs`: contém a documentação autogerada pelo [swag](https://github.com/swaggo/swag). Essa pasta só deve ser alterada pela execução do comando `make docs`;
- `internal`: contém a maior parte do código-fonte. Essa é a pasta na qual você mais vai mexer como desenvolvedor;
    - `internal/database`: contém o código-fonte de conexão com o banco de dados;
    - `internal/server`: contém o código-fonte relativo ao servidor, como as suas rotas (`routes.go`), os controladores das rotas (terminados em `_handler.go`), seus *middlewares* (`server.go`);
- `test`: contém os testes do projeto. Essa é a pasta na qual você mais vai mexer como QA;
- `tmp`: contém os arquivos temporários criados pelo [air](https://github.com/cosmtrek/air) (responsável pelo *live reload*). Essa pasta pode ser ignorada;

Dada essa estrutura, é importante seguí-la para facilitar o desenvolvimento e manutenção futuros.

No mais, é importante notar que cada controlador contém em um único arquivo todas informações necessárias para entender seu funcionamento, que são, nessa ordem:

- Modelo de dados
- Definição dos parâmetros
- Validação dos parâmetros
- *Queries* SQL
- Consultas ao banco de dados
- Lógica adicional

## Como criar um novo teste?

O testes serão realizados a partir de relatórios gerados pelo FIPLAN, inicialmente, a fim de garantir que os dados entregados sejam os mesmos. Para isso, usaremos Python e Selenium, mas como no momento não há nenhum teste escrito, esse passo a passo está incompleto.

## Como criar um novo *endpoint*?

Para adicionar um novo *endpoint*, é necessário que seja criados os seguintes componentes:

- [Rota](#rota)
- [Controlador](#controlador)

### Rota

Para criar a rota do seu *endpoint*, basta criar uma nova linha no arquivo `routes.go` da pasta `internal/server`. O trecho de código para esta nova rota será semelhante ao seguinte exemplo:

```go
// Você não deve criar essa função, ela já está disponível no arquivo `routes.go`
func (s *Server) RegisterRoutes() http.Handler {
    ...
    e.GET("/exemplo", s.ControladorDesseEndpoint)
    ...
}
```

Certifique-se de escolher um nome descritivo para o *endpoint* (por exemplo, `/seu-endpoint`) que não entre em conflito com outros *endpoints* definidos neste arquivo. Além disso, observe que o método HTTP para este *endpoint* deve ser `GET` ou `HEAD`, pois esta API é designada exclusivamente para consultas de dados no FIPLAN.

É essencial evitar conflitos de nomenclatura entre os `endpoints` neste arquivo. Caso opte por um mesmo nome, certifique-se de que o método HTTP associado seja diferente para evitar ambiguidades.

Por último, é fundamental que o nome do controlador (`s.ControladorDoSeuEndpoint`) seja o mesmo do controlador que será criado posteriormente. Isso garante a correta associação entre a rota e seu manipulador correspondente. Continue abaixo para criar o controlador correspondente a este *endpoint*.

### Controlador

1. Acesse o relatório no FIPLAN e se atenha aos campos que devem ser fornecidos como entrada (eles serão os parâmetros do *endpoint*)
2. Acesse o relatório no FIPLAN e se atenha aos campos que são extraídos (eles serão a base para montar o modelo de dados)
3. Acesse o código-fonte do FIPLAN e procure pelos arquivos que são usados na consulta daquele relatório (eles serão a base para a validação dos parâmetros, as *queries* SQL e para qualquer lógica adicional requisitada pelo relatório)
4. Crie um arquivo terminado em `_handler.go` na pasta `internal/server` para conter a lógica do seu controlador
5. Documente o relatório de acordo com o proposto pelo [swag](https://github.com/swaggo/swag)
6. Codifique o modelo de dados do relatório
7. Codifique os parâmetros e suas validações
8. Codifique os templates de *query* SQL
9. Codifique qualquer lógica adicional pendente
10. Teste o novo controlador para verificar se os dados são iguais aos obtidos no FIPLAN
