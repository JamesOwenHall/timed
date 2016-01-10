DROP KEYSPACE IF EXISTS demo;
CREATE KEYSPACE timed WITH replication = {'class': 'SimpleStrategy', 'replication_factor': 1};
