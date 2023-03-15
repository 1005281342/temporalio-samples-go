package helloworld

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.temporal.io/sdk/activity"
	"go.temporal.io/sdk/workflow"

	// TODO(cretz): Remove when tagged
	_ "go.temporal.io/sdk/contrib/tools/workflowcheck/determinism"
)

// Workflow is a Hello World workflow definition.
func Workflow(ctx workflow.Context, name string) (string, error) {
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: 10 * time.Second,
	}
	ctx = workflow.WithActivityOptions(ctx, ao)

	logger := workflow.GetLogger(ctx)
	logger.Info("HelloWorld workflow started", "name", name)

	var result string
	err := workflow.ExecuteActivity(ctx, Activity, name).Get(ctx, &result)
	if err != nil {
		logger.Error("Activity failed.", "Error", err)
		return "", err
	}

	logger.Info("HelloWorld workflow completed.", "result", result)

	return result, nil
}

func Activity(ctx context.Context, name string) (string, error) {
	logger := activity.GetLogger(ctx)
	logger.Info("Activity", "name", name)
	return "Hello " + name + "!", nil
}

func MyWorkflow(ctx workflow.Context, name string) (string, error) {

	ao := workflow.ActivityOptions{
		StartToCloseTimeout: 12 * time.Second,
	}
	ctx = workflow.WithActivityOptions(ctx, ao)

	// 启动任务 A
	aFuture := workflow.ExecuteActivity(ctx, TaskA)
	_ = aFuture

	// 启动任务 B 和任务 C，并等待它们执行完成
	bFuture := workflow.ExecuteActivity(ctx, TaskB)
	cFuture := workflow.ExecuteActivity(ctx, TaskC)

	// 等待任务 B 和任务 C 执行完成
	var bResult, cResult string
	bErr := bFuture.Get(ctx, &bResult)
	cErr := cFuture.Get(ctx, &cResult)
	if bErr != nil {
		return "", bErr
	}
	if cErr != nil {
		return "", cErr
	}

	// 启动任务 D，并等待它执行完成
	dFuture := workflow.ExecuteActivity(ctx, TaskD, bResult, cResult)
	var dResult string
	if err := dFuture.Get(ctx, &dResult); err != nil {
		return "", err
	}

	// 打印任务 D 的执行结果
	fmt.Println("Task D result:", dResult)

	return dResult, nil
}

func TaskA(ctx context.Context) (string, error) {
	// 执行任务 A 的逻辑
	return "Task A result", nil
}

func TaskB(ctx context.Context) (string, error) {
	// 执行任务 B 的逻辑
	time.Sleep(time.Second * 10)
	return "Task B result", nil
}

func TaskC(ctx context.Context) (string, error) {
	// 执行任务 C 的逻辑
	time.Sleep(time.Second * 10)
	return "Task C result", nil
}

func TaskD(ctx context.Context, bResult string, cResult string) (string, error) {
	log.Printf("b:%s, c:%s\n", bResult, cResult)
	// 执行任务 D 的逻辑
	return "Task D result", nil
}
