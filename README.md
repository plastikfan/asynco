# asynco

Test repo to learn parallel/concurrent programming in go

⚠️ Even though to project works, it does so in a very limited fashion. Only a single batch of work can be dispatched to the work pool, which may be ok in some circumstances, but a real application is more likely to require being able to stream jobs in multiple sub-batches to avoid excessive memory consumption. Without this facility, if a client program needs to run a large batch, then this will require the entire batch being loaded up front and dispatched in a single operation.

Therefore, this project needs to re-design to accommodate this.
