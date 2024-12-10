all: build run

name = forum

build:
	@sudo docker build -t $(name) .

create:
	@sudo docker create --name $(name) -p 8080:8080 $(name) ./main &> /dev/null && \
	echo "'$(name)' container created successfully" || \
	echo "'$(name)' container already exists" \

start:
	@sudo docker container start -a $(name)

run: create start

stop:
	@sudo docker container stop $(name) &> /dev/null

delete:
	@sudo docker container rm $(name) &> /dev/null && \
	echo "$(name) container deleted correctly" || \
	echo "No container to remove"
	@sudo docker rmi $(name) &> /dev/null && \
	echo "$(name) image deleted correctly" || \
	echo "No image to remove"

clean: stop delete