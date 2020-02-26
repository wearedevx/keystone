#!/bin/bash

# Output the current keystone environment and project

ks_status() {
	if command -v ks 2>/dev/null; then
		command_output="$(ks status | awk '{print $NF}')"
		# readarray -d ' ' -t strarr <<<"$command_output"
		# printf "\n"
		# echo ${#strarr[@]}
		# # Print each value of the array by using loop
		# for ((n = 0; n < ${#strarr[@]}; n++)); do
		# 	echo $n
		# 	if [[ ${strarr[n]} =~ ^.*$ ]]; then
		# 		#if [[ ${strarr[n]} =~ ^.*\b[0-9a-f]{8}\b-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-\b[0-9a-f]{12}\b$ ]]; then
		# 		echo "valid"
		# 	else
		# 		echo "invalid"
		# 	fi
		# 	echo "${strarr[n]}"
		# done
		echo $command_output
		IFS=$(echo -en "\n\b")                            # space is set as delimiter
		echo ${#command_ouput[@]}
#		for ROW in ${command_output[@]}
#			do
#  				echo "$ROW"
#			done

#		read -ra ADDR <<<"$command_output" # str is read into an array as tokens separated by IFS
#		echo ${#ADDR[@]}
#		for i in "${ADDR[@]}"; do # access each element of array
#			echo "$i"
#		done
	fi
}

ks_status
