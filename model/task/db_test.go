package task

import (
	"context"
	"testing"
	"time"

	"github.com/evergreen-ci/evergreen"
	"github.com/evergreen-ci/evergreen/apimodels"
	"github.com/evergreen-ci/evergreen/db"
	"github.com/evergreen-ci/evergreen/model/annotations"
	"github.com/evergreen-ci/utility"
	adb "github.com/mongodb/anser/db"
	"github.com/pkg/errors"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
)

func checkStatuses(t *testing.T, expected string, toCheck Task) {
	var dbTasks []Task
	aggregation := []bson.M{
		{"$match": bson.M{
			IdKey: toCheck.Id,
		}},
		addDisplayStatus,
	}
	err := db.Aggregate(Collection, aggregation, &dbTasks)
	assert.NoError(t, err)
	assert.Equal(t, expected, dbTasks[0].DisplayStatus)
	assert.Equal(t, expected, toCheck.GetDisplayStatus())
}

func TestFindMergeTaskForVersion(t *testing.T) {
	assert.NoError(t, db.ClearCollections(Collection))
	t1 := &Task{
		Id:               "t1",
		Version:          "abcdef123456",
		CommitQueueMerge: false,
	}
	assert.NoError(t, t1.Insert())

	_, err := FindMergeTaskForVersion("abcdef123456")
	assert.Error(t, err)
	assert.True(t, adb.ResultsNotFound(err))

	t2 := &Task{
		Id:               "t2",
		Version:          "abcdef123456",
		CommitQueueMerge: true,
	}
	assert.NoError(t, t2.Insert())
	t2Db, err := FindMergeTaskForVersion("abcdef123456")
	assert.NoError(t, err)
	assert.Equal(t, "t2", t2Db.Id)
}

func TestFindTasksByIds(t *testing.T) {
	Convey("When calling FindTasksByIds...", t, func() {
		So(db.Clear(Collection), ShouldBeNil)
		Convey("only tasks with the specified ids should be returned", func() {

			tasks := []Task{
				{
					Id: "one",
				},
				{
					Id: "two",
				},
				{
					Id: "three",
				},
			}

			for _, task := range tasks {
				So(task.Insert(), ShouldBeNil)
			}

			dbTasks, err := Find(ByIds([]string{"one", "two"}))
			So(err, ShouldBeNil)
			So(len(dbTasks), ShouldEqual, 2)
			So(dbTasks[0].Id, ShouldNotEqual, "three")
			So(dbTasks[1].Id, ShouldNotEqual, "three")
		})
	})
}
func TestDisplayTasksByVersion(t *testing.T) {
	Convey("When calling DisplayTasksByVersion...", t, func() {
		So(db.Clear(Collection), ShouldBeNil)
		Convey("only tasks that are display tasks should be returned", func() {
			tasks := []Task{
				{
					Id:      "one",
					Version: "v1",
				},
				{
					Id:          "two",
					Version:     "v1",
					DisplayOnly: true,
				},
				{
					Id:            "three",
					Version:       "v1",
					DisplayTaskId: utility.ToStringPtr(""),
				},
				{
					Id:             "four",
					Version:        "v1",
					ExecutionTasks: []string{"execution_task_one, execution_task_two"},
				},
				{
					Id:            "five",
					Version:       "v1",
					ActivatedTime: utility.ZeroTime,
				},
				{
					Id:            "execution_task_one",
					Version:       "v1",
					DisplayTaskId: utility.ToStringPtr("four"),
				},
				{
					Id:            "execution_task_two",
					Version:       "v1",
					DisplayTaskId: utility.ToStringPtr("four"),
				},
			}

			for _, task := range tasks {
				So(task.Insert(), ShouldBeNil)
			}

			dbTasks, err := FindAll(db.Query(DisplayTasksByVersion("v1", false)))
			So(err, ShouldBeNil)
			So(len(dbTasks), ShouldEqual, 4)
			So(dbTasks[0].Id, ShouldNotEqual, "execution_task_one")
			So(dbTasks[1].Id, ShouldNotEqual, "execution_task_one")
			So(dbTasks[2].Id, ShouldNotEqual, "execution_task_one")
			So(dbTasks[3].Id, ShouldNotEqual, "execution_task_one")

			So(dbTasks[0].Id, ShouldNotEqual, "execution_task_two")
			So(dbTasks[1].Id, ShouldNotEqual, "execution_task_two")
			So(dbTasks[2].Id, ShouldNotEqual, "execution_task_two")
			So(dbTasks[3].Id, ShouldNotEqual, "execution_task_two")

			So(dbTasks[0].Id, ShouldNotEqual, "five")
			So(dbTasks[1].Id, ShouldNotEqual, "five")
			So(dbTasks[2].Id, ShouldNotEqual, "five")
			So(dbTasks[3].Id, ShouldNotEqual, "five")

		})
	})
}

func TestNonExecutionTasksByVersion(t *testing.T) {
	assert.NoError(t, db.Clear(Collection))
	displayTask := Task{
		Id:             "dt",
		Version:        "v1",
		DisplayTaskId:  nil, // legacy, not populated
		ExecutionTasks: []string{"exec_task", "legacy_task"},
	}
	regularTask := Task{
		Id:            "t1",
		Version:       "v1",
		DisplayTaskId: utility.ToStringPtr(""),
	}
	wrongVersionTask := Task{
		Id:            "lame_task",
		Version:       "lame_version",
		DisplayTaskId: utility.ToStringPtr(""),
	}
	execTask := Task{
		Id:            "exec_task",
		Version:       "v1",
		DisplayTaskId: utility.ToStringPtr("dt"),
	}
	legacyTask := Task{
		Id:            "legacy_task",
		Version:       "v2",
		DisplayTaskId: nil, // legacy, not populated
	}
	assert.NoError(t, db.InsertMany(Collection, displayTask, regularTask, wrongVersionTask, execTask, legacyTask))

	tasks, err := Find(NonExecutionTasksByVersions([]string{"v1", "v2"}))
	assert.NoError(t, err)
	assert.Len(t, tasks, 3) // doesn't include wrong version or execution task with DisplayTaskId cached
	for _, task := range tasks {
		assert.NotEqual(t, task.Id, "exec_task")
		assert.NotEqual(t, task.Version, "lame_version")
	}
}

func TestFailedTasksByVersion(t *testing.T) {
	Convey("When calling FailedTasksByVersion...", t, func() {
		So(db.Clear(Collection), ShouldBeNil)
		Convey("only tasks with the failed statuses should be returned", func() {

			tasks := []Task{
				{
					Id:      "one",
					Version: "v1",
					Status:  evergreen.TaskFailed,
				},
				{
					Id:      "two",
					Version: "v1",
					Status:  evergreen.TaskSetupFailed,
				},
				{
					Id:      "three",
					Version: "v1",
					Status:  evergreen.TaskSucceeded,
				},
			}

			for _, task := range tasks {
				So(task.Insert(), ShouldBeNil)
			}

			dbTasks, err := Find(FailedTasksByVersion("v1"))
			So(err, ShouldBeNil)
			So(len(dbTasks), ShouldEqual, 2)
			So(dbTasks[0].Id, ShouldNotEqual, "three")
			So(dbTasks[1].Id, ShouldNotEqual, "three")
		})
	})
}

func TestFindTasksByBuildIdAndGithubChecks(t *testing.T) {
	tasks := []Task{
		{
			Id:            "t1",
			BuildId:       "b1",
			IsGithubCheck: true,
		},
		{
			Id:      "t2",
			BuildId: "b1",
		},
		{
			Id:            "t3",
			BuildId:       "b2",
			IsGithubCheck: true,
		},
		{
			Id:            "t4",
			BuildId:       "b2",
			IsGithubCheck: true,
		},
	}

	for _, task := range tasks {
		assert.NoError(t, task.Insert())
	}
	dbTasks, err := FindAll(db.Query(ByBuildIdAndGithubChecks("b1")))
	assert.NoError(t, err)
	assert.Len(t, dbTasks, 1)
	dbTasks, err = FindAll(db.Query(ByBuildIdAndGithubChecks("b2")))
	assert.NoError(t, err)
	assert.Len(t, dbTasks, 2)
	dbTasks, err = FindAll(db.Query(ByBuildIdAndGithubChecks("b3")))
	assert.NoError(t, err)
	assert.Len(t, dbTasks, 0)
}

