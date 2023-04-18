#!/bin/sh

set -e

echo "******************************* Listing Env Variables..."
printenv
echo "******************************* Starting single cockroach node..."

./cockroach start-single-node --insecure --log-config-file=files/logs.yaml --background

# echo "******************************* Creating user"
# cockroach user set ${COCKROACH_USER} --password 1234 --echo-sql
# cockroach user ls

echo "******************************* Init database"

cat > init.sql <<EOF
-- Create Database
GRANT CONNECT ON DATABASE ${COCKROACH_DATABASE} TO ${COCKROACH_USER};
GRANT USAGE ON SCHEMA public TO ${COCKROACH_USER};
GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO ${COCKROACH_USER};
SET CLUSTER SETTING sql.trace.log_statement_execute = true;
EOF

cat init.sql

echo "*******************************  |=> Applying init.sql"

./cockroach sql --insecure --file init.sql

echo "******************************* To the moon"

# tail logs to make them accesible with docker logs
tail -f cockroach-data/logs/cockroach.log