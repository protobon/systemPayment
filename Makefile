.PHONY: stage

TEST := test
PROD := prod

COMPOSE-TEST=sudo docker compose -f docker-compose-test.yml
COMPOSE-PROD=sudo docker compose -f docker-compose-prod.yml
LOGS-TEST=sudo docker logs -f system_payment_test
LOGS-PROD=sudo docker logs -f system_payment

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


# ------- Up ----------------------------------------------------
up:
	@echo $(stage)
ifeq ($(stage), $(TEST))
	$(COMPOSE-TEST) up --remove-orphans
endif
ifeq ($(stage), $(PROD))
	$(COMPOSE-PROD) up --remove-orphans
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

# ------- Start ----------------------------------------------------
start:
	@echo $(stage)
ifeq ($(stage), $(TEST))
	$(COMPOSE-TEST) start
endif
ifeq ($(stage), $(PROD))
	$(COMPOSE-PROD) start
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

# ------- Logs ----------------------------------------------------
logs:
	@echo $(stage)
ifeq ($(stage), $(TEST))
	$(LOGS-TEST)
endif
ifeq ($(stage), $(PROD))
	$(LOGS-PROD)
endif