func TestFindOneIdAndExecutionWithDisplayStatus(t *testing.T) {
	assert := assert.New(t)
	assert.NoError(db.ClearCollections(Collection, OldCollection))
	taskDoc := Task{
		Id:        "task",
		Status:    evergreen.TaskSucceeded,
		Activated: true,
	}
	assert.NoError(taskDoc.Insert())
	task, err := FindOneIdAndExecutionWithDisplayStatus(taskDoc.Id, utility.ToIntPtr(0))
	assert.NoError(err)
	assert.NotNil(task)
	assert.Equal(task.DisplayStatus, evergreen.TaskSucceeded)

	// Should fetch tasks from the old collection
	assert.NoError(taskDoc.Archive())
	task, err = FindOneOldByIdAndExecution(taskDoc.Id, 0)
	assert.NoError(err)
	assert.NotNil(task)
	task, err = FindOneIdAndExecutionWithDisplayStatus(taskDoc.Id, utility.ToIntPtr(0))
	assert.NoError(err)
	assert.NotNil(task)
	assert.Equal(task.OldTaskId, taskDoc.Id)

	// Should fetch recent executions by default
	task, err = FindOneIdAndExecutionWithDisplayStatus(taskDoc.Id, nil)
	assert.NoError(err)
	assert.NotNil(task)
	assert.Equal(task.Execution, 1)
	assert.Equal(task.DisplayStatus, evergreen.TaskSucceeded)

	taskDoc = Task{
		Id:        "task2",
		Status:    evergreen.TaskUndispatched,
		Activated: false,
	}
	assert.NoError(taskDoc.Insert())
	task, err = FindOneIdAndExecutionWithDisplayStatus(taskDoc.Id, utility.ToIntPtr(0))
	assert.NoError(err)
	assert.NotNil(task)
	assert.Equal(task.DisplayStatus, evergreen.TaskUnscheduled)
}

func TestFindOldTasksByID(t *testing.T) {
	assert := assert.New(t)
	assert.NoError(db.ClearCollections(Collection, OldCollection))

	taskDoc := Task{
		Id:     "task",
		Status: evergreen.TaskSucceeded,
	}
	assert.NoError(taskDoc.Insert())
	assert.NoError(taskDoc.Archive())
	taskDoc.Execution += 1
	assert.NoError(taskDoc.Archive())
	taskDoc.Execution += 1

	tasks, err := FindOld(ByOldTaskID("task"))
	assert.NoError(err)
	assert.Len(tasks, 2)
	assert.Equal(0, tasks[0].Execution)
	assert.Equal("task_0", tasks[0].Id)
	assert.Equal("task", tasks[0].OldTaskId)
	assert.Equal(1, tasks[1].Execution)
	assert.Equal("task_1", tasks[1].Id)
	assert.Equal("task", tasks[1].OldTaskId)
}

func TestFindAllFirstExecution(t *testing.T) {
	require.NoError(t, db.ClearCollections(Collection, OldCollection))
	tasks := []Task{
		{Id: "t0"},
		{Id: "t1", Execution: 1},
		{Id: "t2", DisplayOnly: true},
	}
	for _, task := range tasks {
		require.NoError(t, task.Insert())
	}
	oldTask := Task{Id: MakeOldID("t1", 0)}
	require.NoError(t, db.Insert(OldCollection, &oldTask))

	foundTasks, err := FindAllFirstExecution(All)
	assert.NoError(t, err)
	assert.Len(t, foundTasks, 3)
	expectedIDs := []string{"t0", MakeOldID("t1", 0), "t2"}
	for _, task := range foundTasks {
		assert.Contains(t, expectedIDs, task.Id)
		assert.Equal(t, 0, task.Execution)
	}
}

func TestFindOneIdOldOrNew(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	require.NoError(db.ClearCollections(Collection, OldCollection))

	taskDoc := Task{
		Id:     "task",
		Status: evergreen.TaskSucceeded,
	}
	require.NoError(taskDoc.Insert())
	require.NoError(taskDoc.Archive())

	task00, err := FindOneIdOldOrNew("task", 0)
	assert.NoError(err)
	require.NotNil(task00)
	assert.Equal(0, task00.Execution)

	task01, err := FindOneIdOldOrNew("task", 1)
	assert.NoError(err)
	require.NotNil(task01)
	assert.Equal("task", task01.Id)
	assert.Equal(1, task01.Execution)
}

func TestAddHostCreateDetails(t *testing.T) {
	assert.NoError(t, db.ClearCollections(Collection))
	task := Task{Id: "t1", Execution: 0}
	assert.NoError(t, task.Insert())
	errToSave := errors.Wrapf(errors.New("InsufficientCapacityError"), "error trying to start host")
	assert.NoError(t, AddHostCreateDetails(task.Id, "h1", 0, errToSave))
	dbTask, err := FindOneId(task.Id)
	assert.NoError(t, err)
	assert.NotNil(t, dbTask)
	require.Len(t, dbTask.HostCreateDetails, 1)
	assert.Equal(t, dbTask.HostCreateDetails[0].HostId, "h1")
	assert.Contains(t, dbTask.HostCreateDetails[0].Error, "InsufficientCapacityError")

	assert.NoError(t, AddHostCreateDetails(task.Id, "h2", 0, errToSave))
	dbTask, err = FindOneId(task.Id)
	assert.NoError(t, err)
	assert.NotNil(t, dbTask)
	assert.Len(t, dbTask.HostCreateDetails, 2)
}

