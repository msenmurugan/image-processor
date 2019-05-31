/*
Copyright (C) 2019 Synopsys, Inc.
Licensed to the Apache Software Foundation (ASF) under one
or more contributor license agreements. See the NOTICE file
distributed with this work for additional information
regarding copyright ownership. The ASF licenses this file
to you under the Apache License, Version 2.0 (the
"License"); you may not use this file except in compliance
with the License. You may obtain a copy of the License at
http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing,
software distributed under the License is distributed on an
"AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
KIND, either express or implied. See the License for the
specific language governing permissions and limitations
under the License.
*/

package main

import (
	"time"

	docker "github.com/fsouza/go-dockerclient"
	log "github.com/sirupsen/logrus"
)

const workerTimeout = 60 * time.Second

// Handler defines the Docker handler interface
type Handler interface {
	Handle(*docker.APIEvents) error
}

// EventRouter defines the Docker handler event routes
type EventRouter struct {
	handlers      map[string][]Handler
	dockerClient  *docker.Client
	listener      chan *docker.APIEvents
	workers       chan *worker
	workerTimeout time.Duration
}

// NewEventRouter returns the Docker handler event routes configuration
func NewEventRouter(bufferSize int, workerPoolSize int, dockerClient *docker.Client, handlers map[string][]Handler) (*EventRouter, error) {
	workers := make(chan *worker, workerPoolSize)
	for i := 0; i < workerPoolSize; i++ {
		workers <- &worker{}
	}

	eventRouter := &EventRouter{
		handlers:      handlers,
		dockerClient:  dockerClient,
		listener:      make(chan *docker.APIEvents, bufferSize),
		workers:       workers,
		workerTimeout: workerTimeout,
	}

	return eventRouter, nil
}

// Start listens for an events
func (e *EventRouter) Start() error {
	log.Info("Starting event router.")
	go e.routeEvents()
	if err := e.dockerClient.AddEventListener(e.listener); err != nil {
		return err
	}
	return nil
}

// Stop removes the listener for an event
func (e *EventRouter) Stop() error {
	if e.listener == nil {
		return nil
	}
	if err := e.dockerClient.RemoveEventListener(e.listener); err != nil {
		return err
	}
	return nil
}

// routeEvents routes the incoming events
func (e *EventRouter) routeEvents() {
	for {
		event := <-e.listener
		timer := time.NewTimer(e.workerTimeout)
		gotWorker := false
		for !gotWorker {
			select {
			case w := <-e.workers:
				go w.doWork(event, e)
				gotWorker = true
			case <-timer.C:
				log.Infof("Timed out waiting for worker. Re-initializing wait.")
			}
		}
	}
}

type worker struct{}

// doWork perform for the incoming event
func (w *worker) doWork(event *docker.APIEvents, e *EventRouter) {
	defer func() { e.workers <- w }()
	if event == nil {
		return
	}
	if handlers, ok := e.handlers[event.Status]; ok {
		log.Debugf("Processing event: %#v", event)
		for _, handler := range handlers {
			if err := handler.Handle(event); err != nil {
				log.Errorf("Error processing event %#v. Error: %v", event, err)
			}
		}
	}
}
