# StateFlow

StateFlow is a workflow orchestration engine written in Go that executes long-running business processes through a series of durable, state-driven jobs.

The project is designed around a simple principle:

A workflow advances one job at a time until completion.

While the initial implementation focuses on order fulfillment, the orchestration engine is intended to support any workflow that can be represented as a sequence of state transitions.

## Current Status

🚧 StateFlow is currently under active development.

## Planned documentation and project deliverables include:

* Architecture diagrams and workflow visualizations
* Package and project structure documentation
* Local development and testing instructions
* API documentation
* Deployment manifests and charts
* Observability and monitoring examples
* Retry and recovery strategy documentation

### Example Workflow

The current reference implementation models an order fulfillment process.

## Entry Point

`POST /order`

Creates a new workflow and enqueues the initial job.

## Workflow Steps

```
create_order
    ↓
reserve_inventory
    ↓
create_shipment
    ↓
send_confirmation
```

Each workflow step is represented by a durable job stored in the database and processed asynchronously by workers.

## Worker Engine

Workers continuously poll for pending jobs and process them according to their current step.

Responsibilities include:

* Claiming pending jobs
* Executing workflow steps
* Advancing workflow state
* Creating subsequent jobs
* Handling failures
* Supporting retry strategies

## Current Job Types

Core Workflow

* create_order
* reserve_inventory
* create_shipment
* send_confirmation

Planned Extensions

* retry_failed_job
* cancel_order
* reconcile_order

## Design Goals

* Durable workflow execution
* State-driven orchestration
* Retry and failure recovery
* Pluggable workflow implementations
* Horizontal worker scalability
* Database-backed job persistence
* Observability and operational visibility

### High-Level Flow

```
POST /order
    ↓
Create Workflow
    ↓
Create Initial Job
    ↓
Worker Claims Job
    ↓
Execute Step
    ↓
Create Next Job
    ↓
Workflow Complete
```