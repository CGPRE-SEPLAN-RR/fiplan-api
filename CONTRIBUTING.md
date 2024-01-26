# Documento de Apoio aos Desenvolvedores

Este documento trata de algumas práticas comuns no desenvolvimento desse projeto, e deve servir de guia inicial para começar o desenvolvimento de novas funcionalidades ou manutenção das já existentes. Ele será muito alterado nos estágios iniciais do projeto, onde a arquitetura do sistema e os processos de desenvolvimento padrão não estão precisamente definidos.

A fazer:

- [ ] Ser mais descritivo na construção de [Modelo](#modelo) e [Controlador](#controlador)
    - [ ] Como fazer a *godoc* do jeito certo para gerar o resultado pretendido ao executar o `make docs` (Consultar documentação do [swag](https://github.com/swaggo/swag))?
    - [ ] Como avaliar quais são os tipos de dados que devo definir no meu modelo (Consultar documentação do [godror](https://github.com/godror/godror))?
    - [ ] Como fazer a validação dos parâmetros corretamente?
    - [ ] Como saber qual o tipo de dado dos parâmetros que vou utilizar?
    - [ ] Quais os tipos de consulta que podem ser realizados?

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

## Como criar um novo teste?

Para criar um novo teste, você deve criar um novo arquivo terminado em `_test.go` na pasta `test` e inserir a lógica do seu teste nele. Alguns recursos que podem ajudar na escrita de novos testes são:

- [Aprenda Go com Testes](https://larien.gitbook.io/aprenda-go-com-testes/)
- [Test Driven Development: TDD Simples e Prático](https://www.devmedia.com.br/test-driven-development-tdd-simples-e-pratico/18533)

## Como criar um novo *endpoint*?

Para adicionar um novo *endpoint*, é necessário que seja criados os seguintes componentes:

- [Rota](#rota)
- [Modelo](#modelo)
- [Controlador](#controlador)
- [Query SQL](#query-sql) (Se necessário)

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

### Modelo

No contexto de desenvolvimento de software, um modelo refere-se a uma representação estruturada de dados que descreve as propriedades e comportamentos de uma entidade específica. Aqui, usaremos como exemplo o modelo [`RelatorioFIP215`](internal/model/relatorio.go#L15).

#### `DadoRelatorioFIP215`

O tipo [`DadoRelatorioFIP215`](internal/model/relatorio.go#L3) é um modelo que representa uma entrada individual no relatório FIP215. Cada instância deste tipo contém informações detalhadas sobre uma unidade orçamentária específica, incluindo seu código, nome, identificação da conta contábil, identificação da conta contábil de explosão, código da conta contábil, nome da conta contábil, saldo anterior, valor de crédito e valor de débito.

A utilização das tags `json` em cada campo especifica como esses campos serão serializados quando convertidos para JSON, garantindo uma representação adequada quando os dados são enviados como resposta a uma solicitação HTTP.

#### `RelatorioFIP215`

O tipo `RelatorioFIP215` representa o relatório FIP215 como um todo. Ele contém um campo chamado `Dados`, que é uma lista (ou slice) de instâncias do tipo `DadoRelatorioFIP215`. Essa estrutura permite organizar e representar todas as entradas do relatório como uma coleção.

A tag `@name RelatorioFIP215` é uma anotação que pode ser usada por ferramentas de documentação (como o [swag](https://github.com/swaggo/swag)) para identificar o nome da estrutura no contexto da documentação da API.

#### Utilização

Esses modelos são úteis porque proporcionam uma maneira estruturada de organizar os dados do relatório FIP215. Durante a execução do controlador que manipula as solicitações HTTP, os dados consultados no banco de dados são mapeados para instâncias do tipo `DadoRelatorioFIP215`, que são então agrupadas em uma instância do tipo `RelatorioFIP215`.

Logo, para criar um modelo, deve ser observada a *query* SQL que será executada, e quais campos ela retorna.

### Controlador

O controlador será responsável por lidar com a lógica que deve ser executada quando o usuário consultar aquele *endpoint* específico. Como sua implementação é diferente para cada ação do usuário, ele pode ser um pouco mais complexo de se escrever. Mas, aqui destacaremos alguns pontos cruciais no desenvolvimento de um controlador, tomando como exemplo o arquivo [`relatorio_fip215_handler.go`](./internal/server/relatorio_fip215_handler.go) (Aliás, note que todos os controladores ficam em um arquivo dedicado para eles, terminado em `_handler.go`, na pasta `internal/server`):

#### Estrutura

O controlador segue um padrão de estrutura geral. A função `FIP215Handler` é responsável por lidar com as solicitações HTTP direcionadas ao endpoint associado ao relatório FIP215. A documentação do endpoint é fornecida no formato de anotações GoDoc para facilitar a compreensão e utilização, além de permitir a geração automática de documentação via [swag](https://github.com/swaggo/swag).

#### Parâmetros

Os parâmetros do *endpoint* são capturados da solicitação HTTP usando o pacote [Echo](https://echo.labstack.com/). Eles são validados para garantir que sejam fornecidos corretamente e estejam dentro dos limites esperados. Os parâmetros incluem informações como ano de exercício, unidade gestora, mês de referência, tipo de encerramento e vários outros, proporcionando flexibilidade na consulta do relatório.

#### Validação de Parâmetros

A seção de validação de parâmetros garante que as entradas fornecidas pelo usuário sejam válidas. Se houver algum problema na validação, o controlador retorna um erro HTTP apropriado com uma mensagem descritiva.

#### Consulta ao Banco de Dados

Após validar os parâmetros, o controlador realiza uma consulta ao banco de dados utilizando uma consulta SQL previamente definida. O resultado da consulta é processado e mapeado para uma estrutura de dados apropriada (`model.RelatorioFIP215`). Esses dados são então retornados como uma resposta JSON para o cliente.

#### Tratamento de Erros

O código inclui tratamento adequado para possíveis erros durante a leitura do SQL, execução da consulta no banco de dados e processamento dos resultados. Erros específicos são registrados nos logs, e respostas HTTP apropriadas são enviadas em caso de falha.

### Query SQL

A *query* SQL fica em um arquivo `.sql` na pasta `ìnternal/database/queries`. Ela contém a solicitação que será enviada ao banco de dados, e se necessário, deve conter campos como `:1` e `:2`, que são utilizados pelo controlador para passar parâmetros personalizados pelas funções `db.Query` ou `db.QueryRow`.