func TestDisplayStatus(t *testing.T) {
	assert.NoError(t, db.ClearCollections(Collection))
	t1 := Task{
		Id:     "t1",
		Status: evergreen.TaskSucceeded,
	}
	assert.NoError(t, t1.Insert())
	checkStatuses(t, evergreen.TaskSucceeded, t1)
	t2 := Task{
		Id:        "t2",
		Status:    evergreen.TaskUndispatched,
		Activated: true,
	}
	assert.NoError(t, t2.Insert())
	checkStatuses(t, evergreen.TaskWillRun, t2)
	t3 := Task{
		Id:        "t3",
		Status:    evergreen.TaskFailed,
		Activated: true,
	}
	assert.NoError(t, t3.Insert())
	checkStatuses(t, evergreen.TaskFailed, t3)
	t4 := Task{
		Id:     "t4",
		Status: evergreen.TaskFailed,
		Details: apimodels.TaskEndDetail{
			Type: evergreen.CommandTypeSetup,
		},
	}
	assert.NoError(t, t4.Insert())
	checkStatuses(t, evergreen.TaskSetupFailed, t4)
	t5 := Task{
		Id:     "t5",
		Status: evergreen.TaskFailed,
		Details: apimodels.TaskEndDetail{
			Type: evergreen.CommandTypeSystem,
		},
	}
	assert.NoError(t, t5.Insert())
	checkStatuses(t, evergreen.TaskSystemFailed, t5)
	t6 := Task{
		Id:     "t6",
		Status: evergreen.TaskFailed,
		Details: apimodels.TaskEndDetail{
			Type:     evergreen.CommandTypeSystem,
			TimedOut: true,
		},
	}
	assert.NoError(t, t6.Insert())
	checkStatuses(t, evergreen.TaskSystemTimedOut, t6)
	t7 := Task{
		Id:     "t7",
		Status: evergreen.TaskFailed,
		Details: apimodels.TaskEndDetail{
			Type:        evergreen.CommandTypeSystem,
			TimedOut:    true,
			Description: evergreen.TaskDescriptionHeartbeat,
		},
	}
	assert.NoError(t, t7.Insert())
	checkStatuses(t, evergreen.TaskSystemUnresponse, t7)
	t8 := Task{
		Id:        "t8",
		Status:    evergreen.TaskStarted,
		Activated: true,
	}
	assert.NoError(t, t8.Insert())
	checkStatuses(t, evergreen.TaskStarted, t8)
	t9 := Task{
		Id:        "t9",
		Status:    evergreen.TaskUndispatched,
		Activated: false,
	}
	assert.NoError(t, t9.Insert())
	checkStatuses(t, evergreen.TaskUnscheduled, t9)
	t10 := Task{
		Id:        "t10",
		Status:    evergreen.TaskUndispatched,
		Activated: true,
		DependsOn: []Dependency{
			{
				TaskId:       "t9",
				Unattainable: true,
				Status:       "success",
			},
			{
				TaskId:       "t8",
				Unattainable: false,
				Status:       "success",
			},
		},
	}
	assert.NoError(t, t10.Insert())
	checkStatuses(t, evergreen.TaskStatusBlocked, t10)
	t11 := Task{
		Id:        "t11",
		Status:    evergreen.TaskUndispatched,
		Activated: true,
		DependsOn: []Dependency{
			{
				TaskId:       "t8",
				Unattainable: false,
				Status:       "success",
			},
		},
	}
	assert.NoError(t, t11.Insert())
	checkStatuses(t, evergreen.TaskWillRun, t11)
	t12 := Task{
		Id:                   "t12",
		Status:               evergreen.TaskUndispatched,
		Activated:            true,
		OverrideDependencies: true,
		DependsOn: []Dependency{
			{
				TaskId:       "t9",
				Unattainable: true,
				Status:       "success",
			},
		},
	}
	// No CheckStatuses for t12 to avoid paradox
	assert.NoError(t, t12.Insert())
	checkStatuses(t, evergreen.TaskWillRun, t11)
	t13 := Task{
		Id:                 "t13",
		Status:             evergreen.TaskUndispatched,
		Activated:          true,
		ContainerAllocated: false,
	}
	require.NoError(t, t13.Insert())
	checkStatuses(t, evergreen.TaskWillRun, t13)
	t14 := Task{
		Id:                 "t14",
		Status:             evergreen.TaskUndispatched,
		Activated:          true,
		ContainerAllocated: true,
	}
	require.NoError(t, t14.Insert())
	checkStatuses(t, evergreen.TaskWillRun, t14)
	t15 := Task{
		Id:                 "t15",
		Status:             evergreen.TaskUndispatched,
		Activated:          false,
		ContainerAllocated: false,
	}
	require.NoError(t, t15.Insert())
	checkStatuses(t, evergreen.TaskUnscheduled, t15)
}

func TestFindTaskNamesByBuildVariant(t *testing.T) {
	Convey("Should return unique task names for a given build variant", t, func() {
		assert.NoError(t, db.ClearCollections(Collection))
		t1 := Task{
			Id:                  "t1",
			Status:              evergreen.TaskSucceeded,
			BuildVariant:        "ubuntu1604",
			DisplayName:         "dist",
			Project:             "evergreen",
			Requester:           evergreen.RepotrackerVersionRequester,
			RevisionOrderNumber: 1,
		}
		assert.NoError(t, t1.Insert())
		t2 := Task{
			Id:                  "t2",
			Status:              evergreen.TaskSucceeded,
			BuildVariant:        "ubuntu1604",
			DisplayName:         "test-agent",
			Project:             "evergreen",
			Requester:           evergreen.RepotrackerVersionRequester,
			RevisionOrderNumber: 1,
		}
		assert.NoError(t, t2.Insert())
		t3 := Task{
			Id:                  "t3",
			Status:              evergreen.TaskSucceeded,
			BuildVariant:        "ubuntu1604",
			DisplayName:         "test-graphql",
			Project:             "evergreen",
			Requester:           evergreen.RepotrackerVersionRequester,
			RevisionOrderNumber: 1,
		}
		assert.NoError(t, t3.Insert())
		t4 := Task{
			Id:                  "t4",
			Status:              evergreen.TaskFailed,
			BuildVariant:        "ubuntu1604",
			DisplayName:         "test-graphql",
			Project:             "evergreen",
			Requester:           evergreen.RepotrackerVersionRequester,
			RevisionOrderNumber: 1,
		}
		assert.NoError(t, t4.Insert())
		buildVariantTask, err := FindTaskNamesByBuildVariant("evergreen", "ubuntu1604", 1)
		assert.NoError(t, err)
		assert.Equal(t, []string{"dist", "test-agent", "test-graphql"}, buildVariantTask)

	})
	Convey("Should only include tasks that appear on mainline commits", t, func() {
		assert.NoError(t, db.ClearCollections(Collection))
		t1 := Task{
			Id:                  "t1",
			Status:              evergreen.TaskSucceeded,
			BuildVariant:        "ubuntu1604",
			DisplayName:         "test-patch-only",
			Project:             "evergreen",
			Requester:           evergreen.PatchVersionRequester,
			RevisionOrderNumber: 1,
		}
		assert.NoError(t, t1.Insert())
		t2 := Task{
			Id:                  "t2",
			Status:              evergreen.TaskSucceeded,
			BuildVariant:        "ubuntu1604",
			DisplayName:         "test-graphql",
			Project:             "evergreen",
			Requester:           evergreen.RepotrackerVersionRequester,
			RevisionOrderNumber: 1,
		}
		assert.NoError(t, t2.Insert())
		t3 := Task{
			Id:                  "t3",
			Status:              evergreen.TaskSucceeded,
			BuildVariant:        "ubuntu1604",
			DisplayName:         "dist",
			Project:             "evergreen",
			Requester:           evergreen.PatchVersionRequester,
			RevisionOrderNumber: 1,
		}
		assert.NoError(t, t3.Insert())
		t4 := Task{
			Id:                  "t4",
			Status:              evergreen.TaskFailed,
			BuildVariant:        "ubuntu1604",
			DisplayName:         "test-something",
			Project:             "evergreen",
			Requester:           evergreen.RepotrackerVersionRequester,
			RevisionOrderNumber: 1,
		}
		assert.NoError(t, t4.Insert())
		buildVariantTasks, err := FindTaskNamesByBuildVariant("evergreen", "ubuntu1604", 1)
		assert.NoError(t, err)
		assert.Equal(t, []string{"test-graphql", "test-something"}, buildVariantTasks)
	})

}

