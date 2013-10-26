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


Examples

		glog.Info("Always printed")
		
		if glog.DebugEnabled() {
			glog.Info("Printed only if at least on DEBUG level")
		}

		glog.Debug("Printed only if at least on DEBUG level")

		if glog.TraceEnabled() {
			glog.Info("Printed only if at least on TRACE level")
		}
		
		glog.Trace("Printed only is at least on TRACE level")

- There is no Debugf, use glog.DebugEnabled() to wrap the Infof call to avoid computation of expensive arguments when severity level is lower than DEBUG.

- There is no Tracef, use glog.TraceEnabled() to wrap the Infof call to avoid computation of expensive arguments when severity level is lower than DEBUG.

* * *
glog is copyright 2013 Google Inc. All Rights Reserved.

glog modifications are copyright 2013 Ernest Micklei. All Rights Reserved.