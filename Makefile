# Diretório-raíz
ROOT_DIR := $(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))

# Cores
GREEN   := $(shell tput -Txterm setaf 2)
WHITE   := $(shell tput -Txterm setaf 7)
YELLOW  := $(shell tput -Txterm setaf 3)
RED     := $(shell tput -Txterm setaf 1)
RESET   := $(shell tput -Txterm sgr0)

# Comandos
DOCKER_COMPOSE_FILE := ./infra/docker/compose.yml
DOCKER_COMPOSE_CMD  := docker compose -f $(DOCKER_COMPOSE_FILE)

# Constrói a documentação de cada script, visualizável via 'make' ou 'make help'
# A documentação dee cada script é feita através de uma string começando por '\#\#'
# Uma categoria de comandos pode ser adicionada em uma string iniciando a mesma com @category
HELP_FUN = \
    %help; \
    while(<>) { \
		push @{$$help{$$2 // 'Comandos'}}, [$$1, $$3] if /^([a-zA-Z\-]+)\s*:.*\#\#(?:@([a-zA-Z\-]+))?\s(.*)$$/ \
	}; \
	print "Utilização: make [comando]\n\n"; \
	for (sort keys %help) { \
		print "${WHITE}$$_:${RESET}\n"; \
		for (@{$$help{$$_}}) { \
				$$sep = " " x (32 - length $$_->[0]); \
				print "  ${YELLOW}$$_->[0]${RESET}$$sep${GREEN}$$_->[1]${RESET}\n"; \
		}; \
		print "\n"; \
	}

.PHONY: help \
		header \
		format \
		lint \
		test \
		lservices \
		build \
		confirm \
		clean \
		destroy \
		logs \
		restart \
		start \
		init \
		status \
		stop \
		up \
		run \
		exec \
		ps \
		imganalysisci \
		imganalysisui \
		topology

.DEFAULT_GOAL := help

info: header

define HEADER
+--------------------------------------------------------------------------------------------------------------------------+
8""8""8 8""""8   8"""8  8""""8 8""""8          8""""8
8  8  8 8    8   8   8  8    8 8               8    " eeee eeee e   e eeee e   e e     eeeeeee eeee eeeee eeeee eeee eeeee
8e 8  8 8eeee8ee 8eee8e 8eeee8 8eeeee          8e     8    8    8   8 8    8   8 8     8  8  8 8    8   " 8   " 8    8   8
88 8  8 88     8 88   8 88   8     88   eeee   88  ee 8eee 8eee 8e  8 8eee 8eee8 8e    8e 8  8 8eee 8eeee 8eeee 8eee 8eee8e
88 8  8 88     8 88   8 88   8 e   88          88   8 88   88   88  8 88   88  8 88    88 8  8 88      88    88 88   88   8
88 8  8 88eeeee8 88   8 88   8 8eee88          88eee8 88ee 88   88ee8 88ee 88  8 88eee 88 8  8 88ee 8ee88 8ee88 88ee 88   8
+--------------------------------------------------------------------------------------------------------------------------+
endef
export HEADER

header: ##@Outros Mostra o header deste help, formado com caracteres ASCII.
	clear
	@echo "$$HEADER"

help: ##@Outros Mostra esta documentação.
	clear
	@echo "$$HEADER"
	@perl -e '$(HELP_FUN)' $(MAKEFILE_LIST)

lservices: ## Lista todos os nomes de serviços declarados no YAML do Docker Compose, dado um env=<dev | prod> ambiente de infra
	$(DOCKER_COMPOSE_CMD) config --services

build: ## Realiza a build de todas as imagens Docker, ou para um c=<node de serviço> específico, dado um env=<dev | prod> ambiente de infra
	$(DOCKER_COMPOSE_CMD) build $(c)

confirm:
	@( read -p "$(RED)Tem certeza? [y/N]$(RESET): " sure && case "$$sure" in [sSyY]) true;; *) false;; esac )

clean: confirm ## Realiza a limpeza de todos os dados associados aos conteineres, dado um env=<dev | prod> ambiente de infra
	$(DOCKER_COMPOSE_CMD) down)

destroy: confirm ## Remove todas as imagens, volumes, networks e conteineres não utilizados. Use com cautela!
	@docker system prune --all --volumes --force
	@docker volume prune --all --force
	@docker network prune --force
	@docker image prune --all --force

logs: ## Adiciona captura de logs para todos os conteineres ou para um c=<nome de serviço>, dado um env=<dev | prod> ambiente de infra
	$(DOCKER_COMPOSE_CMD) logs --follow $(c)

restart: ## Reinicia todos os conteineres ou apenas um c=<nome de serviço>, dado um env=<dev | prod> ambiente de infra
	$(DOCKER_COMPOSE_CMD) stop $(c)
	@make init c=$(c)

start: ## Inicia todos os conteineres em background (detached mode) ou apenas um c=<nome de serviço>, dado um env=<dev | prod> ambiente de infra
	$(DOCKER_COMPOSE_CMD) up -d $(c)

init: ## Inicia um conteiner em detached mode, com captura de logs, dado um env=<dev | prod> ambiente de infra
	@make start env=$(env) c=$(c) && make logs env=$(env) c=$(c)

status: ## Lista os status dos conteineres em execução, dado um env=<dev | prod> ambiente de infra
	$(DOCKER_COMPOSE_CMD) ps

stop: ## Encerra a execução de todos os conteineres ou de apenas um c=<nome de serviço>, dado um env=<dev | prod> ambiente de infra
	$(DOCKER_COMPOSE_CMD) stop $(c)

up: ## Inicia todos os conteineres em modo "attached" ou apenas um c=<nome de serviço>, dado um env=<dev | prod> ambiente de infra
	$(DOCKER_COMPOSE_CMD) up $(c)

run: ## Roda um comando (o que seria especificado em 'CMD' na imagem), dado um c=<nome de serviço> e um env=<dev | prod> ambiente de infra
	$(DOCKER_COMPOSE_CMD) run --rm $(c) $(cmd)

exec: ## Executa um comando em um container já iniciado, dado um c=<nome de serviço> e um s=<script> e um env=<dev | prod> ambiente de infra
	$(DOCKER_COMPOSE_CMD) exec -it $(c) $(s)

ps: status ## Alias do comando 'status'

imganalysisui: ## Executa a análise de uma imagem Docker, em modo UI, dado uma img=<imagem Docker>
	@./scripts/docker-analysis.sh ui $(img)

imganalysisci: ## Executa a análise de uma imagem Docker, em modo CI, dado uma img=<imagem Docker>
	@./scripts/docker-analysis.sh ci $(img)

topology: ## Gera um diagrama dos serviços listados no arquivo YML do Docker Compose
	@./scripts/generate-topology.sh topology $(env)

swagger: ## Gera ou atualiza package Go associada à geração da documentação da API (Swagger UI)
	@./webservice/bin/swag init -d ./webservice/cmd/api,./webservice -o ./webservice/docs
