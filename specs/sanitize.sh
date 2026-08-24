#!/bin/bash
# sanitize.sh — transforms openapi_current.yaml → openapi_corrected.yaml
#
# The raw IBM Verify OpenAPI spec uses several inconsistent error schema names
# across its 385 paths. This script normalises all of them to a single
# VerifyError schema so Fern generates one consistent error type.
#
# Usage:
#   bash specs/sanitize.sh
#
# Run this whenever IBM releases a new version of their OpenAPI spec:
#   1. Replace specs/openapi_current.yaml with the new IBM spec
#   2. Run this script
#   3. Review the diff on specs/openapi_corrected.yaml
#   4. Re-run: make generate

set -e

INPUT="$(dirname "$0")/openapi_current.yaml"
OUTPUT="$(dirname "$0")/openapi_corrected.yaml"

cp "$INPUT" "$OUTPUT"

# Normalise all IBM Verify error schema refs to a single VerifyError type.
# The raw spec uses at least five different names for the same concept.
sed -i '' 's|\$ref": "#/components/schemas/BadRequest"|\$ref": "#/components/schemas/VerifyError"|g' "$OUTPUT"
sed -i '' 's|\$ref": "#/components/schemas/Forbidden"|\$ref": "#/components/schemas/VerifyError"|g' "$OUTPUT"
sed -i '' 's|\$ref": "#/components/schemas/NotFound"|\$ref": "#/components/schemas/VerifyError"|g' "$OUTPUT"
sed -i '' 's|\$ref": "#/components/schemas/ErrorResponse"|\$ref": "#/components/schemas/VerifyError"|g' "$OUTPUT"
sed -i '' 's|\$ref": "#/components/schemas/ErrorResponseBean"|\$ref": "#/components/schemas/VerifyError"|g' "$OUTPUT"

echo "✓ Sanitized $(basename $INPUT) → $(basename $OUTPUT)"
echo "  Review changes: git diff specs/openapi_corrected.yaml"
