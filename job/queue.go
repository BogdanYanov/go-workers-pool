package job

// Queue is the queue that the work enters
type Queue chan Job

// QueueInit initializes the queue with the initial queue capacity parameter.
func QueueInit(maxJobs int) Queue {
	if maxJobs < 1 {
		return make(chan Job)
	}
	return make(chan Job, maxJobs)
}

// SendJob sending job to queue.
func (q Queue) SendJob(job Job) {
	q <- job
}
