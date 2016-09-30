#Kafka Connect, Kafka Streams, Spark Streams, Mapr Streams Lab

####This is a lab for testing out different ingestion methods with tools like Kafka, Confluent, Mapr and Spark

###Create topic
```
root@kafka1:/opt/kafka_2.11-0.10.0.1/bin# ./kafka-topics.sh --create --topic test --partitions 3 --replication-factor 3 --zookeeper kafka1.lan:2181,kafka2.lan:2181,kafka3.lan:2181
```
###Producer
```
./kafka-console-producer.sh --broker-list kafka1.lan:9092,kafka2.lan:9092,kafka3.lan:9092 --topic test
```
###Consumer
```
./kafka-console-consumer.sh --zookeeper kafka1.lan:2181,kafka2.lan:2181,kafka3.lan:2181 --topic test --from-beginning
```

###Producer Performance
```
kafka@kafka3:/opt/kafka_2.11-0.10.0.1/bin$ ./kafka-producer-perf-test.sh --topic test --num-records 5000 --producer-props bootstrap.servers=kafka1.lan:9092,kafka2.lan:9092,kafka3.lan:9092 --throughput 10 --record-size 1024
```

###Confluent Kafka AvroProducer
```
kafka-avro-console-producer --broker-list kafka1.lan:9092,kafka2.lan:9092,kafka3.lan:9092  --topic test --property value.schema='{"type":"record","name":"myrecord","fields":[{"name":"f1","type":"string"}]}'
```

##Connect HDFS
####connect-standalone now, will be distributed soon:
```
kafka-avro-console-producer --broker-list kafka1.lan:9092,kafka2.lan:9092,kafka3.lan:9092 --topic test_hdfs --property value.schema='{"type":"record","name":"myrecord","fields":[{"name":"f1","type":"string"}]}'
```

####Manually starting the connect-hdfs This embeds a REST API Server and accepts additional properties at the command line. The distributed connector does not. See below. (This works with hadoop..although I cant get it working with maprfs:/// yet.):
```
connect-standalone /etc/schema-registry/connect-avro-standalone.properties /etc/kafka-connect-hdfs/quickstart-hdfs.properties
```

####Connect Distributed with HDFS, must use the REST API to pass properties to the kafka connect-distributed process:

On hadoop1.lan, create the HDFS filesystem. Additionally, setup the /topics and /logs folders so Kafka can write to HDFS.
```
hdfs namenode -format
/opt/hadoop-2.7.2/sbin/start-dfs.sh
/opt/hadoop-2.7.2/sbin/start-yarn.sh
jps
hadoop fs -mkdir /topics
hadoop fs -chmod 0777 /topics
hadoop fs -mkdir /logs
hadoop fs -chmod 0777 /logs
```

The connect distributed process has an upstart script running on the kafka cluster /etc/init/connect-standalone.conf, like so:
```
connect-distributed /etc/schema-registry/connect-avro-distributed.properties
```

Have to use the embedded REST API to submit connectors (The following json is POSTed to the connect-distributed embedded Jetty API server):
```
{
    "name": "hdfs-sink-connector-distributed",
    "config": {
        "connector.class": "io.confluent.connect.hdfs.HdfsSinkConnector",
        "tasks.max": "1",
        "topics": "test_distributed",
        "topics.dir": "topics",
        "logs.dir": "logs",
        "hdfs.url": "hdfs://hadoop1.lan:9000",
        "hadoop.home": "/opt/hadoop-2.7.2/",
        "hadoop.conf.dir": "/opt/hadoop-2.7.2/etc/hadoop/",
        "flush.size": "3",
        "partitioner.class": "io.confluent.connect.hdfs.partitioner.DefaultPartitioner",
        "format.class": "io.confluent.connect.hdfs.avro.AvroFormat",
        "storage.class": "io.confluent.connect.hdfs.storage.HdfsStorage",
        "hive.integration": "false",
        "filename.offset.zero.pad.width": "10",
        "rotate.interval.ms": "-1",
        "shutdown.timeout.ms": "3000",
        "hdfs.authentication.kerberos": "false",
        "schema.cache.size": "1000",
        "schema.compatibility": "NONE",
        "partition.duration.ms": "-1",
        "retry.backoff.ms": "5000"
    }
}
```
Additional standalone values that may need to be brought into connect-distributed json data:
```
	filename.offset.zero.pad.width = 10
	topics.dir = topics
	flush.size = 3
	connect.hdfs.principal = 
	timezone = 
	hive.home = 
	hive.database = default
	rotate.interval.ms = -1
	retry.backoff.ms = 5000
	locale = 
	hadoop.home = 
	logs.dir = logs
	schema.cache.size = 1000
	format.class = io.confluent.connect.hdfs.avro.AvroFormat
	hive.integration = false
	hdfs.namenode.principal = 
	hive.conf.dir = 
	partition.duration.ms = -1
	hadoop.conf.dir = 
	schema.compatibility = NONE
	connect.hdfs.keytab = 
	hdfs.url = hdfs://hadoop1.lan:9000
	hdfs.authentication.kerberos = false
	hive.metastore.uris = 
	partition.field.name = 
	kerberos.ticket.renew.period.ms = 3600000
	shutdown.timeout.ms = 3000
	partitioner.class = io.confluent.connect.hdfs.partitioner.DefaultPartitioner
	storage.class = io.confluent.connect.hdfs.storage.HdfsStorage
```

