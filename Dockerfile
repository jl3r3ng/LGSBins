FROM golang
RUN mkdir -p /snippetbox
WORKDIR /snippetbox
COPY . /snippetbox
CMD ["go", "run", "./cmd/web"]