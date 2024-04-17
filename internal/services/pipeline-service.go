package services

import (
	"encoding/json"
	apitypes "github.com/PiotrFerenc/mash2/api/types"
	"github.com/PiotrFerenc/mash2/internal/queues"
	"github.com/PiotrFerenc/mash2/internal/types"
	"github.com/google/uuid"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
)

// PipelineService is an interface for executing a pipeline of tasks.
// It defines a single method Run, which takes a pipeline and returns
// a unique identifier for the execution and an error if any errors occur.
type PipelineService interface {
	Run(pipeline *apitypes.Pipeline) (uuid.UUID, error)
}

type pipelineService struct {
	queue          queues.MessageQueue
	processService ProcessService
}

// CreatePipelineService creates a PipelineService instance by initializing the queue, processService,
// and starting a goroutine to watch for messages. It returns the PipelineService instance.
// The watchForMessages function is called with queue, onSuccess, onFail, onFinish, and processService as arguments.
// The pipelineService struct is created with the initialized queue and processService.
// The method returns the reference to the pipelineService instance.
//
// Example usage:
//
//	pipeLineService := CreatePipelineService(messageQueue, processService, CreateOnSuccessFunc(),
//	  CreateOnFailFunc(), CreateOnFinishFunc())
//
// Parameters:
//
//	queue: The message queue that implements the MessageQueue interface.
//	processService: The process service that implements the ProcessService interface.
//	onSuccess: The function to be called when a message is received successfully.
//	onFail: The function to be called when a message fails to be processed.
//	onFinish: The function to be called when all steps in a message have been processed.
//
// Returns:
//
//	PipelineService: The created PipelineService instance.
func CreatePipelineService(queue queues.MessageQueue, processService ProcessService, onSuccess OnMessageReceivedFunc, onFail OnMessageReceivedFunc, onFinish OnMessageReceivedFunc) PipelineService {
	go watchForMessages(queue, onSuccess, onFail, onFinish, processService)

	return &pipelineService{
		queue:          queue,
		processService: processService,
	}
}

// Run runs the pipeline by creating a new process from the provided pipeline,
// marking the process as started, adding the process to the queue, and returning
// the process ID if successful.
//
// Parameters:
//
//	pipeline: The pipeline to be executed.
//
// Returns:
//
//	uuid.UUID: The ID of the created process.
//	error: An error if the process fails to be added to the queue.
func (p *pipelineService) Run(pipeline *apitypes.Pipeline) (uuid.UUID, error) {
	process := types.NewProcessFromPipeline(pipeline)
	p.processService.MarkAsStarted(process)

	err := p.queue.AddTaskToQueue(*process)
	if err != nil {
		return uuid.Nil, err
	}

	return process.Id, nil
}

// watchForMessages is a function that listens for messages from the message queue
// and processes them using the provided onSuccess, onFail, onFinish, and processService functions.
// It starts three goroutines to process messages for each of the provided functions.
// The function blocks until it receives a signal on the forever channel.
//
// Parameters:
//
//	queue: The message queue that implements the MessageQueue interface.
//	onSuccess: The function to be called when a message is received successfully.
//	onFail: The function to be called when a message fails to be processed.
//	onFinish: The function to be called when all steps in a message have been processed.
//	processService: The process service that implements the ProcessService interface.
func watchForMessages(queue queues.MessageQueue, onSuccess OnMessageReceivedFunc, onFail OnMessageReceivedFunc, onFinish OnMessageReceivedFunc, processService ProcessService) {
	forever := make(chan struct{})
	go ProcessTasks(queue, onSuccess, processService, queue.WaitingForSucceedTask)
	go ProcessTasks(queue, onFail, processService, queue.WaitingForFailedTask)
	go ProcessTasks(queue, onFinish, processService, queue.WaitingForFinishedTask)
	<-forever
}

