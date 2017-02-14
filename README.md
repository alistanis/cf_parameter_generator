# cf_parameter_generator
Do you need to write a ton of CloudFormation instead of Terraform for some reason? Is it annoying to have to type ParameterKey and ParameterValue a million times and copying the Parameter names into a parameters file? Well, now you don't have to.

## Examples

Printing output to command line:

    $ cf_parameter_generator -f /Users/cmc/some_cf_file.json
    [
    	{
    		"ParameterKey": "AWSaccount",
    		"ParameterValue": ""
    	},
    	{
    		"ParameterKey": "AmazonEC2FullAccessARN",
    		"ParameterValue": ""
    	},
    	{
    		"ParameterKey": "AmazonRoute53FullAccessARN",
    		"ParameterValue": ""
    	},
    	{
    		"ParameterKey": "AppTierFleetSize",
    		"ParameterValue": ""
    	},
    	{
    		"ParameterKey": "AppTierInstanceType",
    		"ParameterValue": ""
    	},
    	{
    		"ParameterKey": "ApplicationName",
    		"ParameterValue": ""
    	}
    ]

Saving output to a new file (will overwrite existing)

    $ cf_parameter_generator -f local_file.json -o /Path/to/output.json
    	