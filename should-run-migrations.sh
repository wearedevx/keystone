#!/bin/bash
if [[ -z "$MIGRATE" ]]; then
	MIGRATE=migrate		 
fi

MIGRATION_DIR=api/db/migrations;
echo "Migration directory: ${MIGRATION_DIR}";
N=();

for file in $MIGRATION_DIR/*; do
	f=$(echo $file | sed "s#${MIGRATION_DIR}/##; s/_.*//; s/^00000//; s/^0000//; s/^000//; s/^00//; s/^0//");
	N+=($f);
done;

IFS=$'\n';
latest_migration_on_disk=$(echo "${N[*]}" | sort -rn | head -n1);
echo "The latest available migration on disk is ${latest_migration_on_disk}";
on_disk=$(($latest_migration_on_disk+0));

latest_migration_ran=$($MIGRATE -database "$DATABASE_URL" -path "$MIGRATION_DIR" version 2>&1);
echo "The latest ran migration in db is ${latest_migration_ran}";
ran=$(($latest_migration_ran+0));

if (( on_disk > ran )); then
	echo "Migrations should be run";
	# Success case, it’s the truthy case in a if statement
	exit 0;
else
	echo "No need for migrations";
	# Error case, it’s the falsy case in a if statement
	exit 1;
fi;