// ProcessTasks processes tasks from the message queue by getting the tasks using getTasksFunc,
// and then calling processTasks to process each task using taskFunc, processService, and the queue.
// If an error occurs during processing, it panics with the error.
//
// Parameters:
//
// queue: The message queue that implements the MessageQueue interface.
// taskFunc: The function to be called on each task received from the task channel.
// processService: The process service that implements the ProcessService interface.
// getTasksFunc: The function that returns the task channel and an error.
func ProcessTasks(queue queues.MessageQueue, taskFunc OnMessageReceivedFunc, processService ProcessService, getTasksFunc getTasksFunc) {
	tasks, _ := getTasksFunc()
	err := processTasks(tasks, queue, taskFunc, processService)
	if err != nil {
		panic(err)
	}
}

type getTasksFunc func() (<-chan amqp.Delivery, error)

// processTasks processes tasks from the message queue by receiving tasks from the tasks channel.
// It unmarshals each task message and calls the taskFunc to process the task.
// If an error occurs during unmarshaling or processing, it returns the error.
//
// Parameters:
//
// tasks: The channel that delivers the tasks to be processed.
// queue: The message queue that implements the MessageQueue interface.
// taskFunc: The function to be called to process each task received from the tasks channel.
// processService: The process service that implements the ProcessService interface.
//
// Returns:
//
// error: The error that occurred during unmarshaling or processing, if any.
func processTasks(tasks <-chan amqp.Delivery, queue queues.MessageQueue, taskFunc OnMessageReceivedFunc, processService ProcessService) error {
	for d := range tasks {
		message, err := unmarshalMessage(d.Body)
		if err != nil {
			return err
		}
		taskFunc(message, queue, processService)
	}
	return nil
}

// unmarshalMessage unmarshals a byte array into a pointer to a types.Process object.
// It returns the unmarshaled message and any error that occurred during the unmarshaling process.
//
// Parameters:
//
// body: The byte array to be unmarshaled into a types.Process object.
//
// Returns:
//
// *types.Process: The unmarshaled types.Process object.
// error: The error that occurred during unmarshaling, if any.
func unmarshalMessage(body []byte) (*types.Process, error) {
	var message types.Process
	err := json.Unmarshal(body, &message)
	if err != nil {
		return nil, err
	}
	return &message, nil
}

type OnMessageReceivedFunc func(process *types.Process, queue queues.MessageQueue, service ProcessService)

// CreateOnSuccessFunc creates a function of type OnMessageReceivedFunc that processes the given process,
// queue, and service based on the logic provided in the implementation. If the process has no steps,
// it adds the process to the queue as a finished task. Otherwise, it sets the current step of the process,
// removes the current step from the steps list, adds the process to the queue, and calls the TaskFinished
// method of the service. The function returns void.
func CreateOnSuccessFunc() OnMessageReceivedFunc {
	return func(process *types.Process, queue queues.MessageQueue, service ProcessService) {
		if len(process.Steps) == 0 {
			err := queue.AddTaskAsFinished(*process)
			if err != nil {
				panic(err)
			}
		} else {
			currentStep := process.Steps[0]
			process.CurrentStep = currentStep
			process.Steps = process.Steps[1:]
			err := queue.AddTaskToQueue(*process)
			if err != nil {
				panic(err)
			}
			service.TaskFinished(process)
		}
	}

}

// CreateOnFailFunc creates a function of type OnMessageReceivedFunc that logs the failure of a process,
// marks the process as failed in the service, and returns void.
//
// Parameters:
//
//	process: The process that failed.
//	queue: The message queue that implements the MessageQueue interface.
//	service: The process service that implements the ProcessService interface.
//
// Returns:
//
//	void
func CreateOnFailFunc() OnMessageReceivedFunc {
	return func(process *types.Process, queue queues.MessageQueue, service ProcessService) {
		log.Printf("Fail %s => %+v\\n", process.CurrentStep.Name, process.Error)
		service.MarkAsFailed(process, process.Error)
	}
}

// CreateOnFinishFunc creates a function of type OnMessageReceivedFunc that logs the successful completion of a process,
// marks the process as done in the service.
//
// Parameters:
//
//	process: The process that has completed.
//	queue: The message queue that implements the MessageQueue interface.
//	service: The process service that implements the ProcessService interface.
func CreateOnFinishFunc() OnMessageReceivedFunc {
	return func(process *types.Process, queue queues.MessageQueue, service ProcessService) {
		log.Printf("Done %s ", process.Id)
		service.MarkAsDone(process)
	}
}
