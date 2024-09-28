#!/bin/bash

# Variables
COUCHBASE_HOST="couchbase-server"
ADMIN_USERNAME=${COUCHBASE_ADMINISTRATOR_USERNAME}
ADMIN_PASSWORD=${COUCHBASE_ADMINISTRATOR_PASSWORD}
BUCKET_NAME="user"

# Function to execute Couchbase CLI commands
execute_cb_cli() {
  couchbase-cli "$@"
}

# Wait until Couchbase server is up
echo "Waiting for Couchbase Server to be ready..."
until curl -s http://${COUCHBASE_HOST}:8091/pools > /dev/null; do
  echo -n "."
  sleep 1
done
echo "Couchbase Server is up!"

# Create the 'user' bucket if it doesn't exist
echo "Creating bucket '$BUCKET_NAME'..."
execute_cb_cli bucket-create -c ${COUCHBASE_HOST}:8091 \
  --username=${ADMIN_USERNAME} \
  --password=${ADMIN_PASSWORD} \
  --bucket=${BUCKET_NAME} \
  --bucket-type=couchbase \
  --bucket-ramsize=100 \
  --bucket-replica=1

# Wait for the bucket to be ready
echo "Waiting for bucket '$BUCKET_NAME' to be ready..."
until curl -s http://${COUCHBASE_HOST}:8091/pools/default/buckets/${BUCKET_NAME} > /dev/null; do
  echo -n "."
  sleep 1
done
echo "Bucket '$BUCKET_NAME' is ready!"

# Create a primary index for the 'user' bucket
echo "Creating primary index on '$BUCKET_NAME'..."
curl -s -X POST \
  -u ${ADMIN_USERNAME}:${ADMIN_PASSWORD} \
  http://${COUCHBASE_HOST}:8093/query/service \
  -d "statement=CREATE PRIMARY INDEX \`#primary\` ON \`${BUCKET_NAME}\`"

# Insert sample data into the 'user' bucket
echo "Inserting sample data into '$BUCKET_NAME'..."
curl -s -X POST \
  -u ${ADMIN_USERNAME}:${ADMIN_PASSWORD} \
  http://${COUCHBASE_HOST}:8093/query/service \
  -d "statement=INSERT INTO \`${BUCKET_NAME}\` (KEY, VALUE) VALUES ('user::1', {\"id\": \"1\", \"name\": \"John Doe\", \"email\": \"john.doe@example.com\"})"

curl -s -X POST \
  -u ${ADMIN_USERNAME}:${ADMIN_PASSWORD} \
  http://${COUCHBASE_HOST}:8093/query/service \
  -d "statement=INSERT INTO \`${BUCKET_NAME}\` (KEY, VALUE) VALUES ('user::2', {\"id\": \"2\", \"name\": \"Jane Smith\", \"email\": \"jane.smith@example.com\"})"

echo "Sample data inserted into '$BUCKET_NAME' bucket."
