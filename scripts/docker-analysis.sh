#!/bin/bash

set -e

DEFAULT_DIVE_VERSION=latest
DEFAULT_DIVE_IMAGE=docker.io/wagoodman/dive
DEFAULT_DIVE_UI_CONFIG_PATH="$(pwd)/.dive-ui.yaml"
DEFAULT_DIVE_CI_CONFIG_PATH="$(pwd)/.dive-ci.yaml"
DEFAULT_DOCKER_SOCK_PATH=/var/run/docker.sock

function usage() {
    MODE=$1

    case $MODE in
        "ui")
            CONFIG_PATH=$DEFAULT_DIVE_UI_CONFIG_PATH
            WHAT_IT_DOES="Performs the Docker image analysis, showing the console user interface.";;
        "ci")
            CONFIG_PATH=$DEFAULT_DIVE_CI_CONFIG_PATH
            WHAT_IT_DOES="Performs the Docker image analysis, without the user interface, i.e. pass/fail mode.";;
        *)
            echo "Modo de execução não suportado. Valores suportados: ui | ci"
            exit 1;;
    esac

    echo "----------------------------------------------------------------------------------------------------------------------------------"
    echo -e "\033[1m${MODE^^}: $WHAT_IT_DOES\033[0m"
    echo "----------------------------------------------------------------------------------------------------------------------------------"
    echo "Command: ./docker-analysis.sh $MODE [target-image] [config-file-path] [dive-version]"
    echo "----------------------------------------------------------------------------------------------------------------------------------"
    echo "Where:"
    echo "- target-image (string):       The image to be analyzed, just the repository part (without the tag)"
    echo "- config-file-path (string):   The absolute path to the config file."
    echo "                               Default: \"$CONFIG_PATH\""
    echo "- dive-version (string):       The (semantic) version of Dive application."
    echo "                               Default: \"latest\""
    echo
    echo "Example: ./docker-analysis.sh $MODE docker.io/wagoodman/dive /home/foo/.dive.yaml 0.9.2"
    echo "----------------------------------------------------------------------------------------------------------------------------------"
}

function trim() {
    echo "$(echo -e "${1}" | sed -e 's/^[[:space:]]*//' -e 's/[[:space:]]*$//')"
}

function validate_input() {
    MODE=$1
    TARGET_IMAGE=$2

    if [ ! $TARGET_IMAGE ]
	then
        echo -e "ERROR: The repository for the target image must be provided!\n"
		usage $MODE
		exit 1
	fi
}

function ui() {
    TARGET_IMAGE=$(trim $1)
    CONFIG_PATH=$(trim ${2:-DEFAULT_DIVE_UI_CONFIG_PATH})
    DIVE_VERSION=$(trim ${3:-$DEFAULT_DIVE_VERSION})

    validate_input "ui" $TARGET_IMAGE

    DIVE_IMAGE="$DEFAULT_DIVE_IMAGE:$DIVE_VERSION"

    docker run --rm -it \
        -v $DEFAULT_DOCKER_SOCK_PATH:$DEFAULT_DOCKER_SOCK_PATH \
        -v  "$(pwd)":"$(pwd)" \
        -v $CONFIG_PATH:"$HOME/.dive.yaml" \
        $DIVE_IMAGE \
            $TARGET_IMAGE

    exit 0
}

function ci() {
    TARGET_IMAGE=$(trim $1)
    CONFIG_PATH=$(trim ${2:-DEFAULT_DIVE_CI_CONFIG_PATH})
    DIVE_VERSION=$(trim ${3:-$DEFAULT_DIVE_VERSION})

    validate_input "ci" $TARGET_IMAGE

    DIVE_IMAGE="$DEFAULT_DIVE_IMAGE:$DIVE_VERSION"

    docker run \
        --rm \
        -it \
        -e CI=true \
        -v $DEFAULT_DOCKER_SOCK_PATH:$DEFAULT_DOCKER_SOCK_PATH \
        $DIVE_IMAGE \
            --ci-config $CONFIG_PATH \
            $TARGET_IMAGE

    exit 0
}

$@
