# PDFFEAR

Docx to PDF Converter with HTTP Server in Docker using LibreOffice

<p align="center">
    <img src="https://img.shields.io/badge/golang-v1.22-lightblue" height="25"/>
    <img src="https://img.shields.io/badge/libreoffice-v7.6-00a500" height="25"/>
    <img src="https://img.shields.io/badge/dockerfile-support-green" height="25"/>
</p>

## Features

- Fast
- Simple
- Extendable

## How to use

1. **Clone the repository**:

```bash
git clone https://github.com/Altynboy/pdffear.git
cd pdffear
```

2. **Build image**:

```bash
docker run build -t pdffear .
```

3. **Run the container**:

```bash
docker run --rm -d -p 8080:8080/tcp pdffear:latest
```

4. **Send docx**:

```bash
curl --location 'localhost:8080/upload' \
--form 'myFile=@"/C:/your-folder/test.docx"'
```
