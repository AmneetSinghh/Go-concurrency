package main

/*
timeouts:

* if request is in queue for long time, then it should be timed out.
* if request is taking long time to process, then it should be timed out.
* If you don’t have the resources to store the requests (e.g., memory for in- memory queues, disk space for persisted queues).
* Sometimes data has a window within which it must be processed before more relevant data is available, or the need to process the data has expired.
  If a concur‐ rent process takes longer to process the data than this window, we would want to time out and cancel the concurrent process.
*

cancellation:
* its a good practice for giving user control over the cancel long running transactions. Lets say in app, we show button make pdf, but now its taking time, loading = append(* its a good practice for giving user control over the cancel long running transactions. Lets say in app, we show button make pdf, but now its taking time, loading,
 we give user button to cancel the process.



 When concurrent process should be cancelled :

 * long running goroutines should have a timeout at each level and its preemtable.
 * if our goroutine happens to modify shared state—e.g., a database, a file, an in-memory data structure—what
   happens when the goroutine is canceled? Does your goroutine try and roll back the intermedi‐ ary work it’s done
* Duplicate messages
	- stage A sends message to Stage B, now Stage B received but doesn't processed, now now Stage B cancelled. Now its up again, and stage A again sends message.
    - For solving this problem make stage B idompotnet or make heartbeat mechanism


*/