func TestFindNeedsContainerAllocation(t *testing.T) {
	defer func() {
		assert.NoError(t, db.Clear(Collection))
	}()
	getTaskThatNeedsContainerAllocation := func() Task {
		return Task{
			Id:                 utility.RandomString(),
			Activated:          true,
			ActivatedTime:      time.Now(),
			Status:             evergreen.TaskUndispatched,
			ContainerAllocated: false,
			ExecutionPlatform:  ExecutionPlatformContainer,
		}
	}
	for tName, tCase := range map[string]func(t *testing.T){
		"IncludesOneContainerTaskWaitingForAllocation": func(t *testing.T) {
			needsAllocation := getTaskThatNeedsContainerAllocation()
			require.NoError(t, needsAllocation.Insert())

			found, err := FindNeedsContainerAllocation()
			require.NoError(t, err)
			require.Len(t, found, 1)
			assert.Equal(t, needsAllocation.Id, found[0].Id)
		},
		"IncludesAllContainerTasksWaitingForAllocation": func(t *testing.T) {
			needsAllocation0 := getTaskThatNeedsContainerAllocation()
			require.NoError(t, needsAllocation0.Insert())
			needsAllocation1 := getTaskThatNeedsContainerAllocation()
			needsAllocation1.ActivatedTime = time.Now().Add(-time.Hour)
			require.NoError(t, needsAllocation1.Insert())
			doesNotNeedAllocation := getTaskThatNeedsContainerAllocation()
			doesNotNeedAllocation.Activated = false
			require.NoError(t, doesNotNeedAllocation.Insert())

			found, err := FindNeedsContainerAllocation()
			require.NoError(t, err)
			require.Len(t, found, 2)
			assert.Equal(t, needsAllocation1.Id, found[0].Id, "tasks should be sorted by activation time, so task with earlier activation time should be first")
			assert.Equal(t, needsAllocation0.Id, found[1].Id, "tasks should be sorted by activation time, so task with later activation time should be second")
		},
		"IncludesTasksWithAllDependenciesMet": func(t *testing.T) {
			needsAllocation := getTaskThatNeedsContainerAllocation()
			needsAllocation.DependsOn = []Dependency{
				{
					TaskId:   "dependency0",
					Finished: true,
				},
				{
					TaskId:   "dependency1",
					Status:   evergreen.TaskFailed,
					Finished: true,
				},
			}
			require.NoError(t, needsAllocation.Insert())

			found, err := FindNeedsContainerAllocation()
			require.NoError(t, err)
			require.Len(t, found, 1)
			assert.Equal(t, needsAllocation.Id, found[0].Id)
		},
		"IncludesTasksWithOverriddenDependencies": func(t *testing.T) {
			overriddenDependencies := getTaskThatNeedsContainerAllocation()
			overriddenDependencies.DependsOn = []Dependency{
				{
					TaskId: "dependency0",
				},
			}
			overriddenDependencies.OverrideDependencies = true
			require.NoError(t, overriddenDependencies.Insert())

			found, err := FindNeedsContainerAllocation()
			require.NoError(t, err)
			require.Len(t, found, 1)
			assert.Equal(t, overriddenDependencies.Id, found[0].Id)
		},
		"IgnoresTasksWithUnmetDependencies": func(t *testing.T) {
			unmetDependencies := getTaskThatNeedsContainerAllocation()
			unmetDependencies.DependsOn = []Dependency{
				{
					TaskId: "dependency0",
				},
			}
			require.NoError(t, unmetDependencies.Insert())

			found, err := FindNeedsContainerAllocation()
			require.NoError(t, err)
			assert.Empty(t, found)
		},
		"IgnoresTasksWithoutExecutionPlatform": func(t *testing.T) {
			doesNotNeedAllocation := getTaskThatNeedsContainerAllocation()
			doesNotNeedAllocation.ExecutionPlatform = ""
			require.NoError(t, doesNotNeedAllocation.Insert())

			found, err := FindNeedsContainerAllocation()
			require.NoError(t, err)
			assert.Empty(t, found)
		},
		"IgnoresHostTasks": func(t *testing.T) {
			doesNotNeedAllocation := getTaskThatNeedsContainerAllocation()
			doesNotNeedAllocation.ExecutionPlatform = ExecutionPlatformHost
			require.NoError(t, doesNotNeedAllocation.Insert())

			found, err := FindNeedsContainerAllocation()
			require.NoError(t, err)
			assert.Empty(t, found)
		},
		"IgnoresDeactivatedTasks": func(t *testing.T) {
			doesNotNeedAllocation := getTaskThatNeedsContainerAllocation()
			doesNotNeedAllocation.Activated = false
			doesNotNeedAllocation.ActivatedTime = utility.ZeroTime
			require.NoError(t, doesNotNeedAllocation.Insert())

			found, err := FindNeedsContainerAllocation()
			require.NoError(t, err)
			assert.Empty(t, found)
		},
		"IgnoresTasksWithContainerAlreadyAllocated": func(t *testing.T) {
			doesNotNeedAllocation := getTaskThatNeedsContainerAllocation()
			doesNotNeedAllocation.ContainerAllocated = true
			require.NoError(t, doesNotNeedAllocation.Insert())

			found, err := FindNeedsContainerAllocation()
			require.NoError(t, err)
			assert.Empty(t, found)
		},
		"IgnoresDisabledTasks": func(t *testing.T) {
			doesNotNeedAllocation := getTaskThatNeedsContainerAllocation()
			doesNotNeedAllocation.Priority = evergreen.DisabledTaskPriority
			require.NoError(t, doesNotNeedAllocation.Insert())

			found, err := FindNeedsContainerAllocation()
			require.NoError(t, err)
			assert.Empty(t, found)
		},
	} {
		t.Run(tName, func(t *testing.T) {
			require.NoError(t, db.Clear(Collection))
			tCase(t)
		})
	}
}

func TestFindByStaleRunningTask(t *testing.T) {
	defer func() {
		assert.NoError(t, db.ClearCollections(Collection))
	}()
	for tName, tCase := range map[string]func(t *testing.T){
		"ReturnsDispatchedStaleTask": func(t *testing.T) {
			tsk := Task{
				Id:            "task",
				Status:        evergreen.TaskDispatched,
				LastHeartbeat: time.Now().Add(-time.Hour),
			}
			require.NoError(t, tsk.Insert())

			found, err := Find(ByStaleRunningTask(30 * time.Minute))
			require.NoError(t, err)
			require.Len(t, found, 1)
			assert.Equal(t, tsk.Id, found[0].Id)
		},
		"ReturnsRunningStaleTask": func(t *testing.T) {
			tsk := Task{
				Id:            "task",
				Status:        evergreen.TaskStarted,
				LastHeartbeat: time.Now().Add(-time.Hour),
			}
			require.NoError(t, tsk.Insert())

			found, err := Find(ByStaleRunningTask(30 * time.Minute))
			require.NoError(t, err)
			require.Len(t, found, 1)
			assert.Equal(t, tsk.Id, found[0].Id)
		},
		"ReturnsMultipleStaleTasks": func(t *testing.T) {
			tasks := []Task{
				{
					Id:            "task0",
					Status:        evergreen.TaskDispatched,
					LastHeartbeat: time.Now().Add(-time.Hour),
				},
				{
					Id:            "task1",
					Status:        evergreen.TaskStarted,
					LastHeartbeat: time.Now().Add(-time.Minute),
				},
				{
					Id:            "task2",
					Status:        evergreen.TaskStarted,
					LastHeartbeat: time.Now().Add(-time.Hour),
				},
			}
			for _, tsk := range tasks {
				require.NoError(t, tsk.Insert())
			}

			found, err := Find(ByStaleRunningTask(30 * time.Minute))
			require.NoError(t, err)
			require.Len(t, found, 2)
			for _, tsk := range found {
				assert.True(t, utility.StringSliceContains([]string{tasks[0].Id, tasks[2].Id}, tsk.Id))
			}
		},
		"IgnoresRunningTaskWithRecentHeartbeat": func(t *testing.T) {
			tsk := Task{
				Id:            "task",
				Status:        evergreen.TaskStarted,
				LastHeartbeat: time.Now().Add(-time.Minute),
			}
			require.NoError(t, tsk.Insert())

			found, err := Find(ByStaleRunningTask(30 * time.Minute))
			require.NoError(t, err)
			assert.Empty(t, found)
		},
		"IgnoresDisplayTasksWithNoHeartbeat": func(t *testing.T) {
			tsk := Task{
				Id:          "id",
				DisplayOnly: true,
				Status:      evergreen.TaskDispatched,
			}
			require.NoError(t, tsk.Insert())

			found, err := Find(ByStaleRunningTask(0))
			require.NoError(t, err)
			assert.Empty(t, found)
		},
	} {
		t.Run(tName, func(t *testing.T) {
			require.NoError(t, db.ClearCollections(Collection))
			tCase(t)
		})
	}
}

