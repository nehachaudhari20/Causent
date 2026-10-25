#!/usr/bin/env bash
# Template — customize DRUNIX_ROOT and org env before use.
set -euo pipefail

DRUNIX_ROOT="${DRUNIX_ROOT:-../drunix}"
CHAINCODE_NAME="${CHAINCODE_NAME:-causent-assurance}"
CHAINCODE_PATH="${CHAINCODE_PATH:-$(cd "$(dirname "$0")/.." && pwd)/drunix/chaincode/causent-assurance}"
CHANNEL_NAME="${CHANNEL_NAME:-causent-channel}"
SEQUENCE="${SEQUENCE:-1}"
VERSION="${VERSION:-1.0}"

echo "Packaging ${CHAINCODE_NAME} from ${CHAINCODE_PATH}"
peer lifecycle chaincode package "${CHAINCODE_NAME}.tar.gz" \
  --path "${CHAINCODE_PATH}" \
  --lang golang \
  --label "${CHAINCODE_NAME}_${VERSION}"

echo "Next steps (per org): install, approveformyorg, commit — see Drunix test-network docs."
echo "  peer lifecycle chaincode install ${CHAINCODE_NAME}.tar.gz"
echo "  peer lifecycle chaincode approveformyorg -o orderer.example.com:7050 --channelID ${CHANNEL_NAME} --name ${CHAINCODE_NAME} --version ${VERSION} --package-id <PACKAGE_ID> --sequence ${SEQUENCE}"
echo "  peer lifecycle chaincode commit -o orderer.example.com:7050 --channelID ${CHANNEL_NAME} --name ${CHAINCODE_NAME} --version ${VERSION} --sequence ${SEQUENCE}"
