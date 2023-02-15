.PHONY: stage

TEST := test
PROD := prod

COMPOSE=sudo docker compose -f docker-compose.yml
COMPOSE-TEST=docker compose -f docker-compose-test.yml
COMPOSE-PROD=docker compose -f docker-compose-prod.yml
LOGS=sudo docker logs -f system_payment_test
LOGS-TEST=docker logs -f system_payment_test
LOGS-PROD=docker logs -f system_payment

# ------- pull de codigo sobre la rama actual ------------------------
pull:
	git pull

# ------- Build ----------------------------------------------------
build:
	@echo $(stage)
ifeq ($(stage), $(TEST))
	$(COMPOSE-TEST) build
endif
ifeq ($(stage), $(PROD))
	$(COMPOSE-PROD) build
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
ifeq ($(stage),)
	$(COMPOSE) restart
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
ifeq ($(stage),)
	$(LOGS)
endif