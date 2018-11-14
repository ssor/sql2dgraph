package main

import (
    "fmt"
    "github.com/ssor/sql2graphql/cmd/chain/schedule"
    "github.com/ssor/sql2graphql/cmd/chain/server"
    "os"
    "os/signal"
)

func main() {
    dataSource := server.NewSimulateDataSource("abc123")
    dataSource.Run()
    updateServer := server.CreateLegerUpdateServer(dataSource)
    updateServer.Run()

    scheduler := schedule.NewScheduler(updateServer.AddBlockTask)
    scheduler.Run()
    dataSource.SubscribeBlockChangedEvent(scheduler.EventHandler)

    waitForStopSignal()
}

func waitForStopSignal() {
    c := make(chan os.Signal, 1)
    signal.Notify(c, os.Interrupt)

    // Block until a signal is received.
    <-c
    fmt.Println("stopped")
}
