#!/bin/bash
input=$1
if [ $# -eq 0 ]; then
    echo "No arguments supplied"
    exit 1
fi
if [[ $input == "-h" ]]; then
	echo usage: "bash setup.sh src/configuration/FILENAME.go"
	exit 0
fi
if [ ! -f $1 ]; then
	echo "enter a valid .go file"
	echo usage: "bash setup.sh src/configuration/FILENAME.go"
	exit 1
fi
echo "starting ..."
found="false"
counter=1
IPvars=()
Ports=()
while read -r line
do
	if [[ $line == *"differentIPs"* ]]; then
		#echo "$line"
		read differentIPs <<<${line//[^0-9]/ }
		#echo $differentIPs
		found="true"
	fi
	if [[ $found = "true" ]]; then
		#echo $found
		while [ $counter -lt $differentIPs ] && read -r newLine; 
		do
			if [[ $newLine == *"const otherIP"* ]]; then
				IPvar=$( echo "$newLine" | cut -d '"' -f2 )
				echo adding ..."$IPvar"
				IPvars+=($IPvar)
				let counter=counter+1
			fi			
		done < "$input"
		while read -r lineThisMachine
		do
			if [[ $lineThisMachine == *"THIS MACHINE"* ]]; then
			read lineThisMachine
			read lineThisMachine
				while [[ $lineThisMachine != *"// --"* ]];
				do
					if [[ $lineThisMachine == *"IP"* ]]; then
						port=$( echo "$lineThisMachine" | cut -d '"' -f2 )
						#echo $port
						port=$( echo "$port" | cut -d ':' -f2 )
						echo $port
						Ports+=($port)
					fi
					read lineThisMachine
				done
				for ip in "${IPvars[@]}";
				do
					for currentPort in "${Ports[@]}";
					do
						sudo iptables -I INPUT -p tcp -s $ip --dport $currentPort -j ACCEPT
					done
				done
				exit 0
			fi	
		done < "$input"
		exit 0
	
	fi
done < "$input"
if [[ $found = "false" ]]; then
	echo ERROR please have the variable differentIPs with the number of different IP addresses you are using		
	exit 1
fi
echo "ending ..."
#sudo iptables -I INPUT -p tcp -s 10.0.0.11 --dport 8022 -j ACCEPT

