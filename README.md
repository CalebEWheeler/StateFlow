# StateFlow
Readme is a work in progress. Definition of done includes fully documented architecture via flowchart diagram with text, table of contents for package structure, local testing instruction, API documentation links, deployment charts, observability, etc. 

Order fulfillment system - design it so that it's a workflow orchestrator engine that can have any use case "plugged in"

"state machine that advances one job at a time..."

flow:
1. POST /create/order
1. Workers poll jobs
1. Jobs have steps
1. Jobs Transition state
1. Retries handled by engine

Jobs:
- create_order
- reserve_inventory
- create_shipment
- send_confirmation

(BONUS)
- retry_failed_job
- cancel_order
- reconcile_order

ENTRY POINT
POST /order - creates workflow
