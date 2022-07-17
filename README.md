# myq-webhooks
This repository provides a way to control myq devices via an AWS Lambda Function fronted by AWS API Gateway. 

The API has a single endpoint - POST /devices with the following body:
```json
{
    "serial_number": "The serial number of the device you wish to control",
    "desired_state": "The desired state of the device",
    "action": "Action to take to get to desired state"
}
````

The function will retry the action up to 3 times with 30 second intervals but these are configurable via the code, 