func TestGetTasksByVersionExecTasks(t *testing.T) {
	require.NoError(t, db.ClearCollections(Collection))
	// test that we can handle the different kinds of tasks
	t1 := Task{
		Id:            "execWithDisplayId",
		Version:       "v1",
		DisplayTaskId: utility.ToStringPtr("displayTask"),
	}
	t2 := Task{
		Id:            "notAnExec",
		Version:       "v1",
		DisplayTaskId: utility.ToStringPtr(""),
	}

	t3 := Task{
		Id:      "execWithNoId",
		Version: "v1",
	}
	t4 := Task{
		Id:      "notAnExecWithNoId",
		Version: "v1",
	}
	dt := Task{
		Id:             "displayTask",
		Version:        "v1",
		DisplayOnly:    true,
		ExecutionTasks: []string{"execWithDisplayId", "execWithNoId"},
	}
	assert.NoError(t, db.InsertMany(Collection, t1, t2, t3, t4, dt))

	ctx := context.TODO()
	// execution tasks have been filtered outs
	opts := GetTasksByVersionOptions{}
	tasks, count, err := GetTasksByVersion(ctx, "v1", opts)
	assert.NoError(t, err)
	assert.Equal(t, count, 3)
	// alphabetical order
	assert.Equal(t, dt.Id, tasks[0].Id)
	assert.Equal(t, t2.Id, tasks[1].Id)
	assert.Equal(t, t4.Id, tasks[2].Id)
}

func TestGetTasksByVersionIncludeNeverActivatedTasks(t *testing.T) {
	require.NoError(t, db.ClearCollections(Collection))

	inactiveTask := Task{
		Id:            "inactiveTask",
		Version:       "v1",
		ActivatedTime: utility.ZeroTime,
	}

	assert.NoError(t, inactiveTask.Insert())

	ctx := context.TODO()

	// inactive tasks should be included
	opts := GetTasksByVersionOptions{IncludeNeverActivatedTasks: true}
	_, count, err := GetTasksByVersion(ctx, "v1", opts)
	assert.NoError(t, err)
	assert.Equal(t, count, 1)
	// inactive tasks should be excluded
	opts = GetTasksByVersionOptions{IncludeNeverActivatedTasks: false}
	_, count, err = GetTasksByVersion(ctx, "v1", opts)
	assert.NoError(t, err)
	assert.Equal(t, count, 0)
}

func TestGetTasksByVersionAnnotations(t *testing.T) {
	require.NoError(t, db.ClearCollections(Collection, annotations.Collection))
	t1 := Task{
		Id:        "t1",
		Version:   "v1",
		Execution: 2,
		Status:    evergreen.TaskSucceeded,
	}
	t2 := Task{
		Id:        "t2",
		Version:   "v1",
		Execution: 3,
		Status:    evergreen.TaskFailed,
	}
	t3 := Task{
		Id:        "t3",
		Version:   "v1",
		Execution: 1,
		Status:    evergreen.TaskFailed,
	}
	assert.NoError(t, db.InsertMany(Collection, t1, t2, t3))

	a := annotations.TaskAnnotation{
		Id:            "myAnnotation",
		TaskId:        t2.Id,
		TaskExecution: t2.Execution,
		Issues: []annotations.IssueLink{
			{IssueKey: "EVG-1212"},
		},
	}
	assert.NoError(t, a.Upsert())

	ctx := context.TODO()

	opts := GetTasksByVersionOptions{}
	tasks, count, err := GetTasksByVersion(ctx, "v1", opts)
	assert.NoError(t, err)
	assert.Equal(t, count, 3)
	assert.Equal(t, tasks[0].Id, "t1")
	assert.Equal(t, evergreen.TaskSucceeded, tasks[0].DisplayStatus)
	assert.Equal(t, tasks[1].Id, "t2")
	assert.Equal(t, evergreen.TaskKnownIssue, tasks[1].DisplayStatus)
	assert.Equal(t, tasks[2].Id, "t3")
	assert.Equal(t, evergreen.TaskFailed, tasks[2].DisplayStatus)
}

func TestGetTasksByVersionBaseTasks(t *testing.T) {
	require.NoError(t, db.ClearCollections(Collection))

	t1 := Task{
		Id:                  "t1",
		Version:             "v1",
		BuildVariant:        "bv",
		DisplayName:         "displayName",
		Execution:           0,
		Status:              evergreen.TaskSucceeded,
		RevisionOrderNumber: 1,
		Requester:           evergreen.RepotrackerVersionRequester,
		Revision:            "abc123",
	}
	t2 := Task{
		Id:           "t2",
		Version:      "v2",
		BuildVariant: "bv",
		DisplayName:  "displayName",
		Execution:    0,
		Status:       evergreen.TaskFailed,
		Requester:    evergreen.GithubPRRequester,
		Revision:     "abc123",
	}

	t3 := Task{
		Id:                  "t3",
		Version:             "v3",
		BuildVariant:        "bv",
		DisplayName:         "displayName",
		Execution:           0,
		Status:              evergreen.TaskFailed,
		RevisionOrderNumber: 2,
		Requester:           evergreen.RepotrackerVersionRequester,
		Revision:            "abc125",
	}
	assert.NoError(t, db.InsertMany(Collection, t1, t2, t3))

	ctx := context.TODO()

	// Normal Patch builds
	opts := GetTasksByVersionOptions{
		IncludeBaseTasks: true,
		IsMainlineCommit: false,
	}
	tasks, count, err := GetTasksByVersion(ctx, "v2", opts)
	assert.NoError(t, err)
	assert.Equal(t, 1, count)
	assert.Equal(t, "t2", tasks[0].Id)
	assert.Equal(t, evergreen.TaskFailed, tasks[0].DisplayStatus)
	assert.NotNil(t, tasks[0].BaseTask)
	assert.Equal(t, "t1", tasks[0].BaseTask.Id)
	assert.Equal(t, t1.Status, tasks[0].BaseTask.Status)

	// Mainline builds
	opts = GetTasksByVersionOptions{
		IncludeBaseTasks: true,
		IsMainlineCommit: true,
	}
	tasks, count, err = GetTasksByVersion(ctx, "v3", opts)
	assert.NoError(t, err)
	assert.Equal(t, 1, count)
	assert.Equal(t, "t3", tasks[0].Id)
	assert.Equal(t, evergreen.TaskFailed, tasks[0].DisplayStatus)
	assert.NotNil(t, tasks[0].BaseTask)
	assert.Equal(t, "t1", tasks[0].BaseTask.Id)
	assert.Equal(t, t1.Status, tasks[0].BaseTask.Status)
}

