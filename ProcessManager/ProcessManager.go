/*
ProcessManager

uses to build DependencyInjection Container
 */
package ProcessManager

type FinishPoint string

// if some condition is not exist ,they are error return
type Process struct{
	BeforeFinishPoint FinishPoint   // this process must happen after this FinishPoint come out  (may be empty)
	AfterFinishPoint FinishPoint    // this process must happen before this FinishPoint come out (may be empty)
	ResultCondition FinishPoint // this process will product a FinishPoint ,if it runs       (may be empty)
	Work func()error            // the work to do
}

type Manager struct{
	processes []*Process
}
func (manager *Manager)AddProcess(process *Process){
	manager.processes = append(manager.processes,process)
}
func (manager *Manager)Run()error{
	return nil
}
