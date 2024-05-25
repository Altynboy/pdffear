FROM golang:alpine

RUN apk update && apk add --no-cache libreoffice
RUN apk add --no-cache msttcorefonts-installer fontconfig
RUN update-ms-fonts


RUN apk --no-cache add openjdk11
ENV JAVA_HOME /usr/lib/jvm/default-jvm


# COPY go.mod ./
# COPY go.mod go.sum ./
# RUN go mod download

RUN mkdir /tmp/generated_pdfs
RUN mkdir /tmp/uploaded_docx
RUN mkdir /tmp/libreoffice_profiles

WORKDIR /pdffear

COPY . .

ENV TMP_PDF_PATH="/tmp/generated_pdfs/"
ENV TMP_DOCX_PATH="/tmp/uploaded_docx/"
ENV LIBREOFFICE_PROFILES="/tmp/libreoffice_profiles/"

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -o /docker-gs-ping

EXPOSE 8080

# Run
CMD ["/docker-gs-ping"]
