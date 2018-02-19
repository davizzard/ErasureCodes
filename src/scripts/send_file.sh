#!/bin/bash
# NAME
#	send_file.sh
#
# SYNOPSIS
#	send_file.sh [OPTIONS]... SOURCE ...
#
# DESCRIPTION
#	Script which encodes a file into shards and sends them to a server list.
# 	It also creates parity shards to restore the file if some shards are missing.
#	Shards are temporarily stored in the /tmp directory.
#
#	-data <shards>
#		Number of shards to split the data into (default 4)
#	-par <parity-shards>
#    		Number of parity shards (default 2)
#
#	-s1 <server1-ip-address>
#		First server ip address (default 127.0.0.1)
#
#	-s2 <server2-ip-address>
#		Second server ip address (default 127.0.0.1)
#
#	-s3 <server3-ip-address>
#		Thirds server ip address (default 127.0.0.1)
#	
#	-s4 <server4-ip-address>
#		Fourth server ip address (default 127.0.0.1)
#
#	-h
#		Shows script usage


if [ $# -lt 1 ]; then				# Throwing error if there isn't any source file specified
    echo "Error: no source file specified."
    echo "Usage: send_file.sh [OPTIONS] ... SOURCE ..."
fi
if [ "$1" == "-h" ] || [ $# -lt 1 ]; then	# Shows script help if specified or due to incorrect usage
    echo -e "\nNAME"
    echo -e "\tsend_file.sh"
    echo -e "\nSYNOPSIS"
    echo -e "\tsend_file.sh [OPTIONS]... SOURCE ..."
    echo -e "\nDESCRIPTION"
    echo -e "\tScript which encodes a file into shards and sends them to a server list. It also creates parity shards to restore the file if some shards are missing."
    echo -e "\n\nOPTIONS"
    echo -e "\n-data <shards>"
    echo -e "\tNumber of shards to split the data into (default 4)"
    echo -e "\n-par <parity-shards>"
    echo -e "\tNumber of parity shards (default 2)"
    echo -e "\n-s1 <server1-ip-address>"
    echo -e "\tFirst server ip address (default 127.0.0.1)"
    echo -e "\n-s2 <server2-ip-address>"
    echo -e "\tSecond server ip address (default 127.0.0.1)"
    echo -e "\n-s3 <server3-ip-address>"
    echo -e "\tThird server ip address (default 127.0.0.1)"
    echo -e "\n-s4 <server4-ip-address>"
    echo -e "\tFourth server ip address (default 127.0.0.1)"
    echo -e "\n-h"
    echo -e "\tShows how to use the script"

else 
	DATA=4
	PAR=2
	S1="127.0.0.1"
	S2="127.0.0.1"
	S3="127.0.0.1"
	S4="127.0.0.1"
	SHARDS=$(($DATA + $PAR))
	FILEPATH="$1"
	shift	

	while [[ $# -gt 1 ]]
	do
	key="$1"

	case $key in
	    -data)
		    DATA="$2"
		    shift # past argument
	    ;;
	    -par)
		    PAR="$2"
		    shift # past argument
	    ;;
	    -s1)
		    S1="$2"
		    shift # past argument
	    ;;
	    -s2)
		    S2="$2"
		    shift # past argument
	    ;;
	    -s3)
		    S3="$2"
		    shift # past argument
	    ;;
	    -s4)
		    S4="$2"
		    shift # past argument
	    ;;
	    *)
		    # unknown option
		    shift
	    ;;
	esac
	done
	
	# We split the file into shards and store them in the /tmp directory
	echo "$(go run src/github.com/klauspost/reedsolomon/examples/stream-encoder.go -data $DATA -par $PAR -out /tmp $FILEPATH)"

	# We send each shard to one of the servers
	for ((i=0; i<$SHARDS; i+=1))
	do
		echo $i
		s_num=$i%4
		case $s_num in
		0)
			$server=$S1
		;;
		1)
			$server=$S2
		;;
		2)
			$server=$S3
		;;
		3)
			$server=$S4
		;;
		esac
		file_shard="${FILEPATH##*/}.$i"
		echo "$(curl --user root:toortoor --upload-file /tmp/$file_shard http://127.0.0.1:8080/api/latest/file/random/folders/$file_shard )"
	done

fi
