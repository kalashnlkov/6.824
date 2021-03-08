package mr

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"os"
	"sync"
)

// Job balabala
type Job struct {
	mu       sync.Mutex
	WorkerID int
	Status   int
	Filename string
}

// Jobs contains job
type Jobs struct {
	Job []*Job
}

// Coordinator is in charge of worker.
type Coordinator struct {
	// Your definitions here.
	Jobs   Jobs
	Status bool
}

// Your code here -- RPC handlers for the worker to call.

//
// an example RPC handler.
//
// the RPC argument and reply types are defined in rpc.go.
//
func (c *Coordinator) Example(args *ExampleArgs, reply *ExampleReply) error {
	reply.Y = args.X + 1
	return nil
}

// GetJob called by worker. suck for job.
func (c *Coordinator) GetJob(workerID int, filename *string) error {
	jobs := c.GetJobs()
	for i := len(jobs); i >= 0; i-- {
		fmt.Print(i)
		job := jobs[i]
		if job.Status == 0 {
			*filename = job.Filename
			job.Status = 1
			job.WorkerID = workerID
			return nil
		}
	}
	return fmt.Errorf("no work here")
}

// HandinJob called by worker. submit for job.
func (c *Coordinator) HandinJob(workerID int, filename *string) error {
	jobs := c.GetJobs()
	for i := len(jobs); i >= 0; i-- {
		fmt.Print(i)
		job := jobs[i]
		if job.Status == 0 {
			*filename = job.Filename
			job.Status = 1
			job.WorkerID = workerID
			return nil
		}
	}
	return fmt.Errorf("we received a wrong job")
}

//
// start a thread that listens for RPCs from worker.go
//
func (c *Coordinator) server() {
	rpc.Register(c)
	rpc.HandleHTTP()
	//l, e := net.Listen("tcp", ":1234")
	sockname := coordinatorSock()
	os.Remove(sockname)
	l, e := net.Listen("unix", sockname)
	if e != nil {
		log.Fatal("listen error:", e)
	}
	go http.Serve(l, nil)
}

//
// main/mrcoordinator.go calls Done() periodically to find out
// if the entire job has finished.
//
func (c *Coordinator) Done() bool {
	return c.Status
	// Your code here.
}

// GetJobs return Coordinator's all job
func (c *Coordinator) GetJobs() []*Job {
	return c.Jobs.Job
}

//
// create a Coordinator.
// main/mrcoordinator.go calls this function.
// nReduce is the number of reduce tasks to use.
//
func MakeCoordinator(files []string, nReduce int) *Coordinator {
	c := Coordinator{Jobs{make([]*Job, len(files))}, false}
	for i, file := range files {
		fmt.Printf("index %d, filename %s", i, file)
		j := new(Job)
		j.WorkerID, j.Status, j.Filename = -1, -1, file
		c.Jobs.Job[i] = j
		fmt.Println(c.Jobs.Job)
	}

	c.server()
	return &c
}
