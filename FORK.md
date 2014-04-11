glog (fork)
==

This fork adds the following features:
- writes logstash JSON messages asynchronuously to a provided io.Writer.
- convencience methods for DEBUG and TRACE level logging.

Additional flags

- -logstash=false
	
	Logs are also written to the Writer that is setup by SetLogstashWriter.

- -logstash.capacity=1000

	How many messages can be queued for asynchronuous writes.

Setup the logstash destination

		glog.SetLogstashWriter(aWriter)

> Provide an io.Writer to write the JSON representation of log events.
> This can a file, an UDP connection or any other implementation.

Passing extra fields to log messages (will be part of @fields)

		ExtraFields["instance"] = "ps34"
		ExtraFields["role"] = "webservice"
		
Sample

		{"@source":"MacErnest"
		,"@type":"glog","@timestamp":"2014-03-21T10:52:05.495118455+01:00"
		,"@fields":{"level":"INFO","threadid":02628,"file":"glog_logstash_test.go","line":60,"instance":"ps34"
		,"role":"webservice"
		}
		,"@message":"hello"
		}				

Examples of severity levels DEBUG(=10) and TRACE(=100)

		glog.Info("Always printed")
		glog.Infof("Printed on %v", time.Now())
		
		glog.Debug("Printed only if at least on DEBUG level")
		
		if glog.DebugEnabled() {
			glog.Debugf("Printed only if at least on %s level", "DEBUG")
		}

		glog.Trace("Printed only is at least on TRACE level")
		
		if glog.TraceEnabled() {
			glog.Tracef("Printed only if at least on %s level", "TRACE")
		}		

* * *
glog is copyright 2013 Google Inc. All Rights Reserved.

glog modifications are copyright 2013 Ernest Micklei. All Rights Reserved.