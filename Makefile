.PHONY: stage

TEST := test
PROD := prod
LOCAL := local

# ------- Test (server) ------- 
COMPOSE-TEST=sudo docker compose -f docker-compose.yml
LOGS-TEST=sudo docker logs -f system_payment_test

# ------- Prod (server) ------- 
COMPOSE-PROD=docker compose -f docker-compose-prod.yml
LOGS-PROD=docker logs -f system_payment

# ------- LOCAL (postgres instance + run app locally) -------
COMPOSE-DB-LOCAL=sudo docker run --name system_payment_db_local -v systempayment_data_test:/var/lib/postgresql/data \
					-p 5432:5432 -e POSTGRES_USER=spuser -e POSTGRES_PASSWORD=SPuser96 -e POSTGRES_DB=system_payment_test -d postgres
START-LOCAL=sudo docker start system_payment_db_local
RESTART-LOCAL=sudo docker restart system_payment_db_local
STOP-LOCAL=sudo docker stop system_payment_db_local

# ------- Local Compose -------
COMPOSE=sudo docker compose -f docker-compose.yml
LOGS=sudo docker logs -f system_payment_test

# ------- pull from current branch -------
pull:
	git pull

# ------- levanta la aplicacion en maquina local --------------------
run:
	DLOCAL_URL=https://sandbox.dlocal.com DLOCAL_X_LOGIN=id9LdRTxgd DLOCAL_X_TRANS_KEY=MLYS0sI3qt DLOCAL_SECRET=b1pHjMu99d6Y7YLmgok9KEfQDt1d4KCuI POSTGRES_USER=spuser POSTGRES_PASSWORD=SPuser96 POSTGRES_DB=system_payment_test APPLICATION_PORT=:8080 DATABASE_HOST=localhost:5432 go run main.go

# ------- Build ----------------------------------------------------
build:
	@echo $(stage)
ifeq ($(stage), $(TEST))
	$(COMPOSE-TEST) build
endif
ifeq ($(stage), $(PROD))
	$(COMPOSE-PROD) build
endif
ifeq ($(stage), $(LOCAL))
	$(COMPOSE-DB-LOCAL)
endif
ifeq ($(stage),)
	$(COMPOSE) build
endif


# ------- Up ----------------------------------------------------
up:
	@echo $(stage)
ifeq ($(stage), $(TEST))
	$(COMPOSE-TEST) up --remove-orphans
endif
ifeq ($(stage), $(PROD))
	$(COMPOSE-PROD) up --remove-orphans
endif
ifeq ($(stage),)
	$(COMPOSE) up --remove-orphans
endif


# ------- Up detached ----------------------------------------------------
dup:
	@echo $(stage)
ifeq ($(stage), $(TEST))
	$(COMPOSE-TEST) up -d --remove-orphans
endif
ifeq ($(stage), $(PROD))
	$(COMPOSE-PROD) up -d --remove-orphans
endif
ifeq ($(stage),)
	$(COMPOSE) up -d --remove-orphans
endif


# ------- Start ----------------------------------------------------
start:
	@echo $(stage)
ifeq ($(stage), $(TEST))
	$(COMPOSE-TEST) start
endif
ifeq ($(stage), $(PROD))
	$(COMPOSE-PROD) start
endif
ifeq ($(stage), $(LOCAL))
	$(START-LOCAL)
endif
ifeq ($(stage),)
	$(COMPOSE) start
endif


# ------- Stop ----------------------------------------------------
stop:
	@echo $(stage)
ifeq ($(stage), $(TEST))
	$(COMPOSE-TEST) stop
endif
ifeq ($(stage), $(PROD))
	$(COMPOSE-PROD) stop
endif
ifeq ($(stage), $(LOCAL))
	$(STOP-LOCAL)
endif
ifeq ($(stage),)
	$(COMPOSE) stop
endif


# ------- Restart ----------------------------------------------------
restart:
	@echo $(stage)
ifeq ($(stage), $(TEST))
	$(COMPOSE-TEST) restart
endif
ifeq ($(stage), $(PROD))
	$(COMPOSE-PROD) restart
endif
ifeq ($(stage), $(LOCAL))
	$(RESTART-LOCAL)
endif
ifeq ($(stage),)
	$(COMPOSE) restart
endif


# ------- Down ----------------------------------------------------
down:
	@echo $(stage)
ifeq ($(stage), $(TEST))
	$(COMPOSE-TEST) down
endif
ifeq ($(stage), $(PROD))
	$(COMPOSE-PROD) down
endif
ifeq ($(stage),)
	$(COMPOSE) down
endif


# ------- Logs ----------------------------------------------------
logs:
	@echo $(stage)
ifeq ($(stage), $(TEST))
	$(LOGS-TEST)
endif
ifeq ($(stage), $(PROD))
	$(LOGS-PROD)
endif
ifeq ($(stage), $(LOCAL))
	$(LOGS-LOCAL)
endif
ifeq ($(stage),)
	$(LOGS)
endif