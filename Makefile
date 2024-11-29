########################################
# Commands
########################################

# Services
# make run      - Run service
# make db       - Run docker database
# make clean-db - Clean DB volume

########################################

run:
	go run main.go

db:
	cd docker && docker-compose -f glower-db.yaml up

clean-db:
	cd docker && docker-compose -f glower-db.yaml down -v