func TestGetTasksByVersionSorting(t *testing.T) {
	require.NoError(t, db.ClearCollections(Collection))

	t1 := Task{
		Id:           "t1",
		Version:      "v1",
		BuildVariant: "bv_foo",
		DisplayName:  "displayName_foo",
		Execution:    0,
		Status:       evergreen.TaskSucceeded,
		BaseTask:     BaseTaskInfo{Id: "t1_base", Status: evergreen.TaskSucceeded},
		StartTime:    time.Date(2022, time.April, 7, 23, 0, 0, 0, time.UTC),
		TimeTaken:    time.Minute,
	}
	t2 := Task{
		Id:           "t2",
		Version:      "v1",
		BuildVariant: "bv_bar",
		DisplayName:  "displayName_bar",
		Execution:    0,
		Status:       evergreen.TaskFailed,
		BaseTask:     BaseTaskInfo{Id: "t2_base", Status: evergreen.TaskFailed},
		StartTime:    time.Date(2022, time.April, 7, 23, 0, 0, 0, time.UTC),
		TimeTaken:    25 * time.Minute,
	}
	t3 := Task{
		Id:           "t3",
		Version:      "v1",
		BuildVariant: "bv_qux",
		DisplayName:  "displayName_qux",
		Execution:    0,
		Status:       evergreen.TaskStarted,
		BaseTask:     BaseTaskInfo{Id: "t3_base", Status: evergreen.TaskSucceeded},
		StartTime:    time.Date(2021, time.November, 10, 23, 0, 0, 0, time.UTC),
		TimeTaken:    0,
	}
	t4 := Task{
		Id:           "t4",
		Version:      "v1",
		BuildVariant: "bv_baz",
		DisplayName:  "displayName_baz",
		Execution:    0,
		Status:       evergreen.TaskSetupFailed,
		BaseTask:     BaseTaskInfo{Id: "t4_base", Status: evergreen.TaskSucceeded},
		StartTime:    time.Date(2022, time.April, 7, 23, 0, 0, 0, time.UTC),
		TimeTaken:    2 * time.Hour,
	}

	assert.NoError(t, db.InsertMany(Collection, t1, t2, t3, t4))

	ctx := context.TODO()

	// Sort by display name, asc
	opts := GetTasksByVersionOptions{
		Sorts: []TasksSortOrder{
			{Key: DisplayNameKey, Order: 1},
		},
	}
	tasks, count, err := GetTasksByVersion(ctx, "v1", opts)
	assert.NoError(t, err)
	assert.Equal(t, 4, count)
	assert.Equal(t, "t2", tasks[0].Id)
	assert.Equal(t, "t4", tasks[1].Id)
	assert.Equal(t, "t1", tasks[2].Id)
	assert.Equal(t, "t3", tasks[3].Id)

	// Sort by build variant name, asc
	opts = GetTasksByVersionOptions{
		Sorts: []TasksSortOrder{
			{Key: BuildVariantKey, Order: 1},
		},
	}
	tasks, count, err = GetTasksByVersion(ctx, "v1", opts)
	assert.NoError(t, err)
	assert.Equal(t, 4, count)
	assert.Equal(t, "t2", tasks[0].Id)
	assert.Equal(t, "t4", tasks[1].Id)
	assert.Equal(t, "t1", tasks[2].Id)
	assert.Equal(t, "t3", tasks[3].Id)

	// Sort by display status, asc
	opts = GetTasksByVersionOptions{
		Sorts: []TasksSortOrder{
			{Key: DisplayStatusKey, Order: 1},
		},
	}
	tasks, count, err = GetTasksByVersion(ctx, "v1", opts)
	assert.NoError(t, err)
	assert.Equal(t, 4, count)
	assert.Equal(t, "t2", tasks[0].Id)
	assert.Equal(t, "t4", tasks[1].Id)
	assert.Equal(t, "t3", tasks[2].Id)
	assert.Equal(t, "t1", tasks[3].Id)

	// Sort by base task status, asc
	opts = GetTasksByVersionOptions{
		Sorts: []TasksSortOrder{
			{Key: BaseTaskStatusKey, Order: 1},
		},
	}
	tasks, count, err = GetTasksByVersion(ctx, "v1", opts)
	assert.NoError(t, err)
	assert.Equal(t, 4, count)
	assert.Equal(t, "t2", tasks[0].Id)
	assert.Equal(t, "t1", tasks[1].Id)
	assert.Equal(t, "t3", tasks[2].Id)
	assert.Equal(t, "t4", tasks[3].Id)

	// Sort by duration, asc
	opts = GetTasksByVersionOptions{
		Sorts: []TasksSortOrder{
			{Key: TimeTakenKey, Order: 1},
		},
	}
	tasks, count, err = GetTasksByVersion(ctx, "v1", opts)
	assert.NoError(t, err)
	assert.Equal(t, 4, count)
	assert.Equal(t, "t1", tasks[0].Id)
	assert.Equal(t, "t2", tasks[1].Id)
	assert.Equal(t, "t4", tasks[2].Id)
	assert.Equal(t, "t3", tasks[3].Id)
}

func TestGetTaskStatsByVersion(t *testing.T) {
	assert.NoError(t, db.ClearCollections(Collection))
	t1 := Task{
		Id:               "t1",
		Version:          "v1",
		Execution:        0,
		Status:           evergreen.TaskStarted,
		ExpectedDuration: time.Minute,
		StartTime:        time.Date(2009, time.November, 10, 12, 0, 0, 0, time.UTC),
	}
	t2 := Task{
		Id:               "t2",
		Version:          "v1",
		Execution:        0,
		Status:           evergreen.TaskStarted,
		ExpectedDuration: 150 * time.Minute,
		StartTime:        time.Date(2009, time.November, 10, 12, 0, 0, 0, time.UTC),
	}
	t3 := Task{
		Id:        "t3",
		Version:   "v1",
		Execution: 1,
		Status:    evergreen.TaskSucceeded,
	}
	t4 := Task{
		Id:        "t4",
		Version:   "v1",
		Execution: 1,
		Status:    evergreen.TaskFailed,
	}
	t5 := Task{
		Id:        "t5",
		Version:   "v1",
		Execution: 2,
		Status:    evergreen.TaskStatusPending,
	}
	t6 := Task{
		Id:        "t6",
		Version:   "v1",
		Execution: 2,
		Status:    evergreen.TaskFailed,
	}
	assert.NoError(t, db.InsertMany(Collection, t1, t2, t3, t4, t5, t6))
	opts := GetTasksByVersionOptions{}
	stats, err := GetTaskStatsByVersion("v1", opts)
	assert.NoError(t, err)
	assert.Equal(t, 4, len(stats.Counts))
	assert.True(t, stats.ETA.Equal(time.Date(2009, time.November, 10, 14, 30, 0, 0, time.UTC)))

	assert.NoError(t, db.ClearCollections(Collection))
	assert.NoError(t, db.InsertMany(Collection, t3, t4, t5, t6))
	stats, err = GetTaskStatsByVersion("v1", opts)
	assert.NoError(t, err)
	assert.Nil(t, stats.ETA)
}

func TestGetGroupedTaskStatsByVersion(t *testing.T) {
	assert.NoError(t, db.ClearCollections(Collection))

	t1 := Task{
		Id:                      "t1",
		Version:                 "v1",
		Execution:               0,
		Status:                  evergreen.TaskSucceeded,
		BuildVariant:            "bv1",
		BuildVariantDisplayName: "Build Variant 1",
	}
	t2 := Task{
		Id:                      "t2",
		Version:                 "v1",
		Execution:               0,
		Status:                  evergreen.TaskFailed,
		BuildVariant:            "bv1",
		BuildVariantDisplayName: "Build Variant 1",
	}
	t3 := Task{
		Id:                      "t3",
		Version:                 "v1",
		Execution:               1,
		Status:                  evergreen.TaskSucceeded,
		BuildVariant:            "bv1",
		BuildVariantDisplayName: "Build Variant 1",
	}
	t4 := Task{
		Id:                      "t4",
		Version:                 "v1",
		Execution:               1,
		Status:                  evergreen.TaskFailed,
		BuildVariant:            "bv2",
		BuildVariantDisplayName: "Build Variant 2",
	}
	t5 := Task{
		Id:                      "t5",
		Version:                 "v1",
		Execution:               2,
		Status:                  evergreen.TaskStatusPending,
		BuildVariant:            "bv2",
		BuildVariantDisplayName: "Build Variant 2",
	}
	t6 := Task{
		Id:                      "t6",
		Version:                 "v1",
		Execution:               2,
		Status:                  evergreen.TaskFailed,
		BuildVariant:            "bv2",
		BuildVariantDisplayName: "Build Variant 2",
	}
	assert.NoError(t, db.InsertMany(Collection, t1, t2, t3, t4, t5, t6))

	t.Run("Fetch GroupedTaskStats with no filters applied", func(t *testing.T) {

		opts := GetTasksByVersionOptions{}
		variants, err := GetGroupedTaskStatsByVersion("v1", opts)
		assert.NoError(t, err)
		assert.Equal(t, 2, len(variants))
		expectedValues := []*GroupedTaskStatusCount{
			{
				Variant:     "bv1",
				DisplayName: "Build Variant 1",
				StatusCounts: []*StatusCount{
					{
						Status: evergreen.TaskFailed,
						Count:  1,
					},
					{
						Status: evergreen.TaskSucceeded,
						Count:  2,
					},
				},
			},
			{
				Variant:     "bv2",
				DisplayName: "Build Variant 2",
				StatusCounts: []*StatusCount{
					{
						Status: evergreen.TaskFailed,
						Count:  2,
					},
					{
						Status: evergreen.TaskStatusPending,
						Count:  1,
					},
				},
			},
		}

		compareGroupedTaskStatusCounts(t, expectedValues, variants)
	})
	t.Run("Fetch GroupedTaskStats with filters applied", func(t *testing.T) {

		opts := GetTasksByVersionOptions{
			Variants: []string{"bv1"},
		}

		variants, err := GetGroupedTaskStatsByVersion("v1", opts)
		assert.NoError(t, err)
		assert.Equal(t, 1, len(variants))
		expectedValues := []*GroupedTaskStatusCount{
			{
				Variant:     "bv1",
				DisplayName: "Build Variant 1",
				StatusCounts: []*StatusCount{
					{
						Status: evergreen.TaskFailed,
						Count:  1,
					},
					{
						Status: evergreen.TaskSucceeded,
						Count:  2,
					},
				},
			},
		}
		compareGroupedTaskStatusCounts(t, expectedValues, variants)
	})

}

