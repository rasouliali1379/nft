FROM p-repo.780.ir/alpine:latest

WORKDIR /server

COPY application .

EXPOSE 8080
EXPOSE 8082

CMD /server/application migrate init ; /server/application migrate up ; /server/application