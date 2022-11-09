# !/bin/bash
sleep 1

echo "Importing csv file into mongodb "
mongoimport --host localhost --db bookInfo --collection bookInfos --file=/docker-entrypoint-initdb.d/bookInfo.csv --type=csv --headerline