func compareGroupedTaskStatusCounts(t *testing.T, expected, actual []*GroupedTaskStatusCount) {
	// reflect.DeepEqual does not work here, it was failing because of the slice ptr values for StatusCounts.
	for i, e := range expected {
		a := actual[i]
		assert.Equal(t, e.Variant, a.Variant)
		assert.Equal(t, e.DisplayName, a.DisplayName)
		assert.Equal(t, len(e.StatusCounts), len(a.StatusCounts))
		for j, expectedCount := range e.StatusCounts {
			actualCount := a.StatusCounts[j]
			assert.Equal(t, expectedCount.Status, actualCount.Status)
			assert.Equal(t, expectedCount.Count, actualCount.Count)
		}
	}
}

func TestGetBaseStatusesForActivatedTasks(t *testing.T) {
	assert.NoError(t, db.ClearCollections(Collection))
	t1 := Task{
		Id:            "t1",
		Version:       "v1",
		Status:        evergreen.TaskStarted,
		ActivatedTime: time.Time{},
		DisplayName:   "task_1",
		BuildVariant:  "bv_1",
	}
	t2 := Task{
		Id:            "t2",
		Version:       "v1",
		Status:        evergreen.TaskSetupFailed,
		ActivatedTime: time.Time{},
		DisplayName:   "task_2",
		BuildVariant:  "bv_2",
	}
	t3 := Task{
		Id:            "t1_base",
		Version:       "v1_base",
		Status:        evergreen.TaskSucceeded,
		ActivatedTime: time.Time{},
		DisplayName:   "task_1",
		BuildVariant:  "bv_1",
	}
	t4 := Task{
		Id:            "t2_base",
		Version:       "v1_base",
		Status:        evergreen.TaskStarted,
		ActivatedTime: time.Time{},
		DisplayName:   "task_2",
		BuildVariant:  "bv_2",
	}
	t5 := Task{
		Id:            "only_on_base",
		Version:       "v1_base",
		Status:        evergreen.TaskFailed,
		ActivatedTime: time.Time{},
		DisplayName:   "only_on_base",
		BuildVariant:  "bv_2",
	}
	assert.NoError(t, db.InsertMany(Collection, t1, t2, t3, t4, t5))
	statuses, err := GetBaseStatusesForActivatedTasks("v1", "v1_base")
	assert.NoError(t, err)
	assert.Equal(t, 2, len(statuses))
	assert.Equal(t, statuses[0], evergreen.TaskStarted)
	assert.Equal(t, statuses[1], evergreen.TaskSucceeded)

	assert.NoError(t, db.ClearCollections(Collection))
	assert.NoError(t, db.InsertMany(Collection, t1, t2, t5))
	statuses, err = GetBaseStatusesForActivatedTasks("v1", "v1_base")
	assert.NoError(t, err)
	assert.Equal(t, 0, len(statuses))
}

func TestHasMatchingTasks(t *testing.T) {
	require.NoError(t, db.ClearCollections(Collection))
	t1 := Task{
		Id:           "t1",
		Version:      "v1",
		BuildVariant: "bv1",
		Execution:    0,
		Status:       evergreen.TaskSucceeded,
	}
	t2 := Task{
		Id:           "t2",
		Version:      "v1",
		BuildVariant: "bv1",
		Execution:    0,
		Status:       evergreen.TaskFailed,
	}
	// TODO: Reenable this test once https://jira.mongodb.org/browse/EVG-16918 is complete
	// bv1 := build.Build{
	// 	Id:           "bv1",
	// 	BuildVariant: "bv1",
	// 	DisplayName:  "Build Variant 1",
	// }
	t3 := Task{
		Id:           "t3",
		Version:      "v1",
		BuildVariant: "bv2",
		Execution:    1,
		Status:       evergreen.TaskSucceeded,
	}
	t4 := Task{
		Id:           "t4",
		Version:      "v1",
		BuildVariant: "bv2",
		Execution:    1,
		Status:       evergreen.TaskFailed,
	}
	// bv2 := build.Build{
	// 	Id:           "bv2",
	// 	BuildVariant: "bv2",
	// 	DisplayName:  "Build Variant 2",
	// }
	t5 := Task{
		Id:           "t5",
		Version:      "v1",
		BuildVariant: "bv3",
		Execution:    2,
		Status:       evergreen.TaskStatusPending,
	}
	t6 := Task{
		Id:           "t6",
		Version:      "v1",
		BuildVariant: "bv3",
		Execution:    2,
		Status:       evergreen.TaskFailed,
	}
	// bv3 := build.Build{
	// 	Id:           "bv3",
	// 	BuildVariant: "bv3",
	// 	DisplayName:  "Build Variant 3",
	// }
	// assert.NoError(t, db.InsertMany("build", bv1, bv2, bv3))
	assert.NoError(t, db.InsertMany(Collection, t1, t2, t3, t4, t5, t6))
	opts := HasMatchingTasksOptions{
		Statuses: []string{evergreen.TaskFailed},
	}
	ctx := context.TODO()
	hasMatchingTasks, err := HasMatchingTasks(ctx, "v1", opts)
	assert.NoError(t, err)
	assert.True(t, hasMatchingTasks)

	opts.Statuses = []string{evergreen.TaskWillRun}

	hasMatchingTasks, err = HasMatchingTasks(ctx, "v1", opts)
	assert.NoError(t, err)
	assert.False(t, hasMatchingTasks)

	opts = HasMatchingTasksOptions{
		Variants: []string{"bv1"},
	}
	hasMatchingTasks, err = HasMatchingTasks(ctx, "v1", opts)
	assert.NoError(t, err)
	assert.True(t, hasMatchingTasks)

	// TODO: Reenable this test once https://jira.mongodb.org/browse/EVG-16918 is complete
	// opts.Variants = []string{"Build Variant 2"}
	// hasMatchingTasks, err = HasMatchingTasks("v1", opts)
	// assert.NoError(t, err)
	// assert.True(t, hasMatchingTasks)

	opts.Variants = []string{"DNE"}
	hasMatchingTasks, err = HasMatchingTasks(ctx, "v1", opts)
	assert.NoError(t, err)
	assert.False(t, hasMatchingTasks)
}

