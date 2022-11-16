docker build -t serdarkalayci/bookinfo-info:$1.1 ./bookInfoAPI/.
docker build -t serdarkalayci/bookinfo-stock:$1.1 ./bookStockAPI/. -f ./bookStockAPI/v1.dockerfile
docker build -t serdarkalayci/bookinfo-stock:$1.2 ./bookStockAPI/. -f ./bookStockAPI/v2.dockerfile
docker push serdarkalayci/bookinfo-info:$1.1
docker push serdarkalayci/bookinfo-stock:$1.1
docker push serdarkalayci/bookinfo-stock:$1.2