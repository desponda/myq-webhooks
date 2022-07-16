FROM golang:1.18 as build
WORKDIR /build
ADD go.mod .
RUN go mod download
COPY . /build
RUN go build -o /main cmd/main.go
# copy artifacts to a clean image
FROM public.ecr.aws/lambda/provided:al2
COPY --from=build /main /main
ENTRYPOINT [ "/main" ]     