func TestFindAllUnmarkedBlockedDependencies(t *testing.T) {
	assert := assert.New(t)
	require.NoError(t, db.ClearCollections(Collection))

	t1 := &Task{
		Id:     "t1",
		Status: evergreen.TaskFailed,
	}

	tasks := []Task{
		{
			Id: "t2",
			DependsOn: []Dependency{
				{
					TaskId: "t1",
					Status: evergreen.TaskSucceeded,
				},
			},
		},
		{
			Id: "t3",
			DependsOn: []Dependency{
				{
					TaskId: "t1",
					Status: evergreen.TaskFailed,
				},
			},
		},
		{
			Id: "t4",
			DependsOn: []Dependency{
				{
					TaskId:       "t1",
					Status:       evergreen.TaskSucceeded,
					Unattainable: true,
				},
			},
		},
		{
			Id: "t5",
			DependsOn: []Dependency{
				{
					TaskId: "t1",
					Status: evergreen.TaskFailed,
				},
				{
					TaskId: "t2",
					Status: evergreen.TaskSucceeded,
				},
			},
		},
	}
	for _, task := range tasks {
		assert.NoError(task.Insert())
	}

	deps, err := t1.FindAllUnmarkedBlockedDependencies()
	assert.NoError(err)
	assert.Len(deps, 1)
}

func TestFindAllMarkedUnattainableDependencies(t *testing.T) {
	assert := assert.New(t)
	require.NoError(t, db.ClearCollections(Collection))

	t1 := &Task{Id: "t1"}
	tasks := []Task{
		{
			Id: "t2",
			DependsOn: []Dependency{
				{
					TaskId:       "t1",
					Unattainable: true,
				},
			},
		},
		{
			Id: "t3",
			DependsOn: []Dependency{
				{
					TaskId: "t1",
				},
				{
					TaskId:       "t2",
					Unattainable: true,
				},
			},
		},
	}

	for _, task := range tasks {
		assert.NoError(task.Insert())
	}

	unattainableTasks, err := t1.FindAllMarkedUnattainableDependencies()
	assert.NoError(err)
	assert.Len(unattainableTasks, 1)
}

func TestCountNumExecutionsForInterval(t *testing.T) {
	assert.NoError(t, db.ClearCollections(Collection, OldCollection))

	now := time.Now()
	earlier := time.Now().Add(-time.Hour)
	reallyEarly := now.Add(-12 * time.Hour)
	tasks := []Task{
		{
			Id:           "notFinished",
			Project:      "myProject",
			Status:       evergreen.TaskStarted,
			Requester:    evergreen.RepotrackerVersionRequester,
			BuildVariant: "bv1",
			DisplayName:  "task1",
			Execution:    1,
		},
		{
			Id:           "finished",
			Project:      "myProject",
			Status:       evergreen.TaskFailed,
			Requester:    evergreen.RepotrackerVersionRequester,
			BuildVariant: "bv1",
			DisplayName:  "task1",
			FinishTime:   now,
			Execution:    1,
		},
		{
			Id:           "finishedEarlier",
			Project:      "myProject",
			Status:       evergreen.TaskFailed,
			Requester:    evergreen.RepotrackerVersionRequester,
			BuildVariant: "bv1",
			DisplayName:  "task1",
			FinishTime:   earlier,
			Execution:    1,
		},
		{
			Id:           "patch",
			Project:      "myProject",
			Status:       evergreen.TaskSucceeded,
			Requester:    evergreen.PatchVersionRequester,
			BuildVariant: "bv1",
			DisplayName:  "task1",
			FinishTime:   now,
			Execution:    1,
		},
		{
			Id:           "tooEarly",
			Project:      "myProject",
			Status:       evergreen.TaskSucceeded,
			Requester:    evergreen.RepotrackerVersionRequester,
			BuildVariant: "bv1",
			DisplayName:  "task1",
			FinishTime:   reallyEarly,
			Execution:    1,
		},
		{
			Id:           "wrongTask",
			Project:      "myProject",
			Status:       evergreen.TaskFailed,
			Requester:    evergreen.RepotrackerVersionRequester,
			BuildVariant: "bv1",
			DisplayName:  "task2",
			FinishTime:   now,
			Execution:    1,
		},
		{
			Id:           "wrongVariant",
			Project:      "myProject",
			Status:       evergreen.TaskFailed,
			Requester:    evergreen.RepotrackerVersionRequester,
			BuildVariant: "bv2",
			DisplayName:  "task1",
			FinishTime:   now,
			Execution:    1,
		},
	}
	for _, each := range tasks {
		assert.NoError(t, each.Insert())
		each.Execution = 0
		// Duplicate everything for the old task collection to ensure this is working.
		assert.NoError(t, db.Insert(OldCollection, each))
	}

	for testName, test := range map[string]func(*testing.T){
		"nothingInRange": func(t *testing.T) {
			input := NumExecutionsForIntervalInput{
				ProjectId:    "myProject",
				BuildVarName: "bv1",
				TaskName:     "task1",
				StartTime:    time.Now().Add(-20 * time.Hour),
				EndTime:      time.Now().Add(-18 * time.Hour),
			}
			numExecutions, err := CountNumExecutionsForInterval(input)
			assert.NoError(t, err)
			assert.Equal(t, 0, numExecutions)
		},
		"lotsInRange": func(t *testing.T) {
			input := NumExecutionsForIntervalInput{
				ProjectId:    "myProject",
				BuildVarName: "bv1",
				TaskName:     "task1",
				StartTime:    now.Add(-20 * time.Hour),
			}
			// Should include the finished tasks in both new and old.
			numExecutions, err := CountNumExecutionsForInterval(input)
			assert.NoError(t, err)
			assert.Equal(t, 6, numExecutions)
		},
		"lessInRange": func(t *testing.T) {
			input := NumExecutionsForIntervalInput{
				ProjectId:    "myProject",
				BuildVarName: "bv1",
				TaskName:     "task1",
				StartTime:    now.Add(-2 * time.Hour),
			}
			// Should include the finished tasks in both new and old except reallyEarly.
			numExecutions, err := CountNumExecutionsForInterval(input)
			assert.NoError(t, err)
			assert.Equal(t, 4, numExecutions)
		},
		"onlyPatches": func(t *testing.T) {
			input := NumExecutionsForIntervalInput{
				ProjectId:    "myProject",
				BuildVarName: "bv1",
				TaskName:     "task1",
				Requesters:   evergreen.PatchRequesters,
				StartTime:    now.Add(-2 * time.Hour),
			}
			// Should include the patch requester.
			numExecutions, err := CountNumExecutionsForInterval(input)
			assert.NoError(t, err)
			assert.Equal(t, 2, numExecutions)
		},
	} {
		t.Run(testName, test)
	}
}

func TestHasActivatedDependentTasks(t *testing.T) {
	assert.NoError(t, db.Clear(Collection))
	t1 := Task{
		Id:        "activeDependent",
		Activated: true,
		DependsOn: []Dependency{
			{TaskId: "current"},
		},
	}
	t2 := Task{
		Id: "inactiveDependent",
		DependsOn: []Dependency{
			{TaskId: "inactive"},
		},
	}
	t3 := Task{
		Id:        "manyDependencies",
		Activated: true,
		DependsOn: []Dependency{
			{TaskId: "current"},
			{TaskId: "secondTask"},
		},
	}
	assert.NoError(t, db.InsertMany(Collection, t1, t2, t3))

	hasDependentTasks, err := HasActivatedDependentTasks("current")
	assert.NoError(t, err)
	assert.True(t, hasDependentTasks)

	hasDependentTasks, err = HasActivatedDependentTasks("secondTask")
	assert.NoError(t, err)
	assert.True(t, hasDependentTasks)

	// Tasks overriding dependencies don't count as dependent.
	assert.NoError(t, t3.SetOverrideDependencies("me"))
	hasDependentTasks, err = HasActivatedDependentTasks("secondTask")
	assert.NoError(t, err)
	assert.False(t, hasDependentTasks)

	hasDependentTasks, err = HasActivatedDependentTasks("inactive")
	assert.NoError(t, err)
	assert.False(t, hasDependentTasks)

}
