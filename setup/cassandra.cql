DROP KEYSPACE IF EXISTS timed;
CREATE KEYSPACE timed WITH replication = {'class': 'SimpleStrategy', 'replication_factor': 1};
USE timed;

CREATE TABLE data(
	partition int,
	period timestamp,
	shop_id bigint,
	sample text,
	PRIMARY KEY (partition, period)
);
CREATE INDEX data_period ON data(period);

INSERT INTO data(partition, period, shop_id, sample) VALUES (1, '2015-01-01T00:00:00Z', 100, 'hola');
INSERT INTO data(partition, period, shop_id, sample) VALUES (1, '2015-01-02T00:00:00Z', 100, 'hola');
INSERT INTO data(partition, period, shop_id, sample) VALUES (1, '2015-01-03T00:00:00Z', 100, 'hola');
INSERT INTO data(partition, period, shop_id, sample) VALUES (1, '2015-01-04T00:00:00Z', 100, 'hola');
INSERT INTO data(partition, period, shop_id, sample) VALUES (1, '2015-01-05T00:00:00Z', 100, 'hola');
