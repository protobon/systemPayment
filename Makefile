.PHONY: stage

TEST := test
PROD := prod

COMPOSE-TEST=sudo docker-compose -f docker-compose-test.yml
COMPOSE-PROD=sudo docker-compose -f docker-compose-prod.yml

# ------- create network --------------------
network:
	docker network create system_payment_network

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
	$(COMPOSE-TEST) up
endif
ifeq ($(stage), $(PROD))
	$(COMPOSE-PROD) up
endif
ifeq ($(stage),)
	$(COMPOSE) up
endif

# ------- Up detached ----------------------------------------------------
dup:
	@echo $(stage)
ifeq ($(stage), $(TEST))
	$(COMPOSE-TEST) up -d
endif
ifeq ($(stage), $(PROD))
	$(COMPOSE-PROD) up -d
endif

# ------- Start ----------------------------------------------------
start:
	@echo $(stage)
ifeq ($(stage), $(TEST))
	$(COMPOSE-TEST) start
endif
ifeq ($(stage), $(DEMO))
	$(COMPOSE-DEMO) start
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

# ------- Limpieza ----------------------------------------------------
down:
	@echo $(stage)
ifeq ($(stage), $(TEST))
	$(COMPOSE-TEST) down -v --remove-orphans
endif
ifeq ($(stage), $(PROD))
	$(COMPOSE-PROD) down -v --remove-orphans
endif
