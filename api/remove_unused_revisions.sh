#!/bin/sh

# This is a utility script to remove all unused revisions.
# It includes:
#  - failed revisions
#  - revisions wihtout a tag and no traffic

count=0

for rev in $(gcloud run revisions list \
	--format=json \
	--project=$GCP_PROJECT \
	--region=$GCP_REGION \
	--flatten="status.conditions" \
	--filter="status.conditions.type=Active AND status.conditions.status=False" \
	| jq '.[] | @base64');
do
  _jq() {
		e=$(echo ${rev} | sed 's/"//g' | base64 --decode | jq -r ${1})
		echo "${e}"
	}

	name=$(_jq '.metadata.name');
	echo "Found inactive revision: ${name}";

	gcloud run revisions delete \
		$name \
		--quiet \
		--project=$GCP_PROJECT \
		--region=$GCP_REGION \
		--async;

	((count++));
done

echo "Deleted ${count} revisions"
