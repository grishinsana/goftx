#!/bin/bash

CURTAG=`git describe --abbrev=0 --tags`;
CURTAG="${CURTAG/v/}"

IFS='.' read -a vers <<< "$CURTAG"

MAJ=${vers[0]}
MIN=${vers[1]}
PATCH=${vers[2]}

echo "Current Tag: v$MAJ.$MIN.$PATCH"

for cmd in "$@"
do
	case $cmd in
		"--major")
			((MAJ+=1))
			MIN=0
			PATCH=0
			echo "Incrementing Major Version"
			;;
		"--minor")
			((MIN+=1))
			PATCH=0
			echo "Incrementing Minor Version"
			;;
		"--patch")
			((PATCH+=1))
			echo "Incrementing Bug Version"
			;;
	esac
done
NEWTAG="v$MAJ.$MIN.$PATCH"
echo "Adding Tag: $NEWTAG";
git tag $NEWTAG && git push --tags
