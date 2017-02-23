# cf_parameter_generator
Do you need to write a ton of CloudFormation instead of Terraform for some reason? (I may have an obvious bias in favor of Terraform) Is it annoying to have to type ParameterKey and ParameterValue a million times and copying the Parameter names into a parameters file? Well, now you don't have to.

Get it: `go get -u github.com/alistanis/cf_parameter_generator ./...`

## Examples

Usage

    cf_parameter_generator --help
    Usage of cf_parameter_generator:
      -f string
        	The file to read from to generate parameters.
      -inyaml
        	Will expect input as yaml instead of json.
      -min
        	If given, will generate minified output.
      -o string
        	Optional: Specify a file name to write out parameters.
      -outyaml
        	Will output in yaml instead of json.
      -overwrite
        	By default, will update an existing parameters file with newly found parameters, but will not overwrite.
      -r	Removes old entries from parameters found in old parameters files.
      -spaces int
        	The number of spaces used to indent the file if not generating minified output. (default 2)

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

Saving output to a new file (will update an existing file or overwrite it if it is blank (0 bytes))

    $ cf_parameter_generator -f local_file.json -o /Path/to/output.json
    	