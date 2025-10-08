# 1. Gefuehlmesser - Medidor de Sentimentos

![Imagem Docker](https://github.com/lwglg/gefuehlmesser/actions/workflows/image-analysis.yml/badge.svg)
![Código-fonte](https://github.com/lwglg/gefuehlmesser/actions/workflows/static-analysis.yml/badge.svg)

Plataforma de software construída para realizar análises de sentimentos, em tempo real, processando mensagens de feeds para calcular métricas de engajamento usando algorítmos determinísticos.
Avaliação técnica para a posição de desenvolvedor de software na [MBRAS Empreeendimentos](https://github.com/MBRAS-Emprendimentos).

## 1.1. TOC

- [1. Gefuehlmesser - Medidor de Sentimentos](#1-gefuehlmesser---medidor-de-sentimentos)
  - [1.1. TOC](#11-toc)
  - [1.2. Preliminares](#12-preliminares)
  - [1.3. Topologia do projeto (WiP)](#13-topologia-do-projeto-wip)
  - [1.4. Cálculo de sentimentos](#14-cálculo-de-sentimentos)
  - [1.5. Manipulação do projeto](#15-manipulação-do-projeto)
    - [1.5.1. Startup](#151-startup)
    - [1.5.2. Shutdown](#152-shutdown)
  - [1.6. TL;DR](#16-tldr)
  - [1.7. O que deseja fazer?](#17-o-que-deseja-fazer)


## 1.2. Preliminares

Esse projeto foi desenvolvido utilizando as seguintes ferramentas:

- [Go (`go`)](https://tip.golang.org/doc/go1.25): v1.25.1 ou superior
- [Docker Engine](https://docs.docker.com/engine/install/ubuntu/): v28.0.1 ou superior:
    - [Comando `docker` sem ser `sudo`](https://docs.docker.com/engine/install/linux-postinstall/). Opcional;
- [Docker Compose](https://docs.docker.com/compose/install/linux/): v2.29.7-desktop.1 ou superior;
- [GNU Make](https://www.gnu.org/software/make/): v4.3 ou superior.


## 1.3. Topologia do projeto (WiP)


Em linhas gerais, o projeto é composto de:

- `webservice`: Um servidor de aplicação HTTP, implementado em Go, utilizando uma estrutura idiomática baseada em um padrào de projeto orientado à recursos de API (_resource-oriented design_). No tocante a essa camada, a aplicação serve um API RESTful com os seguintes recursos:
  - API utilitária (`/utilitary`):
    - Requisições de _healthchecking_, de modo a verificar a sanidade do servidor de aplicação;
  - API de analise de sentimentos (`/api/v1/sentiment`):
    - Requisições para determinar a análise de sentimento para:
      - Um feed com um conjunto de mensagens;
      - Uma única mensagem; 
- `infra`: Contem as declarações do [manifesto](./infra/docker/compose.yml) do Docker Compose da [imagem](./infra/docker/webservice/Dockerfile) Docker do serviço `webservice`, além de arquivos subsidiários ao funcionamento do servidor;
- `Makefile`: Comandos executados via `make`, implementados para facilitar a manipulação da camada de automação do projeto;
- `scripts`: Scripts que são consumidos pelos comandos executados via `make`;
- `.github` Uma camada de infraestrutura como código (IaC). Na pasta `workflows`, foram implementadas duas pipelines, executadas automaticamente pelo Github Actions:
  - `static-analysis.yml`: Realiza a execução de testes unitários, verificação de vulnerabilidades e linting do código-fonte;
  - `image-analysis.yml`: Realiza a análise de sanidade da imagem Docker do serviço `webservice`, de modo a garantir o uso eficiente dos recursos do contêiner, em tempo de execução.


![Topologia](./resources/docs/images/docker-topology.png)

> [!NOTE]
> Para gerar ou atualizar o diagrama acima, basta executar na raíz do projeto o comando `make topology`.

Associado ao manifesto do Docker Compose, o serviço declarado é

 **Serviço**            | **Descrição**                                                       | **Portas expostas** | **URL base externa**    | **URL base interna**              |
|-----------------------|---------------------------------------------------------------------|---------------------|-------------------------|-----------------------------------|
| `webservice`          | Servidor de aplicação (Go) que fornece a API RESTful                | 8080 (TCP,HTTP)     | `http://localhost:8080` | `http://webservice:8080`          |

> [!IMPORTANT]
> Concernindo as URLs supracitadas, as **externas** são classificadas assim pois são **externas à rede do Docker**, e.g. quando a aplicação Fastify está sendo rodada localmente na máquina.
> 
> Similarmente, uma URL é dita **interna** pois é **interna à rede do Docker**, e.g. quando tanto o app Fastify quanto o servidor PostgreSQL estão rodando em seus contêineres, se comunicando via TCP.


## 1.4. Cálculo de sentimentos

Com base nas [especificações](https://github.com/MBRAS-Emprendimentos/backend-challenge-092025/blob/main/docs/algorithm_examples.md) do teste técnico, o modelo determinístico para o cálculo de sentimentos de mensagens de feed foi implementado como uma biblioteca, em `./webservices/libs/sentiment`, seguindo a estrutura abaixo:

```bash
./webservice/libs
├── ... # Demais bibliotecas
├── sentiment
│   ├── analyzer.go         # Onde se encontram os métodos principais para cálculo de sentimentos para feed e para uma mensagem isolada
│   ├── anomalies.go        # Função auxiliar para detecção de anomalias (várias mensagens dentro de uma janela de 2s), com base em todas as mensagens do feed
│   ├── definitions.go      # Declarações de constantes comuns a todos os módulos da lib
│   ├── engagement.go       # Funções auxiliares para o cálculo do escore de engajamento e de influência que as mensagens causaram, para todos os usuários no feed
│   ├── helpers.go          # Funções auxiliares comuns aos demais módulos da lib, e.g. remoção de acentos e tokenização de mensagens, validaçòes de data/hora
│   ├── models.go           # Declarações de tipos de I/O mais complexos, utilizados nas assinaturas de funções
│   └── trending_topics.go  # Função auxiliar para a extração das "trending topics" (tópicos com engajamento mais notável) dentre as mensagens do feed
```

> [!IMPORTANT]
> Por questões de restrição de tempo, não foram conduzidos testes exploratórios mais extensos da API e, por extensão, da biblioteca em questão. Entretanto, ambos pode ser considerados funcionais.

## 1.5. Manipulação do projeto

Através do `make`, via scripts de automação do Docker Compose implementados em um [Makefile](./Makefile), na raíz do projeto. Para conferir a documentação de cada script, basta executar no terminal

```bash
make                                # Sem nenhum comando, executa o fallback 'help'
make help                           # Explicitamente, mostra a documentação
```

### 1.5.1. Startup

Considerando uma instalação inicial, na raíz do projeto, execute os seguintes comandos:

```bash
$ make swagger                   # Gera ou atualiza a pasta `./webservice/docs`, a qual contém 
$ make build c=webservice        # Realiza a build da imagem Docker, cujo manifesto encontra em `./infra/docker/webservice`
$ make init c=webservice         # Inicia o container do serviço `webservice`, em modo detached, e inicia a captura de logs
```

> [!NOTE]
> Para se certificar de que as imagens foram geradas pelo processo de build, basta executar o comando `docker image ls`.
> 
> Para se certificar de que os contêineres foram de fato devidamente iniciados e na escuta das portas corretas, basta executar o comando `make ps`.
> 
> Para ver os logs de um conteiner específico, execute `make logs c=[nome-do-serviço]`.


> [!IMPORTANT]
> A documentação da API RESTful pode ser acessada pela Swagger UI, via browser, através da URL [`http://localhost:8080/api/v1/swagger`](http://localhost:8080/api/v1/swagger/index.html#/)
> 
> De modo a conduzir testes exploratórios na API RESTful, foram criados uma _collection_ e um _environment_ do Postman. Ambos se encontram em [`./resources/testing`](./resources/testing/)
 

### 1.5.2. Shutdown

Similarmente, para ambos os ambientes, de modo a encerrar a execução de todos os contêineres, basta rodar:

```bash
make stop          # Interrompe os contêineres que estiverem sendo executados
make clean         # Opcional. Remove os contêineres e a network associadas aos serviços do ambiente
```

## 1.6. TL;DR

Entre em contato com o desenvolvedor do projeto:

- Guilherme Lima Gonçalves
  - [Github](https://github.com/lwglg) 
  - [E-mail](mailto:lwglguilherme@gmail.com)

---

## 1.7. O que deseja fazer?

- [Voltar ao topo](#11-toc)