Now POST the json data to the Embedded Jetty server /connectors endpoint:
```
curl -X POST -H "Content-Type: application/json" --data-binary @hdfs.json http://kafka1.lan:8083/connectors
```

Start the avro producer and enter the following:

```
kafka-avro-console-producer --broker-list kafka1.lan:9092,kafka2.lan:9092,kafka3.lan:9092 --topic test_distributed --property value.schema='{"type":"record","name":"myrecord","fields":[{"name":"f1","type":"string"}]}'
```

```
{"f1":"value1"}
{"f1":"value2"}
{"f1":"value3"}
```

Sarama HTTP Go Server:
```
%7B%22f1%22%3A%22val232323%22%7D

###Confluent Kafka AvroConsumer
```
kafka-avro-console-consumer --zookeeper kafka1.lan:2181,kafka2.lan:2181,kafka3.lan:2181  --topic test
```


```
kafka-topics --zookeeper kafka1.lan:2181,kafka2.lan:2181,kafka3.lan:2181 --create --topic poem --partitions 1 --replication-factor 1
```

###Confluent Rest

###Avro
####Register Schema:
```
curl -X POST -H "Content-Type: application/vnd.kafka.avro.v1+json" --data '{"value_schema": "{\"type\": \"record\", \"name\": \"User\", \"fields\": [{\"name\": \"name\", \"type\": \"string\"}]}", "records": [{"value": {"name": "testUser"}}]}' "http://kafka1.lan:8082/topics/avrotest"
```

####Create Consumer:
```
curl -X POST -H "Content-Type: application/vnd.kafka.v1+json" --data '{"name": "my_consumer_instance", "format": "avro", "auto.offset.reset": "smallest"}' http://kafka1.lan:8082/consumers/my_avro_consumer
```

####Get values
```
curl -X GET -H "Accept: application/vnd.kafka.avro.v1+json" http://kafka1.lan:8082/consumers/my_avro_consumer/instances/my_consumer_instance/topics/avrotest
```

####Clean up resources
```
curl -X DELETE http://kafka1.lan:8082/consumers/my_avro_consumer/instances/my_consumer_instance
```

###Json

####Produce:
```
curl -X POST -H "Content-Type: application/vnd.kafka.json.v1+json" --data '{"records":[{"value":{"foo":"bar"}}]}' "http://kafka1.lan:8082/topics/jsontest"
```

####Create Consumer:
```
curl -X POST -H "Content-Type: application/vnd.kafka.v1+json" --data '{"name": "my_consumer_instance", "format": "json", "auto.offset.reset": "smallest"}' http://kafka1.lan:8082/consumers/my_json_consumer
```
####Get values
```
curl -X GET -H "Accept: application/vnd.kafka.json.v1+json" http://kafka1.lan:8082/consumers/my_json_consumer/instances/my_consumer_instance/topics/jsontest
```
####Clean up resources
```
curl -X DELETE http://kafka1.lan:8082/consumers/my_json_consumer/instances/my_consumer_instance
```

###Binary
####Produce
```
curl -X POST -H "Content-Type: application/vnd.kafka.binary.v1+json" --data '{"records":[{"value":"S2Fma2E="}]}' "http://kafka1.lan:8082/topics/binarytest"
```

####Consumer
```
curl -X POST -H "Content-Type: application/vnd.kafka.v1+json" --data '{"name": "my_consumer_instance", "format": "binary", "auto.offset.reset": "smallest"}' http://kafka1.lan:8082/consumers/my_binary_consumer
```

####Get Records
```
curl -X GET -H "Accept: application/vnd.kafka.binary.v1+json" http://kafka1.lan:8082/consumers/my_binary_consumer/instances/my_consumer_instance/topics/binarytest
```

####Delete
```
curl -X DELETE http://kafka1.lan:8082/consumers/my_binary_consumer/instances/my_consumer_instance
```