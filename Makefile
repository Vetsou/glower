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

test-l1:
	go test -tags=L1 ./tests/...

test-l2:
	go test -tags=L2 ./tests/...

db:
	cd build && docker-compose -f glower-db.yaml up

clean-db:
	cd build && docker-compose -f glower-db.yaml down -v