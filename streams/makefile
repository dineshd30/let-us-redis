


# clean docker containers, images, and volumes
docker-clean:
	-@docker stop $$(docker container ls -qa)
	-@docker container rm $$(docker container ls -aq) -f
	-@docker volume rm $$(docker volume ls -q)
	-@docker image rmi $$(docker image ls -aq)
	-@docker system prune -af