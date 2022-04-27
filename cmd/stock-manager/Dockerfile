FROM jrei/systemd-centos:7

COPY StockManager.viper.yaml /etc/StockManager/StockManager.viper.yaml
COPY stock-manager /usr/local/bin/stock-manager

CMD ["stock-manager", "run"]
