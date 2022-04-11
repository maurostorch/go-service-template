FROM alpine
RUN apk add ca-certificates
COPY app /bin/app

CMD [ "/bin/app" ]
