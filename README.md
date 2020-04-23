# clean

`clean` deletes existing processes in the session, except for PID 0 and a list of process exceptions. 

It is used by `okteto up` to remove its dependency on `ps`.
