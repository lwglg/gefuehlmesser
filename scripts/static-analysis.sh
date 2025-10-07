#!/bin/bash


function check_vulnerabilities() {
    echo "Checking for vulnerabiliites in the source code..."

    ./bin/govulncheck ./...
}


function suspicious_constructs() {
    echo "Checking for suspicious constructs..."

    go vet ./...
}


function unit_tests() {
    echo "Running unit tests..."

    go test ./...
}


function formatting() {
    echo "Running source code formatting..."

    ./bin/gofumpt -l -w . 2>&1 | tee outfile && test -z "$(cat outfile)" && rm outfile
}


function all() {
    echo "Running all static analysis steps..."

    check_vulnerabilities
    suspicious_constructs
    formatting
    unit_tests
}


function main() {
    MODE=$1
    SUPPORTED_MODES="vulnerability | construct | test | format | all"

    cd ./webservice

    case $MODE in
        "vulnerability")
            check_vulnerabilities;;
        "construct")
            suspicious_constructs;;
        "test")
            unit_tests;;
        "format")
            formatting;;
        "all")
            all;;
        *)
            echo "Modo de execução não suportado. Valores suportados: ${SUPPORTED_MODES}"
            exit 1;;
    esac

    echo "Done!"
    cd ..
    exit 0
}

main $@
